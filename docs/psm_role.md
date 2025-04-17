# Terraform Provider Documentation

## Resource: psm_role

The `psm_role` resource allows you to manage Roles in PSM (Pensando Service Mesh).

### Example Usage

```hcl
resource "psm_role" "example" {
  name      = "example-role"
  tenant    = "my-tenant"
  namespace = "my-namespace"

  permissions {
    resource_group = "Network"
    resource_kind  = "NetworkInterface"
    actions        = ["Get", "List", "Create", "Update", "Delete"]
  }

  permissions {
    resource_group = "Auth"
    resource_kind  = "User"
    actions        = ["Get", "List"]
  }
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Role. This must be unique within the tenant and namespace.

* `tenant` - (Optional) The tenant for the Role. Defaults to "default".

* `namespace` - (Optional) The namespace for the Role. Defaults to "default".

* `permissions` - (Required) A list of permission blocks. Each permission block supports:
    * `resource_group` - (Required) The resource group for the permission.
    * `resource_kind` - (Required) The resource kind for the permission.
    * `actions` - (Required) A list of actions allowed for this resource.

### Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the Role.

### Import

Roles can be imported using the `tenant/namespace/name` format, e.g.,

```
$ terraform import psm_role.example my-tenant/my-namespace/example-role
```

### Notes

* The `Kind` and `APIVersion` fields are automatically set to "Role" and "v1" respectively.
* The resource supports full CRUD operations (Create, Read, Update, Delete).
* When creating or updating a Role, all permissions are replaced with the ones specified in the Terraform configuration.
* The `ResourceNamespace` for permissions is automatically set to "*_ALL_*".

### Error Handling

* If the API returns a non-200 status code during any operation, the provider will return an error with details about the failed operation, including the HTTP status code and response body.
* Network errors or issues with JSON marshalling/unmarshalling will also result in an error being returned.
* During deletion, if the Role still exists after a successful delete operation, an error will be returned.

This resource allows for comprehensive management of Roles in PSM. It provides flexibility in defining permissions for different resource groups and kinds, allowing for fine-grained access control within your PSM environment.
