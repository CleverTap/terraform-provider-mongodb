package mongodb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.mongodb_user.tapendra", "emailaddress", "user@gmail.com"),
					resource.TestCheckResourceAttr("data.mongodb_user.tapendra", "firstname", "FirstName"),
					resource.TestCheckResourceAttr("data.mongodb_user.tapendra", "lastname", "LastName"),
					resource.TestCheckResourceAttr("data.mongodb_user.tapendra", "country", "IN"),
				),
			},
		},
	})
}

func testAccCheckUserDataExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		return nil
	}
}
func testAccUserDataSourceConfig() string {
	return fmt.Sprintf(`
	data "mongodb_user" "tapendra" {
		id = "user@gmail.com"
	  }`)
}
