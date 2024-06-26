package handling

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
func GetPatients(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		collection := db.GetCollection(client, "patients")

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
}
