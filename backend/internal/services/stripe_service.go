package services

import (
	"errors"
	"os"
	"strings"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

func GetStripeSecretKey() string {
	return strings.TrimSpace(os.Getenv("STRIPE_SECRET_KEY"))
}

func GetAppDomain() string {
	domain := os.Getenv("APP_DOMAIN")
	if domain == "" {
		domain = "http://localhost:5173"
	}
	return domain
}

func CreateStripeCheckoutSession(lineItems []*stripe.CheckoutSessionLineItemParams, successPath, cancelPath string) (string, error) {
	apiKey := GetStripeSecretKey()
	stripe.Key = apiKey

	domain := GetAppDomain()
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String(domain + successPath),
		CancelURL:          stripe.String(domain + cancelPath),
	}

	sessionStripe, err := session.New(params)
	if err != nil {
		return "", err
	}

	if sessionStripe == nil || sessionStripe.URL == "" {
		return "", errors.New("URL de session Stripe vide")
	}

	return sessionStripe.URL, nil
}