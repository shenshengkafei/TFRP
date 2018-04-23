//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package storage

import (
	"TFRP/pkg/core/consts"
)

// ResourceDataProvider is the base struc of all data provider
type ResourceDataProvider struct {
	baseDataProvider BaseDataProvider
}

// GetResourceDataProvider returns the resource data provider
func GetResourceDataProvider() *ResourceDataProvider {
	return &ResourceDataProvider{
		baseDataProvider: BaseDataProvider{
			Database: consts.StorageDatabase,
			Password: consts.StoragePassword,
		},
	}
}

// Insert inserts the data into collection
func (resourceDataProvider *ResourceDataProvider) Insert(doc interface{}) error {
	return resourceDataProvider.baseDataProvider.Insert(consts.ResourceCollectionName, doc)
}

// Find returns the data
func (resourceDataProvider *ResourceDataProvider) Find(qurey interface{}, result interface{}) error {
	return resourceDataProvider.baseDataProvider.Find(consts.ResourceCollectionName, qurey, result)
}
