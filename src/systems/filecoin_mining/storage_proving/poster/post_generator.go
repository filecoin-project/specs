package poster

// See "Proof-of-Spacetime Parameters" Section
const POST_PROVING_PERIOD = uint(5760)
const POST_CHALLENGE_DEADLINE = uint(480)

// // this must  GetRandFromBlock(self.ProvingPeriodEnd - POST_CHALLENGE_TIME)
//   GetChallenge(minerActor &StorageMinerActor, currBlock) Challenge

//   GeneratePoSt(challenge Challenge, sectors [Sector]) PoStProof {
//       sectorsMetadata := sectors.map(func(sector) { SectorStorage.GetMetadata(sector.CommR) });

//       // Question: Should we pass metadata into FilProofs so it can interact with SectorStore directly?
//       // Like this:
//       PoStReponse := SectorStorageSubsystem.GeneratePoSt(sectorSize, challenge, faults, sectorsMetatada);

//       // Question: Or should we resolve + manifest trees here and pass them in?
//       // Like this:
//       trees := sectorsMetadata.map(func(md) { SectorStorage.GetMerkleTree(md.MerkleTreePath) });
//       // Done this way, we redundantly pass the tree paths in the metadata. At first thought, the other way
//       // seems cleaner.
//       PoStReponse := SectorStorageSubsystem.GeneratePoSt(sectorSize, challenge, faults, sectorsMetadata, trees);
//   }
