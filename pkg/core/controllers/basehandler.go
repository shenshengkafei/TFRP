//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package controllers

import (
	"TFRP/pkg/core/storage"
)

// BaseHandler is the base handler
type BaseHandler struct {
	ProviderRegistrationDataProvider *storage.ProviderRegistrationDataProvider
	ResourceDataProvider             *storage.ResourceDataProvider
}
