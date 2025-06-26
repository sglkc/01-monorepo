package domain

// Response represents a basic API response
// @Description Basic API response structure
type Response struct {
	Code    int    `json:"code" example:"200"`
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Operation completed successfully"`
}

// ResponseSingleData represents an API response with single data item
// @Description API response structure for single data item
type ResponseSingleData[Data any] struct {
	Code    int    `json:"code" example:"200"`
	Status  string `json:"status" example:"success"`
	Data    Data   `json:"data"`
	Message string `json:"message" example:"Operation completed successfully"`
}

// ResponseMultipleData represents an API response with multiple data items
// @Description API response structure for multiple data items
type ResponseMultipleData[Data any] struct {
	Code    int    `json:"code" example:"200"`
	Status  string `json:"status" example:"success"`
	Data    []Data `json:"data"`
	Message string `json:"message" example:"Operation completed successfully"`
}

// Empty represents an empty response data
// @Description Empty response data structure
type Empty struct{}
