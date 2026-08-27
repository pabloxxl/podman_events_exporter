package utils

import (
	"testing"
)

func TestCheckRegex(t *testing.T) {
	c := &ConfigOpts{}
	if !(c.checkRegex(".*")) {
		t.Errorf("Expected regex to match, but it did not")
	}
}

func TestCheckRegexInvalid(t *testing.T) {
	c := &ConfigOpts{}
	if c.checkRegex(".**") {
		t.Errorf("Expected regex to not match, but it did")
	}
}
