package tagalign

type Option func(*Helper)

// WithAutoSort enable auto sort tags.
// Param fixedOrder specify the fixed order of tags, the other tags will be sorted by name.
func WithAutoSort(fixedOrder ...string) Option {
	return func(h *Helper) {
		h.autoSort = true
		h.fixedTagOrder = fixedOrder
	}
}
