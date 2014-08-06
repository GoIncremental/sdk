//Provides helper functions to communicate with the goincremental stack apis
package sdk

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/goincremental/sdk/security"
)

var (
	ErrNotAuthorised         = errors.New("sdk: request not authorised")
	ErrSDKUnexpectedResponse = errors.New("sdk: unexpected response code")
)

type sdkClient struct {
	host   string
	key    string
	secret string
}

type Client interface {
	Post(string, []byte) error
	Get(string) ([]byte, error)
}

func NewClient(host string, key string, secret string) Client {
	return &sdkClient{
		host:   host,
		key:    key,
		secret: secret,
	}
}

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
func (c *sdkClient) Post(url string, data []byte) error {
	json := strings.NewReader(string(data))
	req, err := http.NewRequest("POST", url, json)
	if err != nil {
		log.Fatalf(err.Error())
	}

	resp, err := sendRequest(req, c.key, c.secret)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

//GetJson will send a get request to the url specified after adding
//an expiry header and hmac signature to the request using the credential given
func (c *sdkClient) Get(url string) (json []byte, err error) {
	address := strings.Join([]string{c.host, url}, "")
	log.Printf("getting Json %s\n", address)
	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		log.Fatalf(err.Error())
	}

	resp, err := sendRequest(req, c.key, c.secret)
	if err != nil {
		log.Printf("error: %s\n", err)
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 401:
		err = ErrNotAuthorised
		return nil, err
	case 200:
		log.Println("get json 200")
		// Read the content into a byte array
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		// At this point we're done - simply return the bytes
		return body, nil
	default:
		log.Printf("responseCode: %d\n", resp.StatusCode)
		return nil, ErrSDKUnexpectedResponse
	}
}
