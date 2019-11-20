package filproofs

import "bytes"
import "errors"
import "fmt"
import "math"
import "math/rand"
import big "math/big"
import "encoding/binary"

import util "github.com/filecoin-project/specs/util"
import file "github.com/filecoin-project/specs/systems/filecoin_files/file"
import piece "github.com/filecoin-project/specs/systems/filecoin_files/piece"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
import sectorIndex "github.com/filecoin-project/specs/systems/filecoin_mining/sector_index"

type SHA256Hash Bytes32
type PedersenHash Bytes32
type Bytes32 []byte
type UInt = util.UInt
type PieceInfo = *sector.PieceInfo_I
type Label Bytes32
type Commitment = sector.Commitment

func WinSDRParams(cfg SDRCfg) *WindowedStackedDRG_I {
	inner := SDRParams(cfg)
	innerSealCfg := cfg.SealCfg()

	outerSealCfg := sector.SealCfg_I{
		SubsectorCount_: 1,
		SectorSize_:     innerSealCfg.SectorSize(),
		Partitions_:     innerSealCfg.Partitions(),
	}

	outerCfg := SDRCfg_I{
		SealCfg_:         &outerSealCfg,
		ElectionPoStCfg_: cfg.ElectionPoStCfg(),
		SurprisePoStCfg_: cfg.SurprisePoStCfg(),
	}

	outer := SDRParams(&outerCfg)
	return &WindowedStackedDRG_I{
		Inner_: inner,
		Outer_: outer,
	}
}

func SDRParams(cfg SDRCfg) *StackedDRG_I {
	// TODO: Bridge constants with orient model.
	const LAYERS = 10
	const NODE_SIZE = 32
	const OFFLINE_CHALLENGES = 6666
	const FEISTEL_ROUNDS = 3
	var FEISTEL_KEYS = [FEISTEL_ROUNDS]UInt{1, 2, 3}
	var FIELD_MODULUS = new(big.Int)
	// https://github.com/zkcrypto/pairing/blob/master/src/bls12_381/fr.rs#L4
	FIELD_MODULUS.SetString("52435875175126190479447740508185965837690552500527637822603658699938581184513", 10)

	nodes := UInt(cfg.SealCfg().SectorSize() / NODE_SIZE)

	return &StackedDRG_I{
		Layers_:     StackedDRGLayers(LAYERS),
		Challenges_: StackedDRGChallenges(OFFLINE_CHALLENGES),
		NodeSize_:   StackedDRGNodeSize(NODE_SIZE),
		Nodes_:      StackedDRGNodes(nodes),
		Algorithm_:  &StackedDRG_Algorithm_I{},
		DRGCfg_: &DRGCfg_I{
			Algorithm_: &DRGCfg_Algorithm_I{
				ParentsAlgorithm_: DRGCfg_Algorithm_ParentsAlgorithm_DRSample,
				RNGAlgorithm_:     DRGCfg_Algorithm_RNGAlgorithm_ChaCha20,
			},
			Degree_: 6,
			Nodes_:  DRGNodeCount(nodes),
		},
		ExpanderGraphCfg_: &ExpanderGraphCfg_I{
			Algorithm_: ExpanderGraphCfg_Algorithm_Make_ChungExpanderAlgorithm(
				&ChungExpanderAlgorithm_I{
					PermutationAlgorithm_: ChungExpanderAlgorithm_PermutationAlgorithm_Make_Feistel(&Feistel_I{
						Keys_:         FEISTEL_KEYS[:],
						Rounds_:       FEISTEL_ROUNDS,
						HashFunction_: ChungExpanderPermutationFeistelHashFunction_SHA256,
					}),
				}),
			Degree_: 8,
			Nodes_:  ExpanderGraphNodeCount(nodes),
		},
		Curve_: &EllipticCurve_I{
			FieldModulus_: *FIELD_MODULUS,
		},
		Cfg_: cfg,
	}
}

func (sdr *StackedDRG_I) drg() *DRG_I {
	return &DRG_I{
		Config_: sdr.DRGCfg(),
	}
}

func (sdr *StackedDRG_I) expander() *ExpanderGraph_I {
	return &ExpanderGraph_I{
		Config_: sdr.ExpanderGraphCfg(),
	}
}

func (drg *DRG_I) Parents(node UInt) []UInt {
	config := drg.Config()
	degree := UInt(config.Degree())
	return DRGAlgorithmComputeParents(config.Algorithm().ParentsAlgorithm(), degree, node)
}

