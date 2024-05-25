package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Salgac/wac-patient-be/pkg/db"
	"github.com/Salgac/wac-patient-be/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GET
func GetAmbulances(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		collection := db.GetCollection(client, "ambulances")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		var ambulances []models.Ambulance
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var ambulance models.Ambulance
			cursor.Decode(&ambulance)
			ambulances = append(ambulances, ambulance)
		}

		if err := cursor.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(ambulances)
	}
}
