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
func checkJsonData(postData *NotionData, allData map[string]any) string {
	errMsg := ""

	// 必須パラメータのチェック
	data := allData["data"].(map[string]interface{})

	if postData.Url = getJsonValue[string](data, "url"); postData.Url == "" {
		errMsg += "missing url\n"
	}

	if postData.Title = getJsonValue[[]string](data, "properties", "概要", "title")[2]; postData.Title == "" {
		errMsg += "missing title\n"
	}

	if postData.Status = getJsonValue[string](data, "properties", "進捗", "status", "name"); postData.Status == "" {
		errMsg += "missing status\n"
	}

	if postData.User = getJsonValue[string](data, "properties", "報告者", "created_by", "name"); postData.User == "" {
		errMsg += "missing user\n"
	}

	if postData.Team = getJsonValue[[]string](data, "properties", "Team", "rich_text")[2]; postData.Team == "" {
		errMsg += "missing team\n"
	}

	return errMsg
}

// getJsonValue はJSONデータから指定したキーの値を取得する関数
func getJsonValue[T any](m map[string]any, keys ...string) T {
	var empty T

	if len(keys) == 0 {
		return empty
	}

	key := keys[0]
	value := m[key]

	if len(keys) == 1 {
		v, ok := value.(T)
		if !ok {
			return empty
		}
		return v
	}

	nextMap, ok := value.(map[string]any)
	if !ok {
		return empty
	}
	nextKeys := keys[1:]

	return getJsonValue[T](nextMap, nextKeys...)
}

// PostNotionWebhook はNotionのWebhookを受け取る関数 (CloudFunctionsのエントリーポイント)
func PostNotionWebhook(w http.ResponseWriter, r *http.Request) {
	var postData NotionData

	// JSONデータをひとまず全て受け取る
	var allData map[string]any

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
