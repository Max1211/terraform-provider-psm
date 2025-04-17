# Terraform Provider Documentation

## Resource: psm_workload_group

The `psm_workload_group` resource allows you to manage Workload Groups in PSM (Pensando Service Mesh).

### Example Usage

```hcl
resource "psm_workload_group" "example" {
  name = "example-workload-group"

  workload_selector {
    workload_label_selector {
      workload_label_key = "app"
      operator           = "In"
      values             = ["web", "api"]
    }
    workload_label_selector {
      workload_label_key = "environment"
      operator           = "Equals"
      values             = ["production"]
    }
  }

  ip_collections = ["internal-network", "dmz"]
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Workload Group. This must be unique within the tenant.

* `workload_selector` - (Optional) A list of workload selectors. Each workload selector block supports the following:
    * `workload_label_selector` - (Optional) A list of label selectors. Each label selector block supports:
        * `workload_label_key` - (Required) The key of the workload label to match.
        * `operator` - (Required) The operator to use for matching. Valid operators include "In", "NotIn", "Exists", "DoesNotExist", "Equals", and "NotEquals".
        * `values` - (Required) A list of values to match against.

* `ip_collections` - (Optional) A list of IP collection names associated with this Workload Group.

### Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the Workload Group.

### Import

Workload Groups can be imported using the `name`, e.g.,

```
$ terraform import psm_workload_group.example example-workload-group
```

### Notes

* The tenant for the Workload Group is always set to "default" in the current implementation.
* The `Read` operation populates the Terraform state with the `name`, `workload_label_selector`, and `ip_collections` attributes.

### Limitations

* The current implementation does not support setting or reading all available Workload Group properties. Some properties might not be exposed through this resource.
* Error handling is implemented for non-200 HTTP status codes, but specific error messages from the API are not parsed or exposed in detail.

### Error Handling

* If the API returns a non-200 status code during create, update, or delete operations, the provider will return an error with details about the failed operation.
* Any network errors or issues with marshalling/unmarshalling JSON will also result in an error being returned.

This resource allows for basic management of Workload Groups. For more advanced configurations or to access additional Workload Group properties, you may need to extend this resource or use other means to interact with the PSM API.
