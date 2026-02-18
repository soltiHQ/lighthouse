package agent

import (
	v1 "github.com/soltiHQ/control-plane/api/v1"
	"github.com/soltiHQ/control-plane/ui/templates/component/modal"
)

// editLabelFields builds editable fields from existing agent labels.
// Each label becomes a text field; the handler collects them back into a labels map.
func editLabelFields(a v1.Agent) []modal.Field {
	fields := make([]modal.Field, 0, len(a.Labels))
	for k, v := range a.Labels {
		fields = append(fields, modal.Field{
			ID:          k,
			Label:       k,
			Value:       v,
			Placeholder: "value",
		})
	}
	return fields
}
