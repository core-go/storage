package storage

type Config struct {
	CredentialsFile       string `mapstructure:"credentials_File"`
	BucketName            string `mapstructure:"bucket_name"`
	SubDirectory          string `mapstructure:"sub_directory"`
	PermissionFileRoleAll bool `mapstructure:"permission_file_role_all"`
}
