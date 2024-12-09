package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/AshimKoirala/load-balancer-admin/messaging"
	"github.com/AshimKoirala/load-balancer-admin/pkg/db"
	"github.com/AshimKoirala/load-balancer-admin/utils"
)

// to fetch latest PrequalParametersResponse
func GetPrequalParameters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.NewErrorResponse(w, http.StatusMethodNotAllowed, []string{"Method not allowed"})
		return
	}

	response, err := db.GetPrequalParametersResponse(r.Context())
	if err != nil {
		log.Printf("Error fetching latest prequal parameters response: %v", err)
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Failed to fetch the latest entry"})
		return
	}

	utils.NewSuccessResponse(w, response)
}

// to create a new PrequalParametersResponse
func AddPrequalParameters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.NewErrorResponse(w, http.StatusMethodNotAllowed, []string{"Method not allowed"})
		return
	}

	var payload struct {
		Data       db.PrequalParametersResponse `json:"data"`
		ActivateId *int64                       `json:"activate_id,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.NewErrorResponse(w, http.StatusBadRequest, []string{"Invalid request payload"})
		return
	}

	// Set default status if not provided
	if payload.Data.Status == "" {
		payload.Data.Status = "inactive"
	}

	// Call the database function with the payload
	err := db.AddPrequalParametersResponse(r.Context(), payload.Data, payload.ActivateId)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Failed to create or activate entry"})
		return
	}

	message := "Prequal Parameter added successfully"
	if payload.ActivateId != nil {
		message += fmt.Sprintf(" and entry with ID %d was activated", *payload.ActivateId)
	}

	rabbitmessage := &messaging.Message{
		Name: messaging.NEW_PARAMETERS,
		Body: map[string]interface{}{
			"data":        payload.Data,
			"activate_id": payload.ActivateId,
		},
	}

	// Publish the message to RabbitMQ
	err = messaging.PublishMessage(messaging.PUBLISHING_QUEUE, rabbitmessage)
	if err != nil {
		utils.NewErrorResponse(w, http.StatusInternalServerError, []string{"Failed to publish message to RabbitMQ"})
		return
	}

	utils.NewSuccessResponse(w, message)
}
