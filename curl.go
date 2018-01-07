package curl

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ----------------------
// OPTION HANDLUING
// ----------------------

func optD(v string) []string {
	if v == "" {
		return nil
	}
	split := strings.Split(v, "=")
	if len(split) != 2 {
		return nil
	}
	if strings.HasPrefix(split[1], `"`) {
		split[1] = strings.Trim(split[1], `"`)
	}
	return split
}

func optDataUrlencode(v string) []string {
	if v == "" {
		return nil
	}
	if string(v[0]) == "=" {
		// Unsupported
		return nil
	}
	split := strings.Split(v, "=")
	if len(split) == 2 {
		split[1] = url.QueryEscape(split[1])
		return split
	}
	// Everything else is unsupported
	return nil
}

// ----------------------
// REQUEST
// ----------------------

type request struct {
	method            string
	url               string
	contentType       string
	basicAuthUser     string
	basicAuthPassword string
	headers           map[string]string
	formData          map[string]string
}

func newRequest() *request {
	headers := make(map[string]string)
	formData := make(map[string]string)
	return &request{headers: headers, formData: formData}
}

func (r *request) addBasicAuth(composite string) error {
	pair := strings.Split(composite, ":")
	if len(pair) < 1 || len(pair) > 2 {
		return badRequestErr
	}
	r.basicAuthUser = pair[0]
	if len(pair) > 1 {
		r.basicAuthPassword = pair[1]
	}
	return nil
}

func (r *request) addHeader(composite string) error {
	if composite == "" {
		return badRequestErr
	}
	pos := strings.Index(composite, ":")
	if pos < 1 {
		return badRequestErr
	}
	key := composite[:pos]
	value := strings.TrimSpace(composite[pos+1:])
	if key == "" || value == "" {
		return badRequestErr
	}
	r.headers[key] = value
	return nil
}

func (r *request) addFormData(kv []string) error {
	if len(kv) != 2 {
		return badRequestErr
	}
	if r.contentType == "" {
		r.contentType = "application/x-www-form-urlencoded"
	}
	r.formData[kv[0]] = kv[1]
	return nil
}

func (r *request) httpRequest() (*http.Request, error) {
	var body io.Reader
	if len(r.formData) > 0 {
		body_str := ""
		for k, v := range r.formData {
			if body_str != "" {
				body_str = body_str + "&"
			}
			body_str = body_str + k + "=" + v
		}
		body = strings.NewReader(body_str)
	}
	req, err := http.NewRequest(r.method, r.url, body)
	if err != nil {
		return nil, err
	}
	if r.contentType != "" {
		req.Header.Add("Content-Type", r.contentType)
	}
	if len(r.headers) > 0 {
		for k, v := range r.headers {
			req.Header.Add(k, v)
		}
	}
	if r.basicAuthUser != "" {
		req.SetBasicAuth(r.basicAuthUser, r.basicAuthPassword)
	}
	return req, nil
}
