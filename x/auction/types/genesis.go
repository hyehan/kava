package types

import (
	"bytes"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/kava-labs/kava/x/auction/exported"
)

// DefaultNextAuctionID is the starting point for auction IDs.
const DefaultNextAuctionID int64 = 1

var (
	_ exported.Auction        = &GenesisAuction{}
	_ exported.GenesisAuction = &GenesisAuction{}
)

func (a *GenesisAuction) GetBidder() string { return a.GetBidder() }

func (a *GenesisAuction) GetInitiator() string { return a.GetInitiator() }

func (a *GenesisAuction) GetEndTime() time.Time { return a.GetEndTime() }

func (a *GenesisAuction) GetMaxEndTime() time.Time { return a.GetMaxEndTime() }

func (a *GenesisAuction) GetBid() sdk.Coin { return a.GetBid() }

func (a *GenesisAuction) GetLot() sdk.Coin { return a.GetLot() }

func (a *GenesisAuction) GetId() int64 { return a.GetId() }

func (a *GenesisAuction) WithID(id int64) exported.Auction { return a.WithID(id) }

// GetPhase returns the direction of a surplus auction, which never changes.
func (a GenesisAuction) GetPhase() string { return a.GetPhase() }

// GetType returns the auction type. Used to identify auctions in event attributes.
func (a GenesisAuction) GetType() string { return a.GetType() }

func (a *GenesisAuction) GetModuleAccountCoins() sdk.Coins { return a.GetModuleAccountCoins() }

func (a *GenesisAuction) Validate() error { return a.Validate() }

func (a GenesisAuction) String() string {
	return fmt.Sprintf(`Auction %d:
  Initiator:              %s
  Lot:               			%s
  Bidder:            		  %s
  Bid:        						%s
  End Time:   						%s
  Max End Time:      			%s`,
		a.GetId(), a.GetInitiator(), a.GetLot(),
		a.GetBidder(), a.GetBid(), a.GetEndTime().String(),
		a.GetMaxEndTime().String(),
	)
}

// NewGenesisState returns a new genesis state object for auctions module.
func NewGenesisState(nextID int64, ap Params, ga []*GenesisAuction) GenesisState {
	return GenesisState{
		NextAuctionId: nextID,
		Params:        ap,
		Auctions:      ga,
	}
}

// DefaultGenesisState returns the default genesis state for auction module.
func DefaultGenesisState() GenesisState {
	return NewGenesisState(
		DefaultNextAuctionID,
		DefaultParams(),
		[]*GenesisAuction{},
	)
}

// Equal checks whether two GenesisState structs are equivalent.
func (gs GenesisState) Equal(gs2 GenesisState) bool {
	b1 := ModuleCdc.Amino.MustMarshalBinaryBare(&gs)
	b2 := ModuleCdc.Amino.MustMarshalBinaryBare(&gs2)
	return bytes.Equal(b1, b2)
}

// IsEmpty returns true if a GenesisState is empty.
func (gs GenesisState) IsEmpty() bool {
	return gs.Equal(GenesisState{})
}

// Validate validates genesis inputs. It returns error if validation of any input fails.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	ids := map[int64]bool{}
	for _, a := range gs.Auctions {

		if err := a.Validate(); err != nil {
			return fmt.Errorf("found invalid auction: %w", err)
		}

		if ids[a.GetId()] {
			return fmt.Errorf("found duplicate auction ID (%d)", a.GetId())
		}
		ids[a.GetId()] = true

		if a.GetId() >= gs.NextAuctionId {
			return fmt.Errorf("found auction ID ≥ the nextAuctionID (%d ≥ %d)", a.GetId(), gs.NextAuctionId)
		}
	}
	return nil
}
