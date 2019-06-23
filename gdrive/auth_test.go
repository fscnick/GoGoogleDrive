package gdrive

import "testing"
import "runtime"
import "path/filepath"

// import "fmt"

func getProjectRoot() string {

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	basepath = filepath.Dir(basepath)
	return basepath
}

func TestConnect2Gdrive(tt *testing.T) {

	projectRoot := getProjectRoot()

	clientSecret := projectRoot + "/config/client_secret.json"
	tokenPath := projectRoot + "/config/token.json"

	_, err := NewGoogleDriveClient(clientSecret, tokenPath)
	if err != nil {
		tt.Error(err)
		return
	}

}

func TestNewClient(tt *testing.T) {

	projectRoot := getProjectRoot()

	clientSecret := projectRoot + "/config/client_secret.json"
	tokenPath := projectRoot + "/config/token.json"

	_, err := NewClient(clientSecret, tokenPath)
	if err != nil {
		tt.Error(err)
		return
	}

}
