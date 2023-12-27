package domain

import "encoding/json"

type UserDTO struct {
	UserCredentials
}

func (u *UserDTO) MarshalBinary() ([]byte, error) {
	u.ClearPassword()
	return json.Marshal(u)
}

func (u *UserDTO) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