// TODO: Verify this. Both the port from impl and the algorithm.
func DRGAlgorithmComputeParents(alg DRGCfg_Algorithm_ParentsAlgorithm, degree UInt, node UInt) (parents []UInt) {
	switch alg {
	case DRGCfg_Algorithm_ParentsAlgorithm_DRSample:
		util.Assert(node > 0)
		parents = append(parents, node-1)

		m := degree - 1

		var k UInt
		for k = 0; k < m; k++ {
			logi := int(math.Floor(math.Log2(float64(node * m))))
			// FIXME: Make RNG parameterizable and specify it.
			j := rand.Intn(logi)
			jj := math.Min(float64(node*m+k), float64(UInt(1)<<uint(j+1)))
			backDist := randInRange(int(math.Max(float64(UInt(jj)>>1), 2)), int(jj+1))
			out := (node*m + k - backDist) / m

			parents = append(parents, out)
		}

		return parents

	default:
		panic(fmt.Sprintf("DRG algorithm not supported: %v", alg))
	}
}

func randInRange(lowInclusive int, highExclusive int) UInt {
	// NOTE: current implementation uses a more sophisticated method for repeated sampling within a range.
	// We need to converge on and fully specify the actual method, since this must be deterministic.
	return UInt(rand.Intn(highExclusive-lowInclusive) + lowInclusive)
}

func (exp *ExpanderGraph_I) Parents(node UInt) []UInt {
	d := exp.Config().Degree()

	// TODO: How do we handle choice of algorithm generically?
	return exp.Config().Algorithm().As_ChungExpanderAlgorithm().Parents(node, d, exp.Config().Nodes())
}

func (chung *ChungExpanderAlgorithm_I) Parents(node UInt, d ExpanderGraphDegree, nodes ExpanderGraphNodeCount) []UInt {
	var parents []UInt
	var i UInt
	for i = 0; i < UInt(d); i++ {
		parent := chung.ithParent(node, i, d, nodes)
		parents = append(parents, parent)
	}
	return parents
}

func (chung *ChungExpanderAlgorithm_I) ithParent(node UInt, i UInt, degree ExpanderGraphDegree, nodes ExpanderGraphNodeCount) UInt {
	// ithParent generates one of d parents of node.
	d := UInt(degree)

	// This is done by operating on permutations of a set with d elements per node.
	setSize := UInt(nodes) * d

	// There are d ways of mapping each node into the set, and we choose the ith.
	// Note that we can project the element back to the original node: element / d == node.
	element := node*d + i

	// Permutations of the d elements corresponding to each node yield d new elements,
	permuted := chung.PermutationAlgorithm().As_Feistel().Permute(setSize, element)

	// each of which can be projected back to a node.
	projected := permuted / d

	// We have selected the ith such parent of node.
	return projected
}

func (f *Feistel_I) Permute(size UInt, i UInt) UInt {
	// Call into feisel.go.
	panic("TODO")
}

func (sdr *StackedDRG_I) Seal(sid sector.SectorID, subsectorData [][]byte, randomness sector.SealRandomness) SealSetupArtifacts {
	subsectorCount := int(sdr.SubsectorCount())
	nodeSize := int(sdr.NodeSize())
	nodes := int(sdr.Nodes())
	subsectorSize := nodes * nodeSize

	util.Assert(len(subsectorData) == subsectorCount)

	var allKeyLayers [][]byte
	var replicas [][]byte
	var sealSeeds []sector.SealSeed
	var subsectorCommDs []sector.UnsealedSectorCID
	var subsectorCommDTreePaths []file.Path

	for i, data := range subsectorData {
		_ = data
		keyLayers, replica, sealSeed, subsectorCommD, subsectorCommDTreePath := sdr.SealSubsector(i, sid, data, randomness)

		allKeyLayers = append(allKeyLayers, keyLayers...)
		replicas = append(replicas, replica)
		sealSeeds = append(sealSeeds, sealSeed)
		subsectorCommDs = append(subsectorCommDs, subsectorCommD)
		subsectorCommDTreePaths = append(subsectorCommDTreePaths, subsectorCommDTreePath)
	}

	var replicaColumns [][]Label
	for _, replica := range replicas {
		for i := 0; i < subsectorSize; i += nodeSize {
			replicaColumns[i] = append(replicaColumns[i], Label(replica[i:i+nodeSize]))
		}
	}
	var columnReplica []byte
	for i, column := range replicaColumns {
		copy(columnReplica[i*nodeSize:(i+1)*nodeSize], []byte(hashColumn(column)))
	}

	commC, commRLast, commR, commCTreePath, commRLastTreePath := sdr.GenerateCommitments(columnReplica, allKeyLayers)

	result := SealSetupArtifacts_I{
		CommD_:             Commitment(commC),
		CommR_:             SealedSectorCID(commR),
		CommC_:             Commitment(commC),
		CommRLast_:         Commitment(commRLast),
		CommDTreePaths_:    subsectorCommDTreePaths,
		CommCTreePath_:     commCTreePath,
		CommRLastTreePath_: commRLastTreePath,
		Seeds_:             sealSeeds,
		KeyLayers_:         allKeyLayers,
		Replicas_:          replicas,
	}
	return &result
}

