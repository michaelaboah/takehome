package pkg

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	rdstypes "github.com/aws/aws-sdk-go-v2/service/rdsdata/types"
)

// Possible Payment Methods
// - FedNow = 0
// - Card = 1
type PaymentMethod int

const (
	FedNow PaymentMethod = 0
	Card   PaymentMethod = 1
)

// Possible Interaction Types
// - Mobile = 1
type InteractionType int

const (
	Mobile InteractionType = 0
)

// Possible Payment Status'
// Completed = 0
// Pending = 1
type PaymentStatus int

const (
	Completed PaymentStatus = 4
	Pending   PaymentStatus = 5
)

type PromotionInfo struct {
	PromoAmount int
}

// Compatible with IPaymentParam
type CardParams struct {
	PayerId     []byte        `json:"payerId"`
	PayeeId     []byte        `json:"payeeId"`
	Amount      float64       `json:"amount"`
	Interaction float64       `json:"interaction"`
	PaymentId   []byte        `json:"paymentId"`
	DatePaid    string        `json:"datePaid"`
	LedgerId    []byte        `json:"ledgerId"`
	DevId       []byte        `json:"devId"`
	Method      PaymentMethod `json:"method"`
	Status      PaymentStatus `json:"status"`
}

func ToSQLParams(p IPaymentParam) []rdstypes.SqlParameter {
	var in map[string]interface{}
	var params []rdstypes.SqlParameter
	var sqlValue rdstypes.Field
	inres, _ := json.Marshal(p)
	json.Unmarshal(inres, &in)

	for field, val := range in {
		switch v := val.(type) {
		case int:
			sqlValue = &rdstypes.FieldMemberLongValue{
				Value: int64(v),
			}
		case float64:
			sqlValue = &rdstypes.FieldMemberDoubleValue{
				Value: v,
			}
		case string:
			sqlValue = &rdstypes.FieldMemberStringValue{
				Value: v,
			}
		case []byte:
			sqlValue = &rdstypes.FieldMemberBlobValue{
				Value: v,
			}
		}

		s := rdstypes.SqlParameter{
			Name:  aws.String(field),
			Value: sqlValue,
		}

		params = append(params, s)
	}

	return params
}

// Compatible with IPaymentParam
type FedNowParams struct {
	FedNowPaymentId []byte `json:"fedNowPaymentId"`
	PayerAccountId  []byte `json:"payerAccountId"`
	PayeeAccountId  []byte `json:"payeeAccountId"`
}

// Compatible with IPaymentParam
type LedgerParams struct {
	PayerId     []byte  `json:"payerId"`
	PayeeId     []byte  `json:"payeeId"`
	Amount      float64 `json:"amount"`
	Interaction float64 `json:"interaction"`
	LedgerId    []byte  `json:"ledgerId"`
	DevId       []byte  `json:"devId"`
}

// Interface for working with paramaters
type IPaymentParam interface {
	GetId()
}

func (pp CardParams) GetId()   {}
func (lp LedgerParams) GetId() {}
func (fp FedNowParams) GetId() {}
