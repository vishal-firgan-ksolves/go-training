package main

import (
	"fmt"
)

type Notifier interface{
	Send(message string) error
}

type EmailNotifier struct{
	Email string
}

type SMSNotifier struct{
	ContactNo int64
}

func (emailNotifier	EmailNotifier) Send(message string) error{
	fmt.Printf("Sending Email notification on email %s with message %s",emailNotifier.Email,message)
	return nil
}

func (smsNotifier SMSNotifier) Send(message string) error{
	fmt.Printf("Sending message notification on contact number %d with message %s",smsNotifier.ContactNo,message)
	return nil
}

func sendNotification(notifer Notifier,message string) error{
	
	if err := notifer.Send(message); 
	
	err != nil {
		fmt.Printf("Notification failed to send: %v\n", err)
	}
	
	return nil
}

func main(){

	emailService := EmailNotifier{
		Email: "vishal@gmail.com",
	}
	
	smsService := SMSNotifier{
	    ContactNo: 9876543210,
	}

	err := sendNotification(emailService, "You won a car! \n")

	if err != nil {
		fmt.Println("CRITICAL: Failed to send email:", err)
	}

    err = sendNotification(smsService, "Your OTP is 1234 \n")
	
	if err != nil {
		fmt.Println("CRITICAL: Failed to send SMS:", err)
	}
}