func (sdr *StackedDRG_I) SealSubsector(subsectorIndex int, sid sector.SectorID, data []byte, randomness sector.SealRandomness) ([][]byte, []byte, sector.SealSeed, sector.UnsealedSectorCID, file.Path) {
	commD, commDTreePath := ComputeDataCommitment(data)
	sealSeed := computeSealSeed(sid, subsectorIndex, randomness, commD)
	nodeSize := int(sdr.NodeSize())
	nodes := len(data) / nodeSize
	curveModulus := sdr.Curve().FieldModulus()
	layers := int(sdr.Layers())

	keyLayers := generateSDRKeyLayers(sdr.drg(), sdr.expander(), sealSeed, nodes, layers, nodeSize, curveModulus)
	key := keyLayers[len(keyLayers)-1]
	replica := encodeData(data, key, nodeSize, &curveModulus)

	return keyLayers, replica, sealSeed, commD, commDTreePath
}

func (sdr *StackedDRG_I) GenerateCommitments(replica []byte, keyLayers [][]byte) (commC PedersenHash, commRLast PedersenHash, commR PedersenHash, commCTreePath file.Path, commRLastTreePath file.Path) {
	commRLast, commRLastTreePath = BuildTree_PedersenHash(replica)
	commC, commCTreePath = computeCommC(keyLayers, int(sdr.NodeSize()))
	commR = BinaryHash_PedersenHash(commC, commRLast)

	return commC, commRLast, commR, commCTreePath, commRLastTreePath
}

func computeSealSeed(sid sector.SectorID, subsectorIndex int, randomness sector.SealRandomness, commD sector.UnsealedSectorCID) sector.SealSeed {
	var proverId []byte // TODO: Derive this from sid.MinerID()
	sectorNumber := sid.Number()

	var preimage []byte
	preimage = append(preimage, proverId...)
	preimage = append(preimage, littleEndianBytesFromUInt(UInt(sectorNumber), 8)...)
	preimage = append(preimage, littleEndianBytesFromInt(subsectorIndex, 8)...)
	preimage = append(preimage, randomness...)
	preimage = append(preimage, Commitment_UnsealedSectorCID(commD)...)

	// FIXME: Implement
	return sector.SealSeed{}
}

func generateSDRKeyLayers(drg *DRG_I, expander *ExpanderGraph_I, sealSeed sector.SealSeed, nodes int, layers int, nodeSize int, modulus big.Int) [][]byte {
	var keyLayers [][]byte
	var prevLayer []byte

	for i := 0; i <= layers; i++ {
		currentLayer := labelLayer(drg, expander, sealSeed, nodes, nodeSize, prevLayer)
		keyLayers = append(keyLayers, currentLayer)
		prevLayer = currentLayer
	}

	return keyLayers
}

func labelLayer(drg *DRG_I, expander *ExpanderGraph_I, sealSeed sector.SealSeed, nodeSize int, nodes int, prevLayer []byte) []byte {
	size := nodes * nodeSize
	labels := make([]byte, size)

	for i := 0; i < nodes; i++ {
		var parents []Label

		// The first node of every layer has no DRG Parents.
		if i > 0 {
			for parent := range drg.Parents(UInt(i)) {
				start := parent * nodeSize
				parents = append(parents, labels[start:start+nodeSize])
			}
		}

		// The first layer has no expander parents.
		if prevLayer != nil {
			for parent := range expander.Parents(UInt(i)) {
				start := parent * nodeSize
				parents = append(parents, labels[start:start+nodeSize])
			}
		}

		label := generateLabel(sealSeed, i, parents)
		labels = append(labels, label...)
	}

	return labels
}

func encodeData(data []byte, key []byte, nodeSize int, modulus *big.Int) []byte {
	if len(data) != len(key) {
		panic("Key and data must be same length.")
	}

	encoded := make([]byte, len(data))
	for i := 0; i < len(data); i += nodeSize {
		copy(encoded[i:i+nodeSize], encodeNode(data[i:i+nodeSize], key[i:i+nodeSize], modulus, nodeSize))
	}

	return encoded
}

func generateLabel(sealSeed sector.SealSeed, node int, dependencies []Label) []byte {
	nodeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nodeBytes, uint64(node))

	preimage := append(sealSeed, nodeBytes...)
	for _, dependency := range dependencies {
		preimage = append(preimage, dependency...)
	}

	return deriveLabel(preimage)
}

