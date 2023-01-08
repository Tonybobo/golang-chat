package utils

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/tonybobo/go-chat/config"
	"github.com/tonybobo/go-chat/pkg/global/log"
	"google.golang.org/api/option"
)

type StorageClient struct {
	storageClient *storage.Client
	bucket        string
	projectID     string
}

var Uploader *StorageClient

func init() {
	client, err := storage.NewClient(context.Background(), option.WithCredentialsFile("keys.json"))
	if err != nil {
		log.Logger.Error("GCP init error", log.Any("error :", err))
		panic(err)
	}

	Uploader = &StorageClient{
		storageClient: client,
		bucket:        config.GetConfig().Bucket,
		projectID:     config.GetConfig().ProjectID,
	}
}

func (s *StorageClient) UploadImage(file multipart.File, path string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	writer := s.storageClient.Bucket(s.bucket).Object(path).NewWriter(ctx)

	if _, err := io.Copy(writer, file); err != nil {
		log.Logger.Error("Upload Error", log.Any("error", err))
		return "", err
	}

	if err := writer.Close(); err != nil {
		log.Logger.Error("Upload Error", log.Any("error", err))
		return "", err
	}

	url, err := url.Parse("/" + s.bucket + "/" + writer.Attrs().Name)

	if err != nil {
		log.Logger.Error("Upload Error", log.Any("error", err))
		return "", err
	}
	return url.EscapedPath(), nil

}

func (s *StorageClient) DeleteImage(path string) error {
	object := strings.Split(path, "go-chat/")[1]
	fmt.Println(object)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	o := s.storageClient.Bucket(s.bucket).Object(object)

	attrs, err := o.Attrs(ctx)

	if err != nil {
		log.Logger.Error("error", log.Any("error", err))
		return err
	}

	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err := o.Delete(ctx); err != nil {
		log.Logger.Error("Delete image error", log.Any("error", err))
		return err
	}

	return nil

}
