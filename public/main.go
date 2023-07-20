package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/gin-gonic/gin"
)

const (
	oauthClientID     = "client_id"
	oauthClientSecret = "client_secret"
	redirectURL       = "http://localhost:8001/here"
)

type OAuthAccessToken struct {
	AccessToken   string `json:"access_token,omitempty"`
	WorkspaceName string `json:"workspace_name,omitempty"`
	WorkspaceIcon string `json:"workspace_icon,omitempty"`
	BotID         string `json:"bot_id,omitempty"`
}

func main() {
	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 5,
	}

	var token string
	r := gin.Default()

	r.GET("/here", func(c *gin.Context) {
		code := c.Query("code")
		b, err := json.Marshal(&struct {
			GrantType   string `json:"grant_type,omitempty"`
			Code        string `json:"code,omitempty"`
			RedirectURI string `json:"redirect_uri,omitempty"`
		}{
			GrantType:   "authorization_code",
			Code:        code,
			RedirectURI: redirectURL,
		})
		if err != nil {
			log.Fatal(err)
			return
		}
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://api.notion.com/v1/oauth/token", bytes.NewReader(b))
		if err != nil {
			log.Fatal(err)
			return
		}
		req.SetBasicAuth(oauthClientID, oauthClientSecret)
		req.Header.Add("Content-Type", "application/json")

		rsp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
			return
		}

		defer rsp.Body.Close()

		var body OAuthAccessToken
		if err = json.NewDecoder(rsp.Body).Decode(&body); err != nil {
			fmt.Println(err)
			return
		}
		token = body.AccessToken
		notionClient := notion.NewClient(token)

		db, err := notionClient.FindPageByID(context.Background(), "my_page_id")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(db.URL)
	})

	r.Run(":8001")
}
