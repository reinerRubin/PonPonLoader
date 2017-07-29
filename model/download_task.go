package model

import (
	"fmt"
	"net/url"
	"path"

	"github.com/PonPonLoader/definition"
)

// DownloadTask TBD
type DownloadTask struct {
	Source *url.URL
	Target string
}

// NewDownloadTask TBD
func NewDownloadTask(post *Post, baseTargetPath string) (*DownloadTask, error) {
	if !post.HasImage() || len(post.BoardName) == 0 {
		return nil, fmt.Errorf("post number (%d) does not have pic", post.No)
	}

	imageServerName, err := post.ImagePath()
	if err != nil {
		return nil, err
	}

	imageURLString := fmt.Sprintf("%s/%s/%s", definition.ImageHost, post.BoardName, imageServerName)
	imageURL, err := url.Parse(imageURLString)
	if err != nil {
		return nil, err
	}

	pathForSave, err := post.ImageName()
	if err != nil {
		return nil, err
	}

	return &DownloadTask{
		Source: imageURL,
		Target: path.Join(baseTargetPath, pathForSave),
	}, nil
}
