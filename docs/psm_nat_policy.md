# Resource: psm_nat_policy

Manages Network Address Translation (NAT) policies in the PSM system.

## Example Usage

```hcl
resource "psm_nat_policy" "example" {
  display_name = "Example NAT Policy"

  rule {
    name = "Rule 1"
    type = "static"
    source {
      addresses = ["10.0.0.0/24"]
    }
    destination {
      addresses = ["0.0.0.0/0"]
    }
    destination_proto_port {
      protocol = "tcp"
      ports    = "80,443"
    }
    translated_source {
      addresses = ["192.168.1.100"]
    }
  }

  rule {
    name = "Rule 2"
    type = "dynamic"
    source {
      ipcollections = ["internal_network"]
    }
    destination {
      addresses = ["0.0.0.0/0"]
    }
    destination_proto_port {
      protocol = "any"
      ports    = "any"
    }
    translated_source {
      addresses = ["203.0.113.1-203.0.113.100"]
    }
  }

  policy_distribution_targets = ["DSC1", "DSC2"]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) The display name of the NAT policy.

* `rule` - (Required) One or more rule blocks defining the NAT rules. Each rule block supports:
  * `name` - (Required) The name of the rule.
  * `disable` - (Optional) Whether the rule is disabled. Defaults to `false`.
  * `type` - (Optional) The type of NAT rule. Defaults to "static".
  * `source` - (Required) A source block defining the source of the traffic.
    * `addresses` - (Optional) A list of source IP addresses or subnets.
    * `ipcollections` - (Optional) A list of IP collections to use as sources.
  * `destination` - (Required) A destination block defining the destination of the traffic.
    * `addresses` - (Optional) A list of destination IP addresses or subnets.
    * `ipcollections` - (Optional) A list of IP collections to use as destinations.
  * `destination_proto_port` - (Required) A block defining the destination protocol and port.
    * `protocol` - (Required) The protocol (e.g., "tcp", "udp", "any").
    * `ports` - (Required) The port or port range (e.g., "80", "1024-2048", "any").
  * `translated_source` - (Optional) A block defining the translated source.
    * `addresses` - (Optional) A list of translated source IP addresses.
    * `ipcollections` - (Optional) A list of IP collections to use as translated sources.
  * `translated_destination` - (Optional) A block defining the translated destination.
    * `addresses` - (Optional) A list of translated destination IP addresses.
    * `ipcollections` - (Optional) A list of IP collections to use as translated destinations.
  * `translated_destination_port` - (Optional) The translated destination port.

* `policy_distribution_targets` - (Required) A list of distribution targets for the policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the NAT policy.

## Import

NAT policies can be imported using the policy ID, e.g.,

```
$ terraform import psm_nat_policy.example 12345
```

## Notes

1. The `source` and `destination` blocks in each rule must have either `addresses` or `ipcollections` defined. If both are empty, the rule will apply to any source/destination.

2. The `translated_source` and `translated_destination` blocks are optional. If not provided, no translation will be performed for that direction.

3. The `type` field in each rule can be "static" or "dynamic". "static" is used for one-to-one NAT, while "dynamic" is used for many-to-one NAT.

4. The `policy_distribution_targets` define where this NAT policy will be applied. These are typically the names or IDs of network devices where the policy should be distributed.

5. When updating a NAT policy, be cautious as changes may impact existing network traffic flows.

## Best Practices

1. Use meaningful names for your NAT policies and rules to easily identify their purpose.

2. Order your rules carefully. Rules are processed in the order they appear in the configuration.

3. Use IP collections where possible to make your policies more manageable and easier to update.

4. Regularly review and audit your NAT policies to ensure they align with your current network architecture and security requirements.

5. Be cautious when using "any" for protocols or ports, as this may create overly permissive rules.

6. Consider using variables for IP addresses and ports to make your Terraform configurations more flexible and reusable across different environments.

7. Always test NAT policy changes in a non-production environment before applying them to production systems.