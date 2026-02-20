package modal

import "github.com/soltiHQ/control-plane/ui/templates/component/button"

// Variant controls the visual style of modal action buttons.
type Variant string

const (
	VariantDefault Variant = "default"
	VariantDanger  Variant = "danger"
)

// buttonVariant maps modal Variant to button.Variant.
func buttonVariant(v Variant) button.Variant {
	if v == VariantDanger {
		return button.VariantDanger
	}
	return button.VariantPrimary
}

// Method is the HTTP method for modal form submissions.
type Method string

const (
	MethodPost   Method = "post"
	MethodDelete Method = "delete"
	MethodPut    Method = "put"
	MethodPatch  Method = "patch"
)
