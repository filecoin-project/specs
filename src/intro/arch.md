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
    participant StorageProvider

    participant PaymentChannelActor
    participant PaymentsSubsystem

    participant BlockchainSubsystem
    participant BlockSyncer
    participant BlockProducer

    participant StoragePowerConsensusSubsystem
    participant StoragePowerActor

    participant StorageMiningSubsystem
    participant StorageMinerActor
    participant SectorIndexerSubsystem
    participant StorageProvingSubsystem

    participant FilecoinProofsSubsystem
    participant ClockSubsystem
    participant libp2p

    Note over RetrievalClient,RetrievalProvider: RetrievalMarketSubsystem
    Note over StorageClient,StorageProvider: StorageMarketSubsystem
    Note over BlockchainSubsystem,StoragePowerActor: BlockchainGroup
    Note over StorageMiningSubsystem,StorageProvingSubsystem: MiningGroup

    opt RetrievalDealMake
        RetrievalClient ->>+ RetrievalProvider: DealProposal
        RetrievalProvider -->>- RetrievalClient: {Accepted, Rejected}
    end

    opt RetrievalQuery
        RetrievalClient ->>+ RetrievalProvider: Query(CID)
        RetrievalProvider -->>- RetrievalClient: MinPrice, Unavail
    end

    opt RegisterStorageMiner
        StorageMiningSubsystem ->> StorageMiningSubsystem: CreateMiner(ownerPubKey PubKey, workerPubKey PubKey, pledgeAmt TokenAmount)
        StorageMiningSubsystem ->>+ StoragePowerActor: RegisterMiner(OwnerAddr, WorkerPubKey)
        StoragePowerActor -->>- StorageMiningSubsystem: StorageMinerActor
    end

    opt StorageDealMake
        Note left of StorageClient: Piece, PieceCID
        StorageClient ->> StorageProvider: ProposeStorageDeal(StorageDealProposal)
        StorageClient ->>+ StorageProvider: QueryStorageDealStatus(StorageDealQuery)
        StorageProvider -->>- StorageClient: StorageDealResponse{StorageDealAccepted, Deal} 
         
        Note left of StorageClient: Piece, PieceCID, Deal
        Note right of StorageProvider: Piece, PieceCID, Deal
        StorageClient ->>+ StorageProvider: QueryStorageDealStatus(StorageDealQuery)
        StorageProvider -->>- StorageClient: StorageDealResponse{StorageDealStarted, Deal}
    end

    opt AddingDealToSector
        StorageProvider ->>+ StorageMiningSubsystem: HandleStorageDeal(Deal, PieceRef)
        StorageMiningSubsystem ->>+ SectorIndexerSubsystem: AddPieceToSector(Deal, SectorID)
        SectorIndexerSubsystem ->> SectorIndexerSubsystem: PIP ← GetPieceInclusionProof(Deal)
        SectorIndexerSubsystem -->>- StorageMiningSubsystem: SectorID
        StorageMiningSubsystem ->>- StorageProvider: NotifyStorageDealStaged(Deal,PieceRef,PIP,SectorID)
    end

    opt ClientQuery
        StorageClient ->>+ StorageProvider: QueryStorageDealStatus(StorageDealQuery)
        StorageProvider -->>- StorageClient: StorageDealResponse{StorageDealStaged,Deal,PIP}
    end

    opt SealingSector
        StorageMiningSubsystem ->>+ StoragePowerConsensusSubsystem: GetSealSeed(Chain, Epoch)
        StoragePowerConsensusSubsystem -->>- StorageMiningSubsystem: Seed
        StorageMiningSubsystem ->>+ StorageProvingSubsystem: SealSector(Seed, SectorID, ReplicaCfg)
        StorageProvingSubsystem ->>+ SectorSealer: Seal(Seed, SectorID, ReplicaCfg)
        SectorSealer -->>- StorageProvingSubsystem: SealOutputs
        StorageProvingSubsystem ->>- StorageMiningSubsystem: SealOutputs
        opt CommitSector
            StorageMiningSubsystem ->> StorageMinerActor: CommitSector(Seed, SectorID, SealCommitment, SealProof)
            StorageMinerActor ->>+ FilecoinProofsSubsystem: VerifySeal(SectorID, OnSectorInfo)
            FilecoinProofsSubsystem -->>- StorageMinerActor: {1,0}
            alt 1 - success
                StorageMinerActor ->> StoragePowerActor: IncrementPower(StorageMiner.WorkerPubKey)
            else 0 - failure
                StorageMinerActor -->> StorageMiningSubsystem: CommitSectorError
            end
        end
    end

    loop PoStSubmission
        Note Right of PoStSubmission: in every proving period
        Note Right of PoStSubmission: DoneSet
        StorageMiningSubsystem ->>+ StoragePowerConsensusSubsystem: GetPoStChallenge(Chain, Epoch)
        StoragePowerConsensusSubsystem -->>- StorageMiningSubsystem: challenge
        StorageMiningSubsystem ->>+ StorageProvingSubsystem: GeneratePoSt(challenge, [SectorID])
        StorageProvingSubsystem ->>+ StorageProvingSubsystem: GeneratePoSt(challenge, [SectorID])
        StorageMiningSubsystem ->>+ StorageProvingSubsystem: GeneratePoSt(challenge, [SectorID])
        StorageProvingSubsystem -->>- StorageMiningSubsystem: PoStProof
        StorageMiningSubsystem ->> StorageMinerActor: SubmitPoSt(PoStProof, DoneSet)
    end

    opt ClientQuery
        StorageClient ->>+ StorageProvider: QueryStorageDealStatus(StorageDealQuery)
        StorageProvider -->>- StorageClient: StorageDealResponse{SealingParams,DealComplete,...}
    end
    
    loop StorageDealCollect
        Note Right of StorageProvider: Deal
        alt Via Client
            StorageProvider ->>+ StorageClient: RequestVouchersApproval(Deal, [Voucher])
            opt If Client Does Not Have PIP
                StorageClient ->>+ StorageProvider: QueryStorageDealStatus(StorageDealQuery)
                StorageProvider -->>- StorageClient: StorageDealResponse{SealingParams,DealComplete,SectorID, PIP}
                StorageClient ->> StorageClient: {0, 1} ← VerifyPIP(SectorID, PIP)
            end
            StorageClient ->>+ BlockchainSubsystem: VerifySectorExists(SectorID)
            BlockchainSubsystem ->>- StorageClient: {0, 1}
            StorageClient ->> StorageClient: VouchersApprovalResponse ← ApproveVouchers([Voucher])
            StorageClient -->>- StorageProvider: VouchersApprovalResponse
            StorageProvider ->> PaymentChannelActor: RedeemVoucherWithApproval(VoucherApprovalResponse.Voucher)
        else Via Blockchain
            StorageProvider ->> PaymentChannelActor: RedeemVoucherWithPIP(Voucher, PIP)
        end
    end

    loop BlockReception
        BlockSyncer ->>+ libp2p: Subscribe(OnNewBlock)
        libp2p -->>- BlockSyncer: Event(OnNewBlock, block)
        BlockSyncer ->> BlockSyncer: ValidateSyntax(block)
        BlockSyncer ->>+ BlockchainSubsystem: HandleBlock(block)
        BlockchainSubsystem ->> BlockchainSubsystem: ValidateBlock(block)
        BlockchainSubsystem ->> StoragePowerConsensusSubsystem: ValidateBlock(block)
        BlockchainSubsystem ->> FilecoinProofsSubsystem: ValidateBlock(block)
        BlockchainSubsystem ->>- BlockchainSubsystem: StateTree ← TryGenerateStateTree(block)

        alt Round Cutoff
            WallClock -->> BlockchainSubsystem: AssembleTipsets()
            BlockchainSubsystem ->> BlockchainSubsystem: [Tipset] ← AssembleTipsets()
            BlockchainSubsystem ->> BlockchainSubsystem: Tipset ← ChooseTipset([Tipset])
            BlockchainSubsystem ->> Blockchain: ApplyStateTree(StateTree)
        end
    end

    loop BlockProduction
        alt New Tipset
            BlockchainSubsystem ->> StorageMiningSubsystem: OnNewTipset(Chain, Epoch)
        else Null block last round
            WallClock ->> StorageMiningSubsystem: OnNewRound()
            Note Right of WallClock: epoch is incremented by 1
        end
        StorageMiningSubsystem ->>+ StoragePowerConsensusSubsystem: GetElectionArtifacts(Chain, Epoch)
        StoragePowerConsensusSubsystem ->> StoragePowerConsensusSubsystem: TK ← TicketAtEpoch(Chain, Epoch-k)
        StoragePowerConsensusSubsystem ->> StoragePowerConsensusSubsystem: T1 ← TicketAtEpoch(Chain, Epoch-1)
        StoragePowerConsensusSubsystem -->>- StorageMiningSubsystem: TK, T1
       
        loop forEach StorageMiningSubsystem.StorageMiner
            StorageMiningSubsystem ->> StorageMiningSubsystem: EP ← DrawElectionProof(TK.randomness(), StorageMiner.WorkerKey)
            alt New Tipset
                StorageMiningSubsystem ->> StorageMiningSubsystem: T0 ← GenerateNextTicket(T1.randomness(), StorageMiner.WorkerKey)            
            else Null block last round
                StorageMiningSubsystem ->> StorageMiningSubsystem: T1 ← GenerateNextTicket(T0.randomness(), StorageMiner.WorkerKey)   
                Note Right of StorageMiningSubsystem: Using tickets derived in failed election proof in last epoch
            end
            StorageMiningSubsystem ->>+ StoragePowerConsensusSubsystem: TryLeaderElection(EP)
            StoragePowerConsensusSubsystem -->>- StorageMiningSubsystem: {1, 0}
            opt 1- success
                StorageMiningSubsystem ->> BlockProducer: GenerateBlock(EP, T0, Tipset, StorageMiner.Address)
                BlockProducer ->>+ MessagePool: GetMostProfitableMessages(StorageMiner.Address)
                MessagePool -->>- BlockProducer: [Message]
                BlockProducer ->> BlockProducer: block ← AssembleBlock([Message], Tipset, EP, T0, StorageMiner.Address)
                BlockProducer ->> BlockSyncer: PropagateBlock(block)
            end
        end
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
            StorageMinerActor -->>  StorageMinerActor: SuspendMiner(Address)
        else Undeclared Storage Fault
            ClockSubsystem -->>  StorageMinerActor: SuspendMiner(Address)
        end

        alt Recovery in Grace Period
            StorageMinerActor -->> StorageMinerActor: SubmitPoSt(PoStProof, DoneSet)
            StorageMinerActor -->> StorageMinerActor: UpdatePower()
        else Recovery past Grace Period
            Clock -->>  StorageMinerActor: SlashStorageFault()
            StorageMinerActor -->> StorageMinerActor: AddCollateral()
            StorageMinerActor -->> StorageMinerActor: SubmitPoSt(PoStProof, DoneSet)
            StorageMinerActor -->> StoragePowerActor: UpdatePower()
        else Recovery Past Sector Failure Timeout
            Clock -->>  StorageMinerActor: SlashStorageFault()
            StorageMinerActor -->> StorageMinerActor: AddCollateral()
            StorageMiningSubsystem-->>StorageProvingSubsystem: SealSector(SectorID, ReplicaCfg)
        end
    end

    opt Consensus Fault
        StorageMinerActor -->> StoragePowerActor: DeclareConsensusFault(Proof)
        StoragePowerActor -->+ StoragePowerConsensusSubsystem: ValidateFault(Proof)

        alt Valid Fault
            StoragePowerConsensusSubsystem -->> StoragePowerActor: TerminateMiner()
            StoragePowerConsensusSubsystem -->> StoragePowerActor: SlashPledgeCollateral(Address)
            StoragePowerConsensusSubsystem -->- StorageMinerActor: Reward ← DeclareConsensusFault(Proof)
        end
    end

{{% /mermaid %}}
