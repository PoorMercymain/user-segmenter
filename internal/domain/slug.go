package domain

type Slug struct {
	Slug           string `json:"slug"`
	PercentOfUsers int    `json:"percent,omitempty"`
}

type SlugNoPercent struct {
	Slug string `json:"slug"`
}
