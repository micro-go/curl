package curl

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)


// func ResponseBody() performs the request and reads the response body.
func ResponseBody(req *http.Request, out interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	// XXX Default to JSON-formatted responses for now, but need to check the type.
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}
