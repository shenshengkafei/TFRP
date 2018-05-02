//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package entities

import (
	"github.com/hashicorp/terraform/terraform"
	"gopkg.in/mgo.v2/bson"
)

// ResourcePackage is the package stored in storag
type ResourcePackage struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	ResourceID   string
	StateID      string
	State        *terraform.InstanceState
	Config       string
	ResourceType string
	ProviderType string
}

// ResourcePackageDefinition is the package definition
type ResourcePackageDefinition struct {
	Properties ResourcePackage
}

// ToDefinition returns the definition
func (resourcePackage *ResourcePackage) ToDefinition() *ResourcePackageDefinition {
	resourcePackage.Config = ""
	return &ResourcePackageDefinition{
		Properties: ResourcePackage{
			ID:           resourcePackage.ID,
			ResourceID:   resourcePackage.ResourceID,
			StateID:      resourcePackage.StateID,
			State:        resourcePackage.State,
			ResourceType: resourcePackage.ResourceType,
			ProviderType: resourcePackage.ProviderType,
		},
	}
}
