package response_body

type Id struct {
	Id      int    `json:"id"`
	Message string `json:"message,omitempty"`
}
