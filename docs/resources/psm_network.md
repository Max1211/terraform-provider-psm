# Resource: psm_network

Manages network configurations in the PSM system.

## Example Usage

```hcl
resource "psm_network" "example" {
  name                      = "example-network"
  tenant                    = "default"
  vlan_id                   = 100
  virtual_router            = "custom-router"
  ingress_security_policy   = "ingress-policy"
  egress_security_policy    = "egress-policy"
  connection_tracking_mode  = "per-connection"
  allow_session_reuse       = "enabled"
  service_bypass            = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the network. This must be unique within the PSM system.

* `tenant` - (Optional) The tenant for this network. Defaults to "default".

* `vlan_id` - (Optional) The VLAN ID for this network. Defaults to 0.

* `virtual_router` - (Optional) The virtual router for this network. Defaults to "default".

* `ingress_security_policy` - (Optional) The ingress security policy to apply to this network.

* `egress_security_policy` - (Optional) The egress security policy to apply to this network.

* `connection_tracking_mode` - (Optional) The connection tracking mode for this network. Can be set to "inherit from vrf" or other values supported by your PSM system.

* `allow_session_reuse` - (Optional) Whether to allow session reuse. Can be set to "inherit from vrf", "enabled", or "disabled".

* `service_bypass` - (Optional) Whether to enable service bypass for this network. Defaults to false.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the network (UUID).

## Import

Networks can be imported using the `name`, e.g.,

```
$ terraform import psm_network.example example-network
```

## Lifecycle Management

### Creation

When creating a new network, you must provide a unique `name`. Other attributes will use their default values if not specified.

### Read

The resource can be read at any time to sync the Terraform state with the actual configuration in the PSM system.

### Update

Networks can be updated after creation. Most fields can be modified, but some (like `name`, `tenant`, and `vlan_id`) are ForceNew and will result in the creation of a new resource if changed.

### Deletion

When a network is deleted through Terraform, it will be removed from the PSM system.

## Notes

1. The `name` field is used as the identifier for the network and cannot be changed after creation. If you need to rename a network, you must create a new resource and delete the old one.

2. The `tenant` field is ForceNew, meaning changing it will result in the creation of a new resource.

3. The `vlan_id` is also ForceNew. Changing it will create a new network resource.

4. When updating a network, be cautious as changes may impact existing network configurations and connected systems.

5. The `connection_tracking_mode` and `allow_session_reuse` fields default to "inherit from vrf" if not specified.

6. The `service_bypass` field is a boolean. Set it to `true` to enable service bypass for the network.

## Best Practices

1. Use meaningful names for your networks to easily identify their purpose and scope.

2. Consider using variables for VLAN IDs, virtual router names, and policy names to make your Terraform configurations more flexible and reusable across different environments.

3. Regularly review and update your network configurations to ensure they align with your current security and connectivity requirements.

4. Be cautious when modifying existing networks, especially in production environments. Consider the impact on connected systems and applications.

5. Use the `connection_tracking_mode` and `allow_session_reuse` settings judiciously, as they can impact network performance and security.

6. Document the purpose and configuration of each network resource in your Terraform code comments or external documentation.

7. When using `service_bypass`, ensure you understand the implications for your network security and monitoring capabilities.
