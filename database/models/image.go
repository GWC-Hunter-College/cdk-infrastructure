package models

type Image struct {
	ImageID string `db:"id"`
	Purpose string `db:"purpose"`
	// S3 Object Key. i.e. the path from the root of the bucket to the file.
	ObjectKey string `db:"object_key"`
	Filename  string `db:"filename"`
	Type      string `db:"mimetype"`
	CreatedAt string `db:"created_at"`
}
