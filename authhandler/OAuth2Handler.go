package authhandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

//OAuth2Handler Handles oauth2
type OAuth2Handler struct {
	UserInfoEndpointURL string
}

// Init Initializes the auth handler object
func Init() (*OAuth2Handler, error) {
	endpointURL := viper.GetString("Auth.UserInfoEndpointURL")
	if endpointURL == "" {
		err := errors.New("Endpoint URL has to be provided in config as 'Auth.UserInfoEndpointURL'")
		log.Println(err.Error())
		return nil, err
	}

	oauth2Handler := OAuth2Handler{
		UserInfoEndpointURL: endpointURL,
	}

	return &oauth2Handler, nil
}

func (handler *OAuth2Handler) getUserIDFromOAuth2(accessToken string) (string, error) {
	req, err := http.NewRequest(
		"GET",
		handler.UserInfoEndpointURL,
		http.NoBody,
	)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		err := fmt.Errorf("bad reponse when requesting userinfo: %v", response.Status)
		log.Println(err)
		return "", err
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed reading response body: %s", err.Error())
	}

	parsedContents := make(map[string]interface{})
	err = json.Unmarshal(contents, &parsedContents)
	if err != nil {
		log.Println(err.Error()) // Lists all datasets
		return "", err
	}

	var ok bool
	var userID interface{}
	if userID, ok = parsedContents["sub"]; !ok {
		return "", fmt.Errorf("Could not read sub claim from userinfo response")
	}

	userIDString := userID.(string)
	return userIDString, nil
}
