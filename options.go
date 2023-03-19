package tagalign

type Option func(*Helper)

// WithMode specify the mode of tagalign.
func WithMode(mode Mode) Option {
	return func(h *Helper) {
		h.mode = mode
	}
}

// WithAutoSort enable auto sort tags.
// Param fixedOrder specify the fixed order of tags, the other tags will be sorted by name.
func WithAutoSort(fixedOrder ...string) Option {
	return func(h *Helper) {
		h.autoSort = true
		h.fixedTagOrder = fixedOrder
	}
}
