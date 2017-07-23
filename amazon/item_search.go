package amazon

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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

type Client struct {
	AccessKeyID     string
	SecretAccessKey string
	AssociateTag    string
	Region          string
}

func NewClient() *Client {
	awsAccessKey := flag.String("awsAccessKey", os.Getenv("AWS_ACCESS_KEY_ID"), "aws acces key")
	awsSecretKey := flag.String("awsSecretKey", os.Getenv("AWS_SECRET_ACCESS_KEY"), "aws secret acces key")
	awsAssociateTag := flag.String("aswsAssociateTag", os.Getenv("AWS_ASSOCIATE_TAG"), "asws associate tag")
	awsRegion := flag.String("awsRegion", os.Getenv("AWS_PRODUCT_REGION"), "aws product region")
	flag.Parse()

	return &Client{*awsAccessKey, *awsSecretKey, *awsAssociateTag, *awsRegion}
}

func (c *Client) ItemSearch(searchIndex, keywords string) (*ItemSearchResponse, error) {
	params := map[string]string{
		"SearchIndex":   searchIndex,
		"Keywords":      keywords,
		"ResponseGroup": "Images,ItemAttributes",
	}
	url, err := c.buildURL(params)
	if err != nil {
		return nil, err
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	respObj := ItemSearchResponse{}
	if data, _ := ioutil.ReadAll(res.Body); data != nil {
		if err = xml.Unmarshal(data, &respObj); err != nil {
			return nil, err
		}
	}

	return &respObj, err
}

func (c *Client) buildURL(params map[string]string) (string, error) {
	endpoint := fmt.Sprintf("http://%s/onca/xml", domains[c.Region])
	parsedURL, _ := url.Parse(endpoint)
	queryValues := url.Values{}
	queryValues.Set("Service", "AWSECommerceService")
	queryValues.Set("AWSAccessKeyId", c.AccessKeyID)
	queryValues.Set("Version", "2013-08-01")
	queryValues.Set("Operation", "ItemSearch")
	queryValues.Set("AssociateTag", c.AssociateTag)
	queryValues.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))

	queryValues, err := c.signURL(params, queryValues, parsedURL)

	if err != nil {
		return "", err
	}
	parsedURL.RawQuery = queryValues.Encode()

	return parsedURL.String(), nil
}

func (c *Client) signURL(params map[string]string, queryValues url.Values, parsedURL *url.URL) (url.Values, error) {
	// ADD PARAMS TO THE QUERY
	for key, value := range params {
		queryValues.Set(key, value)
	}

	// ORDER ALL PARAMS IN THE QUERY BECAUSE AMAZON NEED IT IN ORDER
	queryKeys := make([]string, 0, len(queryValues))
	for key := range queryValues {
		queryKeys = append(queryKeys, key)
	}
	sort.Strings(queryKeys)

	// ESCAPE ALL ORDERED KEY/VALUES IN QUERY AND SIGN THEM
	queryKeysAndValues := make([]string, len(queryKeys))
	for i, key := range queryKeys {
		escapedKey := strings.Replace(url.QueryEscape(key), "+", "%20", -1)
		escapedValue := strings.Replace(url.QueryEscape(queryValues.Get(key)), "+", "%20", -1)
		queryKeysAndValues[i] = escapedKey + "=" + escapedValue
	}
	query := strings.Join(queryKeysAndValues, "&")

	msg := fmt.Sprintf("GET\n%s\n%s\n%s", parsedURL.Host, parsedURL.Path, query)
	hasher := hmac.New(sha256.New, []byte(c.SecretAccessKey))
	_, err := hasher.Write([]byte(msg))
	if err != nil {
		return url.Values{}, err
	}

	signature := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	queryValues.Set("Signature", signature)

	return queryValues, nil
}
