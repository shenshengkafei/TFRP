//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package engines

import (
	"TFRP/pkg/core/apierror"
	"TFRP/pkg/core/consts"
	"TFRP/pkg/core/entities"
	"fmt"
	"strings"
)

// ValidateProviderRegistrationDefinition validates the provider registration definition
func ValidateProviderRegistrationDefinition(providerRegistrationDefinition *entities.ProviderRegistrationDefinition) *apierror.ErrorResponse {
	if providerRegistrationDefinition.Properties == nil {
		return apierror.New(
			apierror.ClientError,
			apierror.BadRequest,
			fmt.Sprintf("Request content is missing properties."))
	}
	if len(strings.TrimSpace(providerRegistrationDefinition.Properties.ProviderType)) == 0 {
		return apierror.New(
			apierror.ClientError,
			apierror.BadRequest,
			fmt.Sprintf("Request content is missing property 'providerType'."))
	}

	supportedProviders := []string{consts.KubernetesProvider, consts.DatadogProvider, consts.CloudflareProvider}
	isSupported := false
	for _, provider := range supportedProviders {
		if strings.EqualFold(provider, providerRegistrationDefinition.Properties.ProviderType) {
			isSupported = true
			break
		}
	}

	if !isSupported {
		return apierror.New(
			apierror.ClientError,
			apierror.BadRequest,
			fmt.Sprintf("The provider type %s is not supported. Supported providers are %s.", providerRegistrationDefinition.Properties.ProviderType, supportedProviders))
	}

	return nil
}
