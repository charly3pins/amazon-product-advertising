package amazon

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

var (
	domains = map[string]string{
		"CA": "ecs.amazonaws.ca",
		"CN": "webservices.amazon.cn",
		"DE": "ecs.amazonaws.de",
		"ES": "webservices.amazon.es",
		"FR": "ecs.amazonaws.fr",
		"IT": "webservices.amazon.it",
		"JP": "ecs.amazonaws.jp",
		"UK": "ecs.amazonaws.co.uk",
		"US": "ecs.amazonaws.com",
	}
	defaultHTTPClient = http.DefaultClient
)

type HTTPClient interface {
	Do(*http.Request, map[string]string) (*http.Response, error)
}

type awsHTTPClient struct {
	Config ClientConfig
}

func NewAWSHTTPClient(config ClientConfig) HTTPClient {
	return awsHTTPClient{config}
}

func (a awsHTTPClient) Do(req *http.Request, params map[string]string) (*http.Response, error) {
	url, err := a.signURL(req.URL, params)
	if err != nil {
		log.Println("Error:", err.Error())
		return nil, err
	}

	res, err := defaultHTTPClient.Get(url)
	if err != nil {
		log.Println("Error:", err.Error())
		return nil, err
	}

	return res, nil
}

func (a awsHTTPClient) signURL(parsedURL *url.URL, params map[string]string) (string, error) {
	urlValues := a.buildUrlValues()
	// ADD PARAMS TO THE QUERY
	for key, value := range params {
		urlValues.Set(key, value)
	}

	// ORDER ALL PARAMS IN THE QUERY BECAUSE AMAZON NEED IT IN ORDER
	queryKeys := make([]string, 0, len(urlValues))
	for key := range urlValues {
		queryKeys = append(queryKeys, key)
	}
	sort.Strings(queryKeys)

	// ESCAPE ALL ORDERED KEY/VALUES IN QUERY AND SIGN THEM
	queryKeysAndValues := make([]string, len(queryKeys))
	for i, key := range queryKeys {
		escapedKey := strings.Replace(url.QueryEscape(key), "+", "%20", -1)
		escapedValue := strings.Replace(url.QueryEscape(urlValues.Get(key)), "+", "%20", -1)
		queryKeysAndValues[i] = escapedKey + "=" + escapedValue
	}
	query := strings.Join(queryKeysAndValues, "&")

	msg := fmt.Sprintf("GET\n%s\n%s\n%s", parsedURL.Host, parsedURL.Path, query)
	hasher := hmac.New(sha256.New, []byte(a.Config.SecretAccessKey))
	_, err := hasher.Write([]byte(msg))
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	urlValues.Set("Signature", signature)

	parsedURL.RawQuery = urlValues.Encode()

	return parsedURL.String(), nil
}

func (a awsHTTPClient) buildUrlValues() url.Values {
	urlValues := url.Values{}
	urlValues.Set("Service", "AWSECommerceService")
	urlValues.Set("AWSAccessKeyId", a.Config.AccessKeyID)
	urlValues.Set("Version", "2013-08-01")
	urlValues.Set("Operation", "ItemSearch")
	urlValues.Set("AssociateTag", a.Config.AssociateTag)
	urlValues.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	return urlValues
}
