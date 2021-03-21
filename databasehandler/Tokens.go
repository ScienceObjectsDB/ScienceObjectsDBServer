package databasehandler

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ScienceObjectsDB/go-api/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const tokenLen = 64

//TokenActionHandler Handler for token related database actions
type TokenActionHandler struct {
	*DBUtilsHandler
}

// CreateToken Creates a token for a project
func (handler *TokenActionHandler) CreateToken(userID string, rights []models.Right, resource models.Resource) (*models.TokenEntry, error) {
	b := make([]byte, tokenLen)
	_, err := rand.Read(b)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	uuidString := uuid.New().String()
	secretString := base64.RawStdEncoding.EncodeToString(b)

	user := models.User{
		Resource: resource,
		Rights:   rights,
		UserID:   userID,
	}

	expireTime := timestamppb.New(time.Now().Add(365 * 24 * time.Hour))

	token := models.TokenEntry{
		ID:       uuidString,
		Created:  timestamppb.Now(),
		Token:    secretString,
		UserID:   &user,
		Expires:  expireTime,
		Resource: resource,
	}

	insertResult, err := handler.GetTokenCollection().InsertOne(handler.MongoDefaultContext, &token)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	insertedToken := models.TokenEntry{}

	err = handler.parseInsertResult(insertResult, &insertedToken, handler.GetTokenCollection())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &insertedToken, nil
}

//ValidateTokenForResourceAction Validates an action on a specific resource from a token
func (handler *TokenActionHandler) ValidateTokenForResourceAction(
	token string, resourceID string,
	ressourceType models.Resource,
	requiredRight models.Right) (bool, error) {

	authorized := false

	ressourceType = models.Resource_Project

	result := handler.GetTokenCollection().FindOne(handler.MongoDefaultContext, bson.M{
		"Token":         token,
		"UserID.Rights": requiredRight,
		"Resource":      ressourceType,
	})

	if result.Err() != nil && result.Err() != mongo.ErrNoDocuments {
		log.Println(result.Err().Error())
		return false, result.Err()
	}

	if result.Err() == mongo.ErrNoDocuments {
		return false, nil
	}

	authorized = true

	return authorized, nil
}

// GetTokenUser Returns the user of this token
func (handler *TokenActionHandler) GetTokenUser(accessToken string) (*models.TokenEntry, error) {
	queryResults := handler.GetTokenCollection().FindOne(handler.MongoDefaultContext, bson.M{
		"Token": accessToken,
	})

	if queryResults.Err() != nil {
		log.Println(queryResults.Err().Error())
		return nil, queryResults.Err()
	}

	var token models.TokenEntry

	err := queryResults.Decode(&token)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &token, err
}

// GetUserDatasetTokens Returns all tokens of a user for a dataset
func (handler *TokenActionHandler) GetUserDatasetTokens(userID string) ([]*models.TokenEntry, error) {
	csr, err := handler.GetTokenCollection().Find(handler.MongoDefaultContext, bson.M{
		"UserID.UserID": userID,
	})

	if err != nil && err != mongo.ErrNoDocuments {
		log.Println(err.Error())
		return nil, err
	}

	var tokens []*models.TokenEntry

	for csr.Next(handler.MongoDefaultContext) {
		token := models.TokenEntry{}

		err := csr.Decode(&token)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		tokens = append(tokens, &token)
	}

	return tokens, nil
}
