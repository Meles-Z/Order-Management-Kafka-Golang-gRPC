package dto

type ProductEvent struct {
	EventType string      `json:"event_type"` // "create", "update", "delete"
	Payload   interface{} `json:"payload"`    // Can be CreateProductEvent, UpdateProductEvent, or DeleteProductEvent
}

type CreateProductEvent struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	IsActive    bool    `json:"is_active"`
}

type UpdateProductEvent struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	IsActive    bool    `json:"is_active"`
}

type DeleteProductEvent struct {
	ID string `json:"id"`
}
