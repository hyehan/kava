package auction_test

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/kava-labs/kava/app"
	"github.com/kava-labs/kava/x/auction"
	"github.com/kava-labs/kava/x/auction/types"
)

var _, testAddrs = app.GeneratePrivKeyAddressPairs(2)
var testTime = time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)
var testAuction = types.NewCollateralAuction(
	"seller",
	c("lotdenom", 10),
	testTime,
	c("biddenom", 1000),
	types.WeightedAddresses{Addresses: testAddrs, Weights: []sdk.Int{sdk.OneInt(), sdk.OneInt()}},
	c("debt", 1000),
).WithID(3).(types.GenesisAuction)

func TestInitGenesis(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		// setup keepers
		tApp := app.NewTestApp()
		ctx := tApp.NewContext(true, tmproto.Header{Height: 1})

		// setup module account
		modBaseAcc := authtypes.NewBaseAccount(authtypes.NewModuleAddress(types.ModuleName), nil, 0, 0)
		modAcc := authtypes.NewModuleAccount(modBaseAcc, types.ModuleName, []string{authtypes.Minter, authtypes.Burner}...)
		tApp.GetAccountKeeper().SetModuleAccount(ctx, modAcc)
		tApp.GetBankKeeper().MintCoins(ctx, types.ModuleName, testAuction.GetModuleAccountCoins())

		// set up auction genesis state with module account
		auctionGS, err := types.NewGenesisState(
			10,
			types.DefaultParams(),
			[]types.GenesisAuction{testAuction},
		)
		require.NoError(t, err)

		// run init
		keeper := tApp.GetAuctionKeeper()
		require.NotPanics(t, func() {
			auction.InitGenesis(ctx, keeper, tApp.GetBankKeeper(), tApp.GetAccountKeeper(), auctionGS)
		})

		// check state is as expected
		actualID, err := keeper.GetNextAuctionID(ctx)
		require.NoError(t, err)
		require.Equal(t, auctionGS.NextAuctionId, actualID)

		require.Equal(t, auctionGS.Params, keeper.GetParams(ctx))

		genesisAuctions, err := types.UnpackGenesisAuctions(auctionGS.Auctions)
		if err != nil {
			panic(err)
		}

		// TODO is there a nicer way of comparing state?
		sort.Slice(genesisAuctions, func(i, j int) bool {
			return genesisAuctions[i].GetID() > genesisAuctions[j].GetID()
		})
		i := 0
		keeper.IterateAuctions(ctx, func(a types.Auction) bool {
			// types.UnpackGenesisAuctions()
			require.Equal(t, auctionGS.Auctions[i], a)
			i++
			return false
		})
	})
	// t.Run("invalid (invalid nextAuctionID)", func(t *testing.T) {
	// 	// setup keepers
	// 	tApp := app.NewTestApp()
	// 	ctx := tApp.NewContext(true, abci.Header{})

	// 	// create invalid genesis
	// 	gs := auction.NewGenesisState(
	// 		0, // next id < testAuction ID
	// 		auction.DefaultParams(),
	// 		auction.GenesisAuctions{testAuction},
	// 	)

	// 	// check init fails
	// 	require.Panics(t, func() {
	// 		auction.InitGenesis(ctx, tApp.GetAuctionKeeper(), tApp.GetSupplyKeeper(), gs)
	// 	})
	// })
	// t.Run("invalid (missing mod account coins)", func(t *testing.T) {
	// 	// setup keepers
	// 	tApp := app.NewTestApp()
	// 	ctx := tApp.NewContext(true, abci.Header{})

	// 	// create invalid genesis
	// 	gs := auction.NewGenesisState(
	// 		10,
	// 		auction.DefaultParams(),
	// 		auction.GenesisAuctions{testAuction},
	// 	)
	// 	// invalid as there is no module account setup

	// 	// check init fails
	// 	require.Panics(t, func() {
	// 		auction.InitGenesis(ctx, tApp.GetAuctionKeeper(), tApp.GetSupplyKeeper(), gs)
	// 	})
	// })
}

// func TestExportGenesis(t *testing.T) {
// 	t.Run("default", func(t *testing.T) {
// 		// setup state
// 		tApp := app.NewTestApp()
// 		ctx := tApp.NewContext(true, abci.Header{})
// 		tApp.InitializeFromGenesisStates()

// 		// export
// 		gs := auction.ExportGenesis(ctx, tApp.GetAuctionKeeper())

// 		// check state matches
// 		require.Equal(t, auction.DefaultGenesisState(), gs)
// 	})
// 	t.Run("one auction", func(t *testing.T) {
// 		// setup state
// 		tApp := app.NewTestApp()
// 		ctx := tApp.NewContext(true, abci.Header{})
// 		tApp.InitializeFromGenesisStates()
// 		tApp.GetAuctionKeeper().SetAuction(ctx, testAuction)

// 		// export
// 		gs := auction.ExportGenesis(ctx, tApp.GetAuctionKeeper())

// 		// check state matches
// 		expectedGenesisState := auction.DefaultGenesisState()
// 		expectedGenesisState.Auctions = append(expectedGenesisState.Auctions, testAuction)
// 		require.Equal(t, expectedGenesisState, gs)
// 	})
// }
