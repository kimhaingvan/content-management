package loanguideline

import (
	"time"
)

type GetListResponse struct {
	LoanGuidelines []*LoanGuideline `json:"loan_guidelines"`
}

type LoanGuideline struct {
	ID        int       `json:"id"`
	HTMLCode  string    `json:"html_code"`
	CSSCode   string    `json:"css_code"`
	Medias    []*Media  `json:"medias"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Media struct {
	ID              int       `json:"id"`
	URL             string    `json:"url"`
	Description     string    `json:"description"`
	LoanGuidelineID *int      `json:"loan_guideline_id,omitempty"`
	OrderID         *int      `json:"order_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
