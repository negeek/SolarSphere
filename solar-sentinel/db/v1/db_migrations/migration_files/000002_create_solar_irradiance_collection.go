package main

import (
	"os"
	"time"
	"context"
	"log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/joho/godotenv"
)

func connect(connString string, dbName string)(context.Context, context.CancelFunc, error) {
	var err error
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(connString).SetServerAPIOptions(serverAPI))

	if err != nil {
		return ctx, cancel, err
	}
	
	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(ctx, bson.D{{"ping", 1}}).Decode(&result); err != nil {
		return ctx, cancel, err
	}

	//db
	mongoDB = client.Database(dbName)

	return ctx, cancel, nil
}

func disconnect(ctx context.Context, cancel context.CancelFunc){
	defer cancel()
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func getEnv(){
	appEnv:=os.Getenv("APP_ENV")
	if appEnv=="dev"{
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

// The above is perequisite for your migration functionality.
// Write your logic below
var (
	irrSchema primitive.M = bson.M{
									"$jsonSchema":bson.M{
										"bsonType": "object",
										"required": []string{"_id", "device_id"},
										"properties": bson.M{
											"_id": bson.M{
												"bsonType": "string",
												"description": "id is required and must be a string",
											},
											"device_id": bson.M{
												"bsonType": "string",
												"description": "device_id is required and must be a string",
											},
										},
									},
								}

	irrOptions = &options.CreateCollectionOptions{}
	err error
	IRR_COLLECTION string = "solar-irradiance-data"
	client *mongo.Client
	mongoDB *mongo.Database
)

func MakeMigration(){
	// Optional
	getEnv()

	// Connect to DB
	dbctx, dbcancel, err:= connect(os.Getenv("DATABASE_URL"), os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}
	
	// Create solar-irradiance collection
	irrOptions.SetValidator(irrSchema)
	err = mongoDB.CreateCollection(context.Background(), IRR_COLLECTION, irrOptions)
	if err != nil {
		log.Fatal(err)
	}
	
	// End connection
	defer disconnect(dbctx,dbcancel)
}