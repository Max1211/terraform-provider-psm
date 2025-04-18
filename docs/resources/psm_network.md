# Resource: psm_network

Manages network configurations in the PSM system.

## Example Usage

```hcl
resource "psm_network" "example" {
  name                      = "example-network"
  tenant                    = "default"
  vlan_id                   = 100
  virtual_router            = "example-vrf"
  ingress_security_policy   = "ingress-policy"
  egress_security_policy    = "egress-policy"
  connection_tracking_mode  = "enable"
  allow_session_reuse       = "enable"
  service_bypass            = true
  ip_fragments_forwarding   = "enable"
  ingress_mirror_session    = "traffic-monitoring"
  egress_mirror_session     = "traffic-monitoring"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the network. This must be unique within the PSM system.

* `tenant` - (Optional) The tenant for this network. Defaults to "default".

* `vlan_id` - (Required) The VLAN ID for this network. This must be unique within the PSM system.  
  Defaults to 0.

* `virtual_router` - (Optional) The virtual router (VRF) for this network. Defaults to "default".

* `ingress_security_policy` - (Optional) The ingress security policy to apply to this network.

* `egress_security_policy` - (Optional) The egress security policy to apply to this network.

* `connection_tracking_mode` - (Optional) The connection tracking mode for this network. Defaults to `inherit from vrf`.  
  Possible values: `enable`, `disable`, `inherit from vrf`.

* `allow_session_reuse` - (Optional) Whether to allow session reuse. Defaults to `inherit from vrf`.
  Possible values: `enable`, `disable`, `inherit from vrf`.

* `service_bypass` - (Optional) Whether to enable service bypass for this network. Defaults to false.  
  Possible values: `true`, `false`.

* `ip_fragments_forwarding` - (Optional) Whether to allow ip fragments forwarding for this network.  
Possible values: `enable`, `disable`, `inherit from vrf`.

* `ingress_mirror_session` - (Optional) Mirror session to export traffic in ingress direction.  
* `egress_mirror_session` - (Optional) Mirror session to export traffic in egress direction.  

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the network (UUID).

## Import

Networks can be imported using the `name`, e.g.,

```text
terraform import psm_network.example example-network
```
