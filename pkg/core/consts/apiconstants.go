//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package consts

// Case insensitive literals
const (
	SubscriptionsLiteral     = "{sb:(?i)subscriptions}"
	ResourceGroupsLiteral    = "{rg:(?i)resourcegroups}"
	ProvidersLiteral         = "{pv:(?i)providers}"
	ResourcesLiteral         = "{rs:(?i)resources}"
	ContainerServicesLiteral = "{cs:(?i)containerservices}"
	LocationsLiteral         = "{lc:(?i)locations}"
	OperationResultsLiteral  = "{or:(?i)operationresults}"
	OperationsLiteral        = "{op:(?i)operations}"
	DeploymentsLiteral       = "{dp:(?i)deployments}"
	PreflightLiteral         = "{pf:(?i)preflight}"
	InternalLiteral          = "{in:(?i)internal}"
	ManagedClustersLiteral   = "{mc:(?i)managedclusters}"
	OrchestratorsLiteral     = "{or:(?i)orchestrators}"
	UpgradeProfilesLiteral   = "{us:(?i)upgradeprofiles}"
	AccessProfilesLiteral    = "{ap:(?i)accessprofiles}"
	AdminLiteral             = "{ad:(?i)admin}"
	PodsLiteral              = "{po:(?i)pods}"
	LogLiteral               = "{lo:(?i)log}"
	EventLiteral             = "{ev:(?i)events}"
	KubectlLiteral           = "{ku:(?i)kubectl}"
	ContainersLiteral        = "{co:(?i)containers}"
	UnderlaysLiteral         = "{un:(?i)underlays}"
	DefaultLiteral           = "{up:(?i)default}"
	ListCredentialLiteral    = "{li:(?i)listcredential}"
)

const (
	// PathSubscriptionIDParameter is the path parameter name used in routing for the subscription id
	PathSubscriptionIDParameter = "subscriptionId"
	// PathResourceGroupNameParameter is the path parameter name used in routing for the resource group name
	PathResourceGroupNameParameter = "resourceGroupName"
	// PathResourceNameParameter is the path parameter name used in routing for the resource name
	PathResourceNameParameter = "resourceName"
	// RequestAPIVersionParameterName is the query string parameter name ARM adds for the api version
	RequestAPIVersionParameterName = "api-version"
	// TerraformRPNamespace is the ARM namespace for Terraform RP
	TerraformRPNamespace = "Microsoft.Terraform-OSS"
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
	// /{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Terraform-OSS/resources/{resourceName}
	ResourceOperationRoute = SubscriptionResourceOperationRoute + "/" + ResourceGroupsLiteral + "/{" +
		PathResourceGroupNameParameter +
		"}/" + ProvidersLiteral + "/" + TerraformRPNamespace + "/" + ResourcesLiteral + "/{" +
		PathResourceNameParameter + "}"
)

const (
	// GetResourceControllerName is the constant logged for get resource calls
	GetResourceControllerName = "GetResourceController"

	// PutResourceControllerName is the constant logged for get resource calls
	PutResourceControllerName = "PutResourceController"

	// DeleteResourceControllerName is the constant logged for get resource calls
	DeleteResourceControllerName = "DeleteResourceController"
)
