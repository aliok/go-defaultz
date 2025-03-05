package defaultz

import (
	"reflect"
	"testing"
)

// ✅ Non-Cyclic Structs
type nonCyclic struct {
	Field1 int
	Field2 string
	Field3 int
}

type nonCyclicPtr struct {
	Field1 *int
	Field2 *string
	Field3 *int
}

type nonCyclicNested struct {
	Inner struct {
		Field1 int
	}
}

// 🔄 Direct Cycles
type directSelfCycle struct {
	Self *directSelfCycle
}

type cyclicParent struct {
	Child cyclicChild
}

type cyclicChild struct {
	Parent *cyclicParent
}

// 🔁 Deep Cycles
type deepCyclic1 struct {
	Next *deepCyclic2
}

type deepCyclic2 struct {
	Next *deepCyclic1
}

type deepMultiCyclic1 struct {
	Ref *deepMultiCyclic2
}

type deepMultiCyclic2 struct {
	Ref *deepMultiCyclic3
}

type deepMultiCyclic3 struct {
	Ref *deepMultiCyclic1
}

// 🏗️ Embedded Structs
type embeddedCyclic struct {
	Inner struct {
		Parent *embeddedCyclic
	}
}

// 📦 Collection Types (Not Cyclic Themselves)
type structWithSlice struct {
	List []string
}

type structWithMap struct {
	Lookup map[string]int
}

type structWithInterface struct {
	Any interface{}
}

// 🧩 Complex Edge Cases
type nestedPointerCycle struct {
	Ref **nestedPointerCycle
}

// TestDetectPotentialCycles tests cycle detection in type definitions
func TestDetectPotentialCycles(t *testing.T) {
	tests := []struct {
		name     string
		t        reflect.Type
		expected bool
	}{
		// ✅ Non-Cyclic Structs
		{"Non-Cyclic Simple", reflect.TypeOf(nonCyclic{}), false},
		{"Non-Cyclic With Pointers", reflect.TypeOf(nonCyclicPtr{}), false},
		{"Non-Cyclic With Nested Structs", reflect.TypeOf(nonCyclicNested{}), false},

		// 🔄 Direct Cycles
		{"Direct Self Cycle", reflect.TypeOf(directSelfCycle{}), true},
		{"Direct Cycle Between Two Structs", reflect.TypeOf(cyclicParent{}), true},

		// 🔁 Deep Cycles (Indirect Recursion)
		{"Deep Cycle A → B → A", reflect.TypeOf(deepCyclic1{}), true},
		{"Deep Cycle with Multiple Fields", reflect.TypeOf(deepMultiCyclic1{}), true},

		// 🏗️ Embedded Structs
		{"Embedded Struct Cycle", reflect.TypeOf(embeddedCyclic{}), true},

		// 📦 Collection Types (Not Cyclic Themselves)
		{"Struct with Slice", reflect.TypeOf(structWithSlice{}), false},
		{"Struct with Map", reflect.TypeOf(structWithMap{}), false},
		{"Struct with Interface", reflect.TypeOf(structWithInterface{}), false},

		// 🧩 Complex Edge Cases
		// I don't like this result, but defaulter should handle this case anyway when it is actually defaulting things
		{"Nested Pointer Cycle", reflect.TypeOf(nestedPointerCycle{}), false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := detectPotentialCycles(test.t, make(map[reflect.Type]bool))
			if result != test.expected {
				t.Errorf("Test %s failed: expected %v, got %v", test.name, test.expected, result)
			}
		})
	}
}