func deriveLabel(elements []byte) []byte {
	return HashBytes_SHA256Hash(elements)
}

func computeCommC(keyLayers [][]byte, nodeSize int) (PedersenHash, file.Path) {
	leaves := make([]byte, len(keyLayers[0]))

	// For each node in the graph,
	for start := 0; start < len(leaves); start += nodeSize {
		end := start + nodeSize

		var column []Label
		// Concatenate that node's label at each layer, in order, into a column.
		for i := 0; i < len(keyLayers); i++ {
			label := keyLayers[i][start:end]
			column = append(column, label)
		}

		// And hash that column to create the leaf of a new tree.
		hashed := hashColumn(column)
		copy(leaves[start:end], hashed[:])
	}

	// Return the root of and path to the column tree.
	return BuildTree_PedersenHash(leaves)
}

func hashColumn(column []Label) PedersenHash {
	var preimage []byte
	for _, label := range column {
		preimage = append(preimage, label...)
	}
	return HashBytes_PedersenHash(preimage)
}

type PrivateOfflineSDRProof []OfflineSDRChallengeProof

func (sdr *StackedDRG_I) CreateSealProof(challengeSeed sector.InteractiveSealRandomness, aux sector.ProofAuxTmp) sector.SealProof {
	privateProof := sdr.CreatePrivateSealProof(challengeSeed, aux)

	// Sanity check: newly-created proofs must pass verification.
	util.Assert(sdr.VerifyPrivateProof(privateProof, aux.Seeds(), challengeSeed, aux.CommD(), aux.CommR()))

	return sdr.CreateOfflineCircuitProof(privateProof, aux)
}

func (sdr *StackedDRG_I) CreatePrivateSealProof(randomness sector.InteractiveSealRandomness, aux sector.ProofAuxTmp) (challengeProofs PrivateOfflineSDRProof) {
	sealSeeds := aux.Seeds()
	nodeSize := UInt(sdr.NodeSize())
	challenges := sdr.GenerateOfflineChallenges(sealSeeds, randomness, sdr.Challenges())

	columnTree := LoadMerkleTree(aux.CommCTreePath())
	for c := range challenges {
		challengeProof := CreateChallengeProof(sdr.drg(), sdr.expander(), sealSeeds, UInt(c), nodeSize, columnTree, aux)
		challengeProofs = append(challengeProofs, challengeProof)
	}

	privateProof := challengeProofs

	return privateProof
}

// Verify a private proof.
// NOTE: Verification of a private proof is exactly the computation we will prove we have performed in a zk-SNARK.
// If we can verifiably prove that we have performed the verification of a private proof, then we need not reveal the proof itself.
// Since the zk-SNARK circuit proof is much smaller than the private proof, this allows us to save space on the chain (at the cost of increased computation to generate the zk-SNARK proof).
func (sdr *StackedDRG_I) VerifyPrivateProof(privateProof []OfflineSDRChallengeProof, sealSeeds []sector.SealSeed, randomness sector.InteractiveSealRandomness, commD Commitment, commR sector.SealedSectorCID) bool {
	subsectorCount := int(sdr.SubsectorCount())
	layers := int(sdr.Layers())
	curveModulus := sdr.Curve().FieldModulus()
	challenges := sdr.GenerateOfflineChallenges(sealSeeds, randomness, sdr.Challenges())

	// commC and commRLast must be the same for all challenge proofs, so we can arbitrarily verify against the first.
	firstChallengeProof := privateProof[0]
	commRLast := firstChallengeProof.ReplicaProof.InclusionProof.root()
	commC := firstChallengeProof.ColumnProofs[0].InclusionProof.root()

	for i, challenge := range challenges {
		// Verify one OfflineSDRChallengeProof.
		challengeProof := privateProof[i]
		dataProofs := challengeProof.DataProofs
		columnProofs := challengeProof.ColumnProofs
		replicaProof := challengeProof.ReplicaProof

		// Verify column proofs and that they represent the right columns.
		columnElements := getColumnElements(sdr.drg(), sdr.expander(), challenge)

		for i, columnElement := range columnElements {
			columnProof := columnProofs[i]

			// The provided column proofs must correspond to the expected columns.
			if !columnProof.verify(commC, UInt(columnElement)) {
				return false
			}
		}

		for layer := 0; layer < int(sdr.Layers()); layer++ {
			var sealSeed sector.SealSeed // FIXME: Get the right seed
			var parents []Label

			for _, parentProof := range columnProofs[1:] {
				parent := parentProof.Column[layer]
				parents = append(parents, parent)
			}
			providedLabel := columnProofs[columnElements[0]].Column[layer]
			calculatedLabel := generateLabel(sealSeed, i, parents)

			if !bytes.Equal(calculatedLabel, providedLabel) {
				return false
			}
		}

		for i, dataProof := range dataProofs {
			if !dataProof.verify(commD, challenge) {
				return false
			}
			dataNode := dataProof.leaf()
			dataColumn := columnProofs[0]

			keyNode := dataColumn.Column[i*subsectorCount+layers-1]

			replicaNode := replicaProof.Column[i]

			encodedNode := encodeData(dataNode, keyNode, int(sdr.NodeSize()), &curveModulus)
			if !bytes.Equal(encodedNode, replicaNode) {
				return false
			}

		}

		if !replicaProof.verify(commRLast, challenge) {
			return false
		}
	}

	commRCalculated := BinaryHash_PedersenHash(commC, commRLast)

	if !bytes.Equal(commRCalculated, AsBytes_SealedSectorCID(commR)) {
		return false
	}

	return true
}

