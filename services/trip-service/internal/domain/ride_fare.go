package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)
type RideFareModel struct {
	ID primitive.ObjectID `bson:"_id"`
	UserID string `bson:"user_id"`
	PackageSlug string `bson:"package_slug"`
	TotalPriceInCents int64 `bson:"total_price_in_cents"`
	ExpiresAt time.Time `bson:"expires_at"`
}