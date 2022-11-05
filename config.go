package storage

type Config struct {
	Bucket                      string `yaml:"bucket" mapstructure:"bucket" gorm:"column:bucket" bson:"bucket,omitempty" dynamodbav:"bucket,omitempty" firestore:"bucket,omitempty"`
	Directory                   string `yaml:"directory" mapstructure:"directory" json:"directory" gorm:"column:directory" bson:"directory,omitempty" dynamodbav:"directory,omitempty" firestore:"directory,omitempty"`
	Public                      *bool  `yaml:"public" mapstructure:"public" json:"allUsersReader" gorm:"column:public" bson:"public,omitempty" dynamodbav:"public,omitempty" firestore:"public,omitempty"`
	Private                     *bool  `yaml:"private" mapstructure:"private" json:"allUsersReader" gorm:"column:private" bson:"private,omitempty" dynamodbav:"private,omitempty" firestore:"private,omitempty"`
	AllAuthenticatedUsersWriter *bool  `yaml:"all_authenticated_users_writer" mapstructure:"all_authenticated_users_writer" json:"allAuthenticatedUsersWriter" gorm:"column:allauthenticateduserswriter" bson:"allAuthenticatedUsersWriter,omitempty" dynamodbav:"allAuthenticatedUsersWriter,omitempty" firestore:"allAuthenticatedUsersWriter,omitempty"`
}
