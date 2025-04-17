# Terraform Provider Documentation

## Resource: psm_vrf

The `psm_vrf` resource allows you to manage Virtual Routing and Forwarding (VRF) instances in PSM (Pensando Service Mesh).

### Example Usage

```hcl
resource "psm_vrf" "example" {
  name                     = "example-vrf"
  ingress_security_policy  = "ingress-policy-name"
  egress_security_policy   = "egress-policy-name"
  connection_tracking_mode = "per-session"
  allow_session_reuse      = "disable"
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the VRF instance. This must be unique within the tenant.
* `ingress_security_policy` - (Optional) The name of the ingress security policy to apply to this VRF.
* `egress_security_policy` - (Optional) The name of the egress security policy to apply to this VRF.
* `connection_tracking_mode` - (Optional) The connection tracking mode for this VRF.
* `allow_session_reuse` - (Optional) Whether to allow session reuse for this VRF.

### Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the VRF instance.

### Import

VRF instances can be imported using the `name`, e.g.,

```
$ terraform import psm_vrf.example example-vrf
```

### Notes

* The `default` VRF is a special case. When creating or updating a VRF named "default", the provider will not perform any API operations and will simply set the `id` to "default".
* The VRF type is always set to "unknown" when creating or updating a VRF.
* The tenant for the VRF is always set to "default".

### Limitations

* The current implementation does not support setting or reading all available VRF properties. Some properties like `router_mac_address`, `vxlan_vni`, `default_ipam_policy`, and others are not exposed through this resource.
* The `Read` operation does not currently populate all the VRF properties back into the Terraform state. Only `name`, `kind`, and `api_version` are set after a read operation.

### Error Handling

* If the API returns a non-200 status code during create, update, or delete operations, the provider will return an error with details about the failed operation.
* Any network errors or issues with marshalling/unmarshalling JSON will also result in an error being returned.

This resource allows for basic management of VRF instances. For more advanced configurations or to access additional VRF properties, you may need to extend this resource or use other means to interact with the PSM API.
