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
	ResourceID   string        `json:",omitempty"`
	StateID      string        `json:",omitempty"`
	State        *terraform.InstanceState
	Config       string `json:",omitempty"`
	ResourceType string `json:",omitempty"`
	ProviderType string `json:",omitempty"`
}

// ResourcePackageDefinition is the package definition
type ResourcePackageDefinition struct {
	Properties ResourcePackage
}

// ToDefinition returns the definition
func (resourcePackage *ResourcePackage) ToDefinition() *ResourcePackageDefinition {
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
