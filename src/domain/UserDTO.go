package domain

import (
	"encoding/json"
)

type UserDTO struct {
	UserCredentials
}

func (u UserDTO) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *UserDTO) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

//TODO Fix unmarshalling from Redis
//func (u UserDTO) MarshalJSON() ([]byte, error) {
//	customMap := map[string]string{
//		"username": u.Username,
//	}
//	return json.Marshal(customMap)
//}
