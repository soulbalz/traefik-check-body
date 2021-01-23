package checkbodyplugin

import (
	"strings"

	"github.com/albrow/forms"
)

//SingleBody contains a single body keypair
type SingleBody struct {
	Name      string   `json:"name,omitempty"`
	Values    []string `json:"values,omitempty"`
	MatchType string   `json:"matchtype,omitempty"`
	Required  *bool    `json:"required,omitempty"`
	Contains  *bool    `json:"contains,omitempty"`
	URLDecode *bool    `json:"urldecode,omitempty"`
}

//IsValid checks is validation valid
func (s *SingleBody) IsValid(data *forms.Data) bool {
	isValid := true

	reqVal := data.Get(s.Name)
	if s.IsContains() && reqVal != "" {
		isValid = s.checkContains(&reqVal)
	} else {
		isValid = s.checkRequired(&reqVal)
	}
	return isValid
}

//IsContains checks whether a body value should contain the configured value
func (s *SingleBody) IsContains() bool {
	if s.Contains == nil || *s.Contains == false {
		return false
	}
	return true
}

//IsRequired checks whether a body is mandatory in the request, defaults to 'true'
func (s *SingleBody) IsRequired() bool {
	if s.Required == nil || *s.Required != false {
		return true
	}
	return false
}

func (s *SingleBody) checkContains(reqValue *string) bool {
	matchCount := 0
	for _, value := range s.Values {
		if strings.Contains(*reqValue, value) {
			matchCount++
		}
	}

	if matchCount == 0 {
		return false
	} else if s.MatchType == string(MatchAll) && matchCount != len(s.Values) {
		return false
	}
	return true
}

func (s *SingleBody) checkRequired(reqValue *string) bool {
	matchCount := 0
	for _, value := range s.Values {
		if *reqValue == value {
			matchCount++
		}

		if !s.IsRequired() && *reqValue == "" {
			matchCount++
		}
	}

	if matchCount == 0 {
		return false
	}
	return true
}
