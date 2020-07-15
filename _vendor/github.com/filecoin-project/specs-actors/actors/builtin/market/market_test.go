package market_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	"github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/market"
	"github.com/filecoin-project/specs-actors/actors/builtin/miner"
	"github.com/filecoin-project/specs-actors/actors/builtin/verifreg"
	"github.com/filecoin-project/specs-actors/actors/crypto"
	"github.com/filecoin-project/specs-actors/actors/runtime"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	"github.com/filecoin-project/specs-actors/actors/util/adt"
	"github.com/filecoin-project/specs-actors/support/mock"
	tutil "github.com/filecoin-project/specs-actors/support/testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustCbor(o runtime.CBORMarshaler) []byte {
	buf := new(bytes.Buffer)
	if err := o.MarshalCBOR(buf); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func TestExports(t *testing.T) {
	mock.CheckActorExports(t, market.Actor{})
}

func TestRemoveAllError(t *testing.T) {
	marketActor := tutil.NewIDAddr(t, 100)
	builder := mock.NewBuilder(context.Background(), marketActor)
	rt := builder.Build(t)
	store := adt.AsStore(rt)

	smm := market.MakeEmptySetMultimap(store)

	if err := smm.RemoveAll(42); err != nil {
		t.Fatalf("expected no error, got: %s", err)
	}
}

func TestMarketActor(t *testing.T) {
	owner := tutil.NewIDAddr(t, 101)
	provider := tutil.NewIDAddr(t, 102)
	worker := tutil.NewIDAddr(t, 103)
	client := tutil.NewIDAddr(t, 104)
	minerAddrs := &minerAddrs{owner, worker, provider}

	var st market.State

	t.Run("simple construction", func(t *testing.T) {
		actor := market.Actor{}
		receiver := tutil.NewIDAddr(t, 100)
		builder := mock.NewBuilder(context.Background(), receiver).
			WithCaller(builtin.SystemActorAddr, builtin.InitActorCodeID)

		rt := builder.Build(t)

		rt.ExpectValidateCallerAddr(builtin.SystemActorAddr)

		ret := rt.Call(actor.Constructor, nil).(*adt.EmptyValue)
		assert.Nil(t, ret)
		rt.Verify()

		store := adt.AsStore(rt)

		emptyMap, err := adt.MakeEmptyMap(store).Root()
		assert.NoError(t, err)

		emptyArray, err := adt.MakeEmptyArray(store).Root()
		assert.NoError(t, err)

		emptyMultiMap, err := market.MakeEmptySetMultimap(store).Root()
		assert.NoError(t, err)

		var state market.State
		rt.GetState(&state)

		assert.Equal(t, emptyArray, state.Proposals)
		assert.Equal(t, emptyArray, state.States)
		assert.Equal(t, emptyMap, state.EscrowTable)
		assert.Equal(t, emptyMap, state.LockedTable)
		assert.Equal(t, abi.DealID(0), state.NextID)
		assert.Equal(t, emptyMultiMap, state.DealOpsByEpoch)
		assert.Equal(t, abi.ChainEpoch(-1), state.LastCron)
	})

	t.Run("AddBalance", func(t *testing.T) {
		t.Run("adds to provider escrow funds", func(t *testing.T) {
			testCases := []struct {
				delta int64
				total int64
			}{
				{10, 10},
				{20, 30},
				{40, 70},
			}

			// Test adding provider funds from both worker and owner address
			for _, callerAddr := range []address.Address{owner, worker} {
				rt, actor := basicMarketSetup(t, owner, provider, worker, client)

				for _, tc := range testCases {
					rt.SetCaller(callerAddr, builtin.AccountActorCodeID)
					rt.SetReceived(abi.NewTokenAmount(tc.delta))
					actor.expectProviderControlAddressesAndValidateCaller(rt, provider, owner, worker)

					rt.Call(actor.AddBalance, &provider)

					rt.Verify()

					rt.GetState(&st)
					assert.Equal(t, abi.NewTokenAmount(tc.total), st.GetEscrowBalance(rt, provider))
				}
			}
		})

		t.Run("fails unless called by an account actor", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)

			rt.SetReceived(abi.NewTokenAmount(10))
			actor.expectProviderControlAddressesAndValidateCaller(rt, provider, owner, worker)

			rt.SetCaller(provider, builtin.StorageMinerActorCodeID)
			rt.ExpectAbort(exitcode.ErrForbidden, func() {
				rt.Call(actor.AddBalance, &provider)
			})

			rt.Verify()
		})

		t.Run("adds to non-provider escrow funds", func(t *testing.T) {
			testCases := []struct {
				delta int64
				total int64
			}{
				{10, 10},
				{20, 30},
				{40, 70},
			}

			// Test adding non-provider funds from both worker and client addresses
			for _, callerAddr := range []address.Address{client, worker} {
				rt, actor := basicMarketSetup(t, owner, provider, worker, client)

				for _, tc := range testCases {
					rt.SetCaller(callerAddr, builtin.AccountActorCodeID)
					rt.SetReceived(abi.NewTokenAmount(tc.delta))
					rt.ExpectValidateCallerType(builtin.CallerTypesSignable...)

					rt.Call(actor.AddBalance, &callerAddr)

					rt.Verify()

					rt.GetState(&st)
					assert.Equal(t, abi.NewTokenAmount(tc.total), st.GetEscrowBalance(rt, callerAddr))
				}
			}
		})
	})

	t.Run("WithdrawBalance", func(t *testing.T) {
		startEpoch := abi.ChainEpoch(10)
		endEpoch := abi.ChainEpoch(20)
		publishEpoch := abi.ChainEpoch(5)

		t.Run("fails with a negative withdraw amount", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)

			params := market.WithdrawBalanceParams{
				ProviderOrClientAddress: provider,
				Amount:                  abi.NewTokenAmount(-1),
			}

			rt.ExpectAbort(exitcode.ErrIllegalArgument, func() {
				rt.Call(actor.WithdrawBalance, &params)
			})

			rt.Verify()
		})

		t.Run("withdraws from provider escrow funds and sends to owner", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)

			actor.addProviderFunds(rt, abi.NewTokenAmount(20), minerAddrs)

			rt.GetState(&st)
			assert.Equal(t, abi.NewTokenAmount(20), st.GetEscrowBalance(rt, provider))

			// worker calls WithdrawBalance, balance is transferred to owner
			withdrawAmount := abi.NewTokenAmount(1)
			actor.withdrawProviderBalance(rt, withdrawAmount, withdrawAmount, minerAddrs)

			rt.GetState(&st)
			assert.Equal(t, abi.NewTokenAmount(19), st.GetEscrowBalance(rt, provider))
		})

		t.Run("withdraws from non-provider escrow funds", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)
			actor.addParticipantFunds(rt, client, abi.NewTokenAmount(20))

			rt.GetState(&st)
			assert.Equal(t, abi.NewTokenAmount(20), st.GetEscrowBalance(rt, client))

			withdrawAmount := abi.NewTokenAmount(1)
			actor.withdrawClientBalance(rt, client, withdrawAmount, withdrawAmount)

			rt.GetState(&st)
			assert.Equal(t, abi.NewTokenAmount(19), st.GetEscrowBalance(rt, client))
		})

		t.Run("client withdrawing more than escrow balance limits to available funds", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)
			actor.addParticipantFunds(rt, client, abi.NewTokenAmount(20))

			// withdraw amount greater than escrow balance
			withdrawAmount := abi.NewTokenAmount(25)
			expectedAmount := abi.NewTokenAmount(20)
			actor.withdrawClientBalance(rt, client, withdrawAmount, expectedAmount)

			rt.GetState(&st)
			assert.Equal(t, abi.NewTokenAmount(0), st.GetEscrowBalance(rt, client))
		})

		t.Run("worker withdrawing more than escrow balance limits to available funds", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)
			actor.addProviderFunds(rt, abi.NewTokenAmount(20), minerAddrs)

			rt.GetState(&st)
			assert.Equal(t, abi.NewTokenAmount(20), st.GetEscrowBalance(rt, provider))

			// withdraw amount greater than escrow balance
			withdrawAmount := abi.NewTokenAmount(25)
			actualWithdrawn := abi.NewTokenAmount(20)
			actor.withdrawProviderBalance(rt, withdrawAmount, actualWithdrawn, minerAddrs)

			rt.GetState(&st)
			assert.Equal(t, abi.NewTokenAmount(0), st.GetEscrowBalance(rt, provider))
		})

		t.Run("balance after withdrawal must ALWAYS be greater than or equal to locked amount", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)

			// create the deal to publish
			deal := actor.generateDealAndAddFunds(rt, client, minerAddrs, startEpoch, endEpoch)

			// publish the deal so that client AND provider collateral is locked
			rt.SetEpoch(publishEpoch)
			actor.publishDeals(rt, minerAddrs, deal)
			rt.GetState(&st)
			require.Equal(t, deal.ProviderCollateral, st.GetLockedBalance(rt, provider))
			require.Equal(t, deal.ClientBalanceRequirement(), st.GetLockedBalance(rt, client))

			withDrawAmt := abi.NewTokenAmount(1)
			withDrawableAmt := abi.NewTokenAmount(0)
			// client cannot withdraw any funds since all it's balance is locked
			actor.withdrawClientBalance(rt, client, withDrawAmt, withDrawableAmt)
			//  provider cannot withdraw any funds since all it's balance is locked
			actor.withdrawProviderBalance(rt, withDrawAmt, withDrawableAmt, minerAddrs)

			// add some more funds to the provider & ensure withdrawal is limited by the locked funds
			withDrawAmt = abi.NewTokenAmount(30)
			withDrawableAmt = abi.NewTokenAmount(25)
			actor.addProviderFunds(rt, withDrawableAmt, minerAddrs)
			actor.withdrawProviderBalance(rt, withDrawAmt, withDrawableAmt, minerAddrs)

			// add some more funds to the client & ensure withdrawal is limited by the locked funds
			actor.addParticipantFunds(rt, client, withDrawableAmt)
			actor.withdrawClientBalance(rt, client, withDrawAmt, withDrawableAmt)
		})

		t.Run("worker balance after withdrawal must account for slashed funds", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)

			// create the deal to publish
			deal := actor.generateDealAndAddFunds(rt, client, minerAddrs, startEpoch, endEpoch)

			// publish the deal
			rt.SetEpoch(publishEpoch)
			dealID := actor.publishDeals(rt, minerAddrs, deal)[0]

			// activate the deal
			actor.activateDeals(rt, endEpoch+1, provider, publishEpoch, dealID)
			st := actor.getDealState(rt, dealID)
			require.EqualValues(t, publishEpoch, st.SectorStartEpoch)

			// slash the deal
			newEpoch := publishEpoch + 1
			rt.SetEpoch(newEpoch)
			actor.terminateDeals(rt, provider, dealID)
			st = actor.getDealState(rt, dealID)
			require.EqualValues(t, publishEpoch+1, st.SlashEpoch)

			// provider cannot withdraw any funds since all it's balance is locked
			withDrawAmt := abi.NewTokenAmount(1)
			actualWithdrawn := abi.NewTokenAmount(0)
			actor.withdrawProviderBalance(rt, withDrawAmt, actualWithdrawn, minerAddrs)

			// add some more funds to the provider & ensure withdrawal is limited by the locked funds
			actor.addProviderFunds(rt, abi.NewTokenAmount(25), minerAddrs)
			withDrawAmt = abi.NewTokenAmount(30)
			actualWithdrawn = abi.NewTokenAmount(25)

			actor.withdrawProviderBalance(rt, withDrawAmt, actualWithdrawn, minerAddrs)
		})
	})
}

