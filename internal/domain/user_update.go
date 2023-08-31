package domain

type UserUpdate struct {
	SlugsToAdd    []string `json:"slugs_to_add" example:"SEGMENT_NAME"`
	SlugsToDelete []string `json:"slugs_to_delete" example:"ANOTHER_SEGMENT_NAME"`
	UserID        string   `json:"user_id" example:"1"`
	TTL           []string `json:"ttl,omitempty" example:"2023-09-30T20:19:05+03:00"`
}
