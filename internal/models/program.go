package models

type Program struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Category    string `json:"category"` // e.g., 'unggulan', 'kajian', 'pendidikan', 'remaja', 'sosial'
	ArabicTitle string `json:"arabic_title"`
	Description string `json:"description"`
	Ustadz      string `json:"ustadz"`
	Schedule    string `json:"schedule"`
	Level       string `json:"level"`
	Quota       string `json:"quota"`
	IsFeatured  bool   `json:"is_featured"`
	ShowOnHome  bool   `json:"show_on_home"`
}
