# Resource: psm_workload

Manages workloads in the PSM system.

## Example Usage

```hcl
resource "psm_workload" "example" {
  name               = "example-workload"
  host_name          = "example-host.domain.com"

  interface {
    mac_address    = "0011.2233.4455"
    external_vlan  = 100
    ip_addresses   = ["10.0.0.10"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the workload. Changing this creates a new resource.
* `host_name` - (Required) The hostname of the workload.
* `interface` - (Required) One or more `interface` blocks as defined below.

The `interface` block supports:

* `mac_address` - (Required) The MAC address of the interface.
* `external_vlan` - (Required) The external VLAN ID for the interface.
* `ip_addresses` - (Required) A list of IP addresses assigned to the interface.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the workload (same as `name`).

## Import

Workloads can be imported using the `name`, e.g.,

```
$ terraform import psm_workload.example example-workload
```