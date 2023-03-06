package hypermedia

import (
	"context"
	"fmt"
	"net/http"
)

func (hypermedia *HyperMedia) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(r.Context(), "links", hypermedia.Links)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Handler(options HyperMediaOptions) func(next http.Handler) http.Handler {
	c := New(options)
	return c.Handler
}
func New(options HyperMediaOptions) *HyperMedia {
	opt := &HyperMedia{
		Links: options.Links,
	}
	return opt
}

type HypermediaLink struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
	Type string `json:"type"`
}

type HyperMedia struct {
	Links []HypermediaLink `json:"links"`
}

type HyperMediaOptions struct {
	AutoId bool
	Links  []HypermediaLink
}

func CreateHyperMedia(rel string, href string, request_type string) HypermediaLink {
	hypermedia := HypermediaLink{
		Rel:  rel,
		Href: href,
		Type: request_type,
	}
	return hypermedia
}

func CreateHyperMediaLinksFor(ID int64, context context.Context) []HypermediaLink {
	links, _ := context.Value("links").([]HypermediaLink)
	var hypermediaLink []HypermediaLink
	for _, link := range links {
		link.Href = fmt.Sprintf(link.Href, ID)
		hypermediaLink = append(hypermediaLink, link)
	}
	return hypermediaLink
}
