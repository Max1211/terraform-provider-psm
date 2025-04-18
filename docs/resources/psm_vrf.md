# Terraform Provider Documentation

## Resource: psm_vrf

The `psm_vrf` resource allows you to manage Virtual Routing and Forwarding (VRF) instances in PSM.

### Example Usage Host Leaf mnoe

```hcl
resource "psm_vrf" "example" {
  name                        = "example-vrf"
  egress_security_policy      = ["vrf-policy-example-egress"]
  connection_tracking_mode    = "enable"
  allow_session_reuse         = "enable"
  ip_fragments_forwarding     = "enable"
}
```

### Example Usage Border Leaf mode

```hcl
resource "psm_vrf" "example" {
  name                        = "example-vrf"
  ingress_security_policy     = ["vrf-policy-example-ingress"]
  egress_security_policy      = ["vrf-policy-example-egress"]
  connection_tracking_mode    = "enable"
  allow_session_reuse         = "enable"
  ip_fragments_forwarding     = "enable"
  ingress_nat_policy          = ["nat-policy-example-ingress"]
  egress_nat_policy           = ["nat-policy-example-egress"]
  ipsec_policy                = ["vpn-policy-example"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the VRF. This must be unique within the tenant. Changing this forces a new resource to be created.
* `ingress_security_policy` - (Optional) A list of ingress security policies to apply to the VRF.
* `egress_security_policy` - (Optional) A list of egress security policies to apply to the VRF.  

-> If several security policies are bound, each must be for a separate PDT (Policy Distribution Target)

* `connection_tracking_mode` - (Optional) The connection tracking mode for the VRF.  
Defaults to enable.
* `allow_session_reuse` - (Optional) Whether to allow session reuse for the VRF.  
Defaults to disable
* `ip_fragments_forwarding` - (Optional) Whether to allow ip fragments forwarding for the VRF.  
Defaults to disable
* `ingress_nat_policy` - (Optional) A list of ingress NAT policies to apply to the VRF.
* `egress_nat_policy` - (Optional) A list of egress NAT policies to apply to the VRF.
* `ipsec_policy` - (Optional) A list of IPsec policies to apply to the VRF.  

-> Binding NAT or IPSec policies requires the CX 10000 to be in BL (Border Leaf) mode.

* `flow_export_policy` - (Optional) A list of flow export policies to apply to the VRF.
* `maximum_cps_per_network` - (Optional) The maximum connections per second (CPS) per network for the VRF.  
Defaults to 0 (unlimited).
* `maximum_sessions_per_network` - (Optional) The maximum sessions per network for the VRF.  
Defaults to 0 (unlimited).

### Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the VRF instance.

### Import

VRF instances can be imported using the `name`, e.g.,

```
terraform import psm_vrf.example example-vrf
```
