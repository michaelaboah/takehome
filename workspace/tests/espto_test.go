package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	o "purs.example.com/m/ops"
	"purs.example.com/m/pkg"
	"purs.example.com/m/pkg/gen"
	"purs.example.com/m/types"
)

// Ensure that parameters aren't deleted
func TestParamsAppend(t *testing.T) {
	var ops o.PTOperations

	u := exUserPurchaseInfo(pkg.Mobile, pkg.Card)

	ops.Init("", "", "", "")
	ops.BuildCardParams(u, "")

	ops.BuildFedNowParams(u, gen.RandBytes(32), "")

	assert.Equal(t, 13, len(*ops.GetParams()))

}

func TestCard(t *testing.T) {
	var ops o.PTOperations
	u := exUserPurchaseInfo(pkg.Mobile, pkg.Card)

	ops.Init("", "", "", "")
	ops.BuildCardParams(u, "")

	assert.Equal(t, 10, len(*ops.GetParams()))
}

func TestFedNow(t *testing.T) {
	var ops o.PTOperations
	id := gen.RandBytes(32)
	u := exUserPurchaseInfo(pkg.Mobile, pkg.FedNow)

	ops.Init("", "", "", "")
	ops.BuildFedNowParams(u, id, "")

	assert.Equal(t, 3, len(*ops.GetParams()))
}

func TestLedger(t *testing.T) {
	var ops o.PTOperations
	u := exUserPurchaseInfo(pkg.Mobile, pkg.Card)

	id := gen.RandBytes(32)

	ops.Init("", "", "", "")
	ops.BuildLedgerParams(u, id, 100, "")

	assert.Equal(t, 6, len(*ops.GetParams()))

}

func exUserPurchaseInfo(interactionType pkg.InteractionType, paymentMethod pkg.PaymentMethod) types.UserPurchaseInformation {
	u := types.UserPurchaseInformation{
		Payor:              gen.RandString(32),
		Payee:              gen.RandString(32),
		PayorBankAccountID: gen.RandString(32),
		PayeeBankAccountID: gen.RandString(32),
		Dev:                gen.RandString(32),
		Amount:             0,
		InteractionType:    interactionType,
		PaymentMethod:      paymentMethod,
	}

	return u
}
