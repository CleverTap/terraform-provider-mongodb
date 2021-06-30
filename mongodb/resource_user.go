package mongodb

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"terraform-provider-mongodb/client"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validateEmail(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value := v.(string)

	var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	if !(emailRegex.MatchString(value)) {
		errs = append(errs, fmt.Errorf("Expected EmailId is not valid  %s", k))
		return warns, errs
	}
	return
}

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceUserImporter,
		},
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateEmail,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"invitationid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"userid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"inviterusername": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"orgname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"country": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"email_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"roles": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	userRoles := d.Get("roles").(*schema.Set).List()
	roles := make([]string, len(userRoles))

	for i, role := range userRoles {
		roles[i] = role.(string)
	}
	user := client.NewUser{
		Username: d.Get("username").(string),
		Roles:    roles,
	}

	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		newuser, err := apiClient.AddNewUser(&user)
		if err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		d.SetId(newuser.Username)
		d.Set("inviterusername", newuser.InviterUsername)
		d.Set("orgname", newuser.OrgName)
		d.Set("username", newuser.Username)
		return nil
	})
	if retryErr != nil {
		if strings.Contains(retryErr.Error(), "Unautharized Access") == true {
			d.SetId("")
			return diags
		}
		return diag.FromErr(retryErr)
	}
	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	if os.Getenv("ORGID") != "" {
		return diags
	}
	apiClient := m.(*client.Client)
	userId := d.Id()

	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		user, err := apiClient.GetUser(userId)
		if err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if len(user.ID) > 0 {
			fmt.Println("Set all values")
			d.SetId(user.Username)
			d.Set("userid", user.ID)
			d.Set("country", user.Country)
			d.Set("email_address", user.EmailAddress)
			d.Set("first_name", user.FirstName)
			d.Set("last_name", user.LastName)
			roles := make([]string, 0)
			for _, v := range user.Roles {
				if v.OrgID != "" {
					roles = append(roles, v.RoleName)
				}
			}
			d.Set("roles", roles)
			d.Set("username", user.Username)
		}
		return nil
	})
	if retryErr != nil {
		if strings.Contains(retryErr.Error(), "User Does Not Exist") == true {
			d.SetId("")
			return diags
		}
		return diag.FromErr(retryErr)
	}
	return diags

}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var _ diag.Diagnostics
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	if d.HasChange("username") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "User not allowed to change username",
			Detail:   "User not allowed to change username",
		})

		return diags
	}
	userId := d.Get("username").(string)
	userInfo, error := apiClient.GetUser(userId)
	if error != nil {
		log.Println("[ERROR]: ", error)
		if strings.Contains(error.Error(), "not found") {
			d.SetId("")
		} else {
			return diag.FromErr(error)
		}
	}
	userId = userInfo.ID
	userRoles := d.Get("roles").(*schema.Set).List()
	roles := make([]string, len(userRoles))

	for i, role := range userRoles {
		roles[i] = role.(string)
	}

	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, error = apiClient.UpdateUser(roles, userId)
		if error != nil {
			if apiClient.IsRetry(error) {
				return resource.RetryableError(error)
			}
			return resource.NonRetryableError(error)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}
	if error != nil {
		return diag.FromErr(error)
	}
	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	var err error
	userId := d.Get("username").(string)
	userInfo, error := apiClient.GetUser(userId)
	if error != nil {
		log.Println("[ERROR]: ", error)
		if strings.Contains(error.Error(), "not found") {
			d.SetId("")
		} else {
			return diag.FromErr(error)
		}
	}
	userId = userInfo.ID

	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		if err = apiClient.DeleteUser(userId); err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceUserImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	apiClient := m.(*client.Client)
	userId := d.Id()
	user, err := apiClient.GetUser(userId)
	if err != nil {
		return nil, err
	}
	if len(user.ID) > 0 {
		fmt.Println("Set all values")
		d.SetId(user.Username)
		d.Set("userid", user.ID)
		d.Set("country", user.Country)
		d.Set("email_address", user.EmailAddress)
		d.Set("first_name", user.FirstName)
		d.Set("last_name", user.LastName)
		roles := make([]string, 0)
		for _, v := range user.Roles {
			if v.OrgID != "" {
				roles = append(roles, v.RoleName)
			}
		}

		d.Set("roles", roles)
		d.Set("username", user.Username)

	}
	return []*schema.ResourceData{d}, nil
}
