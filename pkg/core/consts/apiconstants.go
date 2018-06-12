//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package consts

// Case insensitive literals
const (
	SubscriptionsLiteral         = "{sb:(?i)subscriptions}"
	ResourceGroupsLiteral        = "{rg:(?i)resourcegroups}"
	ResourcesLiteral             = "{rs:(?i)resources}"
	OperationStatusLiteral       = "{os:(?i)operationstatus}"
	ProvidersLiteral             = "{pv:(?i)providers}"
	ProviderRegistrationsLiteral = "{pr:(?i)providerregistrations}"
	ContainerServicesLiteral     = "{cs:(?i)containerservices}"
	LocationsLiteral             = "{lc:(?i)locations}"
	OperationResultsLiteral      = "{or:(?i)operationresults}"
	OperationsLiteral            = "{op:(?i)operations}"
	DeploymentsLiteral           = "{dp:(?i)deployments}"
	PreflightLiteral             = "{pf:(?i)preflight}"
	InternalLiteral              = "{in:(?i)internal}"
	ManagedClustersLiteral       = "{mc:(?i)managedclusters}"
	OrchestratorsLiteral         = "{or:(?i)orchestrators}"
	UpgradeProfilesLiteral       = "{us:(?i)upgradeprofiles}"
	AccessProfilesLiteral        = "{ap:(?i)accessprofiles}"
	AdminLiteral                 = "{ad:(?i)admin}"
	PodsLiteral                  = "{po:(?i)pods}"
	LogLiteral                   = "{lo:(?i)log}"
	EventLiteral                 = "{ev:(?i)events}"
	KubectlLiteral               = "{ku:(?i)kubectl}"
	ContainersLiteral            = "{co:(?i)containers}"
	UnderlaysLiteral             = "{un:(?i)underlays}"
	DefaultLiteral               = "{up:(?i)default}"
	ListSettingsLiteral          = "{li:(?i)listsettings}"
)

const (
	// PathSubscriptionIDParameter is the path parameter name used in routing for the subscription id
	PathSubscriptionIDParameter = "subscriptionId"
	// PathResourceGroupNameParameter is the path parameter name used in routing for the resource group name
	PathResourceGroupNameParameter = "resourceGroupName"
	// PathResourceNameParameter is the path parameter name used in routing for the resource name
	PathResourceNameParameter = "resourceName"
	// PathOperationStatusParameter is the path parameter name used in routing for the operation id
	PathOperationStatusParameter = "operationId"
	// PathProviderRegistrationParameter is the path parameter name used in routing for the provider registration
	PathProviderRegistrationParameter = "providerRegistration"
	// RequestAPIVersionParameterName is the query string parameter name ARM adds for the api version
	RequestAPIVersionParameterName = "api-version"
	// TerraformRPNamespace is the ARM namespace for Terraform RP
	TerraformRPNamespace = "Microsoft.TerraformOSS"
)

// subscription and common routes.
const (
	// SubscriptionsURLPrefix is the base route prefix for all subscription based operations.
	SubscriptionsURLPrefix = "/" + SubscriptionsLiteral

	// SubscriptionResourceOperationRoute is the route used to perform PUT/GET on one Subscription resource
	// /{subscriptionId}
	SubscriptionResourceOperationRoute = "/{" +
		PathSubscriptionIDParameter + "}"
)

// resource operation routes
const (
	// ResourceOperationRoute is the route used to perform PUT/GET/DELETE on one resource
	// /{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.TerraformOSS/resources/{resourceName}
	ResourceOperationRoute = SubscriptionResourceOperationRoute + "/" + ResourceGroupsLiteral + "/{" +
		PathResourceGroupNameParameter +
		"}/" + ProvidersLiteral + "/" + TerraformRPNamespace + "/" + ResourcesLiteral + "/{" +
		PathResourceNameParameter + "}"

	// ProviderRegistrationOperationRoute is the route used to perform PUT/GET/DELETE on one provider registration
	// /{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.TerraformOSS/providerregistrations/{providerRegistration}
	ProviderRegistrationOperationRoute = SubscriptionResourceOperationRoute + "/" + ResourceGroupsLiteral + "/{" +
		PathResourceGroupNameParameter +
		"}/" + ProvidersLiteral + "/" + TerraformRPNamespace + "/" + ProviderRegistrationsLiteral + "/{" +
		PathProviderRegistrationParameter + "}"

	// ProviderRegistrationListSettingsRoute is the route used to perform POST on one provider registration settings
	// /{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.TerraformOSS/providerregistrations/{providerRegistration}/listsettings
	ProviderRegistrationListSettingsRoute = SubscriptionResourceOperationRoute + "/" + ResourceGroupsLiteral + "/{" +
		PathResourceGroupNameParameter +
		"}/" + ProvidersLiteral + "/" + TerraformRPNamespace + "/" + ProviderRegistrationsLiteral + "/{" +
		PathProviderRegistrationParameter + "}" + "/" + ListSettingsLiteral

	// OperationStatusRoute is the route used to perform GET on an operation
	// /{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.TerraformOSS/operationstatus/{opeartionId}
	OperationStatusRoute = SubscriptionResourceOperationRoute + "/" + ResourceGroupsLiteral + "/{" +
		PathResourceGroupNameParameter +
		"}/" + ProvidersLiteral + "/" + TerraformRPNamespace + "/" + OperationStatusLiteral + "/{" +
		PathOperationStatusParameter + "}"
)

