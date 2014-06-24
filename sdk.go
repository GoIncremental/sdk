//Provides helper functions to communicate with the goincremental stack apis
package sdk

import (
	"github.com/goincremental/sdk/security"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func sendRequest(req *http.Request, key string, secret string) (resp *http.Response, err error) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("access-key", key)
	security.SignRequest(req, secret)

	client := &http.Client{}
	resp, err = client.Do(req)
	return
}

//PostJson will send a post message to the url specified after adding
//an expiry header and hmac signature to the request using the credentials given
func PostJson(url string, json *strings.Reader, key string, secret string) error {
	req, err := http.NewRequest("POST", url, json)
	if err != nil {
		log.Fatalf(err.Error())
	}

	resp, err := sendRequest(req, key, secret)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

//GetJson will send a get request to the url specified after adding
//an expiry header and hmac signature to the request using the credential given
func GetJson(url string, key string, secret string) (json []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf(err.Error())
	}

	resp, err := sendRequest(req, key, secret)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	// Read the content into a byte array
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// At this point we're done - simply return the bytes
	return body, nil

}
