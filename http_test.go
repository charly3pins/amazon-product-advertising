package amazon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAwsHTTPClient_Do(t *testing.T) {
	server := httptest.NewServer(&myHandler{func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			t.Error("unexpected method:", req.Method)
			return
		}
		fmt.Fprintf(w,
			`{
				"ItemSearchResponse": {
					"Items": {
						"Request": {
							"IsValid": true
						}
					},
					"TotalResult": 122879,
					"TotalPages": 12288,
					"MoreSearchResultsUrl": "url",
					"Item": {
						"ASIN": "ASIN_CODE",
						"DetailPageURL": "Detail_Page_Url",
						"ItemLinks": [{
							"ItemLink": {
								"Description": "Desc",
								"URL": "URL"
							}
						}]
					},
					"SmallImage": {
						"URL": "url_small_img",
						"Height": {
							"Value": 75,
							"Units": "pixels"
						}
					},
					"MediumImage": {
						"URL": "url_medium_img",
						"Height": {
							"Value": 160,
							"Units": "pixels"
						}
					},
					"LargImage": {
						"URL": "url_large_img",
						"Height": {
							"Value": 500,
							"Units": "pixels"
						}
					},
					"ItemAttributes": {
						"Binding": "Toy",
						"Brand": "Best Brand",
						"EAN": "EAN_NUMBER",
						"Label": "The label ",
						"Title": "My fantastic object",
						"ListPrice": {
							"Amount": "100",
							"Currency": "EUR",
							"FormattedPrice": "EUR 9,88"
						}
					}
				}
			}`)
	}})
	defer server.Close()

	config := ClientConfig{
		AccessKeyID:     "some_ID",
		SecretAccessKey: "some_secret",
		AssociateTag:    "some_tag",
		Region:          "ES",
		AWSEndpoint:     server.URL,
	}

	c := awsHTTPClient{config}
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Error(err)
		return
	}
	resp, err := c.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	data := ItemSearchResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Error("decoding the resp body:", err.Error())
		return
	}
	defer resp.Body.Close()
}
