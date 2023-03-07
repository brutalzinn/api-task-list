package scope

import (
	"strings"
)

type ScopePermission struct {
	Scopes []string
}

func New(scope string) (scopes ScopePermission) {
	scopeStrings := strings.Split(scope, ",")
	scopes = ScopePermission{
		Scopes: scopeStrings,
	}
	return scopes
}

func (s *ScopePermission) HasScope(scope []string) bool {
	valid := true
	for _, v := range s.Scopes {
		for _, t := range scope {
			if v == t {
				continue
			}
			valid = false
		}
	}

	return false

}
