//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package storage

import (
	"TFRP/pkg/core/consts"
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

// Insert inserts the data into collection
func (providerRegistrationDataProvider *ProviderRegistrationDataProvider) Insert(doc interface{}) error {
	return providerRegistrationDataProvider.baseDataProvider.Insert(consts.ProviderRegistrationCollectionName, doc)
}

// Find returns the data
func (providerRegistrationDataProvider *ProviderRegistrationDataProvider) Find(qurey interface{}, result interface{}) error {
	return providerRegistrationDataProvider.baseDataProvider.Find(consts.ProviderRegistrationCollectionName, qurey, result)
}
