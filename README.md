> Delete, Activate, Deactivate API's are not supported.

- This terraform providre allows to perform Create ,Read ,Update ,Delete and Import MongoDB user

## Requirements

- [Go Lang](https://golang.org/doc/install)>=1.11 <br>
- [Terraform](https://www.terraform.io/downloads.html)>=v0.13.0 <br/>
- MongoDB Basic Account
- [MongoDB API Credentials](https://docs.atlas.mongodb.com/configure-api-access/)
- [MongoDB API documentation](https://docs.cloudmanager.mongodb.com/reference/api/users/)

## Application Account
This provider can be successfully tested on a basic mongoDB cloud manager account.

## Setup 
1. Create a MongoDB account [here](https://account.mongodb.com/account/)

## API Authentication
1. Login  to mongodb account
2. Click the Organization Settings icon next to the Organizations menu.
3. Toggle the Require IP Access List for Public API setting to On
4. Click Access Manager in the sidebar, or click Access Manager in the navigation bar, then click your organization.
5. Click Create API Key.
6. Enter the API Key Information
7. Click Next and save the Public and Private key.


## Building The Provider
1. Clone the repository, add all the dependencies and create a vendor directory that contains all dependencies. For this, run the following commands:
```bash
go mod init terraform-provider-mongodb
go mod tidy
go mod vendor
```
## Managing terraform plugins

- Clone or Download the repositor. <br>
- Add the Public key,Private key and Organization id generated in API Key
- Run the following command :

```bash
go mod init terraform-provider-mongodb
go mod tidy
```

- Run go mod vendor to create a vendor directory that contains all the provider's dependencies.

## Managing plugins for terraform
1. Run the following command to create a vendor subdirectory which will comprise of all provider dependencies. 
```
%APPDATA%/terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${target}
```
Command:
```
mkdir -p ~/.terraform.d/plugins/mongdb.com/edu/mongodb/0.2/[OS_ARCH]
```
For eg. mkdir -p %APPDATA%/.terraform.d/plugins/mongdb.com/edu/mongodb/0.2/windows_amd64

2. Run `go build -o terraform-provider-mongodb.exe`. This will save the binary (`.exe`) file in the main/root directory. <br>
3. Run this command to move this binary file to the appropriate location.

```
move terraform-provider-mongodb.exe %APPDATA%\terraform.d\plugins\mongdb.com\edu\mongodb\0.2\[OS_ARCH]
```
[OR]
1. Manually move the file from current directory to destination directory.<br>

## Working with terraform
### Application Credential Integration in terraform
1. Add terraform block and provider block as shown in [example usage](#example-usage).  
2. Get credentials :Public Key, Private Key and Organization ID  
3. Assign the above credentials to the respective field in the provider block.` <br>

### Basic Terraform Commands

- `terraform init` -To initialize a working directory containing Terraform configuration files.
- `terraform plan` - To create an execution plan. Displays the changes to be done.
- `terraform apply` - To execute the actions proposed in a Terraform plan. Apply the changes.
- `terraform destroy` - Destroy previously-created
  infrastructure

### Invite User

1.  Add the user email and [Organization roles](https://docs.atlas.mongodb.com/reference/api/user-update/) in the respective field as shown in [example usage](#example-usage)
2. Run the basic terraform commands.
3.  You will see that a user has been successfully created and an invitation mail has been sent to the user.

### Update the user

*If you own an organization or project, you can update the user roles for any user with membership in that organization or project. You cannot modify any other user profile information.*

1. Update the data of the user resource block as show in [example usage](#example-usage) un the basic terraform commands to update user.User is not allowed to update email.
<br>

### Read the User Data

1. Add `data` and `output` blocks as shown in the [example usage](#example-usage) and run the basic terraform commands.

### Delete the user


1. Delete the `resource` block of the user and run `terraform apply`.

### Import a User command

1.  Write manually a `resource` configuration block for the user as shown in [example usage](#example-usage). Imported user will be mapped to this block
2. Run the command `terraform import mongodb_user.resource_name [EMAIL_ID]`
3.Run `terraform plan`, if output shows `0 to add, 0 to change and 0 to destroy` user import is successful, otherwise recheck the user data in `resource` block with user data in mongoDB website.

## Example Usage <a id="example-usage"></a>

```terraform
terraform {
  required_providers {
    mongodb = {
      source = "Tapendrakmr/mongodb"
      version = "0.2.4"
    }
  }
}

provider "mongodb" {
  mongodb_public_key = "[MONGODBCLOUD_PUBLIC_KEY]"
  mongodb_private_key = "[MONGODBCLOUD_PRIVATE_KEY]"
  mongodb_orgid = "[MONGODBCLOUD_ORGID]"
}

resource "mongodb_user" "user1" {
  username="user@gmail.com"
  roles=["ORG_Role1","ORG_Role2",..]
}

data "mongodb_user" "user" {
  id = "user@gmail.com"
}

output "datasource_user" {
  value = data.mongodb_user.user
}
```

## Argument Reference
- `mongodb_public_key` (required, string) - The MongoDB Public Key. This may also be set via the `"MONGODBCLOUD_PUBLIC_KEY"` environment variable.
- `mongodb_private_key` (required, string) - The MongoDB Private Key. This may also be set via the `"MONGODBCLOUD_PRIVATE_KEY"` environment variable.
- `mongodb_orgid`  (required, string) - The MongoDB OrgId. This may also be set via the `"MONGODBCLOUD_ORGID"` environment variable.
- `email` (required, string) - The email id associated with the user account.
- `roles` (required, list) - Each string in the array represents a organiation role you want to assign to the user.You must specify an array even if you are only associating a single role with the team. The following are valid roles:
 
  - ORG_OWNER
  - ORG_GROUP_CREATOR
  - ORG_BILLING_ADMIN
  - ORG_READ_ONLY
  - ORG_MEMBER



