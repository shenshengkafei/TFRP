//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package entities

// ResourceDefinition is the resource definition
type ResourceDefinition struct {
	Location   string
	Properties *ResourceDefinitionProperties
}

// ResourceDefinitionProperties is the resouce definition properties
type ResourceDefinitionProperties struct {
	ProviderID   string
	ResourceType string
	Settings     interface{}
}
