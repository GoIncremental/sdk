//Implementation of security methods for the goincremental API
package security

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func addExpiry(req *http.Request) error {
	expiry := strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10)
	req.Header.Add("expiry-date", expiry)
	return nil
}

func stringToSign(req *http.Request) string {
	var buffer bytes.Buffer
	buffer.WriteString(req.Method)
	buffer.WriteString("\n")
	buffer.WriteString(req.URL.Host)
	buffer.WriteString("\n")
	buffer.WriteString(req.URL.Path)
	buffer.WriteString("\n")
	if accessKey := req.Header.Get("access-key"); accessKey != "" {
		buffer.WriteString("access-key=")
		buffer.WriteString(url.QueryEscape(accessKey))
	}
	if expiryDate := req.Header.Get("expiry-date"); expiryDate != "" {
		buffer.WriteString("&expiry-date=")
		buffer.WriteString(url.QueryEscape(expiryDate))
	}

	return buffer.String()
}

func signString(key string, payload string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(payload))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

//SignRequest takes an http request that is ready to send to the api and
//adds expiry-date and signature headers to the request
func SignRequest(req *http.Request, secret string) error {
	addExpiry(req)
	signature := signString(secret, stringToSign(req))
	req.Header.Add("signature", signature)
	return nil
}
