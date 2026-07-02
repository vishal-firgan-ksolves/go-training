package main

import "fmt"

type Paymenter interface {
	makePayment(amount float32) string
	refund(amount float32) string
}

type RazorPayGateway struct{}
type PaypalPayGateway struct{}

func (pgw RazorPayGateway) makePayment(amount float32) string {
	return fmt.Sprintf("Making payment with RazorPay: %.2f", amount)
}

func (pgw RazorPayGateway) refund(amount float32) string {
	return fmt.Sprintf("Taking refund with RazorPay: %.2f", amount)
}

func (pgw PaypalPayGateway) makePayment(amount float32) string {
	return fmt.Sprintf("Making payment with Paypal: %.2f", amount)
}

func (pgw PaypalPayGateway) refund(amount float32) string {
	return fmt.Sprintf("Taking refund with Paypal: %.2f", amount)
}

func checkout(p Paymenter, amount float32) string {
	return p.makePayment(amount)
}

func main() {
	var RazorPay Paymenter = RazorPayGateway{}
	PaypalPay := PaypalPayGateway{}

	fmt.Println(RazorPay.makePayment(33))
	fmt.Println(PaypalPay.makePayment(55))

	fmt.Println(checkout(RazorPay, 44))
}