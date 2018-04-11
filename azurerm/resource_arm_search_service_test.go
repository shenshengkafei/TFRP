package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jen20/riviera/search"
)

func TestAccAzureRMSearchService_basic(t *testing.T) {
	resourceName := "azurerm_search_service.test"
	ri := acctest.RandInt()
	config := testAccAzureRMSearchService_basic(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSearchServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSearchServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMSearchService_updateReplicaCountAndTags(t *testing.T) {
	resourceName := "azurerm_search_service.test"
	ri := acctest.RandInt()
	preConfig := testAccAzureRMSearchService_basic(ri)
	postConfig := testAccAzureRMSearchService_updated(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSearchServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSearchServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "replica_count", "1"),
				),
			},

			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSearchServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "replica_count", "2"),
				),
			},
		},
	})
}

func testCheckAzureRMSearchServiceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).rivieraClient

		readRequest := conn.NewRequestForURI(rs.Primary.ID)
		readRequest.Command = &search.GetSearchService{}

		readResponse, err := readRequest.Execute()
		if err != nil {
			return fmt.Errorf("Bad: GetSearchService: %+v", err)
		}
		if !readResponse.IsSuccessful() {
			return fmt.Errorf("Bad: GetSearchService: %+v", readResponse.Error)
		}

		return nil
	}
}

func testCheckAzureRMSearchServiceDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).rivieraClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_search_service" {
			continue
		}

		readRequest := conn.NewRequestForURI(rs.Primary.ID)
		readRequest.Command = &search.GetSearchService{}

		readResponse, err := readRequest.Execute()
		if err != nil {
			return fmt.Errorf("Bad: GetSearchService: %+v", err)
		}

		if readResponse.IsSuccessful() {
			return fmt.Errorf("Bad: Search Service still exists: %+v", readResponse.Error)
		}
	}

	return nil
}

func testAccAzureRMSearchService_basic(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}

resource "azurerm_search_service" "test" {
    name = "acctestsearchservice%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    sku = "standard"

    tags {
    	environment = "staging"
    	database = "test"
    }
}
`, rInt, rInt)
}

func testAccAzureRMSearchService_updated(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestRG_%d"
    location = "West US"
}
resource "azurerm_search_service" "test" {
    name = "acctestsearchservice%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    location = "${azurerm_resource_group.test.location}"
    sku = "standard"
    replica_count = 2

    tags {
    	environment = "production"
    }
}
`, rInt, rInt)
}