func TestPublishStorageDeals(t *testing.T) {
	owner := tutil.NewIDAddr(t, 101)
	provider := tutil.NewIDAddr(t, 102)
	worker := tutil.NewIDAddr(t, 103)
	client := tutil.NewIDAddr(t, 104)
	mAddr := &minerAddrs{owner, worker, provider}
	var st market.State

	t.Run("publish a deal after activating a previous deal which has a start epoch far in the future", func(t *testing.T) {
		startEpoch := abi.ChainEpoch(1000)
		endEpoch := abi.ChainEpoch(2000)
		publishEpoch := abi.ChainEpoch(1)

		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		deal1 := actor.generateDealAndAddFunds(rt, client, mAddr, startEpoch, endEpoch)

		// publish the deal and activate it
		rt.SetEpoch(publishEpoch)
		deal1ID := actor.publishDeals(rt, mAddr, deal1)[0]
		actor.activateDeals(rt, endEpoch, provider, publishEpoch, deal1ID)
		st := actor.getDealState(rt, deal1ID)
		require.EqualValues(t, publishEpoch, st.SectorStartEpoch)

		// now publish a second deal and activate it
		newEpoch := publishEpoch + 1
		deal2 := actor.generateDealAndAddFunds(rt, client, mAddr, startEpoch+1, endEpoch+1)
		rt.SetEpoch(newEpoch)
		deal2ID := actor.publishDeals(rt, mAddr, deal2)[0]
		actor.activateDeals(rt, endEpoch+1, provider, newEpoch, deal2ID)
	})

	t.Run("publish multiple deals for different clients and ensure balances are correct", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		client1 := tutil.NewIDAddr(t, 900)
		client2 := tutil.NewIDAddr(t, 901)
		client3 := tutil.NewIDAddr(t, 902)

		// generate first deal for
		deal1 := actor.generateDealAndAddFunds(rt, client1, mAddr, abi.ChainEpoch(42), abi.ChainEpoch(100))

		// generate second deal
		deal2 := actor.generateDealAndAddFunds(rt, client2, mAddr, abi.ChainEpoch(42), abi.ChainEpoch(100))

		// generate third deal
		deal3 := actor.generateDealAndAddFunds(rt, client3, mAddr, abi.ChainEpoch(42), abi.ChainEpoch(100))

		actor.publishDeals(rt, mAddr, deal1, deal2, deal3)

		// assert locked balance for all clients and provider
		providerLocked := big.Sum(deal1.ProviderCollateral, deal2.ProviderCollateral, deal3.ProviderCollateral)
		client1Locked := actor.getLockedBalance(rt, client1)
		client2Locked := actor.getLockedBalance(rt, client2)
		client3Locked := actor.getLockedBalance(rt, client3)
		require.EqualValues(t, deal1.ClientBalanceRequirement(), client1Locked)
		require.EqualValues(t, deal2.ClientBalanceRequirement(), client2Locked)
		require.EqualValues(t, deal3.ClientBalanceRequirement(), client3Locked)
		require.EqualValues(t, providerLocked, actor.getLockedBalance(rt, provider))

		// assert locked funds states
		rt.GetState(&st)
		totalClientCollateralLocked := big.Sum(deal3.ClientCollateral, deal1.ClientCollateral, deal2.ClientCollateral)
		require.EqualValues(t, totalClientCollateralLocked, st.TotalClientLockedCollateral)
		require.EqualValues(t, providerLocked, st.TotalProviderLockedCollateral)
		totalStorageFee := big.Sum(deal1.TotalStorageFee(), deal2.TotalStorageFee(), deal3.TotalStorageFee())
		require.EqualValues(t, totalStorageFee, st.TotalClientStorageFee)

		// publish two more deals for same clients with same provider
		deal4 := actor.generateDealAndAddFunds(rt, client3, mAddr, abi.ChainEpoch(1000), abi.ChainEpoch(10000))
		deal5 := actor.generateDealAndAddFunds(rt, client3, mAddr, abi.ChainEpoch(100), abi.ChainEpoch(1000))
		actor.publishDeals(rt, mAddr, deal4, deal5)

		// assert locked balances for clients and provider
		rt.GetState(&st)
		providerLocked = big.Sum(providerLocked, deal4.ProviderCollateral, deal5.ProviderCollateral)
		require.EqualValues(t, providerLocked, actor.getLockedBalance(rt, provider))

		client3LockedUpdated := actor.getLockedBalance(rt, client3)
		require.EqualValues(t, big.Sum(client3Locked, deal4.ClientBalanceRequirement(), deal5.ClientBalanceRequirement()), client3LockedUpdated)

		client1Locked = actor.getLockedBalance(rt, client1)
		client2Locked = actor.getLockedBalance(rt, client2)
		require.EqualValues(t, deal1.ClientBalanceRequirement(), client1Locked)
		require.EqualValues(t, deal2.ClientBalanceRequirement(), client2Locked)

		// assert locked funds states
		totalClientCollateralLocked = big.Sum(totalClientCollateralLocked, deal4.ClientCollateral, deal5.ClientCollateral)
		require.EqualValues(t, totalClientCollateralLocked, st.TotalClientLockedCollateral)
		require.EqualValues(t, providerLocked, st.TotalProviderLockedCollateral)

		totalStorageFee = big.Sum(totalStorageFee, deal4.TotalStorageFee(), deal5.TotalStorageFee())
		require.EqualValues(t, totalStorageFee, st.TotalClientStorageFee)

		// PUBLISH DEALS with a different provider
		provider2 := tutil.NewIDAddr(t, 109)
		miner := &minerAddrs{owner, worker, provider2}

		// generate first deal for second provider
		deal6 := actor.generateDealAndAddFunds(rt, client1, miner, abi.ChainEpoch(20), abi.ChainEpoch(50))

		// generate second deal for second provider
		deal7 := actor.generateDealAndAddFunds(rt, client1, miner, abi.ChainEpoch(25), abi.ChainEpoch(60))

		// publish both the deals for the second provider
		actor.publishDeals(rt, miner, deal6, deal7)

		// assertions
		rt.GetState(&st)
		provider2Locked := big.Add(deal6.ProviderCollateral, deal7.ProviderCollateral)
		require.EqualValues(t, provider2Locked, actor.getLockedBalance(rt, provider2))
		client1LockedUpdated := actor.getLockedBalance(rt, client1)
		require.EqualValues(t, big.Add(deal7.ClientBalanceRequirement(), big.Add(client1Locked, deal6.ClientBalanceRequirement())), client1LockedUpdated)

		// assert first provider's balance as well
		require.EqualValues(t, providerLocked, actor.getLockedBalance(rt, provider))

		totalClientCollateralLocked = big.Add(totalClientCollateralLocked, big.Add(deal6.ClientCollateral, deal7.ClientCollateral))
		require.EqualValues(t, totalClientCollateralLocked, st.TotalClientLockedCollateral)
		require.EqualValues(t, big.Add(providerLocked, provider2Locked), st.TotalProviderLockedCollateral)
		totalStorageFee = big.Add(totalStorageFee, big.Add(deal6.TotalStorageFee(), deal7.TotalStorageFee()))
		require.EqualValues(t, totalStorageFee, st.TotalClientStorageFee)
	})
}

