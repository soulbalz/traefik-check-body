package checkbody

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/albrow/forms"
)

// BodyMatch demonstrates a BodyMatch plugin.
type BodyMatch struct {
	name     string
	next     http.Handler
	body     []SingleBody
	response ResponseError
}

func (a *BodyMatch) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	isBodyValid := true

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		a.response.Response(rw)
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	data, err := forms.Parse(r)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		a.response.Response(rw)
		return
	}

	for _, vBody := range a.body {
		isBodyValid = vBody.IsValid(data)
		if !isBodyValid {
			break
		}
	}

	if isBodyValid {
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		a.next.ServeHTTP(rw, r)
	} else {
		a.response.Response(rw)
		return
	}
}
