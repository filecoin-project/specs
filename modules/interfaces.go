package filmodules

// dependencies

// Multiformats
type mf_Multihash []byte
type mf_Multiaddr []byte

// IPLD
type ipld_CID string
type ipld_Path interface {
  CID() CID
  String() string
  Components() []string
}
type ipld_Selector interface {
  CID() CID
  String() string
}

type ipld_Block []byte
type ipld_Object interface {
  CID() ipld_CID
}

type ipld_Store interface {
  Get(cid ipld_CID) ipld_Object
  Put(ipld_Object) error
}

type ipld_PathGetter interface {
  GetPath(path ipld_Path) ipld_Object
}

type ipld_SelectorGetter interface {
  GetSelector(sel ipld_Selector) ipld_Object
}

// IPLD Immutable ShardedMap (today, a HAMT)
type ipld_ShardedMap interface {
  Get(string) ([]byte, error)
  GetCID(string) (ipld.CID, error)
  GetObject(string) (ipld.Object, error)
}

// IPLD MutableShardedMap
type ipld_MutableShardedMap interface {
  Put(string, []byte) error
  PutCID(string, ipld_CID) error
  PutInline(string, ipld_Object) error
}

type ipld_ShardedArray interface {
}

// BigInt
type BigInt interface{}

// pc - things dealing with libp2p crypto algos (keys, signatures)

// Keys
// (ic -> pc, ipfs crypto -> libp2p crypto)
type pc_SigAlgo interface {
  Sign([]byte, SigPrivateKey) (Signature, error)
  Verify([]byte, SigPublicKey, Signature) (ok bool, err error)
  Recover([]byte, Signature) (SigPublicKey, error)
}

type pc_AggregateSignature interface {
  pc_Signature
}

type pc_UnbiasedRandomness []byte

// fil_proofs - things dealing with the cryptographic proofs

// Proofs
type fil_proofs_CRH []byte
type fil_proofs_Commitment []byte
type fil_proofs_Proof []byte

type fil_proofs_PoRepProof fil_proofs_Proof
type fil_proofs_PoStProof fil_proofs_Proof
type fil_proofs_SealProof fil_proofs_Proof
type fil_proofs_PieceInclusionProof fil_proofs_Proof

type fil_proofs_CommD fil_proofs_Commitment
type fil_proofs_CommP fil_proofs_Commitment
type fil_proofs_CommR fil_proofs_Commitment
type fil_proofs_CommRStar fil_proofs_Commitment

type fil_proofs_PoRepAlgo interface {
}

type fil_proofs_PoStAlgo interface {
}

type fil_proofs_SealAlgo interface {
}

// fil_sectors - things dealing with filecoin sectors

type fil_sectors_SectorID struct {
  MinerActorID fil_Address
  SectorNumber int
}

type fil_sectors_SectorInfo struct {
  SectorID    fil_sectors_SectorID
  CommD       fil_proofs_CommD
  CommR       fil_proofs_CommR
  CommRStar   fil_proofs_CommRStar
  SealProof   fil_proofs_SealProof
  PayloadSize int
  NumPieces   int
}

type fil_sectors_FaultSet struct { // onchain
  // Index is a block height offset from the start of the miner's proving period,
  // The index is used to make the representation of the FaultSet more compact.
  Index    UVarint
  BitField BitField
}

type fil_sectors_Sector interface {
  SectorID() SectorID
  Info() SectorInfo
  Bytes() io.Reader
}

type fil_proofs_SealCfg struct {
  Partitions int
  SectorSize int
  SealAlgo   int
}

type fil_proofs_SectorPath string

type fil_proofs_Sealer interface {
  Seal(cfg fil_proofs_SealCfg, usrc, sdst SectorPath, sid fil_sectors_SectorID, r UnbiasedRand)
  VerifySeal(info fil_sectors_SectorInfo, r UnbiasedRand)
  Unseal(cfg fil_proofs_SealCfg, ssrc, udst SectorPath, sid fil_sectors_SectorID)
  UnsealPartial(cfg fil_proofs_SealCfg, ssrc, udst SectorPath, sid fil_sectors_SectorID, startOffset int, length int)
}

// fil_Address
type fil_Address interface {
  Network() Varint
  Protocol() Varint
  Payload() []byte
  Checksum() []byte
  Bytes() []byte
}
