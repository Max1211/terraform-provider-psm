# Resource: psm_authpolicy

Manages the authentication policy in the PSM system.

## Example Usage

```hcl
resource "psm_authpolicy" "example" {
  token_expiry = "144h"
  authenticator_order = ["local", "ldap", "radius"]

  local {
    password_length = 12
    allowed_failed_login_attempts = 10
    failed_login_attempts_duration = "30m"
  }

  ldap {
    base_dn = "dc=example,dc=com"
    bind_dn = "cn=admin,dc=example,dc=com"
    bind_password = "password123"

    attribute_mapping {
      user = "uid"
      user_object_class = "inetOrgPerson"
      tenant = "o"
      group = "memberOf"
      group_object_class = "groupOfNames"
      email = "mail"
      fullname = "cn"
    }

    servers {
      url = "ldap.example.com:389"
      tls_options {
        start_tls = true
        server_name                   = "server.example.com"
        skip_server_cert_verification = false
        trusted_certs                 = file("./server_crt.pem")
      }
    }

    tag = "primary"
    skip_nested_groups = false
  }

  radius {
    nas_id = "psm-server"
    servers {
      url = "radius.example.com:1812"
      secret = "radiussecret"
      auth_method = "pap"
    }
    tag = "backup"
  }
}
```

## Argument Reference

The following arguments are supported:

* `token_expiry` - (Optional) The expiration time for authentication tokens. Default is "144h".  

-> When LDAP or RADIUS resources are defined, they must be defined in the `authentication_order`.

* `authenticator_order` - (Optional) The order in which authentication methods are tried. Valid values are "local", "ldap", and "radius".

* `local` - (Optional) Configuration for local authentication.
  * `password_length` - (Optional) Minimum password length. Default is 9.
  * `allowed_failed_login_attempts` - (Optional) Number of allowed failed login attempts. Default is 10.
  * `failed_login_attempts_duration` - (Optional) Duration for counting failed login attempts. Default is "15m".

* `ldap` - (Optional) Configuration for LDAP authentication. Can be specified multiple times for multiple LDAP domains.
  * `base_dn` - (Required) The base DN for LDAP searches.
  * `bind_dn` - (Required) The DN to bind with for LDAP operations.
  * `bind_password` - (Required) The password for the bind DN.
  * `attribute_mapping` - (Required) Mapping of LDAP attributes to PSM attributes.
    * `user` - (Required) LDAP attribute for user identification.
    * `user_object_class` - (Required) LDAP object class for users.
    * `tenant` - (Optional) LDAP attribute for tenant information.
    * `group` - (Optional) LDAP attribute for group membership.
    * `group_object_class` - (Optional) LDAP object class for groups.
    * `email` - (Optional) LDAP attribute for user email.
    * `fullname` - (Optional) LDAP attribute for user's full name.
  * `servers` - (Required) List of LDAP servers.
    * `url` - (Required) URL of the LDAP server.
    * `tls_options` - (Optional) TLS configuration for the LDAP connection.
      * `start_tls` - (Optional) Whether to use StartTLS.
      * `skip_server_cert_verification` - (Optional) Whether to skip server certificate verification.
      * `server_name` - (Optional) The expected server name for certificate verification.
      * `trusted_certs` - (Optional) Trusted certificates for LDAP server verification.  

-> `skip_server_cert_verification` "false" requires the use of `start_tls`.

* `tag` - (Optional) A tag for the LDAP configuration.
* `skip_nested_groups` - (Optional) Whether to skip nested group resolution.

* `radius` - (Optional) Configuration for RADIUS authentication. Can be specified multiple times for multiple RADIUS domains.
  * `nas_id` - (Required) The NAS identifier for RADIUS requests.
  * `servers` - (Required) List of RADIUS servers.
    * `url` - (Required) URL of the RADIUS server.
    * `secret` - (Optional) Shared secret for RADIUS authentication.
    * `auth_method` - (Optional) Authentication method for RADIUS.
    * `trusted_certs` - (Optional) Trusted certificates for RADIUS server verification.
  * `tag` - (Optional) A tag for the RADIUS configuration.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the authentication policy.

## Import

The authentication policy can be imported using a placeholder ID:

```text
terraform import psm_authpolicy.example authn-policy
```
