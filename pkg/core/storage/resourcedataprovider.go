//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package storage

import (
	"TFRP/pkg/core/consts"
	"TFRP/pkg/core/entities"

	"gopkg.in/mgo.v2/bson"
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

// Insert inserts a doc into collection
func (resourceDataProvider *ResourceDataProvider) Insert(doc *entities.ResourcePackage) error {
	return resourceDataProvider.baseDataProvider.Insert(consts.ResourceCollectionName, bson.M{"resourceid": doc.ResourceID}, doc)
}

// Find returns a doc from colletion
func (resourceDataProvider *ResourceDataProvider) Find(resourceID string, result interface{}) error {
	return resourceDataProvider.baseDataProvider.Find(consts.ResourceCollectionName, bson.M{"resourceid": resourceID}, result)
}

// Remove deletes a doc from collection
func (resourceDataProvider *ResourceDataProvider) Remove(resourceID string) error {
	return resourceDataProvider.baseDataProvider.Remove(consts.ResourceCollectionName, bson.M{"resourceid": resourceID})
}
