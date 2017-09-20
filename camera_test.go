package main

import (
	"testing"
)

var DummyConfig = []byte(`
{
	"storagepath": "test/",
	"cameras": [{
		"name": "Camera1",
		"url": "http://example.com/test.jpg"
	}]
}
`)
var c Config

func TestDecodeConfigJson(t *testing.T) {
	c.Decode(DummyConfig)
	if len(c.Cameras) != 1 {
		t.Fatalf("expected 1 item, got %d", len(c.Cameras))
	}
	if c.Cameras[0].Name != "Camera1" {
		t.Fatalf("expected name 'Camera1', got '%v'", c.Cameras[0].Name)
	}
}

func TestStoragePath(t *testing.T) {
	if c.StoragePath != "test/" {
		t.Fatalf("expected storage path 'test/', got '%s'", c.StoragePath)
	}
}

func TestCameraOne(t *testing.T) {
	var cam *Camera
	cam = c.Cameras[0]
	if cam.Auth != nil {
		t.Fatalf("no auth section expected, found one")
	}
	if cam.Name != "Camera1" {
		t.Fatalf("expected name 'Camera1', got '%v'", cam.Name)
	}
	if cam.URL != "http://example.com/test.jpg" {
		t.Fatalf("expected URL 'http://example.com/test.jpg', got '%v'", cam.URL)
	}
	if cam.SaveTo != "" {
		t.Fatalf("expected SaveTo '', got '%v'", cam.SaveTo)
	}
}
