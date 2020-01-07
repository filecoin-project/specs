package filproofs

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/rand"

	"encoding/binary"
	big "math/big"

	file "github.com/filecoin-project/specs/systems/filecoin_files/file"
	piece "github.com/filecoin-project/specs/systems/filecoin_files/piece"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	sector_index "github.com/filecoin-project/specs/systems/filecoin_mining/sector_index"
	addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
	util "github.com/filecoin-project/specs/util"
)

type SHA256Hash Bytes32
type PedersenHash Bytes32
type Bytes32 []byte
type UInt = util.UInt
type PieceInfo = sector.PieceInfo
type Label Bytes32
type Commitment = sector.Commitment

type PrivatePostCandidateProof = sector.PrivatePoStCandidateProof

const WRAPPER_LAYER_WINDOW_INDEX = -1

const NODE_SIZE = 32
const ELECTION_POST_PARTITIONS = 1
const SURPRISE_POST_PARTITIONS = 1
const POST_LEAF_CHALLENGE_COUNT = 66
const POST_CHALLENGE_RANGE_SIZE = 1

func PoStCfg(pType sector.PoStType, sectorSize sector.SectorSize, partitions UInt) *sector.PoStCfg_I {
	nodes := UInt(sectorSize / NODE_SIZE)
	return &sector.PoStCfg_I{
		Type_:               pType,
		Nodes_:              nodes,
		Partitions_:         partitions,
		LeafChallengeCount_: POST_LEAF_CHALLENGE_COUNT,
		ChallengeRangeSize_: POST_CHALLENGE_RANGE_SIZE,
	}
}
func SurprisePoStCfg(sectorSize sector.SectorSize) *sector.PoStCfg_I {
	return PoStCfg(sector.PoStType_SurprisePoSt, sectorSize, SURPRISE_POST_PARTITIONS)
}

func ElectionPoStCfg(sectorSize sector.SectorSize) *sector.PoStCfg_I {
	return PoStCfg(sector.PoStType_ElectionPoSt, sectorSize, ELECTION_POST_PARTITIONS)
}

func ElectionPoStVerifier(cfg sector.PoStCfg) *PoStVerifier_I {
	return &PoStVerifier_I{
		ElectionPoStCfg_: cfg,
	}
}

func SurprisePoStVerifier(cfg sector.PoStCfg) *PoStVerifier_I {
	return &PoStVerifier_I{
		SurprisePoStCfg_: cfg,
	}
}

