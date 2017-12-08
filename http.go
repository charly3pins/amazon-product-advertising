package amazon

import (
	"fmt"
	"net/http"
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

func GetEndpoint(region string) string {
	return fmt.Sprintf("http://%s/onca/xml", domains[region])
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type awsHTTPClient struct {
	Config ClientConfig
}

func NewAWSHTTPClient(config ClientConfig) HTTPClient {
	return awsHTTPClient{config}
}

func (a awsHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return defaultHTTPClient.Do(req)
}
