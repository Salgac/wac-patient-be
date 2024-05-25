package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Patient struct {
	Id            int
	FirstName     string
	LastName      string
	HeathStatuses []HealthStatus
	Visits        []Visit
}

type Visit struct {
	Id     int
	Time   string // todo
	Reason string
}

type HealthStatus struct {
	Id int
}

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

	var patients []Patient
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var patient Patient
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

	collection := GetCollection(client)
	collection.DeleteMany(ctx, bson.M{})
	collection.InsertOne(ctx, Patient{Id: 69, FirstName: "Jozko", LastName: "Mrkvicka"})
	collection.InsertOne(ctx, Patient{Id: 420, FirstName: "Anka", LastName: "Mrkvickova"})
}
