package domain

type Slug struct {
	Slug           string `json:"slug" example:"SEGMENT_NAME"`
	PercentOfUsers int    `json:"percent,omitempty" example:"10"`
}

type SlugNoPercent struct {
	Slug string `json:"slug" example:"SEGMENT_NAME"`
}
