package address

import "context"

func (s Service) Delete(ctx context.Context, userID uint, addressID uint) error {
	const op = "address-service.Delete"
	// get address with this id
	a, err := s.repository.GetOne(ctx, userID, addressID)
	if err != nil {
		return err
	}
	// call delete methd from repository
	err = s.repository.Delete(ctx, a.UserID, a.ID)
	if err != nil {
		return err
	}
	return nil
}
