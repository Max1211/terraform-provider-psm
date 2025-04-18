# Terraform Provider Documentation

## Resource: psm_syslog_policy

The `psm_syslog_policy` resource allows you to manage Syslog Policies in PSM (Pensando Service Mesh).

### Example Usage

```hcl
resource "psm_syslog_policy" "example" {
  name   = "example-syslog-policy"
  format = "rfc5424"
  filter = ["INFO", "WARNING"]

  syslogconfig {
    facility         = "local0"
    disable_batching = false
  }

  psm_target {
    enable = true
  }

  targets {
    destination = "192.168.1.100:514"
    transport   = "udp"
  }

  targets {
    destination = "syslog.example.com:6514"
    transport   = "tcp"
  }
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Syslog Policy. This must be unique within the tenant.

* `format` - (Required) The format of the syslog messages.

* `filter` - (Required) A list of filter strings to apply to the syslog messages.

* `syslogconfig` - (Required) A block to specify syslog configuration. It supports:
    * `facility` - (Required) The syslog facility to use.
    * `disable_batching` - (Required) Whether to disable batching of syslog messages.

* `psm_target` - (Optional) A block to specify PSM target configuration. It supports:
    * `enable` - (Required) Whether to enable the PSM target.

* `targets` - (Required) A list of target blocks. Each block supports:
    * `destination` - (Required) The destination for syslog messages.
    * `transport` - (Required) The transport protocol to use (e.g., "udp" or "tcp").

### Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the Syslog Policy.

### Import

Syslog Policies can be imported using the `name`, e.g.,

```
$ terraform import psm_syslog_policy.example example-syslog-policy
```

### Notes

* The resource supports full CRUD operations (Create, Read, Update, Delete).
* All operations are performed in the context of the "default" tenant.
* The `Read` operation populates the Terraform state with the current configuration of the Syslog Policy.

### Error Handling

* If the API returns a non-200 status code during any operation, the provider will return an error with details about the failed operation, including the HTTP status code and response body.
* Network errors or issues with JSON marshalling/unmarshalling will also result in an error being returned.

This resource allows for comprehensive management of Syslog Policies in PSM. It provides options to configure syslog message format, filtering, facility settings, and multiple targets for syslog message delivery.
