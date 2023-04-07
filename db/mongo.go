package db

import (
	"context"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var ctx context.Context
var collection *mongo.Collection
var router *mux.Router

func MongoClient(collection_name string) (*mux.Router, context.Context, *mongo.Client, *mongo.Collection) {
	router = mux.NewRouter()
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://admin:CEMt1994@sandbox.0sac2.mongodb.net/test"))
	collection = client.Database("horsea_db").Collection(collection_name) //user_collection

	return router, ctx, client, collection
}
