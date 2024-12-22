package orders

import (
	"strings"

	// You'll need to implement/import SDK equivalents
	"github.com/your-org/your-project/entities"
)

// LimitOrder represents a limit order in the system
type LimitOrder struct {
	Inner     sdk.DutchOrder
	Signature string
	ChainId   int64
	QuoteId   string
	RequestId string
}

// NewLimitOrder creates a new LimitOrder instance
func NewLimitOrder(inner sdk.DutchOrder, signature string, chainId int64, quoteId, requestId string) *LimitOrder {
	return &LimitOrder{
		Inner:     inner,
		Signature: signature,
		ChainId:   chainId,
		QuoteId:   quoteId,
		RequestId: requestId,
	}
}

// OrderType returns the type of the order
func (l *LimitOrder) OrderType() sdk.OrderType {
	return sdk.OrderTypeLimit
}

// ToEntity converts the LimitOrder to a UniswapXOrderEntity
func (l *LimitOrder) ToEntity(orderStatus entities.OrderStatus) entities.UniswapXOrderEntity {
	input := l.Inner.Info.Input
	outputs := l.Inner.Info.Outputs

	// Convert outputs to entity format
	entityOutputs := make([]entities.OrderOutput, len(outputs))
	for i, output := range outputs {
		entityOutputs[i] = entities.OrderOutput{
			Token:       output.Token,
			StartAmount: output.StartAmount.String(),
			EndAmount:   output.EndAmount.String(),
			Recipient:   strings.ToLower(output.Recipient),
		}
	}

	return entities.UniswapXOrderEntity{
		Type:         sdk.OrderTypeDutch,
		EncodedOrder: l.Inner.Serialize(),
		Signature:    l.Signature,
		Nonce:        l.Inner.Info.Nonce.String(),
		OrderHash:    strings.ToLower(l.Inner.Hash()),
		ChainId:      l.Inner.ChainId,
		OrderStatus:  orderStatus,
		Offerer:      strings.ToLower(l.Inner.Info.Swapper),
		Input: entities.OrderInput{
			Token:       input.Token,
			StartAmount: input.StartAmount.String(),
			EndAmount:   input.EndAmount.String(),
		},
		Outputs:        entityOutputs,
		Reactor:        strings.ToLower(l.Inner.Info.Reactor),
		DecayStartTime: l.Inner.Info.DecayStartTime,
		DecayEndTime:   l.Inner.Info.Deadline,
		Deadline:       l.Inner.Info.Deadline,
		Filler:         strings.ToLower(l.Inner.Info.ExclusiveFiller),
	}
}
