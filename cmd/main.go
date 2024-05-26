package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Salgac/wac-patient-be/pkg/db"
	"github.com/Salgac/wac-patient-be/pkg/handling"
	"github.com/Salgac/wac-patient-be/pkg/models"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func main() {
	router := mux.NewRouter()

	// Define CORS options
	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH"}),
		handlers.AllowedHeaders([]string{"Origin", "Authorization", "Content-Type", "Accept", "Cache-Control"}),
		handlers.ExposedHeaders([]string{""}),
		handlers.MaxAge(86400),
	)

	// Database
	client = db.ConnectDB()
	setupDb()

	// Handling
	router.HandleFunc("/ambulances", handling.GetAmbulances(client)).Methods("GET")

	router.HandleFunc("/patients", handling.GetPatients(client)).Methods("GET")

	router.HandleFunc("/patients/{id}/conditions", handling.AddHealthCondition(client)).Methods("POST")

	router.HandleFunc("/patients/{id}/visits", handling.AddVisit(client)).Methods("POST")
	router.HandleFunc("/patients/{patientId}/visits/{visitId}", handling.DeleteVisit(client)).Methods("DELETE")
	router.HandleFunc("/patients/{patientId}/visits/{visitId}", handling.UpdateVisit(client)).Methods("PUT")

	// Setup port
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	// Run
	log.Println("API is running!")
	loggedRouter := handlers.LoggingHandler(log.Writer(), router)
	http.ListenAndServe(":"+httpPort, corsOptions(loggedRouter))
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
	collection.InsertOne(ctx, models.Ambulance{Name: "Ambulance B"})
	collection.InsertOne(ctx, models.Ambulance{Name: "Ambulance C"})

	// Create visits
	timestamp := time.Now().Truncate(60*time.Second).Truncate(60*time.Minute)
	visits := []models.Visit{
		{Id: primitive.NewObjectID(), Ambulance: ambulance, Timestamp: timestamp.AddDate(0, 0, 5).Format(time.RFC3339), Reason: "Routine Checkup", Status: "requested"},
		{Id: primitive.NewObjectID(), Ambulance: ambulance, Timestamp: timestamp.AddDate(0, 0, 15).Format(time.RFC3339), Reason: "Emergency", Status: "done"},
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
