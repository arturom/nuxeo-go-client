package nuxeogoclient

import (
	"fmt"
	"testing"
)

func TestImplementsInterface(t *testing.T) {
	nuxeo, err := getTestClient()
	if err != nil {
		t.Error(err)
	}
	user, err := nuxeo.Users().GetUser("Administrator")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(user.ID, user.IsAdministrator, user.IsAnonymous, user.Properties.Email)
}
