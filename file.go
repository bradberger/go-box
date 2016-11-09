package box

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"golang.org/x/net/context"
)

// UploadFile uploads the given fileContents to fileName in the given folder. If the parent folder is nil, then
// the base folder for the box account will be selected.
func (a *API) UploadFile(ctx context.Context, f *FolderDetails, fileName string, fileContents []byte) (*PathCollection, error) {

	folderID := "0"
	if f != nil {
		folderID = f.ID
	}
	fileName = filepath.Base(fileName)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(part, bytes.NewBuffer(fileContents)); err != nil {
		return nil, err
	}

	if err = writer.WriteField("attributes", fmt.Sprintf(`{"name":"%s", "parent":{"id":"%s"}}`, fileName, folderID)); err != nil {
		return nil, err
	}

	if err = writer.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://upload.box.com/api/2.0/files/content", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	var pc PathCollection
	if _, err := a.do(ctx, req, &pc); err != nil {
		return nil, err
	}

	return &pc, nil
}
