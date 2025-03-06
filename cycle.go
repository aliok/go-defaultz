package defaultz

import "reflect"

// detectPotentialCycles checks if a type definition allows cycles.
// even though a struct *instance* may not have a cycle, we check the type definition here
// since we are creating new instances of the struct on the go in the defaultz.
func detectPotentialCycles(t reflect.Type, seen map[reflect.Type]bool) bool {
	// Dereference pointer types (if any)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// If already seen, it's a recursive reference (potential cycle)
	if seen[t] {
		return true
	}

	// Only process structs
	if t.Kind() != reflect.Struct {
		return false
	}

	// Mark this type as seen for the current recursion path
	seen[t] = true
	defer delete(seen, t) // Remove after processing

	// Check each field
	for i := range t.NumField() {
		fieldType := t.Field(i).Type
		if detectPotentialCycles(fieldType, seen) {
			return true
		}
	}

	return false
}
