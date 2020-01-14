package filproofs

import (
	"bytes"
	big "math/big"

	file "github.com/filecoin-project/specs/systems/filecoin_files/file"
	sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
	util "github.com/filecoin-project/specs/util"
	"github.com/ipfs/go-cid"
)

func WinSDRParams(c sector.SealInstanceCfg) *WinStackedDRG_I {
	cfg := c.As_WinStackedDRGCfgV1()
	// TODO: Bridge constants with orient model.
	const LAYERS = 10
	const OFFLINE_CHALLENGES = 6666
	const OFFLINE_WINDOW_CHALLENGES = 1111
	const FEISTEL_ROUNDS = 3
	var FEISTEL_KEYS = [FEISTEL_ROUNDS]UInt{1, 2, 3}
	var FIELD_MODULUS = new(big.Int)
	// https://github.com/zkcrypto/pairing/blob/master/src/bls12_381/fr.rs#L4
	FIELD_MODULUS.SetString("52435875175126190479447740508185965837690552500527637822603658699938581184513", 10)

	nodes := UInt(cfg.SectorSize() / NODE_SIZE)

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
		Cfg_: c,
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

func (sdr *WinStackedDRG_I) Seal(registeredProof sector.RegisteredProof, sid sector.SectorID, data []byte, randomness sector.SealRandomness) SealSetupArtifacts {

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
		CommD_:             cid.Cid(commD).Bytes(),
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
	windowSize := int(uint64(sdr.Cfg().As_WinStackedDRGCfgV1().SectorSize()) / UInt(sdr.WindowCount()))

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
	windowSize := int(UInt(sdr.Cfg().As_WinStackedDRGCfgV1().SectorSize()) / UInt(sdr.WindowCount())) // TOOD: Make this a function.
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

func (sdr *WinStackedDRG_I) CreateOfflineCircuitProof(proof PrivateOfflineProof, aux sector.ProofAuxTmp) sector.SealProof {
	// partitions := sdr.Partitions()
	// publicInputs := GeneratePublicInputs()

	panic("TODO")
	var proofBytes []byte
	panic("TODO")

	sealProof := sector.SealProof_I{
		ProofBytes_: proofBytes,
	}

	return &sealProof
}

func (sdr *WinStackedDRG_I) _generateOfflineChallenges(sealSeed sector.SealSeed, randomness sector.InteractiveSealRandomness, wrapperChallengeCount WinStackedDRGChallenges, windowChallengeCount WinStackedDRGWindowChallenges) (windowChallenges []UInt, wrapperChallenges []UInt) {
	wrapperChallenges = generateOfflineChallenges(int(sdr.Nodes()), sealSeed, randomness, int(wrapperChallengeCount))
	windowChallenges = generateOfflineChallenges(int(sdr.WindowDRGCfg().Nodes()), sealSeed, randomness, int(windowChallengeCount))

	return windowChallenges, wrapperChallenges
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

func (sdr *WinStackedDRG_I) _verifyOfflineCircuitProof(commD sector.UnsealedSectorCID, commR sector.SealedSectorCID, sealSeed sector.SealSeed, windowChallenges []UInt, wrapperChallenges []UInt, sv sector.SealProof) bool {
	//publicInputs := GeneratePublicInputs()
	panic("TODO")
}

////////////////////////////////////////////////////////////////////////////////
// PoSt

func (sdr *WinStackedDRG_I) _generateCandidate(postCfg sector.PoStInstanceCfg, randomness sector.PoStRandomness, aux sector.PersistentProofAux, sectorID sector.SectorID, sectorChallengeIndex UInt) sector.PoStCandidate {
	cfg := postCfg.As_PoStCfgV1()

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
		PrivateProof_:   privateProof.externalize(sector.ProofAlgorithm_WinStackedDRGSeal),
		SectorID_:       sectorID,
		ChallengeIndex_: sectorChallengeIndex,
	}
	return &candidate
}

func (sdr *WinStackedDRG_I) VerifyInternalPrivateCandidateProof(postCfg sector.PoStInstanceCfg, p *InternalPrivateCandidateProof, challengeSeed sector.PoStRandomness, candidate sector.PoStCandidate, commRLast Commitment) bool {
	cfg := postCfg.As_PoStCfgV1()
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

func (sdr *WinStackedDRG_I) VerifyPrivatePoStProof(cfg sector.PoStInstanceCfg, privateProof PrivatePoStProof, candidates []sector.PoStCandidate, sectorIDs []sector.SectorID, sectorCommitments sector.SectorCommitments) bool {
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

func (sdr *WinStackedDRG_I) _createPoStCircuitProof(postCfg sector.PoStInstanceCfg, privateProof PrivatePoStProof) sector.PoStProof {
	panic("TODO")

	var proofBytes []byte
	panic("TODO")

	postProof := sector.PoStProof_I{
		ProofBytes_: proofBytes,
	}

	return &postProof
}
