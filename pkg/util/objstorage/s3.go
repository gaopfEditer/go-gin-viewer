package objstorage

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"log"
	"time"
)

var S3c *S3Client
var S3cExternal *S3Client

// S3Client 实现 OSSClient 接口
type S3Client struct {
	client *s3.S3
	info   *S3Info
}

type S3Info struct {
	region          string
	bucketName      string
	accessKeyID     string
	accessKeySecret string
	endpoint        string
}

// NewS3Client 创建一个新的 S3Client 实例
func NewS3Client(product, region, bucketName, accessKeyID, secretKeyAccess, endpoint string) *S3Client {
	s3info := S3Info{
		region:          region,
		bucketName:      bucketName,
		accessKeyID:     accessKeyID,
		accessKeySecret: secretKeyAccess,
		endpoint:        endpoint,
	}

	// 创建 AWS 会话
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(s3info.region),
		Credentials: credentials.NewStaticCredentials(s3info.accessKeyID, s3info.accessKeySecret, ""),
		Endpoint:    aws.String(s3info.endpoint),
	})
	if err != nil {
		log.Fatal(err)
	}

	return &S3Client{
		client: s3.New(sess),
		info:   &s3info,
	}
}

// InitS3Clients 初始化内部和外部的 S3Client 实例
func InitS3Clients(product, region, bucketName, accessKeyID, secretKeyAccess string) {
	var internalEndpoint, externalEndpoint string

	switch product {
	case "aliyun":
		internalEndpoint = "https://oss-" + region + "-internal.aliyuncs.com"
		externalEndpoint = "https://oss-" + region + ".aliyuncs.com"
	default:
		log.Fatal("Unsupported product")
	}

	S3c = NewS3Client(product, region, bucketName, accessKeyID, secretKeyAccess, internalEndpoint)
	S3cExternal = NewS3Client(product, region, bucketName, accessKeyID, secretKeyAccess, externalEndpoint)
}

// UploadObject 上传对象
func (c *S3Client) UploadObject(objectName string, reader io.Reader) error {
	// 读取数据到内存中
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return err
	}

	// 使用 bytes.Reader，它实现了 io.ReadSeeker
	_, err = c.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(c.info.bucketName),
		Key:    aws.String(objectName),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	return err
}

// DownloadObject 下载对象
func (c *S3Client) DownloadObject(objectName string) (io.ReadCloser, error) {
	result, err := c.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(c.info.bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

// DeleteObject 删除对象
func (c *S3Client) DeleteObject(objectName string) error {
	_, err := c.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(c.info.bucketName),
		Key:    aws.String(objectName),
	})
	return err
}

// ListObjects 列出对象
func (c *S3Client) ListObjects(prefix string) ([]string, error) {
	result, err := c.client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(c.info.bucketName),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return nil, err
	}

	var keys []string
	for _, item := range result.Contents {
		keys = append(keys, *item.Key)
	}
	return keys, nil
}

// GetTemporaryURL 获取临时URL
func (c *S3Client) GetTemporaryURL(objectName string) (string, error) {
	req, _ := c.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(c.info.bucketName),
		Key:    aws.String(objectName),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", err
	}
	return urlStr, nil
}

// PutObjectRequest 生成预签名上传链接供客户端直传
func (c *S3Client) PutObjectRequest(objectName string, expires ...time.Duration) (string, error) {
	duration := time.Minute * 15
	if len(expires) > 0 {
		duration = expires[0]
	}
	req, _ := c.client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(c.info.bucketName),
		Key:    aws.String(objectName),
	})
	urlStr, err := req.Presign(duration)
	if err != nil {
		return "", fmt.Errorf("failed to sign request, %v", err)
	}
	return urlStr, nil
}
