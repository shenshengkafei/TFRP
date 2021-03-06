//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package entities

import (
	"TFRP/pkg/core/consts"

	"github.com/hashicorp/terraform/terraform"
	"gopkg.in/mgo.v2/bson"
)

// ResourcePackage is the package stored in storag
type ResourcePackage struct {
	ID                       bson.ObjectId `bson:"_id,omitempty"`
	Location                 string        `json:",omitempty"`
	ResourceID               string        `json:",omitempty"`
	StateID                  string        `json:",omitempty"`
	State                    *terraform.InstanceState
	ProvisioningState        string `json:",omitempty"`
	ProvisioningErrorCode    string `json:",omitempty"`
	ProvisioningErrorMessage string `json:",omitempty"`
	Config                   string `json:",omitempty"`
	ResourceType             string `json:",omitempty"`
	ProviderType             string `json:",omitempty"`
}

// ResourcePackageDefinition is the package definition
type ResourcePackageDefinition struct {
	Type       string
	Location   string
	Properties ResourcePackage
}

// ToDefinition returns the definition
func (resourcePackage *ResourcePackage) ToDefinition() *ResourcePackageDefinition {
	return &ResourcePackageDefinition{
		Location: resourcePackage.Location,
		Type:     consts.TerraformResourceType,
		Properties: ResourcePackage{
			ID:                resourcePackage.ID,
			ResourceID:        resourcePackage.ResourceID,
			StateID:           resourcePackage.StateID,
			State:             resourcePackage.State,
			ProvisioningState: resourcePackage.ProvisioningState,
			ResourceType:      resourcePackage.ResourceType,
			ProviderType:      resourcePackage.ProviderType,
		},
	}
}

// ToAsyncOperationResult returns the AsyncOperationResult
func (resourcePackage *ResourcePackage) ToAsyncOperationResult() *AsyncOperationResult {
	return &AsyncOperationResult{
		Status: resourcePackage.ProvisioningState,
		Error: &ExtendedErrorInfo{
			Code:    resourcePackage.ProvisioningErrorCode,
			Message: resourcePackage.ProvisioningErrorMessage,
		},
	}
}
