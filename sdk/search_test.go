package sdk

import (
	"context"
	"fmt"
	"testing"

	corekugou "github.com/rushingrain/kugou-music-api/core/kugou"
)

func TestSearchMatchesJSRouteByType(t *testing.T) {
	tests := []struct {
		searchType string
		wantURL    string
	}{
		{searchType: "song", wantURL: "/v3/search/song"},
		{searchType: "album", wantURL: "/v1/search/album"},
		{searchType: "author", wantURL: "/v1/search/author"},
		{searchType: "mv", wantURL: "/v1/search/mv"},
		{searchType: "lyric", wantURL: "/v1/search/lyric"},
		{searchType: "special", wantURL: "/v1/search/special"},
	}

	for _, test := range tests {
		t.Run(test.searchType, func(t *testing.T) {
			client, err := New(WithCookie(map[string]string{"token": "token-value"}))
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			client.requester = func(_ context.Context, cfg corekugou.RequestConfig) (corekugou.Response, error) {
				if got := cfg.URL; got != test.wantURL {
					t.Fatalf("URL = %q, want %q", got, test.wantURL)
				}
				if got := fmt.Sprint(cfg.Params["keyword"]); got != "周杰伦" {
					t.Fatalf("keyword = %q, want %q", got, "周杰伦")
				}
				if _, ok := cfg.Params["type"]; ok {
					t.Fatal("unexpected type query parameter")
				}
				if _, ok := cfg.Params["cookie"]; ok {
					t.Fatal("unexpected raw cookie query parameter")
				}
				if got := cfg.Cookie["token"]; got != "token-value" {
					t.Fatalf("cookie token = %q, want token-value", got)
				}
				return stubCoreResponse(`{"status":1,"error_code":0,"data":{"info":[]}}`), nil
			}

			_, err = client.Search(context.Background(), SearchRequest{
				Keywords: "周杰伦",
				Type:     test.searchType,
				Cookie:   client.Cookie(),
				Extra:    map[string]any{"cookie": "token=token-value"},
			})
			if err != nil {
				t.Fatalf("Search() error = %v", err)
			}
		})
	}
}
