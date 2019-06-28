package gdrive

import (
	"bytes"
	"fmt"
	"strconv"

	// "fmt"
	"testing"

	"google.golang.org/api/drive/v3"
)

const TEST_FOLDER = "GoGoogleDriveUnitTest"

var gdrive *GoogleDriveClient
var testFolder *drive.File

func TestListFile(tt *testing.T) {
	var err error

	projectRoot := getProjectRoot()

	clientSecret := projectRoot + "/config/client_secret.json"
	tokenPath := projectRoot + "/config/token.json"

	gdrive, err = NewGoogleDriveClient(clientSecret, tokenPath)
	if err != nil {
		tt.Error(err)
		return
	}

	files, err := gdrive.ListFile(nil, 2)
	if err != nil {
		tt.Error(err)
		return
	}

	if len(files) != 2 {
		tt.Log("The count of files is " + strconv.Itoa(len(files)) + " instead of 10")
		tt.Fail()
		return
	}

}

func TestFileByName(tt *testing.T) {
	var err error

	testFolder, err = gdrive.GetFileByName(TEST_FOLDER, nil)
	if err != nil {
		tt.Error(err)
		return
	}

	if testFolder.Name != TEST_FOLDER {
		tt.Log("File name is not StockTest")
		tt.Fail()
		return
	}

}

func TestDownloadFile(tt *testing.T) {

	files, err := gdrive.ListFile(testFolder, 3)
	if err != nil {
		tt.Error(err)
		return
	}

	if len(files) != 3 {
		tt.Log("The count of files is " + strconv.Itoa(len(files)) + " instead of 3")
		tt.Fail()
		return
	}

	_, err = gdrive.downloadFileById(files[1].Id)
	if err != nil {
		tt.Log(err)
		tt.Fail()
		return
	}

}

func TestMkdirAndDeleteIt(tt *testing.T) {
	folderName := "CreatTestFolder"

	folder, err := gdrive.Mkdir(folderName, "This is part of unit test", testFolder)
	if err != nil {
		tt.Log(err)
		tt.Fail()
		return
	}

	if folder.Name != folderName {
		tt.Log("created folder is not correct.")
		tt.Fail()
		return
	}

	err = gdrive.DeleteFileById(folder.Id)
	if err != nil {
		tt.Log(err)
		tt.Fail()
		return
	}
}

func TestUploadFileAndDeleteIt(tt *testing.T) {
	fileName := "TestUpload"

	testBytes := []byte("this is TestUploadFileAndDeleteIt content")

	content := bytes.NewReader(testBytes)

	file, err := gdrive.UploadFile(fileName, "Unit test upload file", content, testFolder)
	if err != nil {
		tt.Log(err)
		tt.Fail()
		return
	}

	if fileName != file.Name || len(file.Id) == 0 {
		tt.Log("Upload file fails. Something wrong with name or id")
		tt.Fail()
		return
	}

	err = gdrive.DeleteFileById(file.Id)
	if err != nil {
		tt.Log(err)
		tt.Fail()
		return
	}

}

func TestListAllFile(tt *testing.T) {
	var err error

	files, err := gdrive.ListAllFile(testFolder)
	if err != nil {
		tt.Error(err)
		return
	}

	if len(files) != 10 {
		tt.Log("The count of files is " + strconv.Itoa(len(files)) + " instead of 10")
		tt.Fail()
		return
	}

	for i := 0; i < len(files); i++ {
		fmt.Println(files[i].Name)
	}

}
