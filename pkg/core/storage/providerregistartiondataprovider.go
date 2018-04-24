//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package storage

import (
	"TFRP/pkg/core/consts"
	"TFRP/pkg/core/entities"

	"gopkg.in/mgo.v2/bson"
)

// ProviderRegistrationDataProvider is the base struc of all data provider
type ProviderRegistrationDataProvider struct {
	baseDataProvider BaseDataProvider
}

// GetProviderRegistrationDataProvider returns the provider registration provider
func GetProviderRegistrationDataProvider() *ProviderRegistrationDataProvider {
	return &ProviderRegistrationDataProvider{
		baseDataProvider: BaseDataProvider{
			Database: consts.StorageDatabase,
			Password: consts.StoragePassword,
		},
	}
}

// Insert inserts a doc into collection
func (providerRegistrationDataProvider *ProviderRegistrationDataProvider) Insert(doc *entities.ProviderRegistrationPackage) error {
	return providerRegistrationDataProvider.baseDataProvider.Insert(consts.ProviderRegistrationCollectionName, bson.M{"resourceid": doc.ResourceID}, doc)
}

// Find returns a doc from collection
func (providerRegistrationDataProvider *ProviderRegistrationDataProvider) Find(resourceID string, result interface{}) error {
	return providerRegistrationDataProvider.baseDataProvider.Find(consts.ProviderRegistrationCollectionName, bson.M{"resourceid": resourceID}, result)
}

// Remove deletes a doc from collection
func (providerRegistrationDataProvider *ProviderRegistrationDataProvider) Remove(resourceID string) error {
	return providerRegistrationDataProvider.baseDataProvider.Remove(consts.ProviderRegistrationCollectionName, bson.M{"resourceid": resourceID})
}
