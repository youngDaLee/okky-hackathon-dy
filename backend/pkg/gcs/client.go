package gcs

import (
	"context"
	"io"
	"time"
)

type Client struct{}

func NewClient(ctx context.Context, bucket string) (*Client, error) {
	return nil, nil
}

func (c *Client) Upload(ctx context.Context, objectName string, file io.Reader, contentType string) (string, error) {
	return "", nil
}

func (c *Client) GenerateSignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	return "", nil
}

func (c *Client) Delete(ctx context.Context, objectName string) error {
	return nil
}
