package utils

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"io"
	"log"
	"strings"
)

func UploadPictures(fileName, uploadFilePath string, file io.Reader) (string, error) {
	//cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(Settings.CloudFlareRegion))

	// Configure credentials using the credentials package
	r2Credentials := aws.NewCredentialsCache(
		credentials.NewStaticCredentialsProvider(
			Settings.CloudFlareAccessKeyID,     // Cloudflare Access Key ID
			Settings.CloudFlareAccessSecretKey, // Cloudflare Secret Access Key
			"",                                 // No session token
		),
	)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(r2Credentials),
		config.WithRegion("auto"), // "auto" to let the resolver manage it
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			if service == s3.ServiceID && region == "auto" {
				return aws.Endpoint{
					URL:           Settings.CloudFlareBucketUrl, // Replace with your Cloudflare account ID
					SigningRegion: "auto",                       // Set signing region for Cloudflare
				}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})),
	)
	if err != nil {
		log.Println("unable to load SDK config, " + err.Error())

		return "", err
	}
	client := s3.NewFromConfig(cfg)
	uploadKey := fmt.Sprintf("%s/%s/%s", Settings.Environment, uploadFilePath, fileName)
	_, err = client.PutObject(context.TODO(),
		&s3.PutObjectInput{
			Bucket:      aws.String(Settings.CloudFlareBucket),
			Key:         aws.String(uploadKey),
			Body:        file,
			ContentType: aws.String("image/png"),
		})

	if err != nil {
		log.Println("Unable to upload file %v", err)
		return "", err
	}
	cloudflareURL := fmt.Sprintf("%s/%s", Settings.CloudFlareCdnUrl, uploadKey)
	return cloudflareURL, nil
}

func DeletePicture(fileURL string) error {

	// Configure credentials using the credentials package
	r2Credentials := aws.NewCredentialsCache(
		credentials.NewStaticCredentialsProvider(
			Settings.CloudFlareAccessKeyID,     // Cloudflare Access Key ID
			Settings.CloudFlareAccessSecretKey, // Cloudflare Secret Access Key
			"",                                 // No session token
		),
	)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(r2Credentials),
		config.WithRegion("auto"), // "auto" to let the resolver manage it
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			if service == s3.ServiceID && region == "auto" {
				return aws.Endpoint{
					URL:           Settings.CloudFlareBucketUrl, // Replace with your Cloudflare account ID
					SigningRegion: "auto",                       // Set signing region for Cloudflare
				}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})),
	)
	if err != nil {
		log.Println("unable to load SDK config, " + err.Error())

		return err
	}
	// Create S3 client
	client := s3.NewFromConfig(cfg)

	// Extract the key from the file URL by stripping the bucket's base URL
	baseURL := fmt.Sprintf("%s/", Settings.CloudFlareBucketUrl)
	key := strings.TrimPrefix(fileURL, baseURL)
	// Delete the file from S3
	_, err = client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(Settings.CloudFlareBucket),
		Key:    aws.String(key),
	})

	return err
}
