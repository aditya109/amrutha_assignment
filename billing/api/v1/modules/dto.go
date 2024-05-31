package modules

type CustomerDto struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}

type TransitionLoanRequestDto struct {
	CustomerId      string `json:"-"`
	ConfigurationId *int   `json:"configurationId"`
}

type MakePaymentRequestDto struct {
	CustomerId             string `json:"-"`
	Amount                 string `json:"amount" binding:"required"`
	DateOfTransaction      string `json:"dateOfTransaction" binding:"required"`
	TransactionReferenceId string `json:"transactionReferenceId" binding:"required"`
}
