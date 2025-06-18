package auth

import (
	"context"
	"distributed-web-scrapper/services/scraper/internal/config"
	"fmt"

	"golang.org/x/oauth2"
)

type OAuthClient struct {
	linkedin  *oauth2.Config
	youtube   *oauth2.Config
	instagram *oauth2.Config
}

func NewOAuthClient(cfg config.OAuthConfig) (*OAuthClient, error) {
	return &OAuthClient{
		linkedin: &oauth2.Config{
			ClientID:     cfg.LinkedInClientID,
			ClientSecret: cfg.LinkedInClientSecret,
			Endpoint:     oauth2.Endpoint{ /* LinkedIn OAuth endpoints */ },
		},
		youtube: &oauth2.Config{
			ClientID:     cfg.YouTubeClientID,
			ClientSecret: cfg.YouTubeClientSecret,
			Endpoint:     oauth2.Endpoint{ /* YouTube OAuth endpoints */ },
		},
		instagram: &oauth2.Config{
			ClientID:     cfg.InstagramClientID,
			ClientSecret: cfg.InstagramClientSecret,
			Endpoint:     oauth2.Endpoint{ /* Instagram OAuth endpoints */ },
		},
	}, nil
}

func (c *OAuthClient) GetToken(ctx context.Context, platform string) (*oauth2.Token, error) {
	switch platform {
	case "linkedin":
		return c.linkedin.TokenSource(ctx, nil).Token()
	case "youtube":
		return c.youtube.TokenSource(ctx, nil).Token()
	case "instagram":
		return c.instagram.TokenSource(ctx, nil).Token()
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}