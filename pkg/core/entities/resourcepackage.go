//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package entities

import "gopkg.in/mgo.v2/bson"

// ResourcePackage is the package stored in storag
type ResourcePackage struct {
	Id           bson.ObjectId `bson:"_id,omitempty"`
	ResourceID   string
	StateID      string
	Config       string
	ResourceType string
	ProviderType string
}
