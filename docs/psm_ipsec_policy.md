# Resource: psm_ipsec_policy

Manages IPSec policies in the PSM system. These policies define the configuration for IPSec tunnels.

## Example Usage

```hcl
resource "psm_ipsec_policy" "example" {
  display_name = "Example IPSec Policy"

  tunnel {
    policy_distribution_targets = ["DSC1", "DSC2"]
    ha_mode                     = "active_standby"
    disable_tcp_mss_adjust      = false

    tunnel_endpoints {
      interface_name = "eth0"
      dse             = "DSE1"
      ike_version     = "v2"

      ike_sa {
        encryption_algorithms = ["aes256"]
        hash_algorithms       = ["sha256"]
        dh_groups             = ["modp2048"]
        rekey_lifetime        = "8h"
        pre_shared_key        = "your-preshared-key"
        auth_type             = "psk"
      }

      ipsec_sa {
        encryption_algorithms = ["aes256"]
        dh_groups             = ["modp2048"]
        rekey_lifetime        = "1h"
      }

      local_identifier {
        type  = "fqdn"
        value = "local.example.com"
      }

      remote_identifier {
        type  = "fqdn"
        value = "remote.example.com"
      }
    }

    lifetime {
      sa_lifetime  = "8h"
      ike_lifetime = "24h"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) The display name of the IPSec policy.

* `tunnel` - (Required) A block that defines the tunnel configuration. It supports the following:
  * `policy_distribution_targets` - (Required) List of distribution targets for the policy.
  * `ha_mode` - (Optional) High availability mode. Default is "no_ha".
  * `disable_tcp_mss_adjust` - (Optional) Whether to disable TCP MSS adjustment. Default is false.
  * `tunnel_endpoints` - (Required) A list of tunnel endpoint configurations. Each endpoint supports:
    * `interface_name` - (Required) The name of the interface.
    * `dse` - (Required) The DSE (Distributed Services Engine) for this endpoint.
    * `ike_version` - (Required) The IKE version to use.
    * `ike_sa` - (Required) IKE Security Association configuration.
    * `ipsec_sa` - (Required) IPSec Security Association configuration.
    * `local_identifier` - (Required) Local identifier configuration.
    * `remote_identifier` - (Required) Remote identifier configuration.
  * `lifetime` - (Optional) Lifetime configuration for the tunnel.

## Nested Blocks

### `ike_sa`

* `encryption_algorithms` - (Required) List of encryption algorithms.
* `hash_algorithms` - (Required) List of hash algorithms.
* `dh_groups` - (Required) List of Diffie-Hellman groups.
* `rekey_lifetime` - (Optional) Rekey lifetime. Default is "8h".
* `pre_shared_key` - (Optional) Pre-shared key for authentication.
* `reauth_lifetime` - (Optional) Reauthentication lifetime. Default is "24h".
* `dpd_delay` - (Optional) Dead Peer Detection delay. Default is "60s".
* `ikev1_dpd_timeout` - (Optional) IKEv1 Dead Peer Detection timeout. Default is "180s".
* `ike_initiator` - (Optional) Whether this side is the IKE initiator. Default is true.
* `auth_type` - (Optional) Authentication type. Default is "psk".
* `local_identity_certificates` - (Optional) Local identity certificates.
* `remote_ca_certificates` - (Optional) List of remote CA certificates.

### `ipsec_sa`

* `encryption_algorithms` - (Required) List of encryption algorithms.
* `dh_groups` - (Required) List of Diffie-Hellman groups.
* `rekey_lifetime` - (Optional) Rekey lifetime. Default is "1h".

### `local_identifier` and `remote_identifier`

* `type` - (Required) The type of identifier.
* `value` - (Required) The value of the identifier.

### `lifetime`

* `sa_lifetime` - (Optional) Security Association lifetime.
* `ike_lifetime` - (Optional) IKE lifetime.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the IPSec policy.
* `kind` - The kind of the resource.
* `api_version` - The API version of the resource.

## Import

IPSec policies can be imported using the policy ID, e.g.,

```
$ terraform import psm_ipsec_policy.example 12345
```

## Notes

1. The `pre_shared_key` in the `ike_sa` block is sensitive and will not be displayed in logs or console output.

2. When updating an IPSec policy, be cautious as changes may impact existing VPN connections.

3. The `policy_distribution_targets` define where this IPSec policy will be applied. These are typically the names or IDs of network devices where the policy should be distributed.

4. The `ha_mode` can be set to "no_ha", "active_standby", or other values supported by your PSM system.

## Best Practices

1. Use meaningful names for your IPSec policies to easily identify their purpose.

2. Regularly review and update your IPSec policies to ensure they align with your current security requirements.

3. Use variables for sensitive information like pre-shared keys to keep them out of your main configuration.

4. Consider using stronger encryption algorithms and longer key lengths for better security.

5. Implement proper access controls and review processes before applying changes to IPSec policies, as these changes can have significant impacts on network connectivity.

6. Regularly rotate pre-shared keys and certificates used in IPSec configurations.

7. Monitor the lifetimes set for SA and IKE to ensure they meet your security and performance requirements.
