package amazon

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
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
)

type AmazonRequest interface {
	Do() (*http.Response, error)
}

type HTTPClient struct {
	Client     Client
	Parameters map[string]string
}

func (h HttpClient) Do() (*http.Response, error) {
	urlValues := h.buildUrlValues()
	url, err := h.signURL(urlValues)
	if err != nil {
		return nil, err
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h HTTPClient) buildUrlValues() url.Values {
	urlValues := url.Values{}
	urlValues.Set("Service", "AWSECommerceService")
	urlValues.Set("AWSAccessKeyId", h.Client.AccessKeyID)
	urlValues.Set("Version", "2013-08-01")
	urlValues.Set("Operation", "ItemSearch")
	urlValues.Set("AssociateTag", h.Client.AssociateTag)
	urlValues.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	return urlValues
}

func (h HTTPClient) signURL(urlValues url.Values) (string, error) {
	endpoint := fmt.Sprintf("http://%s/onca/xml", domains[h.Client.Region])
	parsedURL, _ := url.Parse(endpoint)

	// ADD PARAMS TO THE QUERY
	for key, value := range h.Parameters {
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
	hasher := hmac.New(sha256.New, []byte(h.Client.SecretAccessKey))
	_, err := hasher.Write([]byte(msg))
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	urlValues.Set("Signature", signature)

	parsedURL.RawQuery = urlValues.Encode()

	return parsedURL.String(), nil
}
