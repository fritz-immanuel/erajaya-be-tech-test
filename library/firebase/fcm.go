package firebase

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"github.com/fritz-immanuel/erajaya-be-tech-test/configs"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

// initFirebase initializes the Firebase app and returns a storage bucket reference.
func initFirebase(ctx context.Context, config *configs.Config) (*storage.BucketHandle, *types.Error) {
	app, err := firebase.NewApp(ctx, &firebase.Config{StorageBucket: config.FirebaseStorageBucketURL}, option.WithCredentialsFile(config.FirebaseAuthFilePath))
	if err != nil {
		return nil, newError(".initFirebase", err)
	}

	client, err := app.Storage(ctx)
	if err != nil {
		return nil, newError(".initFirebase", err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return nil, newError(".initFirebase", err)
	}
	return bucket, nil
}

// docTypeChecker checks if the uploaded file is an image or document.
func docTypeChecker(fileHeader *multipart.FileHeader) (string, *types.Error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", newError(".docTypeChecker->OpenFile", err)
	}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err = file.Read(buffer); err != nil {
		return "", newError(".docTypeChecker->ReadFile", err)
	}

	contentType := http.DetectContentType(buffer)
	if strings.Contains(contentType, "image") {
		return "images", nil
	}

	return "documents", nil
}

// UploadFile uploads a file to Firebase Storage.
func UploadFile(c *gin.Context, file *multipart.FileHeader, bucketName string) (string, *types.Error) {
	config, errConfig := configs.GetConfiguration()
	if errConfig != nil {
		log.Fatalf("failed to get configuration: %v", errConfig)
	}

	ctx := context.Background()
	bucket, err := initFirebase(ctx, config)
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("%s_%d_%s", bucketName, library.UTCPlus7().Unix(), filepath.Base(file.Filename))
	localPath := filepath.Join("filestore", fileName)

	if err := c.SaveUploadedFile(file, localPath); err != nil {
		return "", newError(".UploadFile->SaveUploadedFile", err)
	}
	defer os.Remove(localPath)

	fileType, err := docTypeChecker(file)
	if err != nil {
		return "", err
	}

	storagePath := fmt.Sprintf("%s/%s/%s", bucketName, fileType, fileName)

	if err := uploadToStorage(ctx, bucket, localPath, storagePath); err != nil {
		return "", err
	}

	return storagePath, nil // Return the signed URL
}

// uploadToStorage uploads the file to the specified path in Firebase Storage.
func uploadToStorage(ctx context.Context, bucket *storage.BucketHandle, localPath, storagePath string) *types.Error {
	file, err := os.Open(localPath)
	if err != nil {
		return newError(".uploadToStorage->OpenFile", err)
	}
	defer file.Close()

	wc := bucket.Object(storagePath).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return newError(".uploadToStorage->Write", err)
	}

	if err := wc.Close(); err != nil {
		return newError(".uploadToStorage->CloseWriter", err)
	}

	return nil
}

// DeleteFile removes a file from Firebase Storage.
func DeleteFile(ctx context.Context, storagePath string) *types.Error {
	config, errConfig := configs.GetConfiguration()
	if errConfig != nil {
		log.Fatalf("failed to get configuration: %v", errConfig)
	}

	bucket, err := initFirebase(ctx, config)
	if err != nil {
		return err
	}

	if err := bucket.Object(storagePath).Delete(ctx); err != nil {
		return newError(".DeleteFile->DeleteObject", err)
	}

	return nil
}

func GenerateSignedURL(fileName string) (string, *types.Error) {
	config, errConfig := configs.GetConfiguration()
	if errConfig != nil {
		log.Fatalf("failed to get configuration: %v", errConfig)
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(config.FirebaseAuthFilePath))
	if err != nil {
		return "", newError(".GenerateSignedURL->NewClient", err)
	}
	defer client.Close()

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	}

	url, err := client.Bucket(config.FirebaseStorageBucketURL).SignedURL(fileName, opts)
	if err != nil {
		return "", newError(".GenerateSignedURL->SignedURL", err)
	}

	return url, nil
}

// newError simplifies error creation.
func newError(path string, err error) *types.Error {
	return &types.Error{
		Path:       path,
		Message:    err.Error(),
		Error:      err,
		StatusCode: http.StatusInternalServerError,
		Type:       "golang-error",
	}
}
