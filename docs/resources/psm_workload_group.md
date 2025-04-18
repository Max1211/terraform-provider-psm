# Terraform Provider Documentation

## Resource: psm_workload_group

The `psm_workload_group` resource allows you to manage Workload Groups in the PSM system.

### Example Usage

```hcl
resource "psm_workload_group" "example" {
  name = "example-workload-group"

  workload_selector {
    workload_label_selector {
      workload_label_key = "namespace"
      operator           = "equals"
      values             = ["example-value"]
    }
  }
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

```text
terraform import psm_workload_group.example example-workload-group
```
