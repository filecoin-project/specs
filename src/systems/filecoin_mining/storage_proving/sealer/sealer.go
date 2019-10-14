package sealer

import filproofs "github.com/filecoin-project/specs/libraries/filcrypto/filproofs"
import file "github.com/filecoin-project/specs/systems/filecoin_files/file"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
import . "github.com/filecoin-project/specs/util"
import "math/big"

func (s *SectorSealer_I) SealSector(si SealInputs) *SectorSealer_SealSector_FunRet_I {
	sid := si.SectorID()

	commD := sector.UnsealedSectorCID(s.ComputeDataCommitment(si.UnsealedPath()).As_commD())

	buf := make(Bytes, si.SealCfg().SectorSize())
	f := file.FromPath(si.SealedPath())
	length, _ := f.Read(buf)

	// TODO: How do we meant to handle errors in implementation methods? This could get tedious fast.

	if UInt(length) != UInt(si.SealCfg().SectorSize()) {
		panic("Sector file is wrong size.")
	}

	return &SectorSealer_SealSector_FunRet_I{
		rawValue: Seal(sid, si.RandomSeed(), commD, buf),
	}
}

func (s *SectorSealer_I) VerifySeal(sv sector.SealVerifyInfo) *SectorSealer_VerifySeal_FunRet_I {
	return &SectorSealer_VerifySeal_FunRet_I{}
}

func (s *SectorSealer_I) ComputeDataCommitment(unsealedPath file.Path) *SectorSealer_ComputeDataCommitment_FunRet_I {
	return &SectorSealer_ComputeDataCommitment_FunRet_I{}
}

func ComputeReplicaID(sid sector.SectorID, commD sector.UnsealedSectorCID, seed sector.SealRandomSeed) *SectorSealer_ComputeReplicaID_FunRet_I {

	_, _ = sid.MinerID(), (sid.Number())

	return &SectorSealer_ComputeReplicaID_FunRet_I{}
}

// type SealOutputs struct {
//     SealInfo  sector.SealVerifyInfo
//     ProofAux  sector.ProofAux
// }

// type SealVerifyInfo struct {
//     SectorID
//     OnChain OnChainSealVerifyInfo
// }

func SDRParams() *filproofs.StackedDRG_I  {
	return &filproofs.StackedDRG_I { }
}

func Seal(sid sector.SectorID, randomSeed sector.SealRandomSeed, commD sector.UnsealedSectorCID, data Bytes) *SealOutputs_I {
	replicaID := ComputeReplicaID(sid, commD, randomSeed).As_replicaID()

	params := SDRParams();

	nodeSize := int(params.NodeSize().Size());
	nodes := len(data) / nodeSize;
	curveModulus := params.Curve().Modulus();
	layers := int(params.Layers().Layers());
	key := generateSDRKey(replicaID, nodes, layers, curveModulus);
	
	replica := encodeData(data, key, nodeSize, curveModulus);
	
	_ = replica
	return &SealOutputs_I{}
}

func generateSDRKey(replicaID Bytes, nodes int, layers int, modulus UInt) Bytes {
	return []byte{}
}

func encodeData(data Bytes, key Bytes, nodeSize int, modulus UInt) Bytes {
	bigMod := big.NewInt(int64(modulus));
	
	if len(data) != len(key) {
		panic("Key and data must be same length.")
	}

	encoded := make(Bytes, len(data))
	for i := 0; i < len(data); i += nodeSize {
		copy(encoded[i:i+nodeSize], encodeNode(data[i:i+nodeSize], key[i:i+nodeSize], bigMod, nodeSize));
	}
	
	return encoded
}

func encodeNode(data Bytes, key Bytes, modulus *big.Int, nodeSize int) Bytes {

	// TODO: Allow this to vary by algorithm variant.
	return addEncode(data, key, modulus, nodeSize);
}


func reverse(bytes []byte) {
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
}

func addEncode(data Bytes, key Bytes, modulus *big.Int, nodeSize int) Bytes {
	// FIXME: Check correct endianness.
	sum := new(big.Int);
	reverse(data); // Reverse for little-endian 
	reverse(key);  // Reverse for little-endian 

	d := new(big.Int).SetBytes(data); // Big-endian
	k := new(big.Int).SetBytes(key);  // Big-endian
	
	sum = sum.Add(d, k);
	
	result := new(big.Int);
	resultBytes := result.Mod(sum, modulus).Bytes()[0:nodeSize]; // Big-endian
	reverse(resultBytes); // Reverse for little-endian 

	return resultBytes; 
}
