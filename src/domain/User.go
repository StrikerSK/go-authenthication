package domain

import "encoding/json"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}
