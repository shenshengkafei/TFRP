package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestValidateArmStorageAccountType(t *testing.T) {
	testCases := []struct {
		input       string
		shouldError bool
	}{
		{"standard_lrs", false},
		{"invalid", true},
	}

	for _, test := range testCases {
		_, es := validateArmStorageAccountType(test.input, "account_type")

		if test.shouldError && len(es) == 0 {
			t.Fatalf("Expected validating account_type %q to fail", test.input)
		}
	}
}

func TestValidateArmStorageAccountName(t *testing.T) {
	testCases := []struct {
		input       string
		shouldError bool
	}{
		{"ab", true},
		{"ABC", true},
		{"abc", false},
		{"123456789012345678901234", false},
		{"1234567890123456789012345", true},
		{"abc12345", false},
	}

	for _, test := range testCases {
		_, es := validateArmStorageAccountName(test.input, "name")

		if test.shouldError && len(es) == 0 {
			t.Fatalf("Expected validating name %q to fail", test.input)
		}
	}
}

func TestAccAzureRMStorageAccount_basic(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	preConfig := testAccAzureRMStorageAccount_basic(ri, rs)
	postConfig := testAccAzureRMStorageAccount_update(ri, rs)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "account_type", "Standard_LRS"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "tags.%", "1"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "tags.environment", "production"),
				),
			},

			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "account_type", "Standard_GRS"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "tags.%", "1"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "tags.environment", "staging"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_disappears(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	preConfig := testAccAzureRMStorageAccount_basic(ri, rs)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "account_type", "Standard_LRS"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "tags.%", "1"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "tags.environment", "production"),
					testCheckAzureRMStorageAccountDisappears("azurerm_storage_account.testsa"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMStorageAccount_blobConnectionString(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	preConfig := testAccAzureRMStorageAccount_basic(ri, rs)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttrSet("azurerm_storage_account.testsa", "primary_blob_connection_string"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_blobEncryption(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	preConfig := testAccAzureRMStorageAccount_blobEncryption(ri, rs)
	postConfig := testAccAzureRMStorageAccount_blobEncryptionDisabled(ri, rs)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "enable_blob_encryption", "true"),
				),
			},

			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "enable_blob_encryption", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_enableHttpsTrafficOnly(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	preConfig := testAccAzureRMStorageAccount_enableHttpsTrafficOnly(ri, rs)
	postConfig := testAccAzureRMStorageAccount_enableHttpsTrafficOnlyDisabled(ri, rs)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "enable_https_traffic_only", "true"),
				),
			},

			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "enable_https_traffic_only", "false"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_blobStorageWithUpdate(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	preConfig := testAccAzureRMStorageAccount_blobStorage(ri, rs)
	postConfig := testAccAzureRMStorageAccount_blobStorageUpdate(ri, rs)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "account_kind", "BlobStorage"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "access_tier", "Hot"),
				),
			},

			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
					resource.TestCheckResourceAttr("azurerm_storage_account.testsa", "access_tier", "Cool"),
				),
			},
		},
	})
}

func TestAccAzureRMStorageAccount_NonStandardCasing(t *testing.T) {
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	preConfig := testAccAzureRMStorageAccountNonStandardCasing(ri, rs)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMStorageAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMStorageAccountExists("azurerm_storage_account.testsa"),
				),
			},

			{
				Config:             preConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testCheckAzureRMStorageAccountExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		storageAccount := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		// Ensure resource group exists in API
		conn := testAccProvider.Meta().(*ArmClient).storageServiceClient

		resp, err := conn.GetProperties(resourceGroup, storageAccount)
		if err != nil {
			return fmt.Errorf("Bad: Get on storageServiceClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: StorageAccount %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMStorageAccountDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		storageAccount := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		// Ensure resource group exists in API
		conn := testAccProvider.Meta().(*ArmClient).storageServiceClient

		_, err := conn.Delete(resourceGroup, storageAccount)
		if err != nil {
			return fmt.Errorf("Bad: Delete on storageServiceClient: %s", err)
		}

		return nil
	}
}

func testCheckAzureRMStorageAccountDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).storageServiceClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_storage_account" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.GetProperties(resourceGroup, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Storage Account still exists:\n%#v", resp.AccountProperties)
		}
	}

	return nil
}

func testAccAzureRMStorageAccount_basic(rInt int, rString string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "westus"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "westus"
    account_type = "Standard_LRS"

    tags {
        environment = "production"
    }
}`, rInt, rString)
}

func testAccAzureRMStorageAccount_update(rInt int, rString string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "westus"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "westus"
    account_type = "Standard_GRS"

    tags {
        environment = "staging"
    }
}`, rInt, rString)
}

func testAccAzureRMStorageAccount_blobEncryption(rInt int, rString string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "westus"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "westus"
    account_type = "Standard_LRS"
    enable_blob_encryption = true

    tags {
        environment = "production"
    }
}`, rInt, rString)
}

func testAccAzureRMStorageAccount_blobEncryptionDisabled(rInt int, rString string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "westus"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "westus"
    account_type = "Standard_LRS"
    enable_blob_encryption = false

    tags {
        environment = "production"
    }
}`, rInt, rString)
}

func testAccAzureRMStorageAccount_enableHttpsTrafficOnly(rInt int, rString string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "westus"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "westus"
    account_type = "Standard_LRS"
    enable_https_traffic_only = true

    tags {
        environment = "production"
    }
}`, rInt, rString)
}

func testAccAzureRMStorageAccount_enableHttpsTrafficOnlyDisabled(rInt int, rString string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "westus"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "westus"
    account_type = "Standard_LRS"
    enable_https_traffic_only = false

    tags {
        environment = "production"
    }
}`, rInt, rString)
}

// BlobStorage accounts are not available in WestUS
func testAccAzureRMStorageAccount_blobStorage(rInt int, rString string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "northeurope"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "northeurope"
	account_kind = "BlobStorage"
    account_type = "Standard_LRS"

    tags {
        environment = "production"
    }
}
`, rInt, rString)
}

func testAccAzureRMStorageAccount_blobStorageUpdate(rInt int, rString string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "northeurope"
}

resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"

    location = "northeurope"
	account_kind = "BlobStorage"
    account_type = "Standard_LRS"
	access_tier = "Cool"

    tags {
        environment = "production"
    }
}
`, rInt, rString)
}

func testAccAzureRMStorageAccountNonStandardCasing(ri int, rs string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "testrg" {
    name = "testAccAzureRMSA-%d"
    location = "westus"
}
resource "azurerm_storage_account" "testsa" {
    name = "unlikely23exst2acct%s"
    resource_group_name = "${azurerm_resource_group.testrg.name}"
    location = "westus"
    account_type = "standard_LRS"
    tags {
        environment = "production"
    }
}`, ri, rs)
}
