package filproofs

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/rand"

	"encoding/binary"
	big "math/big"

	abi "github.com/filecoin-project/specs/actors/abi"
	file "github.com/filecoin-project/specs/systems/filecoin_files/file"
	piece "github.com/filecoin-project/specs/systems/filecoin_files/piece"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	sector_index "github.com/filecoin-project/specs/systems/filecoin_mining/sector_index"
	util "github.com/filecoin-project/specs/util"
	"github.com/ipfs/go-cid"
)

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

const GIB_32 = 32 * 1024 * 1024 * 1024

func PoStCfg(pType sector.PoStType, sectorSize sector.SectorSize, partitions UInt) *sector.PoStCfg_I {
	nodes := UInt(sectorSize / NODE_SIZE)

	var registeredProof sector.RegisteredProof

	switch pType {
	case sector.PoStType_ElectionPoSt:
		registeredProof = sector.RegisteredProof_WinStackedDRG32GiBPoSt
	case sector.PoStType_SurprisePoSt:
		registeredProof = sector.RegisteredProof_WinStackedDRG32GiBPoSt

	}

	return &sector.PoStCfg_I{
		InstanceCfg_: sector.PoStInstanceCfg_Make_PoStCfgV1(&sector.PoStCfgV1_I{
			Type_:               pType,
			Nodes_:              nodes,
			Partitions_:         partitions,
			LeafChallengeCount_: POST_LEAF_CHALLENGE_COUNT,
			ChallengeRangeSize_: POST_CHALLENGE_RANGE_SIZE,
		}),
		ProofInstance_: sector.PROOFS[UInt(registeredProof)],
	}
}

func MakeSealVerifier(cfg sector.SealCfg) *SealVerifier_I {
	return &SealVerifier_I{
		SealCfg_: cfg,
	}
}

func SurprisePoStCfg(sectorSize sector.SectorSize) *sector.PoStCfg_I {
	return PoStCfg(sector.PoStType_SurprisePoSt, sectorSize, SURPRISE_POST_PARTITIONS)
}

func ElectionPoStCfg(sectorSize sector.SectorSize) *sector.PoStCfg_I {
	return PoStCfg(sector.PoStType_ElectionPoSt, sectorSize, ELECTION_POST_PARTITIONS)
}

func MakeElectionPoStVerifier(cfg sector.PoStCfg) *PoStVerifier_I {
	return &PoStVerifier_I{
		ElectionPoStCfg_: cfg,
	}
}

