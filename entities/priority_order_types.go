package entities

import "time"

// PriorityOrderEntity extends UniswapXOrderEntity with priority-specific fields
type PriorityOrderEntity struct {
	UniswapXOrderEntity
	AuctionStartBlock      int64
	BaselinePriorityFeeWei string
	Input                  PriorityOrderInput
	Outputs                []PriorityOrderOutput
	CosignerData           CosignerData
	Cosignature            string
	CreatedAt              *time.Time
}

// PriorityOrderInput represents the input for a priority order
type PriorityOrderInput struct {
	Token                string
	Amount               string
	MpsPerPriorityFeeWei string
}

// PriorityOrderOutput represents an output in a priority order
type PriorityOrderOutput struct {
	Token                string
	Amount               string
	MpsPerPriorityFeeWei string
	Recipient            string
}

// CosignerData represents the cosigner-specific data
type CosignerData struct {
	AuctionTargetBlock int64
}
