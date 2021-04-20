package storage

type File struct {
	ContentType string   `mapstructure:"content_type" json:"contentType" gorm:"column:contenttype" bson:"contentType,omitempty" dynamodbav:"contentType,omitempty" firestore:"contentType,omitempty"`
	Name        string   `mapstructure:"name" json:"permissionFileRoleAll" gorm:"column:permissionFileRoleAll" bson:"permissionFileRoleAll,omitempty" dynamodbav:"permissionFileRoleAll,omitempty" firestore:"permissionFileRoleAll,omitempty"`
	Tags        []string `mapstructure:"tags" json:"tags" gorm:"column:tags" bson:"tags,omitempty" dynamodbav:"tags,omitempty" firestore:"tags,omitempty"`
	Bytes       []byte   `mapstructure:"bytes" json:"bytes" gorm:"column:bytes" bson:"bytes,omitempty" dynamodbav:"bytes,omitempty" firestore:"bytes,omitempty"`
}

type StorageResult struct {
	Status     int64  `mapstructure:"status" json:"status" gorm:"column:status" bson:"status,omitempty" dynamodbav:"status,omitempty" firestore:"status,omitempty"`
	Name       string `mapstructure:"name" json:"name" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	MediaLink  string `mapstructure:"media_link" json:"mediaLink" gorm:"column:medialink" bson:"mediaLink,omitempty" dynamodbav:"mediaLink,omitempty" firestore:"mediaLink,omitempty"`
	Link       string `mapstructure:"status" json:"link" gorm:"column:link" bson:"link,omitempty" dynamodbav:"link,omitempty" firestore:"link,omitempty"` // link file public
	BucketName string `mapstructure:"bucket_name" json:"bucketName" gorm:"column:bucketname" bson:"bucketName,omitempty" dynamodbav:"bucketName,omitempty" firestore:"bucketName,omitempty"`
}
