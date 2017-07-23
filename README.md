# go-amazon-product-api

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
	"github.com/charly3pins/go-amazon-product-api/amazon"
)

func main() {
	client := amazon.NewClient()
	res, err := client.ItemSearch("Books", "Clean Code")
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
In order to work correctly, you must set the following environment variables before execute the code:
```sh
export AWS_ACCESS_KEY_ID={YOUR_AWS_ACCESS_KEY_ID}
export AWS_SECRET_ACCESS_KEY={YOUR_AWS_SECRET_ACCESS_KEY}
export AWS_ASSOCIATE_TAG=collectus-21
export AWS_PRODUCT_REGION={YOUR_AWS_PRODUCT_REGION}
```
Finally, you can execute the example created:
```
go run example.go
```

## Author
[Charly3Pins](http://github.com/charly3pins)