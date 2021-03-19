package objectstoragehandler

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/ScienceObjectsDB/go-api/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
)

type S3Handler struct {
	S3Client      *s3.Client
	PresignClient *s3.PresignClient
	S3Endpoint    string
}

//NewS3Handler Creates a new  handler to handle S3 related operations
func NewS3Handler() (*S3Handler, error) {

	endpoint := "http://localhost:9000"

	if viper.GetString("Config.S3.Endpoint") != "" {
		endpoint = viper.GetString("Config.S3.Endpoint")
	}

	if os.Getenv("MINIO_ENDPOINT_URL") != "" {
		endpoint = os.Getenv("MINIO_ENDPOINT_URL")
	}

	hostnameImmutable := false
	if strings.HasSuffix(os.Args[0], ".test") {
		hostnameImmutable = true
	}

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("RegionOne"),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: endpoint, HostnameImmutable: hostnameImmutable}, nil
			})),
	)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	presignClient := s3.NewPresignClient(client)

	handler := S3Handler{
		S3Client:      client,
		PresignClient: presignClient,
		S3Endpoint:    endpoint,
	}

	return &handler, nil
}

// CreatePresignedUploadLink Creates new upload link
func (s3handler *S3Handler) CreatePresignedUploadLink(object *models.DatasetObjectEntry) (string, error) {
	presignedRequestURL, err := s3handler.PresignClient.PresignPutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(object.GetLocation().GetBucket()),
		Key:    aws.String(object.GetLocation().GetKey()),
	})

	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return presignedRequestURL.URL, nil
}

// CreatePresignedDownloadLink Creates a  new download link
func (s3handler *S3Handler) CreatePresignedDownloadLink(object *models.DatasetObjectEntry) (string, error) {
	presignedRequestURL, err := s3handler.PresignClient.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(object.GetLocation().GetBucket()),
		Key:    aws.String(object.GetLocation().GetKey()),
	})

	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return presignedRequestURL.URL, nil
}
