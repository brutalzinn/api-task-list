package request_entities

type ApiKeyRequest struct {
	Name     string     `json:"name"`
	ExpireAt ExpireDays `json:"expire_at"`
}

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