func TestPublishStorageDealsFailures(t *testing.T) {
	owner := tutil.NewIDAddr(t, 101)
	provider := tutil.NewIDAddr(t, 102)
	worker := tutil.NewIDAddr(t, 103)
	client := tutil.NewIDAddr(t, 104)
	mAddrs := &minerAddrs{owner, worker, provider}

	currentEpoch := abi.ChainEpoch(5)
	startEpoch := abi.ChainEpoch(10)
	endEpoch := abi.ChainEpoch(20)

	// simple failures because of invalid deal params
	{
		tcs := map[string]struct {
			setup                      func(*mock.Runtime, *marketActorTestHarness, *market.DealProposal)
			exitCode                   exitcode.ExitCode
			signatureVerificationError error
		}{
			"deal end after deal start": {
				setup: func(_ *mock.Runtime, _ *marketActorTestHarness, d *market.DealProposal) {
					d.StartEpoch = 10
					d.EndEpoch = 9
				},
				exitCode: exitcode.ErrIllegalArgument,
			},
			"current epoch greater than start epoch": {
				setup: func(_ *mock.Runtime, _ *marketActorTestHarness, d *market.DealProposal) {
					d.StartEpoch = currentEpoch - 1
				},
				exitCode: exitcode.ErrIllegalArgument,
			},
			"deal duration greater than max deal duration": {
				setup: func(_ *mock.Runtime, _ *marketActorTestHarness, d *market.DealProposal) {
					d.StartEpoch = abi.ChainEpoch(10)
					d.EndEpoch = d.StartEpoch + (1 * builtin.EpochsInYear) + 1
				},
				exitCode: exitcode.ErrIllegalArgument,
			},
			"negative price per epoch": {
				setup: func(_ *mock.Runtime, _ *marketActorTestHarness, d *market.DealProposal) {
					d.StoragePricePerEpoch = abi.NewTokenAmount(-1)
				},
				exitCode: exitcode.ErrIllegalArgument,
			},
			"price per epoch greater than total filecoin": {
				setup: func(_ *mock.Runtime, _ *marketActorTestHarness, d *market.DealProposal) {
					d.StoragePricePerEpoch = big.Add(abi.TotalFilecoin, big.NewInt(1))
				},
				exitCode: exitcode.ErrIllegalArgument,
			},
			"negative provider collateral": {
				setup: func(_ *mock.Runtime, _ *marketActorTestHarness, d *market.DealProposal) {
					d.ProviderCollateral = big.NewInt(-1)
				},
				exitCode: exitcode.ErrIllegalArgument,
			},
			"provider collateral greater than max collateral": {
				setup: func(_ *mock.Runtime, _ *marketActorTestHarness, d *market.DealProposal) {
					d.ProviderCollateral = big.Add(abi.TotalFilecoin, big.NewInt(1))
				},
				exitCode: exitcode.ErrIllegalArgument,
			},
			"negative client collateral": {
				setup: func(_ *mock.Runtime, _ *marketActorTestHarness, d *market.DealProposal) {
					d.ClientCollateral = big.NewInt(-1)
				},
				exitCode: exitcode.ErrIllegalArgument,
			},
			"client collateral greater than max collateral": {
				setup: func(_ *mock.Runtime, _ *marketActorTestHarness, d *market.DealProposal) {
					d.ClientCollateral = big.Add(abi.TotalFilecoin, big.NewInt(1))
				},
				exitCode: exitcode.ErrIllegalArgument,
			},
			"client does not have enough balance for collateral": {
				setup: func(rt *mock.Runtime, a *marketActorTestHarness, d *market.DealProposal) {
					a.addParticipantFunds(rt, client, big.Sub(d.ClientBalanceRequirement(), big.NewInt(1)))
					a.addProviderFunds(rt, d.ProviderCollateral, mAddrs)
				},
				exitCode: exitcode.ErrInsufficientFunds,
			},
			"provider does not have enough balance for collateral": {
				setup: func(rt *mock.Runtime, a *marketActorTestHarness, d *market.DealProposal) {
					a.addParticipantFunds(rt, client, d.ClientBalanceRequirement())
					a.addProviderFunds(rt, big.Sub(d.ProviderCollateral, big.NewInt(1)), mAddrs)
				},
				exitCode: exitcode.ErrInsufficientFunds,
			},
			"unable to resolve client address": {
				setup: func(_ *mock.Runtime, a *marketActorTestHarness, d *market.DealProposal) {
					d.Client = tutil.NewBLSAddr(t, 1)
				},
				exitCode: exitcode.ErrNotFound,
			},
			"signature is invalid": {
				setup: func(_ *mock.Runtime, a *marketActorTestHarness, d *market.DealProposal) {

				},
				exitCode:                   exitcode.ErrIllegalArgument,
				signatureVerificationError: errors.New("error"),
			},
			"no entry for client in locked  balance table": {
				setup: func(rt *mock.Runtime, a *marketActorTestHarness, d *market.DealProposal) {
					a.addProviderFunds(rt, d.ProviderCollateral, mAddrs)
				},
				exitCode: exitcode.ErrInsufficientFunds,
			},
			"no entry for provider in locked  balance table": {
				setup: func(rt *mock.Runtime, a *marketActorTestHarness, d *market.DealProposal) {
					a.addParticipantFunds(rt, client, d.ClientBalanceRequirement())
				},
				exitCode: exitcode.ErrInsufficientFunds,
			},
		}

		for name, tc := range tcs {
			t.Run(name, func(t *testing.T) {
				rt, actor := basicMarketSetup(t, owner, provider, worker, client)
				dealProposal := generateDealProposal(client, provider, startEpoch, endEpoch)
				rt.SetEpoch(currentEpoch)
				tc.setup(rt, actor, &dealProposal)
				params := mkPublishStorageParams(dealProposal)

				rt.ExpectValidateCallerType(builtin.AccountActorCodeID, builtin.MultisigActorCodeID)
				rt.ExpectSend(provider, builtin.MethodsMiner.ControlAddresses, nil, abi.NewTokenAmount(0), &miner.GetControlAddressesReturn{Worker: worker, Owner: owner}, 0)
				rt.SetCaller(worker, builtin.AccountActorCodeID)
				rt.ExpectVerifySignature(crypto.Signature{}, dealProposal.Client, mustCbor(&dealProposal), tc.signatureVerificationError)
				rt.ExpectAbort(tc.exitCode, func() {
					rt.Call(actor.PublishStorageDeals, params)
				})

				rt.Verify()
			})
		}
	}

	// fails when client or provider has some funds but not enough to cover a deal
	{
		t.Run("fail when client has some funds but not enough for a deal", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)

			//
			actor.addParticipantFunds(rt, client, abi.NewTokenAmount(100))
			deal1 := generateDealProposal(client, provider, abi.ChainEpoch(42), abi.ChainEpoch(100))
			actor.addProviderFunds(rt, deal1.ProviderCollateral, mAddrs)
			params := mkPublishStorageParams(deal1)

			rt.ExpectValidateCallerType(builtin.AccountActorCodeID, builtin.MultisigActorCodeID)
			rt.ExpectSend(provider, builtin.MethodsMiner.ControlAddresses, nil, abi.NewTokenAmount(0), &miner.GetControlAddressesReturn{Worker: worker, Owner: owner}, 0)
			rt.SetCaller(worker, builtin.AccountActorCodeID)
			rt.ExpectVerifySignature(crypto.Signature{}, deal1.Client, mustCbor(&deal1), nil)
			rt.ExpectAbort(exitcode.ErrInsufficientFunds, func() {
				rt.Call(actor.PublishStorageDeals, params)
			})

			rt.Verify()
		})

		t.Run("fail when provider has some funds but not enough for a deal", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)

			actor.addProviderFunds(rt, abi.NewTokenAmount(1), mAddrs)
			deal1 := generateDealProposal(client, provider, abi.ChainEpoch(42), abi.ChainEpoch(100))
			actor.addParticipantFunds(rt, client, deal1.ClientBalanceRequirement())

			params := mkPublishStorageParams(deal1)

			rt.ExpectValidateCallerType(builtin.AccountActorCodeID, builtin.MultisigActorCodeID)
			rt.ExpectSend(provider, builtin.MethodsMiner.ControlAddresses, nil, abi.NewTokenAmount(0), &miner.GetControlAddressesReturn{Worker: worker, Owner: owner}, 0)
			rt.SetCaller(worker, builtin.AccountActorCodeID)
			rt.ExpectVerifySignature(crypto.Signature{}, deal1.Client, mustCbor(&deal1), nil)
			rt.ExpectAbort(exitcode.ErrInsufficientFunds, func() {
				rt.Call(actor.PublishStorageDeals, params)
			})

			rt.Verify()
		})
	}

	// fail when deals have different providers
	{
		t.Run("fail when deals have different providers", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)
			deal1 := actor.generateDealAndAddFunds(rt, client, mAddrs, abi.ChainEpoch(42), abi.ChainEpoch(100))
			m2 := &minerAddrs{owner, worker, tutil.NewIDAddr(t, 1000)}

			deal2 := actor.generateDealAndAddFunds(rt, client, m2, abi.ChainEpoch(1), abi.ChainEpoch(5))

			params := mkPublishStorageParams(deal1, deal2)

			rt.ExpectValidateCallerType(builtin.AccountActorCodeID, builtin.MultisigActorCodeID)
			rt.ExpectSend(provider, builtin.MethodsMiner.ControlAddresses, nil, abi.NewTokenAmount(0), &miner.GetControlAddressesReturn{Worker: worker, Owner: owner}, 0)
			rt.SetCaller(worker, builtin.AccountActorCodeID)
			rt.ExpectVerifySignature(crypto.Signature{}, deal1.Client, mustCbor(&deal1), nil)
			rt.ExpectVerifySignature(crypto.Signature{}, deal2.Client, mustCbor(&deal2), nil)
			rt.ExpectAbort(exitcode.ErrIllegalArgument, func() {
				rt.Call(actor.PublishStorageDeals, params)
			})

			rt.Verify()
		})

		//  failures because of incorrect call params
		t.Run("fail when caller is not of signable type", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)
			params := mkPublishStorageParams(generateDealProposal(client, provider, abi.ChainEpoch(1), abi.ChainEpoch(5)))
			w := tutil.NewIDAddr(t, 1000)
			rt.SetCaller(w, builtin.StorageMinerActorCodeID)
			rt.ExpectValidateCallerType(builtin.AccountActorCodeID, builtin.MultisigActorCodeID)
			rt.ExpectAbort(exitcode.ErrForbidden, func() {
				rt.Call(actor.PublishStorageDeals, params)
			})
		})

		t.Run("fail when no deals in params", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)
			params := mkPublishStorageParams()
			rt.SetCaller(worker, builtin.AccountActorCodeID)
			rt.ExpectValidateCallerType(builtin.AccountActorCodeID, builtin.MultisigActorCodeID)
			rt.ExpectAbort(exitcode.ErrIllegalArgument, func() {
				rt.Call(actor.PublishStorageDeals, params)
			})
		})

		t.Run("fail to resolve provider address", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)
			deal := generateDealProposal(client, provider, abi.ChainEpoch(1), abi.ChainEpoch(5))
			deal.Provider = tutil.NewBLSAddr(t, 100)

			params := mkPublishStorageParams(deal)
			rt.SetCaller(worker, builtin.AccountActorCodeID)
			rt.ExpectValidateCallerType(builtin.AccountActorCodeID, builtin.MultisigActorCodeID)
			rt.ExpectAbort(exitcode.ErrNotFound, func() {
				rt.Call(actor.PublishStorageDeals, params)
			})
		})

		t.Run("caller is not the same as the worker address for miner", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)
			deal := generateDealProposal(client, provider, abi.ChainEpoch(1), abi.ChainEpoch(5))
			params := mkPublishStorageParams(deal)
			rt.ExpectValidateCallerType(builtin.AccountActorCodeID, builtin.MultisigActorCodeID)
			rt.ExpectSend(provider, builtin.MethodsMiner.ControlAddresses, nil, abi.NewTokenAmount(0), &miner.GetControlAddressesReturn{Worker: tutil.NewIDAddr(t, 999), Owner: owner}, 0)
			rt.SetCaller(worker, builtin.AccountActorCodeID)
			rt.ExpectAbort(exitcode.ErrForbidden, func() {
				rt.Call(actor.PublishStorageDeals, params)
			})

			rt.Verify()
		})
	}
}

func TestActivateDeals(t *testing.T) {

	owner := tutil.NewIDAddr(t, 101)
	provider := tutil.NewIDAddr(t, 102)
	worker := tutil.NewIDAddr(t, 103)
	client := tutil.NewIDAddr(t, 104)
	mAddrs := &minerAddrs{owner, worker, provider}

	startEpoch := abi.ChainEpoch(10)
	endEpoch := abi.ChainEpoch(20)
	currentEpoch := abi.ChainEpoch(5)
	sectorExpiry := abi.ChainEpoch(100)

	t.Run("active deals multiple times with different providers", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		rt.SetEpoch(currentEpoch)

		// provider 1 publishes deals1 and deals2 and deal3
		dealId1 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch)
		dealId2 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch+1)
		dealId3 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch+2)

		// provider2 publishes deal4 and deal5
		provider2 := tutil.NewIDAddr(t, 401)
		mAddrs.provider = provider2
		dealId4 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch)
		dealId5 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch+1)

		// provider1 activates deal 1 and deal2 but that does not activate deal3 to deal5
		actor.activateDeals(rt, sectorExpiry, provider, currentEpoch, dealId1, dealId2)
		actor.assertDealsNotActivated(rt, currentEpoch, dealId3, dealId4, dealId5)

		// provider3 activates deal5 but that does not activate deal3 or deal4
		actor.activateDeals(rt, sectorExpiry, provider2, currentEpoch, dealId5)
		actor.assertDealsNotActivated(rt, currentEpoch, dealId3, dealId4)

		// provider1 activates deal3
		actor.activateDeals(rt, sectorExpiry, provider, currentEpoch, dealId3)
		actor.assertDealsNotActivated(rt, currentEpoch, dealId4)
	})
}

func TestActivateDealFailures(t *testing.T) {
	owner := tutil.NewIDAddr(t, 101)
	provider := tutil.NewIDAddr(t, 102)
	worker := tutil.NewIDAddr(t, 103)
	client := tutil.NewIDAddr(t, 104)
	mAddrs := &minerAddrs{owner, worker, provider}

	startEpoch := abi.ChainEpoch(10)
	endEpoch := abi.ChainEpoch(20)
	sectorExpiry := abi.ChainEpoch(100)

	// caller is not the provider
	{
		t.Run("fail when caller is not the provider of the deal", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)
			provider2 := tutil.NewIDAddr(t, 201)
			mAddrs2 := &minerAddrs{owner, worker, provider2}
			dealId := actor.generateAndPublishDeal(rt, client, mAddrs2, startEpoch, endEpoch)

			params := mkActivateDealParams(sectorExpiry, dealId)

			rt.ExpectValidateCallerType(builtin.StorageMinerActorCodeID)
			rt.SetCaller(provider, builtin.StorageMinerActorCodeID)
			rt.ExpectAbort(exitcode.ErrIllegalState, func() {
				rt.Call(actor.ActivateDeals, params)
			})

			rt.Verify()
		})
	}

	// caller is not a StorageMinerActor
	{
		t.Run("fail when caller is not a StorageMinerActor", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)
			rt.ExpectValidateCallerType(builtin.StorageMinerActorCodeID)
			rt.SetCaller(provider, builtin.AccountActorCodeID)
			rt.ExpectAbort(exitcode.ErrForbidden, func() {
				rt.Call(actor.ActivateDeals, &market.ActivateDealsParams{})
			})

			rt.Verify()
		})
	}

	// deal has not been published before
	{
		t.Run("fail when deal has not been published before", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)
			params := mkActivateDealParams(sectorExpiry, abi.DealID(42))

			rt.ExpectValidateCallerType(builtin.StorageMinerActorCodeID)
			rt.SetCaller(provider, builtin.StorageMinerActorCodeID)
			rt.ExpectAbort(exitcode.ErrIllegalState, func() {
				rt.Call(actor.ActivateDeals, params)
			})

			rt.Verify()
		})
	}

	// deal has ALREADY been activated
	{
		t.Run("fail when deal has already been activated", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)
			dealId := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch)
			actor.activateDeals(rt, sectorExpiry, provider, 0, dealId)

			rt.ExpectValidateCallerType(builtin.StorageMinerActorCodeID)
			rt.SetCaller(provider, builtin.StorageMinerActorCodeID)
			rt.ExpectAbort(exitcode.ErrIllegalArgument, func() {
				rt.Call(actor.ActivateDeals, mkActivateDealParams(sectorExpiry, dealId))
			})

			rt.Verify()
		})
	}

	// deal has invalid params
	{
		t.Run("fail when current epoch greater than start epoch of deal", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)
			dealId := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch)

			rt.ExpectValidateCallerType(builtin.StorageMinerActorCodeID)
			rt.SetCaller(provider, builtin.StorageMinerActorCodeID)
			rt.SetEpoch(startEpoch + 1)
			rt.ExpectAbort(exitcode.ErrIllegalState, func() {
				rt.Call(actor.ActivateDeals, mkActivateDealParams(sectorExpiry, dealId))
			})

			rt.Verify()
		})

		t.Run("fail when end epoch of deal greater than sector expiry", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)
			dealId := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch)

			rt.ExpectValidateCallerType(builtin.StorageMinerActorCodeID)
			rt.SetCaller(provider, builtin.StorageMinerActorCodeID)
			rt.ExpectAbort(exitcode.ErrIllegalState, func() {
				rt.Call(actor.ActivateDeals, mkActivateDealParams(endEpoch-1, dealId))
			})

			rt.Verify()
		})
	}

	// all fail if one fails
	{
		t.Run("fail to activate all deals if one deal fails", func(t *testing.T) {
			rt, actor := basicMarketSetup(t, owner, provider, worker, client)

			// activate deal1 so it fails later
			dealId1 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch)
			actor.activateDeals(rt, sectorExpiry, provider, 0, dealId1)

			dealId2 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch+1)

			rt.ExpectValidateCallerType(builtin.StorageMinerActorCodeID)
			rt.SetCaller(provider, builtin.StorageMinerActorCodeID)
			rt.ExpectAbort(exitcode.ErrIllegalArgument, func() {
				rt.Call(actor.ActivateDeals, mkActivateDealParams(sectorExpiry, dealId1, dealId2))
			})
			rt.Verify()

			// no state for deal2 means deal2 activation has failed
			var st market.State
			rt.GetState(&st)

			states, err := market.AsDealStateArray(adt.AsStore(rt), st.States)
			require.NoError(t, err)

			_, found, err := states.Get(dealId2)
			require.NoError(t, err)
			require.False(t, found)
		})
	}

}

