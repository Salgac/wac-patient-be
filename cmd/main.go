package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Salgac/wac-patient-be/pkg/db"
	"github.com/Salgac/wac-patient-be/pkg/handlers"
	"github.com/Salgac/wac-patient-be/pkg/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func main() {
	router := mux.NewRouter()

	// Database
	client = db.ConnectDB()
	setupDb()

	// Handling
	router.HandleFunc("/ambulances", handlers.GetAmbulances(client)).Methods("GET")

	router.HandleFunc("/patients", handlers.GetPatients(client)).Methods("GET")

	// Setup port
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	// Run
	log.Println("API is running!")
	http.ListenAndServe(":"+httpPort, router)
}

// setup default db values
func setupDb() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create health statuses
	healthConditions := []models.HealthCondition{
		{Timestamp: time.Now().Format(time.RFC3339), Description: "Healthy"},
		{Timestamp: time.Now().AddDate(0, 0, -2).Format(time.RFC3339), Description: "Recovered from flu"},
	}

	// Create an ambulance
	ambulance := models.Ambulance{Name: "Ambulance A"}
	collection := db.GetCollection(client, "ambulances")
	collection.DeleteMany(ctx, bson.M{})
	collection.InsertOne(ctx, ambulance)

	// Create visits
	visits := []models.Visit{
		{Ambulance: ambulance, Timestamp: time.Now().Format(time.RFC3339), Reason: "Routine Checkup"},
		{Ambulance: ambulance, Timestamp: time.Now().AddDate(0, 0, 15).Format(time.RFC3339), Reason: "Emergency"},
	}

	collection = db.GetCollection(client, "patients")
	collection.DeleteMany(ctx, bson.M{})
	collection.InsertOne(ctx, models.Patient{
		FirstName:        "Jozko",
		LastName:         "Mrkvicka",
		Age:              69,
		HealthConditions: healthConditions,
		Visits:           visits,
	})
	collection.InsertOne(ctx, models.Patient{
		FirstName:        "Anka",
		LastName:         "Mrkvickova",
		Age:              14,
		HealthConditions: []models.HealthCondition{},
		Visits:           []models.Visit{},
	})
}
