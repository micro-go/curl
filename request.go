package curl

import (
	"github.com/micro-go/parse"
	"net/http"
)

// func Request answers a new HTTP request from the arguments.
func Request(args ...string) (*http.Request, error) {
	req := newRequest()
	var err error
	token := parse.NewStringToken(args...)
	for a, err := token.Next(); err == nil; a, err = token.Next() {
		// see https://curl.haxx.se/docs/manpage.html for the rules
		switch a {
		case "-d":
			req.method = "POST"
			n, _ := token.Next()
			mergeError(err, req.addFormData(optD(n)))
		case "--data-urlencode":
			n, _ := token.Next()
			mergeError(err, req.addFormData(optDataUrlencode(n)))
		case "-G", "--get":
			req.method = "GET"
		case "-H":
			n, _ := token.Next()
			mergeError(err, req.addHeader(n))
		case "-u", "--user", "--basic":
			n, _ := token.Next()
			mergeError(err, req.addBasicAuth(n))
		case "-X":
			req.method = a
		default:
			// Options start with dashes. The non-option is the URL.
			if a == "" || string(a[0]) == "-" {
				return nil, badRequestErr
			}
			req.url = a
		}
	}
	if err != nil {
		return nil, err
	}
	return req.httpRequest()
}
