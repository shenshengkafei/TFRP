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
	BaseDataProvider
}

// NewResourceDataProvider creates a new resource data provider
func NewResourceDataProvider(database, password string) (resourceDataProvider *ResourceDataProvider) {
	resourceDataProvider = new(ResourceDataProvider)
	resourceDataProvider.Database = database
	resourceDataProvider.Password = password
	return resourceDataProvider
}

// InsertPackage inserts a doc into collection
func (resourceDataProvider *ResourceDataProvider) InsertPackage(doc *entities.ResourcePackage) error {
	return resourceDataProvider.Insert(consts.ResourceCollectionName, bson.M{"resourceid": doc.ResourceID}, doc)
}

// FindPackage returns a doc from colletion
func (resourceDataProvider *ResourceDataProvider) FindPackage(resourceID string, result interface{}) error {
	return resourceDataProvider.Find(consts.ResourceCollectionName, bson.M{"resourceid": resourceID}, result)
}

// RemovePackage deletes a doc from collection
func (resourceDataProvider *ResourceDataProvider) RemovePackage(resourceID string) error {
	return resourceDataProvider.Remove(consts.ResourceCollectionName, bson.M{"resourceid": resourceID})
}
