package handling

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Salgac/wac-patient-be/pkg/db"
	"github.com/Salgac/wac-patient-be/pkg/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// POST
func AddHealthCondition(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		collection := db.GetCollection(client, "patients")

		patientID, err := primitive.ObjectIDFromHex(params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var healthCondition models.HealthCondition
		_ = json.NewDecoder(r.Body).Decode(&healthCondition)

		filter := bson.M{"_id": patientID}
		update := bson.M{"$push": bson.M{"healthconditions": healthCondition}}

		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(healthCondition)
	}
}
