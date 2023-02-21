package entities

type ExpireDays int64

const (
	SevenDays ExpireDays = iota
	ThirtyDays
	NinetyDays
	YearDays
)

func (s ExpireDays) Days() int {
	switch s {
	case SevenDays:
		return 7
	case ThirtyDays:
		return 30
	case NinetyDays:
		return 90
	case YearDays:
		return 365
	}
	return 0
}

type ApiKey struct {
	ID             int64
	ApiKey         string
	Scopes         string
	Name           string
	NameNormalized string
	UserId         int64
	CreateAt       *string
	UpdateAt       *string
	ExpireAt       string
}

type ApiKeyRequest struct {
	Name     string
	ExpireAt ExpireDays
}

type ApiKeyResponse struct {
	Name     string
	ExpireAt string
	CreateAt *string
	Scopes   string
	UpdateAt *string
}
