# Terraform Provider Documentation

## Resource: psm_user

The `psm_user` resource allows you to manage Users in PSM System.

### Example Usage

```hcl
resource "psm_user" "example" {
  name     = "john.doe"
  fullname = "John Doe"
  email    = "john.doe@example.com"
  password = "Pensando0$"
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

```text
terraform import psm_user.example my-tenant/default/john.doe
```
