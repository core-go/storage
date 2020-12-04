package storage

type File struct {
	ContentType string
	FileName    string
	Tags        []string
	BytesData   []byte
}

type StorageResult struct {
	Status     int    `json:"status"`
	Name       string `json:"name"`
	MediaLink  string `json:"mediaLink"`
	Link       string `json:"link"` // link file public
	BucketName string `json:"bucketName"`
}
