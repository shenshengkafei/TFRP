//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package entities

import "gopkg.in/mgo.v2/bson"

// ProviderRegistrationPackage is the package stored in storage
type ProviderRegistrationPackage struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	ResourceID   string
	ProviderType string
	Credentials  []byte
}
