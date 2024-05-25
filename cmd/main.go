package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Salgac/wac-patient-be/pkg/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func main() {
	router := mux.NewRouter()

	// Database
	client = ConnectDB()
	setupDb()

	// Handling
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("Hello World")
	})

	router.HandleFunc("/patients", getPatients).Methods("GET")

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	log.Println("API is running!")
	http.ListenAndServe(":"+httpPort, router)
}

func getPatients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	collection := GetCollection(client)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var patients []models.Patient
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var patient models.Patient
		cursor.Decode(&patient)
		patients = append(patients, patient)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(patients)
}

// setup default db values
func setupDb() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create health statuses
	healthStatuses := []models.HealthStatus{
		{Id: 1, Description: "Healthy"},
		{Id: 2, Description: "Recovered from flu"},
	}

	// Create an ambulance
	ambulance := models.Ambulance{Id: 1, Name: "Ambulance A"}

	// Create visits
	visits := []models.Visit{
		{Id: 1, Ambulance: ambulance, Timestamp: time.Now().Format(time.RFC3339), Reason: "Routine Checkup"},
		{Id: 2, Ambulance: ambulance, Timestamp: time.Now().AddDate(0, 0, 15).Format(time.RFC3339), Reason: "Emergency"},
	}

	collection := GetCollection(client)
	collection.DeleteMany(ctx, bson.M{})
	collection.InsertOne(ctx, models.Patient{
		Id:             69,
		FirstName:      "Jozko",
		LastName:       "Mrkvicka",
		HealthStatuses: healthStatuses,
		Visits:         visits,
	})
	collection.InsertOne(ctx, models.Patient{
		Id:             420,
		FirstName:      "Anka",
		LastName:       "Mrkvickova",
		HealthStatuses: []models.HealthStatus{},
		Visits:         []models.Visit{},
	})
}
