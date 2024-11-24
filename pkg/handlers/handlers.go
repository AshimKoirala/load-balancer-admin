package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/AshimKoirala/load-balancer-admin/middleware"
	"github.com/AshimKoirala/load-balancer-admin/pkg/db"
)

func Handler() {

	if err := db.InitDB(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	// Routes setup with CORS
	mux := http.NewServeMux()

	// Route setup
	mux.HandleFunc("/admin/register", AuthRegister)
	mux.HandleFunc("/admin/login", AuthLogin)
	mux.Handle("/admin/protected", middleware.AuthMiddleware(http.HandlerFunc(ProtectedRoute)))
	mux.HandleFunc("/admin/users", GetUsers)
	mux.HandleFunc("/admin/update", UpdateUser)
	mux.HandleFunc("/admin/forgotpassword", ForgotPassword)
	mux.HandleFunc("/admin/resetpassword", ResetPassword)
	mux.HandleFunc("/admin/add_replica", AddReplica)
    mux.HandleFunc("/admin/remove_replica", RemoveReplica)
    mux.HandleFunc("/admin/change_status", ChangeStatus)

	// Wrap the entire mux with CORS
	handlerWithCORS := middleware.CORS(mux)

	// Start the server
	port := os.Getenv("SERVER_PORT")
    if port == "" {
    port = "8080" 
       }
	log.Println("Server is running on : %s",port)
	if err := http.ListenAndServe(":"+port, handlerWithCORS); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
