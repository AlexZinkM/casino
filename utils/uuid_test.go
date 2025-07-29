package utils

import (
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	uuid1 := GenerateUUID()
	uuid2 := GenerateUUID()

	if uuid1 == "" {
		t.Error("Expected UUID to be generated")
	}

	if uuid2 == "" {
		t.Error("Expected UUID to be generated")
	}

	if uuid1 == uuid2 {
		t.Error("Expected UUIDs to be different")
	}

	if len(uuid1) != 36 {
		t.Errorf("Expected UUID length 36, got %d", len(uuid1))
	}

	if len(uuid2) != 36 {
		t.Errorf("Expected UUID length 36, got %d", len(uuid2))
	}
}

func TestGenerateUUID_Format(t *testing.T) {
	uuid := GenerateUUID()

	if len(uuid) != 36 {
		t.Errorf("Expected UUID length 36, got %d", len(uuid))
	}

	if uuid[8] != '-' || uuid[13] != '-' || uuid[18] != '-' || uuid[23] != '-' {
		t.Error("Expected UUID to have proper format with dashes")
	}
}
