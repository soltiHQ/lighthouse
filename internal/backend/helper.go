package backend

import "github.com/soltiHQ/control-plane/internal/storage"

func NormalizeListLimit(qlimit, dlimit int) int {
	if qlimit <= 0 {
		qlimit = dlimit
	}
	if qlimit > storage.MaxListLimit {
		qlimit = storage.MaxListLimit
	}
	return qlimit
}
