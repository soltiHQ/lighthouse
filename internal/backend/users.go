package backend

import (
	"context"

	"github.com/soltiHQ/control-plane/domain/model"
	"github.com/soltiHQ/control-plane/internal/storage"
)

type UsersListResult struct {
	Items      []*model.User
	NextCursor string
	Total      int
}

type Users struct {
	store storage.UserStore
}

func NewUsers(store storage.UserStore) *Users {
	return &Users{store: store}
}

func (x *Users) List(ctx context.Context, limit int, cursor string, filter storage.UserFilter) (*UsersListResult, error) {
	if limit <= 0 {
		limit = 5
	}

	res, err := x.store.ListUsers(ctx, filter, storage.ListOptions{
		Limit:  limit,
		Cursor: cursor,
	})
	if err != nil {
		return nil, err
	}

	return &UsersListResult{
		Items:      res.Items,
		NextCursor: res.NextCursor,
		Total:      len(res.Items),
	}, nil
}
