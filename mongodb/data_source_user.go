package mongodb

import (
	"fmt"
	"strings"
	"terraform-provider-mongodb/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"country": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"emailaddress": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"firstname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"lastname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"roles": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceUserRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	userId := d.Get("id").(string)
	user, err := apiClient.GetUser(userId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding Item with ID %s", userId)
		}
	}

	d.SetId(user.ID)
	d.Set("country", user.Country)
	d.Set("id", user.ID)
	d.Set("emailaddress", user.EmailAddress)
	d.Set("firstname", user.FirstName)
	d.Set("lastname", user.LastName)
	roles := make([]string, 0)
	for _, v := range user.Roles {
		if v.OrgID != "" {
			roles = append(roles, v.RoleName)
		}
	}
	d.Set("roles", roles)
	d.Set("username", user.Username)

	return nil
}
