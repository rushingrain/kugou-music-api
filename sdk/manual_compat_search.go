// Search compatibility wrappers live here because both the plain and the
// complex search endpoints need hand-assembled request maps that the generated
// wrappers cannot express. Keep them next to the generated search models in
// generated_api_search_sheet.go.
package sdk

import (
	"context"
	"fmt"
)

// Search mirrors the upstream JS client: the requested type selects both the
// endpoint version and the path, so song goes to /v3 while the other types go
// to /v1. The type is a path segment, never a query parameter.
func (c *Client) Search(ctx context.Context, req SearchRequest) (*SearchResponse, error) {
	params := structToMap(req)
	delete(params, "Cookie")
	delete(params, "Extra")
	if compat, ok := buildCompatParams("search", params, req.Cookie); ok {
		params = compat
	}
	for k, v := range req.Extra {
		if k == "cookie" {
			continue
		}
		params[k] = v
	}
	searchType := compatFirstAnyString(params["type"], "song")
	if !isSearchType(searchType) {
		searchType = "song"
	}
	delete(params, "type")
	version := "v1"
	if searchType == "song" {
		version = "v3"
	}
	cookie := applyCompatCookie("search", req.Cookie)
	resp, err := c.Call(ctx, RouteSearch, Request{
		URL:    fmt.Sprintf("/%s/search/%s", version, searchType),
		Params: params,
		Cookie: cookie,
	})
	if err != nil {
		return nil, err
	}
	out := SearchResponse(*resp)
	return &out, nil
}

// isSearchType reports whether the caller asked for one of the search endpoints
// the upstream JS client exposes.
func isSearchType(searchType string) bool {
	switch searchType {
	case "song", "album", "author", "mv", "lyric", "special":
		return true
	default:
		return false
	}
}

// SearchComplex mirrors the upstream JS client: the endpoint expects the
// singular keyword parameter and its own fixed page defaults. Authentication
// stays in Cookie, because forwarding a raw cookie query parameter changes the
// signed request and the upstream rejects it.
func (c *Client) SearchComplex(ctx context.Context, req SearchComplexRequest) (*SearchComplexResponse, error) {
	params := map[string]any{
		"platform": "AndroidFilter",
		"keyword":  compatFirstAnyString(req.Keywords, ""),
		"page":     firstNonZero(req.Page, 1),
		"pagesize": firstNonZero(req.Pagesize, 30),
		"cursor":   0,
	}
	for k, v := range req.Extra {
		if k == "cookie" {
			continue
		}
		params[k] = v
	}
	resp, err := c.Call(ctx, RouteSearchComplex, Request{
		Params:      params,
		Cookie:      req.Cookie,
		EncryptType: "android",
	})
	if err != nil {
		return nil, err
	}
	out := SearchComplexResponse(*resp)
	return &out, nil
}
