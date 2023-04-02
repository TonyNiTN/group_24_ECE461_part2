package db

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type DB struct {
	client *storage.Client
	bucket *storage.BucketHandle
	ctx    context.Context
}

func NewBucketClient(ctx context.Context, projectID, bucketName string) (*DB, error) {
	client, err := storage.NewClient(ctx)
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

func (db *DB) SearchPackage(packageName string) ([]string, error) {
	var packages []string

	query := &storage.Query{
		Prefix: packageName,
	}

	it := db.bucket.Objects(db.ctx, query)
	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		packages = append(packages, objAttrs.Name)
	}

	return packages, nil
}

func (db *DB) DownloadPackage(packageName string, w io.Writer) error {
	obj := db.bucket.Object(packageName)
	rc, err := obj.NewReader(db.ctx)
	if err != nil {
		return err
	}
	defer rc.Close()

	_, err = io.Copy(w, rc)
	if err != nil {
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

func (db *DB) UploadPackage(packageName string, r multipart.File) error {
	obj := db.bucket.Object(packageName)

	wc := obj.NewWriter(db.ctx)
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
