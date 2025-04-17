# Terraform Provider Documentation

## Resource: psm_rules

The `psm_rules` resource allows you to manage Network Security Policies in PSM (Pensando Service Mesh).

### Example Usage

```hcl
resource "psm_rules" "example" {
  policy_name                 = "example-policy"
  tenant                      = "default"
  policy_distribution_target  = "default"

  rule {
    rule_name           = "allow-http"
    action              = "permit"
    description         = "Allow HTTP traffic"
    apps                = ["http"]
    from_ip_addresses   = ["10.0.0.0/24"]
    to_ip_addresses     = ["192.168.1.0/24"]
    proto_ports {
      protocol = "tcp"
      ports    = "80"
    }
  }

  rule {
    rule_name           = "deny-telnet"
    action              = "deny"
    description         = "Block Telnet traffic"
    from_workloadgroups = ["frontend"]
    to_workloadgroups   = ["backend"]
    proto_ports {
      protocol = "tcp"
      ports    = "23"
    }
  }
}
```

### Argument Reference

The following arguments are supported:

* `policy_name` - (Required) The name of the Network Security Policy.
* `tenant` - (Optional) The tenant for the policy. Defaults to "default".
* `policy_distribution_target` - (Optional) The distribution target for the policy. Defaults to "default".
* `rule` - (Optional) A list of rules for the policy. Each rule block supports the following:
    * `rule_name` - (Optional) The name of the rule.
    * `action` - (Required) The action to take. Must be either "permit" or "deny".
    * `description` - (Optional) A description of the rule.
    * `apps` - (Optional) A list of applications the rule applies to.
    * `disable` - (Optional) Whether the rule is disabled. Defaults to false.
    * `from_ip_addresses` - (Optional) A list of source IP addresses or CIDRs.
    * `to_ip_addresses` - (Optional) A list of destination IP addresses or CIDRs.
    * `from_ip_collections` - (Optional) A list of source IP collections.
    * `to_ip_collections` - (Optional) A list of destination IP collections.
    * `from_workloadgroups` - (Optional) A list of source workload groups.
    * `to_workloadgroups` - (Optional) A list of destination workload groups.
    * `rule_profile` - (Optional) The profile for the rule.
    * `proto_ports` - (Optional) A list of protocol and port combinations. Each block supports:
        * `protocol` - (Required) The protocol (e.g., "tcp", "udp").
        * `ports` - (Optional) The port or port range.
    * `labels` - (Optional) A map of key/value labels for the rule.

### Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the Network Security Policy.
* `meta` - A block containing metadata about the policy, including:
    * `name` - The name of the policy.
    * `tenant` - The tenant of the policy.
    * `namespace` - The namespace of the policy.
    * `generation_id` - The generation ID of the policy.
    * `resource_version` - The resource version of the policy.
    * `uuid` - The UUID of the policy.
    * `labels` - Any labels associated with the policy.
    * `self_link` - The self-link of the policy.
* `spec` - A block containing the specification of the policy, including:
    * `attach_tenant` - Whether the policy is attached to a tenant.
    * `rules` - The list of rules in the policy.
    * `priority` - The priority of the policy.
    * `policy_distribution_targets` - The distribution targets for the policy.

### Import

Network Security Policies can be imported using the `id`, e.g.,

```
$ terraform import psm_rules.example &lt;policy-uuid&gt;
```

Note: The import ID should be the UUID of the Network Security Policy, not its name.
