package visual

// Variant controls the color scheme of visual elements (badges, dots, etc.).
type Variant string

const (
	VariantPrimary   Variant = "primary"
	VariantSecondary Variant = "secondary"
	VariantSuccess   Variant = "success"
	VariantDanger    Variant = "danger"
	VariantMuted     Variant = "muted"
)
