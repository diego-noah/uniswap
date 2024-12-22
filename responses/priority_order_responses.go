package responses

import "time"

// GetPriorityOrderResponse represents the API response for priority orders
type GetPriorityOrderResponse struct {
	Type                   string                         `json:"type"`
	OrderStatus            string                         `json:"orderStatus"`
	Signature              string                         `json:"signature"`
	EncodedOrder           string                         `json:"encodedOrder"`
	ChainId                int64                          `json:"chainId"`
	Nonce                  string                         `json:"nonce"`
	TxHash                 string                         `json:"txHash,omitempty"`
	OrderHash              string                         `json:"orderHash"`
	Swapper                string                         `json:"swapper"`
	Reactor                string                         `json:"reactor"`
	Deadline               int64                          `json:"deadline"`
	AuctionStartBlock      int64                          `json:"auctionStartBlock"`
	BaselinePriorityFeeWei string                         `json:"baselinePriorityFeeWei"`
	Input                  *PriorityOrderInputResponse    `json:"input"`
	Outputs                []*PriorityOrderOutputResponse `json:"outputs"`
	CosignerData           *CosignerDataResponse          `json:"cosignerData"`
	Cosignature            string                         `json:"cosignature"`
	QuoteId                string                         `json:"quoteId,omitempty"`
	RequestId              string                         `json:"requestId,omitempty"`
	CreatedAt              *time.Time                     `json:"createdAt,omitempty"`
}

type PriorityOrderInputResponse struct {
	Token                string `json:"token"`
	Amount               string `json:"amount"`
	MpsPerPriorityFeeWei string `json:"mpsPerPriorityFeeWei"`
}

type PriorityOrderOutputResponse struct {
	Token                string `json:"token"`
	Amount               string `json:"amount"`
	MpsPerPriorityFeeWei string `json:"mpsPerPriorityFeeWei"`
	Recipient            string `json:"recipient"`
}

type CosignerDataResponse struct {
	AuctionTargetBlock int64 `json:"auctionTargetBlock"`
}
