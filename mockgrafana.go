package mockgrafana

import (
	"fmt"
	"github.com/grafana/grafana-api-golang-client"
	"math/rand"
	"time"
)

// MockClient is a substitute for the real grafana api token so we can
// simulate the same behavior
type MockClient struct {
	ServiceAccountsDTO []gapi.ServiceAccountDTO
	Tokens             []Token
	CloudAPIKeys       []*gapi.CloudAPIKey
}

// Token  is a simulation of a grafana api token
type Token struct {
	ID               int64      `json:"id"`
	Name             string     `json:"name"`
	Created          time.Time  `json:"created,omitempty"`
	Key              string     `json:"key"`
	Expiration       *time.Time `json:"expiration,omitempty"`
	ServiceAccountID int64      `json:"-"`
	SecondsToLive    int64      `json:"secondsToLive,omitempty"`
}

func (ClientWrapper *MockClient) Initialize(org string) error {
	return nil
}

// NewClient returns a MockClient for use in simulating the grafana api key
func NewClient() *MockClient {
	return &MockClient{}
}

// CreateServiceAccount is a Mock of the grafana api method, that will take a CreateServieAccountRequest and will create and return
// the grafana service account created
func (client *MockClient) CreateServiceAccount(request gapi.CreateServiceAccountRequest) (*gapi.ServiceAccountDTO, error) {
	for _, sa := range client.ServiceAccountsDTO {
		if sa.Name == request.Name {
			return nil, fmt.Errorf("service account name must be unique")
		}
	}

	serviceAccount := gapi.ServiceAccountDTO{
		ID:     int64(len(client.ServiceAccountsDTO) + 1),
		Name:   request.Name,
		Login:  fmt.Sprintf("sa-%s", request.Name),
		Role:   request.Role,
		Tokens: 0,
	}
	client.ServiceAccountsDTO = append(client.ServiceAccountsDTO, serviceAccount)
	return &serviceAccount, nil
}