func CreateChallengeProof(drg *DRG_I, expander *ExpanderGraph_I, sealSeeds []sector.SealSeed, challenge UInt, nodeSize UInt, columnTree *MerkleTree, aux sector.ProofAuxTmp) (proof OfflineSDRChallengeProof) {
	columnElements := getColumnElements(drg, expander, challenge)

	var columnProofs []SDRColumnProof
	for c := range columnElements {
		columnProof := createColumnProof(UInt(c), nodeSize, columnTree, aux)
		columnProofs = append(columnProofs, columnProof)
	}

	var dataProofs []InclusionProof

	for _, treePath := range aux.CommDTreePaths() {
		dataTree := LoadMerkleTree(treePath)
		dataProof := dataTree.proveInclusion(challenge)
		dataProofs = append(dataProofs, dataProof)
	}

	replicas := aux.Replicas()
	var replicaColumn []Label
	for _, replica := range replicas {
		replicaColumn = append(replicaColumn, replica[challenge*nodeSize:(challenge+1)*nodeSize])
	}

	replicaTree := LoadMerkleTree(aux.PersistentAux().CommRLastTreePath())
	replicaProof := SDRColumnProof{
		Column:         replicaColumn,
		InclusionProof: replicaTree.proveInclusion(challenge),
	}

	proof = OfflineSDRChallengeProof{
		DataProofs:   dataProofs,
		ColumnProofs: columnProofs,
		ReplicaProof: replicaProof,
	}

	return proof
}

func getColumnElements(drg *DRG_I, expander *ExpanderGraph_I, challenge UInt) (columnElements []UInt) {
	columnElements = append(columnElements, challenge)
	columnElements = append(columnElements, drg.Parents(challenge)...)
	columnElements = append(columnElements, expander.Parents(challenge)...)

	return columnElements
}

func createColumnProof(c UInt, nodeSize UInt, columnTree *MerkleTree, aux sector.ProofAuxTmp) (columnProof SDRColumnProof) {
	layers := aux.KeyLayers()
	var column []Label

	for i := 0; i < len(layers); i++ {
		column = append(column, layers[i][c*nodeSize:(c+1)*nodeSize])
	}

	columnProof = SDRColumnProof{
		Column:         column,
		InclusionProof: columnTree.proveInclusion(c),
	}

	return columnProof
}

type OfflineSDRChallengeProof struct {
	CommRLast Commitment
	CommC     Commitment

	// TODO: these proofs need to depend on hash function.
	DataProofs   []InclusionProof // SHA256
	ColumnProofs []SDRColumnProof
	ReplicaProof SDRColumnProof // Pedersen

}

type InclusionProof struct{}

func (ip *InclusionProof) leaf() []byte {
	panic("TODO")
}

func (ip *InclusionProof) leafIndex() UInt {
	panic("TODO")
}

func (ip *InclusionProof) root() Commitment {
	panic("TODO")
}

type MerkleTree struct{}

func (mt *MerkleTree) proveInclusion(challenge UInt) InclusionProof {
	panic("TODO")
}

func LoadMerkleTree(path file.Path) *MerkleTree {
	panic("TODO")
}

func (ip *InclusionProof) verify(root []byte, challenge UInt) bool {
	panic("TODO")
}

type SDRColumnProof struct {
	Column         []Label
	InclusionProof InclusionProof
}

func (proof *SDRColumnProof) verify(root []byte, challenge UInt) bool {
	if !bytes.Equal(hashColumn(proof.Column), proof.InclusionProof.leaf()) {
		return false
	}

	if proof.InclusionProof.leafIndex() != challenge {
		return false
	}

	return proof.InclusionProof.verify(root, challenge)
}

func (sdr *StackedDRG_I) CreateOfflineCircuitProof(challengeProofs []OfflineSDRChallengeProof, aux sector.ProofAuxTmp) sector.SealProof {
	// partitions := sdr.Partitions()
	// publicInputs := GeneratePublicInputs()
	panic("TODO")
}

