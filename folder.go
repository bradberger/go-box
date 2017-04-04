package box

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"
)

// Folder access levels as defined in the Box API documentation
const (
	AccessLevelOpen          = "open"
	AccessLevelCompany       = "company"
	AccessLevelCollaborators = "collaborators"
	AccessLevelDefault       = ""
)

type CreateFolderRequest struct {
	Name         string       `json:"name"`
	ParentFolder ParentFolder `json:"parent"`
}

type ParentFolder struct {
	ID string `json:"id"`
}

type Folder struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	SequenceID string `json:"sequence_id"`
	ETag       string `json:"etag"`
	Name       string `json:"name"`
}

type FolderDetails struct {
	Folder
	SharedLink        SharedLink        `json:"shared_link"`
	CreatedAt         *time.Time        `json:"created_at,omitempty"`
	ModifiedAt        *time.Time        `json:"modified_at,omitempty"`
	Size              int64             `json:"size,omitempty"`
	FolderUploadEmail FolderUploadEmail `json:"folder_upload_email,omitempty"`
	Parent            Folder
	ItemStatus        string         `json:"item_status,omitempty"`
	ItemCollection    ItemCollection `json:"item_collection,omitempty"`
}

type PathCollection struct {
	TotalCount int64    `json:"total_count"`
	Entries    []Folder `json:"entries"`
}

type ItemCollection struct {
	PathCollection
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}

type FolderUploadEmail struct {
	Access string `json:"acceess"`
	Email  string `json:"email"`
}

type SharedLink struct {
	URL               string      `json:"url,omitempty"`
	DownloadURL       string      `json:"download_url,omitempty"`
	VanityURL         string      `json:"vanity_url,omitempty"`
	IsPasswordEnabled bool        `json:"is_password_enabled,omitempty"`
	UnsharedAt        *time.Time  `json:"unshared_at,omitempty"`
	DownloadCount     int64       `json:"download_count,omitempty"`
	PreviewCount      int64       `json:"preview_count,omitempty"`
	Access            string      `json:"access,omitempty"`
	Permissions       Permissions `json:"permissions,omitempty"`
}

type SharedLinkCreateRequest struct {
	SharedLinkPermissions SharedLinkPermissions `json:"shared_link"`
}

type SharedLinkPermissions struct {
	Access          string       `json:"access,omitempty"`
	UnsharedAt      *time.Time   `json:"unshared_at,omitempty"`
	Password        string       `json:"password,omitempty"`
	EffectiveAccess string       `json:"effective_access,omitempty"`
	Permissions     *Permissions `json:"permissions,omitempty"`
}

type Permissions struct {
	CanDownload bool `json:"can_download"`
	CanPreview  bool `json:"can_preview"`
}

func (a *API) CreateFolder(ctx context.Context, folderName string, parentID string) (*FolderDetails, *ErrorCodeResponse, error) {
	var f FolderDetails
	fr := &CreateFolderRequest{Name: folderName, ParentFolder: ParentFolder{ID: parentID}}
	if resp, err := a.PostJSON(ctx, "/folders", fr, &f); err != nil {
		log.Printf("Error: %+v", resp)
		return nil, resp, err
	}
	return &f, nil, nil
}

func (a *API) GetFolderDetails(ctx context.Context, folderID string) (*FolderDetails, *ErrorCodeResponse, error) {
	var f FolderDetails
	if resp, err := a.Get(ctx, fmt.Sprintf("/folders/%s", folderID), &f); err != nil {
		return nil, resp, err
	}
	return &f, nil, nil
}

func (a *API) CreateSharedLink(ctx context.Context, folderID string, l *SharedLinkCreateRequest) (*FolderDetails, *ErrorCodeResponse, error) {
	var f FolderDetails
	if l == nil {
		l = &SharedLinkCreateRequest{}
	}
	if resp, err := a.Put(ctx, fmt.Sprintf("/folders/%s", folderID), &l, &f); err != nil {
		return nil, resp, err
	}
	return &f, nil, nil
}
