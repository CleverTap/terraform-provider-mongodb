package mongodb

import (
	"terraform-provider-mongodb/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"mongodb_public_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MONGODBCLOUD_PUBLIC_KEY", ""),
			},
			"mongodb_private_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MONGODBCLOUD_PRIVATE_KEY", ""),
			},
			"mongodb_orgid": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MONGODBCLOUD_ORGID", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"mongodb_user": resourceUser(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"mongodb_user": dataSourceUser(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	publickey := d.Get("mongodb_public_key").(string)
	privatekey := d.Get("mongodb_private_key").(string)
	orgid := d.Get("mongodb_orgid").(string)
	return client.NewClient(publickey, privatekey, orgid), nil

}
