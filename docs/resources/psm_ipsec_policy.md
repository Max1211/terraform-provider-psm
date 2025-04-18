# Resource: psm_ipsec_policy

Manages IPSec policies in the PSM system. These policies define the configuration for IPSec tunnels.

## Example Usage single tunnel with PSK based authentication

```hcl
resource "psm_ipsec_policy" "example" {
  display_name = "Example IPSec Policy"

  tunnel {
    policy_distribution_targets = ["example-pdt"]
    ha_mode                     = "no_ha"
    disable_tcp_mss_adjust      = false

    tunnel_endpoints {
      interface_name = "tunnel100"
      dse             = "1234.1234.1234"
      ike_version    = "ikev2"

      ike_sa {
        encryption_algorithms = ["aes_gcm_128", "aes_gcm_256"]
        hash_algorithms       = ["sha_512", "sha_384", "sha_256"]
        dh_groups             = ["group20", "group19"]
        rekey_lifetime        = "8h"
        pre_shared_key        = "ExampleSecretPreSharedKeys"
        reauth_lifetime       = "24h"
        dpd_delay             = "60s"
        ikev1_dpd_timeout     = "180s"
        ike_initiator         = true
        auth_type             = "psk"
      }

      ipsec_sa {
        encryption_algorithms = ["aes_gcm_128", "aes_gcm_256"]
        dh_groups             = ["group19", "group20"]
        rekey_lifetime        = "1h"
      }

      local_identifier {
        type  = "ip"
        value = "1.2.3.4"
      }

      remote_identifier {
        type  = "ip"
        value = "5.6.7.8"
      }
    }

    lifetime {
      sa_lifetime  = "1h"
      ike_lifetime = "8h"
    }
  }
}
```

## Example Usage active_active tunnel with certificate based authentication

