package traefik_check_body

import (
	"context"
	"fmt"
	"net/http"
)

// New created a new BodyMatch plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	err := config.Validate()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return &BodyMatch{
		name:     name,
		next:     next,
		body:     config.Body,
		response: config.Response,
	}, nil

}
