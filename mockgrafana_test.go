package mockgrafana

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"testing"

	"github.com/grafana/grafana-api-golang-client"
)

func TestCreateCloudAccessPolicy(t *testing.T) {
	t.Run("should not create access policy if invalid type", func(t *testing.T) {
		client := NewClient()
		policyNameArg := "TestPolicyName"
		realmTypeArg := "InvalidType"
		realmIdentifierArg := "clabs"
		policyRealmsArg := NewRealm(realmTypeArg, realmIdentifierArg, `{env="dev"}`)
		policyScopesArg := []string{"testScope"}
		regionArg := "us"

		input := gapi.CreateCloudAccessPolicyInput{
			Name:        policyNameArg,
			DisplayName: policyNameArg,
			Scopes:      policyScopesArg,
			Realms:      []gapi.CloudAccessPolicyRealm{policyRealmsArg},
		}

		_, err := client.CreateCloudAccessPolicy(regionArg, input)

		if err == nil {
			t.Errorf("expected to get an error but was nil")
		}
	})

	t.Run("should return error if region doesn't exist", func(t *testing.T) {
		client := NewClient()
		policyNameArg := "TestPolicyName"
		realmTypeArg := "stack"
		realmIdentifierArg := "clabs"
		policyRealmsArg := NewRealm(realmTypeArg, realmIdentifierArg, `{env="dev"}`)
		policyScopesArg := []string{"testScope"}
		regionArg := ""

		input := gapi.CreateCloudAccessPolicyInput{
			Name:        policyNameArg,
			DisplayName: policyNameArg,
			Scopes:      policyScopesArg,
			Realms:      []gapi.CloudAccessPolicyRealm{policyRealmsArg},
		}

		_, err := client.CreateCloudAccessPolicy(regionArg, input)

		if err == nil {
			t.Errorf("expected to get an error but was nil")
		}
	})

	t.Run("should create access policy", func(t *testing.T) {
		client := NewClient()
		policyNameArg := "TestPolicyName"
		policyRealmsArg := NewRealm("org", "clabs", `{env="dev"}`)
		policyScopesArg := []string{"testScope"}
		regionArg := "us"

		input := gapi.CreateCloudAccessPolicyInput{
			Name:        policyNameArg,
			DisplayName: policyNameArg,
			Scopes:      policyScopesArg,
			Realms:      []gapi.CloudAccessPolicyRealm{policyRealmsArg},
		}

		client.CreateCloudAccessPolicy(regionArg, input)
		want := policyNameArg
		got := client.CloudAccessPolicyItems[0].Name

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestCloudAccessPolicies(t *testing.T) {
	t.Run("should return error if no region ", func(t *testing.T) {
		client := NewClient()
		policyNameArg := "TestPolicyName"
		policyRealmsArg := NewRealm("org", "clabs", `{env="dev"}`)
		policyScopesArg := []string{"testScope"}
		regionArg := "us"

		input := gapi.CreateCloudAccessPolicyInput{
			Name:        policyNameArg,
			DisplayName: policyNameArg,
			Scopes:      policyScopesArg,
			Realms:      []gapi.CloudAccessPolicyRealm{policyRealmsArg},
		}

		client.CreateCloudAccessPolicy(regionArg, input)

		_, err := client.CloudAccessPolicies("")

		if err == nil {
			t.Errorf("expected error but got none")
		}

	})

	t.Run("should list access policies", func(t *testing.T) {
		client := NewClient()
		policyNameArg := "TestPolicyName"
		policyRealmsArg := NewRealm("org", "clabs", `{env="dev"}`)
		policyScopesArg := []string{"testScope"}
		regionArg := "us"

		input := gapi.CreateCloudAccessPolicyInput{
			Name:        policyNameArg,
			DisplayName: policyNameArg,
			Scopes:      policyScopesArg,
			Realms:      []gapi.CloudAccessPolicyRealm{policyRealmsArg},
		}

		client.CreateCloudAccessPolicy(regionArg, input)

		items, _ := client.CloudAccessPolicies(regionArg)

		want := policyNameArg
		got := items.Items[0].Name

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestDeleteCloudAccessPolicy(t *testing.T) {
	t.Run("should delete policy if it exists", func(t *testing.T) {
		client := NewClient()
		policyNameArg := "TestPolicyName"
		policyRealmsArg := NewRealm("org", "clabs", `{env="dev"}`)
		policyScopesArg := []string{"testScope"}
		regionArg := "us"

		input := gapi.CreateCloudAccessPolicyInput{
			Name:        policyNameArg,
			DisplayName: policyNameArg,
			Scopes:      policyScopesArg,
			Realms:      []gapi.CloudAccessPolicyRealm{policyRealmsArg},
		}

		policy, _ := client.CreateCloudAccessPolicy(regionArg, input)
		client.DeleteCloudAccessPolicy(regionArg, policy.ID)

		if len(client.CloudAccessPolicyItems) > 0 {
			t.Errorf("expected access policies to be empty but found %v", client.CloudAccessPolicyItems)
		}
	})
	t.Run("should delete tokens from policy", func(t *testing.T) {
    	client := NewClient()
		policyNameArg := "TestPolicyName"
		policyRealmsArg := NewRealm("org", "clabs", `{env="dev"}`)
		policyScopesArg := []string{"testScope"}
		regionArg := "us"
        count := 100

		input := gapi.CreateCloudAccessPolicyInput{
			Name:        policyNameArg,
			DisplayName: policyNameArg,
			Scopes:      policyScopesArg,
			Realms:      []gapi.CloudAccessPolicyRealm{policyRealmsArg},
		}
		policy, _ := client.CreateCloudAccessPolicy(regionArg, input)
        client.GenerateCloudAccessPolicyTokens(count, "", policy.ID)
    	client.DeleteCloudAccessPolicy(regionArg, policy.ID)

		if len(client.CloudAccessPolicyTokenItems) > 0 {
			t.Errorf("expected access policy tokens to be empty but found %v", client.CloudAccessPolicyTokenItems)
		}
	})

	t.Run("should return error if region doesn't exist", func(t *testing.T) {
		client := NewClient()
		policyNameArg := "TestPolicyName"
		policyRealmsArg := NewRealm("org", "clabs", `{env="dev"}`)
		policyScopesArg := []string{"testScope"}
		regionArg := "us"

		input := gapi.CreateCloudAccessPolicyInput{
			Name:        policyNameArg,
			DisplayName: policyNameArg,
			Scopes:      policyScopesArg,
			Realms:      []gapi.CloudAccessPolicyRealm{policyRealmsArg},
		}

		policy, _ := client.CreateCloudAccessPolicy(regionArg, input)
		err := client.DeleteCloudAccessPolicy(regionArg, policy.ID)

		if err != nil {
			t.Errorf("expected an error but got none")
		}
	})

}

func TestCreateCloudAccessPolicyToken(t *testing.T) {
	t.Run("should return error if access policy doesn't exist", func(t *testing.T) {
		client := NewClient()
		tokenNameArg := "TestTokenName"

		tokenInput := gapi.CreateCloudAccessPolicyTokenInput{
			AccessPolicyID: "3",
			Name:           tokenNameArg,
			DisplayName:    tokenNameArg,
		}

		_, err := client.CreateCloudAccessPolicyToken("", tokenInput)

		if err == nil {
			t.Errorf("expected error, but got none")
		}
	})

	t.Run("should return error if region doesn't exist", func(t *testing.T) {
		client := NewClient()
		policyNameArg := "TestPolicyName"
		tokenNameArg := "TestTokenName"
		policyRealmsArg := NewRealm("org", "clabs", `{env="dev"}`)
		policyScopesArg := []string{"testScope"}
		regionArg := "us"

		policyInput := gapi.CreateCloudAccessPolicyInput{
			Name:        policyNameArg,
			DisplayName: policyNameArg,
			Scopes:      policyScopesArg,
			Realms:      []gapi.CloudAccessPolicyRealm{policyRealmsArg},
		}

		policy, _ := client.CreateCloudAccessPolicy(regionArg, policyInput)

		tokenInput := gapi.CreateCloudAccessPolicyTokenInput{
			AccessPolicyID: policy.ID,
			Name:           tokenNameArg,
			DisplayName:    tokenNameArg,
		}

		_, err := client.CreateCloudAccessPolicyToken("", tokenInput)

		if err == nil {
			t.Errorf("expected error, but got none")
		}
	})

	t.Run("should create access policy token", func(t *testing.T) {
		client := NewClient()
		policyNameArg := "TestPolicyName"
		tokenNameArg := "TestTokenName"
		policyRealmsArg := NewRealm("org", "clabs", `{env="dev"}`)
		policyScopesArg := []string{"testScope"}
		regionArg := "us"

		policyInput := gapi.CreateCloudAccessPolicyInput{
			Name:        policyNameArg,
			DisplayName: policyNameArg,
			Scopes:      policyScopesArg,
			Realms:      []gapi.CloudAccessPolicyRealm{policyRealmsArg},
		}

		policy, _ := client.CreateCloudAccessPolicy(regionArg, policyInput)

		tokenInput := gapi.CreateCloudAccessPolicyTokenInput{
			AccessPolicyID: policy.ID,
			Name:           tokenNameArg,
			DisplayName:    tokenNameArg,
		}

		client.CreateCloudAccessPolicyToken(regionArg, tokenInput)
		want := "MockToken"
		got := client.CloudAccessPolicyTokenItems[0].Token

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestCloudAccessPolicyTokenByID(t *testing.T) {
	t.Run("should return error if no region ", func(t *testing.T) {
		client := NewClient()
		policyNameArg := "TestPolicyName"
		tokenNameArg := "TestTokenName"
		policyRealmsArg := NewRealm("org", "clabs", `{env="dev"}`)
		policyScopesArg := []string{"testScope"}
		regionArg := "us"

		policyInput := gapi.CreateCloudAccessPolicyInput{
			Name:        policyNameArg,
			DisplayName: policyNameArg,
			Scopes:      policyScopesArg,
			Realms:      []gapi.CloudAccessPolicyRealm{policyRealmsArg},
		}
		policy, _ := client.CreateCloudAccessPolicy(regionArg, policyInput)

		tokenInput := gapi.CreateCloudAccessPolicyTokenInput{
			AccessPolicyID: policy.ID,
			Name:           tokenNameArg,
			DisplayName:    tokenNameArg,
		}

		client.CreateCloudAccessPolicyToken(regionArg, tokenInput)
		_, err := client.CloudAccessPolicyTokens("", policy.ID)

		if err == nil {
			t.Errorf("expected error but got none")
		}
	})

	t.Run("should get access policy token", func(t *testing.T) {
		client := NewClient()
		policyNameArg := "TestPolicyName"
		tokenNameArg := "TestTokenName"
		policyRealmsArg := NewRealm("org", "clabs", `{env="dev"}`)
		policyScopesArg := []string{"testScope"}
		regionArg := "us"

		policyInput := gapi.CreateCloudAccessPolicyInput{
			Name:        policyNameArg,
			DisplayName: policyNameArg,
			Scopes:      policyScopesArg,
			Realms:      []gapi.CloudAccessPolicyRealm{policyRealmsArg},
		}
		policy, _ := client.CreateCloudAccessPolicy(regionArg, policyInput)

		tokenInput := gapi.CreateCloudAccessPolicyTokenInput{
			AccessPolicyID: policy.ID,
			Name:           tokenNameArg,
			DisplayName:    tokenNameArg,
		}

		token, _ := client.CreateCloudAccessPolicyToken(regionArg, tokenInput)
		foundToken, _ := client.CloudAccessPolicyTokenByID(regionArg, token.ID)

		want := tokenNameArg
		got := foundToken.Name

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestCloudAccessPolicyTokens(t *testing.T) {
	t.Run("should return error if no region ", func(t *testing.T) {
		client := NewClient()
		policyNameArg := "TestPolicyName"
		tokenNameArg := "TestTokenName"
		policyRealmsArg := NewRealm("org", "clabs", `{env="dev"}`)
		policyScopesArg := []string{"testScope"}
		regionArg := "us"

		policyInput := gapi.CreateCloudAccessPolicyInput{
			Name:        policyNameArg,
			DisplayName: policyNameArg,
			Scopes:      policyScopesArg,
			Realms:      []gapi.CloudAccessPolicyRealm{policyRealmsArg},
		}
		policy, _ := client.CreateCloudAccessPolicy(regionArg, policyInput)

		tokenInput := gapi.CreateCloudAccessPolicyTokenInput{
			AccessPolicyID: policy.ID,
			Name:           tokenNameArg,
			DisplayName:    tokenNameArg,
		}

		client.CreateCloudAccessPolicyToken(regionArg, tokenInput)
		_, err := client.CloudAccessPolicyTokens("", policy.ID)

		if err == nil {
			t.Errorf("expected error but got none")
		}
	})

	t.Run("should not return anything if access policy ID is empty", func(t *testing.T) {
		client := NewClient()
		policyNameArg := "TestPolicyName"
		tokenNameArg := "TestTokenName"
		policyRealmsArg := NewRealm("org", "clabs", `{env="dev"}`)
		policyScopesArg := []string{"testScope"}
		regionArg := "us"

		policyInput := gapi.CreateCloudAccessPolicyInput{
			Name:        policyNameArg,
			DisplayName: policyNameArg,
			Scopes:      policyScopesArg,
			Realms:      []gapi.CloudAccessPolicyRealm{policyRealmsArg},
		}
		policy, _ := client.CreateCloudAccessPolicy(regionArg, policyInput)

		tokenInput := gapi.CreateCloudAccessPolicyTokenInput{
			AccessPolicyID: policy.ID,
			Name:           tokenNameArg,
			DisplayName:    tokenNameArg,
		}

		client.CreateCloudAccessPolicyToken(regionArg, tokenInput)
		items, _ := client.CloudAccessPolicyTokens(regionArg, "")

		if len(items.Items) > 0 {
			t.Errorf("expected empty tokens, but found %+v", items.Items)
		}
	})

	t.Run("should list access policy tokens", func(t *testing.T) {
		client := NewClient()
		policyNameArg := "TestPolicyName"
		tokenNameArg := "TestTokenName"
		policyRealmsArg := NewRealm("org", "clabs", `{env="dev"}`)
		policyScopesArg := []string{"testScope"}
		regionArg := "us"

		policyInput := gapi.CreateCloudAccessPolicyInput{
			Name:        policyNameArg,
			DisplayName: policyNameArg,
			Scopes:      policyScopesArg,
			Realms:      []gapi.CloudAccessPolicyRealm{policyRealmsArg},
		}
		policy, _ := client.CreateCloudAccessPolicy(regionArg, policyInput)

		tokenInput := gapi.CreateCloudAccessPolicyTokenInput{
			AccessPolicyID: policy.ID,
			Name:           tokenNameArg,
			DisplayName:    tokenNameArg,
		}

		client.CreateCloudAccessPolicyToken(regionArg, tokenInput)
		items, _ := client.CloudAccessPolicyTokens(regionArg, policy.ID)

		want := tokenNameArg
		got := items.Items[0].Name

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestDeleteCloudAccessPolicyToken(t *testing.T) {
	t.Run("should return error if token doesn't exist", func(t *testing.T) {
		client := NewClient()
		err := client.DeleteCloudAccessPolicyToken("", "3")
		if err == nil {
			t.Errorf("Expected an error got but none")
		}
	})

	t.Run("should return error if region doesn't exist", func(t *testing.T) {
		client := NewClient()
		policyNameArg := "TestPolicyName"
		tokenNameArg := "TestTokenName"
		policyRealmsArg := NewRealm("org", "clabs", `{env="dev"}`)
		policyScopesArg := []string{"testScope"}
		regionArg := "us"

		policyInput := gapi.CreateCloudAccessPolicyInput{
			Name:        policyNameArg,
			DisplayName: policyNameArg,
			Scopes:      policyScopesArg,
			Realms:      []gapi.CloudAccessPolicyRealm{policyRealmsArg},
		}

		policy, _ := client.CreateCloudAccessPolicy(regionArg, policyInput)

		tokenInput := gapi.CreateCloudAccessPolicyTokenInput{
			AccessPolicyID: policy.ID,
			Name:           tokenNameArg,
			DisplayName:    tokenNameArg,
		}

		token, _ := client.CreateCloudAccessPolicyToken(regionArg, tokenInput)
		err := client.DeleteCloudAccessPolicyToken("", token.ID)

		if err == nil {
			t.Errorf("expected error got none")
		}
	})

	t.Run("should delete token from the access policy", func(t *testing.T) {
		client := NewClient()
		policyNameArg := "TestPolicyName"
		tokenNameArg := "TestTokenName"
		policyRealmsArg := NewRealm("org", "clabs", `{env="dev"}`)
		policyScopesArg := []string{"testScope"}
		regionArg := "us"

		policyInput := gapi.CreateCloudAccessPolicyInput{
			Name:        policyNameArg,
			DisplayName: policyNameArg,
			Scopes:      policyScopesArg,
			Realms:      []gapi.CloudAccessPolicyRealm{policyRealmsArg},
		}

		policy, _ := client.CreateCloudAccessPolicy("", policyInput)

		tokenInput := gapi.CreateCloudAccessPolicyTokenInput{
			AccessPolicyID: policy.ID,
			Name:           tokenNameArg,
			DisplayName:    tokenNameArg,
		}

		token, _ := client.CreateCloudAccessPolicyToken(regionArg, tokenInput)
		client.DeleteCloudAccessPolicyToken(regionArg, token.ID)

		if len(client.CloudAccessPolicyTokenItems) > 0 {
			t.Errorf("tokens is %v", client.CloudAccessPolicyTokenItems)
		}
	})
}

func TestCreateServiceAccount(t *testing.T) {
	t.Run("should return error if service account exists", func(t *testing.T) {
		client := NewClient()
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
		client := NewClient()

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
		client := NewClient()

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
		client := NewClient()

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
		client := NewClient()

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
		client := NewClient()

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
		client := NewClient()

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
		client := NewClient()

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
		client := NewClient()

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
		client := NewClient()

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
		client := NewClient()

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
		client := NewClient()

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
		client := NewClient()

		nonexistentID := int64(rand.Intn(1000))
		_, err := client.DeleteServiceAccount(nonexistentID)

		if err == nil {
			t.Errorf("expected an error, got none")
		}
	})

	t.Run("service account should be deleted", func(t *testing.T) {
		client := NewClient()

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
		client := NewClient()

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
		client := NewClient()

		nonexistentID := int64(rand.Intn(1000))
		sa, _ := client.GenerateServiceAccount("", "")

		_, err := client.DeleteServiceAccountToken(sa.ID, nonexistentID)

		if err == nil {
			t.Errorf("expected an error, got none")
		}
	})
	t.Run("should return error if service account doesn't exist", func(t *testing.T) {
		client := NewClient()

		randomInt := int64(rand.Intn(1000))

		_, err := client.DeleteServiceAccountToken(randomInt, randomInt)

		if err == nil {
			t.Errorf("expected an error, got none")
		}
	})

	t.Run("token should be deleted", func(t *testing.T) {
		client := NewClient()

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
		client := NewClient()

		countArg := 10
		sa, _ := client.GenerateServiceAccount("", "")

		token, _ := client.GenerateServiceAccountToken("", sa.ID)
		_, err := client.GenerateServiceAccountTokens(sa.ID, countArg)
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
		client := NewClient()

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
		client := NewClient()

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
		client := NewClient()

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
	t.Run("created keys should have correct prefix", func(t *testing.T) {
		client := NewClient()

		countArg := 10
		client.GenerateCloudAPIKeys(countArg, "testprefix", "")

		want := countArg
		got, _ := client.ListCloudAPIKeys("")

		if len(got.Items) != want {
			t.Errorf("got %v want %v", got, want)
		}

	})
}

func TestDeleteCloudAPIKey(t *testing.T) {
	t.Run("should have right number of keys", func(t *testing.T) {
		client := NewClient()

		countArg := 20

		key, _ := client.GenerateCloudAPIKey("", "")
		client.GenerateCloudAPIKeys(countArg, "", "")

		client.DeleteCloudAPIKey("", key.Name)

		want := countArg
		got := len(client.CloudAPIKeys)

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("key should be deleted", func(t *testing.T) {
		client := NewClient()

		countArg := 10

		key, _ := client.GenerateCloudAPIKey("", "")
		client.GenerateCloudAPIKeys(countArg, "", "")

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
		client := NewClient()

		countArg := 10

		key, _ := client.GenerateCloudAPIKey("", "")
		client.GenerateCloudAPIKeys(countArg, "", "")

		client.DeleteCloudAPIKey("", key.Name)

		want := countArg
		got := len(client.CloudAPIKeys)

		if got != want {
			t.Errorf("Expected to have %d key(s) after deletion, but got %d", want, got)
		}
	})

}

func TestGenerateServiceAccount(t *testing.T) {
	client := NewClient()

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
	client := NewClient()

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
	client := NewClient()

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

func TestGenerateCloudAccessPolicy(t *testing.T) {
	t.Run("should generate policy and add it to the client", func(t *testing.T) {
		client := NewClient()
		nameArg := "TestName"
		count := 100
		for i := 1; i <= count; i++ {
			client.GenerateCloudAccessPolicy(fmt.Sprintf("%v-%d", nameArg, count))
		}

		want := count
		got := len(client.CloudAccessPolicyItems)
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("should generate policy with the correct name", func(t *testing.T) {
		client := NewClient()
		nameArg := "TestName"

		want := nameArg
		got := client.GenerateCloudAccessPolicy(nameArg).Name

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("should generate policy and return it", func(t *testing.T) {
		client := NewClient()
		nameArg := "TestName"

		want := "*gapi.CloudAccessPolicy"
		got := client.GenerateCloudAccessPolicy(nameArg)

		if reflect.TypeOf(got).String() != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

func TestGenerateCloudAccessPolicies(t *testing.T) {
	t.Run("should generate the specified number of policies", func(t *testing.T) {
		client := NewClient()
		nameArg := "TestName"
		count := 10
		client.GenerateCloudAccessPolicies(count, nameArg)

		want := count
		got := len(client.CloudAccessPolicyItems)

		if got != want {
			t.Errorf("got %v policies but want %v", got, want)
		}
	})
}

func TestGenerateCloudAccessPolicyToken(t *testing.T) {
    t.Run("should generate tokens and add it to the policy", func(t *testing.T) {
        client := NewClient()
        nameArg := "TestArg"
        count := 100

        accessPolicy := client.GenerateCloudAccessPolicy(nameArg)
        client.GenerateCloudAccessPolicyTokens(count, fmt.Sprintf("%v", nameArg), accessPolicy.ID )

        want := count
        got := len(client.CloudAccessPolicyTokenItems)

        if got != want {
            t.Errorf("got %q want %q", got, want)
        }

    })

    t.Run("should generate token with the correct name", func(t *testing.T) {
        client := NewClient()
        nameArg := "TestName"
        accessPolicy := client.GenerateCloudAccessPolicy(nameArg)

        want := nameArg
        got := client.GenerateCloudAccessPolicyToken(nameArg, accessPolicy.ID).Name

        if got != want {
            t.Errorf("got %q want %q", got, want)
        }
    })

    t.Run("should generate policy and return it", func(t *testing.T) {
        client := NewClient()
        nameArg := "TestName"
        accessPolicy := client.GenerateCloudAccessPolicy(nameArg)

        want := "*gapi.CloudAccessPolicyToken"
        got := client.GenerateCloudAccessPolicyToken(nameArg, accessPolicy.ID)

        if reflect.TypeOf(got).String() != want {
            t.Errorf("got %q want %q", got, want)
        }
    })
}

func TestGenerateCloudAccessPolicyTokens(t *testing.T) {
	t.Run("should generate the specified number of tokens", func(t *testing.T) {
		client := NewClient()
		nameArg := "TestName"
        accessPolicy := client.GenerateCloudAccessPolicy(nameArg)
		count := 10

		client.GenerateCloudAccessPolicyTokens(count, nameArg, accessPolicy.ID)

		want := count
		got := len(client.CloudAccessPolicyTokenItems)

		if got != want {
			t.Errorf("got %v policies but want %v", got, want)
		}
	})
}
