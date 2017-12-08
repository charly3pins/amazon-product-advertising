package amazon

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type ClientConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	AssociateTag    string
	Region          string
	AWSEndpoint     string
}

type Criteria struct {
	SearchIndex string
	Keywords    string
}

type Client interface {
	ItemSearch(Criteria) (*ItemSearchResponse, error)
}

var ErrBadStatusCode = fmt.Errorf("wrong status code")

func NewClient(config ClientConfig) Client {
	if (config.Region != "" || config.AWSEndpoint != "") && config.AccessKeyID != "" && config.SecretAccessKey != "" {
		return client{config, NewAWSHTTPClient(config)}
	}
	return client{config, defaultHTTPClient}
}

type client struct {
	config ClientConfig
	client HTTPClient
}

func (c client) ItemSearch(criteria Criteria) (*ItemSearchResponse, error) {
	itemSearchResponse := ItemSearchResponse{}
	params := map[string]string{
		"SearchIndex":   criteria.SearchIndex,
		"Keywords":      criteria.Keywords,
		"ResponseGroup": "Images,ItemAttributes",
	}

	if c.config.AWSEndpoint == "" {
		c.config.AWSEndpoint = GetEndpoint(c.config.Region)
	}
	parsedURL, err := url.Parse(c.config.AWSEndpoint)
	if err != nil {
		log.Println("Error obtaining parsed URL:", err)
		return nil, err
	}
	signedURL, err := c.signURL(parsedURL, params)
	if err != nil {
		log.Println("Error signing URL:", err)
		return nil, err
	}
	req, err := http.NewRequest("GET", signedURL, nil)
	if err != nil {
		log.Println("Error making new request:", err.Error())
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, ErrBadStatusCode
	}
	if data, _ := ioutil.ReadAll(res.Body); data != nil {
		if len(data) > 0 {
			if err = xml.Unmarshal(data, &itemSearchResponse); err != nil {
				log.Println("Error reading body:", err.Error())
				return nil, err
			}
		}
	}

	return &itemSearchResponse, err
}

func (c client) signURL(parsedURL *url.URL, params map[string]string) (string, error) {
	urlValues := c.buildUrlValues()
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
	hasher := hmac.New(sha256.New, []byte(c.config.SecretAccessKey))
	_, err := hasher.Write([]byte(msg))
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	urlValues.Set("Signature", signature)

	parsedURL.RawQuery = urlValues.Encode()

	return parsedURL.String(), nil
}

func (c client) buildUrlValues() url.Values {
	urlValues := url.Values{}
	urlValues.Set("Service", "AWSECommerceService")
	urlValues.Set("AWSAccessKeyId", c.config.AccessKeyID)
	urlValues.Set("Version", "2013-08-01")
	urlValues.Set("Operation", "ItemSearch")
	urlValues.Set("AssociateTag", c.config.AssociateTag)
	urlValues.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	return urlValues
}
