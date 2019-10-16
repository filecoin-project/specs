package sealer

import "math/big"
import "encoding/binary"
import . "github.com/filecoin-project/specs/util"

import filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
import file "github.com/filecoin-project/specs/systems/filecoin_files/file"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"

func (s *SectorSealer_I) SealSector(si SealInputs) *SectorSealer_SealSector_FunRet_I {
	sid := si.SectorID()
	commD := sector.UnsealedSectorCID(s.ComputeDataCommitment(si.UnsealedPath()).As_commD())

	buf := make(Bytes, si.SealCfg().SectorSize())
	f := file.FromPath(si.SealedPath())
	length, _ := f.Read(buf)

	if UInt(length) != UInt(si.SealCfg().SectorSize()) {
		return &SectorSealer_SealSector_FunRet_I{
			rawValue: "Sector file is wrong size",
			which:    SectorSealer_SealSector_FunRet_Case_err,
		}
	}

	return &SectorSealer_SealSector_FunRet_I{
		rawValue: Seal(sid, commD, buf),
		which:    SectorSealer_SealSector_FunRet_Case_so,
	}
}

func (s *SectorSealer_I) CreateSealProof(si CreateSealProofInputs) *SectorSealer_CreateSealProof_FunRet_I {
	sid := si.SectorID()
	commD := si.SealOutputs().SealInfo().OnChain().UnsealedCID()
	layers := si.SealOutputs().ProofAuxTmp().Layers()

	return &SectorSealer_CreateSealProof_FunRet_I{
		rawValue: CreateSealProof(sid, si.RandomSeed(), commD, layers[len(layers)-1]),
	}
}

func (s *SectorSealer_I) VerifySeal(sv sector.SealVerifyInfo) *SectorSealer_VerifySeal_FunRet_I {
	return &SectorSealer_VerifySeal_FunRet_I{}
}

func (s *SectorSealer_I) ComputeDataCommitment(unsealedPath file.Path) *SectorSealer_ComputeDataCommitment_FunRet_I {
	// TODO: Generate merkle tree using appropriate hash.
	return &SectorSealer_ComputeDataCommitment_FunRet_I{}
}

func ComputeReplicaID(sid sector.SectorID, commD sector.UnsealedSectorCID) *SectorSealer_ComputeReplicaID_FunRet_I {

	_, _ = sid.MinerID(), (sid.Number())

	// FIXME: Implement
	return &SectorSealer_ComputeReplicaID_FunRet_I{}
}

func UnsealedSectorCID(h filproofs.Blake2sHash) sector.UnsealedSectorCID {
	panic("not implemented -- re-arrange bits")
}

func SealedSectorCID(h filproofs.PedersenHash) sector.SealedSectorCID {
	panic("not implemented -- re-arrange bits")
}

func SDRParams() *filproofs.StackedDRG_I {
	return &filproofs.StackedDRG_I{}
}

func Seal(sid sector.SectorID, commD sector.UnsealedSectorCID, data Bytes) *SealOutputs_I {
	replicaID := ComputeReplicaID(sid, commD).As_replicaID()

	params := SDRParams()

	drg := filproofs.DRG_I{}                // FIXME: Derive from params
	expander := filproofs.ExpanderGraph_I{} // FIXME: Derive from params
	nodeSize := int(params.NodeSize().Size())
	nodes := len(data) / nodeSize
	curveModulus := params.Curve().FieldModulus()
	layers := int(params.Layers().Layers())
	keyLayers := generateSDRKeyLayers(&drg, &expander, replicaID, nodes, layers, nodeSize, curveModulus)
	key := keyLayers[len(keyLayers)-1]

	replica := encodeData(data, key, nodeSize, curveModulus)

	var cachedMerkleTreePath file.Path // FIXME: get this

	commR, cachedMerkleTreePath := repHash(replica)

	var proof sector.SealProof

	return &SealOutputs_I{
		SealInfo_: &sector.SealVerifyInfo_I{
			SectorID_: sid,
			OnChain_: &sector.OnChainSealVerifyInfo_I{
				SealedCID_:   SealedSectorCID(commR),
				UnsealedCID_: commD,
				Proof_:       proof,
			},
		},
		ProofAuxTmp_: &sector.ProofAuxTmp_I{
			CommRLast_:            sector.Commitment{},
			CommC_:                sector.Commitment{},
			CachedMerkleTreePath_: cachedMerkleTreePath,
		},
	}
}

func CreateSealProof(sid sector.SectorID, randomSeed sector.SealRandomSeed, commD sector.UnsealedSectorCID, replica Bytes) *CreateSealProofOutputs_I {
	//replicaID := ComputeReplicaID(sid, commD).As_replicaID()

	//params := SDRParams()

	//	drg := filproofs.DRG_I{}                // FIXME: Derive from params
	//expander := filproofs.ExpanderGraph_I{} // FIXME: Derive from params
	//nodeSize := int(params.NodeSize().Size())
	//	nodes := len(replica) / nodeSize

	var cachedMerkleTreePath file.Path // FIXME: get this

	commR, cachedMerkleTreePath := repHash(replica)

	var proof sector.SealProof

	return &CreateSealProofOutputs_I{
		SealInfo_: &sector.SealVerifyInfo_I{
			SectorID_: sid,
			OnChain_: &sector.OnChainSealVerifyInfo_I{
				SealedCID_:   SealedSectorCID(commR),
				UnsealedCID_: commD,
				RandomSeed_:  randomSeed,
				Proof_:       proof,
			},
		},
		ProofAux_: &sector.ProofAux_I{
			CommRLast_:            sector.Commitment{},
			CommC_:                sector.Commitment{},
			CachedMerkleTreePath_: cachedMerkleTreePath,
		},
	}
}

func repHash(data Bytes) (filproofs.PedersenHash, file.Path) {
	return Bytes{}, file.Path("") // FIXME
}

func generateSDRKeyLayers(drg *filproofs.DRG_I, expander *filproofs.ExpanderGraph_I, replicaID Bytes, nodes int, layers int, nodeSize int, modulus UInt) []Bytes {
	keyLayers := make([]Bytes, layers)
	var prevLayer Bytes

	for i := 0; i <= layers; i++ {
		keyLayers[i] = labelLayer(drg, expander, replicaID, nodes, nodeSize, prevLayer)
	}
	return keyLayers
}

func labelLayer(drg *filproofs.DRG_I, expander *filproofs.ExpanderGraph_I, replicaID Bytes, nodeSize int, nodes int, prevLayer Bytes) Bytes {
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
