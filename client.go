package amazon

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type ClientConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	AssociateTag    string
	Region          string
}

type Client interface {
	ItemSearch(string, string) (*ItemSearchResponse, error)
}

func NewClient(config ClientConfig) Client {
	return client{config.Region, NewAWSHTTPClient(config)}
}

type client struct {
	region string
	client HTTPClient
}

func (c client) ItemSearch(searchIndex, keywords string) (*ItemSearchResponse, error) {
	itemSearchResponse := ItemSearchResponse{}
	params := map[string]string{
		"SearchIndex":   searchIndex,
		"Keywords":      keywords,
		"ResponseGroup": "Images,ItemAttributes",
	}

	endpoint := fmt.Sprintf("http://%s/onca/xml", domains[c.region])
	parsedURL, _ := url.Parse(endpoint)
	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		log.Println("Error:", err.Error())
		return nil, err
	}

	res, err := c.client.Do(req, params)
	if data, _ := ioutil.ReadAll(res.Body); data != nil {
		if err = xml.Unmarshal(data, &itemSearchResponse); err != nil {
			log.Println("Error:", err.Error())
			return nil, err
		}
	}

	return &itemSearchResponse, err
}
