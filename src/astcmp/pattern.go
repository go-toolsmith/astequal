package astcmp

type patternKind uint16

type pattern struct {
	kind patternKind

	// Specifies index in shared []string storage that
	// contains bound value.
	// Index 0 always contains empty string.
	valueIndex uint16

	// Index in shared []pattern storage that
	// contains first sub-pattern.
	subIndex uint16

	// Sub-patterns count.
	subCount uint16
}