const (
	// GetResourceControllerName is the constant logged for get resource calls
	GetResourceControllerName = "GetResourceController"
	// PutResourceControllerName is the constant logged for put resource calls
	PutResourceControllerName = "PutResourceController"
	// DeleteResourceControllerName is the constant logged for delete resource calls
	DeleteResourceControllerName = "DeleteResourceController"

	// GetProviderRegistrationControllerName is the constant logged for get provider registration calls
	GetProviderRegistrationControllerName = "GetProviderRegistrationController"
	// PutProviderRegistrationControllerName is the constant logged for put provider registration calls
	PutProviderRegistrationControllerName = "PutProviderRegistrationController"
	// DeleteProviderRegistrationControllerName is the constant logged for delete provider registration calls
	DeleteProviderRegistrationControllerName = "DeleteProviderRegistrationController"

	// PostProviderRegistrationListSettingsControllerName is the constant logged for post provider registration settings calls
	PostProviderRegistrationListSettingsControllerName = "PostProviderRegistrationListSettingsController"

	// GetSubscriptionControllerName is the constant logged for get subscription calls
	GetSubscriptionControllerName = "GetSubscriptionController"
	// PutSubscriptionControllerName is the constant logged for put subscription calls
	PutSubscriptionControllerName = "PutSubscriptionControllerName"

	// GetOperationStatusControllerName is the constant logged for get resource operation status calls
	GetOperationStatusControllerName = "GetOperationStatus"
)

// Long running operation status constants
const (
	ProvisioningStateAccepted  = "Accepted"
	ProvisioningStateDeleting  = "Deleting"
	ProvisioningStateDeleted   = "Deleted"
	ProvisioningStateFailed    = "Failed"
	ProvisioningStateSucceeded = "Succeeded"
)

// Headers
const (
	// AcceptLanguageHeader is the standard http header name used so that we don't have to pass in the http request
	AcceptLanguageHeader = "Accept-Language"
	// HostHeader is the standard http header Host used to indicate the target host name
	HostHeader = "Host"
	// ContentTypeHeader is the standard http header Content-Type
	ContentTypeHeader = "Content-Type"
	// RequestResourceTypeParameterName is the query string parameter name to distiguish between cluster types
	RequestResourceTypeParameterName = "resource-type"
	// RequestCorrelationIDHeader is the http header name ARM adds for the correlationID
	RequestCorrelationIDHeader = "x-ms-correlation-request-id"
	// RequestClientRequestIDHeader is the http header name ARM adds for the client request id
	RequestClientRequestIDHeader = "x-ms-client-request-id"
	// RequestClientApplicationIDHeader is the http header name ARM adds for the client app id
	RequestClientApplicationIDHeader = "x-ms-client-app-id"
	// RequestClientPrincipalNameHeader is the http header name ARM adds for the client principal name
	RequestClientPrincipalNameHeader = "x-ms-client-principal-name"
	// ResponseRequestIDHeader is the http header name our rp adds to uniquely identify this request
	ResponseRequestIDHeader = "x-ms-request-id"
	// SkipTokenParameterName is the query string parameter name, optional
	SkipTokenParameterName = "skipToken"
	// RefererHeader is the refer
	RefererHeader = "Referer"

	// LocationHeader is the http header name required for Accepted status
	LocationHeader = "Location"
	// AzureAsyncOperationHeader is the http header name required for async operation
	AzureAsyncOperationHeader = "Azure-AsyncOperation"
	// RetryAfterHeader is the http header name for client's back off duration
	RetryAfterHeader = "Retry-After"
)

const (
	// OperationStatusAPIVersion is the operation status api version
	OperationStatusAPIVersion = "2018-05-01-preview"
)
