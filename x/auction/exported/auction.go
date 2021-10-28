package exported

import (
	"time"

	"github.com/gogo/protobuf/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Auction is an interface for handling common actions on auctions.
type Auction interface {
	proto.Message

	GetId() int64
	WithID(int64) Auction

	GetInitiator() string
	GetLot() sdk.Coin
	GetBidder() string
	GetBid() sdk.Coin
	GetEndTime() time.Time
	GetMaxEndTime() time.Time

	GetType() string
	GetPhase() string
}

// GenesisAuction extends the auction interface to add functionality
// needed for initializing auctions from genesis.
type GenesisAuction interface {
	Auction
	GetModuleAccountCoins() sdk.Coins
	Validate() error
}