func TestOnMinerSectorsTerminate(t *testing.T) {
	owner := tutil.NewIDAddr(t, 101)
	provider := tutil.NewIDAddr(t, 102)
	worker := tutil.NewIDAddr(t, 103)
	client := tutil.NewIDAddr(t, 104)
	mAddrs := &minerAddrs{owner, worker, provider}

	startEpoch := abi.ChainEpoch(10)
	endEpoch := abi.ChainEpoch(20)
	currentEpoch := abi.ChainEpoch(5)
	sectorExpiry := abi.ChainEpoch(100)

	t.Run("terminate multiple deals from multiple providers", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		rt.SetEpoch(currentEpoch)

		// provider1 publishes deal1,2 and 3
		dealId1 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch)
		dealId2 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch+1)
		dealId3 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch+2)
		actor.activateDeals(rt, sectorExpiry, provider, currentEpoch, dealId1, dealId2, dealId3)

		// provider2 publishes deal4 and deal5
		provider2 := tutil.NewIDAddr(t, 501)
		maddrs2 := &minerAddrs{owner, worker, provider2}
		dealId4 := actor.generateAndPublishDeal(rt, client, maddrs2, startEpoch, endEpoch)
		dealId5 := actor.generateAndPublishDeal(rt, client, maddrs2, startEpoch, endEpoch+1)
		actor.activateDeals(rt, sectorExpiry, provider2, currentEpoch, dealId4, dealId5)

		// provider1 terminates deal1 but that does not terminate deals2-5
		actor.terminateDeals(rt, provider, dealId1)
		actor.assertDealsTerminated(rt, currentEpoch, dealId1)
		actor.assertDeaslNotTerminated(rt, dealId2, dealId3, dealId4, dealId5)

		// provider2 terminates deal5 but that does not terminate delals 2-4
		actor.terminateDeals(rt, provider2, dealId5)
		actor.assertDealsTerminated(rt, currentEpoch, dealId5)
		actor.assertDeaslNotTerminated(rt, dealId2, dealId3, dealId4)

		// provider1 terminates deal2 and deal3
		actor.terminateDeals(rt, provider, dealId2, dealId3)
		actor.assertDealsTerminated(rt, currentEpoch, dealId2, dealId3)
		actor.assertDeaslNotTerminated(rt, dealId4)

		// provider2 terminates deal4
		actor.terminateDeals(rt, provider2, dealId4)
		actor.assertDealsTerminated(rt, currentEpoch, dealId4)
	})

	t.Run("ignore deal proposal that does not exist", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		rt.SetEpoch(currentEpoch)

		// deal1 will be terminated and the other deal will be ignored because it does not exist
		dealId1 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch)
		actor.activateDeals(rt, sectorExpiry, provider, currentEpoch, dealId1)

		actor.terminateDeals(rt, provider, dealId1, abi.DealID(42))
		st := actor.getDealState(rt, dealId1)
		require.EqualValues(t, currentEpoch, st.SlashEpoch)
	})

	t.Run("terminate valid deals along with expired deals - only valid deals are terminated", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		rt.SetEpoch(currentEpoch)

		// provider1 publishes deal1 and 2 and deal3 -> deal3 has the lowest endepoch
		dealId1 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch)
		dealId2 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch+1)
		dealId3 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch-1)
		actor.activateDeals(rt, sectorExpiry, provider, currentEpoch, dealId1, dealId2, dealId3)

		// set current epoch such that deal3 expires but the other two do not
		newEpoch := endEpoch - 1
		rt.SetEpoch(newEpoch)

		// terminating all three deals ONLY terminates deal1 and deal2 because deal3 has expired
		actor.terminateDeals(rt, provider, dealId1, dealId2, dealId3)
		actor.assertDealsTerminated(rt, newEpoch, dealId1, dealId2)
		actor.assertDeaslNotTerminated(rt, dealId3)

	})

	t.Run("terminating a deal the second time does not change it's slash epoch", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		rt.SetEpoch(currentEpoch)

		dealId1 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch)
		actor.activateDeals(rt, sectorExpiry, provider, currentEpoch, dealId1)

		// terminating the deal so slash epoch is the current epoch
		actor.terminateDeals(rt, provider, dealId1)

		// set a new epoch and terminate again -> however slash epoch will still be the old epoch.
		newEpoch := currentEpoch + 1
		rt.SetEpoch(newEpoch)
		actor.terminateDeals(rt, provider, dealId1)
		st := actor.getDealState(rt, dealId1)
		require.EqualValues(t, currentEpoch, st.SlashEpoch)
	})

	t.Run("terminating new deals and an already terminated deal only terminates the new deals", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		rt.SetEpoch(currentEpoch)

		// provider1 publishes deal1 and 2 and deal3 -> deal3 has the lowest endepoch
		dealId1 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch)
		dealId2 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch+1)
		dealId3 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch-1)
		actor.activateDeals(rt, sectorExpiry, provider, currentEpoch, dealId1, dealId2, dealId3)

		// terminating the deal so slash epoch is the current epoch
		actor.terminateDeals(rt, provider, dealId1)

		// set a new epoch and terminate again -> however slash epoch will still be the old epoch.
		newEpoch := currentEpoch + 1
		rt.SetEpoch(newEpoch)
		actor.terminateDeals(rt, provider, dealId1, dealId2, dealId3)

		st := actor.getDealState(rt, dealId1)
		require.EqualValues(t, currentEpoch, st.SlashEpoch)

		st2 := actor.getDealState(rt, dealId2)
		require.EqualValues(t, newEpoch, st2.SlashEpoch)

		st3 := actor.getDealState(rt, dealId3)
		require.EqualValues(t, newEpoch, st3.SlashEpoch)
	})

	t.Run("do not terminate deal if end epoch is equal to or less than current epoch", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		rt.SetEpoch(currentEpoch)

		// deal1 has endepoch equal to current epoch when terminate is called
		dealId1 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch)
		actor.activateDeals(rt, sectorExpiry, provider, currentEpoch, dealId1)
		rt.SetEpoch(endEpoch)
		actor.terminateDeals(rt, provider, dealId1)
		actor.assertDeaslNotTerminated(rt, dealId1)

		// deal2 has end epoch less than current epoch when terminate is called
		rt.SetEpoch(currentEpoch)
		dealId2 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch+1, endEpoch)
		actor.activateDeals(rt, sectorExpiry, provider, currentEpoch, dealId2)
		rt.SetEpoch(endEpoch + 1)
		actor.terminateDeals(rt, provider, dealId2)
		actor.assertDeaslNotTerminated(rt, dealId2)
	})

	t.Run("fail when caller is not a StorageMinerActor", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		rt.ExpectValidateCallerType(builtin.StorageMinerActorCodeID)
		rt.SetCaller(provider, builtin.AccountActorCodeID)
		rt.ExpectAbort(exitcode.ErrForbidden, func() {
			rt.Call(actor.OnMinerSectorsTerminate, &market.OnMinerSectorsTerminateParams{})
		})

		rt.Verify()
	})

	t.Run("fail when caller is not the provider of the deal", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		rt.SetEpoch(currentEpoch)

		dealId := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch)
		actor.activateDeals(rt, sectorExpiry, provider, currentEpoch, dealId)

		params := mkTerminateDealParams(dealId)

		provider2 := tutil.NewIDAddr(t, 501)
		rt.ExpectValidateCallerType(builtin.StorageMinerActorCodeID)
		rt.SetCaller(provider2, builtin.StorageMinerActorCodeID)
		rt.ExpectAssertionFailure("caller is not the provider of the deal", func() {
			rt.Call(actor.OnMinerSectorsTerminate, params)
		})

		rt.Verify()
	})

	t.Run("fail when deal has been published but not activated", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		rt.SetEpoch(currentEpoch)

		dealId := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch)

		params := mkTerminateDealParams(dealId)
		rt.ExpectValidateCallerType(builtin.StorageMinerActorCodeID)
		rt.SetCaller(provider, builtin.StorageMinerActorCodeID)
		rt.ExpectAbort(exitcode.ErrIllegalArgument, func() {
			rt.Call(actor.OnMinerSectorsTerminate, params)
		})

		rt.Verify()
	})

	t.Run("termination of all deals should fail when one deal fails", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		rt.SetEpoch(currentEpoch)

		// deal1 would terminate but deal2 will fail because deal2 has not been activated
		dealId1 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch)
		actor.activateDeals(rt, sectorExpiry, provider, currentEpoch, dealId1)
		dealId2 := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch+1)

		params := mkTerminateDealParams(dealId1, dealId2)
		rt.ExpectValidateCallerType(builtin.StorageMinerActorCodeID)
		rt.SetCaller(provider, builtin.StorageMinerActorCodeID)
		rt.ExpectAbort(exitcode.ErrIllegalArgument, func() {
			rt.Call(actor.OnMinerSectorsTerminate, params)
		})

		rt.Verify()

		// verify deal1 has not been terminated
		actor.assertDeaslNotTerminated(rt, dealId1)
	})
}

