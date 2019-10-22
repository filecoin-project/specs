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
type PieceInfo = sector.PieceInfo

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

func (sdr *StackedDRG_I) Seal(sid sector.SectorID, data []byte, randomness sector.SealRandomness) SealSetupArtifacts {
	commD, commDTreePath := ComputeDataCommitment(data)
	sealSeed := ComputeSealSeed(sid, commD, randomness)

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

	keyLayers := generateSDRKeyLayers(&drg, &expander, sealSeed, nodes, layers, nodeSize, curveModulus)
	key := keyLayers[len(keyLayers)-1]

	replica := encodeData(data, key, nodeSize, &curveModulus)

	commRLast, commRLastTreePath := RepHash_PedersenHash(replica)
	commC, commCTreePath := computeCommC(keyLayers, nodeSize)
	commR := RepCompress_PedersenHash(commC, commRLast)

	result := SealSetupArtifacts_I{
		CommD_:             sector.Commitment(commC),
		CommR_:             SealedSectorCID(commR),
		CommC_:             sector.Commitment(commC),
		CommRLast_:         sector.Commitment(commRLast),
		CommDTreePath_:     commDTreePath,
		CommCTreePath_:     commCTreePath,
		CommRLastTreePath_: commRLastTreePath,
		Seed_:              sealSeed,
		KeyLayers_:         keyLayers,
		Replica_:           replica,
	}
	return &result
}

func ComputeSealSeed(sid sector.SectorID, commD sector.Commitment, randomness sector.SealRandomness) sector.SealSeed {
	_, _ = sid.MinerID(), (sid.Number())

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

func generateLabel(sealSeed sector.SealSeed, node int, dependencies []byte) []byte {
	nodeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nodeBytes, uint64(node))

	preimage := append(sealSeed, nodeBytes...)
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
	sealSeed := aux.Seed()

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
		challengeProofs = append(challengeProofs, CreateChallengeProof(&drg, &expander, sealSeed, UInt(c), nodeSize, aux))
	}

	return sdr.CreateOfflineCircuitProof(challengeProofs, aux)
}

func CreateChallengeProof(drg *DRG_I, expander *ExpanderGraph_I, sealSeed sector.SealSeed, challenge UInt, nodeSize UInt, aux sector.ProofAuxTmp) (proof OfflineSDRChallengeProof) {
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
	columnProof = SDRColumnProof{
		ColumnElements: column,
		InclusionProof: createInclusionProof(leaf, commC),
	}

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

type SDRColumnProof struct {
	ColumnElements []byte
	InclusionProof InclusionProof
}

func (sdr *StackedDRG_I) CreateOfflineCircuitProof(challengeProofs []OfflineSDRChallengeProof, aux sector.ProofAuxTmp) sector.SealProof {
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
// Verification

func (sdr *StackedDRG_I) VerifySeal(sv sector.SealVerifyInfo) bool {
	onChain := sv.OnChain()

	sealProof := onChain.Proof()

	pieceInfos := sv.PieceInfos()

	commR := Commitment_SealedSectorCID(sector.SealedSectorCID(onChain.SealedCID()))

	sealCfg := sv.SealCfg()
	util.Assert(sealCfg.SubsectorCount() == 1) // A more sophisticated accounting of CommD for verification purposes will be required when supersectors are considered.
	sectorSize := sealCfg.SectorSize()

	rootPieceInfo := sdr.ComputeRootPieceInfo(pieceInfos)
	rootSize := rootPieceInfo.Size()
	commD := rootPieceInfo.CommP()

	if rootSize != sectorSize {
		return false
	}

	return sdr.VerifyOfflineCircuitProof(commD, commR, sealProof)
}

func (sdr *StackedDRG_I) ComputeRootPieceInfo(pieceInfos []PieceInfo) PieceInfo {
	// Construct root PieceInfo by (shift-reduce) parsing the constituent pieceInfo array.
	// Later pieces must always be joined with equal-sized predecessors to create a new root twice their size.
	// So if a piece is larger than the current root (top of stack), add padding until it is not.
	// If a piece is smaller than the root, let it be the new root (top of stack) until reduced to a new root that can be joined
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
		if peek().Size() == peek2().Size() {
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
		Size_:  left.Size() + right.Size(),
		CommP_: UnsealedSectorCID(RepCompress_Blake2sHash(AsBytes_UnsealedSectorCID(left.CommP()), AsBytes_UnsealedSectorCID(right.CommP()))), // FIXME: make this whole function generic?
	}
}

func (sdr *StackedDRG_I) VerifyOfflineCircuitProof(commD sector.UnsealedSectorCID, commR sector.Commitment, sv sector.SealProof) bool {

	panic("TODO")
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
	result := Blake2sHash{}
	return trimToFr32(result)
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
	// Digest is truncated to 254 bits.
	result := Blake2sHash{}
	return trimToFr32(result)
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

// RepHash<PedersenHash>
func RepHash_PedersenHash(data []byte) (PedersenHash, file.Path) {
	return PedersenHash{}, file.Path("") // FIXME
}

//  RepHash<Blake2sHash>
func RepHash_Blake2sHash(data []byte) (Blake2sHash, file.Path) {
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

func UnsealedSectorCID(h Blake2sHash) sector.UnsealedSectorCID {
	panic("not implemented -- re-arrange bits")
}

func SealedSectorCID(h PedersenHash) sector.SealedSectorCID {
	panic("not implemented -- re-arrange bits")
}

func Commitment_UnsealedSectorCID(cid sector.UnsealedSectorCID) sector.Commitment {
	panic("not implemented -- re-arrange bits")
}

func Commitment_SealedSectorCID(cid sector.SealedSectorCID) sector.Commitment {
	panic("not implemented -- re-arrange bits")
}

func ComputeDataCommitment(data []byte) ([]byte, file.Path) {
	// TODO: make hash parameterizable
	return RepHash_Blake2sHash(data)
}

// Compute CommP or CommD.
func ComputeUnsealedSectorCID(data []byte) (sector.UnsealedSectorCID, file.Path) {
	// TODO: check that len(data) > minimum piece size and is a power of 2.
	hash, treePath := RepHash_Blake2sHash(data)
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

func AsBytes_T(t util.T) []byte {
	panic("Unimplemented for T")

	return []byte{}
}

func AsBytes_UnsealedSectorCID(cid sector.UnsealedSectorCID) []byte {
	panic("Unimplemented for T")

	return []byte{}
}

func AsBytes_SealedSectorCID(CID sector.SealedSectorCID) []byte {
	panic("Unimplemented for T")

	return []byte{}
}

func fromBytes_T(_ interface{}) util.T {
	panic("Unimplemented for T")
	return util.T{}
}
