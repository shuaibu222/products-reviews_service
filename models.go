package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

const mongoURL = "mongodb://mongo:27017"

// connect with MongoDB
func init() {
	credential := options.Credential{
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}
	clientOpts := options.Client().ApplyURI(mongoURL).SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Println("Error connecting to MongoDB")
		return
	}

	collection = client.Database("reviews").Collection("reviews")

	// collection instance
	log.Println("Collections instance is ready")
}

type ReviewsEntry struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Msg       string             `bson:"msg" json:"msg"`
	UserId    string             `bson:"user_id" json:"user_id"`
	ProductId string             `bson:"product_id" json:"product_id"`
	Time      string             `bson:"time" json:"time"`
	Date      string             `bson:"date" json:"date"`
}

func (r *ReviewsEntry) CreateReview(entry ReviewsEntry) (*mongo.InsertOneResult, error) {
	review, err := collection.InsertOne(context.Background(), entry)
	if err != nil {
		log.Println("Error inserting into reviews:", err)
		return nil, err
	}

	return review, nil
}

func (u *ReviewsEntry) GetAllReviews(ctx context.Context, productId string) ([]*ReviewsEntry, error) {

	cursor, err := collection.Find(ctx, bson.M{"product_id": productId})
	if err != nil {
		log.Println("Finding this product reviews error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviews []*ReviewsEntry

	for cursor.Next(ctx) {
		var item ReviewsEntry

		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding log into slice:", err)
			return nil, err
		} else {
			reviews = append(reviews, &item)
		}
	}

	return reviews, nil
}

func (u *ReviewsEntry) UpdateReviewById(id string, reviewToUpdate ReviewsEntry, ctx context.Context) (*mongo.UpdateResult, error) {

	Id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Failed to convert id", err)
		return nil, err
	}

	updatedCount, err := collection.UpdateOne(ctx, bson.M{"_id": Id}, bson.M{"$set": reviewToUpdate})
	if err != nil {
		log.Println("Failed to update the review: ", err)
		return nil, err
	}
	return updatedCount, nil
}

func (u *ReviewsEntry) DeleteReview(id string, ctx context.Context) (*mongo.DeleteResult, error) {
	Id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Failed to convert id", err)
		return nil, err
	}
	deleteCount, err := collection.DeleteOne(ctx, bson.M{"_id": Id})
	if err != nil {
		log.Println("Failed to delete review", err)
		return nil, err
	}

	return deleteCount, nil
}
