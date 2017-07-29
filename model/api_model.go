package model

// JSONThread TBD
type JSONThread struct {
	Posts []*JSONPost `json:"posts"`
}

// ToPosts TBD
func (jt *JSONThread) ToPosts() []*Post {
	posts := make([]*Post, 0, len(jt.Posts))
	for _, jp := range jt.Posts {
		posts = append(posts, jp.ToPost())
	}

	return posts
}

// JSONPost TBD
type JSONPost struct {
	No          int64  `json:"no"`
	FileName    string `json:"filename"`
	Tim         int64  `json:"tim"`
	Ext         string `json:"ext"`
	MD5         string `json:"md5"`
	FileDeleted int    `json:"filedeleted"`
}

// ToPost TBD
func (jp *JSONPost) ToPost() *Post {
	return &Post{
		No:          jp.No,
		Tim:         jp.Tim,
		FileName:    jp.FileName,
		Ext:         jp.Ext,
		MD5:         jp.MD5,
		FileDeleted: jp.FileDeleted > 0,
	}
}