func TestCronTick(t *testing.T) {
	owner := tutil.NewIDAddr(t, 101)
	provider := tutil.NewIDAddr(t, 102)
	worker := tutil.NewIDAddr(t, 103)
	client := tutil.NewIDAddr(t, 104)
	mAddrs := &minerAddrs{owner, worker, provider}

	startEpoch := abi.ChainEpoch(50)
	endEpoch := abi.ChainEpoch(300)
	sectorExpiry := abi.ChainEpoch(400)

	t.Run("fail when deal is activated but proposal is not found", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		dealId := actor.publishAndActivateDeal(rt, client, mAddrs, startEpoch, endEpoch, 0, sectorExpiry)

		// delete the deal proposal
		actor.deleteDealProposal(rt, dealId)

		// move the current epoch to the start epoch of the deal
		rt.SetEpoch(startEpoch)
		rt.ExpectAbort(exitcode.ErrIllegalState, func() {
			actor.cronTick(rt)
		})
	})

	t.Run("fail when deal update epoch is in the future", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		dealId := actor.publishAndActivateDeal(rt, client, mAddrs, startEpoch, endEpoch, 0, sectorExpiry)

		// move the current epoch such that the deal's last updated field is set to the start epoch of the deal
		// and the next tick for it is scheduled at the endepoch.
		rt.SetEpoch(startEpoch)
		actor.cronTick(rt)

		// update last updated to some time in the future
		actor.updateLastUpdated(rt, dealId, startEpoch+1000)

		// set current epoch of the deal to the end epoch so it's picked up for "processing" in the next cron tick.
		rt.SetEpoch(endEpoch)

		rt.ExpectAssertionFailure("assertion failed", func() {
			actor.cronTick(rt)
		})
	})

	t.Run("crontick for a deal at it's start epoch results in zero payment and no slashing", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		dealId := actor.publishAndActivateDeal(rt, client, mAddrs, startEpoch, endEpoch, 0, sectorExpiry)

		// move the current epoch to startEpoch
		current := startEpoch
		rt.SetEpoch(current)
		pay, slashed := actor.cronTickAndAssertBalances(rt, client, provider, current, dealId)
		require.EqualValues(t, big.Zero(), pay)
		require.EqualValues(t, big.Zero(), slashed)

		// deal proposal and state should NOT be deleted
		require.NotNil(t, actor.getDealProposal(rt, dealId))
		require.NotNil(t, actor.getDealState(rt, dealId))
	})
}

func TestLockedFundTrackingStates(t *testing.T) {
	owner := tutil.NewIDAddr(t, 101)
	worker := tutil.NewIDAddr(t, 103)

	p1 := tutil.NewIDAddr(t, 201)
	p2 := tutil.NewIDAddr(t, 202)
	p3 := tutil.NewIDAddr(t, 203)

	c1 := tutil.NewIDAddr(t, 104)
	c2 := tutil.NewIDAddr(t, 105)
	c3 := tutil.NewIDAddr(t, 106)

	m1 := &minerAddrs{owner, worker, p1}
	m2 := &minerAddrs{owner, worker, p2}
	m3 := &minerAddrs{owner, worker, p3}

	startEpoch := abi.ChainEpoch(50)
	endEpoch := abi.ChainEpoch(300)
	sectorExpiry := abi.ChainEpoch(400)

	var st market.State

	// assert values are zero
	rt, actor := basicMarketSetup(t, owner, p1, worker, c1)
	rt.GetState(&st)
	require.True(t, st.TotalClientLockedCollateral.IsZero())
	require.True(t, st.TotalProviderLockedCollateral.IsZero())
	require.True(t, st.TotalClientStorageFee.IsZero())

	// Publish deal1, deal2 and deal3  with different client and provider
	dealId1 := actor.generateAndPublishDeal(rt, c1, m1, startEpoch, endEpoch)
	d1 := actor.getDealProposal(rt, dealId1)

	dealId2 := actor.generateAndPublishDeal(rt, c2, m2, startEpoch, endEpoch)
	d2 := actor.getDealProposal(rt, dealId2)

	dealId3 := actor.generateAndPublishDeal(rt, c3, m3, startEpoch, endEpoch)
	d3 := actor.getDealProposal(rt, dealId3)

	csf := big.Sum(d1.TotalStorageFee(), d2.TotalStorageFee(), d3.TotalStorageFee())
	plc := big.Sum(d1.ProviderCollateral, d2.ProviderCollateral, d3.ProviderCollateral)
	clc := big.Sum(d1.ClientCollateral, d2.ClientCollateral, d3.ClientCollateral)

	actor.assertLockedFundStates(rt, csf, plc, clc)

	// activation dosen't change anything
	curr := startEpoch - 1
	rt.SetEpoch(curr)
	actor.activateDeals(rt, sectorExpiry, p1, curr, dealId1)
	actor.activateDeals(rt, sectorExpiry, p2, curr, dealId2)

	actor.assertLockedFundStates(rt, csf, plc, clc)

	// make payment for p1 and p2, p3 times out as it has not been activated
	curr = 51 // startEpoch + 1
	rt.SetEpoch(curr)
	rt.ExpectSend(builtin.BurntFundsActorAddr, builtin.MethodSend, nil, d3.ProviderCollateral, nil, exitcode.Ok)
	actor.cronTick(rt)
	payment := big.Product(big.NewInt(2), d1.StoragePricePerEpoch)
	csf = big.Sub(big.Sub(csf, payment), d3.TotalStorageFee())
	plc = big.Sub(plc, d3.ProviderCollateral)
	clc = big.Sub(clc, d3.ClientCollateral)
	actor.assertLockedFundStates(rt, csf, plc, clc)

	// deal1 and deal2 will now be charged at epoch = 51 + 100 = 151, so nothing changes before that
	rt.SetEpoch(150)
	actor.cronTick(rt)
	actor.assertLockedFundStates(rt, csf, plc, clc)

	// one more round of payment for deal1 and deal2
	rt.SetEpoch(200)
	payment = big.Product(big.NewInt(2), d1.StoragePricePerEpoch, big.NewInt(149)) // 149 = 200 - 51 (epoch when payment was last made)
	csf = big.Sub(csf, payment)
	actor.cronTick(rt)
	actor.assertLockedFundStates(rt, csf, plc, clc)

	// slash deal1 at 201
	rt.SetEpoch(201)
	actor.terminateDeals(rt, m1.provider, dealId1)

	// cron tick at 300 to slash deal1 and expire deal2
	rt.SetEpoch(300)
	csf = big.Zero()
	clc = big.Zero()
	plc = big.Zero()
	rt.ExpectSend(builtin.BurntFundsActorAddr, builtin.MethodSend, nil, d1.ProviderCollateral, nil, exitcode.Ok)
	actor.cronTick(rt)
	actor.assertLockedFundStates(rt, csf, plc, clc)
}

func TestCronTickTimedoutDeals(t *testing.T) {
	owner := tutil.NewIDAddr(t, 101)
	provider := tutil.NewIDAddr(t, 102)
	worker := tutil.NewIDAddr(t, 103)
	client := tutil.NewIDAddr(t, 104)
	mAddrs := &minerAddrs{owner, worker, provider}

	startEpoch := abi.ChainEpoch(50)
	endEpoch := abi.ChainEpoch(300)

	t.Run("timed out deal is slashed and deleted", func(t *testing.T) {
		// publish a deal but do NOT activate it
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		dealId := actor.generateAndPublishDeal(rt, client, mAddrs, startEpoch, endEpoch)
		d := actor.getDealProposal(rt, dealId)

		cEscrow := actor.getEscrowBalance(rt, client)

		// do a cron tick for it -> should time out and get slashed
		rt.SetEpoch(startEpoch)
		rt.ExpectSend(builtin.BurntFundsActorAddr, builtin.MethodSend, nil, d.ProviderCollateral, nil, exitcode.Ok)
		actor.cronTick(rt)

		require.Equal(t, cEscrow, actor.getEscrowBalance(rt, client))
		require.Equal(t, big.Zero(), actor.getLockedBalance(rt, client))
		require.Equal(t, big.Zero(), actor.getEscrowBalance(rt, provider))
		require.Equal(t, big.Zero(), actor.getLockedBalance(rt, provider))

		actor.assertDealDeleted(rt, dealId)
	})

	t.Run("timed out and verified deals are slashed, deleted AND sent to the Registry actor", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		// deal1 and deal2 are verified
		deal1 := actor.generateDealAndAddFunds(rt, client, mAddrs, startEpoch, endEpoch)
		deal1.VerifiedDeal = true
		deal2 := actor.generateDealAndAddFunds(rt, client, mAddrs, startEpoch, endEpoch+1)
		deal2.VerifiedDeal = true

		// deal3 is NOT verified
		deal3 := actor.generateDealAndAddFunds(rt, client, mAddrs, startEpoch, endEpoch+2)

		//  publishing verified deals
		dealIds := actor.publishDeals(rt, mAddrs, deal1, deal2, deal3)

		// do a cron tick for it -> all should time out and get slashed
		// ONLY deal1 and deal2 should be sent to the Registry actor
		rt.SetEpoch(startEpoch)

		// expected sends to the registry actor
		param1 := &verifreg.RestoreBytesParams{
			Address:  deal1.Client,
			DealSize: big.NewIntUnsigned(uint64(deal1.PieceSize)),
		}
		param2 := &verifreg.RestoreBytesParams{
			Address:  deal2.Client,
			DealSize: big.NewIntUnsigned(uint64(deal2.PieceSize)),
		}

		rt.ExpectSend(builtin.VerifiedRegistryActorAddr, builtin.MethodsVerifiedRegistry.RestoreBytes, param1,
			abi.NewTokenAmount(0), nil, exitcode.Ok)
		rt.ExpectSend(builtin.VerifiedRegistryActorAddr, builtin.MethodsVerifiedRegistry.RestoreBytes, param2,
			abi.NewTokenAmount(0), nil, exitcode.Ok)

		expectedBurn := big.Mul(big.NewInt(3), deal1.ProviderCollateral)
		rt.ExpectSend(builtin.BurntFundsActorAddr, builtin.MethodSend, nil, expectedBurn, nil, exitcode.Ok)
		actor.cronTick(rt)

		actor.assertDealDeleted(rt, dealIds[0])
		actor.assertDealDeleted(rt, dealIds[1])
		actor.assertDealDeleted(rt, dealIds[2])
	})
}