func (sdr *StackedDRG_I) GenerateOfflineChallenges(sealSeeds []sector.SealSeed, randomness sector.InteractiveSealRandomness, challengeCount StackedDRGChallenges) (challenges []UInt) {
	nodeSize := int(sdr.NodeSize())
	nodes := sdr.Nodes()

	challengeRangeSize := nodes - 1 // Never challenge the first node.
	challengeModulus := new(big.Int)
	challengeModulus.SetUint64(uint64(challengeRangeSize))

	count := int(challengeCount)
	for i := 0; i < count; i++ {
		var bytes []byte
		for _, sealSeed := range sealSeeds {
			bytes = append(bytes, sealSeed...)
		}
		bytes = append(bytes, randomness...)
		bytes = append(bytes, littleEndianBytesFromInt(i, nodeSize)...)

		hash := HashBytes_SHA256Hash(bytes)
		bigChallenge := bigIntFromLittleEndianBytes(hash)
		bigChallenge = bigChallenge.Mod(bigChallenge, challengeModulus)

		// Sectors nodes must be 64-bit addressable, always a safe assumption.
		challenge := bigChallenge.Uint64()
		challenge += 1 // Never challenge the first node.
		challenges = append(challenges, challenge)
	}
	return challenges
}

func encodeNode(data []byte, key []byte, modulus *big.Int, nodeSize int) []byte {
	// TODO: Make this a method of StackedDRG.
	return addEncode(data, key, modulus, nodeSize)
}

func addEncode(data []byte, key []byte, modulus *big.Int, nodeSize int) []byte {

	d := bigIntFromLittleEndianBytes(data)
	k := bigIntFromLittleEndianBytes(key)

	sum := new(big.Int).Add(d, k)
	result := new(big.Int).Mod(sum, modulus)

	return littleEndianBytesFromBigInt(result, nodeSize)
}

////////////////////////////////////////////////////////////////////////////////
// Seal Verification

func (sdr *StackedDRG_I) VerifySeal(sv sector.SealVerifyInfo) bool {
	onChain := sv.OnChain()
	sealProof := onChain.Proof()
	commR := sector.SealedSectorCID(onChain.SealedCID())
	commD := sector.UnsealedSectorCID(sv.UnsealedCID())
	sealCfg := sv.SealCfg()

	// A more sophisticated accounting of CommD for verification purposes will be required when supersectors are considered.
	util.Assert(sealCfg.SubsectorCount() == 1)

	subsectorCount := int(sealCfg.SubsectorCount())

	var sealSeeds []sector.SealSeed
	for i := 0; i < subsectorCount; i++ {
		sealSeed := computeSealSeed(sv.SectorID(), i, sv.Randomness(), commD)
		sealSeeds = append(sealSeeds, sealSeed)
	}
	challenges := sdr.GenerateOfflineChallenges(sealSeeds, sv.InteractiveRandomness(), sdr.Challenges())

	return sdr.VerifyOfflineCircuitProof(commD, commR, sealSeeds, challenges, sealProof)
}

func ComputeUnsealedSectorCIDFromPieceInfos(sectorSize UInt, pieceInfos []PieceInfo) (unsealedCID sector.UnsealedSectorCID, err error) {
	rootPieceInfo := computeRootPieceInfo(pieceInfos)
	rootSize := rootPieceInfo.Size()

	if rootSize != sectorSize {
		return unsealedCID, errors.New("Wrong sector size.")
	}

	return UnsealedSectorCID(AsBytes_PieceCID(rootPieceInfo.PieceCID())), nil
}

// commD := rootPieceInfo.CommP()

func computeRootPieceInfo(pieceInfos []PieceInfo) PieceInfo {
	// Construct root PieceInfo by (shift-reduce) parsing the constituent PieceInfo array.
	// Later pieces must always be joined with equal-sized predecessors to create a new root twice their size.
	// So if a piece is larger than the current root (top of stack), add padding until it is not.
	// If a piece is smaller than the root, let it be the new root (top of stack) until reduced to a replacement that can be joined
	// with the previous.
	var stack []PieceInfo

	shift := func(p PieceInfo) {
		stack = append(stack, p)
	}
	peek := func() PieceInfo {
		return stack[len(stack)-1]
	}
	peek2 := func() PieceInfo {
		return stack[len(stack)-2]
	}
	pop := func() PieceInfo {
		stack = stack[:len(stack)-1]
		return stack[len(stack)-1]
	}
	reduce1 := func() bool {
		if len(stack) > 1 && peek().Size() == peek2().Size() {
			right := pop()
			left := pop()
			joined := joinPieceInfos(left, right)
			shift(joined)
			return true
		}
		return false
	}
	reduce := func() {
		for reduce1() {
		}
	}
	shiftReduce := func(p PieceInfo) {
		shift(p)
		reduce()
	}

	// Prime the pump with first pieceInfo.
	shift(pieceInfos[0])

	// Consume the remainder.
	for _, pieceInfo := range pieceInfos[1:] {
		// TODO: Assert that pieceInfo.Size() is a power of 2.

		// Add padding until top of stack is large enough to receive current pieceInfo.
		for peek().Size() < pieceInfo.Size() {
			shiftReduce(zeroPadding(peek().Size()))
		}

		// Add the current piece.
		shiftReduce(pieceInfo)
	}

	// Add any necessary final padding.
	for len(stack) > 1 {
		shiftReduce(zeroPadding(peek().Size()))
	}
	util.Assert(len(stack) == 1)

	return pop()
}

