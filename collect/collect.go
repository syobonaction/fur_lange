package collect

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/syobonaction/fur_lange/datasources"
	model "github.com/syobonaction/fur_lange/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var ctx = context.TODO()

func collect() {
	loadEnv()
	connectMongo()

	var partnerData []*model.Partner

	//set parameters here for the API call
	startRecord := 7000
	size := 50 //seems to cap at 50
	endRecord := 7500

	for i := startRecord; i < endRecord; i += size {
		partnerData = append(partnerData, datasources.GetAWSPartners(i, size)...)
		fmt.Println(i)
	}

	for _, p := range partnerData {
		createRecord(p)
	}
}

func connectMongo() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	url := fmt.Sprintf("mongodb://%s:%s", host, port)

	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("thesis").Collection("partners")
}

func createRecord(partner *model.Partner) error {
	_, err := collection.InsertOne(ctx, partner)
	return err
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
