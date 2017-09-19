package main

import (
	"testing"
)

func TestDecodeCameraJson(t *testing.T) {
	var file_contents = []byte(`
	[{
		"name": "Camera1",
		"url": "http://example.com/test.jpg"
	}]
	`)
	var c Cameras
	c.Decode(file_contents)
	if len(c) != 1 {
		t.Fatalf("expected 1 item, got %d", len(c))
	}
	if c[0].Name != "Camera1" {
		t.Fatalf("expected name 'Camera1', got '%v'", c[0].Name)
	}
}
