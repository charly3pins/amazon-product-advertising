package amazon

import (
	"testing"
	"time"
)

const signedUrl = "http://webservices.amazon.es/onca/xml?AWSAccessKeyId=AWSK&AssociateTag=collectus-21&Keywords=Clean+Code&Operation=ItemSearch&ResponseGroup=Images%2CItemAttributes&SearchIndex=Books&Service=AWSECommerceService&Signature=%2FBe%2BqvyOpvdc2XA0YCxPYD3Wp1LKqWDBnk60k3%2BLPlM%3D&Timestamp=2017-07-26+23%3A00%3A00+%2B0900+Europe%2FSpain&Version=2013-08-01"

func TestHTTPClientSignUrl(t *testing.T) {
	client := Client{"AWSK", "AWSS", "collectus-21", "ES"}
	params := map[string]string{
		"SearchIndex":   "Books",
		"Keywords":      "Clean Code",
		"ResponseGroup": "Images,ItemAttributes",
	}
	httpClient := HTTPClient{client, params}

	urlValues := httpClient.buildUrlValues()
	urlValues.Set("Timestamp", time.Date(2017, time.July, 26, 23, 00, 0, 0, time.FixedZone("Europe/Spain", 9*60*60)).String())
	url, err := httpClient.signURL(urlValues)
	if err != nil {
		t.Errorf("Error signing URL with %v", urlValues)
	}

	if signedUrl != url {
		t.Errorf(`Expected "%v" but got "%v"`, signedUrl, url)
	}
}

func TestHTTPClientDo(t *testing.T) {
	client := Client{"AWSK", "AWSS", "collectus-21", "ES"}
	params := map[string]string{
		"SearchIndex":   "Books",
		"Keywords":      "Clean Code",
		"ResponseGroup": "Images,ItemAttributes",
	}
	httpClient := HTTPClient{client, params}

	res, err := httpClient.Do() // TODO MOCK THE CALL TO AMAZON
	if err != nil {
		t.Errorf("Expected nil but got %v", err)
	}

	if res.StatusCode != 200 {
		t.Errorf("expected res.StatusCode 200; got %v", res.Status)
	}
}
