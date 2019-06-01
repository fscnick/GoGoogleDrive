package gdrive

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

// import "google.golang.org/api/drive/v3"

type GoogleDriveClient struct {
	isDriveInitialized bool
	driveClient        *drive.Service
}

// NewGoogleDriveClient create a new google dirve client to handle files uploading.
func NewGoogleDriveClient(clientSecretPath string, tokenFilePath string) (*GoogleDriveClient, error) {
	gdClient := new(GoogleDriveClient)

	if len(tokenFilePath) < 5 || tokenFilePath[len(tokenFilePath)-5:] != ".json" {
		return nil, fmt.Errorf("tokenFilePath is not end with json")
	}

	bytes, err := ioutil.ReadFile(clientSecretPath)
	if err != nil {
		return nil, fmt.Errorf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(bytes, drive.DriveScope)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secret file to config: %v", err)
	}

	client, err := getClient(config, tokenFilePath)
	if err != nil {
		return nil, fmt.Errorf("Unable to get client from config file: %v", err)
	}

	gdClient.driveClient, err = drive.New(client)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Drive client: %v", err)
	}

	gdClient.isDriveInitialized = true
	return gdClient, nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, tokenFilePath string) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := tokenFilePath
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}

		err = saveToken(tokFile, tok)
		if err != nil {
			log.Println(err)
		}
	}
	return config.Client(context.Background(), tok), nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve token from web %v", err)
	}
	return tok, nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)

	return nil
}
