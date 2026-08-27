package category

import (
	"context"
	"errors"
	"shop/internal/entity"
	"shop/internal/pkg/richerror"
	"time"

	"github.com/go-sql-driver/mysql"
)

func (r Repository) Create(ctx context.Context, c entity.Category) (entity.Category, error) {
	const op = "category-repository.Create"

	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	const query = `INSERT INTO category(parent_id,title,slug,image,deleted_at,created_at,updated_at) VALUES (?,?,?,?,?,?,?)`

	res, err := r.connection.DB.ExecContext(ctx, query, c.ParentID, c.Title, c.Slug, c.Image, c.DeletedAt, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return entity.Category{}, richerror.New().
				SetOp(op).
				SetMsg("duplicate entry").
				SetKind(richerror.KindConflictErr).
				SetErr(err)
		}
		return entity.Category{}, richerror.New().
			SetOp(op).
			SetMsg("unexpected error").
			SetKind(richerror.KindUnexpectedErr).
			SetErr(err)
	}

	id , _ := res.LastInsertId()
	c.ID = uint(id)
	
	return c , nil
}
