package azurerm

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func buildConfigForSweepers() (*ArmClient, error) {
	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	clientID := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	tenantID := os.Getenv("ARM_TENANT_ID")
	environment := os.Getenv("ARM_ENVIRONMENT")

	if environment == "" {
		environment = "public"
	}

	if subscriptionID == "" || clientID == "" || clientSecret == "" || tenantID == "" {
		return nil, fmt.Errorf("ARM_SUBSCRIPTION_ID, ARM_CLIENT_ID, ARM_CLIENT_SECRET and ARM_TENANT_ID must be set for acceptance tests")
	}

	config := &Config{
		SubscriptionID:           subscriptionID,
		ClientID:                 clientID,
		ClientSecret:             clientSecret,
		TenantID:                 tenantID,
		Environment:              environment,
		SkipProviderRegistration: false,
	}

	return config.getArmClient()
}

func shouldSweepAcceptanceTestResource(name string, resourceLocation string, region string) bool {
	loweredName := strings.ToLower(name)

	prefixesToIgnore := []string{"acctest"}

	for _, prefix := range prefixesToIgnore {
		if !strings.HasPrefix(loweredName, prefix) {
			log.Printf("Ignoring Resource '%s' due to prefix '%s'", name, prefix)
			return false
		}
	}

	normalisedResourceLocation := azureRMNormalizeLocation(resourceLocation)
	normalisedRegion := azureRMNormalizeLocation(region)

	if normalisedResourceLocation != normalisedRegion {
		log.Printf("Region '%s' isn't '%s' - skipping", normalisedResourceLocation, normalisedRegion)
		return false
	}

	return true
}
