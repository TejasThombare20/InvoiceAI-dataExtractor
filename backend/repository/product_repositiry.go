package repository

import (
	"context"

	"github.com/TejasThombare20/backend/config"
	"github.com/TejasThombare20/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		collection: config.DB.Collection("products"),
	}
}

func (r *ProductRepository) FindAll(ctx context.Context) ([]models.Product, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (*models.Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var product models.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// func (r *ProductRepository) Update(ctx context.Context, id string, product *models.Product) error {
// 	objID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	update := bson.M{"$set": product}
// 	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
// 	return err
// }

func (r *ProductRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
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
