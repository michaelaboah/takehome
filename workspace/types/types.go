package types

import "purs.example.com/m/pkg"

type UserPurchaseInformation struct {
	Payor              string              `json:"payor"`
	Payee              string              `json:"payee"`
	PayorBankAccountID string              `json:"payorBankAccount"`
	PayeeBankAccountID string              `json:"payeeBankAccount"`
	Dev                string              `json:"dev"`
	Amount             int                 `json:"amount"`
	InteractionType    pkg.InteractionType `json:"interactionType"`
	PaymentMethod      pkg.PaymentMethod   `json:"paymentMethod"`
}
