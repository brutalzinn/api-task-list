package request_entities

type OauthGenerateRequest struct {
	Callback        string `json:"callback"`
	ApplicationName string `json:"application_name"`
}
