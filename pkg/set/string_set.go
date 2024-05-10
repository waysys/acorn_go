// ----------------------------------------------------------------------------
//
// # String Set
//
// Author: William Shaffer
// Version: 07-May-2024
//
// # Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------
//
// The set package contains implementations of sets for different base types.
package set

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type StringSet map[string]bool

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

func NewStringSet() StringSet {
	var set = make(StringSet)
	return set
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------
//
// IsEmpty return true if the set has zero items.
func (set *StringSet) IsEmpty() bool {
	return len(*set) == 0
}

// Size returns the number of items in the set.
func (set *StringSet) Size() int {
	return len(*set)
}

// Contains returns true if the specified value is in the set.
func (set *StringSet) Contains(value string) bool {
	return (*set)[value]
}

// Add inserts a string into the set.  If the string is already in the set,
// there is no change.
func (set *StringSet) Add(value string) {
	(*set)[value] = true
}

// AddAll inserts a collection of strings into the set.
func (set *StringSet) AddAll(values []string) {
	for _, value := range values {
		set.Add(value)
	}
}

// Delete deletes a string from the set.  If the string is not in the set,
// nothing happens.
func (set *StringSet) Delete(value string) {
	if (*set).Contains(value) {
		delete((*set), value)
	}
}

// ToSlice returns a slice with the strings in the set.
func (set *StringSet) ToSlice() []string {
	var values []string
	for value := range *set {
		values = append(values, value)
	}
	return values
}
