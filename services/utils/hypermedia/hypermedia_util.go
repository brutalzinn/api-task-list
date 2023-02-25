package hypermedia_util

func CreateHyperMedia(links map[string]any, rel string, href string, request_type string) map[string]any {
	links[rel] = map[string]any{
		"href": href,
		"type": request_type,
	}
	return links
}
