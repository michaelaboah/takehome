package pkg

import (
	"fmt"
	"testing"
	"time"

	"purs.example.com/m/pkg/gen"
)

func TestCardToSQLParams(t *testing.T) {
	c := CardParams{
		PayerId:     gen.RandBytes(32),
		PayeeId:     gen.RandBytes(32),
		Amount:      0,
		Interaction: 0,
		PaymentId:   gen.RandBytes(32),
		DatePaid:    time.Now().Format(time.RFC3339Nano),
		LedgerId:    gen.RandBytes(32),
		DevId:       gen.RandBytes(32),
		Method:      Card,
		Status:      Completed,
	}

	fmt.Println(c.Status)
	sqlparms := ToSQLParams(c)

	for _, v := range sqlparms {
		fmt.Print(*v.Name, ", ")
		fmt.Println(v.Value)
	}
}
