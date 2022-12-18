package mockgrafana

import (
	"fmt"
	"github.com/grafana/grafana-api-golang-client"
	"math/rand"
	"time"
)

type MockClient struct {
	OrgID              int64
	ServiceAccountsDTO []gapi.ServiceAccountDTO
	Tokens             []Token
	CloudAPIKeys       []*gapi.CloudAPIKey
}

type Token struct {
	ID               int64      `json:"id"`
	Name             string     `json:"name"`
	Created          time.Time  `json:"created,omitempty"`
	Key              string     `json:"key"`
	Expiration       *time.Time `json:"expiration,omitempty"`
	ServiceAccountID int64      `json:"-"`
	SecondsToLive    int64      `json:"secondsToLive,omitempty"`
}

func NewClient(orgID int64) *MockClient {
	return &MockClient{
		OrgID: orgID,
	}
}

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
		OrgID:  client.OrgID,
		Role:   request.Role,
		Tokens: 0,
	}
	client.ServiceAccountsDTO = append(client.ServiceAccountsDTO, serviceAccount)

	return &serviceAccount, nil

}

func (client *MockClient) CreateServiceAccountToken(request gapi.CreateServiceAccountTokenRequest) (*gapi.CreateServiceAccountTokenResponse, error) {
	var found bool
	for _, sa := range client.ServiceAccountsDTO {
		if sa.ID == request.ServiceAccountID {
			found = true
		}
	}
	if found == false {
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

func (client *MockClient) GetServiceAccounts() ([]gapi.ServiceAccountDTO, error) {
	return client.ServiceAccountsDTO, nil
}

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

func (client *MockClient) DeleteServiceAccount(serviceAccountID int64) (*gapi.DeleteServiceAccountResponse, error) {
	var found bool
	for idx, sa := range client.ServiceAccountsDTO {
		if sa.ID == serviceAccountID {
			client.ServiceAccountsDTO[idx] = client.ServiceAccountsDTO[len(client.ServiceAccountsDTO)-1]
			client.ServiceAccountsDTO[len(client.ServiceAccountsDTO)-1] = gapi.ServiceAccountDTO{}
			client.ServiceAccountsDTO = client.ServiceAccountsDTO[:len(client.ServiceAccountsDTO)-1]
			found = true
		}
	}

	if found != true {
		return nil, fmt.Errorf("could not find token")
	}

	return nil, nil
}

func (client *MockClient) DeleteServiceAccountToken(serviceAccountID, tokenID int64) (*gapi.DeleteServiceAccountResponse, error) {
	var saFound bool
	for _, sa := range client.ServiceAccountsDTO {
		if sa.ID == serviceAccountID {
			saFound = true
		}
	}
	if saFound != true {
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

	if tokenFound != true {
		return nil, fmt.Errorf("token not found")
	}
	return nil, nil
}

func (client *MockClient) ListCloudAPIKeys(org string) (*gapi.ListCloudAPIKeysOutput, error) {
	return &gapi.ListCloudAPIKeysOutput{
		Items: client.CloudAPIKeys,
	}, nil
}

func (client *MockClient) DeleteCloudAPIKey(org string, keyName string) error {

	for idx, key := range client.CloudAPIKeys {
		if keyName == key.Name {
            
            keys := client.CloudAPIKeys
            keys = append(keys[:idx], keys[idx+1:]...)
            client.CloudAPIKeys = keys
/*
			client.CloudAPIKeys[idx] = client.CloudAPIKeys[len(client.CloudAPIKeys)-1]
			client.CloudAPIKeys[len(client.CloudAPIKeys)-1] = nil
			client.CloudAPIKeys = client.CloudAPIKeys[:len(client.CloudAPIKeys)-1]
            */
		}
	}
	return nil
}

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

func RoleGenerator() string {
	rand.Seed(time.Now().UnixNano())
	roles := []string{"Admin", "Viewer", "Editor"}
	return roles[rand.Intn(len(roles))]
}

func StringGenerator(seed int) string {
	rand.Seed(time.Now().UnixNano() + int64(seed))
	return fmt.Sprintf("randomString-%d%d", rand.Intn(99999), rand.Intn(99999))
}
