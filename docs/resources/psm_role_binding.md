# Terraform Provider Documentation

## Resource: psm_role_binding

The `psm_role_binding` resource allows you to manage Role Bindings in PSM (Pensando Service Mesh).

### Example Usage

```hcl
resource "psm_role_binding" "example" {
  name       = "example-role-binding"
  tenant     = "my-tenant"
  namespace  = "my-namespace"
  role       = "reader"
  users      = ["user1@example.com", "user2@example.com"]
  user_groups = ["group1", "group2"]
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

```
$ terraform import psm_role_binding.example my-tenant/example-role-binding
```

### Notes

* The `Kind` and `APIVersion` fields are automatically set to "RoleBindingList" and "v1" respectively.
* The resource supports full CRUD operations (Create, Read, Update, Delete).
* When updating a Role Binding, all fields (except `name` and `tenant`) can be modified.

### Error Handling

* If the API returns a non-200 status code during any operation, the provider will return an error with details about the failed operation, including the HTTP status code and response body.
* Network errors or issues with JSON marshalling/unmarshalling will also result in an error being returned.

This resource allows for comprehensive management of Role Bindings in PSM. It provides flexibility in assigning roles to both individual users and user groups within specific tenants and namespaces.
