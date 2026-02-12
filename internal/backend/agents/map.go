package agents

import (
	"github.com/soltiHQ/control-plane/domain/model"
)

func toView(a *model.Agent) View {
	if a == nil {
		return View{}
	}
	return View{
		UptimeSeconds: a.UptimeSeconds(),

		Metadata: a.MetadataAll(),
		Labels:   a.LabelsAll(),

		Name:     a.Name(),
		Endpoint: a.Endpoint(),
		OS:       a.OS(),
		Arch:     a.Arch(),
		Platform: a.Platform(),
		ID:       a.ID(),
	}
}
