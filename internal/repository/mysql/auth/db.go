package auth

import "shop/internal/repository/mysql"

type Repository struct {
	connection mysql.Connection
}

func New(connection mysql.Connection) Repository {
	return Repository{
		connection: connection,
	}
}
