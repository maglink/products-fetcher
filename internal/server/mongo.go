package server

import (
	"context"
	"github.com/maglink/products-fetcher/pkg/messages"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
	"time"
)

func mongoConnect(cfg MongoConfig, ctx context.Context) (*mongo.Collection, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Url))
	if err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	mongoCollection := client.Database(cfg.Database).Collection(cfg.Collection)

	_, err = mongoCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"name": 1,
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	return mongoCollection, nil
}

func mongoUpsertProducts(collection *mongo.Collection, productEntries []CsvProductEntry, ctx context.Context) error {
	now := time.Now().Unix()
	var operations []mongo.WriteModel
	for _, entry := range productEntries {
		operation := mongo.NewUpdateOneModel()
		operation.SetFilter(bson.M{
			"name": entry.Name,
		})
		operation.SetUpdate(bson.M{
			"$set": bson.D{
				{"name", entry.Name},
				{"price", entry.Price},
				{"last_update", now},
			},
			"$inc": bson.D{
				{"updates_count", 1},
			},
		})
		operation.SetUpsert(true)
		operations = append(operations, operation)
	}

	_, err := collection.BulkWrite(ctx, operations)
	if err != nil {
		return err
	}

	return nil
}

func mongoGetProductsList(in *messages.ListRequest, collection *mongo.Collection, ctx context.Context) ([]*messages.ListResponse_ListEntry, error) {
	findOptions := options.Find()
	findOptions.SetLimit(int64(in.Limit))
	findOptions.SetSkip(int64(in.Offset))

	orderOptions := bson.M{}
	for _, order := range in.Order {
		var field = strings.ToLower(order.Field.String())
		var direction = 1
		if order.Direction == messages.ListRequest_OrderOptions_DESC {
			direction = -1
		}
		orderOptions[field] = direction
	}
	findOptions.SetSort(orderOptions)

	var results []*messages.ListResponse_ListEntry
	cur, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var elem messages.ListResponse_ListEntry
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	err = cur.Close(ctx)
	if err != nil {
		return nil, err
	}

	return results, nil
}
