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
func AddVisit(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		collection := db.GetCollection(client, "patients")

		patientID, err := primitive.ObjectIDFromHex(params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var visit models.Visit
		_ = json.NewDecoder(r.Body).Decode(&visit)

		filter := bson.M{"_id": patientID}
		update := bson.M{"$push": bson.M{"visits": visit}}

		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(visit)
	}
}

// PUT
func UpdateVisit(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		collection := db.GetCollection(client, "patients")

		patientID, err := primitive.ObjectIDFromHex(params["patientId"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		visitID, err := primitive.ObjectIDFromHex(params["visitId"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var visit models.Visit
		_ = json.NewDecoder(r.Body).Decode(&visit)

		filter := bson.M{"_id": patientID, "visits._id": visitID}
		update := bson.M{
			"$set": bson.M{
				"visits.$.ambulance": visit.Ambulance,
				"visits.$.timestamp": visit.Timestamp,
				"visits.$.reason":    visit.Reason,
			},
		}

		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("Visit updated")
	}
}

func DeleteVisit(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		collection := db.GetCollection(client, "patients")

		patientID, err := primitive.ObjectIDFromHex(params["patientId"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		visitID, err := primitive.ObjectIDFromHex(params["visitId"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filter := bson.M{"_id": patientID}
		update := bson.M{"$pull": bson.M{"visits": bson.M{"_id": visitID}}}

		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("Visit deleted")
	}
}
