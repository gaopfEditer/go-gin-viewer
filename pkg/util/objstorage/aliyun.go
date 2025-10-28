package objstorage

import (
	"cambridge-hit.com/gin-base/activateserver/pkg/util/str"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//var AliOSSClient *AliyunOSSClient

// STSCredentials 定义 STS 凭证结构
type STSCredentials struct {
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
	Expiration      time.Time
}

type AliyunOSSClient struct {
	client          *oss.Client
	bucket          *oss.Bucket
	region          string
	accessKeyID     string
	accessKeySecret string
	roleArn         string
}

// NewAliyunOSSClient 创建新的阿里云OSS客户端
func NewAliyunOSSClient(region, bucketName, accessKeyID, accessKeySecret, roleArn string) (*AliyunOSSClient, error) {
	if resource.Conf.Env == "prod" {
		region += "-internal"
	}
	endpoint := "oss-" + region + ".aliyuncs.com"
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	return &AliyunOSSClient{
		client:          client,
		bucket:          bucket,
		region:          region,
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
		roleArn:         roleArn,
	}, nil
}

// UploadObject 上传对象
func (c *AliyunOSSClient) UploadObject(objectName string, reader io.Reader) error {
	return c.bucket.PutObject(objectName, reader)
}

// DownloadObject 下载对象
func (c *AliyunOSSClient) DownloadObject(objectName string) (io.Reader, error) {
	return c.bucket.GetObject(objectName)
}

// DeleteObject 删除对象
func (c *AliyunOSSClient) DeleteObject(objectName string) error {
	return c.bucket.DeleteObject(objectName)
}

// ListObjects 列出对象
func (c *AliyunOSSClient) ListObjects(prefix string) ([]string, error) {
	var fileList []string
	marker := ""
	for {
		lsRes, err := c.bucket.ListObjects(oss.Prefix(prefix), oss.Marker(marker))
		if err != nil {
			return nil, err
		}
		for _, object := range lsRes.Objects {
			fileList = append(fileList, object.Key)
		}
		if !lsRes.IsTruncated {
			break
		}
		marker = lsRes.NextMarker
	}
	return fileList, nil
}

// GetTemporaryURL 获取对象的临时URL，默认24小时有效
func (c *AliyunOSSClient) GetTemporaryURL(objectName string) (string, error) {
	expireTime := int64(86400)
	return c.bucket.SignURL(objectName, oss.HTTPGet, expireTime)
}

func (c *AliyunOSSClient) BatchUpload(localDir, ossDir string) error {
	return filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(localDir, path)
			if err != nil {
				return fmt.Errorf("获取相对路径失败: %v", err)
			}
			objectKey := filepath.Join(ossDir, relPath)

			err = c.bucket.PutObjectFromFile(objectKey, path)
			if err != nil {
				fmt.Printf("上传文件 %s 失败: %v\n", path, err)
				return err // 如果需要在单个文件上传失败时停止整个批量上传，可以返回这个错误
			} else {
				fmt.Printf("成功上传文件 %s \n", path)
			}
		}
		return nil
	})
}

// GenerateSTSToken 生成STS令牌
func (c *AliyunOSSClient) GenerateSTSToken(durationSeconds int64, policy string) (*STSCredentials, error) {
	// 创建 STS 客户端时指定 Endpoint
	//endpoint := "sts." + c.region + ".aliyuncs.com"
	stsRegion := strings.Trim(c.region, "-internal")
	client, err := sts.NewClientWithAccessKey(stsRegion, c.accessKeyID, c.accessKeySecret)
	if err != nil {
		return nil, err
	}

	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	request.RoleArn = c.roleArn
	request.RoleSessionName = "sts" + str.GenerateCode(5) // 使用一个合适的会话名称
	request.DurationSeconds = requests.NewInteger(int(durationSeconds))
	request.Policy = policy

	response, err := client.AssumeRole(request)
	if err != nil {
		return nil, err
	}

	expiration, err := time.Parse("2006-01-02T15:04:05Z", response.Credentials.Expiration)
	if err != nil {
		return nil, fmt.Errorf("解析过期时间失败: %v", err)
	}

	credentials := &STSCredentials{
		AccessKeyId:     response.Credentials.AccessKeyId,
		AccessKeySecret: response.Credentials.AccessKeySecret,
		SecurityToken:   response.Credentials.SecurityToken,
		Expiration:      expiration,
	}

	return credentials, nil
}
