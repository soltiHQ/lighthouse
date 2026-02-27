package apimapv1

import (
	restv1 "github.com/soltiHQ/control-plane/api/rest/v1"
	"github.com/soltiHQ/control-plane/domain/model"
)

// Role maps a domain Role to its REST DTO.
func Role(r *model.Role) restv1.Role {
	if r == nil {
		return restv1.Role{}
	}
	return restv1.Role{
		ID:   r.ID(),
		Name: r.Name(),
	}
}
