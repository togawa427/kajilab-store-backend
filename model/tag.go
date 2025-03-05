package model

type TagGetResponseTag struct {
	Name	string 									`json:"name"`
}

type TagGetResponse struct {
	Logs 	[]TagGetResponseTag     `json:"tags"`
}

type TagPostRequestTag struct {
	Name	*string									`json:"name"`
}

type TagPostRequest struct {
	Tag 	*TagGetResponseTag 			`json:"tag"`
}