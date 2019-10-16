package filproofs

import "math/big"
import "encoding/binary"
import . "github.com/filecoin-project/specs/util"

import file "github.com/filecoin-project/specs/systems/filecoin_files/file"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"

type Blake2sHash Bytes
type PedersenHash Bytes

func (drg *DRG_I) Parents(layer Bytes, node UInt) []UInt {
	return []UInt{} // FIXME
}

func (exp *ExpanderGraph_I) Parents(layer Bytes, node UInt) []UInt {
	return []UInt{} // FIXME
}

func (sdr *StackedDRG_I) Seal(sid sector.SectorID, commD sector.UnsealedSectorCID, data Bytes) SealSetupArtifacts {
	replicaID := ComputeReplicaID(sid, commD)

	// FIXME: Derive these from sdr.
	drg := DRG_I{}
	expander := ExpanderGraph_I{}

	nodeSize := int(sdr.NodeSize().Size())
	nodes := len(data) / nodeSize
	curveModulus := sdr.Curve().FieldModulus()
	layers := int(sdr.Layers().Layers())
	keyLayers := generateSDRKeyLayers(&drg, &expander, replicaID, nodes, layers, nodeSize, curveModulus)
	key := keyLayers[len(keyLayers)-1]

	replica := encodeData(data, key, nodeSize, curveModulus)

	var cachedMerkleTreePath file.Path // FIXME: get this

	commR, cachedMerkleTreePath := RepHash_PedersenHash(replica)

	result := SealSetupArtifacts_I{
		CommR_:     SealedSectorCID(commR),
		CommC_:     sector.Commitment{},
		CommRLast_: sector.Commitment{},
		TreePath_:  cachedMerkleTreePath,
	}
	return &result
}

func (sdr *StackedDRG_I) CreateSealProof(randomSeed sector.SealRandomSeed, aux sector.ProofAuxTmp) sector.SealProof {
	panic("TODO")
}

func ComputeReplicaID(sid sector.SectorID, commD sector.UnsealedSectorCID) Bytes {
	_, _ = sid.MinerID(), (sid.Number())

	// FIXME: Implement
	return Bytes{}
}

func SDRParams() *StackedDRG_I {
	return &StackedDRG_I{}
}

func generateSDRKeyLayers(drg *DRG_I, expander *ExpanderGraph_I, replicaID Bytes, nodes int, layers int, nodeSize int, modulus UInt) []Bytes {
	keyLayers := make([]Bytes, layers)
	var prevLayer Bytes

	for i := 0; i <= layers; i++ {
		keyLayers[i] = labelLayer(drg, expander, replicaID, nodes, nodeSize, prevLayer)
	}
	return keyLayers
}

func encodeData(data Bytes, key Bytes, nodeSize int, modulus UInt) Bytes {
	bigMod := big.NewInt(int64(modulus))

	if len(data) != len(key) {
		panic("Key and data must be same length.")
	}

	encoded := make(Bytes, len(data))
	for i := 0; i < len(data); i += nodeSize {
		copy(encoded[i:i+nodeSize], encodeNode(data[i:i+nodeSize], key[i:i+nodeSize], bigMod, nodeSize))
	}

	return encoded
}

func labelLayer(drg *DRG_I, expander *ExpanderGraph_I, replicaID Bytes, nodeSize int, nodes int, prevLayer Bytes) Bytes {
	size := nodes * nodeSize
	labels := make(Bytes, size)

	for i := 0; i < nodes; i++ {
		var dependencies Bytes

		// The first node of every layer has no DRG Parents.
		if i > 0 {
			for parent := range drg.Parents(labels, UInt(i)) {
				start := parent * nodeSize
				dependencies = append(dependencies, labels[start:start+nodeSize]...)
			}
		}

		// The first layer has no expander parents.
		if prevLayer != nil {
			for parent := range expander.Parents(labels, UInt(i)) {
				start := parent * nodeSize
				dependencies = append(dependencies, labels[start:start+nodeSize]...)
			}
		}

		label := generateLabel(replicaID, i, dependencies)
		labels = append(labels, label...)
	}

	return labels
}

func generateLabel(replicaID Bytes, node int, dependencies Bytes) Bytes {
	nodeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nodeBytes, uint64(node))

	preimage := append(replicaID, nodeBytes...)
	preimage = append(preimage, dependencies...)

	return KDF(preimage)
}

// KDF is a key-derivation functions. In SDR, the derived key is used to generate labels directly, without encoding any data.
func KDF(elements Bytes) Bytes {
	return elements // FIXME: Do something.
}

func encodeNode(data Bytes, key Bytes, modulus *big.Int, nodeSize int) Bytes {
	// TODO: Allow this to vary by algorithm variant.
	return addEncode(data, key, modulus, nodeSize)
}

func addEncode(data Bytes, key Bytes, modulus *big.Int, nodeSize int) Bytes {

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
func RepCompress_T(left Bytes, right Bytes) T {
	return T{}
}

// RepCompress<PedersenHash>
func RepCompress_PedersenHash(left Bytes, right Bytes) PedersenHash {
	return PedersenHash{}
}

// RepCompress<Blake2sHash>
func RepCompress_Blake2sHash(left Bytes, right Bytes) Blake2sHash {
	return Blake2sHash{}
}

////////////////////////////////////////////////////////////////////////////////

/// Binary Merkle-tree generation
// RepHash<T>
func RepHash_T(data Bytes) (T, file.Path) {
	// Plan: define this in terms of RepCompress_T, then copy-paste changes into T-specific specializations, for now.
	return T{}, file.Path("") // FIXME
}

// RepHash<PedersenHash>
func RepHash_PedersenHash(data Bytes) (PedersenHash, file.Path) {
	return PedersenHash{}, file.Path("") // FIXME
}

//  RepHash<Blake2sHash>
func RepHash_Blake2sHash(data Bytes) (Blake2sHash, file.Path) {
	return Bytes{}, file.Path("") // FIXME
}

////////////////////////////////////////////////////////////////////////////////

func UnsealedSectorCID(h Blake2sHash) sector.UnsealedSectorCID {
	panic("not implemented -- re-arrange bits")
}

func SealedSectorCID(h PedersenHash) sector.SealedSectorCID {
	panic("not implemented -- re-arrange bits")
}

// Utilities

func reverse(bytes []byte) {
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
}

func bigIntFromLittleEndianBytes(bytes Bytes) *big.Int {
	reverse(bytes)
	return new(big.Int).SetBytes(bytes)
}

func littleEndianBytesFromBigInt(z *big.Int, size int) Bytes {
	bytes := z.Bytes()[0:size]
	reverse(bytes)

	return bytes
}
