# Terraform Provider Documentation

## Resource: psm_role_binding

The `psm_role_binding` resource allows you to manage Role Bindings in the PSM system.

### Example Usage

```hcl
resource "psm_role_binding" "example" {
  name       = "example-role-binding"
  role       = "AdminRole"
  users      = ["user1", "user2"]
}
```

  ->
  The `user_groups` attribute is the link between the LDAP group distinguished name (DN) or radius group name.
  
### Example Usage with LDAP authentication

```hcl
resource "psm_role_binding" "admin_ldap" {
  name        = "example-ldap-role-binding"
  role        = "AdminRole"
  user_groups = ["CN=example-user,CN=Builtin,DC=example,DC=com"]
}
```

### Example Usage with RADIUS authentication

```hcl
resource "psm_role_binding" "admin_radius" {
  name        = "example-radius-role-binding"
  role        = "AdminRole"
  user_groups = ["PSM-ReadOnly"]
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Role Binding. This must be unique within the tenant.

* `tenant` - (Optional) The tenant for the Role Binding. Defaults to "default".

* `namespace` - (Optional) The namespace for the Role Binding. Defaults to "default".

* `role` - (Required) The role to be bound to the users or user groups.

* `users` - (Optional) A set of user identifiers to be bound to the role.

* `user_groups` - (Optional) A set of user group identifiers to be bound to the role.

### Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the Role Binding.

### Import

Role Bindings can be imported using the `tenant/name` format, e.g.,

```text
terraform import psm_role_binding.example my-tenant/example-role-binding
```
