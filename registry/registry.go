package registry

import (
	"content-management/core/config"
	"content-management/pkg/database"
	"content-management/thirdparty/minio"
)

// Registry khởi tạo các client, config mà ứng dụng sử dụng...
// Dùng các client, config đã được khởi ở trên để khởi tạo các service Aggregate và Query của ứng dụng...
type Registry struct {
	Config      config.Config
	DB          *database.Database
	MinioClient *minio.Client
}

// New ...
func New(c config.Config) (*Registry, error) {
	r := &Registry{
		Config: c,
		DB:     database.New(c.Databases.PostgresConfig),
		MinioClient: minio.New(&minio.Config{
			Endpoint:        c.Minio.Endpoint,
			SecretAccessKey: c.Minio.SecretAccessKey,
			AccessKey:       c.Minio.AccessKey,
			BucketName:      c.Minio.BucketName,
		}),
	}
	return r, nil
}