func TestCronTickDealExpiry(t *testing.T) {
	owner := tutil.NewIDAddr(t, 101)
	provider := tutil.NewIDAddr(t, 102)
	worker := tutil.NewIDAddr(t, 103)
	client := tutil.NewIDAddr(t, 104)
	mAddrs := &minerAddrs{owner, worker, provider}

	startEpoch := abi.ChainEpoch(50)
	endEpoch := abi.ChainEpoch(300)
	sectorExpiry := abi.ChainEpoch(400)

	t.Run("deal expiry -> regular payments till deal expires and then locked funds are unlocked", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		dealId := actor.publishAndActivateDeal(rt, client, mAddrs, startEpoch, endEpoch, 0, sectorExpiry)
		d := actor.getDealProposal(rt, dealId)

		// move the current epoch to startEpoch + 5 so payment is made
		current := startEpoch + 5 // 55
		rt.SetEpoch(current)

		// assert payment
		pay, slashed := actor.cronTickAndAssertBalances(rt, client, provider, current, dealId)
		require.EqualValues(t, pay, big.Mul(big.NewInt(5), d.StoragePricePerEpoch))
		require.EqualValues(t, big.Zero(), slashed)

		// The next epoch for this deal's cron schedule is 155 (50 + 5 + 100).
		// Setting the current epoch to anything less than that wont make any payment
		current = 154
		rt.SetEpoch(current)
		actor.cronTickNoChangeBalances(rt, client, provider)

		// however setting the current epoch to 155 will make the payment
		current = 155
		rt.SetEpoch(current)
		pay, slashed = actor.cronTickAndAssertBalances(rt, client, provider, current, dealId)
		require.EqualValues(t, pay, big.Mul(big.NewInt(100), d.StoragePricePerEpoch))
		require.EqualValues(t, big.Zero(), slashed)

		// next epoch for cron schedule is 155 + 100 = 255
		current = 255
		rt.SetEpoch(current)
		pay, slashed = actor.cronTickAndAssertBalances(rt, client, provider, current, dealId)
		require.EqualValues(t, pay, big.Mul(big.NewInt(100), d.StoragePricePerEpoch))
		require.EqualValues(t, big.Zero(), slashed)

		// next epoch for cron schedule is deal end i.e. 300. An epoch less than that wont do anything
		current = 299
		rt.SetEpoch(current)
		actor.cronTickNoChangeBalances(rt, client, provider)

		// however setting epoch to 300 will expire the deal, make the payment and unlock all funds
		current = 300
		rt.SetEpoch(current)
		pay, slashed = actor.cronTickAndAssertBalances(rt, client, provider, current, dealId)
		require.EqualValues(t, pay, big.Mul(big.NewInt(45), d.StoragePricePerEpoch))
		require.EqualValues(t, big.Zero(), slashed)

		// deal should be deleted as it should have expired
		actor.assertDealDeleted(rt, dealId)
	})

	t.Run("deal expiry -> payment for a deal if deal is already expired before a cron tick", func(t *testing.T) {
		start := abi.ChainEpoch(5)
		end := abi.ChainEpoch(20)

		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		dealId := actor.publishAndActivateDeal(rt, client, mAddrs, start, end, 0, sectorExpiry)
		d := actor.getDealProposal(rt, dealId)

		current := abi.ChainEpoch(25)
		rt.SetEpoch(current)

		pay, slashed := actor.cronTickAndAssertBalances(rt, client, provider, current, dealId)
		require.EqualValues(t, pay, big.Mul(big.NewInt(15), d.StoragePricePerEpoch))
		require.EqualValues(t, big.Zero(), slashed)

		actor.assertDealDeleted(rt, dealId)

		// running cron tick again dosen't do anything
		actor.cronTickNoChangeBalances(rt, client, provider)
	})

	t.Run("expired deal should unlock the remaining client and provider locked balance after payment and deal should be deleted", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		dealId := actor.publishAndActivateDeal(rt, client, mAddrs, startEpoch, endEpoch, 0, sectorExpiry)
		deal := actor.getDealProposal(rt, dealId)

		cEscrow := actor.getEscrowBalance(rt, client)
		pEscrow := actor.getEscrowBalance(rt, provider)

		// move the current epoch so that deal is expired
		rt.SetEpoch(startEpoch + 1000)
		actor.cronTick(rt)

		// assert balances
		payment := deal.TotalStorageFee()

		require.EqualValues(t, big.Sub(cEscrow, payment), actor.getEscrowBalance(rt, client))
		require.EqualValues(t, big.Zero(), actor.getLockedBalance(rt, client))

		require.EqualValues(t, big.Add(pEscrow, payment), actor.getEscrowBalance(rt, provider))
		require.EqualValues(t, big.Zero(), actor.getLockedBalance(rt, provider))

		// deal should be deleted
		actor.assertDealDeleted(rt, dealId)
	})
}

func TestCronTickDealSlashing(t *testing.T) {
	owner := tutil.NewIDAddr(t, 101)
	provider := tutil.NewIDAddr(t, 102)
	worker := tutil.NewIDAddr(t, 103)
	client := tutil.NewIDAddr(t, 104)
	mAddrs := &minerAddrs{owner, worker, provider}
	sectorExpiry := abi.ChainEpoch(400)

	// hairy edge cases
	{

		tcs := map[string]struct {
			dealStart        abi.ChainEpoch
			dealEnd          abi.ChainEpoch
			activationEpoch  abi.ChainEpoch
			terminationEpoch abi.ChainEpoch
			cronTickEpoch    abi.ChainEpoch
			payment          abi.TokenAmount
			assertionMsg     string
		}{
			"deal is slashed after the startepoch and then the first crontick happens": {
				dealStart:        abi.ChainEpoch(10),
				dealEnd:          abi.ChainEpoch(20),
				activationEpoch:  abi.ChainEpoch(5),
				terminationEpoch: abi.ChainEpoch(15),
				cronTickEpoch:    abi.ChainEpoch(16),
				payment:          abi.NewTokenAmount(50), // (15 - 10) * 10 as deal storage fee is 10 per epoch
			},
			"deal is slashed at the startepoch and then the first crontick happens": {
				dealStart:        abi.ChainEpoch(10),
				dealEnd:          abi.ChainEpoch(20),
				activationEpoch:  abi.ChainEpoch(5),
				terminationEpoch: abi.ChainEpoch(10),
				cronTickEpoch:    abi.ChainEpoch(11),
				payment:          abi.NewTokenAmount(0), // (10 - 10) * 10
			},
			"deal is slashed before the startepoch and then the first crontick happens": {
				dealStart:        abi.ChainEpoch(10),
				dealEnd:          abi.ChainEpoch(20),
				activationEpoch:  abi.ChainEpoch(5),
				terminationEpoch: abi.ChainEpoch(6),
				cronTickEpoch:    abi.ChainEpoch(10),
				payment:          abi.NewTokenAmount(0), // (10 - 10) * 10
			},
			"deal is terminated at the activation epoch and then the first crontick happens": {
				dealStart:        abi.ChainEpoch(10),
				dealEnd:          abi.ChainEpoch(20),
				activationEpoch:  abi.ChainEpoch(5),
				terminationEpoch: abi.ChainEpoch(5),
				cronTickEpoch:    abi.ChainEpoch(10),
				payment:          abi.NewTokenAmount(0), // (10 - 10) * 10
			},
			"deal is slashed and then deal expiry happens on crontick, but slashing still occurs": {
				dealStart:        abi.ChainEpoch(10),
				dealEnd:          abi.ChainEpoch(20),
				activationEpoch:  abi.ChainEpoch(5),
				terminationEpoch: abi.ChainEpoch(15),
				cronTickEpoch:    abi.ChainEpoch(25), // deal has expired
				payment:          abi.NewTokenAmount(50),
			},
			"deal slash epoch must NOT be greater than current epoch": {
				dealStart:        abi.ChainEpoch(10),
				dealEnd:          abi.ChainEpoch(20),
				activationEpoch:  abi.ChainEpoch(5),
				terminationEpoch: abi.ChainEpoch(15),
				cronTickEpoch:    abi.ChainEpoch(10), // deal has expired
				payment:          abi.NewTokenAmount(50),
				assertionMsg:     "current epoch less than slash epoch",
			},
			"deal is slashed just BEFORE the end epoch": {
				dealStart:        abi.ChainEpoch(10),
				dealEnd:          abi.ChainEpoch(20),
				activationEpoch:  abi.ChainEpoch(5),
				terminationEpoch: abi.ChainEpoch(19),
				cronTickEpoch:    abi.ChainEpoch(19),
				payment:          abi.NewTokenAmount(90), // (19 - 10) * 10
			},
		}

		for n, tc := range tcs {
			t.Run(n, func(t *testing.T) {
				rt, actor := basicMarketSetup(t, owner, provider, worker, client)

				// publish and activate
				rt.SetEpoch(tc.activationEpoch)
				dealId := actor.publishAndActivateDeal(rt, client, mAddrs, tc.dealStart, tc.dealEnd, tc.activationEpoch, sectorExpiry)
				d := actor.getDealProposal(rt, dealId)

				// terminate
				rt.SetEpoch(tc.terminationEpoch)
				actor.terminateDeals(rt, provider, dealId)

				//  cron tick
				rt.SetEpoch(tc.cronTickEpoch)

				if len(tc.assertionMsg) == 0 {
					pay, slashed := actor.cronTickAndAssertBalances(rt, client, provider, tc.cronTickEpoch, dealId)
					require.EqualValues(t, tc.payment, pay)
					require.EqualValues(t, d.ProviderCollateral, slashed)
					actor.assertDealDeleted(rt, dealId)

					// running cron tick again dosen't do anything
					actor.cronTickNoChangeBalances(rt, client, provider)
				} else {
					rt.ExpectAssertionFailure(tc.assertionMsg, func() {
						rt.ExpectValidateCallerAddr(builtin.CronActorAddr)
						rt.SetCaller(builtin.CronActorAddr, builtin.CronActorCodeID)
						param := adt.EmptyValue{}
						rt.Call(actor.CronTick, &param)
						rt.Verify()
					})
				}
			})
		}
	}

	startEpoch := abi.ChainEpoch(50)
	endEpoch := abi.ChainEpoch(300)

	t.Run("deal is slashed AT the end epoch -> should NOT be slashed and should be considered expired", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		dealId := actor.publishAndActivateDeal(rt, client, mAddrs, startEpoch, endEpoch, 0, sectorExpiry)
		d := actor.getDealProposal(rt, dealId)

		// set current epoch to deal end epoch and attempt to slash it -> should not be slashed
		// as deal is considered to be expired.
		current := endEpoch
		rt.SetEpoch(current)
		actor.terminateDeals(rt, provider, dealId)

		// on the next cron tick, it will be processed as expired
		current = 300
		rt.SetEpoch(current)
		pay, slashed := actor.cronTickAndAssertBalances(rt, client, provider, current, dealId)
		duration := big.NewInt(250) // end - start
		require.EqualValues(t, big.Mul(duration, d.StoragePricePerEpoch), pay)
		require.EqualValues(t, big.Zero(), slashed)

		// deal should be deleted as it should have expired
		actor.assertDealDeleted(rt, dealId)
	})

	// end-end test for slashing
	t.Run("regular payments till deal is slashed and then slashing is processed", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		dealId := actor.publishAndActivateDeal(rt, client, mAddrs, startEpoch, endEpoch, 0, sectorExpiry)
		d := actor.getDealProposal(rt, dealId)

		// move the current epoch to startEpoch + 5 so payment is made
		current := abi.ChainEpoch(55)
		rt.SetEpoch(current)

		// assert payment
		pay, slashed := actor.cronTickAndAssertBalances(rt, client, provider, current, dealId)
		require.EqualValues(t, pay, big.Mul(big.NewInt(5), d.StoragePricePerEpoch))
		require.EqualValues(t, big.Zero(), slashed)

		// Setting the current epoch to less than 155 will NOT make any changes as the deal
		// is still not scheduled
		current = 154
		rt.SetEpoch(current)
		actor.cronTickNoChangeBalances(rt, client, provider)

		//  Setting the current epoch to 155 will make another payment (5 + 100 epochs)
		current = 155
		rt.SetEpoch(current)
		pay, slashed = actor.cronTickAndAssertBalances(rt, client, provider, current, dealId)
		require.EqualValues(t, pay, big.Mul(big.NewInt(100), d.StoragePricePerEpoch))
		require.EqualValues(t, big.Zero(), slashed)

		// now terminate the deal
		current = 200
		rt.SetEpoch(current)
		actor.terminateDeals(rt, provider, dealId)

		// Setting the epoch to anything less than 255 will NOT make any changes even though the deal is slashed (155 + 100)
		current = 254
		rt.SetEpoch(current)
		actor.cronTickNoChangeBalances(rt, client, provider)

		// next epoch for cron schedule is 155 + 100 = 255 -> payment will be made and deal will be slashed
		current = 255
		rt.SetEpoch(current)
		pay, slashed = actor.cronTickAndAssertBalances(rt, client, provider, current, dealId)
		// payment will only be made till the 200th epoch as the deal was slashed at that epoch.
		// so duration = 200 - 155(epoch of last payment) = 45.
		require.EqualValues(t, pay, big.Mul(big.NewInt(45), d.StoragePricePerEpoch))
		require.EqualValues(t, d.ProviderCollateral, slashed)

		// deal should be deleted as it should have expired
		actor.assertDealDeleted(rt, dealId)
	})

	// expired deals should NOT be slashed
	t.Run("regular payments till deal expires and then we attempt to slash it but it will NOT be slashed", func(t *testing.T) {
		rt, actor := basicMarketSetup(t, owner, provider, worker, client)
		dealId := actor.publishAndActivateDeal(rt, client, mAddrs, startEpoch, endEpoch, 0, sectorExpiry)
		d := actor.getDealProposal(rt, dealId)

		// move the current epoch to startEpoch + 5 so payment is made and assert payment
		current := startEpoch + 5 // 55
		rt.SetEpoch(current)
		pay, slashed := actor.cronTickAndAssertBalances(rt, client, provider, current, dealId)
		require.EqualValues(t, pay, big.Mul(big.NewInt(5), d.StoragePricePerEpoch))
		require.EqualValues(t, big.Zero(), slashed)

		//  Setting the current epoch to 155 will make another payment
		current = 155
		rt.SetEpoch(current)
		pay, slashed = actor.cronTickAndAssertBalances(rt, client, provider, current, dealId)
		require.EqualValues(t, pay, big.Mul(big.NewInt(100), d.StoragePricePerEpoch))
		require.EqualValues(t, big.Zero(), slashed)

		// set current epoch to deal end epoch and attempt to slash it -> should not be slashed
		// as deal is considered to be expired.
		current = endEpoch
		rt.SetEpoch(current)
		actor.terminateDeals(rt, provider, dealId)

		// next epoch for cron schedule is 155 + 100 = 255 ->
		// setting epoch to higher than that will cause deal to be expired, payment will be made
		// and deal will NOT be slashed
		current = 300
		rt.SetEpoch(current)
		pay, slashed = actor.cronTickAndAssertBalances(rt, client, provider, current, dealId)
		duration := big.NewInt(145) // 300 which is the end epoch MINUS the previous payment epoch i.e. 155
		require.EqualValues(t, big.Mul(duration, d.StoragePricePerEpoch), pay)
		require.EqualValues(t, big.Zero(), slashed)

		// deal should be deleted as it should have expired
		actor.assertDealDeleted(rt, dealId)
	})
}

