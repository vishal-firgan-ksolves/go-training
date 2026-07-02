package main

import (
	"fmt"
	"testing"
)

// FAKE Gateway for testing
type MockGateway struct{}

func (m MockGateway) makePayment(amount float32) string {
	return fmt.Sprintf("MOCK PAYMENT SUCCESS: %.2f", amount)
}

func (m MockGateway) refund(amount float32) string {
	return fmt.Sprintf("MOCK REFUND SUCCESS: %.2f", amount)
}

func TestCheckout(t *testing.T) {
	fakePaymentProcessor := MockGateway{}
	
	testAmount := float32(100.50)
	expected := "MOCK PAYMENT SUCCESS: 100.50"

	// 3. We pass the FAKE gateway into your real checkout function!
	got := checkout(fakePaymentProcessor, testAmount)

	if got != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, got)
	}
}

func TestRealGateways(t *testing.T) {
	razor := RazorPayGateway{}
	paypal := PaypalPayGateway{}

	if got := razor.makePayment(50); got != "Making payment with RazorPay: 50.00" {
		t.Errorf("RazorPay failed. Got: %s", got)
	}

	if got := paypal.makePayment(25); got != "Making payment with Paypal: 25.00" {
		t.Errorf("Paypal failed. Got: %s", got)
	}
}