package firebase

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
)

const maxImageSize = 200 * 1024 // 200KB

type UploadResult struct {
	URL string
}

func (c *Client) UploadProductImage(ctx context.Context, file io.Reader, fileName string) (*UploadResult, error) {
	bucket, err := c.Storage.DefaultBucket()
	if err != nil {
		return nil, fmt.Errorf("get storage bucket: %w", err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	if len(data) > maxImageSize {
		return nil, fmt.Errorf("image too large: %d bytes (max %d)", len(data), maxImageSize)
	}

	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode image config: %w", err)
	}

	_ = config // could enforce min dimensions later

	ext := strings.ToLower(filepath.Ext(fileName))
	if ext != ".webp" && ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return nil, fmt.Errorf("unsupported image format: %s (use webp, jpg, or png)", ext)
	}

	objectPath := fmt.Sprintf("products/%s%s", strings.TrimSuffix(fileName, ext), ext)

	wc := bucket.Object(objectPath).NewWriter(ctx)
	wc.ContentType = fmt.Sprintf("image/%s", format)
	wc.CacheControl = "public, max-age=31536000"

	if _, err := wc.Write(data); err != nil {
		return nil, fmt.Errorf("upload to storage: %w", err)
	}

	if err := wc.Close(); err != nil {
		return nil, fmt.Errorf("close storage writer: %w", err)
	}

	attrs, err := bucket.Object(objectPath).Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("get object attrs: %w", err)
	}

	return &UploadResult{
		URL: fmt.Sprintf("https://storage.googleapis.com/%s/%s", attrs.Bucket, attrs.Name),
	}, nil
}

func (c *Client) DeleteImage(ctx context.Context, url string) error {
	bucket, err := c.Storage.DefaultBucket()
	if err != nil {
		return fmt.Errorf("get storage bucket: %w", err)
	}

	objectPath := extractObjectPath(url)
	if objectPath == "" {
		return nil
	}

	obj := bucket.Object(objectPath)
	if err := obj.Delete(ctx); err != nil && err != storage.ErrObjectNotExist {
		return fmt.Errorf("delete object: %w", err)
	}

	return nil
}

func extractObjectPath(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return ""
	}
	return strings.Join(parts[len(parts)-2:], "/")
}
