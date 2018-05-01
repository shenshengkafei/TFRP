//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package engines

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/arm/keyvault"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
)

// SecretEngine is the secret engine
type SecretEngine struct {
	tenantID     string
	clientID     string
	clientSecret string
}

// GetSecretFromKeyVault returns a secret
func (secretEngine *SecretEngine) GetSecretFromKeyVault(vaultBaseURI string, secretName string, secretVersion string) string {
	fmt.Printf("uri %s", vaultBaseURI)
	fmt.Printf("name %s", secretName)
	fmt.Printf("version %s", secretVersion)

	clientID := secretEngine.clientID
	clientSecret := secretEngine.clientSecret
	tenantID := secretEngine.tenantID

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
