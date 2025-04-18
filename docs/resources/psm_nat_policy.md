# PSM NAT Policy Resource

The `psm_nat_policy` resource allows you to configure Network Address Translation (NAT) policies in the PSM system.

## Resource Configuration

```hcl
resource "psm_nat_policy" "example" {
  display_name = "Example_NAT_Policy"

  rule {
    # Rule configuration
  }

  policy_distribution_targets = ["PDT01"]
}
```

## Argument Reference

* `display_name` - (Required) The display name of the NAT policy.
* `rule` - (Required) One or more rule blocks defining the NAT rules. See [Rule Configuration](#rule-configuration) below.
* `policy_distribution_targets` - (Optional) List of policy distribution targets. Default to "default"

### Rule Configuration

Each `rule` block supports the following arguments:

* `name` - (Required) Name of the rule.
* `disable` - (Optional) Whether the rule is disabled. Defaults to `false`.  
  
-> Omitting `source`, `destination`, `translated_source` and / or `translated_destination` defaults to an empty value which is translated to "any".  
For the sake of completeness using the value "any" is supported, having the same effect.

* `source` - (Optional) Source address configuration. See [Address Configuration](#address-configuration) below.
* `destination` - (Optional) Destination address configuration. See [Address Configuration](#address-configuration) below.
* `destination_proto_port` - (Optional) Protocol and port configuration. See [Protocol and Port Configuration](#protocol-and-port-configuration) below.
* `translated_source` - (Optional) Translated source address configuration. See [Address Configuration](#address-configuration) below.
* `translated_destination` - (Optional) Translated destination address configuration. See [Address Configuration](#address-configuration) below.
* `translated_destination_port` - (Optional) Translated destination port.  

-> Setting `translated_destination_port` will rewrite the destination port. Only single port rewrite is supported for single port match.

### Address Configuration

The `source`, `destination`, `translated_source`, and `translated_destination` blocks support:

* `addresses` - (Optional) List of IP addresses or CIDR ranges.
* `ipcollections` - (Optional) List of IP collection IDs.

### Protocol and Port Configuration

The `destination_proto_port` block supports:

* `protocol` - (Required) The protocol (only protocols tcp, udp, icmp, gre, esp, ah, any and protocol numbers from 0 to 254 are allowed).
* `ports` - (Optional) Port or port range (e.g., "80", "1024-2048").  

-> Single port or a single port range is supported.

## Usage Examples

### Source NAT (SNAT)

```hcl
resource "psm_nat_policy" "snat_example" {
  display_name = "SNAT_Example"

  rule {
    name = "Outbound_SNAT"
    source {
      addresses = ["192.168.1.0/24"]
    }
    translated_source {
      addresses = ["203.0.113.0/24"]
    }
  }

  policy_distribution_targets = ["PDT01"]
}
```

### Destination NAT (DNAT)

```hcl
resource "psm_nat_policy" "dnat_example" {
  display_name = "DNAT_Example"

  rule {
    name = "Inbound_DNAT"
    destination {
      addresses = ["203.0.113.100"]
    }
    destination_proto_port {
      protocol = "tcp"
      ports    = "80"
    }
    translated_destination {
      addresses = ["192.168.1.10"]
    }
    translated_destination_port = "8080"
  }

  policy_distribution_targets = ["PDT01"]
}
```

### Combined SNAT and DNAT

```hcl
resource "psm_nat_policy" "snat_dnat_example" {
  display_name = "SNAT_DNAT_Example"

  rule {
    name = "Combined_NAT"
    source {
      addresses = ["192.168.1.0/24"]
    }
    destination {
      addresses = ["10.100.2.10"]
    }
    destination_proto_port {
      protocol = "tcp"
      ports    = "80"
    }
    translated_source {
      addresses = ["10.0.0.0/24"]
    }
    translated_destination {
      addresses = ["192.168.2.10"]
    }
    translated_destination_port = "8080"
  }

  policy_distribution_targets = ["PDT01"]
}
```

### Using IP Collections

```hcl
resource "psm_nat_policy" "ipcollection_example" {
  display_name = "IP_Collection_NAT_Example"

  rule {
    name = "NAT_With_IPCollections"
    source {
      ipcollections = [psm_ipcollection.internal.id]
    }
    destination {
      ipcollections = [psm_ipcollection.external.id]
    }
    translated_source {
      addresses = ["203.0.113.0/24"]
    }
  }

  policy_distribution_targets = ["PDT01"]
}
```

### Multiple SNAT Rules in a Policy

```hcl
resource "psm_nat_policy" "multi_rule_snat_example" {
  display_name = "Multi_Rule_SNAT_Example"

  rule {
    name = "SNAT_Internal_Network"
    source {
      addresses = ["192.168.1.0/24"]
    }
    translated_source {
      addresses = ["203.0.113.0/24"]
    }
  }

  rule {
    name = "SNAT_DMZ_Servers_HTTP"
    source {
      addresses = ["172.16.1.0/24"]
    }
    destination_proto_port {
      protocol = "tcp"
      ports    = "80"
    }
    translated_source {
      addresses = ["203.0.114.0/24"]
    }
  }

  rule {
    name = "SNAT_DMZ_Servers_HTTPS"
    source {
      addresses = ["172.16.1.0/24"]
    }
    destination_proto_port {
      protocol = "tcp"
      ports    = "443"
    }
    translated_source {
      addresses = ["203.0.114.0/24"]
    }
  }

  rule {
    name = "SNAT_Development_Network"
    source {
      ipcollections = [psm_ipcollection.int1.id]
    }
    translated_source {
      addresses = ["203.0.115.0/24"]
    }
  }

  policy_distribution_targets = ["PDT01"]
}
```

### Multiple DNAT Rules in a Policy

```hcl
resource "psm_nat_policy" "multi_rule_dnat_example" {
  display_name = "Multi_Rule_DNAT_Example"

  rule {
    name = "DNAT_Web_Server_HTTP"
    destination {
      addresses = ["203.0.113.100"]
    }
    destination_proto_port {
      protocol = "tcp"
      ports    = "80"
    }
    translated_destination {
      addresses = ["192.168.1.10"]
    }
  }

  rule {
    name = "DNAT_Web_Server_HTTPS"
    destination {
      addresses = ["203.0.113.100"]
    }
    destination_proto_port {
      protocol = "tcp"
      ports    = "443"
    }
    translated_destination {
      addresses = ["192.168.1.10"]
    }
  }

  rule {
    name = "DNAT_FTP_Server_Control"
    destination {
      addresses = ["203.0.113.101"]
    }
    destination_proto_port {
      protocol = "tcp"
      ports    = "21"
    }
    translated_destination {
      addresses = ["192.168.1.20"]
    }
  }

  rule {
    name = "DNAT_FTP_Server_Passive"
    destination {
      addresses = ["203.0.113.101"]
    }
    destination_proto_port {
      protocol = "tcp"
      ports    = "50000-51000"
    }
    translated_destination {
      addresses = ["192.168.1.20"]
    }
  }

  rule {
    name = "DNAT_VPN_Gateway_IKE"
    destination {
      addresses = ["203.0.113.102"]
    }
    destination_proto_port {
      protocol = "udp"
      ports    = "500"
    }
    translated_destination {
      addresses = ["192.168.1.30"]
    }
  }

  rule {
    name = "DNAT_VPN_Gateway_IKE_NAT"
    destination {
      addresses = ["203.0.113.102"]
    }
    destination_proto_port {
      protocol = "udp"
      ports    = "4500"
    }
    translated_destination {
      addresses = ["192.168.1.30"]
    }
  }

  rule {
    name = "DNAT_VPN_Gateway_IPSEC"
    destination {
      addresses = ["203.0.113.102"]
    }
    destination_proto_port {
      protocol = "esp"
    }
    translated_destination {
      addresses = ["192.168.1.30"]
    }
  }

  policy_distribution_targets = ["PDT01"]
}
```

### Importing Existing NAT Policies

To import an existing NAT policy into Terraform, use the following command:

```text
terraform import psm_nat_policy.example <policy_id>
```
