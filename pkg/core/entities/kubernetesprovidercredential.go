//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package entities

// KubernetesProviderCredential is the kubernetes provider credential
type KubernetesProviderCredential struct {
	InlineConfig string `json:"inline_config"`
}