func zeroPadding(size UInt) PieceInfo {
	return &sector.PieceInfo_I{
		Size_: size,
		// CommP_: FIXME: Implement.
	}
}

func joinPieceInfos(left PieceInfo, right PieceInfo) PieceInfo {
	util.Assert(left.Size() == right.Size())
	return &sector.PieceInfo_I{
		Size_:     left.Size() + right.Size(),
		PieceCID_: piece.PieceCID(BinaryHash_SHA256Hash(AsBytes_PieceCID(left.PieceCID()), AsBytes_PieceCID(right.PieceCID()))), // FIXME: make this whole function generic?
	}
}

func (sdr *StackedDRG_I) VerifyOfflineCircuitProof(commD sector.UnsealedSectorCID, commR sector.SealedSectorCID, sealSeeds []sector.SealSeed, challenges []UInt, sv sector.SealProof) bool {
	//publicInputs := GeneratePublicInputs()
	panic("TODO")
}

////////////////////////////////////////////////////////////////////////////////
// PoSt

func (sdr *StackedDRG_I) GetChallengedSectors(randomness sector.PoStRandomness, faults sector.FaultSet) (sectors []sector.SectorID, challenges []UInt) {
	panic("TODO")
}

func (sdr *StackedDRG_I) GeneratePoStCandidates(challengeSeed sector.PoStRandomness, faults sector.FaultSet, sectorStore sectorIndex.SectorStore) []sector.ElectionCandidate {
	challengedSectors, challenges := sdr.GetChallengedSectors(challengeSeed, faults)
	var proofAuxs []sector.ProofAux

	for _, sector := range challengedSectors {
		proofAux := sectorStore.GetSectorProofAux(sector)
		proofAuxs = append(proofAuxs, proofAux)
	}

	_ = challenges

	panic("TODO")
}

func (sdr *StackedDRG_I) GeneratePoStProof(privateProofs []sector.PrivatePoStProof) sector.PoStProof {
	panic("TODO")
}

////////////////////////////////////////////////////////////////////////////////
// PoSt Verification

func (sdr *StackedDRG_I) VerifyPoSt(sv sector.PoStVerifyInfo) bool {
	panic("TODO")
}

////////////////////////////////////////////////////////////////////////////////
/// Generic Hashing and Merkle Tree generation

/// Binary hash compression.
// BinaryHash<T>
func BinaryHash_T(left []byte, right []byte) util.T {
	return util.T{}
}

// BinaryHash<PedersenHash>
func BinaryHash_PedersenHash(left []byte, right []byte) PedersenHash {
	return PedersenHash{}
}

// BinaryHash<SHA256Hash>
func BinaryHash_SHA256Hash(left []byte, right []byte) SHA256Hash {
	result := SHA256Hash{}
	return trimToFr32(result)
}

////////////////////////////////////////////////////////////////////////////////

/// Digest
// HashBytes<T>
func HashBytes_T(data []byte) util.T {
	return util.T{}
}

// HashBytes<PedersenHash>
func HashBytes_PedersenHash(data []byte) PedersenHash {
	return PedersenHash{}
}

// HashBytes<SHA256Hash.
func HashBytes_SHA256Hash(data []byte) SHA256Hash {
	// Digest is truncated to 254 bits.
	result := SHA256Hash{}
	return trimToFr32(result)
}

////////////////////////////////////////////////////////////////////////////////

func DigestSize_T() int {
	panic("Unspecialized")
}

func DigestSize_PedersenHash() int {
	return 32
}

func DigestSize_SHA256Hash() int {
	return 32
}

////////////////////////////////////////////////////////////////////////////////
/// Binary Merkle-tree generation

