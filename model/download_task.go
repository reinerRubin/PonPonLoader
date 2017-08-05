package model

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/PonPonLoader/definition"
)

// DownloadTask TBD
type DownloadTask struct {
	Source *url.URL
	Target string
	MD5    string
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
		MD5:    post.MD5,
	}, nil
}

// Run TBD
func (t *DownloadTask) Run() error {
	file, err := os.Create(t.Target)
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := http.Get(t.Source.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}

	return nil
}
