package entities

// OrderStatus represents the status of an order
type OrderStatus string

// UniswapXOrderEntity represents an order entity in the system
type UniswapXOrderEntity struct {
	Type           string
	EncodedOrder   string
	Signature      string
	Nonce          string
	OrderHash      string
	ChainId        int64
	OrderStatus    OrderStatus
	Offerer        string
	Input          OrderInput
	Outputs        []OrderOutput
	Reactor        string
	DecayStartTime int64
	DecayEndTime   int64
	Deadline       int64
	Filler         string
}

// OrderInput represents the input part of an order
type OrderInput struct {
	Token       string
	StartAmount string
	EndAmount   string
}

// OrderOutput represents an output in an order
type OrderOutput struct {
	Token       string
	StartAmount string
	EndAmount   string
	Recipient   string
}
