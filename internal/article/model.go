package article

type Article struct {
	ID          int    `json:"id,omitempty"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Category    string `json:"category"`
	Status      string `json:"status"`
	CreatedDate string `json:"created_date,omitempty"`
	UpdatedDate string `json:"updated_date,omitempty"`
}

var validStatus = map[string]bool{"publish": true, "draft": true, "thrash": true}

func (a Article) Validate() map[string]string {
	errs := map[string]string{}
	if n := len([]rune(a.Title)); n < 20 {
		errs["title"] = "min 20 characters"
	} else if n > 200 {
		errs["title"] = "max 200 characters"
	}
	if n := len([]rune(a.Content)); n < 200 {
		errs["content"] = "min 200 characters"
	}
	if n := len([]rune(a.Category)); n < 3 {
		errs["category"] = "min 3 characters"
	} else if n > 100 {
		errs["category"] = "max 100 characters"
	}
	if !validStatus[a.Status] {
		errs["status"] = "must be publish, draft or thrash"
	}
	return errs
}
