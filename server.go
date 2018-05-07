//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package main

import (
	"TFRP/pkg/core/consts"
	"TFRP/pkg/core/controllers"
	"TFRP/pkg/core/engines"
	"TFRP/pkg/core/storage"
	"crypto/tls"
	"encoding/base64"
	"log"
	"net/http"

	restful "github.com/emicklei/go-restful"
	"github.com/spf13/pflag"
)

var (
	addr       = pflag.String("insecure-address", ":8080", "The <host>:<port> for insecure (HTTP) serving")
	secureAddr = pflag.String("secure-address", ":443", "The <host>:<port> for secure (HTTPS) serving")
)

func main() {
	pflag.Parse()

	secretEngine := engines.GetSecretEngine()

	initRoutes(secretEngine)

	httpServer := &http.Server{
		Addr:      ":443",
		TLSConfig: getTLSConfig(secretEngine),
	}

	go func() {
		go log.Fatal(http.ListenAndServe(":8080", nil))
	}()
	log.Fatal(httpServer.ListenAndServeTLS("", ""))
}

func getTLSConfig(secretEngine *engines.SecretEngine) (config *tls.Config) {
	certPem, err := base64.StdEncoding.DecodeString(secretEngine.GetSecretFromKeyVault(consts.SslCertKVBaseURI, consts.SslCertKVSecretName, consts.SslCertKVSecretVersion))
	keyPem, err := base64.StdEncoding.DecodeString(secretEngine.GetSecretFromKeyVault(consts.SslPrivatekeyKVBaseURI, consts.SslPrivatekeyKVSecretName, consts.SslPrivatekeyKVSecretVersion))
	if err != nil {
		log.Fatal("Failed to decode certs: %v", err)
	}

	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		log.Fatal("Cannot load X509 key pair: %v", err)
	}

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	// ciphersuite requirements:
	// https://requirements.azurewebsites.net/Requirements/Details/6417#guide
	// they have to follow the order in above requirement page
	tlsConfig.CipherSuites = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
	}
	tlsConfig.PreferServerCipherSuites = true
	tlsConfig.MinVersion = tls.VersionTLS12
	tlsConfig.MaxVersion = tls.VersionTLS12

	return tlsConfig
}

func initRoutes(secretEngine *engines.SecretEngine) {
	storagePassword := secretEngine.GetSecretFromKeyVault(consts.StoragePasswordKVBaseURI, consts.StoragePasswordKVSecretName, consts.StoragePasswordKVSecretVersion)
	providerRegistrationDataProvider := storage.NewProviderRegistrationDataProvider(consts.StorageDatabase, storagePassword)
	resourceDataProvider := storage.NewResourceDataProvider(consts.StorageDatabase, storagePassword)
	providerRegistrationManager := controllers.NewProviderRegistrationManager(providerRegistrationDataProvider)
	resourceManager := controllers.NewResourceManager(providerRegistrationDataProvider, resourceDataProvider)

	webService := new(restful.WebService)
	webService.
		Path(consts.SubscriptionsURLPrefix).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	addSubscriptionOperationRoutes(webService)
	addProvidersOperationRoutes(webService, providerRegistrationManager)
	addResourcesOperationRoutes(webService, resourceManager)

	restful.Add(webService)
}

func addProvidersOperationRoutes(webService *restful.WebService, providerRegistrationManager *controllers.ProviderRegistrationManager) {
	webService.Route(webService.
		GET(consts.ProviderRegistrationOperationRoute).
		To(providerRegistrationManager.GetProviderRegistrationController).
		Doc("Get a provider registration").
		Operation(consts.GetProviderRegistrationControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathProviderRegistrationParameter, "Name of provider registration").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))

	webService.Route(webService.
		PUT(consts.ProviderRegistrationOperationRoute).
		To(providerRegistrationManager.PutProviderRegistrationController).
		Doc("Create/update a provider registration").
		Operation(consts.PutProviderRegistrationControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathProviderRegistrationParameter, "Name of provider registration").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))

	webService.Route(webService.
		DELETE(consts.ProviderRegistrationOperationRoute).
		To(providerRegistrationManager.DeleteProviderRegistrationController).
		Doc("Delete a provider registration").
		Operation(consts.DeleteProviderRegistrationControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathProviderRegistrationParameter, "Name of provider registration").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))

	webService.Route(webService.
		POST(consts.ProviderRegistrationListSettingsRoute).
		To(providerRegistrationManager.PostProviderRegistrationListSettings).
		Doc("Get settings of a provider registration").
		Operation(consts.PostProviderRegistrationListSettingsControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathProviderRegistrationParameter, "Name of provider registration").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))
}

func addResourcesOperationRoutes(webService *restful.WebService, resourceManager *controllers.ResourceManager) {
	webService.Route(webService.
		GET(consts.ResourceOperationRoute).
		To(resourceManager.GetResourceController).
		Doc("Get a resource").
		Operation(consts.GetResourceControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceNameParameter, "Name of resource").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))

	webService.Route(webService.
		PUT(consts.ResourceOperationRoute).
		To(resourceManager.PutResourceController).
		Doc("Create/update a resource").
		Operation(consts.PutResourceControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceNameParameter, "Name of resource").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))

	webService.Route(webService.
		DELETE(consts.ResourceOperationRoute).
		To(resourceManager.DeleteResourceController).
		Doc("Delete a resource").
		Operation(consts.DeleteResourceControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceNameParameter, "Name of resource").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))
}

func addSubscriptionOperationRoutes(webService *restful.WebService) {
	// Subscription operations
	webService.Route(webService.
		GET(consts.SubscriptionResourceOperationRoute).
		To(controllers.GetSubscriptionOperationController).
		Doc("get a subscription").
		Operation(consts.GetSubscriptionControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "identifier of the subscription").DataType("string")))

	webService.Route(webService.
		PUT(consts.SubscriptionResourceOperationRoute).
		To(controllers.PutSubscriptionOperationController).
		Doc("Put or update a subscription").
		Operation(consts.PutSubscriptionControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "identifier of the subscription").DataType("string")))
}
