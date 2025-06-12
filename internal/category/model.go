package category

// Category represents a hierarchical item category.
type Category struct {
	ID       string  `db:"id" json:"id"`
	Name     string  `db:"name" json:"name"`
	ParentID *string `db:"parent_id" json:"parent_id"`
}
