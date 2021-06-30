package mongodb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUser_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserBasic(),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("mongodb_user.testuser", "username", "user@gmail.com"),
					resource.TestCheckResourceAttr("mongodb_user.testuser", "roles.#", "1"),
					resource.TestCheckResourceAttr("mongodb_user.testuser", "roles", "ORG_MEMBER"),
				),
			},
		},
	})
}

func testAccCheckUserBasic() string {
	return fmt.Sprintf(`
	resource "mongodb_user" "testuser" {
		username="user@gmail.com"
		roles = [ "ORG_MEMBER"]
	  }
	`)
}

func TestAccUser_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckItemUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					// testAccCheckUserDataExists("mongodb_user.testUser"),
					resource.TestCheckResourceAttr("mongodb_user.testuser", "username", "user@gmail.com"),
					resource.TestCheckResourceAttr("mongodb_user.testuser", "roles.#", "1"),
					resource.TestCheckResourceAttr("mongodb_user.testuser", "roles.0", "ORG_MEMBER"),
				),
			},
			{
				Config: testAccCheckItemUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					// testAccCheckUserDataExists("mongodb_user.testUser"),
					resource.TestCheckResourceAttr("mongodb_user.testuser", "username", "user@gmail.com"),
					resource.TestCheckResourceAttr("mongodb_user.testuser", "roles.#", "1"),
					resource.TestCheckResourceAttr("mongodb_user.testuser", "roles.0", "ORG_READ_ONLY"),
				),
			},
		},
	})
}

func testAccCheckItemUpdatePre() string {
	return fmt.Sprintf(`
	resource "mongodb_user" "testuser" {
		username="user@gmail.com"
		roles = [ "ORG_MEMBER"]
		}
	`)
}

func testAccCheckItemUpdatePost() string {
	return fmt.Sprintf(`
	resource "mongodb_user" "testuser" {
		username="user@gmail.com"
		roles = [ "ORG_READ_ONLY"]
		}
	`)
}
