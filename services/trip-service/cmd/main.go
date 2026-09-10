package main

import (
	"context"
	"fmt"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	ctx := context.Background()

	imMem := repository.NewInmemTripRepository()
	svc := service.NewTripService(imMem)

	t, err := svc.CreateTrip(ctx, &domain.RideFareModel{
		ID:                primitive.NewObjectID(),
		UserID:            "123",
		PackageSlug:       "sedan",
		TotalPriceInCents: 10000,
		ExpiresAt:         time.Now().Add(1 * time.Hour),
	})
	if err != nil {
		log.Fatalf("failed to create trip: %v", err)
	}
	fmt.Println(t)

	for {
		time.Sleep(1 * time.Second)
	}
}
