package nuxeogoclient

func getTestClient() (NuxeoClient, error) {
	baseUrl := "http://localhost:8080/nuxeo"
	userName := "Administrator"
	password := userName

	nuxeo := NewNuxeoClient(baseUrl, NewBasicAuth(userName, password))
	token, err := nuxeo.RequestAuthenticationToken("MyApp", "My Device", "Device Description", "rw")

	if err != nil {
		return nil, err
	}

	return NewNuxeoClient(baseUrl, NewTokenAuth(token)), nil
}
