package mockgrafana

import (
	"log"
	"math/rand"
	"testing"

	"github.com/grafana/grafana-api-golang-client"
)

func TestCreateServiceAccount(t *testing.T) {
	t.Run("should return error if service account exists", func(t *testing.T) {
		client := NewClient(12345678)
		request := gapi.CreateServiceAccountRequest{
			Name: StringGenerator(0),
			Role: RoleGenerator(),
		}
		client.CreateServiceAccount(request)
		_, err := client.CreateServiceAccount(request)

		if err == nil {
			t.Errorf("expected error, got none")
		}
	})

	t.Run("should create one service account", func(t *testing.T) {
		client := NewClient(12345678)
		request := gapi.CreateServiceAccountRequest{
			Name: StringGenerator(0),
			Role: RoleGenerator(),
		}
		client.CreateServiceAccount(request)

		want := 1
		got := len(client.ServiceAccountsDTO)
		if got != want {
			t.Errorf("want %d service accounts, got %d", got, want)
		}
	})

	t.Run("service account should have correct name", func(t *testing.T) {
		client := NewClient(12345678)
		arg := StringGenerator(0)
		request := gapi.CreateServiceAccountRequest{
			Name: arg,
			Role: RoleGenerator(),
		}
		client.CreateServiceAccount(request)

		want := arg
		got := client.ServiceAccountsDTO[0].Name
		if got != want {
			t.Errorf("want %s got %s", got, want)
		}
	})

	t.Run("service account should have correct role", func(t *testing.T) {
		client := NewClient(12345678)
		arg := RoleGenerator()
		request := gapi.CreateServiceAccountRequest{
			Name: StringGenerator(0),
			Role: arg,
		}
		client.CreateServiceAccount(request)

		want := arg
		got := client.ServiceAccountsDTO[0].Role
		if got != want {
			t.Errorf("want %s got %s", got, want)
		}
	})
}

