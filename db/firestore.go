package db

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/19chonm/461_1_23/logger"
	"github.com/19chonm/461_1_23/worker"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Table struct {
	client     *firestore.Client
	collection *firestore.CollectionRef
	ctx        context.Context
}

func NewFirestoreClient(ctx context.Context, projectID, tableName string) (*Table, error) {
	creds := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	client, err := firestore.NewClient(ctx, "trusted-package-registry", creds)
	if err != nil {
		return nil, err
	}

	collection := client.Collection(tableName)

	return &Table{
		client:     client,
		collection: collection,
		ctx:        ctx,
	}, nil
}

func (db *Table) GetClient() *firestore.Client {
	if db.client != nil {
		return db.client
	}

	return nil
}

func (db *Table) GetCollection() *firestore.CollectionRef {
	if db.collection != nil {
		return db.collection
	}

	return nil
}

func (db *Table) GetCtx() context.Context {
	if db.ctx != nil {
		return db.ctx
	}

	return nil
}

// Create a new package
func (db *Table) UploadPackage(ctx context.Context, client *firestore.Client, packageData *Package, id string) error {

	_, err := client.Collection("packages").Doc(id).Set(ctx, packageData)
	if err != nil {
		logger.DebugMsg("error creating package")
		return err
	}

	return nil
}

func (db *Table) SearchPackage(ctx context.Context, client *firestore.Client, name string) ([]*Package, error) {
	var packages []*Package
	collection := client.Collection("packages")
	query := collection.Where("name", ">=", name).Where("name", "<=", name+"\uf8ff")

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		logger.DebugMsg("error searching packages")
		return nil, fmt.Errorf("error searching packages")
	}

	for _, doc := range docs {
		var packageData Package
		err = doc.DataTo(&packageData)
		if err != nil {
			logger.DebugMsg("error copying document into Package structure")
			continue
		}
		packages = append(packages, &packageData)
	}

	return packages, nil

}

// Read a package by ID
func (db *Table) GetPackage(ctx context.Context, client *firestore.Client, packageID string) (*Package, error) {
	docRef := client.Collection("packages").Doc(packageID)

	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	var packageData Package
	if err := docSnapshot.DataTo(&packageData); err != nil {
		return nil, err
	}

	return &packageData, nil
}

func (db *Table) ScorePackage(ctx context.Context, client *firestore.Client, url string, packageInfo *Package) {
	rating := worker.RunTask(url)
	packageInfo.BusFactorScore = fmt.Sprintf("%f", rating.Busfactor)
	packageInfo.RampUpScore = fmt.Sprintf("%f", rating.Rampup)
	packageInfo.CorrectnessScore = fmt.Sprintf("%f", rating.Correctness)
	packageInfo.ResponsivenessScore = fmt.Sprintf("%f", rating.Responsiveness)
	packageInfo.LicenseScore = fmt.Sprintf("%f", rating.License)
	packageInfo.VersionScore = fmt.Sprintf("%f", rating.Version)
	packageInfo.ReviewScore = fmt.Sprintf("%f", rating.Review)
	packageInfo.NetScore = fmt.Sprintf("%f", rating.NetScore)
	err := db.UpdatePackage(ctx, client, packageInfo, packageInfo.ID)
	if err != nil {
		logger.DebugMsg("Error updating scores on the database!")
	}
}

// Update a package
func (db *Table) UpdatePackage(ctx context.Context, client *firestore.Client, packageData *Package, id string) error {

	_, err := client.Collection("packages").Doc(id).Set(ctx, packageData)
	if err != nil {
		return err
	}

	return nil
}

// Delete a package
func (db *Table) DeletePackage(ctx context.Context, client *firestore.Client, packageID string) error {
	_, err := client.Collection("packages").Doc(packageID).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

// List all packages
func (db *Table) ListPackages() ([]*Package, error) {
	var packages []*Package

	iter := db.client.Collection("packages").Documents(db.ctx)
	for {
		docSnapshot, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var packageData Package
		if err := docSnapshot.DataTo(&packageData); err != nil {
			return nil, err
		}

		packages = append(packages, &packageData)
	}

	return packages, nil
}

func (db *Table) Close() error {
	return db.client.Close()
}
