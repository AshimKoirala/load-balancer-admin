package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/uptrace/bun"
)

type PrequalParametersResponse struct {
	bun.BaseModel `bun:"table:prequal_parameters_response"`

	Id                int       `bun:"id,pk,autoincrement" json:"id"`
	MaxLifeTime       int       `bun:"max_life_time" json:"max_life_time"`
	PoolSize          int       `bun:"pool_size" json:"pool_size"`
	ProbeFactor       float64   `bun:"probe_factor" json:"probe_factor"`
	ProbeRemoveFactor int       `bun:"probe_remove_factor" json:"probe_remove_factor"`
	Mu                int       `bun:"mu" json:"mu"`
	CreatedAt         time.Time `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt         time.Time `bun:"updated_at,default:current_timestamp" json:"updated_at"`
	Status            string    `bun:"status,default:inactive" json:"status"`
}

// Fetch latest created row
func GetPrequalParametersResponse(ctx context.Context) (PrequalParametersResponse, error) {
	var response PrequalParametersResponse
	err := db.NewSelect().
		Model(&response).
		Order("created_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		log.Printf("Error fetching latest prequal parameters: %v", err)
		return response, err
	}
	// Check if the latest entry is active, if it's not fetch the last active entry
	if response.Status != "active" {
		var lastActive PrequalParametersResponse
		err := db.NewSelect().
			Model(&lastActive).
			Where("status = ?", "active").
			Order("created_at DESC").
			Limit(1).
			Scan(ctx)
		if err != nil {
			return response, fmt.Errorf("error fetching last active entry: %v", err)
		}
		// If the latest is not active, return the last active entry
		response = lastActive
	}

	return response, nil
}

// Insert new row
func AddPrequalParametersResponse(ctx context.Context, response PrequalParametersResponse, activateId *int64) error {
	var count int
	count, err := db.NewSelect().
		Model((*PrequalParametersResponse)(nil)).
		Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to check table count: %v", err)
	}

	if count == 0 {
		response.Status = "active" // First entry is always active
	} else {
		response.Status = "inactive" // Default status for subsequent entries
	}

	// Insert the new record
	_, err = db.NewInsert().Model(&response).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to add prequal parameters response: %v", err)
	}

	// If an activateID is provided, activate the specified entry
	if activateId != nil {
		// Fetch the specified entry
		var activateResponse PrequalParametersResponse
		err := db.NewSelect().
			Model(&activateResponse).
			Where("id = ?", *activateId).
			Scan(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch prequal parameters with ID %d: %v", *activateId, err)
		}

		// Ensure the entry is currently inactive
		if activateResponse.Status != "inactive" {
			return fmt.Errorf("only 'inactive' entries can be activated")
		}

		// Set all other entries to inactive
		_, err = db.NewUpdate().
			Model((*PrequalParametersResponse)(nil)).
			Set("status = ?", "inactive").
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to deactivate other entries: %v", err)
		}

		// Activate the specified entry
		_, err = db.NewUpdate().
			Model(&activateResponse).
			Set("status = ?", "active").
			Set("updated_at = ?", time.Now()).
			Where("id = ?", *activateId).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed to activate Prequal Parameters with ID %d: %v", *activateId, err)
		}

		// Log
		activationMessage := fmt.Sprintf("Prequal Parameters with ID %d is now active", *activateId)
		if logErr := LogActivity(ctx, "success", activationMessage, activateId); logErr != nil {
			return fmt.Errorf("failed to log activity for activation: %v", logErr)
		}
	}

	// Log
	message := fmt.Sprintf(
		"Prequal request added: MaxLifeTime=%d, PoolSize=%d, ProbeFactor=%.2f, ProbeRemoveFactor=%d, Mu=%d",
		response.MaxLifeTime, response.PoolSize, response.ProbeFactor, response.ProbeRemoveFactor, response.Mu,
	)
	if logErr := LogActivity(ctx, "success", message, nil); logErr != nil {
		return fmt.Errorf("failed to log activity for prequal parameters response: %v", logErr)
	}

	return nil
}
