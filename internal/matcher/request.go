package matcher

import (
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

// Request is a matcher for http.Request.
type Request struct {
	values    url.Values
	lastError error
}

// RequestConfigurator is the matcher option.
type RequestConfigurator func(*Request)

// NewRequest builds a Request matcher.
func NewRequest(options ...RequestConfigurator) *Request {
	output := Request{}

	for _, opt := range options {
		opt(&output)
	}

	return &output
}

// RequestWithValue add a value in the matcher.
func RequestWithValue(key string, value string) RequestConfigurator {
	return func(request *Request) {
		if request.values == nil {
			request.values = url.Values{}
		}

		request.values.Add(key, value)
	}
}

// Matches implements the gomock.Matcher interface.
func (r *Request) Matches(x interface{}) bool {
	switch data := x.(type) {
	case http.Request:
		return r.Matches(&data)
	case *http.Request:
		originalValues, err := url.ParseQuery(data.URL.RawQuery)
		if err != nil {
			r.lastError = err

			return false
		}

		for key, values := range r.values {
			if !originalValues.Has(key) {
				r.lastError = fmt.Errorf("missing key %s", key)

				return false
			}

			for _, value := range values {
				if !slices.Contains(originalValues[key], value) {
					r.lastError = fmt.Errorf("missing value %s: %s - has [%s]", key, value, strings.Join(values, ","))

					return false
				}
			}
		}

		return true
	default:
		return false
	}
}

// String implements the gomock.Matcher interface.
func (r Request) String() string {
	output := []string{}

	for key, values := range r.values {
		output = append(output, fmt.Sprintf("%s:[%s]", key, strings.Join(values, ",")))
	}

	return fmt.Sprintf("%s: %s", r.lastError, strings.Join(output, "\n"))
}
