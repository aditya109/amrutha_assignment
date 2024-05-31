package modules

type CustomerDto struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}

type ActivateLoanRequestDto struct {
	CustomerId      string `json:"-"`
	ConfigurationId *int   `json:"configurationId"`
}
