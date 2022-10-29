package uploads

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Request struct {
	OriginalFileName string `json:"originalFileName,omitempty" gorm:"column:original_filename" bson:"original_filename,omitempty" dynamodbav:"original_filename,omitempty" firestore:"original_filename,omitempty"`
	Filename         string `json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Type             string `json:"type,omitempty" gorm:"column:type" bson:"type,omitempty" dynamodbav:"type,omitempty" firestore:"type,omitempty"`
	Data             []byte `json:"data,omitempty" gorm:"column:data" bson:"data,omitempty" dynamodbav:"data,omitempty" firestore:"data,omitempty"`
}

type Response struct {
	Attachments Attachment `json:"attachment,omitempty" gorm:"column:attachment" bson:"attachment,omitempty" dynamodbav:"attachment,omitempty" firestore:"attachment,omitempty"`
}

type Attachment struct {
	OriginalFileName string `json:"originalFileName,omitempty" gorm:"column:original_filename" bson:"original_filename,omitempty" dynamodbav:"original_filename,omitempty" firestore:"original_filename,omitempty"`
	FileName         string `json:"fileName,omitempty" gorm:"column:fileName" bson:"fileName,omitempty" dynamodbav:"fileName,omitempty" firestore:"fileName,omitempty"`
	Url              string `json:"url,omitempty" gorm:"column:url" bson:"url,omitempty" dynamodbav:"url,omitempty" firestore:"url,omitempty" avro:"url" validate:"required"`
	Type             string `json:"type,omitempty" gorm:"column:type" bson:"type,omitempty" dynamodbav:"type,omitempty" firestore:"type,omitempty"`
}

func (u Attachment) Value() (driver.Value, error) {
	return json.Marshal(u)
}
func (u *Attachment) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &u)
}

type Entity struct {
	Id          string       `json:"id,omitempty" gorm:"column:id" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty" gorm:"column:attachments" bson:"attachments,omitempty" dynamodbav:"attachments,omitempty" firestore:"attachments,omitempty"`
}

type SingleEntity struct {
	Id          string     `json:"id,omitempty" gorm:"column:id" bson:"id,omitempty" dynamodbav:"id,omitempty" firestore:"id,omitempty"`
	Attachments Attachment `json:"attachments,omitempty" gorm:"column:attachments" bson:"attachments,omitempty" dynamodbav:"attachments,omitempty" firestore:"attachments,omitempty"`
}
