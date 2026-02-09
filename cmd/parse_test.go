package cmd

import (
	"maps"
	"testing"
)

func TestParseContent(t *testing.T) {
	content := `---
name: Name
enabledByDefault: True
position: 0
---`

	options := Parse(content)

	wants := make(map[string]string)
	wants["name"] = "Name"
	wants["enabledByDefault"] = "True"
	wants["position"] = "0"

	if !maps.Equal(options, wants) {
		t.Errorf("TestParseContent(content), expected %v, but returned %v", wants, options)
	}

}

func TestParseContentWithSpaces(t *testing.T) {
	content := `---
name: Name Other
enabledByDefault: True
position: 0
---`

	options := Parse(content)

	wants := make(map[string]string)
	wants["name"] = "Name Other"
	wants["enabledByDefault"] = "True"
	wants["position"] = "0"

	if !maps.Equal(options, wants) {
		t.Errorf("TestParseContent(content), expected %v, but returned %v", wants, options)
	}

}