func WinSDRParams(cfg ProofsCfg) *WinStackedDRG_I {
	util.Assert(cfg.SealCfg().Algorithm() == sector.SealAlgorithm_WinStackedDRG)
	// TODO: Bridge constants with orient model.
	const LAYERS = 10
	const OFFLINE_CHALLENGES = 6666
	const OFFLINE_WINDOW_CHALLENGES = 1111
	const FEISTEL_ROUNDS = 3
	var FEISTEL_KEYS = [FEISTEL_ROUNDS]UInt{1, 2, 3}
	var FIELD_MODULUS = new(big.Int)
	// https://github.com/zkcrypto/pairing/blob/master/src/bls12_381/fr.rs#L4
	FIELD_MODULUS.SetString("52435875175126190479447740508185965837690552500527637822603658699938581184513", 10)

	nodes := UInt(cfg.SealCfg().SectorSize() / NODE_SIZE)

	return &WinStackedDRG_I{
		Layers_:           WinStackedDRGLayers(LAYERS),
		Challenges_:       WinStackedDRGChallenges(OFFLINE_CHALLENGES),
		WindowChallenges_: WinStackedDRGWindowChallenges(OFFLINE_WINDOW_CHALLENGES),
		NodeSize_:         WinStackedDRGNodeSize(NODE_SIZE),
		Nodes_:            WinStackedDRGNodes(nodes),
		Algorithm_:        &WinStackedDRG_Algorithm_I{},
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
		WindowDRGCfg_: &DRGCfg_I{
			Algorithm_: &DRGCfg_Algorithm_I{
				ParentsAlgorithm_: DRGCfg_Algorithm_ParentsAlgorithm_DRSample,
				RNGAlgorithm_:     DRGCfg_Algorithm_RNGAlgorithm_ChaCha20,
			},
			Degree_: 0,
			Nodes_:  DRGNodeCount(nodes),
		},
		WindowExpanderGraphCfg_: &ExpanderGraphCfg_I{
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

func (sdr *WinStackedDRG_I) Drg() *DRG_I {
	return &DRG_I{
		Config_: sdr.DRGCfg(),
	}
}

func (sdr *WinStackedDRG_I) Expander() *ExpanderGraph_I {
	return &ExpanderGraph_I{
		Config_: sdr.ExpanderGraphCfg(),
	}
}

func (sdr *WinStackedDRG_I) WindowDrg() *DRG_I {
	return &DRG_I{
		Config_: sdr.WindowDRGCfg(),
	}
}

func (sdr *WinStackedDRG_I) WindowExpander() *ExpanderGraph_I {
	return &ExpanderGraph_I{
		Config_: sdr.WindowExpanderGraphCfg(),
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
		parent := chung._ithParent(node, i, d, nodes)
		parents = append(parents, parent)
	}
	return parents
}

func (chung *ChungExpanderAlgorithm_I) _ithParent(node UInt, i UInt, degree ExpanderGraphDegree, nodes ExpanderGraphNodeCount) UInt {
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
	// Call into feistel.go.
	panic("TODO")
}

func (sdr *WinStackedDRG_I) Seal(sid sector.SectorID, data []byte, randomness sector.SealRandomness) SealSetupArtifacts {
	windowCount := int(sdr.WindowCount())
	nodeSize := int(sdr.NodeSize())
	nodes := int(sdr.Nodes())
	curveModulus := sdr.Curve().FieldModulus()

	var windowData [][]byte

	for i := 0; i < len(data); i += nodeSize {
		windowData = append(windowData, data[i*nodeSize:(i+1)*nodeSize])
	}

	util.Assert(len(windowData) == windowCount)

	var windowKeyLayers [][]byte
	var finalWindowKeyLayer []byte

	commD, commDTreePath := ComputeDataCommitment(data)
	sealSeed := computeSealSeed(sid, randomness, commD)

	for i := 0; i < windowCount; i++ {
		keyLayers := sdr._generateWindowKey(sealSeed, i, sid, commD, nodes, randomness)

		lastIndex := len(keyLayers) - 1
		windowKeyLayers = append(windowKeyLayers, keyLayers[:lastIndex]...)
		finalWindowKeyLayer = append(finalWindowKeyLayer, keyLayers[lastIndex]...)
	}

	qLayer := encodeDataInPlace(data, finalWindowKeyLayer, nodeSize, &curveModulus)
	// NOTE: qLayer and data are now the same, and qLayer is introduced here for descriptive clarity only.

	replica := labelLayer(sdr.Drg(), sdr.Expander(), sealSeed, WRAPPER_LAYER_WINDOW_INDEX, nodes, nodeSize, qLayer)

	commC, commQ, commRLast, commR, commCTreePath, commQTreePath, commRLastTreePath := sdr.GenerateCommitments(replica, windowKeyLayers, qLayer)

	result := SealSetupArtifacts_I{
		CommD_:             Commitment(commD),
		CommR_:             SealedSectorCID(commR),
		CommC_:             Commitment(commC),
		CommQ_:             Commitment(commQ),
		CommRLast_:         Commitment(commRLast),
		CommDTreePath_:     commDTreePath,
		CommCTreePath_:     commCTreePath,
		CommQTreePath_:     commQTreePath,
		CommRLastTreePath_: commRLastTreePath,
		Seed_:              sealSeed,
		KeyLayers_:         windowKeyLayers,
		Replica_:           replica,
	}
	return &result
}

func (sdr *WinStackedDRG_I) _generateWindowKey(sealSeed sector.SealSeed, windowIndex int, sid sector.SectorID, commD sector.UnsealedSectorCID, nodes int, randomness sector.SealRandomness) [][]byte {
	nodeSize := int(sdr.NodeSize())
	curveModulus := sdr.Curve().FieldModulus()
	layers := int(sdr.Layers())

	keyLayers := generateSDRKeyLayers(sdr.WindowDrg(), sdr.WindowExpander(), sealSeed, windowIndex, nodes, layers, nodeSize, curveModulus)

	return keyLayers
}

func (sdr *WinStackedDRG_I) GenerateCommitments(replica []byte, windowKeyLayers [][]byte, qLayer []byte) (commC PedersenHash, commQ PedersenHash, commRLast PedersenHash, commR PedersenHash, commCTreePath file.Path, commQTreePath file.Path, commRLastTreePath file.Path) {
	commC, commCTreePath = computeCommC(windowKeyLayers, int(sdr.NodeSize()))
	commQ, commQTreePath = computeCommQ(qLayer, int(sdr.NodeSize()))
	commRLast, commRLastTreePath = BuildTree_PedersenHash(replica)
	commR = TernaryHash_PedersenHash(commC, commQ, commRLast)

	return commC, commQ, commRLast, commR, commCTreePath, commQTreePath, commRLastTreePath
}

func getProverID(minerID addr.ActorID) []byte {
	// return leb128(minerID)
	panic("TODO")
}
func computeSealSeed(sid sector.SectorID, randomness sector.SealRandomness, commD sector.UnsealedSectorCID) sector.SealSeed {
	proverId := getProverID(sid.Miner())
	sectorNumber := sid.Number()

	var preimage []byte
	preimage = append(preimage, proverId...)
	preimage = append(preimage, bigEndianBytesFromUInt(UInt(sectorNumber), 8)...)
	preimage = append(preimage, randomness...)
	preimage = append(preimage, Commitment_UnsealedSectorCID(commD)...)

	sealSeed := HashBytes_SHA256Hash(preimage)
	return sector.SealSeed(sealSeed)
}

func generateSDRKeyLayers(drg *DRG_I, expander *ExpanderGraph_I, sealSeed sector.SealSeed, window int, nodes int, layers int, nodeSize int, modulus big.Int) [][]byte {
	var keyLayers [][]byte
	var prevLayer []byte

	for i := 0; i < layers; i++ {
		currentLayer := labelLayer(drg, expander, sealSeed, window, nodeSize, nodes, prevLayer)
		keyLayers = append(keyLayers, currentLayer)
		prevLayer = currentLayer
	}

	return keyLayers
}

func labelLayer(drg *DRG_I, expander *ExpanderGraph_I, sealSeed sector.SealSeed, window int, nodeSize int, nodes int, prevLayer []byte) []byte {
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
				parents = append(parents, prevLayer[start:start+nodeSize])
			}
		}

		label := generateLabel(sealSeed, i, window, parents)
		labels = append(labels, label...)
	}

	return labels
}

