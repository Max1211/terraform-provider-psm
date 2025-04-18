# Terraform Provider Documentation

## Resource: psm_syslog_export_policy

The `psm_syslog_export_policy` resource allows you to manage Syslog Policies in PSM.

### Example Usage

```hcl
resource "psm_syslog_export_policy" "example" {
  name   = "example-syslog-policy"
  format = "syslog-rfc5424"
  filter = ["all"]

  syslogconfig {
    facility         = "user"
    disable_batching = false
  }

  psm_target {
    enable = true
  }

  targets {
    destination = "10.10.10.10"
    transport   = "udp/514"
  }

  targets {
    destination = "10.10.10.11"
    transport   = "udp/5514"
  }
}
```

```hcl
resource "psm_syslog_export_policy" "example_tls" {
  name   = "example-syslog-policy_tls"
  format = "syslog-rfc5424"
  filter = ["all"]

  syslogconfig {
    facility         = "user"
    disable_batching = false
  }

  psm_target {
    enable = true
  }

  targets {
    destination = "192.168.61.210"
    transport   = "tcp/5514"
    trusted_certs = "Server_EC_CA"
    client_certificate = "Client_RSA_Cert"
    hostname_verification = "SyslogServerName"
    skip_cert_verification = false
  }
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Syslog Policy. This must be unique within the tenant.

* `format` - (Required) The format of the syslog messages.  
Possible values: `syslog-rfc5424`, `syslog-bsd`.  
Default: `syslog-bsd`

* `filter` - (Required) A list of filter strings to apply to the syslog messages.  
Possible values: `all`, `allow`, `deny`.  
Default: `all`

* `syslogconfig` - (Required) A block to specify syslog configuration. It supports:
    * `facility` - (Required) The syslog facility to use.  
Possible values: `kernel`,` user`, `mail`, `daemon`, `auth`, `syslog`, `lpr`, `news`, `uucp`, `cron`, `authpriv`, `ftp`, `local0`, `local1`, `local2`, `local3`, `local4`, `local5`, `local6`, `local7`
    * `disable_batching` - (Required) Whether to disable batching of syslog messages.  
Possible values: `true`, `false`.

* `psm_target` - (Optional) A block to specify PSM target configuration. It supports:
    * `enable` - (Required) Whether to enable the PSM target.  
Possible values: `true`, `false`.

* `targets` - (Required) A list of target blocks. Each block supports:
    * `destination` - (Required) The destination for syslog messages (IP address notation).
    * `transport` - (Required) The transport protocol to use (e.g., "udp" or "tcp").
    * `trusted_certs` - (Optional) The server certificate.
    * `client_certificate` - (Optional) The client certificate.
    * `hostname_verification` - (Optional) The SAN name of the server certificate.
    * `skip_cert_verification` - (Optional) Verify SAN name in server certificate.
      Possible values: `true`, `false`.

-> Up to 4 syslog targets can be configured.

### Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the Syslog Policy.

### Import

Syslog Policies can be imported using the `name`, e.g.,

```
$ terraform import psm_syslog_export_policy.example example-syslog-policy
```