package delivery

// CreateDTO describes a request to create a delivery record.
type CreateDTO struct {
	Provider     string `json:"provider" validate:"required"`
	TrackingCode string `json:"tracking_code" validate:"required"`
	Status       string `json:"status"`
}

// UpdateDTO describes fields that can be updated.
type UpdateDTO struct {
	Status *string `json:"status"`
}