```hcl
resource "psm_ipsec_policy" "example" {
  display_name = "Example IPSec Policy"

  tunnel {
    ha_mode                     = "active_active"
    policy_distribution_targets = ["example-pdt"]
    disable_tcp_mss_adjust      = false

    tunnel_endpoints {
      interface_name = "tunnel100"
      dse            = "1234.1234.1234"
      ike_version    = "ikev2"

      ike_sa {
        encryption_algorithms       = ["aes_gcm_128", "aes_gcm_256"]
        hash_algorithms             = ["sha_512", "sha_384", "sha_256"]
        dh_groups                   = ["group20", "group19"]
        rekey_lifetime              = "8h"
        reauth_lifetime             = "24h"
        dpd_delay                   = "60s"
        ikev1_dpd_timeout           = "180s"
        ike_initiator               = true
        auth_type                   = "certificates"
        local_identity_certificates = file("./example_certificate.pem")
        remote_ca_certificates      = [file("./example_ca_certificate.pem")]
      }

      ipsec_sa {
        encryption_algorithms = ["aes_gcm_128", "aes_gcm_256"]
        dh_groups             = ["group19", "group20"]
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

      lifetime {
        sa_lifetime  = "1h"
        ike_lifetime = "8h"
      }
    }

    tunnel_endpoints {
      interface_name = "tunnel101"
      dse            = "4321.4321.4321"
      ike_version    = "ikev2"

      ike_sa {
        encryption_algorithms       = ["aes_gcm_128", "aes_gcm_256"]
        hash_algorithms             = ["sha_512", "sha_384", "sha_256"]
        dh_groups                   = ["group20", "group19"]
        rekey_lifetime              = "8h"
        reauth_lifetime             = "24h"
        dpd_delay                   = "60s"
        ikev1_dpd_timeout           = "180s"
        ike_initiator               = true
        auth_type                   = "certificates"
        local_identity_certificates = file("./example_certificate.pem")
        remote_ca_certificates      = [file("./example_ca_certificate.pem")]
      }

      ipsec_sa {
        encryption_algorithms = ["aes_gcm_128", "aes_gcm_256"]
        dh_groups             = ["group19", "group20"]
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

      lifetime {
        sa_lifetime  = "1h"
        ike_lifetime = "8h"
      }
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
    Possible values: `no_ha`, `active_standby`, `active_active`.
  * `disable_tcp_mss_adjust` - (Optional) Whether to disable TCP MSS adjustment. Default is false.
  * `tunnel_endpoints` - (Required) A list of tunnel endpoint configurations. Each endpoint supports:
    * `interface_name` - (Required) The name of the tunnel interface. Must correspond to the name in CXOS.
    * `dse` - (Required) The DSE (Distributed Services Engine, CX 10000 DSS-ID) for this endpoint.
    * `ike_version` - (Required) The IKE version to use. Default is "`ikev2`.  
      Possible values: `prefer_ikev2_support_ikev1`, `ikev1`, `ikev2`.
    * `ike_sa` - (Required) IKE Security Association configuration. 
    * `ipsec_sa` - (Required) IPSec Security Association configuration.
    * `local_identifier` - (Required) Local identifier configuration.
    * `remote_identifier` - (Required) Remote identifier configuration.
    * `lifetime` - (Optional) Lifetime configuration for this specific tunnel endpoint.

## Nested Blocks

### `ike_sa`

* `encryption_algorithms` - (Required) List of encryption algorithms. Default is "`aes_128`.  
  Possible values:  `aes_128`, `aes_256`, `triple_des`, `cast_128`, `aes_gcm_128`, `aes_gcm_256`.
* `hash_algorithms` - (Required) List of hash algorithms. Default is `sha_256`.  
Possible values: `sha_256`, `sha_512`, `sha_384`. 
* `dh_groups` - (Required) List of Diffie-Hellman groups. Default is `group1`.  
Possible values: `group1`, `group2`, `group5`, `group14`, `group15`, `group19`, `group20`.
* `rekey_lifetime` - (Required) Rekey lifetime. Default is "8h".  
Possible values such as: "24h", 2h45m. Valid time units are "s", "m", "h".  
Should be a valid time duration between 15m0s and 24h0m0s.
* `reauth_lifetime` - (Required) Reauthentication lifetime. Default is "24h".  
Possible values such as: "24h", 2h45m. Valid time units are "s", "m", "h".  
Should be a valid time duration between 1h0m0s and 24h0m0s.
* `dpd_delay` - (Required) Dead Peer Detection delay. Default is "60s".  
Possible values such as: "24h", 2h45m. Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h". 
Should be a valid time duration between 1m0s and 1h0m0s.
* `ikev1_dpd_timeout` - (Required) IKEv1 Dead Peer Detection timeout. Default is "180s".  
Possible values such as: "180s" or "2h45m". Valid time units are "s", "m", "h"
Should be a valid time duration between 1m0s and 1h0m0s.
* `ike_initiator` - (Required) Whether this side is the IKE initiator. Default is "false".
* `auth_type` - (Required) Authentication type. Can be "psk" or "certificates". Default is `psk`.  
Possible values: `psk`, `certificates`.
* `pre_shared_key` - (Optional) Pre-shared key for authentication (at least 20 chars). Required if `auth_type` is "psk".
* `local_identity_certificates` - (Optional) Local identity certificates. Required if `auth_type` is "certificates".
* `remote_ca_certificates` - (Optional) List of remote CA certificates. Required if `auth_type` is "certificates".

### `ipsec_sa`

* `encryption_algorithms` - (Required) List of encryption algorithms. Default is `aes_gcm_256`.  
Possible values: `aes_gcm_128`, `aes_gcm_256`.
* `dh_groups` - (Required) List of Diffie-Hellman groups. Default is `group1`.  
Possible values: `group1`, `group2`, `group5`, `group14`, `group15`, `group19`, `group20`.
* `rekey_lifetime` - (Required) Rekey lifetime. Default is "1h".  
Possible values such as: "1h", 2h45m. Valid time units are "s", "m", "h".  
Should be a valid time duration between 15m0s and 24h0m0s.

### `local_identifier` and `remote_identifier`

* `type` - (Required) The type of identifier. Can be "ip" or "fqdn". Default is "ip".  
Possible values: `keyid`, `ip`, `fqdn`, `email`.
* `value` - (Required) The value of the identifier.

### `lifetime`

* `sa_lifetime` - (Required) Security Association lifetime. Default is "8h".   
Possible values such as: "1h" or "2h45m". Valid time units are "s", "m", "h"
Should be a valid time duration between 1h0m0s and 24h0m0s.
* `ike_lifetime` - (Required) IKE lifetime. Default is "24h".  
Possible values such as: "1h" or "2h45m". Valid time units are "s", "m", "h"
Should be a valid time duration between 1h0m0s and 24h0m0s.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the IPSec policy.
* `kind` - The kind of the resource.
* `api_version` - The API version of the resource.

## Import

IPSec policies can be imported using the `id`, e.g.,

```
$ terraform import psm_ipsec_policy.example 12345678-1234-1234-1234-123456789012
```