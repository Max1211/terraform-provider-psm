# Resource: psm_flow_export_policy

Manages Flow Export Policies in the PSM system. These policies define how network flow data is exported from the system.

## Example Usage

```hcl
resource "psm_flow_export_policy" "example" {
  name     = "example-policy"
  interval = "10s"
  format   = "ipfix"

  target {
    destination = "10.10.10.10"
    transport   = "udp/9995"
  }

  target {
    destination = "10.10.10.11"
    transport   = "udp/4739"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Flow Export Policy. This must be unique within the PSM system.

* `interval` - (Required) The interval at which flow data is exported.
Default: 10s.

* `format` - (Required) The format of the exported flow data.
Possible values: `ipfix`.
Default: `ipfix`

* `target` - (Required) One or more export targets. Each target is a block that supports:
  * `destination` - (Required) The destination address and port for the exported data (e.g., "192.168.1.100:4739").
  * `transport` - (Required) The transport protocol used for exporting.  
Possible values: `tcp`, `udp`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Flow Export Policy (UUID).

## Import

Flow Export Policies can be imported using the `name`, e.g.,

```
$ terraform import psm_flow_export_policy.example example-policy
```