package storage

type Storage interface {
	CloseConnection()
}

type User struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}
