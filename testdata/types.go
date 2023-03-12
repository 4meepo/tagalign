package testdata

type StructA struct {
	FieldA int `json:"aaa"   validate:"required"`
	FieldB int `json:"b"   validate:"required"`
}
