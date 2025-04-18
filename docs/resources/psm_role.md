# Terraform Provider Documentation

## Resource: psm_role

The `psm_role` resource allows you to manage Roles in PSM System.

### Example Usage

```hcl
resource "psm_role" "example" {
  name      = "example-role"

  permissions {
    resource_group = "Objstore"
    resource_kind  = "_All_"
    actions        = ["all-actions"]
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
    Possible values: `_All_`, `auth`, `cluster`, `monitoring`, `network`, `objstore`, `orchestration`, `others`, `preferences`, `rollout`, `security`, `staging`, `workload`.
    * `resource_kind` - (Required) The resource kind for the permission, value depending on `resource_group`.  
    Possible values:  
    `_All_`, 
    * `actions` - (Required) A list of actions allowed for this resource.  
    Possible values: `all-actions`, `create`, `read`, `update`, `delete`, `commit`, `clear`.  
    Default value: `all-actions`.

### Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the Role.

### Import

Roles can be imported using the `tenant/namespace/name` format, e.g.,

```
$ terraform import psm_role.example my-tenant/my-namespace/example-role
```