func TestMarketActorDeals(t *testing.T) {
	owner := tutil.NewIDAddr(t, 101)
	provider := tutil.NewIDAddr(t, 102)
	worker := tutil.NewIDAddr(t, 103)
	client := tutil.NewIDAddr(t, 104)
	minerAddrs := &minerAddrs{owner, worker, provider}

	var st market.State

	// Test adding provider funds from both worker and owner address
	rt, actor := basicMarketSetup(t, owner, provider, worker, client)
	actor.addProviderFunds(rt, abi.NewTokenAmount(10000), minerAddrs)
	rt.GetState(&st)
	assert.Equal(t, abi.NewTokenAmount(10000), st.GetEscrowBalance(rt, provider))

	actor.addParticipantFunds(rt, client, abi.NewTokenAmount(10000))

	dealProposal := generateDealProposal(client, provider, abi.ChainEpoch(1), abi.ChainEpoch(5))
	params := &market.PublishStorageDealsParams{Deals: []market.ClientDealProposal{market.ClientDealProposal{Proposal: dealProposal}}}

	// First attempt at publishing the deal should work
	{
		actor.publishDeals(rt, minerAddrs, dealProposal)
	}

	// Second attempt at publishing the same deal should fail
	{
		rt.ExpectValidateCallerType(builtin.AccountActorCodeID, builtin.MultisigActorCodeID)
		rt.ExpectSend(provider, builtin.MethodsMiner.ControlAddresses, nil, abi.NewTokenAmount(0), &miner.GetControlAddressesReturn{Worker: worker, Owner: owner}, 0)

		rt.ExpectVerifySignature(crypto.Signature{}, client, mustCbor(&params.Deals[0].Proposal), nil)
		rt.SetCaller(worker, builtin.AccountActorCodeID)
		rt.ExpectAbort(exitcode.ErrIllegalArgument, func() {
			rt.Call(actor.PublishStorageDeals, params)
		})

		rt.Verify()
	}

	dealProposal.Label = "foo"

	// Same deal with a different label should work
	{
		actor.publishDeals(rt, minerAddrs, dealProposal)
	}
}

type marketActorTestHarness struct {
	market.Actor
	t testing.TB
}

func (h *marketActorTestHarness) constructAndVerify(rt *mock.Runtime) {
	rt.ExpectValidateCallerAddr(builtin.SystemActorAddr)
	ret := rt.Call(h.Constructor, nil)
	assert.Nil(h.t, ret)
	rt.Verify()
}

type minerAddrs struct {
	owner    address.Address
	worker   address.Address
	provider address.Address
}

// addProviderFunds is a helper method to setup provider market funds
func (h *marketActorTestHarness) addProviderFunds(rt *mock.Runtime, amount abi.TokenAmount, minerAddrs *minerAddrs) {
	rt.SetReceived(amount)
	rt.SetAddressActorType(minerAddrs.provider, builtin.StorageMinerActorCodeID)
	rt.SetCaller(minerAddrs.owner, builtin.AccountActorCodeID)
	h.expectProviderControlAddressesAndValidateCaller(rt, minerAddrs.provider, minerAddrs.owner, minerAddrs.worker)

	rt.Call(h.AddBalance, &minerAddrs.provider)

	rt.Verify()

	rt.SetBalance(big.Add(rt.Balance(), amount))
}

// addParticipantFunds is a helper method to setup non-provider storage market participant funds
func (h *marketActorTestHarness) addParticipantFunds(rt *mock.Runtime, addr address.Address, amount abi.TokenAmount) {
	rt.SetReceived(amount)
	rt.SetCaller(addr, builtin.AccountActorCodeID)
	rt.ExpectValidateCallerType(builtin.CallerTypesSignable...)

	rt.Call(h.AddBalance, &addr)

	rt.Verify()

	rt.SetBalance(big.Add(rt.Balance(), amount))
}

func (h *marketActorTestHarness) expectProviderControlAddressesAndValidateCaller(rt *mock.Runtime, provider address.Address, owner address.Address, worker address.Address) {
	rt.ExpectValidateCallerAddr(owner, worker)

	expectRet := &miner.GetControlAddressesReturn{Owner: owner, Worker: worker}

	rt.ExpectSend(
		provider,
		builtin.MethodsMiner.ControlAddresses,
		nil,
		big.Zero(),
		expectRet,
		exitcode.Ok,
	)
}

func (h *marketActorTestHarness) withdrawProviderBalance(rt *mock.Runtime, withDrawAmt, expectedSend abi.TokenAmount, miner *minerAddrs) {
	rt.SetCaller(miner.worker, builtin.AccountActorCodeID)
	h.expectProviderControlAddressesAndValidateCaller(rt, miner.provider, miner.owner, miner.worker)

	params := market.WithdrawBalanceParams{
		ProviderOrClientAddress: miner.provider,
		Amount:                  withDrawAmt,
	}

	rt.ExpectSend(miner.owner, builtin.MethodSend, nil, expectedSend, nil, exitcode.Ok)
	rt.Call(h.WithdrawBalance, &params)
	rt.Verify()
}

func (h *marketActorTestHarness) withdrawClientBalance(rt *mock.Runtime, client address.Address, withDrawAmt, expectedSend abi.TokenAmount) {
	rt.SetCaller(client, builtin.AccountActorCodeID)
	rt.ExpectValidateCallerType(builtin.CallerTypesSignable...)
	rt.ExpectSend(client, builtin.MethodSend, nil, expectedSend, nil, exitcode.Ok)

	params := market.WithdrawBalanceParams{
		ProviderOrClientAddress: client,
		Amount:                  withDrawAmt,
	}

	rt.Call(h.WithdrawBalance, &params)
	rt.Verify()
}

func (h *marketActorTestHarness) cronTickNoChangeBalances(rt *mock.Runtime, client, provider address.Address) {
	// fetch current client and provider escrow balances
	cLocked := h.getLockedBalance(rt, client)
	cEscrow := h.getEscrowBalance(rt, client)
	pLocked := h.getLockedBalance(rt, provider)
	pEscrow := h.getEscrowBalance(rt, provider)

	h.cronTick(rt)

	require.EqualValues(h.t, cEscrow, h.getEscrowBalance(rt, client))
	require.EqualValues(h.t, cLocked, h.getLockedBalance(rt, client))
	require.EqualValues(h.t, pEscrow, h.getEscrowBalance(rt, provider))
	require.EqualValues(h.t, pLocked, h.getLockedBalance(rt, provider))
}

func (h *marketActorTestHarness) cronTickAndAssertBalances(rt *mock.Runtime, client, provider address.Address,
	currentEpoch abi.ChainEpoch, dealId abi.DealID) (payment abi.TokenAmount, amountSlashed abi.TokenAmount) {
	// fetch current client and provider escrow balances
	cLocked := h.getLockedBalance(rt, client)
	cEscrow := h.getEscrowBalance(rt, client)
	pLocked := h.getLockedBalance(rt, provider)
	pEscrow := h.getEscrowBalance(rt, provider)
	amountSlashed = big.Zero()

	s := h.getDealState(rt, dealId)
	d := h.getDealProposal(rt, dealId)

	// end epoch for payment calc
	paymentEnd := d.EndEpoch
	if s.SlashEpoch != -1 {
		rt.ExpectSend(builtin.BurntFundsActorAddr, builtin.MethodSend, nil, d.ProviderCollateral, nil, exitcode.Ok)
		amountSlashed = d.ProviderCollateral

		if s.SlashEpoch < d.StartEpoch {
			paymentEnd = d.StartEpoch
		} else {
			paymentEnd = s.SlashEpoch
		}
	} else if currentEpoch < paymentEnd {
		paymentEnd = currentEpoch
	}

	// start epoch for payment calc
	paymentStart := d.StartEpoch
	if s.LastUpdatedEpoch != -1 {
		paymentStart = s.LastUpdatedEpoch
	}
	duration := paymentEnd - paymentStart
	payment = big.Mul(big.NewInt(int64(duration)), d.StoragePricePerEpoch)

	// expected updated amounts
	updatedClientEscrow := big.Sub(cEscrow, payment)
	updatedProviderEscrow := big.Add(pEscrow, payment)
	updatedProviderEscrow = big.Sub(updatedProviderEscrow, amountSlashed)
	updatedClientLocked := big.Sub(cLocked, payment)
	updatedProviderLocked := pLocked
	// if the deal has expired or been slashed, locked amount will be zero for provider and client.
	isDealExpired := paymentEnd == d.EndEpoch
	if isDealExpired || s.SlashEpoch != -1 {
		updatedClientLocked = big.Zero()
		updatedProviderLocked = big.Zero()
	}

	h.cronTick(rt)

	require.EqualValues(h.t, updatedClientEscrow, h.getEscrowBalance(rt, client))
	require.EqualValues(h.t, updatedClientLocked, h.getLockedBalance(rt, client))
	require.Equal(h.t, updatedProviderLocked, h.getLockedBalance(rt, provider))
	require.Equal(h.t, updatedProviderEscrow.Int64(), h.getEscrowBalance(rt, provider).Int64())

	return
}

func (h *marketActorTestHarness) cronTick(rt *mock.Runtime) {
	rt.ExpectValidateCallerAddr(builtin.CronActorAddr)
	rt.SetCaller(builtin.CronActorAddr, builtin.CronActorCodeID)
	param := adt.EmptyValue{}
	rt.Call(h.CronTick, &param)
	rt.Verify()
}

