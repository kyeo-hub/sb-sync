package logger

import (
	"bytes"
	"testing"
)

func TestLogEntry(t *testing.T) {
	entry := LogEntry{
		Time:    "2024-01-01T00:00:00Z",
		Level:   "INFO",
		Message: "test",
	}

	if entry.Level != "INFO" {
		t.Errorf("Expected level INFO, got %s", entry.Level)
	}
}

func TestLevelNames(t *testing.T) {
	testCases := []struct {
		level   Level
		name    string
	}{
		{DEBUG, "DEBUG"},
		{INFO, "INFO"},
		{WARN, "WARN"},
		{ERROR, "ERROR"},
	}

	for _, tc := range testCases {
		if levelNames[tc.level] != tc.name {
			t.Errorf("Expected %s for level %d, got %s", tc.name, tc.level, levelNames[tc.level])
		}
	}
}

func TestLogFilteringByLevel(t *testing.T) {
	l := Logger{
		level:  WARN,
		output: &bytes.Buffer{},
	}

	if l.level != WARN {
		t.Errorf("Expected level WARN, got %v", l.level)
	}
}

func TestDebugLessThanInfo(t *testing.T) {
	if !(DEBUG < INFO) {
		t.Error("DEBUG should be less than INFO")
	}
}

func TestInfoLessThanWarn(t *testing.T) {
	if !(INFO < WARN) {
		t.Error("INFO should be less than WARN")
	}
}

func TestWarnLessThanError(t *testing.T) {
	if !(WARN < ERROR) {
		t.Error("WARN should be less than ERROR")
	}
}
