package models

type LoanConfig struct {
	Id              uint64       `json:"id" gorm:"column:id;primaryKey"`
	PrincipalAmount string       `json:"principalAmount" gorm:"column:principal_amount"`
	MaxSpan         int          `json:"maxSpan" gorm:"column:max_span"`
	RateOfInterest  string       `json:"rateOfInterest" gorm:"column:rate_of_interest"`
	TypeOfLoan      InterestType `json:"typeOfLoan" gorm:"column:type_of_loan"`
	IsActive        bool         `json:"isActive" gorm:"column:is_active"`
}

type InterestType string

const (
	FixedInterestType    InterestType = "FIXED"
	VariableInterestType InterestType = "VARIABLE"
	SimpleInterestType   InterestType = "SIMPLE"
	CompoundInterestType InterestType = "COMPOUND"
)
