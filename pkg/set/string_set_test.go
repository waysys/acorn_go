// ----------------------------------------------------------------------------
//
// String Set
//
// Author: William Shaffer
// Version: 07-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package set

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"os"
	"strconv"
	"testing"
)

// ----------------------------------------------------------------------------
// Test Main
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

// ----------------------------------------------------------------------------
// Test functions
// ----------------------------------------------------------------------------

// Test_NewStringSet checks the operation of creating a new set, IsEmpty, and
// Size.
func Test_NewStringSet(t *testing.T) {
	var set = NewStringSet()
	var testFunction = func(t *testing.T) {
		if !set.IsEmpty() {
			t.Error("An empty set was not created")
		}
		if set.Size() != 0 {
			t.Error("An empty set must have a size of 0")
		}
	}

	t.Run("Test_Total", testFunction)
}

// Test_Add checks the operation of adding an item to the set.
func Test_Add(t *testing.T) {
	var set = NewStringSet()
	set.Add("Hello")
	var testFunction = func(t *testing.T) {
		if set.IsEmpty() {
			t.Error("set should not be empty after adding a string")
		}
		if set.Size() != 1 {
			t.Error("set size should be one, not: " + strconv.Itoa(set.Size()))
		}
		set.Add("Hello")
		if set.Size() != 1 {
			t.Error("set size was incorrect: " + strconv.Itoa(set.Size()))
		}
		if !set.Contains("Hello") {
			t.Error("Hello was not found in set")
		}
	}

	t.Run("Test_Total", testFunction)
}

// Test_Delete checks the operaton of deleting a string from the set.
func Test_Delete(t *testing.T) {
	var set = NewStringSet()
	set.Add("Hello")
	set.Delete("Hello")
	var testFunction = func(t *testing.T) {
		if !set.IsEmpty() {
			t.Error("Delete did not clear set: " + strconv.Itoa(set.Size()))
		}
		set.Delete("Tobasco")
		if !set.IsEmpty() {
			t.Error("Delete of nonexistent item failed")
		}
		if set.Contains("Hello") {
			t.Error("set was found with item that was deleted")
		}
	}

	t.Run("Test_Total", testFunction)
}

// Test_ToSlice checks the operation of converting a set into a slice of
// strings.
func Test_To_Slice(t *testing.T) {
	var set = NewStringSet()
	set.Add("Hello")
	set.Add("World")
	var values = set.ToSlice()

	var testFunction = func(t *testing.T) {
		if len(values) != 2 {
			t.Error("slice is not the correct length")
		}
	}
	t.Run("Test_Total", testFunction)
}

// Test_AddAll checks the operation of adding a slice of strings to
// the set.
func Test_AddAll(t *testing.T) {
	var set = NewStringSet()
	var values = []string{
		"One",
		"Two",
		"Three",
		"Three",
		"Four",
	}
	set.AddAll(values)
	var testFunction = func(t *testing.T) {
		if set.Size() != 4 {
			t.Error("invalid size of set: " + strconv.Itoa(set.Size()))
		}
	}

	t.Run("Test_Total", testFunction)
}
