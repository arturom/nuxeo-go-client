package nuxeogoclient

import "testing"

func TestImplementsAuthInterface(t *testing.T) {
	NewNuxeoClient("http://localhost", NewBasicAuth("username", "password"))
	NewNuxeoClient("http://localhost", NewTokenAuth("token"))
}

func TestBasicAuthConstructor(t *testing.T) {
	auth := NewBasicAuth("Administrator", "Administrator")
	if auth.authString != "Basic QWRtaW5pc3RyYXRvcjpBZG1pbmlzdHJhdG9y" {
		t.Fail()
	}
}
