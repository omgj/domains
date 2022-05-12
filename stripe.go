package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/account"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/stripe/stripe-go/paymentmethod"
)

type Method struct {
	ID    string `json:"id"`
	Brand string `json:"brand"`
	Last4 string `json:"last4"`
}

type PaymentIntent struct {
	Secret string `json:"secret"`
}

// TODO - card, bank account, afterpay...

func buy(w http.ResponseWriter, r *http.Request) {
	stripe.Key = StripeKEY
	params := &stripe.PaymentIntentParams{
		Customer: stripe.String("cus_L1S94rQ5fszdft"),
		// SetupFutureUsage:   stripe.String("off_session"),
		Amount:             stripe.Int64(1000),
		Currency:           stripe.String(string(stripe.CurrencyAUD)),
		PaymentMethodTypes: []*string{stripe.String("afterpay_clearpay"), stripe.String("card")},
		ReturnURL:          stripe.String(RETURN_URL),
		Confirm:            stripe.Bool(true),
		PaymentMethod:      stripe.String("pm_1KLPEzHfeA9afONrwTmWtbJY"),
	}
	params.AddExtra("payment_method_options[afterpay_clearpay][reference]", "order_123")
	pi, _ := paymentintent.New(params)
	i := PaymentIntent{
		Secret: pi.ClientSecret,
	}
	fmt.Println("Sending: ", i)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(i)
}

func CusMethods(id string) []Method {
	params := &stripe.PaymentMethodListParams{
		Customer: stripe.String(id),
		Type:     stripe.String("card"),
	}
	i := paymentmethod.List(params)
	var ms []Method
	for i.Next() {
		a := string(i.PaymentMethod().Card.Brand)
		ms = append(ms, Method{"i", a, i.PaymentMethod().Card.Last4})
	}
	return ms
}

func StripeCustomerAccount(name, email string) (string, string) {
	stripe.Key = StripeKEY

	cparams := &stripe.CustomerParams{
		Name:  stripe.String(name),
		Email: stripe.String(email),
	}
	c, e := customer.New(cparams)
	if e != nil {
		log.Fatal(e)
	}

	params := &stripe.AccountParams{
		Country:      stripe.String("AU"),
		Type:         stripe.String("custom"),
		Email:        stripe.String(email),
		BusinessType: stripe.String("individual"),
		Individual: &stripe.PersonParams{
			FirstName: stripe.String(""),
			LastName:  stripe.String(""),
			DOB: &stripe.DOBParams{
				Day:   stripe.Int64(1),
				Month: stripe.Int64(1),
				Year:  stripe.Int64(1),
			},
			Address: &stripe.AccountAddressParams{
				Line1:      stripe.String(""),
				PostalCode: stripe.String(""),
				City:       stripe.String(""),
				State:      stripe.String(""),
			},
			Email: stripe.String(email),
			Phone: stripe.String(""),
		},
		RequestedCapabilities: []*string{stripe.String("card_payments"), stripe.String("transfers")},
	}
	account, e := account.New(params)
	if e != nil {
		log.Fatal(e)
	}

	return c.ID, account.ID

}
