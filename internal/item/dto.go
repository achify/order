package item

// ItemCreateDTO defines input for creating an item.
type ItemCreateDTO struct {
	Name       string `json:"name" validate:"required"`
	Price      int64  `json:"price" validate:"required"`
	CategoryID string `json:"category_id" validate:"required"`
}

// ItemUpdateDTO defines input for updating an item.
type ItemUpdateDTO struct {
	Name       *string `json:"name"`
	Price      *int64  `json:"price"`
	CategoryID *string `json:"category_id"`
}