// Encodes data in-place, mutating it.
func encodeDataInPlace(data []byte, key []byte, nodeSize int, modulus *big.Int) []byte {
	if len(data) != len(key) {
		panic("Key and data must be same length.")
	}

	for i := 0; i < len(data); i += nodeSize {
		copy(data[i:i+nodeSize], encodeNode(data[i:i+nodeSize], key[i:i+nodeSize], modulus, nodeSize))
	}

	return data
}

func generateLabel(sealSeed sector.SealSeed, node int, window int, dependencies []Label) []byte {
	preimage := sealSeed

	if window != WRAPPER_LAYER_WINDOW_INDEX {
		windowBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(windowBytes, uint64(window))

		preimage = append(preimage, windowBytes...)
	}

	nodeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nodeBytes, uint64(node))

	preimage = append(preimage, nodeBytes...)
	for _, dependency := range dependencies {
		preimage = append(preimage, dependency...)
	}

	return deriveLabel(preimage)
}

func deriveLabel(elements []byte) []byte {
	return trimToFr32(HashBytes_SHA256Hash(elements))
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

func computeCommQ(layerBytes []byte, nodeSize int) (PedersenHash, file.Path) {
	leaves := make([]byte, len(layerBytes)/nodeSize)
	for i := 0; i < len(leaves); i++ {
		leaves = append(leaves, layerBytes[i*nodeSize:(i+1)*nodeSize]...)
	}

	return BuildTree_PedersenHash(leaves)
}

func hashColumn(column []Label) PedersenHash {
	var preimage []byte
	for _, label := range column {
		preimage = append(preimage, label...)
	}
	return HashBytes_PedersenHash(preimage)
}

func (sdr *WinStackedDRG_I) CreateSealProof(challengeSeed sector.InteractiveSealRandomness, aux sector.ProofAuxTmp) sector.SealProof {
	privateProof := sdr.CreatePrivateSealProof(challengeSeed, aux)

	// Sanity check: newly-created proofs must pass verification.
	util.Assert(sdr.VerifyPrivateSealProof(privateProof, aux.Seed(), challengeSeed, aux.CommD(), aux.CommR()))

	return sdr.CreateOfflineCircuitProof(privateProof, aux)
}

func (sdr *WinStackedDRG_I) CreatePrivateSealProof(randomness sector.InteractiveSealRandomness, aux sector.ProofAuxTmp) (privateProof PrivateOfflineProof) {
	sealSeed := aux.Seed()
	nodeSize := UInt(sdr.NodeSize())
	wrapperChallenges, windowChallenges := sdr._generateOfflineChallenges(sealSeed, randomness, sdr.Challenges(), sdr.WindowChallenges())

	dataTree := LoadMerkleTree(aux.CommDTreePath())
	columnTree := LoadMerkleTree(aux.CommCTreePath())
	replicaTree := LoadMerkleTree(aux.PersistentAux().CommRLastTreePath())
	qTree := LoadMerkleTree(aux.CommQTreePath())

	windows := int(sdr.WindowCount())
	windowSize := int(uint64(sdr.Cfg().SealCfg().SectorSize()) / UInt(sdr.WindowCount()))

	for c := range windowChallenges {
		columnProofs := createColumnProofs(sdr.WindowDrg(), sdr.WindowExpander(), UInt(c), nodeSize, columnTree, aux, windows, windowSize)
		privateProof.ColumnProofs = append(privateProof.ColumnProofs, columnProofs...)

		windowProof := createWindowProof(sdr.WindowDrg(), sdr.WindowExpander(), UInt(c), nodeSize, dataTree, columnTree, qTree, aux, windows, windowSize)
		privateProof.WindowProofs = append(privateProof.WindowProofs, windowProof)
	}

	for c := range wrapperChallenges {
		wrapperProof := createWrapperProof(sdr.Drg(), sdr.Expander(), sealSeed, UInt(c), nodeSize, qTree, replicaTree, aux, windows, windowSize)
		privateProof.WrapperProofs = append(privateProof.WrapperProofs, wrapperProof)
	}

	return privateProof
}

// Verify a private proof.
// NOTE: Verification of a private proof is exactly the computation we will prove we have performed in a zk-SNARK.
// If we can verifiably prove that we have performed the verification of a private proof, then we need not reveal the proof itself.
// Since the zk-SNARK circuit proof is much smaller than the private proof, this allows us to save space on the chain (at the cost of increased computation to generate the zk-SNARK proof).
func (sdr *WinStackedDRG_I) VerifyPrivateSealProof(privateProof PrivateOfflineProof, sealSeed sector.SealSeed, randomness sector.InteractiveSealRandomness, commD Commitment, commR sector.SealedSectorCID) bool {
	nodeSize := int(sdr.NodeSize())
	windowCount := int(sdr.WindowCount())
	windowSize := int(UInt(sdr.Cfg().SealCfg().SectorSize()) / UInt(sdr.WindowCount())) // TOOD: Make this a function.
	layers := int(sdr.Layers())
	curveModulus := sdr.Curve().FieldModulus()
	windowChallenges, wrapperChallenges := sdr._generateOfflineChallenges(sealSeed, randomness, sdr.Challenges(), sdr.WindowChallenges())

	windowProofs := privateProof.WindowProofs
	columnProofs := privateProof.ColumnProofs
	wrapperProofs := privateProof.WrapperProofs

	// commC, commQ, and commRLast must be the same for all challenge proofs, so we can arbitrarily verify against the first.
	firstColumnProof := columnProofs[0]
	firstWrapperProof := wrapperProofs[0]
	commC := firstColumnProof.InclusionProof.Root()
	commQ := firstWrapperProof.QLayerProofs[0].Root()
	commRLast := firstWrapperProof.ReplicaProof.Root()

	windowDrgParentCount := int(sdr.WindowDRGCfg().Degree())
	windowExpanderParentCount := int(sdr.WindowDRGCfg().Degree())
	wrapperExpanderParentCount := int(sdr.ExpanderGraphCfg().Degree())

	for i, challenge := range windowChallenges {
		// Verify one OfflineSDRChallengeProof.
		windowProof := windowProofs[i]
		dataProof := windowProof.DataProof
		qLayerProof := windowProof.QLayerProof

		// Verify column proofs and that they represent the right columns.
		columnElements := getColumnElements(sdr.Drg(), sdr.Expander(), challenge)

		// Check column openings.
		for i, columnElement := range columnElements {
			columnProof := columnProofs[i]

			// The provided column proofs must correspond to the expected columns.
			if !columnProof.Verify(commC, UInt(columnElement)) {
				return false
			}
		}

		// Check labeling.
		for w := 0; w < windowCount; w++ {
			for layer := 0; layer < layers; layer++ {
				var parents []Label

				// First column proof is the challenge.
				// Then the DRG parents.
				for _, drgParentProof := range columnProofs[1 : 1+windowDrgParentCount] {
					parent := drgParentProof.Column[layer]
					parents = append(parents, parent)
				}
				// And the expander parents, if not the first layer.
				if layer > 0 {
					for _, expanderParentProof := range columnProofs[1+windowDrgParentCount : 1+windowExpanderParentCount] {
						parent := expanderParentProof.Column[layer-1]
						parents = append(parents, parent)
					}
				}

				calculatedLabel := generateLabel(sealSeed, i, w, parents)

				if layer == layers-1 {
					// Last layer includes encoding.
					dataNode := dataProof.Leaf()
					qLayerNode := qLayerProof.Leaf()

					if !dataProof.Verify(commD, UInt(windowSize*w)+challenge) {
						return false
					}

					encodedNode := encodeNode(dataNode, calculatedLabel, &curveModulus, nodeSize)

					if !bytes.Equal(encodedNode, qLayerNode) {
						return false
					}

				} else {
					providedLabel := columnProofs[columnElements[0]].Column[layer]

					if !bytes.Equal(calculatedLabel, providedLabel) {
						return false
					}
				}
			}
		}
	}

	for i, challenge := range wrapperChallenges {
		wrapperProof := wrapperProofs[i]
		replicaProof := wrapperProof.ReplicaProof
		qLayerProofs := wrapperProof.QLayerProofs

		if !replicaProof.Verify(commRLast, challenge) {
			return false
		}

		var parents []Label
		for i := 0; i < wrapperExpanderParentCount; i++ {
			parent := qLayerProofs[i].Leaf()
			parents = append(parents, parent)
		}

		label := generateLabel(sealSeed, i, windowCount+1, parents)
		replicaNode := replicaProof.Leaf()

		if !bytes.Equal(label, replicaNode) {
			return false
		}
	}

	commRCalculated := TernaryHash_PedersenHash(commC, commQ, commRLast)

	if !bytes.Equal(commRCalculated, AsBytes_SealedSectorCID(commR)) {
		return false
	}

	return true
}

func createColumnProofs(drg *DRG_I, expander *ExpanderGraph_I, challenge UInt, nodeSize UInt, columnTree MerkleTree, aux sector.ProofAuxTmp, windows int, windowSize int) []SDRColumnProof {
	columnElements := getColumnElements(drg, expander, challenge)

	var columnProofs []SDRColumnProof
	for c := range columnElements {
		chall := UInt(c)

		columnProof := createColumnProof(chall, nodeSize, windows, windowSize, columnTree, aux)
		columnProofs = append(columnProofs, columnProof)
	}

	return columnProofs
}

func createWindowProof(drg *DRG_I, expander *ExpanderGraph_I, challenge UInt, nodeSize UInt, dataTree MerkleTree, columnTree MerkleTree, qLayerTree MerkleTree, aux sector.ProofAuxTmp, windows int, windowSize int) (proof OfflineWindowProof) {
	columnElements := getColumnElements(drg, expander, challenge)

	var columnProofs []SDRColumnProof
	for c := range columnElements {
		chall := UInt(c)

		columnProof := createColumnProof(chall, nodeSize, windows, windowSize, columnTree, aux)
		columnProofs = append(columnProofs, columnProof)
	}

	dataProof := dataTree.ProveInclusion(challenge)
	qLayerProof := qLayerTree.ProveInclusion(challenge)

	proof = OfflineWindowProof{
		DataProof:   dataProof,
		QLayerProof: qLayerProof,
	}

	return proof
}

func createWrapperProof(drg *DRG_I, expander *ExpanderGraph_I, sealSeed sector.SealSeed, challenge UInt, nodeSize UInt, qTree MerkleTree, replicaTree MerkleTree, aux sector.ProofAuxTmp, windows int, windowSize int) (proof OfflineWrapperProof) {
	proof.ReplicaProof = replicaTree.ProveInclusion(challenge)

	parents := expander.Parents(challenge)

	for _, parent := range parents {
		proof.QLayerProofs = append(proof.QLayerProofs, qTree.ProveInclusion(parent))
	}
	return proof
}

func getColumnElements(drg *DRG_I, expander *ExpanderGraph_I, challenge UInt) (columnElements []UInt) {
	columnElements = append(columnElements, challenge)
	columnElements = append(columnElements, drg.Parents(challenge)...)
	columnElements = append(columnElements, expander.Parents(challenge)...)

	return columnElements
}

func createColumnProof(c UInt, nodeSize UInt, windowSize int, windows int, columnTree MerkleTree, aux sector.ProofAuxTmp) (columnProof SDRColumnProof) {
	layers := aux.KeyLayers()
	var column []Label

	for w := 0; w < windows; w++ {
		for i := 0; i < len(layers); i++ {
			start := (w * windowSize) + int(c)
			end := start + int(nodeSize)
			column = append(column, layers[i][start:end])
		}
	}
	columnProof = SDRColumnProof{
		Column:         column,
		InclusionProof: columnTree.ProveInclusion(c),
	}

	return columnProof
}

type PrivateOfflineProof struct {
	ColumnProofs  []SDRColumnProof
	WindowProofs  []OfflineWindowProof
	WrapperProofs []OfflineWrapperProof
}

type OfflineWindowProof struct {
	// TODO: these proofs need to depend on hash function.
	DataProof   InclusionProof // SHA256
	QLayerProof InclusionProof
}

type OfflineWrapperProof struct {
	ReplicaProof InclusionProof // Pedersen
	QLayerProofs []InclusionProof
}

func (ip *InclusionProof_I) Leaf() []byte {
	panic("TODO")
}

func (ip *InclusionProof_I) LeafIndex() UInt {
	panic("TODO")
}

func (ip *InclusionProof_I) Root() Commitment {
	panic("TODO")
}

func (mt *MerkleTree_I) ProveInclusion(challenge UInt) InclusionProof {
	panic("TODO")
}

func (mt *MerkleTree_I) Leaf(index UInt) []byte {
	panic("TODO")
}

func LoadMerkleTree(path file.Path) MerkleTree {
	panic("TODO")
}

func (ip *InclusionProof_I) Verify(root []byte, challenge UInt) bool {
	// FIXME: need to verify proof length of private inclusion proofs.
	panic("TODO")
}

type SDRColumnProof struct {
	Column         []Label
	InclusionProof InclusionProof
}

func (proof *SDRColumnProof) Verify(root []byte, challenge UInt) bool {
	if !bytes.Equal(hashColumn(proof.Column), proof.InclusionProof.Leaf()) {
		return false
	}

	if proof.InclusionProof.LeafIndex() != challenge {
		return false
	}

	return proof.InclusionProof.Verify(root, challenge)
}

func (sdr *WinStackedDRG_I) CreateOfflineCircuitProof(proof PrivateOfflineProof, aux sector.ProofAuxTmp) sector.SealProof {
	// partitions := sdr.Partitions()
	// publicInputs := GeneratePublicInputs()

	panic("TODO")
	var bytes []byte

	sealProof := sector.SealProof_I{
		CircuitType_: sector.SealCircuitType_WinStackedSDR,
		ProofBytes_:  bytes,
	}

	return &sealProof
}

func (sdr *WinStackedDRG_I) _generateOfflineChallenges(sealSeed sector.SealSeed, randomness sector.InteractiveSealRandomness, wrapperChallengeCount WinStackedDRGChallenges, windowChallengeCount WinStackedDRGWindowChallenges) (windowChallenges []UInt, wrapperChallenges []UInt) {
	wrapperChallenges = generateOfflineChallenges(int(sdr.Nodes()), sealSeed, randomness, int(wrapperChallengeCount))
	windowChallenges = generateOfflineChallenges(int(sdr.WindowDRGCfg().Nodes()), sealSeed, randomness, int(windowChallengeCount))

	return windowChallenges, wrapperChallenges
}

func generateOfflineChallenges(challengeRange int, sealSeed sector.SealSeed, randomness sector.InteractiveSealRandomness, challengeCount int) []UInt {
	var challenges []UInt
	challengeRangeSize := challengeRange - 1 // Never challenge the first node.
	challengeModulus := new(big.Int)
	challengeModulus.SetUint64(uint64(challengeRangeSize))

	// Maybe factor this into a separate function, since the logic is the same...

	for i := 0; i < challengeCount; i++ {
		var preimage []byte
		preimage = append(preimage, sealSeed...)
		preimage = append(preimage, randomness...)
		preimage = append(preimage, littleEndianBytesFromInt(i, 4)...)

		hash := HashBytes_SHA256Hash(preimage)
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
	// TODO: Make this a method of WinStackedDRG.
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

func (sdr *WinStackedDRG_I) VerifySeal(sv sector.SealVerifyInfo) bool {
	onChain := sv.OnChain()
	sealProof := onChain.Proof()
	commR := sector.SealedSectorCID(onChain.SealedCID())
	commD := sector.UnsealedSectorCID(sv.UnsealedCID())
	sealSeed := computeSealSeed(sv.SectorID(), sv.Randomness(), commD)

	wrapperChallenges, windowChallenges := sdr._generateOfflineChallenges(sealSeed, sv.InteractiveRandomness(), sdr.Challenges(), sdr.WindowChallenges())
	return sdr._verifyOfflineCircuitProof(commD, commR, sealSeed, windowChallenges, wrapperChallenges, sealProof)
}

func ComputeUnsealedSectorCIDFromPieceInfos(sectorSize sector.SectorSize, pieceInfos []PieceInfo) (unsealedCID sector.UnsealedSectorCID, err error) {
	rootPieceInfo := computeRootPieceInfo(pieceInfos)
	rootSize := rootPieceInfo.Size()

	if rootSize != uint64(sectorSize) {
		return unsealedCID, errors.New("Wrong sector size.")
	}

	return UnsealedSectorCID(AsBytes_PieceCID(rootPieceInfo.PieceCID())), nil
}

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

func (sdr *WinStackedDRG_I) _verifyOfflineCircuitProof(commD sector.UnsealedSectorCID, commR sector.SealedSectorCID, sealSeed sector.SealSeed, windowChallenges []UInt, wrapperChallenges []UInt, sv sector.SealProof) bool {
	//publicInputs := GeneratePublicInputs()
	panic("TODO")
}

////////////////////////////////////////////////////////////////////////////////
// PoSt

func getChallengedSectors(cfg sector.PoStCfg, sectorIDs []sector.SectorID, randomness sector.PoStRandomness, eligibleSectors []sector.SectorID, candidateCount int) (sectors []sector.SectorID) {
	for i := 0; i < candidateCount; i++ {
		sector := generateSectorChallenge(randomness, i, sectorIDs)
		sectors = append(sectors, sector)
	}

	return sectors
}

func generateSectorChallenge(randomness sector.PoStRandomness, n int, sectorIDs []sector.SectorID) (sector sector.SectorID) {
	preimage := append(randomness, littleEndianBytesFromInt(n, 8)...)
	hash := SHA256Hash(preimage)
	sectorChallenge := bigIntFromLittleEndianBytes(hash)

	challengeModulus := new(big.Int)
	challengeModulus.SetUint64(uint64(len(sectorIDs)))

	sectorIndex := sectorChallenge.Mod(sectorChallenge, challengeModulus)
	return sectorIDs[int(sectorIndex.Uint64())]
}

func generateLeafChallenge(randomness sector.PoStRandomness, sectorChallengeIndex UInt, leafChallengeIndex int, nodes int, challengeRangeSize int) UInt {
	preimage := append(randomness, littleEndianBytesFromUInt(sectorChallengeIndex, 8)...)
	preimage = append(preimage, littleEndianBytesFromInt(leafChallengeIndex, 8)...)
	hash := SHA256Hash(preimage)
	bigHash := bigIntFromLittleEndianBytes(hash)

	challengeSpaceSize := nodes / challengeRangeSize
	challengeModulus := new(big.Int)
	challengeModulus.SetUint64(UInt(challengeSpaceSize))

	leafChallenge := bigHash.Mod(bigHash, challengeModulus)

	return leafChallenge.Uint64()
}

func generateCandidate(algorithm sector.SealAlgorithm, cfg sector.PoStCfg, randomness sector.PoStRandomness, aux sector.PersistentProofAux, sectorID sector.SectorID, sectorChallengeIndex UInt) sector.PoStCandidate {
	var candidate sector.PoStCandidate
	switch algorithm {
	case sector.SealAlgorithm_StackedDRG:
		panic("TODO")
	case sector.SealAlgorithm_WinStackedDRG:
		sdr := WinStackedDRG_I{}
		candidate = sdr._generateCandidate(cfg, randomness, aux, sectorID, sectorChallengeIndex)
	}
	return candidate
}

func (sdr *WinStackedDRG_I) _generateCandidate(cfg sector.PoStCfg, randomness sector.PoStRandomness, aux sector.PersistentProofAux, sectorID sector.SectorID, sectorChallengeIndex UInt) sector.PoStCandidate {
	nodes := int(cfg.Nodes())
	leafChallengeCount := int(cfg.LeafChallengeCount())
	challengeRangeSize := int(cfg.ChallengeRangeSize())
	treePath := aux.CommRLastTreePath()
	tree := LoadMerkleTree(treePath)

	var data []byte
	var inclusionProofs []InclusionProof
	for i := 0; i < leafChallengeCount; i++ {
		leafChallenge := generateLeafChallenge(randomness, sectorChallengeIndex, i, nodes, challengeRangeSize)

		for j := 0; j < challengeRangeSize; j++ {
			leafIndex := leafChallenge + UInt(j)
			data = append(data, tree.Leaf(leafIndex)...)
			inclusionProof := tree.ProveInclusion(leafIndex)
			inclusionProofs = append(inclusionProofs, inclusionProof)
		}
	}

	partialTicket := computePartialTicket(randomness, sectorID, data)

	privateProof := InternalPrivateCandidateProof{
		InclusionProofs: inclusionProofs,
	}

	candidate := sector.PoStCandidate_I{
		PartialTicket_:  partialTicket,
		PrivateProof_:   privateProof.externalize(sector.SealAlgorithm_WinStackedDRG),
		SectorID_:       sectorID,
		ChallengeIndex_: sectorChallengeIndex,
	}
	return &candidate
}

func computePartialTicket(randomness sector.PoStRandomness, sectorID sector.SectorID, data []byte) sector.PartialTicket {
	preimage := randomness
	preimage = append(preimage, getProverID(sectorID.Miner())...)
	preimage = append(preimage, littleEndianBytesFromUInt(UInt(sectorID.Number()), 8)...)
	preimage = append(preimage, data...)
	partialTicket := sector.PartialTicket(HashBytes_PedersenHash(preimage))

	return partialTicket
}

type PoStCandidatesMap map[sector.SealAlgorithm][]sector.PoStCandidate

func CreatePoStProof(cfg sector.PoStCfg, privateCandidateProofs []PrivatePostCandidateProof, challengeSeed sector.PoStRandomness) []sector.PoStProof {
	var proofsMap map[sector.SealAlgorithm][]PrivatePostCandidateProof

	for _, proof := range privateCandidateProofs {
		algorithm := proof.Algorithm()
		proofsMap[algorithm] = append(proofsMap[algorithm], proof)
	}

	var circuitProofs []sector.PoStProof
	for algorithm, proofs := range proofsMap {
		privateProof := createPrivatePoStProof(algorithm, proofs, challengeSeed)
		// Hmmmmm: we cannot perform this santiy check without the sectorIDs. Should we require them just for that purpose or perform the check earlier?
		// Sanity check: newly-created proofs must pass verification.
		//util.Assert(sdr.VerifyPrivatePoStProof(privateProof, candidates, sectorIDs, sectorCommitments))

		circuitProof := createPoStCircuitProof(cfg, algorithm, privateProof)
		circuitProofs = append(circuitProofs, circuitProof)
	}

	return circuitProofs
}

type PrivatePoStProof struct {
	ChallengeSeed   sector.PoStRandomness
	CandidateProofs []PrivatePostCandidateProof
}

func createPrivatePoStProof(algorithm sector.SealAlgorithm, candidateProofs []PrivatePostCandidateProof, challengeSeed sector.PoStRandomness) PrivatePoStProof {

	return PrivatePoStProof{
		ChallengeSeed:   challengeSeed,
		CandidateProofs: candidateProofs,
	}
}

type InternalPrivateCandidateProof struct {
	InclusionProofs []InclusionProof
}

// This exists because we need to pass private proofs out of filproofs for winner selection.
// Actually implementing it would (will?) be tedious, since it means doing the same for InclusionProofs.
func (p *InternalPrivateCandidateProof) externalize(algorithm sector.SealAlgorithm) sector.PrivatePoStCandidateProof {
	return &sector.PrivatePoStCandidateProof_I{
		Algorithm_:    algorithm,
		Externalized_: []byte{}, // Unimplemented.
	}
}

// This is the inverse of InternalPrivateCandidateProof.externalize and equally tedious.
func newInternalPrivateProof(externalPrivateProof PrivatePostCandidateProof) InternalPrivateCandidateProof {
	return InternalPrivateCandidateProof{}
}

func (sdr *WinStackedDRG_I) VerifyInternalPrivateCandidateProof(cfg sector.PoStCfg, p *InternalPrivateCandidateProof, challengeSeed sector.PoStRandomness, candidate sector.PoStCandidate, commRLast Commitment) bool {
	util.Assert(candidate.PrivateProof() == nil)

	nodes := int(cfg.Nodes())
	challengeRangeSize := int(cfg.ChallengeRangeSize())

	sectorID := candidate.SectorID()
	claimedPartialTicket := candidate.PartialTicket()

	allInclusionProofs := p.InclusionProofs

	var ticketData []byte

	for _, p := range allInclusionProofs {
		ticketData = append(ticketData, p.Leaf()...)
	}

	// Check partial ticket
	calculatedTicket := computePartialTicket(challengeSeed, sectorID, ticketData)

	if len(calculatedTicket) != len(claimedPartialTicket) {
		return false
	}
	for i, byte := range claimedPartialTicket {
		if byte != calculatedTicket[i] {
			return false
		}
	}

	// Helper to get InclusionProofs sequentially.
	next := func() InclusionProof {
		if len(allInclusionProofs) < 1 {
			return nil
		}

		proof := allInclusionProofs[0]
		allInclusionProofs = allInclusionProofs[1:]
		return proof
	}

	// Check all inclusion proofs.
	for i := 0; i < int(cfg.LeafChallengeCount()); i++ {
		leafChallenge := generateLeafChallenge(challengeSeed, candidate.ChallengeIndex(), i, nodes, challengeRangeSize)
		for j := 0; j < challengeRangeSize; j++ {
			leafIndex := leafChallenge + UInt(j)
			proof := next()
			if proof == nil {
				// All required inclusion proofs must be provided.
				return false
			}
			if !proof.Verify(commRLast, leafIndex) {
				return false
			}
		}
	}

	return true
}

func (sdr *WinStackedDRG_I) VerifyPrivatePoStProof(cfg sector.PoStCfg, privateProof PrivatePoStProof, candidates []sector.PoStCandidate, sectorIDs []sector.SectorID, sectorCommitments sector.SectorCommitments) bool {
	// This is safe by construction.
	challengeSeed := privateProof.ChallengeSeed

	for i, p := range privateProof.CandidateProofs {
		proof := newInternalPrivateProof(p)

		candidate := candidates[i]
		ci := candidate.ChallengeIndex()
		expectedSectorID := sectorIDs[ci]

		challengedSectorID := generateSectorChallenge(challengeSeed, i, sectorIDs)

		if expectedSectorID != challengedSectorID {
			return false
		}

		commRLast := sectorCommitments[expectedSectorID]

		if !sdr.VerifyInternalPrivateCandidateProof(cfg, &proof, challengeSeed, candidate, commRLast) {
			return false
		}
	}
	return true
}

func createPoStCircuitProof(postCfg sector.PoStCfg, algorithm sector.SealAlgorithm, privateProof PrivatePoStProof) (proof sector.PoStProof) {
	switch algorithm {
	case sector.SealAlgorithm_WinStackedDRG:
		sdr := WinStackedDRG_I{}
		proof = sdr._createPoStCircuitProof(postCfg, privateProof)
	}
	return proof
}

func (sdr *WinStackedDRG_I) _createPoStCircuitProof(postCfg sector.PoStCfg, privateProof PrivatePoStProof) sector.PoStProof {
	panic("TODO")

	postType := postCfg.Type()

	var circuitType sector.PoStCircuitType

	switch postType {
	case sector.PoStType_ElectionPoSt:
		circuitType = sector.PoStCircuitType_WinStackedSDRElectionPoSt
	case sector.PoStType_SurprisePoSt:
		circuitType = sector.PoStCircuitType_WinStackedSDRSurprisePoSt
	}

	postProof := sector.PoStProof_I{
		Type_:        postCfg.Type(),
		CircuitType_: circuitType,
	}

	return &postProof
}

func (pv *PoStVerifier_I) _verifyPoStProof(sv sector.PoStVerifyInfo) bool {
	// commT := sv.CommT()
	// candidates := sv.Candidates()
	// randomness := sv.Randomness()
	// postProof := sv.OnChain.Proof()

	// Verify circuit proof.
	panic("TODO")
}

////////////////////////////////////////////////////////////////////////////////
// General PoSt

func generatePoStCandidates(cfg sector.PoStCfg, challengeSeed sector.PoStRandomness, eligibleSectors []sector.SectorID, candidateCount int, sectorStore sector_index.SectorStore) (candidates []sector.PoStCandidate) {
	challengedSectors := getChallengedSectors(cfg, eligibleSectors, challengeSeed, eligibleSectors, candidateCount)

	for i, sectorID := range challengedSectors {
		proofAux := sectorStore.GetSectorPersistentProofAux(sectorID)
		sealAlgorithm := sectorStore.GetSectorSealAlgorithm(sectorID).As_a()

		candidate := generateCandidate(sealAlgorithm, cfg, challengeSeed, proofAux, sectorID, UInt(i))

		candidates = append(candidates, candidate)
	}

	return candidates
}

////////////////////////////////////////////////////////////////////////////////
// Election PoSt

func GenerateElectionPoStCandidates(cfg sector.PoStCfg, challengeSeed sector.PoStRandomness, eligibleSectors []sector.SectorID, candidateCount int, sectorStore sector_index.SectorStore) (candidates []sector.PoStCandidate) {
	return generatePoStCandidates(cfg, challengeSeed, eligibleSectors, candidateCount, sectorStore)
}

func CreateElectionPoStProof(cfg sector.PoStCfg, privateCandidateProofs []PrivatePostCandidateProof, challengeSeed sector.PoStRandomness) []sector.PoStProof {
	return CreatePoStProof(cfg, privateCandidateProofs, challengeSeed)
}

func (pv *PoStVerifier_I) VerifyElectionPoSt(sv sector.PoStVerifyInfo) bool {
	return pv._verifyPoStProof(sv)
}

////////////////////////////////////////////////////////////////////////////////
// Surprise PoSt

func GenerateSurprisePoStCandidates(challengeSeed sector.PoStRandomness, eligibleSectors []sector.SectorID, candidateCount int, sectorStore sector_index.SectorStore) []sector.PoStCandidate {
	panic("TODO")
}

func CreateSurprisePoStProof(cfg sector.PoStCfg, privateCandidateProofs []PrivatePostCandidateProof, challengeSeed sector.PoStRandomness) []sector.PoStProof {
	return CreatePoStProof(cfg, privateCandidateProofs, challengeSeed)
}

func (pv *PoStVerifier_I) VerifySurprisePoSt(sv sector.PoStVerifyInfo) bool {
	return pv._verifyPoStProof(sv)
}

////////////////////////////////////////////////////////////////////////////////
/// Generic Hashing and Merkle Tree generation

/// Binary hash compression.
// BinaryHash<T>
func BinaryHash_T(left []byte, right []byte) util.T {
	var preimage = append(left, right...)
	return HashBytes_T(preimage)
}

func TernaryHash_T(a []byte, b []byte, c []byte) util.T {
	var preimage = append(a, append(b, c...)...)
	return HashBytes_T(preimage)
}

// BinaryHash<PedersenHash>
func BinaryHash_PedersenHash(left []byte, right []byte) PedersenHash {
	return PedersenHash{}
}

func TernaryHash_PedersenHash(a []byte, b []byte, c []byte) PedersenHash {
	return PedersenHash{}
}

// BinaryHash<SHA256Hash>
func BinaryHash_SHA256Hash(left []byte, right []byte) SHA256Hash {
	result := SHA256Hash{}
	return trimToFr32(result)
}

func TernaryHash_SHA256Hash(a []byte, b []byte, c []byte) SHA256Hash {
	return SHA256Hash{}
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

	return result
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

func bigIntFromBigEndianBytes(bytes []byte) *big.Int {
	return new(big.Int).SetBytes(bytes)
}

// size is number of bytes to return
func littleEndianBytesFromBigInt(z *big.Int, size int) []byte {
	bytes := z.Bytes()[0:size]
	reverse(bytes)

	return bytes
}

// size is number of bytes to return
func bigEndianBytesFromBigInt(z *big.Int, size int) []byte {
	return z.Bytes()[0:size]
}

func littleEndianBytesFromInt(n int, size int) []byte {
	z := new(big.Int)
	z.SetInt64(int64(n))
	return littleEndianBytesFromBigInt(z, size)
}

func bigEndianBytesFromInt(n int, size int) []byte {
	z := new(big.Int)
	z.SetInt64(int64(n))
	return bigEndianBytesFromBigInt(z, size)
}

func littleEndianBytesFromUInt(n UInt, size int) []byte {
	z := new(big.Int)
	z.SetUint64(uint64(n))
	return littleEndianBytesFromBigInt(z, size)
}

func bigEndianBytesFromUInt(n UInt, size int) []byte {
	z := new(big.Int)
	z.SetUint64(uint64(n))
	return bigEndianBytesFromBigInt(z, size)
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
