package orders

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/your-org/your-project/entities"
	"github.com/your-org/your-project/responses"
)

// PriorityOrder represents a priority order in the system
type PriorityOrder struct {
	Inner     sdk.PriorityOrder
	Signature string
	ChainId   int64
	Status    entities.OrderStatus
	TxHash    string
	QuoteId   string
	RequestId string
	CreatedAt *time.Time
}

// NewPriorityOrder creates a new PriorityOrder instance
func NewPriorityOrder(inner sdk.PriorityOrder, signature string, chainId int64, status entities.OrderStatus,
	txHash, quoteId, requestId string, createdAt *time.Time) *PriorityOrder {
	return &PriorityOrder{
		Inner:     inner,
		Signature: signature,
		ChainId:   chainId,
		Status:    status,
		TxHash:    txHash,
		QuoteId:   quoteId,
		RequestId: requestId,
		CreatedAt: createdAt,
	}
}

// OrderType returns the type of the order
func (p *PriorityOrder) OrderType() sdk.OrderType {
	return sdk.OrderTypePriority
}

// ToEntity converts the PriorityOrder to a PriorityOrderEntity
func (p *PriorityOrder) ToEntity(orderStatus entities.OrderStatus) entities.PriorityOrderEntity {
	input := p.Inner.Info.Input
	outputs := p.Inner.Info.Outputs

	// Convert outputs to entity format
	entityOutputs := make([]entities.PriorityOrderOutput, len(outputs))
	for i, output := range outputs {
		entityOutputs[i] = entities.PriorityOrderOutput{
			Token:                output.Token,
			Amount:               output.Amount.String(),
			MpsPerPriorityFeeWei: output.MpsPerPriorityFeeWei.String(),
			Recipient:            strings.ToLower(output.Recipient),
		}
	}

	return entities.PriorityOrderEntity{
		Type:                   sdk.OrderTypePriority,
		EncodedOrder:           p.Inner.Serialize(),
		Signature:              p.Signature,
		Nonce:                  p.Inner.Info.Nonce.String(),
		OrderHash:              strings.ToLower(p.Inner.Hash()),
		ChainId:                p.Inner.ChainId,
		OrderStatus:            orderStatus,
		Offerer:                strings.ToLower(p.Inner.Info.Swapper),
		AuctionStartBlock:      p.Inner.Info.AuctionStartBlock.Int64(),
		BaselinePriorityFeeWei: p.Inner.Info.BaselinePriorityFeeWei.String(),
		Input: entities.PriorityOrderInput{
			Token:                input.Token,
			Amount:               input.Amount.String(),
			MpsPerPriorityFeeWei: input.MpsPerPriorityFeeWei.String(),
		},
		Outputs:  entityOutputs,
		Reactor:  strings.ToLower(p.Inner.Info.Reactor),
		Deadline: p.Inner.Info.Deadline,
		CosignerData: entities.CosignerData{
			AuctionTargetBlock: p.Inner.Info.CosignerData.AuctionTargetBlock.Int64(),
		},
		TxHash:      p.TxHash,
		Cosignature: p.Inner.Info.Cosignature,
		QuoteId:     p.QuoteId,
		RequestId:   p.RequestId,
		CreatedAt:   p.CreatedAt,
	}
}

// FromEntity creates a PriorityOrder from a UniswapXOrderEntity
func FromEntity(entity entities.UniswapXOrderEntity) (*PriorityOrder, error) {
	order, err := sdk.ParsePriorityOrder(entity.EncodedOrder, entity.ChainId)
	if err != nil {
		return nil, err
	}

	return NewPriorityOrder(
		order,
		entity.Signature,
		entity.ChainId,
		entity.OrderStatus,
		entity.TxHash,
		entity.QuoteId,
		entity.RequestId,
		entity.CreatedAt,
	), nil
}

// ReparameterizeAndCosign updates and cosigns the order
func (p *PriorityOrder) ReparameterizeAndCosign(client *ethclient.Client, cosigner *KmsSigner) error {
	currentBlock, err := client.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	targetBlock := new(big.Int).Add(
		new(big.Int).SetUint64(currentBlock),
		new(big.Int).SetInt64(PRIORITY_ORDER_TARGET_BLOCK_BUFFER),
	)

	p.Inner.Info.CosignerData = sdk.CosignerData{
		AuctionTargetBlock: targetBlock,
	}

	cosignatureHash, err := p.Inner.CosignatureHash(p.Inner.Info.CosignerData)
	if err != nil {
		return err
	}

	signature, err := cosigner.SignDigest(cosignatureHash)
	if err != nil {
		return err
	}

	p.Inner.Info.Cosignature = signature
	return nil
}

// ToGetResponse converts the order to a response format
func (p *PriorityOrder) ToGetResponse() *responses.GetPriorityOrderResponse {
	return &responses.GetPriorityOrderResponse{
		Type:                   sdk.OrderTypePriority,
		OrderStatus:            p.Status,
		Signature:              p.Signature,
		EncodedOrder:           p.Inner.Serialize(),
		ChainId:                p.ChainId,
		Nonce:                  p.Inner.Info.Nonce.String(),
		TxHash:                 p.TxHash,
		OrderHash:              p.Inner.Hash(),
		Swapper:                p.Inner.Info.Swapper,
		Reactor:                p.Inner.Info.Reactor,
		Deadline:               p.Inner.Info.Deadline,
		AuctionStartBlock:      p.Inner.Info.AuctionStartBlock.Int64(),
		BaselinePriorityFeeWei: p.Inner.Info.BaselinePriorityFeeWei.String(),
		Input: &responses.PriorityOrderInputResponse{
			Token:                p.Inner.Info.Input.Token,
			Amount:               p.Inner.Info.Input.Amount.String(),
			MpsPerPriorityFeeWei: p.Inner.Info.Input.MpsPerPriorityFeeWei.String(),
		},
		Outputs: p.formatOutputsResponse(),
		CosignerData: &responses.CosignerDataResponse{
			AuctionTargetBlock: p.Inner.Info.CosignerData.AuctionTargetBlock.Int64(),
		},
		Cosignature: p.Inner.Info.Cosignature,
		QuoteId:     p.QuoteId,
		RequestId:   p.RequestId,
		CreatedAt:   p.CreatedAt,
	}
}

func (p *PriorityOrder) formatOutputsResponse() []*responses.PriorityOrderOutputResponse {
	outputs := make([]*responses.PriorityOrderOutputResponse, len(p.Inner.Info.Outputs))
	for i, o := range p.Inner.Info.Outputs {
		outputs[i] = &responses.PriorityOrderOutputResponse{
			Token:                o.Token,
			Amount:               p.Inner.Info.Input.Amount.String(),
			MpsPerPriorityFeeWei: p.Inner.Info.Input.MpsPerPriorityFeeWei.String(),
			Recipient:            o.Recipient,
		}
	}
	return outputs
}
