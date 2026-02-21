package vertexai

import "context"

type DetectedLabel struct {
	Name       string
	Confidence float64
}

type VisionResult struct {
	Labels   []DetectedLabel
	FullText string
}

type Client struct{}

func NewClient(ctx context.Context, projectID, location string) (*Client, error) {
	return nil, nil
}

func (c *Client) AnalyzeReceipt(ctx context.Context, imageURL string) (*VisionResult, error) {
	return nil, nil
}

func (c *Client) AnalyzeFridge(ctx context.Context, imageURL string) (*VisionResult, error) {
	return nil, nil
}
