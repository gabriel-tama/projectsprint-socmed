package image

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Service interface {
	UploadFile(fileName string, fileForm io.Reader, contentType string) (*s3.PutObjectOutput, error)
	GetObject(objectKey string) (*s3.GetObjectOutput, error)
	ListObject() (*s3.ListObjectsV2Output, error)
	GetObjectWithUrl(objectKey string) string
}

type S3ServiceImpl struct {
	S3Client   *s3.Client
	bucketName string
	baseUrl    string
}

func NewS3Service(accessKey string, secretKey string, bucketName string, baseUrl string, region string) S3Service {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region), config.WithCredentialsProvider(
		credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
	))
	if err != nil {
		panic(err)
	}

	client := s3.NewFromConfig(cfg) //, func(o *s3.Options) {
	// o.Credentials = aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""))
	// })

	return &S3ServiceImpl{
		S3Client:   client,
		bucketName: bucketName,
		baseUrl:    baseUrl,
	}
}

func (service *S3ServiceImpl) UploadFile(objectKey string, file io.Reader, contentType string) (*s3.PutObjectOutput, error) {
	bucketName := service.bucketName

	obj, err := service.S3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		ACL:         types.ObjectCannedACL((*aws.String("public-read"))),
		Body:        file,
		ContentType: &contentType,
	})

	if err != nil {
		log.Printf("Couldn't upload file %v to %v. Here's why: %v\n",
			bucketName, objectKey, err)
	}

	return obj, err
}

func (service *S3ServiceImpl) ListObject() (*s3.ListObjectsV2Output, error) {
	bucketName := service.bucketName

	obj, err := service.S3Client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		log.Printf("Couldn't list Object on %s. Here's why: %v\n",
			bucketName, err)
	}

	return obj, err
}

func (service *S3ServiceImpl) GetObject(objectKey string) (*s3.GetObjectOutput, error) {
	bucketName := service.bucketName

	obj, err := service.S3Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		log.Printf("Couldn't list file on %s. Here's why: %v\n",
			bucketName, err)
	}

	return obj, err
}

func (service *S3ServiceImpl) GetObjectWithUrl(objectKey string) string {
	return fmt.Sprintf("%s/%s", service.baseUrl, objectKey)
}
