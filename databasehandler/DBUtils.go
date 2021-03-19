package databasehandler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"

	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DatasetCollectionName The name of the mongo collection that stores the dataset
const DatasetCollectionName = "Datasets"

// DatasetVersionCollectionName The name of the mongo collection that stores the dataset versions
const DatasetVersionCollectionName = "DatasetVersions"

// DatasetObjectsCollectionName The name of the mongo collection that stores the dataset objects
const DatasetObjectsCollectionName = "DatasetObjects"

// ProjectsCollectionsName The name of the mongo collection that stores the project information
const ProjectsCollectionsName = "Projects"

//DatsetObjectGroupsCollectionName The Collection name for the dataset object groups
const DatsetObjectGroupsCollectionName = "DatasetObjectGroups"

// APITokenCollectionName The name of the mongo collection that stores the api tokens
const APITokenCollectionName = "APITokens"

//DBUtilsHandler The base handler for MongoDB related database action
type DBUtilsHandler struct {
	MongoClient                 *mongo.Client
	MongoDefaultContext         context.Context
	DatasetDatabaseName         string
	AuthDatabaseName            string
	APITokenCollectionName      string
	DatasetCollName             string
	DatasetVersionCollName      string
	DatasetObjectsCollName      string
	DatasetObjectsGroupCollName string
	AuthProjectCollectionName   string
}

//NewDBUtilsHandler Creates a new handler that handles database interaction
func NewDBUtilsHandler(client *mongo.Client, ctx context.Context) (*DBUtilsHandler, error) {
	handler := DBUtilsHandler{
		MongoClient:                 client,
		MongoDefaultContext:         ctx,
		DatasetDatabaseName:         "Dataset",
		DatasetCollName:             "Dataset",
		DatasetVersionCollName:      "DatasetVersion",
		AuthDatabaseName:            "Authentication",
		APITokenCollectionName:      "APIToken",
		DatasetObjectsCollName:      "Objects",
		DatasetObjectsGroupCollName: "ObjectGroups",
		AuthProjectCollectionName:   "AuthProjects",
	}

	return &handler, nil
}

// GetDatasetCollection Returns the collection that stores the dataset entries
func (handler *DBUtilsHandler) GetDatasetCollection() *mongo.Collection {
	return handler.GetManagementDatabase().Collection(handler.DatasetCollName)
}

// GetDatasetVersionCollection Returns the collection that stores the dataset version entries
func (handler *DBUtilsHandler) GetDatasetVersionCollection() *mongo.Collection {
	return handler.GetManagementDatabase().Collection(handler.DatasetVersionCollName)
}

// GetDatasetObjectsCollection Returns the collection that stores the dataset object entries
func (handler *DBUtilsHandler) GetDatasetObjectsCollection() *mongo.Collection {
	return handler.GetManagementDatabase().Collection(handler.DatasetObjectsCollName)
}

//GetDatasetObjectGroupCollection Returns the collection that stores the dataset object group entries
func (handler *DBUtilsHandler) GetDatasetObjectGroupCollection() *mongo.Collection {
	return handler.GetManagementDatabase().Collection(handler.DatasetObjectsGroupCollName)
}

// GetManagementDatabase Returns a handler to the default mongodb management database
func (handler *DBUtilsHandler) GetManagementDatabase() *mongo.Database {
	return handler.MongoClient.Database(handler.DatasetDatabaseName)
}

//GetProjectCollection Returns the collection for projects
func (handler *DBUtilsHandler) GetProjectCollection() *mongo.Collection {
	return handler.MongoClient.Database(handler.AuthDatabaseName).Collection(handler.AuthProjectCollectionName)
}

//GetTokenCollection Returns the collection for the stored tokens
func (handler *DBUtilsHandler) GetTokenCollection() *mongo.Collection {
	return handler.MongoClient.Database(handler.AuthDatabaseName).Collection(handler.APITokenCollectionName)
}

//Insert Inserts a given value into a given collection and decodes the inserted value into the given decode value
func (handler *DBUtilsHandler) Insert(collection *mongo.Collection, insertValue interface{}, decodeValue interface{}) error {
	insertedResult, err := collection.InsertOne(handler.MongoDefaultContext, &insertValue)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	var oid primitive.ObjectID
	var ok bool

	if oid, ok = insertedResult.InsertedID.(primitive.ObjectID); !ok {
		return errors.New("Error decoding result id")
	}

	result := collection.FindOne(handler.MongoDefaultContext, bson.M{
		"_id": oid,
	})

	if result.Err() != nil {
		log.Println(result.Err().Error())
		return err
	}

	err = result.Decode(decodeValue)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// CreateNewMongoConnection Creates a new mongodb connection
func CreateNewMongoConnection(ctx context.Context) (*mongo.Client, error) {
	mongoDBURL := viper.GetString("Config.Database.Mongo.URL")
	if mongoDBURL == "" {
		err := errors.New("MongoDB URL required to be present in config")
		log.Println(err.Error())
		return nil, err
	}

	authSource := viper.GetString("Config.Database.Mongo.AuthSource")
	if authSource == "" {
		err := errors.New("MongoDB authSource required to be present in config")
		log.Println(err.Error())
		return nil, err
	}

	mongoDBUsername := viper.GetString("Config.Database.Mongo.Username")
	if mongoDBUsername == "" {
		err := errors.New("MongoDB username required to be present in config")
		log.Println(err.Error())
		return nil, err
	}

	parser := JSONStructTagParser{}

	jsonStructCodec, err := bsoncodec.NewStructCodec(parser)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	customStructpbCodecs := CustomStructpbCodecs{}

	newRegistryBuilder := bson.NewRegistryBuilder()
	newRegistryBuilder.RegisterDefaultDecoder(reflect.Struct, jsonStructCodec).RegisterDefaultEncoder(reflect.Struct, jsonStructCodec)
	newRegistryBuilder.RegisterTypeDecoder(reflect.TypeOf(&structpb.Struct{}), customStructpbCodecs).RegisterTypeEncoder(reflect.TypeOf(&structpb.Struct{}), customStructpbCodecs)

	newDecoderRegistry := newRegistryBuilder.Build()

	registryOpts := options.ClientOptions{
		Registry: newDecoderRegistry,
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%v:27017", mongoDBURL)).SetAuth(
		options.Credential{
			AuthSource: authSource,
			Username:   mongoDBUsername,
			Password:   os.Getenv("MongoDBPasswd"),
		},
	), &registryOpts)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return client, nil
}

func (handler *DBUtilsHandler) parseInsertResult(insertResult *mongo.InsertOneResult, insertedModel interface{}, collection *mongo.Collection) error {
	var oid primitive.ObjectID
	var ok bool

	if oid, ok = insertResult.InsertedID.(primitive.ObjectID); !ok {
		return errors.New("Error decoding result id")
	}

	result := collection.FindOne(handler.MongoDefaultContext, bson.M{
		"_id": oid,
	})

	if result.Err() != nil {
		log.Println(result.Err().Error())
		return result.Err()
	}

	err := result.Decode(insertedModel)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
