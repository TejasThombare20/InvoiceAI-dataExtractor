package repository

import (
	"context"

	"github.com/TejasThombare20/backend/config"
	"github.com/TejasThombare20/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerRepository struct {
	collection *mongo.Collection
}

func NewCustomerRepository() *CustomerRepository {
	return &CustomerRepository{
		collection: config.DB.Collection("customers"),
	}
}

func (r *CustomerRepository) FindAll(ctx context.Context) ([]models.Customer, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var customers []models.Customer
	if err = cursor.All(ctx, &customers); err != nil {
		return nil, err
	}

	return customers, nil
}

func (r *CustomerRepository) FindByID(ctx context.Context, id string) (*models.Customer, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var customer models.Customer
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&customer)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

// func (r *CustomerRepository) Update(ctx context.Context, id string, customer *models.Customer) error {
// 	objID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	update := bson.M{"$set": customer}
// 	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
// 	return err
// }

func (r *CustomerRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{"$set": updates}
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		update,
	)
	return err
}
