//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package entities

// ProviderRegistrationDefinition is provider registration definition
type ProviderRegistrationDefinition struct {
	Location   string
	Properties PoviderRegistrationProperties
}

// PoviderRegistrationProperties is provider registration properties
type PoviderRegistrationProperties struct {
	ProviderType string
	Settings     interface{}
}