func TestCreateServiceAccountToken(t *testing.T) {
	t.Run("should return error if service account doesn't exist", func(t *testing.T) {
		client := NewClient(1234568)
		request := gapi.CreateServiceAccountTokenRequest{
			Name:             StringGenerator(0),
			ServiceAccountID: int64(rand.Intn(1000)),
		}
		_, err := client.CreateServiceAccountToken(request)

		if err == nil {
			t.Errorf("expected error, but got none")
		}
	})

	t.Run("should return error if token name already exists", func(t *testing.T) {
		client := NewClient(1234568)
		request := gapi.CreateServiceAccountTokenRequest{
			Name:             StringGenerator(0),
			ServiceAccountID: int64(rand.Intn(1000)),
		}
		client.CreateServiceAccountToken(request)
		_, err := client.CreateServiceAccountToken(request)

		if err == nil {
			t.Errorf("expected error, but got none")
		}
	})

	t.Run("should create one token", func(t *testing.T) {
		client := NewClient(12345678)
		sa, _ := client.GenerateServiceAccount("", "")

		tokenRequest := gapi.CreateServiceAccountTokenRequest{
			Name:             StringGenerator(0),
			ServiceAccountID: sa.ID,
		}
		client.CreateServiceAccountToken(tokenRequest)

		want := 1
		got := len(client.Tokens)

		if got != want {
			t.Errorf("got %d number of tokens, want %d", got, want)
		}
	})

	t.Run("should create token under the proper service account", func(t *testing.T) {
		client := NewClient(12345678)
		sa, _ := client.GenerateServiceAccount("", "")

		tokenRequest := gapi.CreateServiceAccountTokenRequest{
			Name:             StringGenerator(0),
			ServiceAccountID: sa.ID,
		}

		client.CreateServiceAccountToken(tokenRequest)

		want := sa.ID
		got := client.Tokens[0].ServiceAccountID

		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})

	t.Run("token should have correct name from the request", func(t *testing.T) {
		client := NewClient(12345678)
		arg := StringGenerator(0)
		sa, _ := client.GenerateServiceAccount("", "")

		tokenRequest := gapi.CreateServiceAccountTokenRequest{
			Name:             arg,
			ServiceAccountID: sa.ID,
		}

		client.CreateServiceAccountToken(tokenRequest)

		want := arg
		got := client.Tokens[0].Name

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("token should create a key", func(t *testing.T) {
		client := NewClient(12345678)
		sa, _ := client.GenerateServiceAccount("", "")
		tokenRequest := gapi.CreateServiceAccountTokenRequest{
			Name:             StringGenerator(0),
			ServiceAccountID: sa.ID,
		}
		client.CreateServiceAccountToken(tokenRequest)

		want := true
		got := client.Tokens[0].Key != "" || len(client.Tokens[0].Key) == 0

		if got != want {
			t.Errorf("expected a string but got %v", client.Tokens[0].Key)
		}
	})
}

func TestGetServiceAccountTokens(t *testing.T) {
	t.Run("should return the right number of tokens for the service account", func(t *testing.T) {
		client := NewClient(12345678)
        argCount := 5

		sa, _ := client.GenerateServiceAccount("", "")
		sa2, _ := client.GenerateServiceAccount("", "")

        client.GenerateServiceAccountTokens(sa.ID, argCount)
		client.GenerateServiceAccountTokens(sa2.ID, 1)

		tokens, _ := client.GetServiceAccountTokens(sa.ID)

		want := argCount
		got := len(tokens)

		if got != want {
			t.Errorf("got %d tokens but expected %d", got, want)
		}
	})

	t.Run("tokens should be for the right service account", func(t *testing.T) {
		client := NewClient(12345678)
        argCount := 5

		sa, _ := client.GenerateServiceAccount("", "")
		sa2, _ := client.GenerateServiceAccount("", "")

        client.GenerateServiceAccountTokens(sa.ID, argCount)
		client.GenerateServiceAccountTokens(sa2.ID, 1)

		responseTokens, _ := client.GetServiceAccountTokens(sa.ID)

		var incorrectTokens int

		for _, responseToken := range responseTokens {
			for _, clientToken := range client.Tokens {
				if responseToken.ID == clientToken.ID && clientToken.ServiceAccountID != sa.ID {
					incorrectTokens++
					t.Logf("expected token %v to have service account %d, but got %d", clientToken.Name, sa.ID, clientToken.ServiceAccountID)
				}
			}
		}

		want := 0
		got := incorrectTokens

		if got != want {
			t.Errorf("Expected returned tokens to match service account but found %d returned tokens with wrong service account ID", got)
		}
	})
}

func TestDeleteServiceAccount(t *testing.T) {
	t.Run("should return error if service account doesn't exist", func(t *testing.T) {
		client := NewClient(12345678)
		nonexistentID := int64(rand.Intn(1000))
		_, err := client.DeleteServiceAccount(nonexistentID)

		if err == nil {
			t.Errorf("expected an error, got none")
		}
	})

	t.Run("service account should be deleted", func(t *testing.T) {
		client := NewClient(12345678)
		sa, _ := client.GenerateServiceAccount("", "")
		client.GenerateServiceAccounts(1)

		client.DeleteServiceAccount(sa.ID)

		var foundSA *gapi.ServiceAccountDTO

		for _, clientSA := range client.ServiceAccountsDTO {
			if sa.ID == clientSA.ID {
				foundSA = &clientSA
				break
			}
		}

		if foundSA != nil {
			t.Errorf("expected service account %v to be deleted but found service account %v with ID %v", sa.ID, foundSA.Name, foundSA.ID)
		}
	})
	t.Run("extra serviceaccounts should not be deleted", func(t *testing.T) {
		client := NewClient(12345678)
        countArg := 10
		sa, _ := client.GenerateServiceAccount("", "")
		client.GenerateServiceAccounts(countArg)

		client.DeleteServiceAccount(sa.ID)

		want := countArg
		got := len(client.ServiceAccountsDTO)

		if got != want {
			t.Errorf("Expected to have %d service account(s) after deletion, but got %d", want, got)
		}
	})
}

func TestDeleteServiceAccountToken(t *testing.T) {
	t.Run("should return error if token doesn't exist", func(t *testing.T) {
		client := NewClient(12345678)
		nonexistentID := int64(rand.Intn(1000))
		sa, _ := client.GenerateServiceAccount("", "")

		_, err := client.DeleteServiceAccountToken(sa.ID, nonexistentID)

		if err == nil {
			t.Errorf("expected an error, got none")
		}
	})
	t.Run("should return error if service account doesn't exist", func(t *testing.T) {
		client := NewClient(12345678)
		randomInt := int64(rand.Intn(1000))

		_, err := client.DeleteServiceAccountToken(randomInt, randomInt)

		if err == nil {
			t.Errorf("expected an error, got none")
		}
	})

	t.Run("token should be deleted", func(t *testing.T) {
		client := NewClient(12345678)
        countArg := 10
		sa, _ := client.GenerateServiceAccount("", "")

		token, _ := client.GenerateServiceAccountToken("", sa.ID)
		client.GenerateServiceAccountTokens(sa.ID, countArg)

		client.DeleteServiceAccountToken(sa.ID, token.ID)

		var foundToken *Token

		for _, clientToken := range client.Tokens {
			if clientToken.ID == token.ID {
				foundToken = &clientToken
				break
			}
		}

		if foundToken != nil {
			t.Errorf("expected token %v to be deleted but found token %v with ID %v", token.ID, foundToken.Name, foundToken.ID)
		}
	})

	t.Run("extra tokens should not be deleted", func(t *testing.T) {
		client := NewClient(12345678)
        countArg := 10
		sa, _ := client.GenerateServiceAccount("", "")

		token, err := client.GenerateServiceAccountToken("", sa.ID)
		_, err = client.GenerateServiceAccountTokens(sa.ID, countArg)
		if err != nil {
			log.Print("test", err)
		}

		client.DeleteServiceAccountToken(sa.ID, token.ID)

		want := countArg
		got := len(client.Tokens)

		if got != want {
			t.Errorf("Expected to have %d token(s) after deletion, but got %d", want, got)
		}
	})
}

func TestCreateCloudAPIKey(t *testing.T) {
	t.Run("should return error if request name exists", func(t *testing.T) {
		client := NewClient(12345678)
		input := gapi.CreateCloudAPIKeyInput{
			Name: StringGenerator(0),
			Role: RoleGenerator(),
		}
		client.CreateCloudAPIKey("", &input)
		_, err := client.CreateCloudAPIKey("", &input)

		if err == nil {
			t.Errorf("expected error, got none")
		}
	})

	t.Run("should create one api key", func(t *testing.T) {
		client := NewClient(12345678)
		input := gapi.CreateCloudAPIKeyInput{
			Name: StringGenerator(0),
			Role: RoleGenerator(),
		}
		client.CreateCloudAPIKey("", &input)

		want := 1
		got := len(client.CloudAPIKeys)
		if got != want {
			t.Errorf("want %d service accounts, got %d", got, want)
		}
	})

	t.Run("created key should have correct name", func(t *testing.T) {
		client := NewClient(12345678)

		arg := "testName"

		input := gapi.CreateCloudAPIKeyInput{
			Name: arg,
			Role: RoleGenerator(),
		}

		client.CreateCloudAPIKey("", &input)

		want := arg
		got := client.CloudAPIKeys[0].Name
		if got != want {
			t.Errorf("want %s got %s", got, want)
		}
	})
}

func TestListCloudAPIKeys(t *testing.T) {
   	client := NewClient(1234568)
    countArg := 10
    client.GenerateCloudAPIKeys(countArg)
   
    want := countArg
    got, _ := client.ListCloudAPIKeys("")

    if len(got.Items) != want {
        t.Errorf("got %v want %v", got, want)
    }
}	  

func TestDeleteCloudAPIKey(t *testing.T) {
    t.Run("should have right number of keys", func(t *testing.T) {
        client := NewClient(1234568)
        countArg := 5

        key, _ := client.GenerateCloudAPIKey("", "")
        client.GenerateCloudAPIKeys(countArg)
        
        client.DeleteCloudAPIKey("", key.Name)
       
        want := countArg
        got := len(client.CloudAPIKeys)

        if got != want {
            t.Errorf("got %v want %v", got, want)
        }
    })
	t.Run("tkey should be deleted", func(t *testing.T) {
		client := NewClient(12345678)
        countArg := 10

        key, _ := client.GenerateCloudAPIKey("", "")
        client.GenerateCloudAPIKeys(countArg)

		client.DeleteCloudAPIKey("", key.Name)

		var foundKey *gapi.CloudAPIKey

		for _, cloudAPIKey := range client.CloudAPIKeys {
			if cloudAPIKey.ID == key.ID {
				foundKey = cloudAPIKey
				break
			}
		}

		if foundKey != nil {
			t.Errorf("expected key %v to be deleted but found key %v with ID %v", key.ID, foundKey.Name, foundKey.ID)
		}
	})

	t.Run("extra keys should not be deleted", func(t *testing.T) {
		client := NewClient(12345678)
        countArg := 10

		key, err := client.GenerateCloudAPIKey("", "")
		_, err = client.GenerateCloudAPIKeys(countArg)
		if err != nil {
			log.Print("test", err)
		}

		client.DeleteCloudAPIKey("", key.Name)

		want := countArg
		got := len(client.CloudAPIKeys)

		if got != want {
			t.Errorf("Expected to have %d key(s) after deletion, but got %d", want, got)
		}
	})

}	   

func TestGenerateServiceAccount(t *testing.T) {
	client := NewClient(12345678)
	count := 100
	for i := 1; i <= count; i++ {
		_, err := client.GenerateServiceAccount("", "")
		if err != nil {
			break

		}
	}

	want := count
	got := len(client.ServiceAccountsDTO)

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGenerateServiceAccountToken(t *testing.T) {
	client := NewClient(12345678)
	count := 100
	sa, _ := client.GenerateServiceAccount("", "")
	for i := 1; i <= count; i++ {
		_, err := client.GenerateServiceAccountToken("", sa.ID)
		if err != nil {
			break

		}
	}

	want := count
	got := len(client.Tokens)

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGenerateCloudAPIKey(t *testing.T) {
	client := NewClient(1234568)
	count := 100
	for i := 1; i <= count; i++ {
		_, err := client.GenerateCloudAPIKey("", "")
		if err != nil {
			break

		}
	}

	want := count
	got := len(client.CloudAPIKeys)

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
