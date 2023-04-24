package db

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/19chonm/461_1_23/logger"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type DB struct {
	client *storage.Client
	bucket *storage.BucketHandle
	ctx    context.Context
}

type PackageSource struct {
	Name    string `json:"name"`
	Content []byte `json:"content"`
}

func NewBucketClient(ctx context.Context, projectID, bucketName string) (*DB, error) {
	credentials, err := ReadCredentialsFromGCS(ctx, "proj-env", "keys.json")
	if err != nil {
		logger.DebugMsg("error reading credentials from the proj-env bucket")
	}

	creds := option.WithCredentialsJSON(credentials)
	client, err := storage.NewClient(ctx, creds)
	if err != nil {
		return nil, err
	}

	bucket := client.Bucket(bucketName)

	return &DB{
		client: client,
		bucket: bucket,
		ctx:    ctx,
	}, nil
}

func (db *DB) ListAllPackages() ([]string, error) {
	var packages []string
	query := &storage.Query{}
	objects := db.bucket.Objects(context.Background(), query)
	for {
		attrs, err := objects.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return nil, err
		}
		packages = append(packages, attrs.Name)
	}

	//c.JSON(http.StatusOK, gin.H{"packages": packages})
	return packages, nil
}

func (db *DB) GetPackageInfo(packageName string) (string, error) {
	obj := db.bucket.Object(packageName)

	attrs, err := obj.Attrs(db.ctx)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Name: %s, Size: %d, LastModified: %s", attrs.Name, attrs.Size, attrs.Updated), nil
}

// func (db *DB) SearchFile(packageName string) ([]string, error) {
// 	var packages []string

// 	query := &storage.Query{
// 		Prefix: packageName,
// 	}

// 	it := db.bucket.Objects(db.ctx, query)
// 	for {
// 		objAttrs, err := it.Next()
// 		if err == iterator.Done {
// 			break
// 		}
// 		if err != nil {
// 			return nil, err
// 		}

// 		packages = append(packages, objAttrs.Name)
// 	}

// 	return packages, nil
// }

func (db *DB) SearchFile(packageID string) (string, error) {
	objectIterator := db.bucket.Objects(db.ctx, nil)

	// Iterate over each object in the bucket
	for {
		objectAttrs, err := objectIterator.Next()
		if err == iterator.Done {
			// We have iterated over all objects in the bucket
			break
		}
		if err != nil {
			return "", fmt.Errorf("error iterating over files in the bucket")
		}

		// Retrieve the metadata for the object
		objectMetadata := objectAttrs.Metadata

		// Check if the object has the metadata field you're interested in
		if objectMetadata["id"] == packageID {
			return objectAttrs.Name, nil

		}
	}
	return "", fmt.Errorf("file not found in the bucket!")
}

func (db *DB) DownloadFile(packageID string, w http.ResponseWriter, r *http.Request) error {
	filename, err := db.SearchFile(packageID)
	if err != nil {
		return err
	}

	obj := db.bucket.Object(filename)
	rc, err := obj.NewReader(db.ctx)
	if err != nil {
		return err
	}
	defer rc.Close()

	// copy the file content from the object reader to the response writer in chunks
	if _, err := io.Copy(w, rc); err != nil {
		return err
	}

	return nil
}

func (db *DB) RemovePackage(packageName string) error {
	obj := db.bucket.Object(packageName)

	err := obj.Delete(db.ctx)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UploadFile(packageName string, r multipart.File, id string) error {
	obj := db.bucket.Object(packageName)

	wc := obj.NewWriter(db.ctx)
	wc.Metadata = map[string]string{
		"id": id,
	}
	defer wc.Close()

	_, err := io.Copy(wc, r)
	if err != nil {
		return err
	}

	return nil
}

//func (db *DB) ListAllPackages

func (db *DB) Close() error {
	return db.client.Close()
}
