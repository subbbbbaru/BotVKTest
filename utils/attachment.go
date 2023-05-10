package utils

import (
	"encoding/json"
	"fmt"
)

type Attachments interface {
	ToAttachment() string
}

type Attachment struct {
	Response Response `json:"response"`
}

type Response struct {
	Type string `json:"type"`
	Doc  Doc    `json:"doc"`
}

type Doc struct {
	ID       int64  `json:"id"`
	OwnerID  int64  `json:"owner_id"`
	Title    string `json:"title"`
	Size     int64  `json:"size"`
	EXT      string `json:"ext"`
	Date     int64  `json:"date"`
	Type     int64  `json:"type"`
	URL      string `json:"url"`
	IsUnsafe int64  `json:"is_unsafe"`
}

func UnmarshalWelcome(data []byte) (Attachment, error) {
	var r Attachment
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Attachment) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (doc Doc) ToAttachment() string {
	return fmt.Sprintf("doc%d_%d", doc.OwnerID, doc.ID)
}
