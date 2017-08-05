package model

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Thread TBD
type Thread struct {
	No        int64
	BoardName string
}

// NewThreadFromURL TBD
func NewThreadFromURL(threadURLString string) (*Thread, error) {
	threadURL, err := url.Parse(threadURLString)
	if err != nil {
		return nil, err
	}

	paths := strings.Split(threadURL.Path, "/")
	// /c/thread/2942063/madotsuki-thread'
	if len(paths) < 3 {
		return nil, fmt.Errorf("cant parse URL: %s", threadURLString)
	}
	boardName := paths[1]
	boardNoStr := paths[3]

	boardNo, err := strconv.ParseInt(boardNoStr, 10, 0)
	if err != nil {
		return nil, err
	}

	return &Thread{
		No:        boardNo,
		BoardName: boardName,
	}, nil

}

// URLPath TBD
func (th *Thread) URLPath() string {
	return fmt.Sprintf("%s/thread/%d.json", th.BoardName, th.No)
}

// Post TBD
type Post struct {
	No          int64
	BoardName   string
	FileName    string
	Tim         int64
	Ext         string
	MD5         string
	FileDeleted bool
}

// HasImage TBD
func (p *Post) HasImage() bool {
	return p.Tim > 0 && !p.FileDeleted
}

// ImagePath TBD
func (p *Post) ImagePath() (string, error) {
	if !p.HasImage() {
		return "", fmt.Errorf("post number (%d) does not have pic", p.No)
	}
	return fmt.Sprintf("%d%s", p.Tim, p.Ext), nil

}

// ImageName TBD
func (p *Post) ImageName() (string, error) {
	if !p.HasImage() {
		return "", fmt.Errorf("post number (%d) does not have pic", p.No)
	}

	return fmt.Sprintf("%s-%d%s", p.FileName, p.Tim, p.Ext), nil
}
