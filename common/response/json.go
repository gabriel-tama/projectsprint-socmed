package response

type ResponseBody struct {
	Message string      `json:"message"`
	Data    any         `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Pagination `json:"meta,omitempty"`
}

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}
