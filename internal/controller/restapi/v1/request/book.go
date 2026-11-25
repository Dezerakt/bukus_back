package request

type StoreBook struct {
	Name string `json:"name" validate:"required"`
}

type UpdateBook struct {
	Name string `json:"name" validate:"required"`
}
