package db

import (
	"context"

	"cloud.google.com/go/firestore"
)

type Package struct {
	ID                  string `firestore:"id"`
	URL                 string `firestore:"url"`
	Name                string `firestore:"name"`
	NetScore            string `firestore:"net score"`
	LicenseScore        string `firestore:"license score"`
	CorrectnessScore    string `firestore:"correctness score"`
	BusFactorScore      string `firestore:"bus factor score"`
	ResponsivenessScore string `firestore:"responsiveness score"`
	RampUpScore         string `firestore:"ramp up time score"`
	VersionScore        string `firestore:"version pinning score"`
	ReviewScore         string `firestore:"reviewed pr score"`
}

func (p *Package) Save(ctx context.Context, client *firestore.Client) error {
	_, err := client.Collection("packages").Doc(p.Name).Set(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

func NewPackage(url string, name string, netScore string) *Package {
	return &Package{
		URL:      url,
		Name:     name,
		NetScore: netScore,
	}
}
