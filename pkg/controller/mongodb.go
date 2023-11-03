package controller

import (
	"context"
	"os"
	"time"

	"github.com/wreckitkenny/vngitpub/pkg/utils"
	"github.com/wreckitkenny/vngitpub/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// State structs deployment status
type State struct {
	ID    		primitive.ObjectID 	`bson:"_id" json:"id,omitempty"`
	Image		string				`bson:"image" json:"image,omitempty"`
	OldTag		string				`bson:"oldtag" json:"oldtag,omitempty"`
	NewTag		string				`bson:"newtag" json:"newtag,omitempty"`
	Cluster		string				`bson:"cluster" json:"cluster,omitempty"`
	BlobName 	string				`bson:"blob" json:"blob,omitempty"`
	Time 		string				`bson:"time" json:"time,omitempty"`
	Status		string				`bson:"status" json:"status,omitempty"`
	Metadata	string				`bson:"metadata" json:"metadata,omitempty"`
}

func connectMongo() (*mongo.Client, string, ) {
	logger := utils.ConfigZap()

	mongoAddress := os.Getenv("MONGO_ADDRESS")
	mongoDBName := os.Getenv("MONGO_DBNAME")
	mongoUsername := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASS")

	uri := "mongodb://" + mongoUsername + ":" + mongoPassword + "@" + mongoAddress + "/?retryWrites=true&w=majority"

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI).SetConnectTimeout(3*time.Second)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		logger.Errorf("Connecting to MongoDB...FAILED: %s", err)
	}
	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		log.Panic(err)
	// 	}
	// }()

	return client, mongoDBName
}

// ValidateMongoConnection makes sure that the connection to Mongo works
func ValidateMongoConnection() {
	logger := utils.ConfigZap()
	client, mongoDBName := connectMongo()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Check connection to MongoDB
	var result bson.M
	if err := client.Database(mongoDBName).RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		logger.Errorf("Connecting to MongoDB...FAILED: %s", err)
	} else {
		logger.Info("Connecting to MongoDB...OK")
	}
}

// LoadState loads documents from MongoDB
func LoadState() []State {
	logger := utils.ConfigZap()
	client, mongoDBName := connectMongo()

	coll := client.Database(mongoDBName).Collection("status")
	filter := bson.D{}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		logger.Errorf("Loading states from MongoDB...FAILED: %s", err)
	} else {
		logger.Debug("Loading states from MongoDB...OK")
	}

	var states []State
	if err = cursor.All(context.TODO(), &states); err != nil {
		logger.Debugf("Getting all documents from MongoDB...FAILED: %s", err)
	}

	return states
}

// SaveUser writes a new user information to collection
func SaveUser(username string, email string, encodedPass string, fullname string, department string, role string) error {
	client, mongoDBName := connectMongo()

	coll := client.Database(mongoDBName).Collection("users")
	newUser := model.User{Username: username, Email: email, Password: encodedPass, FullName: fullname, Department: department, Role: role}

	_, err := coll.InsertOne(context.TODO(), newUser)
	if err != nil {
		return err
	}

	return nil
}

// FindUser queries collection users for user information
func FindUser(key string, value string,) (model.User, error) {
	logger := utils.ConfigZap()
	client, mongoDBName := connectMongo()

	coll := client.Database(mongoDBName).Collection("users")
	filter := bson.D{{Key: key, Value: value}}

	var result model.User
	err := coll.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Debug(err)
			return result, err
		}
	}

	return result, nil
}