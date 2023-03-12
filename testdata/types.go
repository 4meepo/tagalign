package testdata

type StructA struct {
	FieldA int `json:"a"`
	FieldB int `json:"b" validate:"required"`
}
