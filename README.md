# amazon-product-advertising
[![GoDoc](https://godoc.org/github.com/charly3pins/amazon-product-advertising?status.svg)](https://godoc.org/github.com/charly3pins/amazon-product-advertising)
[![Build Status](https://travis-ci.org/charly3pins/amazon-product-advertising.svg?branch=master)](https://travis-ci.org/charly3pins/amazon-product-advertising)

Go Client Library for [Amazon Product Advertising API](https://affiliate-program.amazon.com/gp/advertising/api/detail/main.html)

## Usage
First of all you need to download the library in your [$GOPATH](https://golang.org/doc/code.html#GOPATH) using the following command:
```sh
go get -u github.com/charly3pins/amazon-product-advertising
```
Then, you can create a simple example.go and call the NewClient(config) constructor passing the ClientConfig with your credentials and call the ItemSearch(criteria) method with your creteria, selecting the SearchIndex and the Keywords, to search by text:
```go
package main

import (
	"fmt"
	amazon "github.com/charly3pins/amazon-product-advertising"
)

func main() {
	criteria := amazon.Criteria{
		SearchIndex: "Books",
		Keywords:    "Clean Code",
	}
	config := amazon.ClientConfig{
		AccessKeyID:     "{YOUR_AWS_ACCESS_KEY_ID}",
		SecretAccessKey: "{YOUR_AWS_SECRET_ACCESS_KEY}",
		AssociateTag:    "collectus-21",
		Region:          "{YOUR_AWS_PRODUCT_REGION}",
		AWSEndpoint:     amazon.GetEndpoint("{YOUR_AWS_PRODUCT_REGION}"),
	}
	client := amazon.NewClient(config)
	res, err := client.ItemSearch(criteria)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%d results found\n\n", res.Items.TotalResults)
	for i, item := range res.Items.Item {
		fmt.Printf("Result: %d\nTitle: %v\nAuthor: %v\nBinding: %v\nLargeImage: %v\nURL: %v\n\n", i, item.ItemAttributes.Title, item.ItemAttributes.Author, item.ItemAttributes.Binding, item.ImageSets.ImageSet[0].LargeImage, item.DetailPageURL)
	}
}
```

Finally, you can execute the example created:
```
go run example.go
```
