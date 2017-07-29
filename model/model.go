package model

import "fmt"

// Thread TBD
type Thread struct {
	No        int64
	BoardName string
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
