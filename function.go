package p

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
)

// NotionData NotionのWebhookが送信するJSONデータの構造体
type NotionData struct {
	Team   string
	Url    string
	Title  string
	Status string
	User   string
}

// postErrMsgToDiscordWebhook はエラーメッセージをDiscordのWebhookにPOSTする関数
func postErrMsgToDiscordWebhook(errMsg string, errLog error) {
	type data struct {
		Content string `json:"content"`
	}
	var d = data{
		Content: errMsg + ":\n" + errLog.Error(),
	}

	// 環境変数からDiscordのWebhook URLを取得
	webhookUrl := os.Getenv("Test")
	if webhookUrl == "" {
		slog.Warn("invalid team: %s", "Test")
	}

	// PostするデータをJSONに変換
	jsonData, err := json.Marshal(&d)
	if err != nil {
		slog.Warn("failed to create a new request: %v", err)
	}

	postDiscord(jsonData, webhookUrl)
}

// createDiscordWebhookData はDiscordのWebhookにPOSTする関数
func createDiscordWebhookData(notionData *NotionData) error {
	// DiscordのWebhookにPOSTするデータの構造体
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

	// DiscordのWebhookにPOSTするデータを作成
	var data = postData{}
	data.Content = "Notionに新しい投稿があります！"
	data.Embeds = embeds{}
	data.Embeds.Title = notionData.Title
	data.Embeds.Url = notionData.Url
	data.Embeds.Description = fmt.Sprintf("進捗: %s\n投稿者: %s", notionData.Status, notionData.User)
	data.Embeds.Color = 5620992

	log.Printf("postData: %v", data)

	// 環境変数からDiscordのWebhook URLを取得
	webhookUrl := os.Getenv(notionData.Team)
	if webhookUrl == "" {
		return fmt.Errorf("invalid team: %s", notionData.Team)
	}

	// PostするデータをJSONに変換
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to create a new request: %v", err)
	}

	postDiscord(jsonData, webhookUrl)

	return nil
}

// postDiscord はDiscordのWebhookにPOSTする関数
func postDiscord(d []byte, webhookUrl string) {
	// POSTリクエストを作成
	req, err := http.NewRequest("POST", webhookUrl, bytes.NewBuffer(d))
	if err != nil {
		slog.Warn("failed to create a new request: %v", err)
	}

	// Content-Typeを設定
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Warn("failed to send a request: %v", err)
	}

	defer resp.Body.Close()
}

// checkJsonData はNotionのWebhookが送信するJSONデータの必須パラメータが揃っているかチェックする関数
func checkJsonData(postData *NotionData, allData map[string]interface{}) string {
	// 必須パラメータのチェック
	data := allData["data"].(map[string]interface{})

	if postData.Url = data["url"].(string); postData.Url == "" {
		return "missing url"
	}

	if postData.Title = data["properties"].(map[string]interface{})["概要"].(map[string]interface{})["title"].(map[string]interface{})["plain_text"].(string); postData.Title == "" {
		return "missing title"
	}

	if postData.Status = data["properties"].(map[string]interface{})["進捗"].(map[string]interface{})["status"].(map[string]interface{})["name"].(string); postData.Status == "" {
		return "missing status"
	}

	if postData.User = data["properties"].(map[string]interface{})["報告者"].(map[string]interface{})["created_by"].(map[string]interface{})["name"].(string); postData.User == "" {
		return "missing user"
	}

	if postData.Team = data["properties"].(map[string]interface{})["Team"].(map[string]interface{})["rich_text"].(map[string]interface{})["plain_text"].(string); postData.Team == "" {
		return "missing team"
	}

	return ""
}

// PostNotionWebhook はNotionのWebhookを受け取る関数 (CloudFunctionsのエントリーポイント)
func PostNotionWebhook(w http.ResponseWriter, r *http.Request) {
	var postData NotionData

	// JSONデータをひとまず全て受け取る
	var allData map[string]interface{}

	// リクエストボディをJSONにデコード
	if err := json.NewDecoder(r.Body).Decode(&allData); err != nil {
		switch err {
		case io.EOF:
			if _, resErr := fmt.Fprint(w, "Success!"); resErr != nil {
				return
			}
			break
		default:
			slog.Warn("json.NewDecoder: %v", err)
			break
		}
	}

	log.Printf("allData: %v", allData)

	// 必須パラメータが揃っているかチェック
	if errLog := checkJsonData(&postData, allData); errLog != "" {
		slog.Warn("checkJsonData: %v", errLog)
		postErrMsgToDiscordWebhook("NotionのWebhookが送信するJSONデータに不備があります", fmt.Errorf(errLog))
	}

	// DiscordのWebhookに通知
	if err := createDiscordWebhookData(&postData); err != nil {
		slog.Warn("createDiscordWebhookData: %v", err)
		postErrMsgToDiscordWebhook("DiscordのWebhookにPOSTするデータの作成に失敗しました", err)
	} else {
		if _, resErr := fmt.Fprint(w, "Success!"); resErr != nil {
			return
		}
	}
}
