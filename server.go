package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"context"
	"net/http"

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

func Server() {
	fmt.Println("Running server")
	c := New()
	c.Connect()

	mux := http.NewServeMux()
	// Endpoint paths should only contain nouns
	mux.HandleFunc("/", c.update)
	mux.HandleFunc("/mongo/add", c.add)
	mux.HandleFunc("/mongo/update", c.update)
	mux.HandleFunc("/mongo/remove", c.remove)
	mux.HandleFunc("/mongo/find", c.find)
	mux.HandleFunc("/mongo/dummy", c.dummy)
	// mux.Handler(r)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("error starting server", err)
	} else {
		fmt.Print("Server started!")
	}
}

// curl -X POST http://localhost:8080/mongo/add -H "Content-Type: application/json" -d '{"Name":"Ash", "Age":21, "City": "Gresham"}'
func (c *Client) add(w http.ResponseWriter, r *http.Request) {
	// io.WriteString(w, "This is my website!\n")
	var t Trainer
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
	}

	// myName := r.PostFormValue("Name")
	err = json.Unmarshal(body, &t)
	if err != nil {
		panic(err)
	}
	fmt.Println("Adding to mongo database", t)
	collection := c.mClient.Database("test").Collection("trainers")
	id, err := c.Insert(collection, []*Trainer{&t})
	if err != nil {
		fmt.Print("Error inserting", err)
	} else {
		io.WriteString(w, fmt.Sprintf("Successfully Inserted! Id is %v\n", id[0]))
	}
}

// Not implemented
func (c *Client) update(w http.ResponseWriter, r *http.Request) {

	io.WriteString(w, "Hello, HTTP!\n")
}

func (c *Client) dummy(w http.ResponseWriter, r *http.Request) {

	c.DummyData()
}

type TrainerDelete struct {
	id string
}

// Not setup
func (c *Client) remove(w http.ResponseWriter, r *http.Request) {
	var d TrainerDelete
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
	}

	// myName := r.PostFormValue("Name")
	err = json.Unmarshal(body, &d)
	if err != nil {
		panic(err)
	}
	io.WriteString(w, fmt.Sprintf("Deleting, %s!\n", d.id))
}

// curl http://localhost:8080/mongo/find?name=Ash
func (c *Client) find(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	fmt.Println("QUERY", query["name"])
	filter := bson.D{{"name", query["name"][0]}}
	collection := c.mClient.Database("test").Collection("trainers")
	trainers, err := c.Find(collection, filter)
	if err != nil {
		panic(err)
	}
	io.WriteString(w, fmt.Sprintf("Hello, trainer %s \n", trainers[0].Name))
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
	} else {
		fmt.Println("Ping Successful!")
	}
}

func (c *Client) Disconnect() {
	err := c.mClient.Disconnect(context.TODO())

	if err != nil {
		panic(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func (c *Client) Insert(collection *mongo.Collection, data []*Trainer) ([]interface{}, error) {
	var ids []interface{}
	if len(data) == 1 {
		insertResult, err := collection.InsertOne(context.TODO(), data[0])
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
		ids = append(ids, insertResult.InsertedID)
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

	return ids, nil
}

func (c *Client) Delete(collection *mongo.Collection, filter bson.D) error {
	deleteResult, err := collection.DeleteMany(context.TODO(), filter)
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
	_, err := c.Insert(collection, trainers)
	if err != nil {
		panic(err)
	}
	trainers = []*Trainer{&misty, &brock}
	_, err = c.Insert(collection, trainers)
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
	fmt.Println("updated", updateResult)
	err = c.Delete(collection, filter)
	if err != nil {
		panic(err)
	}
	fmt.Println("Deleted", filter)
	result, err = c.Find(collection, filter)
	if err == nil {
		panic(err)
	}
	fmt.Println("No one found", result)

}

type Trainer struct {
	Name string
	Age  int
	City string
}
