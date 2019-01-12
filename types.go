package main

// APIResponse is the wrapper response type.
type APIResponse struct {
	Error int        `json:"error"`
	Items []AntiJoke `json:"items"`
}

// AntiJoke is the basic type for this API.
type AntiJoke struct {
	ID         int    `json:"id" db:"id"`
	FirstPart  string `json:"first_part" db:"first_part"`
	SecondPart string `json:"second_part" db:"second_part"`
}
