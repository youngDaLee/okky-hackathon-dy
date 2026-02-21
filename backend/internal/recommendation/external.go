package recommendation

import "context"

type ExternalSearcher interface {
	SearchRecipes(ctx context.Context, query string, limit int) ([]*ExternalRecipe, error)
}

type ExternalClient struct{}

func NewExternalClient(youtubeAPIKey, googleSearchKey, googleSearchCX string) *ExternalClient {
	return &ExternalClient{}
}

func (c *ExternalClient) SearchRecipes(ctx context.Context, query string, limit int) ([]*ExternalRecipe, error) {
	return nil, nil
}
