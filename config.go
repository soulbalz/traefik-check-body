package checkbodyplugin

import (
	"fmt"
	"net/http"
	"strings"
)

// Config the plugin configuration.
type Config struct {
	Body     []SingleBody
	Response ResponseError
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Body: []SingleBody{},
		Response: ResponseError{
			Code:    "400",
			Message: "Invalid Request.",
			Status:  http.StatusBadRequest,
		},
	}
}

//Validate validate config data
func (config *Config) Validate() error {
	if len(config.Body) == 0 {
		return fmt.Errorf("configuration incorrect, missing body")
	}

	for _, vBody := range config.Body {
		if strings.TrimSpace(vBody.Name) == "" {
			return fmt.Errorf("configuration incorrect, missing body name")
		}

		if len(vBody.Values) == 0 {
			return fmt.Errorf("configuration incorrect, missing body values")
		} else {
			for _, value := range vBody.Values {
				if strings.TrimSpace(value) == "" {
					return fmt.Errorf("configuration incorrect, empty value found")
				}
			}
		}

		if !vBody.IsContains() && vBody.MatchType == string(MatchAll) {
			return fmt.Errorf("configuration incorrect for body %v %s", vBody.Name, ", matchall can only be used in combination with 'contains'")
		}

		if strings.TrimSpace(vBody.MatchType) == "" {
			return fmt.Errorf("configuration incorrect, missing match type configuration for body %v", vBody.Name)
		}
	}

	return nil
}
