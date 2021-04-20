package storage

type Config struct {
	BucketName                  string `mapstructure:"bucket_name" gorm:"column:bucketname" bson:"bucketName,omitempty" dynamodbav:"bucketName,omitempty" firestore:"bucketName,omitempty"`
	SubDirectory                string `mapstructure:"sub_directory" json:"subDirectory" gorm:"column:subdirectory" bson:"subDirectory,omitempty" dynamodbav:"subDirectory,omitempty" firestore:"subDirectory,omitempty"`
	AllUsersAreReader           *bool  `mapstructure:"all_users_reader" json:"allUsersReader" gorm:"column:allusersreader" bson:"allUsersReader,omitempty" dynamodbav:"allUsersReader,omitempty" firestore:"allUsersReader,omitempty"`
	AllAuthenticatedUsersReader *bool  `mapstructure:"all_authenticated_users_reader" json:"allAuthenticatedUsersReader" gorm:"column:allauthenticatedusersreader" bson:"allAuthenticatedUsersReader,omitempty" dynamodbav:"allAuthenticatedUsersReader,omitempty" firestore:"allAuthenticatedUsersReader,omitempty"`
	AllAuthenticatedUsersWriter *bool  `mapstructure:"all_authenticated_users_writer" json:"allAuthenticatedUsersWriter" gorm:"column:allauthenticateduserswriter" bson:"allAuthenticatedUsersWriter,omitempty" dynamodbav:"allAuthenticatedUsersWriter,omitempty" firestore:"allAuthenticatedUsersWriter,omitempty"`
}
