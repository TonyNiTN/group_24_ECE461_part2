package db

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
)

type Credentials struct {
	Type                string `json:"type"`
	ProjectID           string `json:"project_id"`
	PrivateKeyID        string `json:"private_key_id"`
	PrivateKey          string `json:"private_key"`
	ClientEmail         string `json:"client_email"`
	ClientID            string `json:"client_id"`
	AuthURI             string `json:"auth_uri"`
	TokenURI            string `json:"token_uri"`
	AuthProviderX509URL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL   string `json:"client_x509_cert_url"`
}

func LoadCredentials(credentialsBytes []byte) (*google.Credentials, error) {
	var credentials Credentials
	fmt.Println(string(credentialsBytes))
	err := json.Unmarshal(credentialsBytes, &credentials)
	if err != nil {
		return nil, err
	}
	return google.CredentialsFromJSON(context.Background(), credentialsBytes, "https://www.googleapis.com/auth/cloud-platform")
}

func ReadCredentialsFromGCS(ctx context.Context, bucketName, objectName string) ([]byte, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objectName)

	rc, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return ioutil.ReadAll(rc)
}
