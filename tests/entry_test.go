package tests

import (
	. "knowledgebase/models"
	"testing"
)

func TestEntry_WithFilledValues_IsValid(t *testing.T) {
	entry := &Entry{
		DocumentLocation: "test",
		Elements: []Element{
			Element{
				Text: "Test",
				Type: "Test",
			},
		},
	}
	err := entry.Validate()
	if err != nil {
		t.Fatal("Entry should be valid")
	}
}

func TestEntry_WithZeroValue_IsNotValid(t *testing.T) {
	entry := &Entry{}
	err := entry.Validate()
	if err == nil {
		t.Fatal("Entry should not be valid")
	}
}

func TestEntry_WithEmptyDocumentLocation_IsNotValid(t *testing.T) {
	entry := &Entry{
		DocumentLocation: "",
		Elements: []Element{
			Element{
				Text: "Test",
				Type: "Test",
			},
		},
	}
	err := entry.Validate()
	if err == nil {
		t.Fatal("Entry should not be valid")
	}
}

func TestEntry_WithNilElements_IsNotValid(t *testing.T) {
	entryNilElement := &Entry{
		DocumentLocation: "test",
		Elements:         []Element{},
	}

	err := entryNilElement.Validate()
	if err == nil {
		t.Fatal("Elements should always have atleast 1 item")
	}
}

func TestEntry_WithEmptyElements_IsNotValid(t *testing.T) {
	entryEmptyElement := &Entry{
		DocumentLocation: "test",
		Elements:         make([]Element, 0),
	}

	err := entryEmptyElement.Validate()
	if err == nil {
		t.Fatal("Elements should always have atleast 1 item")
	}
}
