package main

import (
	"encoding/base64"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

// DecryptEnvVar decrypts an environment variable
func DecryptEnvVar(kmsKeyID, variable string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(os.Getenv(variable))
	if err != nil {
		return "", err
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		return "", err
	}

	svc := kms.New(sess)

	out, err := svc.Decrypt(&kms.DecryptInput{
		CiphertextBlob: decoded,
	})
	if err != nil {
		return "", err
	}

	return string(out.Plaintext), nil
}
