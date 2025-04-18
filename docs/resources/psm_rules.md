# Terraform Provider Documentation

## Resource: psm_rules

The `psm_rules` resource allows you to manage Network Security Policies in the PSM system.

### Example Usage IPv4

```hcl
resource "psm_rules" "example" {
  policy_name           = "example-policy"

  rule {
    rule_name           = "allow-http"
    action              = "permit"
    description         = "Allow HTTP traffic"
    from_ip_addresses   = ["10.0.0.0/24", "10.0.10.0/24"]
    to_ip_addresses     = ["192.168.1.0/24"]
    proto_ports {
      protocol = "tcp"
      ports    = "80"
    }
    labels = {
      "Application" : "HTTP"
      "Ticket-ID" : "INC0017399"
    }
  }

  rule {
    rule_name             = "allow-ssh"
    action                = "permit"
    description           = "Allow SSH traffic"
    apps                  = ["SSH"]
    from_ip_collections   = ["Network_10.0.0.0_24"]
    to_ip_collections     = ["Network_192.168.1.0/24"]
    labels = {
      "Application" : "SSH"
      "Ticket-ID" : "INC0017400"
    }
  }

  rule {
    rule_name             = "allow-ntp"
    action                = "permit"
    description           = "Allow NTP traffic"
    apps                  = ["NTP"]
    from_ip_collections   = ["Network_10.0.0.0_24"]
    to_workloadgroups     = ["WLG-NTP-Server"]
    labels = {
      "Application" : "NTP"
      "Ticket-ID" : "INC0017401"
    }
  }

  rule {
    rule_name           = "deny-any"
    action              = "deny"
    description         = "Deny any traffic"
    from_ip_addresses   = ["any"]
    to_ip_addresses     = ["any"]
    proto_ports {
      protocol = "any"
    }
    labels = {
      "Application" : "ANY"
      "Ticket-ID" : "INC0017402"
    }
  }
}
```

### Example Usage IPv6

```hcl
resource "psm_rules" "example_ipv6" {
  policy_name           = "example-policy_ipv6"
  address_family        = "IPv6"

  rule {
    rule_name           = "allow-http"
    action              = "permit"
    description         = "Allow HTTP traffic"
    from_ip_addresses   = ["::/0"]
    to_ip_addresses     = ["::/0"]
    proto_ports {
      protocol = "tcp"
      ports    = "80"
    }
    labels = {
      "Application" : "HTTP"
      "Ticket-ID" : "INC0017399"
    }
  }

  rule {
    rule_name             = "allow-ssh"
    action                = "permit"
    description           = "Allow SSH traffic"
    apps                  = ["SSH"]
    from_ip_collections = ["Network_2001:db8:1::/48"]
    to_ip_collections   = ["Network_2001:db8:2::/48"]
    labels = {
      "Application" : "SSH"
      "Ticket-ID" : "INC0017400"
    }
  }

  rule {
    rule_name             = "allow-ntp"
    action                = "permit"
    description           = "Allow NTP traffic"
    apps                  = ["NTP"]
    from_ip_collections   = ["Network_2001:db8:3::/48"]
    to_workloadgroups     = ["WLG-NTP-Server"]
    labels = {
      "Application" : "NTP"
      "Ticket-ID" : "INC0017401"
    }
  }

  rule {
    rule_name           = "deny-any"
    action              = "deny"
    description         = "Deny any traffic"
    from_ip_addresses   = ["any"]
    to_ip_addresses     = ["any"]
    proto_ports {
      protocol = "any"
    }
    labels = {
      "Application" : "ANY"
      "Ticket-ID" : "INC0017402"
    }
  }
}
```

### Argument Reference

The following arguments are supported:

* `policy_name` - (Required) The name of the Network Security Policy.
* `tenant` - (Optional) The tenant for the policy. Defaults to "default".
* `policy_distribution_target` - (Optional) The distribution target for the policy. Defaults to "default".
* `tenant` - (Optional) The tenant for the policy. Defaults to "default".
* `address_family` - (Optional) The address family of the security policy. Defaults to "IPv4".
  Possible values: `IPv4`, `IPv6`.
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

### Import

Network Security Policies can be imported using the `id`, e.g.,

```
terraform import psm_rules.example &lt;policy-uuid&gt;
```

Note: The import ID should be the UUID of the Network Security Policy, not its name.