func MakeSurprisePoStVerifier(cfg sector.PoStCfg) *PoStVerifier_I {
	return &PoStVerifier_I{
		SurprisePoStCfg_: cfg,
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

func getProverID(minerID abi.ActorID) []byte {
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

func (sv *SealVerifier_I) VerifySeal(svi sector.SealVerifyInfo) bool {
	switch svi.SealCfg().ProofInstance().ID() {
	case sector.RegisteredProof_WinStackedDRG32GiBSeal:
		{
			sdr := WinSDRParams(svi.SealCfg())

			return sdr.VerifySeal(svi)
		}
	case sector.RegisteredProof_StackedDRG32GiBSeal:
		{
			panic("TODO")
		}
	}

	return false
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

	// FIXME: make this whole function generic?
	// Note: cid.Bytes() isn't actually the payload data that we want input to the binary hash function, for more
	// information see discussion: https://filecoinproject.slack.com/archives/CHMNDCK9P/p1578629688082700
	sectorPieceCID, err := cid.Cast(BinaryHash_SHA256Hash(cid.Cid(left.PieceCID()).Bytes(), cid.Cid(right.PieceCID()).Bytes()))
	util.Assert(err == nil)

	return &sector.PieceInfo_I{
		Size_:     left.Size() + right.Size(),
		PieceCID_: piece.PieceCID(sectorPieceCID),
	}
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

func generateCandidate(algorithm sector.ProofAlgorithm, cfg sector.PoStCfg, randomness sector.PoStRandomness, aux sector.PersistentProofAux, sectorID sector.SectorID, sectorChallengeIndex UInt) sector.PoStCandidate {
	var candidate sector.PoStCandidate
	switch algorithm {
	case sector.ProofAlgorithm_StackedDRGSeal:
		panic("TODO")
	case sector.ProofAlgorithm_WinStackedDRGSeal:
		sdr := WinStackedDRG_I{}
		candidate = sdr._generateCandidate(cfg, randomness, aux, sectorID, sectorChallengeIndex)
	}
	return candidate
}

func computePartialTicket(randomness sector.PoStRandomness, sectorID sector.SectorID, data []byte) sector.PartialTicket {
	preimage := randomness
	preimage = append(preimage, getProverID(sectorID.Miner())...)
	preimage = append(preimage, littleEndianBytesFromUInt(UInt(sectorID.Number()), 8)...)
	preimage = append(preimage, data...)
	partialTicket := sector.PartialTicket(HashBytes_PedersenHash(preimage))

	return partialTicket
}

type PoStCandidatesMap map[sector.ProofAlgorithm][]sector.PoStCandidate

func CreatePoStProof(cfg sector.PoStCfg, privateCandidateProofs []PrivatePostCandidateProof, challengeSeed sector.PoStRandomness) []sector.PoStProof {
	var proofsMap map[sector.ProofAlgorithm][]PrivatePostCandidateProof

	for _, proof := range privateCandidateProofs {
		algorithm := proof.Algorithm()
		proofsMap[algorithm] = append(proofsMap[algorithm], proof)
	}

	var circuitProofs []sector.PoStProof
	for algorithm, proofs := range proofsMap {
		privateProof := createPrivatePoStProof(algorithm, proofs, challengeSeed)
		circuitProof := createPoStCircuitProof(cfg, algorithm, privateProof)
		circuitProofs = append(circuitProofs, circuitProof)
	}

	return circuitProofs
}

type PrivatePoStProof struct {
	Algorithm       sector.ProofAlgorithm
	ChallengeSeed   sector.PoStRandomness
	CandidateProofs []PrivatePostCandidateProof
}

func createPrivatePoStProof(algorithm sector.ProofAlgorithm, candidateProofs []PrivatePostCandidateProof, challengeSeed sector.PoStRandomness) PrivatePoStProof {
	// TODO: Verify that all candidateProofs share algorithm.
	return PrivatePoStProof{
		Algorithm:       algorithm,
		ChallengeSeed:   challengeSeed,
		CandidateProofs: candidateProofs,
	}
}

type InternalPrivateCandidateProof struct {
	InclusionProofs []InclusionProof
}

// This exists because we need to pass private proofs out of filproofs for winner selection.
// Actually implementing it would (will?) be tedious, since it means doing the same for InclusionProofs.
func (p *InternalPrivateCandidateProof) externalize(algorithm sector.ProofAlgorithm) sector.PrivatePoStCandidateProof {
	return &sector.PrivatePoStCandidateProof_I{
		Algorithm_:    algorithm,
		Externalized_: []byte{}, // Unimplemented.
	}
}

// This is the inverse of InternalPrivateCandidateProof.externalize and equally tedious.
func newInternalPrivateProof(externalPrivateProof PrivatePostCandidateProof) InternalPrivateCandidateProof {
	return InternalPrivateCandidateProof{}
}

func createPoStCircuitProof(postCfg sector.PoStCfg, algorithm sector.ProofAlgorithm, privateProof PrivatePoStProof) (proof sector.PoStProof) {
	switch algorithm {
	case sector.ProofAlgorithm_WinStackedDRGSeal:
		sdr := WinStackedDRG_I{}
		proof = sdr._createPoStCircuitProof(postCfg, privateProof)
	}
	return proof
}

func (pv *PoStVerifier_I) _verifyPoStProof(sv sector.PoStVerifyInfo) bool {
	// commT := sv.CommT()
	// candidates := sv.Candidates()
	// randomness := sv.Randomness()
	// postProofs := sv.OnChain.Proofs()

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
