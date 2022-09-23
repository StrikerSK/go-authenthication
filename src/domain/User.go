package domain

import "encoding/json"

type User struct {
	UserCredentials
}

func (u User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}
