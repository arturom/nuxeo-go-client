package nuxeogoclient

import "fmt"

type UserProperties struct {
	FirstName string   `json:"fistName"`
	LastName  string   `json:"lastName"`
	Groups    []string `json:"groups"`
	Company   string   `json:"company"`
	Email     string   `json:"email"`
}

type NuxeoUser struct {
	ID              string         `json:"id"`
	IsAdministrator bool           `json:"isAdministrator"`
	IsAnonymous     bool           `json:"isAnonymous"`
	Properties      UserProperties `json:"properties"`
}

type UsersClient interface {
	GetUser(username string) (*NuxeoUser, error)
}

type usersClient struct {
	nuxeo client
}

func newUserClient(nuxeo client) UsersClient {
	c := new(usersClient)
	c.nuxeo = nuxeo
	return c
}

func (c usersClient) GetUser(username string) (*NuxeoUser, error) {
	user := new(NuxeoUser)
	path := fmt.Sprintf("/api/v1/user/%s", username)
	_, err := c.nuxeo.GetJson(path, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
