package category

import (
	"context"
	"shop/internal/pkg/richerror"
)

func (s Service) Delete(ctx context.Context, slug string) error {
	const op = "category-service.Delete"
	// 1. get category and check exists category
	category, err := s.repository.GetOneBySlug(ctx, slug)
	if err != nil {
		return err
	}
	// 2. check category has children ? if has children
	// you not allowed
	// first you delete child
	if category.ParentID == nil {
		childrens, err := s.repository.GetChildrenByParentID(ctx, category.ID)
		if err != nil {
			return err
		}
		if len(childrens) > 0 {
			return richerror.New().
				SetOp(op).
				SetMsg("first you delete this category childrens").
				SetKind(richerror.KindBadRequestErr)
		}
	}
	// 3. call repository
	err = s.repository.DeleteByID(ctx, category.ID)
	if err != nil {
		return err
	}
	// return
	return nil
}
