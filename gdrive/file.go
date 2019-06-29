package gdrive

import (
	"fmt"
	"io"
	"io/ioutil"

	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

const DirectoryMimeType = "application/vnd.google-apps.folder"
const MaxPageSize = 1000

// ListFile return all file under the specific folder.
func (gdClient *GoogleDriveClient) ListFile(parentFolder *drive.File, maxFiles int64) ([]*drive.File, error) {
	var pageSize int64

	if maxFiles <= 0 {
		return nil, fmt.Errorf("maxFiles is invalid")
	}

	if maxFiles > 0 && maxFiles < 1000 {
		pageSize = maxFiles
	} else {
		pageSize = MaxPageSize
	}

	return gdClient.listFile(parentFolder, maxFiles, pageSize)
}

func (gdClient *GoogleDriveClient) listFile(parentFolder *drive.File, maxFiles int64, pageSize int64) ([]*drive.File, error) {
	var queryStr string
	var err error
	var files []*drive.File

	if parentFolder != nil {
		queryStr = fmt.Sprintf("'%s' in parents", parentFolder.Id)
	}

	controlledStop := fmt.Errorf("Controlled stop")
	apiFields := []googleapi.Field{"nextPageToken", "files(id,name,md5Checksum,mimeType,size,createdTime,parents)"}

	err = gdClient.driveClient.Files.List().Q(queryStr).Fields(apiFields...).PageSize(pageSize).Pages(context.TODO(), func(fl *drive.FileList) error {

		files = append(files, fl.Files...)

		// Stop when we have all the files we need
		if maxFiles > 0 && len(files) >= int(maxFiles) {
			return controlledStop
		}

		return nil
	})

	if err != nil && err != controlledStop {
		return nil, err
	}

	if len(files) <= 0 {
		return nil, nil
	}

	// TODO: check pagetoken if it reached the end.
	return files, nil
}

// ListAllFile list all file under parent folder. Be aware that the files are staying in Trash is also counted into it if you delete these files before.
func (gdClient *GoogleDriveClient) ListAllFile(parentFolder *drive.File) ([]*drive.File, error) {
	return gdClient.listAllFile(parentFolder)
}

func (gdClient *GoogleDriveClient) listAllFile(parentFolder *drive.File) ([]*drive.File, error) {

	var queryStr string
	var err error
	var files []*drive.File

	if parentFolder == nil {
		return nil, fmt.Errorf("parentFolder can't be nil")
	}

	queryStr = fmt.Sprintf("'%s' in parents", parentFolder.Id)
	apiFields := []googleapi.Field{"nextPageToken", "files(id,name,md5Checksum,mimeType,size,createdTime,parents)"}

	err = gdClient.driveClient.Files.List().Q(queryStr).Fields(apiFields...).PageSize(MaxPageSize).Pages(context.TODO(), func(fl *drive.FileList) error {

		files = append(files, fl.Files...)

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(files) <= 0 {
		return nil, nil
	}

	// TODO: check pagetoken if it reached the end.
	return files, nil

	return nil, nil
}
func (gdClient *GoogleDriveClient) getFile(driveFiles []*drive.File, fileName string) *drive.File {
	for i := 0; i < len(driveFiles); i++ {
		driveFile := driveFiles[i]
		if driveFile.Name == fileName {
			return driveFile
		}
	}

	return nil
}

// GetFileByName find the file information according to the name.
func (gdClient *GoogleDriveClient) GetFileByName(fileName string, parentFolder *drive.File) (*drive.File, error) {
	return gdClient.getFileByName(fileName, parentFolder)
}

func (gdClient *GoogleDriveClient) getFileByName(fileName string, parentFolder *drive.File) (*drive.File, error) {
	var queryStr string
	if parentFolder == nil {
		queryStr = fmt.Sprintf("name='%s'", fileName)
	} else if parentFolder != nil {
		queryStr = fmt.Sprintf("name='%s' and '%s' in parents", fileName, parentFolder.Id)
	}

	res, err := gdClient.driveClient.Files.List().PageSize(10).Q(queryStr).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		return nil, err
	}

	// TODO: Should I handle this?
	// if res.IncompleteSearch == true {
	// 	return nil,
	// }

	if len(res.Files) <= 0 {
		return nil, nil
	}

	return res.Files[0], nil
}

func (gdClient *GoogleDriveClient) downloadFileById(id string) ([]byte, error) {

	if len(id) == 0 {
		return nil, fmt.Errorf("id can't be empty")
	}

	res, err := gdClient.driveClient.Files.Get(id).Download()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil

}

func (gdClient *GoogleDriveClient) Mkdir(name string, description string, parentFolder *drive.File) (*drive.File, error) {

	if len(name) == 0 {
		return nil, fmt.Errorf("folder name is invalid")
	}

	return gdClient.mkdir(name, description, parentFolder)
}

func (gdClient *GoogleDriveClient) mkdir(name string, description string, parentFolder *drive.File) (*drive.File, error) {
	dstFile := &drive.File{
		Name:        name,
		Description: description,
		MimeType:    DirectoryMimeType,
	}

	if parentFolder != nil {
		// Set parent folders
		dstFile.Parents = []string{parentFolder.Id}
	}

	// Create directory
	folder, err := gdClient.driveClient.Files.Create(dstFile).Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to create directory: %s", err)
	}

	return folder, nil

}

func (gdClient *GoogleDriveClient) DeleteFileById(id string) error {
	if len(id) == 0 {
		return fmt.Errorf(" file id can't be empty")
	}

	return gdClient.deleteFileById(id)
}

func (gdClient *GoogleDriveClient) deleteFileById(id string) error {
	err := gdClient.driveClient.Files.Delete(id).Do()
	if err != nil {
		return fmt.Errorf("Failed to delete file: %s", err)
	}
	return nil
}

func (gdClient *GoogleDriveClient) UploadFile(fileName string, description string, content io.Reader, parentFolder *drive.File) (*drive.File, error) {
	if len(fileName) == 0 {
		return nil, fmt.Errorf("file name can't be empty.")
	}

	return gdClient.uploadFile(fileName, description, content, parentFolder)
}

func (gdClient *GoogleDriveClient) uploadFile(fileName string, description string, content io.Reader, parentFolder *drive.File) (*drive.File, error) {
	dstFile := &drive.File{
		Name:        fileName,
		Description: description,
	}

	if parentFolder != nil {
		dstFile.Parents = []string{parentFolder.Id}
	}

	file, err := gdClient.driveClient.Files.Create(dstFile).Fields("id", "name", "size", "webContentLink").Media(content).Do()
	if err != nil {
		return nil, err
	}

	return file, nil
}
