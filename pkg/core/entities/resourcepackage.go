//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package entities

import "gopkg.in/mgo.v2/bson"

// ResourcePackage is the package stored in storag
type ResourcePackage struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	ResourceID   string
	StateID      string
	State        interface{}
	Config       string
	ResourceType string
	ProviderType string
}

// ResourcePackageDefinition is the package definition
type ResourcePackageDefinition struct {
	Properties interface{}
}

// ToDefinition returns the definition
func (resourcePackage *ResourcePackage) ToDefinition() *ResourcePackageDefinition {
	return &ResourcePackageDefinition{
		Properties: resourcePackage,
	}
}
