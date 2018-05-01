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
	BaseDataProvider
}

// NewProviderRegistrationDataProvider creates a new provider registration data provider
func NewProviderRegistrationDataProvider(database, password string) (providerRegistrationDataProvider *ProviderRegistrationDataProvider) {
	providerRegistrationDataProvider = new(ProviderRegistrationDataProvider)
	providerRegistrationDataProvider.Database = database
	providerRegistrationDataProvider.Password = password
	return providerRegistrationDataProvider
}

// InsertPackage inserts a doc into collection
func (providerRegistrationDataProvider *ProviderRegistrationDataProvider) InsertPackage(doc *entities.ProviderRegistrationPackage) error {
	return providerRegistrationDataProvider.Insert(consts.ProviderRegistrationCollectionName, bson.M{"resourceid": doc.ResourceID}, doc)
}

// FindPackage returns a doc from collection
func (providerRegistrationDataProvider *ProviderRegistrationDataProvider) FindPackage(resourceID string, result interface{}) error {
	return providerRegistrationDataProvider.Find(consts.ProviderRegistrationCollectionName, bson.M{"resourceid": resourceID}, result)
}

// RemovePackage deletes a doc from collection
func (providerRegistrationDataProvider *ProviderRegistrationDataProvider) RemovePackage(resourceID string) error {
	return providerRegistrationDataProvider.Remove(consts.ProviderRegistrationCollectionName, bson.M{"resourceid": resourceID})
}
