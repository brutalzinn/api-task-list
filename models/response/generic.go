package response_entities

type GenericResponse struct {
	Error   bool
	Message string
	Data    any
}
