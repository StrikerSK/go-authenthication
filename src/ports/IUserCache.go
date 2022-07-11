package ports

type IUserCache interface {
	CreateCache(string, interface{}) error
	RetrieveCache(string) (interface{}, error)
}
