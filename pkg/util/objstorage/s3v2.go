package objstorage

//
//import (
//	"cambridge-hit.com/gin-base/activateserver/pkg"
//	"context"
//	"fmt"
//	"github.com/aws/aws-sdk-go-v2/aws"
//	awsConfig "github.com/aws/aws-sdk-go-v2/config"
//	"github.com/aws/aws-sdk-go-v2/credentials"
//	"github.com/aws/aws-sdk-go-v2/service/s3"
//	"io"
//	"log"
//	"time"
//)
//
//var S3c *S3Client
//
//// S3Client 实现 OSSClient 接口
//type S3Client struct {
//	client        *s3.Client
//	presignClient *s3.PresignClient
//	info          *S3Info
//}
//
//type S3Info struct {
//	roleArn         string
//	region          string
//	bucketName      string
//	accessKeyID     string
//	accessKeySecret string
//	endpoint        string
//}
//
//// NewS3Client 创建一个新的 S3Client 实例
//func NewS3Client(product, region, bucketName, accessKeyID, secretKeyAccess string) *S3Client {
//	s3info := S3Info{
//		region:          region,
//		bucketName:      bucketName,
//		accessKeyID:     accessKeyID,
//		accessKeySecret: secretKeyAccess,
//	}
//
//	switch product {
//	case "aliyun":
//		if pkg.Conf.Env == "prod" {
//			//线上internal
//			s3info.endpoint = "https://oss-" + s3info.region + "-internal.aliyuncs.com"
//		} else {
//			s3info.endpoint = "https://oss-" + s3info.region + ".aliyuncs.com"
//		}
//	}
//	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
//		awsConfig.WithRegion(s3info.region),
//		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s3info.accessKeyID, s3info.accessKeySecret, "")),
//		awsConfig.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
//			return aws.Endpoint{
//				URL: s3info.endpoint,
//			}, nil
//		})),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return &S3Client{
//		client:        s3.NewFromConfig(cfg),
//		presignClient: s3.NewPresignClient(s3.NewFromConfig(cfg)),
//		info:          &s3info,
//	}
//}
//
//// UploadObject 上传对象
//func (c *S3Client) UploadObject(objectName string, reader io.Reader) error {
//	_, err := c.client.PutObject(context.TODO(), &s3.PutObjectInput{
//		Bucket: aws.String(c.info.bucketName),
//		Key:    aws.String(objectName),
//		Body:   reader,
//	})
//	return err
//}
//
//// DownloadObject 下载对象
//func (c *S3Client) DownloadObject(objectName string) (io.Reader, error) {
//	result, err := c.client.GetObject(context.TODO(), &s3.GetObjectInput{
//		Bucket: aws.String(c.info.bucketName),
//		Key:    aws.String(objectName),
//	})
//	if err != nil {
//		return nil, err
//	}
//	return result.Body, nil
//}
//
//// DeleteObject 删除对象
//func (c *S3Client) DeleteObject(objectName string) error {
//	_, err := c.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
//		Bucket: aws.String(c.info.bucketName),
//		Key:    aws.String(objectName),
//	})
//	return err
//}
//
//// ListObjects 列出对象
//func (c *S3Client) ListObjects(prefix string) ([]string, error) {
//	result, err := c.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
//		Bucket: aws.String(c.info.bucketName),
//		Prefix: aws.String(prefix),
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	var keys []string
//	for _, item := range result.Contents {
//		keys = append(keys, *item.Key)
//	}
//	return keys, nil
//}
//
//// GetTemporaryURL 获取临时URL
//func (c *S3Client) GetTemporaryURL(objectName string) (string, error) {
//	presignResult, err := c.presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
//		Bucket: aws.String(c.info.bucketName),
//		Key:    aws.String(objectName),
//	}, func(opts *s3.PresignOptions) {
//		opts.Expires = 15 * time.Minute
//	})
//	if err != nil {
//		return "", err
//	}
//	return presignResult.URL, nil
//}
//
//// PutObjectRequest 生成预签名上传链接供客户端直传
//func (c *S3Client) PutObjectRequest(objectName string, expires ...time.Duration) (string, error) {
//	// 生成预签名URL
//	req, err := c.presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
//		Bucket: aws.String(c.info.bucketName),
//		Key:    aws.String(objectName),
//	}, func(opts *s3.PresignOptions) {
//		if len(expires) > 0 {
//			opts.Expires = expires[0]
//		} else {
//			opts.Expires = time.Minute * 15
//		}
//	})
//	if err != nil {
//		return "", fmt.Errorf("failed to sign request, %v", err)
//	}
//	return req.URL, nil
//}
