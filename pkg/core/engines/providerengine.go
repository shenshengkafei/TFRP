//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package engines

import (
	"TFRP/datadog"
	"TFRP/kubernetes"
	"TFRP/pkg/core/consts"

	"github.com/hashicorp/terraform/helper/schema"
)

// GetProvider returns the provider
func GetProvider(providerType string) *schema.Provider {
	switch providerType {
	case consts.KubernetesProvider:
		return kubernetes.Provider().(*schema.Provider)
	case consts.DatadogProvider:
		return datadog.Provider().(*schema.Provider)
	}

	return nil
}
