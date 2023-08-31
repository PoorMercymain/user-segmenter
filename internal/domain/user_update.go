package domain

type UserUpdate struct {
	SlugsToAdd    []string `json:"slugs_to_add"`
	SlugsToDelete []string `json:"slugs_to_delete"`
	UserID        string   `json:"user_id"`
	TTL           []string `json:"ttl,omitempty"`
}
