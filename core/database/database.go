package graph_database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jeconias/graph-service/core/database/schema"
	"github.com/jeconias/graph-service/core/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database struct {
	Conn *mongo.Client
}

func InitDB() (*Database, error) {
	mongoHost := utils.GetEnvValue("MONGO_HOST")
	mongoPort := utils.GetEnvValue("MONGO_PORT")
	mongoUser := utils.GetEnvValue("MONGO_USER")
	mongoPass := utils.GetEnvValue("MONGO_PASS")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", mongoUser, mongoPass, mongoHost, mongoPort)))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("MongoDB not started: %s\n", err))
	}

	pingErr := client.Ping(context.TODO(), &readpref.ReadPref{})
	if pingErr != nil {
		return nil, errors.New(fmt.Sprintf("Failed on Ping: %s\n", err))
	}

	fmt.Println("MongoDB started")

	return &Database{Conn: client}, nil
}

func (v *Database) GetDB() *mongo.Database {
	return v.Conn.Database("graph_service")
}

func (v *Database) GetCollectionVertice() *mongo.Collection {
	return v.GetDB().Collection("vertices")
}

func (v *Database) UpsertManyVertice(items []schema.VerticeSchema) (*mongo.BulkWriteResult, error) {

	vertice := v.GetCollectionVertice()

	upsert := true
	operations := []mongo.WriteModel{}

	for _, item := range items {
		operation := mongo.NewUpdateOneModel()
		operation.Upsert = &upsert

		operation.Filter = bson.M{"from": item.From, "to": item.To}
		operation.Update = bson.M{"$set": bson.M{"from": item.From, "to": item.To}, "$push": bson.M{"infos": bson.M{"$each": item.Infos}}}

		operations = append(operations, operation)
	}

	result, err := vertice.BulkWrite(context.Background(), operations)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Erro on UpsertManyVertice: %v", err))
	}

	return result, nil
}

func (v *Database) InsertVertice(item schema.VerticeSchema) (*schema.VerticeSchema, error) {

	vertice := v.GetCollectionVertice()

	result := &schema.VerticeSchema{}

	upsert := true
	returnDocument := options.After
	singleResult := vertice.FindOneAndUpdate(context.Background(), bson.M{
		"from": item.From,
		"to":   item.To,
	}, bson.M{"$set": bson.M{"from": item.From, "to": item.To}, "$push": bson.M{"infos": bson.M{"$each": item.Infos}}}, &options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &returnDocument,
	}).Decode(result)

	if singleResult != nil {
		return nil, errors.New(fmt.Sprintf("Erro on InsertVertice: %v", singleResult))
	}

	return &schema.VerticeSchema{
		ID:   result.ID,
		From: item.From,
		To:   item.To,
	}, nil
}

func (v *Database) ListVertice() ([]*schema.VerticeSchema, error) {
	collection := v.GetCollectionVertice()

	ctx := context.Background()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Erro on ListVertice: %v", err))
	}

	defer cursor.Close(ctx)

	result := []*schema.VerticeSchema{}
	for cursor.Next(ctx) {
		vertice := &schema.VerticeSchema{}

		if err = cursor.Decode(vertice); err != nil {
			return nil, errors.New(fmt.Sprintf("Erro on ListVertice: %v", err))
		}

		result = append(result, vertice)
	}

	return result, nil
}
