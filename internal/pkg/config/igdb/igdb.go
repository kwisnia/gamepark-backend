package igdb

import (
	"fmt"
	"github.com/kwisnia/igdb/v2"
	"github.com/kwisnia/inzynierka-backend/internal/pkg/config"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var IgdbClient *igdb.Client

func SetupClient() {
	IgdbClient = igdb.NewClient(os.Getenv("IGDB_CLIENT_ID"), os.Getenv("IGDB_ACCESS_TOKEN"), nil)
}

func SetupWebhooks() {
	apiUrl := "https://api.igdb.com/v4"
	resource := "/games/webhooks"
	data := url.Values{}
	data.Set("url", "https://api.gamepark.space/games/webhooks")
	data.Set("method", "create")
	data.Set("secret", config.GetEnv("WEBHOOK_SECRET"))

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Authorization", "Bearer "+config.GetEnv("IGDB_ACCESS_TOKEN"))
	r.Header.Add("Client-ID", config.GetEnv("IGDB_CLIENT_ID"))
	_, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	fmt.Println("Webhook create created")

	data.Set("method", "update")
	r, _ = http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Authorization", "Bearer "+config.GetEnv("IGDB_ACCESS_TOKEN"))
	r.Header.Add("Client-ID", config.GetEnv("IGDB_CLIENT_ID"))
	_, err = client.Do(r)
	if err != nil {
		panic(err)
	}
	fmt.Println("Webhook update created")

	data.Set("method", "delete")
	r, _ = http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Authorization", "Bearer "+config.GetEnv("IGDB_ACCESS_TOKEN"))
	r.Header.Add("Client-ID", config.GetEnv("IGDB_CLIENT_ID"))
	_, err = client.Do(r)
	if err != nil {
		panic(err)
	}
	fmt.Println("Webhook delete created")
}
