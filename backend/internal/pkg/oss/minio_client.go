package oss

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOConfig MinIO 配置
type MinIOConfig struct {
	Endpoint        string
	AccessKey       string
	SecretKey       string
	Bucket          string
	Region          string
	UseSSL          bool
}

// MinIOClient MinIO 客户端实现
type MinIOClient struct {
	client *minio.Client
	bucket string
}

// TestConnection ?? OSS/MinIO ??
func (c *MinIOClient) TestConnection(ctx context.Context) error {
	_, err := c.client.BucketExists(ctx, c.bucket)
	if err != nil {
		return fmt.Errorf("?? OSS/MinIO ????: %w", err)
	}
	return nil
}


// NewMinIOClient 创建 MinIO 客户端实例
func NewMinIOClient(config *MinIOConfig) (*MinIOClient, error) {
	// 初始化 MinIO 客户端
	// 注意：endpoint 应该只是主机名:端口，不包含 http:// 或 https://
	// SDK 会根据 Secure 参数自动选择协议
	options := &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
		Region: config.Region,
	}

	// 打印调试信息
	fmt.Printf("正在连接 MinIO: endpoint=%s, secure=%v\n", config.Endpoint, config.UseSSL)

	client, err := minio.New(config.Endpoint, options)
	if err != nil {
		return nil, fmt.Errorf("初始化 MinIO 客户端失败: %w", err)
	}

	// 检查 bucket 是否存在，不存在则创建
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, config.Bucket)
	if err != nil {
		return nil, fmt.Errorf("检查 bucket 失败: %w", err)
	}

	if !exists {
		// 创建 bucket
		err = client.MakeBucket(ctx, config.Bucket, minio.MakeBucketOptions{
			Region: config.Region,
		})
		if err != nil {
			return nil, fmt.Errorf("创建 bucket 失败: %w", err)
		}
	}

	return &MinIOClient{
		client: client,
		bucket: config.Bucket,
	}, nil
}

// GeneratePresignedUploadURL 生成预签名上传 URL
func (c *MinIOClient) GeneratePresignedUploadURL(ctx context.Context, fileKey string, expiresIn time.Duration) (string, error) {
	// 设置上传参数
	presignedURL, err := c.client.PresignedPutObject(ctx, c.bucket, fileKey, expiresIn)
	if err != nil {
		return "", fmt.Errorf("生成预签名上传 URL 失败: %w", err)
	}
	return presignedURL.String(), nil
}

// GeneratePresignedDownloadURL 生成预签名下载 URL
func (c *MinIOClient) GeneratePresignedDownloadURL(ctx context.Context, fileKey string, expiresIn time.Duration) (string, error) {
	// 设置下载参数
	reqParams := make(url.Values)
	// 设置响应头，触发浏览器下载
	reqParams.Set("response-content-disposition", "attachment")

	presignedURL, err := c.client.PresignedGetObject(ctx, c.bucket, fileKey, expiresIn, reqParams)
	if err != nil {
		return "", fmt.Errorf("生成预签名下载 URL 失败: %w", err)
	}
	return presignedURL.String(), nil
}

// DeleteFile 删除文件
func (c *MinIOClient) DeleteFile(ctx context.Context, fileKey string) error {
	err := c.client.RemoveObject(ctx, c.bucket, fileKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}
	return nil
}

// GetFile 获取文件
func (c *MinIOClient) GetFile(ctx context.Context, fileKey string) ([]byte, error) {
	obj, err := c.client.GetObject(ctx, c.bucket, fileKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取文件失败: %w", err)
	}
	defer obj.Close()

	// 读取文件内容
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(obj)
	if err != nil {
		return nil, fmt.Errorf("读取文件内容失败: %w", err)
	}

	return buf.Bytes(), nil
}

// FileExists 检查文件是否存在
func (c *MinIOClient) FileExists(ctx context.Context, fileKey string) (bool, error) {
	_, err := c.client.StatObject(ctx, c.bucket, fileKey, minio.StatObjectOptions{})
	if err != nil {
		// 检查是否是"不存在"错误
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("检查文件是否存在失败: %w", err)
	}
	return true, nil
}
