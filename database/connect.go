package database

import (
  "context"
  "fmt"
  "log"
  "os"
  "time"

  "github.com/joho/godotenv"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

var (
  DB_NAME = "goax"
)

func DBInstance() *mongo.Client {
  // load environmental variables
  err := godotenv.Load(".env")
  if err != nil {
    log.Fatalln("Oops! could not load environmental variables")
  }

  // -> mongodb url
  DB_URL := os.Getenv("DB_URL")

  if DB_URL == "" {
    //DB_URL = "mongodb+srv://braswelljr:braswellazu@cluster0.1uecn.mongodb.net/test?retryWrites=true&w=majority"
    DB_URL = "mongodb://localhost:27017"
  }

  // Set client options
  client, err := mongo.NewClient(options.Client().ApplyURI(DB_URL))
  if err != nil {
    log.Fatal(err)
  }

  // client context
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  // cancel connection after timeout
  defer cancel()
  if err = client.Connect(ctx); err != nil {
    log.Fatal(err)
  }
  fmt.Println("Connected to Database on ", DB_URL)

  // return a client instance
  return client
}

// Client instantiate a Database instance
var Client = DBInstance()

// OpenCollection open a MongoDB collection
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
  return client.Database(DB_NAME).Collection(collectionName)
}
