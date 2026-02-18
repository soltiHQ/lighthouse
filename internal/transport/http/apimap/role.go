package apimap

import (
	v1 "github.com/soltiHQ/control-plane/api/v1"
	"github.com/soltiHQ/control-plane/domain/model"
)

func Role(r *model.Role) v1.Role {
	if r == nil {
		return v1.Role{}
	}

	return v1.Role{
		ID:   r.ID(),
		Name: r.Name(),
	}
}
