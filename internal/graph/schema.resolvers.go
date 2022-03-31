package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"delivery-service/internal/core/domain"
	"delivery-service/internal/graph/generated"
	"delivery-service/internal/graph/model"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *deliveryResolver) ID(ctx context.Context, obj *domain.Delivery) (string, error) {
	return obj.ID.String(), nil
}

func (r *deliveryResolver) PickupTime(ctx context.Context, obj *domain.Delivery) (int, error) {
	return int(obj.PickupTime.Unix()), nil
}

func (r *deliveryResolver) Route(ctx context.Context, obj *domain.Delivery) ([]*domain.Location, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *deliveryResolver) DeliveryTime(ctx context.Context, obj *domain.Delivery) (*int, error) {
	unix := int(obj.PickupTime.Unix())
	return &unix, nil
}

func (r *mutationResolver) CreateDelivery(ctx context.Context, input model.DeliveryInput) (*domain.Delivery, error) {
	pickupPoint, err := domain.NewLocation(input.PickupPoint.Latitude, input.PickupPoint.Longitude)

	if err != nil {
		return nil, err
	}

	deliveryPoint, err := domain.NewLocation(input.DeliveryPoint.Latitude, input.DeliveryPoint.Longitude)

	if err != nil {
		return nil, err
	}

	parcelId, err := uuid.Parse(input.ParcelID)

	if err != nil {
		return nil, err
	}

	delivery, err := r.DeliveryService.Create(parcelId, pickupPoint, deliveryPoint, time.Unix(int64(input.PickupTime), 0))

	if err != nil {
		return nil, err
	}

	return &delivery, nil
}

func (r *mutationResolver) AssignRider(ctx context.Context, id string, input *string) (*domain.Delivery, error) {
	deliveryId, err := uuid.Parse(id)

	if err != nil {
		return nil, err
	}

	riderId, err := uuid.Parse(*input)

	if err != nil {
		return nil, err
	}

	delivery, err := r.DeliveryService.AssignRider(deliveryId, riderId)

	if err != nil {
		return nil, err
	}

	return &delivery, nil
}

func (r *mutationResolver) CompleteDelivery(ctx context.Context, id string) (*domain.Delivery, error) {
	deliveryId, err := uuid.Parse(id)

	if err != nil {
		return nil, err
	}

	delivery, err := r.DeliveryService.CompleteDelivery(deliveryId)

	if err != nil {
		return nil, err
	}

	return &delivery, nil
}

func (r *parcelResolver) ID(ctx context.Context, obj *domain.Parcel) (string, error) {
	return obj.ID.String(), nil
}

func (r *queryResolver) Deliveries(ctx context.Context) ([]*domain.Delivery, error) {
	deliveries, err := r.DeliveryService.GetAll()

	if err != nil {
		return nil, err
	}

	var deliveryModels []*domain.Delivery
	for i := range deliveries {
		deliveryModels = append(deliveryModels, &deliveries[i])
	}

	return deliveryModels, nil
}

func (r *queryResolver) Delivery(ctx context.Context, id string) (*domain.Delivery, error) {
	deliveryId, err := uuid.Parse(id)

	if err != nil {
		return nil, err
	}

	delivery, err := r.DeliveryService.Get(deliveryId)

	if err != nil {
		return nil, err
	}

	return &delivery, nil
}

func (r *riderResolver) ID(ctx context.Context, obj *domain.Rider) (string, error) {
	return obj.ID.String(), nil
}

// Delivery returns generated.DeliveryResolver implementation.
func (r *Resolver) Delivery() generated.DeliveryResolver { return &deliveryResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Parcel returns generated.ParcelResolver implementation.
func (r *Resolver) Parcel() generated.ParcelResolver { return &parcelResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Rider returns generated.RiderResolver implementation.
func (r *Resolver) Rider() generated.RiderResolver { return &riderResolver{r} }

type deliveryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type parcelResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type riderResolver struct{ *Resolver }
