package db

import (
	"context"
	"fmt"
	"strings"

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
	credentials, err := ReadCredentialsFromGCS(ctx, "proj-env", "keys.json")
	if err != nil {
		logger.DebugMsg("error reading credentials from the proj-env bucket")
	}

	creds := option.WithCredentialsJSON(credentials)
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
	var query firestore.Query
	collection := client.Collection("packages")
	name = strings.ToLower(name)
	if name == "" {
		query = collection.OrderBy("net score", firestore.Desc)
	} else {
		query = collection.Where("lowercase name", ">=", name).Where("lowercase name", "<=", name+"\uf8ff")
	}

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
	if rating == nil {
		logger.DebugMsg("ratings are nil! make sure you entered a valid github or NPM URL")
		packageInfo.BusFactorScore = "0.00"
		packageInfo.RampUpScore = "0.00"
		packageInfo.CorrectnessScore = "0.00"
		packageInfo.ResponsivenessScore = "0.00"
		packageInfo.LicenseScore = "0.00"
		packageInfo.VersionScore = "0.00"
		packageInfo.ReviewScore = "0.00"
		packageInfo.NetScore = "0.00"
	} else {
		packageInfo.BusFactorScore = fmt.Sprintf("%.2f", rating.Busfactor)
		packageInfo.RampUpScore = fmt.Sprintf("%.2f", rating.Rampup)
		packageInfo.CorrectnessScore = fmt.Sprintf("%.2f", rating.Correctness)
		packageInfo.ResponsivenessScore = fmt.Sprintf("%.2f", rating.Responsiveness)
		packageInfo.LicenseScore = fmt.Sprintf("%.2f", rating.License)
		packageInfo.VersionScore = fmt.Sprintf("%.2f", rating.Version)
		packageInfo.ReviewScore = fmt.Sprintf("%.2f", rating.Review)
		packageInfo.NetScore = fmt.Sprintf("%.2f", rating.NetScore)
	}

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

func (db *Table) ResetTable() error {
	iter := db.client.Collection("packages").Documents(db.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.DebugMsg("error getting next object in the database")
			return err
		}

		_, err = db.client.Collection("packages").Doc(doc.Ref.ID).Delete(db.ctx)
		if err != nil {
			logger.DebugMsg("error deleting pbject from database")
			return err
		}
	}

	return nil
}

func (db *Table) Close() error {
	return db.client.Close()
}
