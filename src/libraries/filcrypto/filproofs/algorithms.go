package filproofs

import "math"
import "math/rand"
import big "math/big"
import "encoding/binary"

import util "github.com/filecoin-project/specs/util"
import file "github.com/filecoin-project/specs/systems/filecoin_files/file"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"

type Blake2sHash Bytes32
type PedersenHash Bytes32
type Bytes32 []byte
type UInt = util.UInt

func SDRParams(sealCfg sector.SealCfg) *StackedDRG_I {
	// TODO: Bridge constants with orient model.
	const LAYERS = 10
	const NODE_SIZE = 32
	const OFFLINE_CHALLENGES = 6666
	const FEISTEL_ROUNDS = 3
	var FEISTEL_KEYS = [FEISTEL_ROUNDS]UInt{1, 2, 3}
	var FIELD_MODULUS = new(big.Int)
	// https://github.com/zkcrypto/pairing/blob/master/src/bls12_381/fr.rs#L4
	FIELD_MODULUS.SetString("52435875175126190479447740508185965837690552500527637822603658699938581184513", 10)

	nodes := UInt(sealCfg.SectorSize() / NODE_SIZE)

	return &StackedDRG_I{
		Layers_:     StackedDRGLayers(LAYERS),
		Challenges_: StackedDRGChallenges(OFFLINE_CHALLENGES),
		NodeSize_:   StackedDRGNodeSize(NODE_SIZE),
		Algorithm_:  &StackedDRG_Algorithm_I{},
		DRGCfg_: &DRGCfg_I{
			Algorithm_: &DRGCfg_Algorithm_I{
				ParentsAlgorithm_: DRGCfg_Algorithm_ParentsAlgorithm_Make_DRSample(&DRGCfg_Algorithm_ParentsAlgorithm_DRSample_I{}),
				RNGAlgorithm_:     DRGCfg_Algorithm_RNGAlgorithm_Make_ChaCha20(&DRGCfg_Algorithm_RNGAlgorithm_ChaCha20_I{}),
			},
			Degree_: 6,
			Nodes_:  DRGNodeCount(nodes),
		},
		ExpanderGraphCfg_: &ExpanderGraphCfg_I{
			Algorithm_: ExpanderGraphCfg_Algorithm_Make_ChungExpanderAlgorithm(
				&ChungExpanderAlgorithm_I{
					PermutationAlgorithm_: ChungExpanderAlgorithm_PermutationAlgorithm_Make_Feistel(&Feistel_I{
						Keys_:   FEISTEL_KEYS[:],
						Rounds_: FEISTEL_ROUNDS,
						HashFunction_: ChungExpanderPermutationFeistelHashFunction_Make_Blake2S(
							&ChungExpanderPermutationFeistelHashFunction_Blake2S_I{}),
					}),
				}),
			Degree_: 8,
			Nodes_:  ExpanderGraphNodeCount(nodes),
		},
		Curve_: &EllipticCurve_I{
			FieldModulus_: *FIELD_MODULUS,
		},
	}
}

func (drg *DRG_I) Parents(node UInt) []UInt {
	config := drg.Config()
	degree := UInt(config.Degree())
	return config.Algorithm().ParentsAlgorithm().As_DRSample().Impl().Parents(degree, node)
}

// TODO: Verify this. Both the port from impl and the algorithm.
func (drs *DRGCfg_Algorithm_ParentsAlgorithm_DRSample_I) Parents(degree, node UInt) (parents []UInt) {
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
}

func randInRange(lowInclusive int, highExclusive int) UInt {
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
	panic("TODO")
}

func (sdr *StackedDRG_I) Seal(sid sector.SectorID, commD sector.UnsealedSectorCID, data []byte) SealSetupArtifacts {
	replicaID := ComputeReplicaID(sid, commD)

	drg := DRG_I{
		Config_: sdr.DRGCfg(),
	}

	expander := ExpanderGraph_I{
		Config_: sdr.ExpanderGraphCfg(),
	}

	nodeSize := int(sdr.NodeSize())
	nodes := len(data) / nodeSize
	curveModulus := sdr.Curve().FieldModulus()
	layers := int(sdr.Layers())

	keyLayers := generateSDRKeyLayers(&drg, &expander, replicaID, nodes, layers, nodeSize, curveModulus)
	key := keyLayers[len(keyLayers)-1]

	replica := encodeData(data, key, nodeSize, &curveModulus)

	commRLast, commRLastTreePath := RepHash_PedersenHash(replica)
	commC, commCTreePath := computeCommC(keyLayers, nodeSize)
	commR := RepCompress_PedersenHash(commC, commRLast)

	result := SealSetupArtifacts_I{
		CommR_:             SealedSectorCID(commR),
		CommC_:             sector.Commitment(commC),
		CommRLast_:         sector.Commitment(commRLast),
		CommRLastTreePath_: commRLastTreePath,
		CommCTreePath_:     commCTreePath,
		KeyLayers_:         keyLayers,
		Replica_:           replica,
	}
	return &result
}

