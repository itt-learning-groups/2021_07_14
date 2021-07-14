package handlers

type Todo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type echoRequest struct {
	Key1 string `json:"key1"`
}

type echoPutResponse struct {
	Param   string      `json:"param"`
	ReqBody echoRequest `json:"reqBody"`
}

type createRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=25"`
	Description string `json:"description,omitempty" validate:"max=100"`
}

type createResponse Todo

type getAllResponse struct {
	List []Todo `json:"list,omitempty"`
}

type editRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=25"`
	Description string `json:"description" validate:"required,max=100"`
}

type EditResponse Todo
