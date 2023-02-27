package main

import (
	"fmt"
	"log"
	"os"

	"context"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Uri string `yaml:"uri", envconfig:"URI"`
}

type Client struct {
	Config
	mClient *mongo.Client
}

func (c *Client) CreateDefaultConfig() *Config {

	return &Config{
		Uri: "",
	}
}

func (c *Client) ValidateServer(config *Config) {

}

func New() *Client {
	var cfg Config
	readFile(&cfg)
	readEnv(&cfg)
	c := &Client{
		Config: cfg,
	}
	return c
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(cfg *Config) {
	f, err := os.Open("config.yml")
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}

// Connect through mongo client
func (c *Client) Connect() {
	clientOptions := options.Client().ApplyURI(c.Config.Uri)
	mongoClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	c.mClient = mongoClient

}

// Ping mongo client
func (c *Client) Ping() {
	// Check the connection
	err := c.mClient.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
}

func (c *Client) Disconnect() {
	err := c.mClient.Disconnect(context.TODO())

	if err != nil {
		panic(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func (c *Client) Insert(collection *mongo.Collection, data []*Trainer) error {
	if len(data) == 1 {
		insertResult, err := collection.InsertOne(context.TODO(), data[0])
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	} else {
		p := make([]interface{}, len(data))
		for i, v := range data {
			p[i] = v
		}
		// p := []interface{}{&data}
		insertManyResult, err := collection.InsertMany(context.TODO(), p)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
	}

	return nil
}

func (c *Client) Delete(collection *mongo.Collection, filter bson.D) error {
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		return nil
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	return nil
}
func (c *Client) Find(collection *mongo.Collection, filter bson.D) ([]*Trainer, error) {
	var results []*Trainer
	var result *Trainer
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	results = append(results, result)
	fmt.Printf("Found documents: %+v\n", results)
	return results, nil
}
func (c *Client) Update(collection *mongo.Collection, filter bson.D, update bson.D) (*mongo.UpdateResult, error) {
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	return updateResult, nil
}

func (c *Client) DummyData() {
	collection := c.mClient.Database("test").Collection("trainers")
	fmt.Println("collection", collection)

	ash := Trainer{"Ash", 10, "Pallet Town"}
	misty := Trainer{"Misty", 10, "Cerulean City"}
	brock := Trainer{"Brock", 15, "Pewter City"}

	filter := bson.D{{"name", "Ash"}}

	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}

	trainers := []*Trainer{&ash}
	err := c.Insert(collection, trainers)
	if err != nil {
		panic(err)
	}
	trainers = []*Trainer{&misty, &brock}
	err = c.Insert(collection, trainers)
	if err != nil {
		panic(err)
	}
	result, err := c.Find(collection, filter)
	if err != nil {
		panic(err)
	}
	fmt.Println("Found", result)
	updateResult, err := c.Update(collection, filter, update)
	if err != nil {
		panic(err)
	}
	fmt.Println("Found", updateResult)
	err = c.Delete(collection, filter)
	if err != nil {
		panic(err)
	}
}

type Trainer struct {
	Name string
	Age  int
	City string
}

func main1() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Disconnect(client)
}
