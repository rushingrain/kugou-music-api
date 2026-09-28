package sdk

import (
	"context"
	"fmt"
	"testing"

	corekugou "github.com/rushingrain/kugou-music-api/core/kugou"
)

func TestSearchComplexMatchesJSRequestAndDecodesTaggedResponse(t *testing.T) {
	client, err := New(WithCookie(map[string]string{
		"token":  "token-value",
		"userid": "123456",
		"dfid":   "dfid-value",
	}))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	client.requester = func(_ context.Context, cfg corekugou.RequestConfig) (corekugou.Response, error) {
		if got := fmt.Sprint(cfg.Params["keyword"]); got != "周杰伦" {
			t.Fatalf("keyword = %q, want %q", got, "周杰伦")
		}
		if _, ok := cfg.Params["keywords"]; ok {
			t.Fatal("unexpected keywords parameter")
		}
		if _, ok := cfg.Params["cookie"]; ok {
			t.Fatal("unexpected raw cookie parameter")
		}
		if got := cfg.Cookie["token"]; got != "token-value" {
			t.Fatalf("cookie token = %q, want token-value", got)
		}
		if got := fmt.Sprint(cfg.Params["page"]); got != "1" {
			t.Fatalf("page = %q, want 1", got)
		}
		if got := fmt.Sprint(cfg.Params["pagesize"]); got != "30" {
			t.Fatalf("pagesize = %q, want 30", got)
		}
		if got := fmt.Sprint(cfg.Params["filter"]); got != "6" {
			t.Fatalf("filter = %q, want 6 from Extra", got)
		}

		return stubCoreResponse("<!--KG_TAG_RES_START-->{\"status\":1,\"error_code\":0,\"data\":{\"lists\":[]}}<!--KG_TAG_RES_END-->"), nil
	}

	response, err := client.SearchComplex(context.Background(), SearchComplexRequest{
		Keywords: "周杰伦",
		Cookie:   client.Cookie(),
		Extra: map[string]any{
			"cookie": "token=token-value;userid=123456;dfid=dfid-value",
			"filter": 6,
		},
	})
	if err != nil {
		t.Fatalf("SearchComplex() error = %v", err)
	}
	if got := response.Body["status"]; got != float64(1) {
		t.Fatalf("response status = %v, want 1", got)
	}
	if got := response.Body["error_code"]; got != float64(0) {
		t.Fatalf("response error_code = %v, want 0", got)
	}
}
