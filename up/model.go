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
	Size             int64  `json:"size,omitempty" gorm:"column:size" bson:"size,omitempty" dynamodbav:"size,omitempty" firestore:"size,omitempty"`
	Data             []byte `json:"data,omitempty" gorm:"column:data" bson:"data,omitempty" dynamodbav:"data,omitempty" firestore:"data,omitempty"`
}

type Upload struct {
	OriginalFileName string `json:"originalFileName,omitempty" gorm:"column:original_filename" bson:"original_filename,omitempty" dynamodbav:"original_filename,omitempty" firestore:"original_filename,omitempty"`
	FileName         string `json:"fileName,omitempty" gorm:"column:fileName" bson:"fileName,omitempty" dynamodbav:"fileName,omitempty" firestore:"fileName,omitempty"`
	Url              string `json:"url,omitempty" gorm:"column:url" bson:"url,omitempty" dynamodbav:"url,omitempty" firestore:"url,omitempty" avro:"url" validate:"required"`
	Type             string `json:"type,omitempty" gorm:"column:type" bson:"type,omitempty" dynamodbav:"type,omitempty" firestore:"type,omitempty"`
	Size             int64  `json:"size,omitempty" gorm:"column:size" bson:"size,omitempty" dynamodbav:"size,omitempty" firestore:"size,omitempty"`
}

func (u Upload) Value() (driver.Value, error) {
	return json.Marshal(u)
}
func (u *Upload) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &u)
}

type FileConfig struct {
	MaxSize           int64  `yaml:"max_size" mapstructure:"max_size" json:"maxSize,omitempty"`
	MaxSizeMemory     int64  `yaml:"max_size_memory" mapstructure:"max_size_memory" json:"maxSizeMemory,omitempty"`
	AllowedExtensions string `yaml:"max_size_memory" mapstructure:"allowed_extensions" json:"allowedExtensions,omitempty"`
}
