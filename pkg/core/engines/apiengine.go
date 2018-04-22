//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

// This file is for code about storing and retrieving api tracking
// info from a context struct

package engines

import (
	"TFRP/pkg/core/consts"

	restful "github.com/emicklei/go-restful"
)

const (
	// Subscriptions identifies a customer/an organization
	Subscriptions = "subscriptions"
	// ResourceGroups contains a group of ARM (RP) resources
	ResourceGroups = "resourceGroups"
	// Resource is the name of resource
	Resource = "resource"
	// Providers is the name of provider
	Providers = "providers"
)

// GetSubscriptionID returns the subscription id if it was on the request else empty string
func GetSubscriptionID(request *restful.Request) string {
	return request.PathParameter(consts.PathSubscriptionIDParameter)
}

// GetAPIVersion returns the api version on the request
func GetAPIVersion(request *restful.Request) string {
	return request.QueryParameter(consts.RequestAPIVersionParameterName)
}

// GetResourceGroupName returns the resourceGroupName if it was on the request else empty string
func GetResourceGroupName(request *restful.Request) string {
	return request.PathParameter(consts.PathResourceGroupNameParameter)
}

// GetResourceName returns the resourceName if it was on the request else empty string
func GetResourceName(request *restful.Request) string {
	return request.PathParameter(consts.PathResourceNameParameter)
}

// GetFullyQualifiedResourceID returns the fully qualified resource id
// "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroup}/providers/Microsoft.Terraform-OSS/resource/{resource}"
func GetFullyQualifiedResourceID(request *restful.Request) string {
	return "/" + Subscriptions + "/" + GetSubscriptionID(request) +
		"/" + ResourceGroups + "/" + GetResourceGroupName(request) +
		"/" + Providers + "/" + consts.TerraformRPNamespace +
		"/" + Resource + "/" + GetResourceName(request)
}
