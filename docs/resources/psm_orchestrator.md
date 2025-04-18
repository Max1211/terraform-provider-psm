# Resource: psm_orchestrator

Manages Orchestrator integrations in the PSM system.

## Example Usage

```hcl
resource "psm_orchestrator" "example" {
  name     = "example-orchestrator"
  type     = "vcenter"
  uri      = "vcenter01.example.com"
  username = "psm@example.com"
  password = "Pensando0$"

  ca_data = file("path/to/ca.crt")
}
```

## Example Usage with explicit namespaces (datacenter in vcenter)

```hcl
resource "psm_orchestrator" "example" {
  name     = "example-orchestrator"
  type     = "vcenter"
  uri      = "vcenter01.example.com"
  username = "psm@example.com"
  password = "Pensando0$"

  ca_data = file("path/to/ca.crt")

  namespaces {
    name = "dc01"
    mode = "smartservicemonitored"
  }

  namespaces {
    name = "dc02"
    mode = "smartserviceunmonitored"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Orchestrator integration.

* `type` - (Required) The type of the Orchestrator.  
Possible values: `vcenter`.  
Default: `vcenter`.

* `uri` - (Required) The URI of the Orchestrator.

* `username` - (Required) The username for authentication.

* `password` - (Required) The password for authentication.

* `ca_data` - (Optional) The CA certificate data for server authentication.

* `disable_server_authentication` - (Optional) Whether to disable server authentication. 
Defaults to `true` if `ca_data` is not provided.

* `namespaces` - (Optional) A list of namespace configurations. If not provided, a default namespace "all_namespaces" with mode "smartservicemonitored" will be used. Each namespace block supports:
  * `name` - (Required) The name of the namespace.
  * `mode` - (Required) The mode of the namespace. Possible values: `smartservicemonitored`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Orchestrator integration (UUID).

## Import

Orchestrator integrations can be imported using the `name`, e.g.,

```
$ terraform import psm_orchestrator.example example-orchestrator
```