// BuildTree<T>
func BuildTree_T(data []byte) (util.T, file.Path) {
	// Plan: define this in terms of BinaryHash_T, then copy-paste changes into T-specific specializations, for now.

	// Nodes are always the digest size so data cannot be compressed to digest for storage.
	nodeSize := DigestSize_T()

	// TODO: Fail if len(dat) is not a power of 2 and a multiple of the node size.

	rows := [][]byte{data}

	for row := []byte{}; len(row) > nodeSize; {
		for i := 0; i < len(data); i += 2 * nodeSize {
			left := data[i : i+nodeSize]
			right := data[i+nodeSize : i+2*nodeSize]

			hashed := BinaryHash_T(left, right)

			row = append(row, AsBytes_T(hashed)...)
		}
		rows = append(rows, row)
	}

	// Last row is the root
	root := rows[len(rows)-1]

	if len(root) != nodeSize {
		panic("math failed us")
	}

	var filePath file.Path // TODO: dump tree to file.
	// NOTE: merkle tree file layout is illustrative, not prescriptive.

	// TODO: Check above more carefully. It's just an untested sketch for the moment.
	return fromBytes_T(root), filePath
}

// BuildTree<PedersenHash>
func BuildTree_PedersenHash(data []byte) (PedersenHash, file.Path) {
	return PedersenHash{}, file.Path("") // FIXME
}

//  BuildTree<SHA256Hash>
func BuildTree_SHA256Hash(data []byte) (SHA256Hash, file.Path) {
	return []byte{}, file.Path("") // FIXME
}

////////////////////////////////////////////////////////////////////////////////

// Destructively trim data so most significant two bits of last byte are 0.
// This ensure data interpreted as little-endian will not exceed a field with 254-bit capacity.
// NOTE: 254 bits is the capacity of BLS12-381, but other curves with ~32-byte field elements
// may have a different capacity. (Example: BLS12-377 has a capacity of 252 bits.)
func trimToFr32(data []byte) []byte {
	util.Assert(len(data) == 32)
	data[31] &= 0x3f // 0x3f = 0b0011_1111
	return data
}

func UnsealedSectorCID(h SHA256Hash) sector.UnsealedSectorCID {
	panic("not implemented -- re-arrange bits")
}

func SealedSectorCID(h PedersenHash) sector.SealedSectorCID {
	panic("not implemented -- re-arrange bits")
}

func Commitment_UnsealedSectorCID(cid sector.UnsealedSectorCID) Commitment {
	panic("not implemented -- re-arrange bits")
}

func Commitment_SealedSectorCID(cid sector.SealedSectorCID) Commitment {
	panic("not implemented -- re-arrange bits")
}

func ComputeDataCommitment(data []byte) (sector.UnsealedSectorCID, file.Path) {
	// TODO: make hash parameterizable
	hash, path := BuildTree_SHA256Hash(data)
	return UnsealedSectorCID(hash), path
}

// Compute CommP or CommD.
func ComputeUnsealedSectorCID(data []byte) (sector.UnsealedSectorCID, file.Path) {
	// TODO: check that len(data) > minimum piece size and is a power of 2.
	hash, treePath := BuildTree_SHA256Hash(data)
	return UnsealedSectorCID(hash), treePath
}

// Utilities

func reverse(bytes []byte) {
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
}

func bigIntFromLittleEndianBytes(bytes []byte) *big.Int {
	reverse(bytes)
	return new(big.Int).SetBytes(bytes)
}

// size is number of bytes to return
func littleEndianBytesFromBigInt(z *big.Int, size int) []byte {
	bytes := z.Bytes()[0:size]
	reverse(bytes)

	return bytes
}

func littleEndianBytesFromInt(n int, size int) []byte {
	z := new(big.Int)
	z.SetInt64(int64(n))
	return littleEndianBytesFromBigInt(z, size)
}

func littleEndianBytesFromUInt(n UInt, size int) []byte {
	z := new(big.Int)
	z.SetUint64(uint64(n))
	return littleEndianBytesFromBigInt(z, size)
}

func AsBytes_T(t util.T) []byte {
	panic("Unimplemented for T")

	return []byte{}
}

func AsBytes_UnsealedSectorCID(cid sector.UnsealedSectorCID) []byte {
	panic("Unimplemented for UnsealedSectorCID")

	return []byte{}
}

func AsBytes_SealedSectorCID(CID sector.SealedSectorCID) []byte {
	panic("Unimplemented for SealedSectorCID")

	return []byte{}
}

func AsBytes_PieceCID(CID piece.PieceCID) []byte {
	panic("Unimplemented for PieceCID")

	return []byte{}
}

func fromBytes_T(_ interface{}) util.T {
	panic("Unimplemented for T")
	return util.T{}
}

func isPow2(n int) bool {
	return n != 0 && n&(n-1) == 0
}
