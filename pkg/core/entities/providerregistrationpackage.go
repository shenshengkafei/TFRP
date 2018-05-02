//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package entities

import "gopkg.in/mgo.v2/bson"

// ProviderRegistrationPackage is the package stored in storage
type ProviderRegistrationPackage struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	ResourceID   string        `json:",omitempty"`
	ProviderType string        `json:",omitempty"`
	Settings     []byte        `json:",omitempty"`
}

// ProviderRegistrationPackageDefinition is the package definition
type ProviderRegistrationPackageDefinition struct {
	Properties ProviderRegistrationPackage
}

// ToDefinition returns the definition
func (providerRegistrationPackage *ProviderRegistrationPackage) ToDefinition() *ProviderRegistrationPackageDefinition {
	return &ProviderRegistrationPackageDefinition{
		Properties: ProviderRegistrationPackage{
			ID:           providerRegistrationPackage.ID,
			ResourceID:   providerRegistrationPackage.ResourceID,
			ProviderType: providerRegistrationPackage.ProviderType,
		},
	}
}

// ToListSettingsDefinition returns the definition
func (providerRegistrationPackage *ProviderRegistrationPackage) ToListSettingsDefinition() *ProviderRegistrationPackageDefinition {
	return &ProviderRegistrationPackageDefinition{
		Properties: ProviderRegistrationPackage{
			Settings: providerRegistrationPackage.Settings,
		},
	}
}
