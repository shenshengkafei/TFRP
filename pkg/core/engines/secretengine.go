//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package engines

import (
	"TFRP/pkg/core/consts"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/arm/keyvault"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
)

// SecretEngine is the secret engine
type SecretEngine struct {
	TenantID     string
	ClientID     string
	ClientSecret string
}

// GetSecretEngine creates a secret engine
func GetSecretEngine() (secretEngine *SecretEngine) {
	tenantID, err := ioutil.ReadFile(consts.ServicePrincipalTenantIDPath)
	if err != nil {
		log.Fatal("Failed to get tenant id: ", err)
	}

	clientID, err := ioutil.ReadFile(consts.ServicePrincipalClientIDPath)
	if err != nil {
		log.Fatal("Failed to get client id: ", err)
	}

	clientSecret, err := ioutil.ReadFile(consts.ServicePrincipalClientSecretPath)
	if err != nil {
		log.Fatal("Failed to get client secret: ", err)
	}

	secretEngine = new(SecretEngine)
	secretEngine.TenantID = string(tenantID)
	secretEngine.ClientID = string(clientID)
	secretEngine.ClientSecret = string(clientSecret)

	fmt.Printf("tenant %s", secretEngine.TenantID)
	fmt.Printf("clientID %s", secretEngine.ClientID)
	fmt.Printf("clientsecret %s", secretEngine.ClientSecret)

	return secretEngine
}

// GetSecretFromKeyVault returns a secret
func (secretEngine *SecretEngine) GetSecretFromKeyVault(vaultBaseURI string, secretName string, secretVersion string) string {
	fmt.Printf("uri %s", vaultBaseURI)
	fmt.Printf("name %s", secretName)
	fmt.Printf("version %s", secretVersion)

	tenantID := secretEngine.TenantID
	clientID := secretEngine.ClientID
	clientSecret := secretEngine.ClientSecret

	fmt.Printf("tenant %s", secretEngine.TenantID)
	fmt.Printf("clientID %s", secretEngine.ClientID)
	fmt.Printf("clientsecret %s", secretEngine.ClientSecret)

	oauthConfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, tenantID)
	updatedAuthorizeEndpoint, err := url.Parse("https://login.windows.net/" + tenantID + "/oauth2/token")
	oauthConfig.AuthorizeEndpoint = *updatedAuthorizeEndpoint
	spToken, err := adal.NewServicePrincipalToken(*oauthConfig, clientID, clientSecret, "https://vault.azure.net")

	if err != nil {
		log.Fatal("failed to create token", err)
	}

	vaultsClient := keyvault.NewWithoutDefaults()
	vaultsClient.Authorizer = autorest.NewBearerAuthorizer(spToken)

	vault, err := vaultsClient.GetSecret(context.Background(), vaultBaseURI, secretName, secretVersion)
	if err != nil {
		log.Fatal("Failed to get secret ", err)
	}
	return *vault.Value
}
