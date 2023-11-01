package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shuaibu222/p-reviews/reviews"
	"go.mongodb.org/mongo-driver/bson"
)

// here should be the grpc handlers
type ReviewsServer struct {
	reviews.UnimplementedReviewsServiceServer
}

func (r *ReviewsServer) CreateReview(ctx context.Context, req *reviews.ReviewRequest) (*reviews.ReviewResponse, error) {
	input := req.GetReviewEntry()

	currentTime := time.Now()

	reviewEntry := ReviewsEntry{
		Msg:       input.Msg,
		UserId:    input.UserId,
		ProductId: input.ProductId,
		Time:      fmt.Sprint(currentTime.Clock()),
		Date:      fmt.Sprint(currentTime.Date()),
	}

	rev, err := reviewEntry.CreateReview(reviewEntry)
	if err != nil {
		log.Printf("Error inserting review: %v", err)
	}

	log.Println("Review inserted successfully: ", rev.InsertedID)

	// get it immediately from the database
	var review ReviewsEntry
	err = collection.FindOne(context.Background(), bson.M{"_id": rev.InsertedID}).Decode(&review)
	if err != nil {
		return nil, err
	}

	reviewInserted := reviews.Review{
		Id:        review.ID.String(),
		Msg:       review.Msg,
		UserId:    review.UserId,
		ProductId: review.ProductId,
		Time:      review.Time,
		Date:      review.Date,
	}

	response := &reviews.ReviewResponse{
		Response: &reviewInserted,
	}

	res := []byte(review.Msg)

	err = app.SendReviewToRabbitmq(res)
	if err != nil {
		log.Println(err)
	}

	return response, nil
}

func (u *ReviewsServer) GetReviews(req *reviews.ProductId, stream reviews.ReviewsService_GetReviewsServer) error {
	reviewEntry := ReviewsEntry{}

	reviewsResult, err := reviewEntry.GetAllReviews(context.Background(), req.Id)
	if err != nil {
		log.Println("Failed to get reviews", err)
	}

	for _, reviewEntry := range reviewsResult {
		reviewResponse := &reviews.Review{
			Id:        reviewEntry.ID.Hex(),
			Msg:       reviewEntry.Msg,
			UserId:    reviewEntry.UserId,
			ProductId: reviewEntry.ProductId,
			Time:      reviewEntry.Time,
			Date:      reviewEntry.Date,
		}

		// send each review as a stream(continiously immediately its posted) to the client
		if err := stream.Send(reviewResponse); err != nil {
			log.Println("Error sending review to the client:", err)
			return err
		}
	}

	// Return the response
	return nil
}

func (u *ReviewsServer) UpdateReview(ctx context.Context, req *reviews.Review) (*reviews.Count, error) {
	reviewEntry := ReviewsEntry{
		Msg:       req.Msg,
		UserId:    req.UserId,
		ProductId: req.ProductId,
	}

	updateCount, err := reviewEntry.UpdateReviewById(req.Id, reviewEntry, ctx)
	if err != nil {
		log.Println(err)
	}

	log.Println(updateCount.ModifiedCount)

	reviewResult := &reviews.Count{
		Count: fmt.Sprint(updateCount.ModifiedCount),
	}

	return reviewResult, nil
}

func (u *ReviewsServer) DeleteReview(ctx context.Context, req *reviews.ReviewId) (*reviews.Count, error) {
	reviewEntry := ReviewsEntry{}

	deletedCount, err := reviewEntry.DeleteReview(req.Id, ctx)
	if err != nil {
		log.Println(err)
	}

	res := &reviews.Count{
		Count: fmt.Sprint(deletedCount.DeletedCount),
	}

	return res, nil
}
