package auth

import (
	"context"
	"errors"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
	"time"

	"github.com/go-sql-driver/mysql"
)

func (r Repository) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	const op = "auth-repository.CreateUser"
	// fill the created_at & updated_at field
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	// insert user in to database
	// with this query
	query := `INSERT INTO user (name , avatar , phone_number , password , created_at , updated_at) VALUES (? , ? , ? , ? , ? , ?)`
	// execute query
	result, err := r.connection.DB.ExecContext(ctx, query, user.Name, user.Avatar, user.PhoneNumber, user.Password, user.CreatedAt, user.UpdatedAt)
	
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return entity.User{}, richerror.New().
				SetOp(op).
				SetMsg("user with this data already exist").
				SetKind(richerror.KindConflictErr).
				SetErr(err)
		}
		return entity.User{}, richerror.New().
			SetOp(op).
			SetMsg("unexpected error in create user in databse").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}
	// this field always nil
	userID, _ := result.LastInsertId()
	user.ID = uint(userID)
	return user, nil
}
