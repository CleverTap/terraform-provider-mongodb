package mongodb

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	os.Setenv("TF_ACC", "1")
	os.Setenv("PUBLIC_KEY", "PUBLIC_KEY value")
	os.Setenv("PRIVATE_KEY", "PRIVATE_KEY value")
	os.Setenv("ORGID", "ORGID value")
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"mongodb": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("PUBLIC_KEY"); v == "" {
		t.Fatal("PUBLIC_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("PRIVATE_KEY"); v == "" {
		t.Fatal("PRIVATE_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ORGID"); v == "" {
		t.Fatal("ORGID must be set for acceptance tests")
	}
}
