package databasehandler

import (
	"encoding/json"
	"reflect"

	log "github.com/sirupsen/logrus"

	structpb "github.com/golang/protobuf/ptypes/struct"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

// CustomStructpbCodecs Struct to hold the custom structpb encoder and decoder for bson
type CustomStructpbCodecs struct {
}

// EncodeValue Encodes a structpb value into the correct bson representation
// Encoding a structpb message type value directly into bson using the default bson encoder will represent the internal
// structpb structure that represents the data types along with their corresponding data
// To have the correct representation an intermediate json representation is used using the internal structp tooling
func (codecs CustomStructpbCodecs) EncodeValue(ctx bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	data := val.Interface().(*structpb.Struct)

	jsonData, err := data.MarshalJSON()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	var encodedBsonStructData bson.M
	err = json.Unmarshal(jsonData, &encodedBsonStructData)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	newEncoder, err := bson.NewEncoder(vw)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = newEncoder.Encode(encodedBsonStructData)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// DecodeValue Custom decoder for protocol buffers structpb message type
// The data is read from a new decoder into a bson.M structure.
// The bson.M struct is marshalled into JSON bytes and unmarshalled into its structpb representation, which is then set as value
func (codecs CustomStructpbCodecs) DecodeValue(r bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	decoder, err := bson.NewDecoder(vr)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	intermediateData := bson.M{}
	err = decoder.Decode(&intermediateData)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	intermediateByteData, err := json.Marshal(&intermediateData)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	var structData structpb.Struct
	err = structData.UnmarshalJSON(intermediateByteData)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	val.Set(reflect.ValueOf(&structData))

	return nil
}
