package schemas

type ErrorPayload struct {
	Detail string `json:"detail"`
}

type MessagePayload struct {
	Message string `json:"message"`
}