func ComputeReplicaID(sid sector.SectorID, commD sector.UnsealedSectorCID) Bytes32 {
	_, _ = sid.MinerID(), (sid.Number())

	// FIXME: Implement
	return Bytes32{}
}

func generateSDRKeyLayers(drg *DRG_I, expander *ExpanderGraph_I, replicaID []byte, nodes int, layers int, nodeSize int, modulus big.Int) [][]byte {
	var keyLayers [][]byte
	var prevLayer []byte

	for i := 0; i <= layers; i++ {
		currentLayer := labelLayer(drg, expander, replicaID, nodes, nodeSize, prevLayer)
		keyLayers = append(keyLayers, currentLayer)
		prevLayer = currentLayer
	}

	return keyLayers
}

func labelLayer(drg *DRG_I, expander *ExpanderGraph_I, replicaID []byte, nodeSize int, nodes int, prevLayer []byte) []byte {
	size := nodes * nodeSize
	labels := make([]byte, size)

	for i := 0; i < nodes; i++ {
		var parents []byte

		// The first node of every layer has no DRG Parents.
		if i > 0 {
			for parent := range drg.Parents(UInt(i)) {
				start := parent * nodeSize
				parents = append(parents, labels[start:start+nodeSize]...)
			}
		}

		// The first layer has no expander parents.
		if prevLayer != nil {
			for parent := range expander.Parents(UInt(i)) {
				start := parent * nodeSize
				parents = append(parents, labels[start:start+nodeSize]...)
			}
		}

		label := generateLabel(replicaID, i, parents)
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

func generateLabel(replicaID []byte, node int, dependencies []byte) []byte {
	nodeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nodeBytes, uint64(node))

	preimage := append(replicaID, nodeBytes...)
	preimage = append(preimage, dependencies...)

	return deriveLabel(preimage)
}

func deriveLabel(elements []byte) []byte {
	return WideRepCompress_Blake2sHash(elements)
}

func computeCommC(keyLayers [][]byte, nodeSize int) (PedersenHash, file.Path) {
	leaves := make([]byte, len(keyLayers[0]))

	// For each node in the graph,
	for start := 0; start < len(leaves); start += nodeSize {
		end := start + nodeSize

		var column []byte
		// Concatenate that node's label at each layer, in order, into a column.
		for i := 0; i < len(keyLayers); i++ {
			label := keyLayers[i][start:end]
			column = append(column, label...)
		}

		// And hash that column to create the leaf of a new tree.
		hashed := hashColumn(column)
		copy(leaves[start:end], hashed[:])
	}

	// Return the root of and path to the column tree.
	return RepHash_PedersenHash(leaves)
}

func hashColumn(column []byte) PedersenHash {
	return WideRepCompress_PedersenHash(column)
}

func (sdr *StackedDRG_I) CreateSealProof(randomSeed sector.SealRandomness, aux sector.ProofAuxTmp) sector.SealProof {
	replicaID := ComputeReplicaID(aux.SectorID(), aux.CommD())

	drg := DRG_I{
		Config_: sdr.DRGCfg(),
	}

	expander := ExpanderGraph_I{
		Config_: sdr.ExpanderGraphCfg(),
	}

	nodeSize := UInt(sdr.NodeSize())
	challenges := sdr.GenerateOfflineChallenges(randomSeed, int(sdr.Challenges()))

	var challengeProofs []OfflineSDRChallengeProof

	for c := range challenges {
		challengeProofs = append(challengeProofs, CreateChallengeProof(&drg, &expander, replicaID, UInt(c), nodeSize, aux))
	}

	return sdr.CreateCircuitProof(challengeProofs, aux)
}

func CreateChallengeProof(drg *DRG_I, expander *ExpanderGraph_I, replicaID []byte, challenge UInt, nodeSize UInt, aux sector.ProofAuxTmp) (proof OfflineSDRChallengeProof) {
	var columnElements []UInt
	columnElements = append(columnElements, challenge)
	columnElements = append(columnElements, drg.Parents(challenge)...)
	columnElements = append(columnElements, expander.Parents(challenge)...)

	var columnProofs []SDRColumnProof
	for c := range columnElements {
		columnProof := CreateColumnProof(UInt(c), nodeSize, aux)
		columnProofs = append(columnProofs, columnProof)
	}

	dataProof := createInclusionProof(aux.Data()[challenge*nodeSize:(challenge+1)*nodeSize], aux.Data())
	replicaProof := createInclusionProof(aux.Replica()[challenge*nodeSize:(challenge+1)*nodeSize], aux.Data())

	proof = OfflineSDRChallengeProof{
		DataProof:    dataProof,
		ColumnProofs: columnProofs,
		ReplicaProof: replicaProof,
	}

	return proof
}

func CreateColumnProof(c UInt, nodeSize UInt, aux sector.ProofAuxTmp) (columnProof SDRColumnProof) {
	commC := aux.PersistentAux().CommC()
	layers := aux.KeyLayers()
	var column []byte

	for i := 0; i < len(layers); i++ {
		column = append(column, layers[i][c*nodeSize:(c+1)*nodeSize]...)
	}

	leaf := hashColumn(column)
	columnProof = SDRColumnProof(createInclusionProof(leaf, commC))

	return columnProof
}

func createInclusionProof(leaf []byte, root []byte) InclusionProof {
	panic("TODO")
}

type OfflineSDRChallengeProof struct {
	CommRLast sector.Commitment
	CommC     sector.Commitment

	// TODO: these proofs need to depend on hash function.
	DataProof    InclusionProof // Blake2s
	ColumnProofs []SDRColumnProof
	ReplicaProof InclusionProof // Pedersen

}

type InclusionProof struct{}
type SDRColumnProof InclusionProof

func (sdr *StackedDRG_I) CreateCircuitProof(challengeProofs []OfflineSDRChallengeProof, aux sector.ProofAuxTmp) sector.SealProof {
	panic("TODO")
}

func (sdr *StackedDRG_I) GenerateOfflineChallenges(randomSeed sector.SealRandomness, challenges int) []UInt {
	panic("TODO")
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
/// Generic Hashing and Merkle Tree generation

/// Binary hash compression.
// RepCompress<T>
func RepCompress_T(left []byte, right []byte) util.T {
	return util.T{}
}

// RepCompress<PedersenHash>
func RepCompress_PedersenHash(left []byte, right []byte) PedersenHash {
	return PedersenHash{}
}

// RepCompress<Blake2sHash>
func RepCompress_Blake2sHash(left []byte, right []byte) Blake2sHash {
	return Blake2sHash{}
}

////////////////////////////////////////////////////////////////////////////////

/// Digest
// WideRepCompress<T>
func WideRepCompress_T(data []byte) util.T {
	return util.T{}
}

// RepCompress<PedersenHash>
func WideRepCompress_PedersenHash(data []byte) PedersenHash {
	return PedersenHash{}
}

// RepCompress<Blake2sHash>
func WideRepCompress_Blake2sHash(data []byte) Blake2sHash {
	return Blake2sHash{}
}

////////////////////////////////////////////////////////////////////////////////

func DigestSize_T() int {
	panic("Unspecialized")
}

func DigestSize_PedersenHash() int {
	return 32
}

func DigestSize_Blake2sHash() int {
	return 32
}

////////////////////////////////////////////////////////////////////////////////
/// Binary Merkle-tree generation

// RepHash<T>
func RepHash_T(data []byte) (util.T, file.Path) {
	// Plan: define this in terms of RepCompress_T, then copy-paste changes into T-specific specializations, for now.

	// Nodes are always the digest size so data cannot be compressed to digest for storage.
	nodeSize := DigestSize_T()

	// TODO: Fail if len(dat) is not a power of 2 and a multiple of the node size.

	rows := [][]byte{data}

	for row := []byte{}; len(row) > nodeSize; {
		for i := 0; i < len(data); i += 2 * nodeSize {
			left := data[i : i+nodeSize]
			right := data[i+nodeSize : i+2*nodeSize]
			hashed := RepCompress_T(left, right)

			row = append(row, asBytes(hashed)...)
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

// RepHash<PedersenHash>
func RepHash_PedersenHash(data []byte) (PedersenHash, file.Path) {
	return PedersenHash{}, file.Path("") // FIXME
}

//  RepHash<Blake2sHash>
func RepHash_Blake2sHash(data []byte) (Blake2sHash, file.Path) {
	return []byte{}, file.Path("") // FIXME
}

////////////////////////////////////////////////////////////////////////////////

func UnsealedSectorCID(h Blake2sHash) sector.UnsealedSectorCID {
	panic("not implemented -- re-arrange bits")
}

func SealedSectorCID(h PedersenHash) sector.SealedSectorCID {
	panic("not implemented -- re-arrange bits")
}

// Compute CommP or CommD.
func ComputeUnsealedSectorCID(data []byte) sector.UnsealedSectorCID {
	// TODO: check that len(data) > minimum piece size and is a power of 2.
	hash, _ := RepHash_Blake2sHash(data)
	return UnsealedSectorCID(hash)
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

func asBytes(t util.T) []byte {
	panic("Unimplemented for T")

	return []byte{}
}

func fromBytes_T(_ interface{}) util.T {
	panic("Unimplemented for T")
	return util.T{}
}
