package databasehandler

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/types/known/structpb"
)

//NewMongoClient Connects to a mongodb
func NewMongoClient(ctx context.Context) (*mongo.Client, error) {
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)

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
	newRegistryBuilder.RegisterTypeDecoder(
		reflect.TypeOf(&structpb.Struct{}),
		customStructpbCodecs).RegisterTypeEncoder(
		reflect.TypeOf(&structpb.Struct{}),
		customStructpbCodecs)

	newDecoderRegistry := newRegistryBuilder.Build()
	registryOpts := options.ClientOptions{
		Registry: newDecoderRegistry,
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%v:27017", mongoDBURL)).SetAuth(
		options.Credential{
			AuthSource: authSource,
			Username:   mongoDBUsername,
			Password:   os.Getenv("MONGO_INITDB_ROOT_PASSWORD"),
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

// JSONStructTagParser BSON tag parser that uses JSON tags as replacement
type JSONStructTagParser struct{}

// ParseStructTags A custom struct tag decoder function that parses the json tag instead of the bson tag
// This is done since protocol buffer will only generate json tags for its fields
// Omitempty is ignored for now TODO: should be reassesed before 1.0
func (parser JSONStructTagParser) ParseStructTags(sf reflect.StructField) (bsoncodec.StructTags, error) {
	key := strings.ToLower(sf.Name)
	tag, ok := sf.Tag.Lookup("json")
	if !ok && !strings.Contains(string(sf.Tag), ":") && len(sf.Tag) > 0 {
		tag = string(sf.Tag)
	}
	var st bsoncodec.StructTags
	if tag == "-" {
		st.Skip = true
		return st, nil
	}

	for idx, str := range strings.Split(tag, ",") {
		if idx == 0 && str != "" {
			key = str
		}
		// no omit empty
		switch str {
		case "minsize":
			st.MinSize = true
		case "truncate":
			st.Truncate = true
		case "inline":
			st.Inline = true
		}
	}

	st.Name = key

	return st, nil

}