func (h *marketActorTestHarness) publishDeals(rt *mock.Runtime, minerAddrs *minerAddrs, deals ...market.DealProposal) []abi.DealID {
	rt.SetCaller(minerAddrs.worker, builtin.AccountActorCodeID)
	rt.ExpectValidateCallerType(builtin.CallerTypesSignable...)
	rt.ExpectSend(
		minerAddrs.provider,
		builtin.MethodsMiner.ControlAddresses,
		nil,
		big.Zero(),
		&miner.GetControlAddressesReturn{Owner: minerAddrs.owner, Worker: minerAddrs.worker},
		exitcode.Ok,
	)

	var params market.PublishStorageDealsParams

	for _, deal := range deals {
		//  create a client proposal with a valid signature
		buf := bytes.Buffer{}
		require.NoError(h.t, deal.MarshalCBOR(&buf), "failed to marshal deal proposal")
		sig := crypto.Signature{Type: crypto.SigTypeBLS, Data: []byte("does not matter")}
		clientProposal := market.ClientDealProposal{deal, sig}
		params.Deals = append(params.Deals, clientProposal)

		// expect a call to verify the above signature
		rt.ExpectVerifySignature(sig, deal.Client, buf.Bytes(), nil)
		if deal.VerifiedDeal {
			param := &verifreg.UseBytesParams{
				Address:  deal.Client,
				DealSize: big.NewIntUnsigned(uint64(deal.PieceSize)),
			}

			rt.ExpectSend(builtin.VerifiedRegistryActorAddr, builtin.MethodsVerifiedRegistry.UseBytes, param, abi.NewTokenAmount(0), nil, exitcode.Ok)
		}
	}

	ret := rt.Call(h.PublishStorageDeals, &params)
	rt.Verify()

	resp, ok := ret.(*market.PublishStorageDealsReturn)
	require.True(h.t, ok, "unexpected type returned from call to PublishStorageDeals")
	require.Len(h.t, resp.IDs, len(deals))

	// assert state after publishing the deals
	dealIds := resp.IDs
	for i, deaId := range dealIds {
		expected := deals[i]
		p := h.getDealProposal(rt, deaId)

		require.Equal(h.t, expected.StartEpoch, p.StartEpoch)
		require.Equal(h.t, expected.EndEpoch, p.EndEpoch)
		require.Equal(h.t, expected.PieceCID, p.PieceCID)
		require.Equal(h.t, expected.PieceSize, p.PieceSize)
		require.Equal(h.t, expected.Client, p.Client)
		require.Equal(h.t, expected.Provider, p.Provider)
		require.Equal(h.t, expected.Label, p.Label)
		require.Equal(h.t, expected.VerifiedDeal, p.VerifiedDeal)
		require.Equal(h.t, expected.StoragePricePerEpoch, p.StoragePricePerEpoch)
		require.Equal(h.t, expected.ClientCollateral, p.ClientCollateral)
		require.Equal(h.t, expected.ProviderCollateral, p.ProviderCollateral)
	}

	return resp.IDs
}

func (h *marketActorTestHarness) assertDealsNotActivated(rt *mock.Runtime, epoch abi.ChainEpoch, dealIDs ...abi.DealID) {
	var st market.State
	rt.GetState(&st)

	states, err := market.AsDealStateArray(adt.AsStore(rt), st.States)
	require.NoError(h.t, err)

	for _, d := range dealIDs {
		_, found, err := states.Get(d)
		require.NoError(h.t, err)
		require.False(h.t, found)
	}
}

func (h *marketActorTestHarness) activateDeals(rt *mock.Runtime, sectorExpiry abi.ChainEpoch, provider address.Address, currentEpoch abi.ChainEpoch, dealIDs ...abi.DealID) {
	rt.SetCaller(provider, builtin.StorageMinerActorCodeID)
	rt.ExpectValidateCallerType(builtin.StorageMinerActorCodeID)

	params := &market.ActivateDealsParams{DealIDs: dealIDs, SectorExpiry: sectorExpiry}

	ret := rt.Call(h.ActivateDeals, params)
	rt.Verify()

	require.Nil(h.t, ret)

	for _, d := range dealIDs {
		s := h.getDealState(rt, d)
		require.EqualValues(h.t, currentEpoch, s.SectorStartEpoch)
	}
}

func (h *marketActorTestHarness) getDealProposal(rt *mock.Runtime, dealID abi.DealID) *market.DealProposal {
	var st market.State
	rt.GetState(&st)

	deals, err := market.AsDealProposalArray(adt.AsStore(rt), st.Proposals)
	require.NoError(h.t, err)

	d, found, err := deals.Get(dealID)
	require.NoError(h.t, err)
	require.True(h.t, found)
	require.NotNil(h.t, d)

	return d
}

func (h *marketActorTestHarness) getEscrowBalance(rt *mock.Runtime, addr address.Address) abi.TokenAmount {
	var st market.State
	rt.GetState(&st)

	return st.GetEscrowBalance(rt, addr)
}

func (h *marketActorTestHarness) getLockedBalance(rt *mock.Runtime, addr address.Address) abi.TokenAmount {
	var st market.State
	rt.GetState(&st)

	return st.GetLockedBalance(rt, addr)
}

func (h *marketActorTestHarness) getDealState(rt *mock.Runtime, dealID abi.DealID) *market.DealState {
	var st market.State
	rt.GetState(&st)

	states, err := market.AsDealStateArray(adt.AsStore(rt), st.States)
	require.NoError(h.t, err)

	s, found, err := states.Get(dealID)
	require.NoError(h.t, err)
	require.True(h.t, found)
	require.NotNil(h.t, s)

	return s
}

func (h *marketActorTestHarness) assertLockedFundStates(rt *mock.Runtime, storageFee, providerCollateral, clientCollateral abi.TokenAmount) {
	var st market.State
	rt.GetState(&st)

	require.Equal(h.t, clientCollateral, st.TotalClientLockedCollateral)
	require.Equal(h.t, providerCollateral, st.TotalProviderLockedCollateral)
	require.Equal(h.t, storageFee, st.TotalClientStorageFee)
}

func (h *marketActorTestHarness) assertDealDeleted(rt *mock.Runtime, dealId abi.DealID) {
	var st market.State
	rt.GetState(&st)

	proposals, err := market.AsDealProposalArray(adt.AsStore(rt), st.Proposals)
	require.NoError(h.t, err)
	_, found, err := proposals.Get(dealId)
	require.NoError(h.t, err)
	require.False(h.t, found)

	states, err := market.AsDealStateArray(adt.AsStore(rt), st.States)
	require.NoError(h.t, err)
	_, found, err = states.Get(dealId)
	require.NoError(h.t, err)
	require.False(h.t, found)
}

func (h *marketActorTestHarness) assertDealsTerminated(rt *mock.Runtime, epoch abi.ChainEpoch, dealIds ...abi.DealID) {
	for _, d := range dealIds {
		s := h.getDealState(rt, d)
		require.EqualValues(h.t, epoch, s.SlashEpoch)
	}
}

func (h *marketActorTestHarness) assertDeaslNotTerminated(rt *mock.Runtime, dealIds ...abi.DealID) {
	for _, d := range dealIds {
		s := h.getDealState(rt, d)
		require.EqualValues(h.t, abi.ChainEpoch(-1), s.SlashEpoch)
	}
}

func (h *marketActorTestHarness) terminateDeals(rt *mock.Runtime, minerAddr address.Address, dealIds ...abi.DealID) {
	rt.SetCaller(minerAddr, builtin.StorageMinerActorCodeID)
	rt.ExpectValidateCallerType(builtin.StorageMinerActorCodeID)

	params := &market.OnMinerSectorsTerminateParams{DealIDs: dealIds}

	ret := rt.Call(h.OnMinerSectorsTerminate, params)
	rt.Verify()
	require.Nil(h.t, ret)
}

func (h *marketActorTestHarness) publishAndActivateDeal(rt *mock.Runtime, client address.Address, minerAddrs *minerAddrs,
	startEpoch, endEpoch, currentEpoch, sectorExpiry abi.ChainEpoch) abi.DealID {
	deal := h.generateDealAndAddFunds(rt, client, minerAddrs, startEpoch, endEpoch)
	dealIds := h.publishDeals(rt, minerAddrs, deal)
	h.activateDeals(rt, sectorExpiry, minerAddrs.provider, currentEpoch, dealIds[0])
	return dealIds[0]
}

func (h *marketActorTestHarness) updateLastUpdated(rt *mock.Runtime, dealId abi.DealID, newLastUpdated abi.ChainEpoch) {
	var st market.State

	rt.Transaction(&st, func() interface{} {
		states, err := market.AsDealStateArray(adt.AsStore(rt), st.States)
		require.NoError(h.t, err)
		s, found, err := states.Get(dealId)
		require.True(h.t, found)
		require.NoError(h.t, err)
		require.NotNil(h.t, s)

		require.NoError(h.t, states.Set(dealId, &market.DealState{s.SectorStartEpoch, newLastUpdated, s.SlashEpoch}))
		st.States, err = states.Root()
		require.NoError(h.t, err)
		return nil
	})
}

func (h *marketActorTestHarness) deleteDealProposal(rt *mock.Runtime, dealId abi.DealID) {
	var st market.State

	rt.Transaction(&st, func() interface{} {
		deals, err := market.AsDealProposalArray(adt.AsStore(rt), st.Proposals)
		require.NoError(h.t, err)
		require.NoError(h.t, deals.Delete(uint64(dealId)))
		st.Proposals, err = deals.Root()
		require.NoError(h.t, err)
		return nil
	})
}

func (h *marketActorTestHarness) generateAndPublishDeal(rt *mock.Runtime, client address.Address, minerAddrs *minerAddrs,
	startEpoch, endEpoch abi.ChainEpoch) abi.DealID {
	deal := h.generateDealAndAddFunds(rt, client, minerAddrs, startEpoch, endEpoch)
	dealIds := h.publishDeals(rt, minerAddrs, deal)
	return dealIds[0]
}

func (h *marketActorTestHarness) generateDealAndAddFunds(rt *mock.Runtime, client address.Address, minerAddrs *minerAddrs,
	startEpoch, endEpoch abi.ChainEpoch) market.DealProposal {
	deal4 := generateDealProposal(client, minerAddrs.provider, startEpoch, endEpoch)
	h.addProviderFunds(rt, deal4.ProviderCollateral, minerAddrs)
	h.addParticipantFunds(rt, client, deal4.ClientBalanceRequirement())

	return deal4
}

func generateDealProposal(client, provider address.Address, startEpoch, endEpoch abi.ChainEpoch) market.DealProposal {
	pieceCid := tutil.MakeCID("1")
	pieceSize := abi.PaddedPieceSize(2048)
	storagePerEpoch := big.NewInt(10)
	clientCollateral := big.NewInt(10)
	providerCollateral := big.NewInt(10)

	return market.DealProposal{pieceCid, pieceSize, false, client, provider, "label", startEpoch,
		endEpoch, storagePerEpoch, providerCollateral, clientCollateral}
}

func basicMarketSetup(t *testing.T, owner, provider, worker, client address.Address) (*mock.Runtime, *marketActorTestHarness) {
	builder := mock.NewBuilder(context.Background(), builtin.StorageMarketActorAddr).
		WithCaller(builtin.SystemActorAddr, builtin.InitActorCodeID).
		WithActorType(owner, builtin.AccountActorCodeID).
		WithActorType(worker, builtin.AccountActorCodeID).
		WithActorType(provider, builtin.StorageMinerActorCodeID).
		WithActorType(client, builtin.AccountActorCodeID)

	rt := builder.Build(t)

	actor := marketActorTestHarness{t: t}
	actor.constructAndVerify(rt)

	return rt, &actor
}

func mkPublishStorageParams(proposals ...market.DealProposal) *market.PublishStorageDealsParams {
	m := &market.PublishStorageDealsParams{}
	for _, p := range proposals {
		m.Deals = append(m.Deals, market.ClientDealProposal{Proposal: p})
	}
	return m
}

func mkActivateDealParams(sectorExpiry abi.ChainEpoch, dealIds ...abi.DealID) *market.ActivateDealsParams {
	return &market.ActivateDealsParams{SectorExpiry: sectorExpiry, DealIDs: dealIds}
}

func mkTerminateDealParams(dealIds ...abi.DealID) *market.OnMinerSectorsTerminateParams {
	return &market.OnMinerSectorsTerminateParams{dealIds}
}
