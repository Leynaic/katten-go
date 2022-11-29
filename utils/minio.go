package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"

	"github.com/gosimple/slug"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioInstance *minio.Client
var bucketName = "cats"

func NewMinioClient(endpoint string, accessKeyID string, secretAccessKey string) {
	var err error

	useSSL := true

	minioInstance, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatalln(err)
	}

}

func GetMinioInstance() *minio.Client {
	return minioInstance
}

func Copy(baseObject string) (filename string, err error) {
	objectName := fmt.Sprintf("%d-%s", time.Now().Unix(), baseObject)

	srcOpts := minio.CopySrcOptions{
		Bucket: bucketName,
		Object: baseObject,
	}

	// Destination object
	dstOpts := minio.CopyDestOptions{
		Bucket: bucketName,
		Object: objectName,
	}

	if uploadInfo, err := GetMinioInstance().CopyObject(context.Background(), dstOpts, srcOpts); err != nil {
		fmt.Println(err)
		return "", err
	} else {
		fmt.Println(uploadInfo)
	}

	return objectName, err
}

func GetUrl(objectName string) (u *url.URL, err error) {
	if cacheUrl := Cache.Get(objectName); cacheUrl == nil {

		if objectUrl, err := GetMinioInstance().PresignedGetObject(context.Background(), bucketName, objectName, time.Second*24*60*60, make(url.Values)); err == nil {
			Cache.Set(objectName, objectUrl, time.Second*24*60*60)
			return objectUrl, nil
		}

		return nil, errors.New("empty URL")
	} else {
		return cacheUrl.Value(), nil
	}
}

func Upload(folder string, file *multipart.FileHeader) (filename string, err error) {
	ctx := context.Background()
	buffer, err := file.Open()

	if err != nil {
		return "", err
	}
	defer buffer.Close()

	extension := filepath.Ext(file.Filename)
	fileName := strings.TrimSuffix(file.Filename, extension)

	fmt.Println(fileName, extension)

	objectName := fmt.Sprintf("%s/%d-%s%s", folder, time.Now().Unix(), slug.Make(fileName), extension)
	fileBuffer := buffer
	contentType := file.Header["Content-Type"][0]
	fileSize := file.Size

	info, err := GetMinioInstance().PutObject(ctx, bucketName, objectName, fileBuffer, fileSize, minio.PutObjectOptions{ContentType: contentType})

	fmt.Println("Successfully uploaded bytes: ", info)

	return objectName, nil
}

func ReplaceUpload(folder string, file *multipart.FileHeader, oldFilePath string) (filename string, err error) {
	if filename, err = Upload(folder, file); err == nil {
		if oldFilePath != "" {
			err = Delete(oldFilePath)
		}
		return filename, err
	} else {
		return "", err
	}
}

func Delete(filePath string) (err error) {
	ctx := context.Background()
	return GetMinioInstance().RemoveObject(ctx, bucketName, filePath, minio.RemoveObjectOptions{})
}

func GetFileContentType(fh *multipart.FileHeader) (string, error) {
	buf := make([]byte, 512)

	if file, err := fh.Open(); err != nil {
		return "", err
	} else {
		_, err = file.Read(buf)
		file.Close()

		if err != nil {
			return "", err
		}
	}

	contentType := http.DetectContentType(buf)

	return contentType, nil
}
