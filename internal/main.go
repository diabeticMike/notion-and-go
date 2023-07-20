package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dstotijn/go-notion"
)

const SECRET = "internal_integration_secret"

func main() {
	client := notion.NewClient(SECRET)

	db, err := client.FindPageByID(context.Background(), "my_page_id")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(db.URL)
}
