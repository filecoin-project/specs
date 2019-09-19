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
    participant StorageProving
    participant FilProofs

    participant Clock

    participant libp2p

    Note over RetrievalClient,RetrievalProvider: RetrievalMarketSubsystem
    Note over StorageClient,StorageProvider: StorageMarketSubsystem
    Note over Blockchain,StoragePowerActor: BlockchainGroup
    Note over StorageMining,StorageProving: MiningGroup

    opt RetrievalDealMake
        RetrievalClient ->> RetrievalProvider: DealProposal
        RetrievalProvider ->> RetrievalClient: Accepted, Rejected
    end

    opt RetrievalQuery
        RetrievalClient ->> RetrievalProvider: Query(CID)
        RetrievalProvider ->> RetrievalClient: MinPrice, Unavail
    end

    opt RegisterStorageMiner
        StorageMining->>StorageMining: CreateMiner(WorkerPubKey, PledgeCollateral)
        StorageMining->StoragePowerActor: RegisterMiner(OwnerAddr, WorkerPubKey)
        StoragePowerActor->StorageMining: StorageMinerActor
    end

    opt StorageDealMake
        Note left of StorageClient: Piece, PieceCID
        StorageClient->>StorageProvider: DealProposal
        StorageProvider->>StorageClient: DealResponse,DealAccepted,Deal
        Note left of StorageClient: Piece, PieceCID, Deal
        Note right of StorageProvider: Piece, PieceCID, Deal
        StorageClient->>StorageProvider: StorageDealQuery
        StorageProvider->>StorageClient: DealResponse,DealAccepted,Deal
    end

    opt AddingDealToSector
        StorageProvider->>StorageMining: MadeDeal(Deal,PieceRef)
        StorageMining->>+SectorIndexing: AddToSector(Deal, PieceRef)
        SectorIndexing-->>SectorIndexing: SectorID ← PackDealIntoSector(Deal)
        SectorIndexing-->>SectorIndexing: PIP ← PackSector(SectorID)
        SectorIndexing->>-StorageMining: SectorID
        StorageMining->>StorageProvider: DealInSector(Deal,PieceRef,PIP,SectorID)
    end

    opt ClientQuery
        StorageClient->>StorageProvider: StorageDealQuery
        StorageProvider->>StorageClient: DealResponse,DealAccepted,Deal,PIP
    end

    opt SealingSector
        StorageMining->>+StorageProving: SealSector(SectorID, ReplicaCfg)
        StorageProving-->>StorageProving: SealOutputs ← Seal(SectorID, ReplicaCfg)
        StorageProving->>-StorageMining: (SectorID,SealOutputs)
        opt CommitSector
            StorageMining-->>StorageMinerActor: CommitSector(SectorID,SealCommitment,SealProof)
            StorageMinerActor-->>+FilProofs: VerifySeal(SectorID, OnSectorInfo)
            FilProofs-->>-StorageMinerActor: {1,0} ← VerifySeal
            alt 1 - success
                StorageMinerActor-->>StorageMinerActor: ...Update State...
                StorageMinerActor-->>StoragePowerActor: UpdatePower(MinerAddr)
                StorageMinerActor-->>StorageMinerActor: 1 ← CommitSector(.)
            else 0 - failure
                StorageMinerActor-->>StorageMinerActor: 0 ← CommitSector(.)
            end
        end
    end

    opt ClientQuery
        StorageClient->>StorageProvider: StorageDealQuery
        StorageProvider->>StorageClient: DealResponse,DealAccepted,Deal,PIP,SealedSectorCID
    end

    loop StorageDealCollect
        Note Right of StorageProvider: Deal
        alt Via Client
            StorageProvider ->> StorageClient: ReconcileRequest(Deal, [Voucher])
            opt If Client Does Not Have PIP
                StorageClient -->> StorageProvider: StorageDealQuery(Deal)
                StorageProvider -->> StorageClient: PieceID, SectorID, PIP
            end
            StorageClient -->> Blockchain: VerifySectorExists(SectorID)
            StorageClient --> StorageClient: VerifyPIP(SectorID, PIP)
            StorageClient -->> StorageClient: ReconcileResponse ← SignVouchers([Voucher])
            StorageClient ->> StorageProvider: ReconcileResponse
            StorageProvider ->> PaymentChannelActor: RedeemVoucher(ReconcileResponse.Voucher)

        else Via Blockchain

            StorageProvider ->> StorageMarketActor: RedeemVoucher(Voucher, PIP)
        end
    end

    loop BlockReception
        libp2p -->> BlockSyncer: block ← Subscription.Next()
        BlockSyncer -->+ BlockSyncer: ValidateBlock(block)
        BlockSyncer -->> FilProofs: Proofs.ValidateBlock(block)
        BlockSyncer -->> StoragePowerConsensus: SPC.ValidateBlock(block)
        BlockSyncer -->- Blockchain: ConsiderBlock(block)
        Blockchain -->> Blockchain: VerifyStateRoot(block)
        Blockchain -->> Blockchain: StateTree ← GetStateTree(block)

        alt Round Cutoff
            Blockchain -->> Blockchain: AssembleTipsets([block])
            Blockchain -->+ StoragePowerConsensus: BestTipset([Tipset])
            Blockchain -->- StoragePowerConsensus: {Tipset} ← BestTipset([Tipset])
            Blockchain -->> Blockchain: ApplyStateTree(StateTree)
        end
    end


    loop BlockProduction
        alt New Tipset
            Blockchain -->> Blockchain: tipset ← onNewBestTipset()
        else Retrying on null block
            Blockchain -->> Blockchain: tipset ← addNullBlock(tipset)
        end
        Blockchain -->> BlockProducer: NewBestTipset(tipset)
        BlockProducer -->+ StorageMining: ScratchTicket(randomness)
        StorageMining -->- BlockProducer: [scratchedTicket] ← ScratchTicket(ticket)
        BlockProducer -->+ StoragePowerConsensus: tryLeaderElection([scratchedTicket])
        StoragePowerConsensus -->- BlockProducer: {[(Address, ElectionProof)]} ← tryLeaderElection([ticket])
        opt tryLeaderElection - success, for each ElectionProof
            BlockProducer -->+ MessagePool: GetMessages()
            MessagePool -->- BlockProducer: [Message] ← GetMessages()
            BlockProducer -->+ BlockProducer: AssembleBlock(ElectionProof, Messages)
            BlockProducer  -->- BlockProducer: block ← AssembleBlock()
            BlockProducer -->> BlockSyncer: SendBlock(block)
        end
    end


    loop PoStSubmission
            Note Right of PostSubmission: in every proving period
            StorageMining -->> Blockchain: GetPoStRandomness()
            Blockchain -->> StorageMining: randomness ← GetPoStRandomness()
            StorageMining -->> StorageProving: GeneratePoSt(randomness)
            StorageProving -->> StorageMining: (PoSt) ← GeneratePoSt(randomness)
            StorageMining -->> StorageMinerActor: SubmitPost(PoStProof, DoneSet)
        alt PoStCompletion
            StorageMining -->> SectorIndexing: DoneSet(Sector)
        end
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
