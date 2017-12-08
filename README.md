# go-amazon-product-api
[![GoDoc](https://godoc.org/github.com/charly3pins/go-amazon-product-api?status.svg)](https://godoc.org/github.com/charly3pins/go-amazon-product-api)
[![Build Status](https://travis-ci.org/charly3pins/go-amazon-product-api.svg?branch=master)](https://travis-ci.org/charly3pins/go-amazon-product-api)

Go Client Library for [Amazon Product API](https://affiliate-program.amazon.com/gp/advertising/api/detail/main.html)

## Usage
First of all you need to download the library in your [$GOPATH](https://golang.org/doc/code.html#GOPATH) using the following command:
```sh
go get -u github.com/charly3pins/go-amazon-product-api/amazon
```
Then, you can create a simple example.go and call the NewClient() constructor and the ItemSearch(searchIndex, keywords string) method to search by text:
```go
package main

import (
	"fmt"
	amazon "github.com/charly3pins/go-amazon-product-api"
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

## Author
[Charly3Pins](http://github.com/charly3pins)
