package main

import (
	"log"

	"github.com/AshimKoirala/load-balancer-admin/pkg/handlers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// 		// Test sending an email
	// 	err = utils.NewEmailResponse("ashimkoirala01@gmail.com", "Test Email", "This is a test email.")
	// 	if err != nil {
	//     log.Printf("Failed to send email: %v", err)
	//    }else{
	// 	log.Println("Email sent successfully!")
	//    }

	// messaging.InitializePublisher()
	// defer messaging.CleanupPublisher()

	// go func() {
	// 	messaging.SetupConsumer()
	// }()

	handlers.Handler()
}
