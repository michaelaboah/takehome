package ops

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rdsdata"
	rdstypes "github.com/aws/aws-sdk-go-v2/service/rdsdata/types"
	"purs.example.com/m/pkg"
	"purs.example.com/m/pkg/gen"
	types "purs.example.com/m/types"
)

// Purs Transaction Operations
// Object for generating parameters
type PTOperations struct {
	statement     *rdsdata.ExecuteStatementInput
	LedgeEntries  [][]byte
	LedgerEntryId []byte
	paymentId     []byte
}

// Randomly generate `paymentId` and `LedgerEntryId`
func (pt *PTOperations) Init(sqlTransctionID, resourceArn, secretArn, database string) {
	pt.paymentId = gen.RandBytes(32)
	pt.LedgerEntryId = gen.RandBytes(32)
	statement := &rdsdata.ExecuteStatementInput{
		Database:      aws.String(database),
		ResourceArn:   aws.String(resourceArn),
		SecretArn:     aws.String(secretArn),
		TransactionId: aws.String(sqlTransctionID),
	}

	pt.statement = statement

}

// Builds the parameters for a Card transaction
// Does not clear prior parameters
func (pt *PTOperations) BuildCardParams(userPuchaseInfo types.UserPurchaseInformation, insertPaymentSQL string) *rdsdata.ExecuteStatementInput {
	if len(pt.paymentId) == 0 || len(pt.LedgerEntryId) == 0 {
		panic("Failed to initialize 'PTOperations'. Call '.Init()' before mutating the object")
	}

	cardparams := pkg.CardParams{
		DevId:       []byte(userPuchaseInfo.Dev),
		Method:      userPuchaseInfo.PaymentMethod,
		LedgerId:    pt.LedgerEntryId,
		PayerId:     []byte(userPuchaseInfo.Payor),
		PayeeId:     []byte(userPuchaseInfo.Payee),
		Amount:      float64(userPuchaseInfo.Amount),
		Interaction: float64(userPuchaseInfo.InteractionType),
	}

	if userPuchaseInfo.PaymentMethod == pkg.FedNow && userPuchaseInfo.Amount > 0 {
		cardparams.DatePaid = ""
	} else {
		cardparams.DatePaid = time.Now().Format(time.RFC3339Nano) // conform to ISO
	}

	if userPuchaseInfo.PaymentMethod != pkg.FedNow || userPuchaseInfo.Amount == 0 {
		cardparams.Status = pkg.Completed
	} else {
		cardparams.Status = pkg.Pending
	}

	pt.statement.Sql = aws.String(insertPaymentSQL)
	pt.statement.Parameters = append(pt.statement.Parameters, pkg.ToSQLParams(cardparams)...)

	return pt.statement
}

// Builds the parameters for a FedNow transaction
// Does not clear prior parameters
func (pt *PTOperations) BuildFedNowParams(userPuchaseInfo types.UserPurchaseInformation, id []byte, insertPaymentSQL string) *rdsdata.ExecuteStatementInput {
	if len(pt.paymentId) == 0 || len(pt.LedgerEntryId) == 0 {
		panic("Failed to initialize 'PTOperations'. Call '.Init()' before mutating the object")
	}
	fedparams := pkg.FedNowParams{
		FedNowPaymentId: id,
		PayerAccountId:  []byte(userPuchaseInfo.PayorBankAccountID),
		PayeeAccountId:  []byte(userPuchaseInfo.PayeeBankAccountID),
	}

	pt.statement.Sql = aws.String(insertPaymentSQL)
	pt.statement.Parameters = append(pt.statement.Parameters, pkg.ToSQLParams(fedparams)...)

	return pt.statement
}

// Builds the parameters for a Ledger
// **Does** override prior parameters
func (pt *PTOperations) BuildLedgerParams(userPuchaseInfo types.UserPurchaseInformation, id []byte, promoAmount int, insertPaymentSQL string) *rdsdata.ExecuteStatementInput {
	if len(pt.paymentId) == 0 || len(pt.LedgerEntryId) == 0 {
		panic("Failed to initialize 'PTOperations'. Call '.Init()' before mutating the object")
	}

	ledgerParms := pkg.LedgerParams{
		PayerId:     []byte(userPuchaseInfo.Dev),
		PayeeId:     []byte(userPuchaseInfo.Payee),
		Amount:      float64(promoAmount),
		Interaction: float64(userPuchaseInfo.InteractionType),
		LedgerId:    id,
		DevId:       []byte(userPuchaseInfo.Dev),
	}

	pt.LedgeEntries = append(pt.LedgeEntries, id)

	pt.statement.Sql = aws.String(insertPaymentSQL)
	pt.statement.Parameters = pkg.ToSQLParams(ledgerParms)

	return pt.statement
}

func (pt *PTOperations) BuildBatch(insertPaymentSQL string, ledgerParams []rdstypes.SqlParameter) rdsdata.BatchExecuteStatementInput {
	batch := rdsdata.BatchExecuteStatementInput{
		ResourceArn:   pt.statement.ResourceArn,
		SecretArn:     pt.statement.SecretArn,
		Sql:           aws.String(insertPaymentSQL),
		ParameterSets: [][]rdstypes.SqlParameter{pt.statement.Parameters, ledgerParams},
		TransactionId: pt.statement.TransactionId,
	}
	return batch
}

// Clears the SQL parameters in the `statement` object
func (pt *PTOperations) ClearParams() {
	pt.statement.Parameters = make([]rdstypes.SqlParameter, 0)
}

func (pt *PTOperations) GetParams() *[]rdstypes.SqlParameter {
	return &pt.statement.Parameters
}
