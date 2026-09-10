package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripModel struct {
	ID       primitive.ObjectID `bson:"_id"`
	UserID   string             `bson:"user_id"`
	Status   string             `bson:"status"`
	RideFare *RideFareModel     `bson:"ride_fare"`
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
}
