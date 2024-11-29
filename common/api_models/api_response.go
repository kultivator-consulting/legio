package api_models

type Failed struct {
	Status       string `json:"status"`
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

type Success struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type Results struct {
	Status       string      `json:"status"`
	TotalRows    int64       `json:"totalRows"`
	TotalPages   int         `json:"totalPages"`
	PreviousPage int         `json:"previousPage"`
	NextPage     int         `json:"nextPage"`
	Data         interface{} `json:"data"`
}
