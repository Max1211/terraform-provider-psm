# Terraform Provider Documentation

## Resource: psm_user

The `psm_user` resource allows you to manage Users in PSM (Pensando Service Mesh).

### Example Usage

```hcl
resource "psm_user" "example" {
  name     = "john.doe"
  fullname = "John Doe"
  email    = "john.doe@example.com"
  password = "securepassword123"
  type     = "local"
  tenant   = "my-tenant"
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The username of the user. This must be unique within the tenant.

* `fullname` - (Required) The full name of the user.

* `email` - (Required) The email address of the user.

* `password` - (Required) The password for the user. This field is sensitive and will not be displayed in logs or regular output.

* `type` - (Optional) The type of user. Defaults to "local".

* `tenant` - (Optional) The tenant for the user. Defaults to "default".

* `namespace` - (Optional) The namespace for the user. Defaults to "default".

### Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the user.

* `authenticators` - The list of authenticators for the user.

* `failed_login_attempts` - The number of failed login attempts for the user.

* `locked` - Whether the user account is locked.

### Import

Users can be imported using the `tenant/namespace/name` format, e.g.,

```
$ terraform import psm_user.example my-tenant/default/john.doe
```

### Notes

* The `Kind` and `APIVersion` fields are automatically set to "User" and "v1" respectively.
* The resource supports full CRUD operations (Create, Read, Update, Delete).
* When updating a user, if the password is not changed, it will not be sent in the update request.
* The `authenticators`, `failed_login_attempts`, and `locked` attributes are read-only and computed.

### Error Handling

* If the API returns a non-200 status code during any operation, the provider will return an error with details about the failed operation, including the HTTP status code and response body.
* During creation, if a user with the same name already exists, a specific error message will be returned suggesting to use a different username or import the existing user.
* During deletion, if the user still exists after a successful delete operation, an error will be returned.
* Network errors or issues with JSON marshalling/unmarshalling will also result in an error being returned.

This resource allows for comprehensive management of Users in PSM. It provides options to create, read, update, and delete users, as well as import existing users into Terraform state.
