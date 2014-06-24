//Provides helper functions to communicate with the goincremental stack apis
package sdk

import (
	"github.com/goincremental/sdk/security"
	"log"
	"net/http"
	"strings"
)

//PostJson will send a post message to the url specified after adding
//an expiry header and hmac signature to the request using the credentials given
func PostJson(url string, json *strings.Reader, key string, secret string) error {
	req, err := http.NewRequest("POST", url, json)
	if err != nil {
		log.Fatalf(err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("access-key", key)
	security.SignRequest(req, secret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
