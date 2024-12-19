package c

import "time"

type NotionJsonData struct {
	Data struct {
		Archived  bool `json:"archived"`
		Cover     any  `json:"cover"`
		CreatedBy struct {
			ID     string `json:"id"`
			Object string `json:"object"`
		} `json:"created_by"`
		CreatedTime  time.Time `json:"created_time"`
		Icon         any       `json:"icon"`
		ID           string    `json:"id"`
		InTrash      bool      `json:"in_trash"`
		LastEditedBy struct {
			ID     string `json:"id"`
			Object string `json:"object"`
		} `json:"last_edited_by"`
		LastEditedTime time.Time `json:"last_edited_time"`
		Object         string    `json:"object"`
		Parent         struct {
			DatabaseID string `json:"database_id"`
			Type       string `json:"type"`
		} `json:"parent"`
		Properties struct {
			Team struct {
				ID       string `json:"id"`
				RichText []struct {
					Annotations struct {
						Bold          bool   `json:"bold"`
						Code          bool   `json:"code"`
						Color         string `json:"color"`
						Italic        bool   `json:"italic"`
						Strikethrough bool   `json:"strikethrough"`
						Underline     bool   `json:"underline"`
					} `json:"annotations"`
					Href      any    `json:"href"`
					PlainText string `json:"plain_text"`
					Text      struct {
						Content string `json:"content"`
						Link    any    `json:"link"`
					} `json:"text"`
					Type string `json:"type"`
				} `json:"rich_text"`
				Type string `json:"type"`
			} `json:"Team"`
			Tag struct {
				ID          string `json:"id"`
				MultiSelect []struct {
					Color string `json:"color"`
					ID    string `json:"id"`
					Name  string `json:"name"`
				} `json:"multi_select"`
				Type string `json:"type"`
			} `json:"タグ"`
			Priority struct {
				ID     string `json:"id"`
				Select any    `json:"select"`
				Type   string `json:"type"`
			} `json:"優先度"`
			Reporter struct {
				CreatedBy struct {
					AvatarURL string `json:"avatar_url"`
					ID        string `json:"id"`
					Name      string `json:"name"`
					Object    string `json:"object"`
					Person    struct {
						Email string `json:"email"`
					} `json:"person"`
					Type string `json:"type"`
				} `json:"created_by"`
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"報告者"`
			Person struct {
				ID     string `json:"id"`
				People []any  `json:"people"`
				Type   string `json:"type"`
			} `json:"担当者"`
			FixedDate struct {
				Date struct {
					End      any       `json:"end"`
					Start    time.Time `json:"start"`
					TimeZone any       `json:"time_zone"`
				} `json:"date"`
				ID   string `json:"id"`
				Type string `json:"type"`
			} `json:"期日"`
			Summary struct {
				ID    string `json:"id"`
				Title []struct {
					Annotations struct {
						Bold          bool   `json:"bold"`
						Code          bool   `json:"code"`
						Color         string `json:"color"`
						Italic        bool   `json:"italic"`
						Strikethrough bool   `json:"strikethrough"`
						Underline     bool   `json:"underline"`
					} `json:"annotations"`
					Href      any    `json:"href"`
					PlainText string `json:"plain_text"`
					Text      struct {
						Content string `json:"content"`
						Link    any    `json:"link"`
					} `json:"text,omitempty"`
					Type    string `json:"type"`
					Mention struct {
						Date struct {
							End      any       `json:"end"`
							Start    time.Time `json:"start"`
							TimeZone any       `json:"time_zone"`
						} `json:"date"`
						Type string `json:"type"`
					} `json:"mention,omitempty"`
				} `json:"title"`
				Type string `json:"type"`
			} `json:"概要"`
			WaitingAsNextTask struct {
				HasMore  bool   `json:"has_more"`
				ID       string `json:"id"`
				Relation []any  `json:"relation"`
				Type     string `json:"type"`
			} `json:"次のタスクにより保留中："`
			WaitingNextTask struct {
				HasMore  bool   `json:"has_more"`
				ID       string `json:"id"`
				Relation []any  `json:"relation"`
				Type     string `json:"type"`
			} `json:"次のタスクを保留中："`
			Progress struct {
				ID     string `json:"id"`
				Status struct {
					Color string `json:"color"`
					ID    string `json:"id"`
					Name  string `json:"name"`
				} `json:"status"`
				Type string `json:"type"`
			} `json:"進捗"`
		} `json:"properties"`
		PublicURL any    `json:"public_url"`
		RequestID string `json:"request_id"`
		URL       string `json:"url"`
	} `json:"data"`
	Source struct {
		ActionID     string `json:"action_id"`
		Attempt      int    `json:"attempt"`
		AutomationID string `json:"automation_id"`
		EventID      string `json:"event_id"`
		Type         string `json:"type"`
	} `json:"source"`
}
