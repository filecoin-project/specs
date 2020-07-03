---
title: "Architecture Diagrams"
---


# Filecoin Systems
<!-- <script type="text/javascript" src="https://unpkg.com/svg-pan-zoom@3.6.1/dist/svg-pan-zoom.min.js">  -->
<script src='https://unpkg.com/panzoom@9.2.4/dist/panzoom.min.js'></script>
<script type="text/javascript">

function statusIndicatorsShow() {
  var $uls = document.querySelectorAll('.statusIcon')
  $uls.forEach(function (el) {
    el.classList.remove('hidden')
  })
  return false; // stop click event
}

function statusIndicatorsHide() {
  var $uls = document.querySelectorAll('.statusIcon')
  $uls.forEach(function (el) {
    el.classList.add('hidden')
  })
  return false; // stop click event
}

</script>


Status Legend:

- 🛑 **Bare** - Very incomplete at this time.
  - **Implementors:** This is far from ready for you.
- ⚠️ **Rough** -- work in progress, heavy changes coming, as we put in place key functionality.
  - **Implementors:** This will be ready for you soon.
- 🔁 **Refining** - Key functionality is there, some small things expected to change. Some big things may change.
  - **Implementors:** Almost ready for you. You can start building these parts, but beware there may be changes still.
- ✅ **Stable** - Mostly complete, minor things expected to change, no major changes expected.
  - **Implementors:** Ready for you. You can build these parts.

*Note that the status relates to the state of the spec either written out either in english or in code. The goal is for the spec to eventually be fleshed out in both language-sets.*

[<a href="#" onclick="return statusIndicatorsShow();">Show</a> / <a href="#" onclick="return statusIndicatorsHide();">Hide</a> ] status indicators


{{</* incTocMap "/docs/systems" 2 "colorful" */>}}


# Overview Diagram

TODO:

