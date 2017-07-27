package amazon

import (
	"encoding/xml"
	"io/ioutil"
)

type Client struct {
	AccessKeyID     string
	SecretAccessKey string
	AssociateTag    string
	Region          string
}

func (c Client) ItemSearch(searchIndex, keywords string) (*ItemSearchResponse, error) {
	itemSearchResponse := ItemSearchResponse{}
	params := map[string]string{
		"SearchIndex":   searchIndex,
		"Keywords":      keywords,
		"ResponseGroup": "Images,ItemAttributes",
	}
	httpClient := HTTPClient{c, params}
	res, err := httpClient.Do()
	if data, _ := ioutil.ReadAll(res.Body); data != nil {
		if err = xml.Unmarshal(data, &itemSearchResponse); err != nil {
			return nil, err
		}
	}

	return &itemSearchResponse, err
}
