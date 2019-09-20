---
title: "Architecture Diagram"
---

## Overview Diagram

<img src="overview.svg" />

## Storage Flow

{{% mermaid %}}
sequenceDiagram

    participant RetrievalClient
    participant RetrievalProvider

    participant StorageClient
    participant StorageMarketActor
    participant StorageProvider

    participant PaymentChannelActor

    participant Blockchain
    participant BlockSyncer

    participant BlockProducer

    participant StoragePowerConsensus
    participant StoragePowerActor

    participant StorageMining
    participant StorageMinerActor
    participant SectorIndexing
    participant FilecoinProofs

    participant Clock

    participant libp2p

    Note over RetrievalClient,RetrievalProvider: RetrievalMarketSubsystem
    Note over StorageClient,StorageProvider: StorageMarketSubsystem
    Note over Blockchain,StoragePowerActor: BlockchainGroup
    Note over StorageMining: MiningGroup

    opt RetrievalDealMake
        RetrievalClient ->> RetrievalProvider: DealProposal
        RetrievalProvider ->> RetrievalClient: Accepted, Rejected
    end

    opt RetrievalQuery
        RetrievalClient ->> RetrievalProvider: Query(CID)
        RetrievalProvider ->> RetrievalClient: MinPrice, Unavail
    end

    opt RegisterStorageMiner
        StorageMining->>StorageMining: CreateMiner(ownerPubKey PubKey, workerPubKey PubKey, pledgeAmt TokenAmount)
        StorageMining->StoragePowerActor: RegisterMiner(OwnerAddr, WorkerPubKey)
        StoragePowerActor->StorageMining: StorageMinerActor
    end

    opt StorageDealMake
        Note left of StorageClient: Piece, PieceCID
        StorageClient->>StorageProvider: proposeStorageDeal(StorageDealProposal)
        StorageProvider->>StorageClient: StorageDealResponse{StorageDealAccepted, Deal} ← QueryStorageDealStatus(StorageDealQuery)
        Note left of StorageClient: Piece, PieceCID, Deal
        Note right of StorageProvider: Piece, PieceCID, Deal
        StorageClient->>StorageProvider: QueryStorageDealStatus(StorageDealQuery)
        StorageProvider->>StorageClient: StorageDealResponse{StorageDealStarted, Deal} ← QueryStorageDealStatus(StorageDealQuery)
    end

    opt AddingDealToSector
        StorageProvider->>StorageMining: HandleStorageDeal(Deal,PieceRef)
        StorageMining->>+SectorIndexing: HandlePiece(Deal, PieceRef) SectorID
        SectorIndexing-->>SectorIndexing: AddPieceToSector(Deal, SectorID)
        SectorIndexing-->>SectorIndexing: PIP ← GetPieceInclusionProof(Deal)
        SectorIndexing->>-StorageMining: SectorID ← HandlePiece(Deal, PieceRef)
        StorageMining->>StorageProvider: NotifyStorageDealStaged(Deal,PieceRef,PIP,SectorID)
    end

    opt ClientQuery
        StorageClient->>StorageProvider: QueryStorageDealStatus(StorageDealQuery)
        StorageProvider->>StorageClient: StorageDealResponse{StorageDealStaged,Deal,PIP} ← QueryStorageDealStatus(StorageDealQuery)
    end

    opt SealingSector
        StorageMining->>+StorageProving: SealSector(Seed, SectorID, ReplicaCfg)
        StorageProving-->>StorageProving: SealOutputs ← Seal(Seed, SectorID, ReplicaCfg)
        StorageProving->>-StorageMining: (SectorID,SealOutputs) ← SealSector(Seed, SectorID, ReplicaCfg)
        opt CommitSector
            StorageMining-->>StorageMinerActor: CommitSector(Seed, SectorID, SealCommitment, SealProof)
            StorageMinerActor-->>+FilecoinProofs: VerifySeal(SectorID, OnSectorInfo)
            FilecoinProofs-->>-StorageMinerActor: {1,0} ← VerifySeal
            alt 1 - success
                StorageMinerActor-->>StoragePowerActor: IncrementPower(MinerAddr)
            else 0 - Fail
                StorageMinerActor-->>StoragePowerActor: CommitSectorError
            end
        end
    end

    opt ClientQuery
        StorageClient->>StorageProvider: QueryStorageDealStatus(StorageDealQuery)
        StorageProvider->>StorageClient: StorageDealResponse{SealingParams,DealComplete,?} ← QueryStorageDealStatus(StorageDealQuery)
    end

    loop StorageDealCollect
        Note Right of StorageProvider: Deal
        alt Via Client
            StorageProvider ->> StorageClient: RequestVouchersApproval(Deal, [Voucher])
            opt If Client Does Not Have PIP
                StorageClient -->> StorageProvider: QueryStorageDealStatus(StorageDealQuery)
                StorageProvider -->> StorageClient: StorageDealResponse{SealingParams,DealComplete,SectorID, PIP} ← QueryStorageDealStatus(StorageDealQuery)
                StorageClient --> StorageClient: bool ← VerifyPIP(SectorID, PIP)
            end
            StorageClient -->> Blockchain: bool ← VerifySectorExists(SectorID)
            StorageClient -->> StorageClient: VoucherApprovalResponse ← ApproveVouchers([Voucher])
            StorageClient -->> StorageProvider: VoucherApprovalResponse
            StorageProvider -->> PaymentChannelActor: RedeemVoucherWithApproval(VoucherApprovalResponse.Voucher)
        else Via Blockchain
            StorageProvider -->> PaymentChannelActor: RedeemVoucherWithPIP(Voucher, PIP)
        end
    end

    loop BlockReception
        BlockSyncer -->> libp2p: Subscribe(blockTopic)
        libp2p -->> BlockSyncer: Event(blockTopic, block)
        BlockSyncer -->> BlockSyncer: ValidateSyntax(block)
        BlockSyncer -->+ Blockchain: HandleBlock(block)
        Blockchain -->> Blockchain: ValidateBlock(block)
        Blockchain -->> StoragePowerConsensus: ValidateBlock(block)
        Blockchain -->> FilecoinProofs: ValidateBlock(block)
        Blockchain -->- Blockchain: StateTree ← TryGenerateStateTree(block)

        alt Round Cutoff
            WallClock -->> Blockchain: [Tipset] ← AssembleTipsets()
            Blockchain -->> Blockchain: Tipset ← ChooseTipset([Tipset])
            Blockchain -->> BlockChain: ApplyStateTree(StateTree)
        end
    end

    loop BlockProduction
        alt New Tipset
            Blockchain -->> BlockChain: ApplyStateTree(StateTree)
        else Retrying on null block
            BlockProducer -->+ StoragePowerConsensus: tryLeaderElection([T0])
        end
        Blockchain -->> StorageMining: OnNewTipset(Chain, Epoch)
        StorageMining -->+ StoragePowerConsensus: GetElectionArtifacts(Chain, Epoch)
        StoragePowerConsensus -->> StoragePowerConsensus: TK ← TicketAtEpoch(Chain, Epoch-k)
        StoragePowerConsensus -->> StoragePowerConsensus: TK ← TicketAtEpoch(Chain, Epoch)
        StoragePowerConsensus -->- StorageMining: TK, T1 ← GetElectionArtifacts(Chain, Epoch)
        loop forEach StorageMining.StorageMiner
            StorageMining -->> StorageMining: EP ← DrawElectionProof(TK.randomness(), StorageMiner.WorkerKey)
            StorageMining -->> StorageMining: T0 ← GenerateNextTicket(T1.randomness(), StorageMiner.WorkerKey)
            StorageMining -->> StoragePowerConsensus: TryLeaderElection(EP)
            StoragePowerConsensus -->> StorageMining: 1/0 ← TryLeaderElection(EP)
            alt Success
                StorageMining -->> BlockProducer: GenerateBlock(EP, T0, Tipset, StorageMiner.Address)
                BlockProducer -->+ MessagePool: MostProfitableMessage(StorageMiner.Address)
                MessagePool -->- BlockProducer: [Message] ← MostProfitableMessage(StorageMiner.Address)
                BlockProducer -->+ BlockProducer: AssembleBlock(ElectionProof, [Message], Tipset, EP, T0, StorageMiner.Address)
                BlockProducer  -->- BlockProducer: block ← AssembleBlock()
                BlockProducer -->> BlockSyncer: PropagateBlock(block)
            end
        end
    end

    loop PoStSubmission
        Note Right of PostSubmission: in every proving period
        Note Right of PostSubmission: DoneSet
        StorageMining -->> StoragePowerConsensus: GetPoStChallenge(Chain, Epoch)
        StoragePowerConsensus -->> StorageMining: challenge ← GetPoStChallenge(Chain, Epoch)
        StorageMining -->> StorageProving: GeneratePoSt(challenge, [SectorID])
        StorageProving -->> StorageMining: (PoStProof) ← GeneratePoSt(challenge, [SectorID])
        StorageMining -->> StorageMinerActor: SubmitPost(PoStProof, DoneSet)
    end
    
    opt MiningScheduler
        Note Right of MiningScheduler: Schedule and resume PoSts
        Note Right of MiningScheduler: Schedule and resume SEALs
        Note Right of MiningScheduler: Process expired deals
        Note Right of MiningScheduler: Process deal payments
        Note Right of MiningScheduler: Maintain FaultSet
        Note Right of MiningScheduler: Maintain DoneSet
    end

    opt Storage Fault

        alt Declared Storage Fault
            StorageMinerActor -->> StorageMinerActor: UpdateFaults(FaultSet)
            StorageMinerActor -->>  StoragePowerConsensus: SuspendMiner(Address)
        else Undeclared Storage Fault
            Clock -->>  StoragePowerConsensus: SuspendMiner(Address)
        end

        alt Recovery in Grace Period
            StorageMinerActor -->> StorageMinerActor: SubmitPost(PoStProof, DoneSet)
            StorageMinerActor -->> StoragePowerConsensus: UpdatePower()
        else Recovery past Grace Period
            Clock -->>  StorageMinerActor: SlashStorageFault()
            StorageMinerActor -->> StorageMinerActor: AddCollateral()
            StorageMinerActor -->> StorageMinerActor: SubmitPost(PoStProof, DoneSet)
            StorageMinerActor -->> StoragePowerActor: UpdatePower()
        else Recovery Past Sector Failure Timeout
            Clock -->>  StorageMinerActor: SlashStorageFault()
            StorageMinerActor -->> StorageMinerActor: AddCollateral()
            StorageMining-->>StorageProving: SealSector(SectorID, ReplicaCfg)
        end
    end

    opt Consensus Fault
        StorageMinerActor -->> StoragePowerActor: DeclareConsensusFault(Proof)
        StoragePowerActor -->+ StoragePowerConsensus: ValidateFault(Proof)

        alt Valid Fault
            StoragePowerConsensus -->> StoragePowerActor: TerminateMiner()
            StoragePowerConsensus -->> StoragePowerActor: SlashPledgeCollateral(Address)
            StoragePowerConsensus -->- StorageMinerActor: Reward ← DeclareConsensusFault(Proof)
        end
    end

{{% /mermaid %}}