// CreateServiceAccountToken is a Mock of the grafana api method, that will take a CreateServiceAccountTokenRequest and will create
// and return the grafana service account token created.
func (client *MockClient) CreateServiceAccountToken(request gapi.CreateServiceAccountTokenRequest) (*gapi.CreateServiceAccountTokenResponse, error) {
	var found bool
	for _, sa := range client.ServiceAccountsDTO {
		if sa.ID == request.ServiceAccountID {
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("service account not found")
	}

	for _, token := range client.Tokens {
		if token.Name == request.Name {
			return nil, fmt.Errorf("token name must be unique")
		}
	}

	token := Token{
		ID:               int64(len(client.Tokens) + 1),
		Name:             request.Name,
		Created:          time.Now(),
		ServiceAccountID: request.ServiceAccountID,
		Key:              fmt.Sprintf("%s-%d", request.Name, int64(rand.Intn(99999))),
	}

	client.Tokens = append(client.Tokens, token)

	for _, sa := range client.ServiceAccountsDTO {
		if sa.ID == request.ServiceAccountID {
			sa.Tokens++
		}
	}

	return &gapi.CreateServiceAccountTokenResponse{
		ID:   token.ID,
		Name: token.Name,
		Key:  token.Key,
	}, nil
}

// GetServiceAccounts is a Mock of the grafana api method, that will list all service accounts
func (client *MockClient) GetServiceAccounts() ([]gapi.ServiceAccountDTO, error) {
	return client.ServiceAccountsDTO, nil
}

// GetServiceAccountTokens is a Mock of the grafana api method, that will take a serviceAccountID and return a GetServiceAccountTokensResponse
func (client *MockClient) GetServiceAccountTokens(serviceAccountID int64) ([]gapi.GetServiceAccountTokensResponse, error) {
	response := make([]gapi.GetServiceAccountTokensResponse, 0)

	for _, token := range client.Tokens {
		if token.ServiceAccountID == serviceAccountID {
			response = append(response, gapi.GetServiceAccountTokensResponse{
				ID:         token.ID,
				Name:       token.Name,
				Created:    token.Created,
				Expiration: token.Expiration,
			})
		}
	}
	return response, nil
}

// DeleteServiceAccount is a Mock of the grafana api method, that will take a serviceAccountID and delete the service account
func (client *MockClient) DeleteServiceAccount(serviceAccountID int64) (*gapi.DeleteServiceAccountResponse, error) {
	for idx, sa := range client.ServiceAccountsDTO {
		if sa.ID == serviceAccountID {
			client.ServiceAccountsDTO[idx] = client.ServiceAccountsDTO[len(client.ServiceAccountsDTO)-1]
			client.ServiceAccountsDTO[len(client.ServiceAccountsDTO)-1] = gapi.ServiceAccountDTO{}
			client.ServiceAccountsDTO = client.ServiceAccountsDTO[:len(client.ServiceAccountsDTO)-1]
			return nil, nil
		}
	}
	return nil, fmt.Errorf("could not find token")
}

// DeleteServiceAccountToken is a Mock of the grafana api method, that will take a serviceAccountID and tokenID, and deletes
// the token from that service account
func (client *MockClient) DeleteServiceAccountToken(serviceAccountID, tokenID int64) (*gapi.DeleteServiceAccountResponse, error) {
	var saFound bool
	for _, sa := range client.ServiceAccountsDTO {
		if sa.ID == serviceAccountID {
			saFound = true
		}
	}
	if !saFound {
		return nil, fmt.Errorf("service account not found")
	}
	var tokenFound bool
	for idx, token := range client.Tokens {
		if token.ServiceAccountID == serviceAccountID && token.ID == tokenID {
			client.Tokens[idx] = client.Tokens[len(client.Tokens)-1]
			client.Tokens[len(client.Tokens)-1] = Token{}
			client.Tokens = client.Tokens[:len(client.Tokens)-1]
			tokenFound = true
		}
	}

	if !tokenFound {
		return nil, fmt.Errorf("token not found")
	}
	return nil, nil
}

// ListCloudAPIKeys is a Mock of the grafana api method, that  will return the list all Cloud API Keys
func (client *MockClient) ListCloudAPIKeys(org string) (*gapi.ListCloudAPIKeysOutput, error) {
	return &gapi.ListCloudAPIKeysOutput{
		Items: client.CloudAPIKeys,
	}, nil
}

// DeleteCloudAPIKey is a Mock of the grafana api method, that will delete the specified key
func (client *MockClient) DeleteCloudAPIKey(org string, keyName string) error {

	for idx, key := range client.CloudAPIKeys {
		if keyName == key.Name {
			copy(client.CloudAPIKeys[idx:], client.CloudAPIKeys[idx+1:])
			client.CloudAPIKeys[len(client.CloudAPIKeys)-1] = &gapi.CloudAPIKey{}
			client.CloudAPIKeys = client.CloudAPIKeys[:len(client.CloudAPIKeys)-1]
		}
	}
	return nil
}

// CreateCloudAPIKey is a Mock of the grafana api method, that will create the specified cloud api key
func (client *MockClient) CreateCloudAPIKey(org string, input *gapi.CreateCloudAPIKeyInput) (*gapi.CloudAPIKey, error) {
	for _, key := range client.CloudAPIKeys {
		if key.Name == input.Name {
			return nil, fmt.Errorf("cloud api key must be unique")
		}
	}
	newKey := &gapi.CloudAPIKey{
		ID:    len(client.CloudAPIKeys) + 1,
		Name:  input.Name,
		Role:  input.Role,
		Token: fmt.Sprintf("%v-%v", input.Name, rand.Intn(99999)),
	}

	client.CloudAPIKeys = append(client.CloudAPIKeys, newKey)
	return newKey, nil
}

// GenerateCloudAPIKeys generates x number of APIKeys (x specified by count) with an option prefix and role.
// if role isn't specified, then it a random one will be generated.
func (client *MockClient) GenerateCloudAPIKeys(count int, prefix, role string) ([]*gapi.CloudAPIKey, error) {
	var keys []*gapi.CloudAPIKey
	var name string

	for i := 0; i < count; i++ {
		if prefix != "" {
			name = prefix + "-" + StringGenerator(len(client.CloudAPIKeys)+1)
		}
		key, err := client.GenerateCloudAPIKey(name, role)
		if err != nil {
			return nil, err
		}

		keys = append(keys, key)
	}
	return keys, nil
}

// GenerateCloudAPIKey generates a CloudAPIKey with the supplied inputs (name/role) or generates them randomly
// if not given
func (client *MockClient) GenerateCloudAPIKey(name, role string) (*gapi.CloudAPIKey, error) {
	if name == "" {
		name = StringGenerator(len(client.CloudAPIKeys) + 1)
	}

	if role == "" {
		role = RoleGenerator()
	}
	tokenRequest := gapi.CreateCloudAPIKeyInput{
		Name: name,
		Role: role,
	}
	return client.CreateCloudAPIKey("", &tokenRequest)
}

// GenerateServiceAccountTokens take a service account ID and count integer, and then Generates
// that many service accounts and returns a CreateServiceAccountTokenResponse
func (client *MockClient) GenerateServiceAccountTokens(saID int64, count int) ([]*gapi.CreateServiceAccountTokenResponse, error) {

	var serviceAccountTokenResponses []*gapi.CreateServiceAccountTokenResponse
	for i := 0; i < count; i++ {
		resp, err := client.GenerateServiceAccountToken("", saID)
		if err != nil {
			return nil, err
		}
		serviceAccountTokenResponses = append(serviceAccountTokenResponses, resp)
	}
	return serviceAccountTokenResponses, nil
}

// GenerateServiceAccountToken Generates a ServiceAccountToken
func (client *MockClient) GenerateServiceAccountToken(name string, saID int64) (*gapi.CreateServiceAccountTokenResponse, error) {
	if name == "" {
		name = StringGenerator(len(client.Tokens) + 1)
	}

	tokenRequest := gapi.CreateServiceAccountTokenRequest{
		Name:             name,
		ServiceAccountID: saID,
	}
	return client.CreateServiceAccountToken(tokenRequest)
}

// GenerateServiceAccounts takes a count integer and Generates that many service accounts
func (client *MockClient) GenerateServiceAccounts(count int) ([]*gapi.ServiceAccountDTO, error) {
	var serviceAccounts []*gapi.ServiceAccountDTO
	for i := 0; i < count; i++ {
		sa, err := client.GenerateServiceAccount("", "")
		if err != nil {
			return nil, err
		}
		serviceAccounts = append(serviceAccounts, sa)
	}

	return serviceAccounts, nil
}

// GenerateServiceAccount takes a name and a role and returns a service account.  If name and role
// aren't specified, it will create with random information
func (client *MockClient) GenerateServiceAccount(name, role string) (*gapi.ServiceAccountDTO, error) {
	if name == "" {
		name = StringGenerator(len(client.ServiceAccountsDTO) + 1)
	}
	if role == "" {
		role = RoleGenerator()
	}
	request := gapi.CreateServiceAccountRequest{
		Name: name,
		Role: role,
	}

	return client.CreateServiceAccount(request)
}

// RoleGenerator returns a random role from the list
func RoleGenerator() string {
	rand.Seed(time.Now().UnixNano())
	roles := []string{"Admin", "Viewer", "Editor", "MetricsPublisher"}
	return roles[rand.Intn(len(roles))]
}

// StringGenerator returns a random string
func StringGenerator(seed int) string {
	rand.Seed(time.Now().UnixNano() + int64(seed))
	return fmt.Sprintf("randomString-%d%d", rand.Intn(99999), rand.Intn(99999))
}