- cleanup / reorganize
  - this diagram is accurate, and helps lots to navigate, but it's still a bit confusing
  - the arrows and lines make it a bit hard to follow. We should have a much cleaner version (maybe based on [C4](https://c4model.com))
- reflect addition of Token system
  - move data_transfers into Token

{{</* diagram src="../diagrams/overview1/overview.dot.svg" title="Protocol Overview Diagram" */>}}

# Protocol Flow Diagram -- deals on chain

{{< mermaid class="text-center">}}
sequenceDiagram

    participant RetrievalClient
    participant RetrievalProvider

    participant StorageMarketParticipant1
    participant StorageMarketParticipant2
    participant StorageMarketActor

    participant PaymentChannelActor
    participant PaymentSubsystem

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
    Note over StorageMarketParticipant,StorageMarketActor: StorageMarketSubsystem
    Note over BlockchainSubsystem,StoragePowerActor: BlockchainGroup
    Note over StorageMiningSubsystem,StorageProvingSubsystem: MiningGroup

    opt RetrievalMarket
        RetrievalClient ->>+ RetrievalProvider: NewRetrievalQuery(RetreivalQuery)
        RetrievalProvider -->>- RetrievalClient: RetrievalQueryResponse
        RetrievalClient ->>+ RetrievalProvider: NewRetrievalDealProposal(RetrievalDealProposal)
        RetrievalProvider -->> RetrievalProvider: AcceptRetrievalDealProposal(RetrievalDealPropsal)
        RetrievalProvider -->>- RetrievalClient: NewPaymentChannel
    end

    opt RegisterStorageMiner
        StorageMiningSubsystem ->> StorageMiningSubsystem: CreateMiner(ownerPubKey PubKey, workerPubKey PubKey, pledgeAmt TokenAmount)
        StorageMiningSubsystem ->>+ StoragePowerActor: CreateStorageMiner(OwnerAddr, WorkerPubKey)
        StoragePowerActor -->>- StorageMiningSubsystem: StorageMinerActor
    end

    opt StorageDealMake
        StorageMarketParticipant1 -->> StorageMarketActor: RegisterParticipant(TokenAmount)
        StorageMarketParticipant1 -->> StorageMarketActor: AddBalance(TokenAmount)
        StorageMarketParticipant1 -->> StorageMarketActor: WithdrawBalance(TokenAmount)
        StorageMarketParticipant1 -->+ StorageMarketActor: CheckLockedBalance(Address)
        StorageMarketActor -->- StorageMarketParticipant1: TokenAmount


        StorageMarketParticipant1 -->+ StorageMarketParticipant2: NewStorageDealProposal(StorageDealProposal)
        StorageMarketParticipant2 --> StorageMarketActor: VerifyBalance(StorageMarketParticipant1)
        StorageMarketParticipant2 -->> StorageMarketParticipant2: SignStorageDealProposal(StorageDealProposal)
        StorageMarketParticipant2 -->- StorageMarketParticipant1: NewStorageDeal(StorageDeal)
        StorageMarketParticipant2 -->> StorageMarketActor: PublishStorageDeal(StorageDeal)
        StorageMarketParticipant1 -->> StorageMarketActor: PublishStorageDeal(StorageDeal)

    end

    opt AddingDealToSector
        StorageMarketActor ->>+ StorageMiningSubsystem: HandleStorageDeal(Deal)
        StorageMiningSubsystem ->>+ SectorIndexerSubsystem: AddDealToSector(Deal, SectorID)
        SectorIndexerSubsystem ->> SectorIndexerSubsystem: IndexSectorByDealExpiration(SectorID, Deal)
        SectorIndexerSubsystem -->>- StorageMiningSubsystem: (SectorID, Deal)
        StorageMiningSubsystem ->>- StorageMarketActor: NotifyStorageDealStaged(Deal,SectorID)
    end

    opt ClientQuery
        StorageMarketParticipant1 ->>+ StorageMarketParticipant2: QueryStorageDealProposalStatus(StorageDealProposalQuery)
        StorageMarketParticipant2 -->>- StorageMarketParticipant1: StorageDealProposalQueryResponse
        StorageMarketParticipant1 ->>+ StorageMarketParticipant2: QueryStorageDealStatus(StorageDealQuery)
        StorageMarketParticipant2 -->>- StorageMarketParticipant1: StorageDealQueryResponse
    end

    opt SealingSector
        StorageMiningSubsystem ->>+ StoragePowerConsensusSubsystem: GetSealSeed(Chain, Epoch)
        StoragePowerConsensusSubsystem -->>- StorageMiningSubsystem: Seed
        StorageMiningSubsystem ->>+ StorageProvingSubsystem: SealSector(Seed, SectorID, ReplicaCfg)
        StorageProvingSubsystem ->>+ SectorSealer: Seal(Seed, SectorID, ReplicaCfg)
        SectorSealer -->>- StorageProvingSubsystem: SealOutputs
        StorageProvingSubsystem ->>- StorageMiningSubsystem: SealOutputs
        opt CommitSector
            StorageMiningSubsystem ->> StorageMinerActor: CommitSector(Seed, SectorID, SealCommitment, SealProof, [&Deal], [Deal])
            StorageMinerActor ->>+ FilecoinProofsSubsystem: VerifySeal(SectorID, OnSectorInfo)
            FilecoinProofsSubsystem -->>- StorageMinerActor: {1,0}
            alt 1 - success
                StorageMinerActor ->> StorageMarketActor: AddDeal(SectorID, [&Deal], DealStatusOnChain)
                StorageMinerActor ->> StorageMarketActor: AddDeal(SectorID, [Deal], DealStatusPending)
                StorageMinerActor ->> StoragePowerActor: IncrementPower(StorageMiner.WorkerPubKey)
            else 0 - failure
                StorageMinerActor -->> StorageMiningSubsystem: CommitSectorError
            end
        end
    end

    loop BlockReception
        BlockSyncer ->>+ libp2p: Subscribe(OnNewBlock)
        libp2p -->>- BlockSyncer: Event(OnNewBlock, block)
        BlockSyncer ->> BlockSyncer: ValidateBlockSyntax(block)
        BlockSyncer ->>+ BlockchainSubsystem: HandleBlock(block)
        BlockchainSubsystem ->> BlockchainSubsystem: ValidateBlock(block)
        BlockchainSubsystem ->> StoragePowerConsensusSubsystem: ValidateBlock(block)
        BlockchainSubsystem ->> FilecoinProofsSubsystem: ValidateBlock(block)
        BlockchainSubsystem ->>- BlockchainSubsystem: StateTree ← TryGenerateStateTree(block)

        opt Round Cutoff
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
                StorageMiningSubsystem ->>+ StoragePowerConsensusSubsystem: GetPoStChallenge(Chain, Epoch)
                StoragePowerConsensusSubsystem -->>- StorageMiningSubsystem: challenge
                StorageMiningSubsystem ->>+ StorageProvingSubsystem: GeneratePoSt(challenge, [SectorID])
                StorageProvingSubsystem -->>- StorageMiningSubsystem: PoStProof
                opt SubmitPoSt
                    StorageMiningSubsystem ->> StorageMinerActor: SubmitPoSt(PoStProof, DoneSet)
                    StorageMinerActor ->>+ FilecoinProofsSubsystem: VerifyPoSt(PoStProof)
                    FilecoinProofsSubsystem -->>- StorageMinerActor: {1,0}
                    alt 1 - success
                        StorageMinerActor ->> StorageMinerActor:  UpdateDoneSet()
                        StorageMinerActor ->> StorageMarketActor: HandleStorageDealPayment()
                        StorageMarketActor ->> StorageMarketActor: CloseExpiredStorageDeal()
                    else 0 - failure
                        StorageMinerActor -->> StorageMiningSubsystem: PoStError
                        StorageMarketActor -->> StorageMarketActor: SlashStorageCollateral()
                    end
                end

                StorageMiningSubsystem ->> BlockProducer: GenerateBlock(EP, T0, Tipset, StorageMiner.Address)
                BlockProducer ->>+ MessagePool: GetMostProfitableMessages(StorageMiner.Address)
                MessagePool -->>- BlockProducer: [Message]
                BlockProducer ->> BlockProducer: block ← AssembleBlock([Message], Tipset, EP, T0, StorageMiner.Address)

                BlockProducer ->> BlockSyncer: PropagateBlock(block)
            end
        end
    end

    opt MiningScheduler
        opt Expired deals
            BlockchainSubsystem ->> SectorIndexerSubsystem: OnNewTipset(Chain, Epoch)
            SectorIndexerSubsystem ->> SectorIndexerSubsystem: [SectorID] ← LookupSectorByDealExpiry(Epoch)
            SectorIndexerSubsystem ->> SectorIndexerSubsystem: PurgeSectorsWithNoLiveDeals([SectorID])
            SectorIndexerSubsystem ->>+ StorageMiningSubsystem: HandleExpiredDeals([SectorID])
            StorageMiningSubsystem ->>- StorageMarketActor: CloseExpiredStorageDeal(StorageDeal)
        end
        Note Right of MiningScheduler: Schedule and resume PoSts
        Note Right of MiningScheduler: Schedule and resume SEALs
        Note Right of MiningScheduler: Maintain FaultSet
        Note Right of MiningScheduler: Maintain DoneSet
    end

    opt ClockSubsystem
        Note Right of ClockSubsystem: PoSt challenge if a node has not mined a block in a ProvingPeriod
        ClockSubsystem ->>+ StoragePowerConsensusSubsystem: GetPoStChallenge(Chain, Epoch)
        StoragePowerConsensusSubsystem -->>- StorageMiningSubsystem: challenge
        StorageMiningSubsystem ->>+ StorageProvingSubsystem: GeneratePoSt(challenge, [SectorID])
        StorageProvingSubsystem -->>- StorageMiningSubsystem: PoStProof
        opt SubmitPoSt
            StorageMiningSubsystem ->> StorageMinerActor: SubmitPoSt(PoStProof, DoneSet)
            StorageMinerActor ->>+ FilecoinProofsSubsystem: VerifyPoSt(PoStProof)
            FilecoinProofsSubsystem -->>- StorageMinerActor: {1,0}
            alt 1 - success
                StorageMinerActor ->> StorageMinerActor:  UpdateDoneSet()
                StorageMinerActor ->> StorageMarketActor: HandleStorageDealPayment()
                StorageMarketActor ->> StorageMarketActor: CloseExpiredStorageDeal()
            else 0 - failure
                StorageMinerActor -->> StorageMiningSubsystem: PoStError
                StorageMarketActor -->> StorageMarketActor: SlashStorageCollateral()
            end
        end
    end

    opt Storage Fault
        opt Declaration before receiving a challenge
            StorageMinerSubsystem ->> StorageMinerActor: UpdateSectorStatus([FaultSet], SectorStateSets)
            StorageMinerActor ->> StoragePowerActor: RecomputeMinerPower()
            StorageMarketActor -->> StorageMarketActor: SlashStorageCollateral()
            Note Right of StorageMinerActor: SectorStateSets := (FaultSet, RecoverSet, ExpireSet)
        end

        loop EveryBlock
            loop forEach StorageMinerActor in StoragePowerActor.Miners
                opt if miner ProvingPeriod ends
                    StoragePowerActor ->> StorageMinerActor: ProvingPeriodUpdate()
                    StorageMinerActor ->> StorageMinerActor: computeProvingPeriodEndSectorState()
                    Note Right of StorageMinerActor: FaultSet is all sectors if no post submitted
                    Note Right of StorageMinerActor: sectors Faulted longer than threshold proving periods are destroyed
                    StorageMinerActor ->> StorageMinerActor: UpdateSectorStatus(newSectorState)
                    StorageMinerActor ->> StoragePowerActor: RecomputeMinerPower()
                    StorageMinerActor ->> StorageMarketActor: HandleFailedDeals([newSectorState.DestroyedSet])
                end
            end
        end
    end
{{< /mermaid >}}

<script>
  setTimeout(() => {
    var element = document.querySelector('.mermaid > svg')
    console.log("element", element)
    panzoom(element, {  bounds: true })
    // svgPanZoom('.mermaid > svg', { 
    //   // contain: true, 
    //   // fit: false, 
    //   controlIconsEnabled: true 
    //   })
  }, 100)
</script>

{{</* diagram src="../diagrams/sequence/full-deals-on-chain.mmd.svg" title="Protocol Sequence Diagram - Deals on Chain" >}}

# Parameter Calculation Dependency Graph

This is a diagram of the model for parameter calculation. This is made with [orient](https://github.com/filecoin-project/orient), our tool for modeling and solving for constraints.

{{</* diagram src="../diagrams/orient/filecoin.dot.svg" title="Parameter Calculation Dependency Graph" */>}}

