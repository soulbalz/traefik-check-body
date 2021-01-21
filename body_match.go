package checkbodyplugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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

//ResponseError contains a failuer message
type ResponseError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
	Raw     string `json:"raw,omitempty"`
}

// Config the plugin configuration.
type Config struct {
	Body     []SingleBody
	Response ResponseError
}

// MatchType defines an enum which can be used to specify the match type for the 'contains' config.
type MatchType string

const (
	//MatchAll requires all values to be matched
	MatchAll MatchType = "all"
	//MatchOne requires only one value to be matched
	MatchOne MatchType = "one"
)

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

// New created a new BodyMatch plugin.
func New(ctx context.Context, next http.Handler, config *Config, _ string) (http.Handler, error) {
	if len(config.Body) == 0 {
		return nil, fmt.Errorf("configuration incorrect, missing body")
	}

	for _, vBody := range config.Body {
		if strings.TrimSpace(vBody.Name) == "" {
			return nil, fmt.Errorf("configuration incorrect, missing body name")
		}

		if len(vBody.Values) == 0 {
			return nil, fmt.Errorf("configuration incorrect, missing body values")
		} else {
			for _, value := range vBody.Values {
				if strings.TrimSpace(value) == "" {
					return nil, fmt.Errorf("configuration incorrect, empty value found")
				}
			}
		}

		if !vBody.IsContains() && vBody.MatchType == string(MatchAll) {
			return nil, fmt.Errorf("configuration incorrect for body %v %s", vBody.Name, ", matchall can only be used in combination with 'contains'")
		}

		if strings.TrimSpace(vBody.MatchType) == "" {
			return nil, fmt.Errorf("configuration incorrect, missing match type configuration for body %v", vBody.Name)
		}
	}

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		bodyValid := true

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(rw, config.Response.getMessage(), config.Response.Status)
			return
		}

		var reqBody map[string]string
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		err = json.NewDecoder(r.Body).Decode(&reqBody)

		if err == nil {
			for _, vBody := range config.Body {
				reqBodyVal := reqBody[vBody.Name]

				if vBody.IsContains() && reqBodyVal != "" {
					bodyValid = checkContains(&reqBodyVal, &vBody)
				} else {
					bodyValid = checkRequired(&reqBodyVal, &vBody)
				}

				if !bodyValid {
					break
				}
			}
		} else {
			r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			r.ParseForm() // Parses the request body

			for _, vBody := range config.Body {
				reqBodyVal := r.Form.Get(vBody.Name)
				if vBody.IsContains() && reqBodyVal != "" {
					bodyValid = checkContains(&reqBodyVal, &vBody)
				} else {
					bodyValid = checkRequired(&reqBodyVal, &vBody)
				}

				if !bodyValid {
					break
				}
			}
		}

		if bodyValid {
			r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			next.ServeHTTP(rw, r)
		} else {
			http.Error(rw, config.Response.getMessage(), config.Response.Status)
		}
	}), nil
}

//ResponseError get message
func (re *ResponseError) getMessage() string {
	var s string
	if re.Raw == "" {
		s = fmt.Sprintf(`{
			"data": null,
			"error": {
				"code": "%s",
				"message": "%s"
			}
		}`, re.Code, re.Message)
	} else {
		s = re.Raw
	}
	return s
}

func checkContains(requestValue *string, vBody *SingleBody) bool {
	matchCount := 0
	for _, value := range vBody.Values {
		if strings.Contains(*requestValue, value) {
			matchCount++
		}
	}

	if matchCount == 0 {
		return false
	} else if vBody.MatchType == string(MatchAll) && matchCount != len(vBody.Values) {
		return false
	}

	return true
}

func checkRequired(requestValue *string, vBody *SingleBody) bool {
	matchCount := 0
	for _, value := range vBody.Values {
		if *requestValue == value {
			matchCount++
		}

		if !vBody.IsRequired() && *requestValue == "" {
			matchCount++
		}
	}

	if matchCount == 0 {
		return false
	}

	return true
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
