package p

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func PostNotionWebhook(w http.ResponseWriter, r *http.Request) {
	// NotionのWebhookが送信するJSONの形式が分からないので、仮組み
	var data struct {
		Team      string `json:"team"`
		NotionUrl string `json:"url"`
		Text      string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		switch err {
		case io.EOF:
			fmt.Fprint(w, "Hello World!")
			return
		default:
			log.Printf("json.NewDecoder: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	// 必須パラメータのチェック
	if data.Team == "" {
		http.Error(w, "missing team", http.StatusBadRequest)
		return
	}

	if data.NotionUrl == "" {
		http.Error(w, "missing url", http.StatusBadRequest)
		return
	}

	if data.Text == "" {
		http.Error(w, "missing text", http.StatusBadRequest)
		return
	}

	// DiscordのWebhookに通知
	if err := postToDiscordWebhook(data.Team); err != nil {
		log.Printf("postToDiscordWebhook: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else {
		fmt.Fprint(w, "Success")
	}
}

func postToDiscordWebhook(team string) error {
	type embeds struct {
		Title       string `json:"title"`
		Url         string `json:"url"`
		Description string `json:"description"`
		Color       int    `json:"color"`
	}

	type postData struct {
		Content string `json:"content"`
		Embeds  embeds `json:"embeds"`
	}

	var data = postData{}
	data.Content = "Notionに新しい投稿があります！"
	data.Embeds = embeds{}
	data.Embeds.Title = "page.title"
	data.Embeds.Color = 5620992
	// 環境変数からDiscordのWebhook URLを取得

	var webhookUrl string
	switch team {
	case "teamA":
		webhookUrl = "Team A Webhook URL"
	case "teamB":
		webhookUrl = "Team B Webhook URL"
	case "teamC":
		webhookUrl = "Team C Webhook URL"
	case "teamD":
		webhookUrl = "Team D Webhook URL"
	case "teamE":
		webhookUrl = "Team E Webhook URL"
	default:
		return fmt.Errorf("invalid team: %s", team)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", webhookUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Content-Typeを設定
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}
