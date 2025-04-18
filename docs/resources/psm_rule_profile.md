# Terraform Provider Documentation

## Resource: psm_rule_profile

Manages rule profiles in the PSM system.

### Example Usage

```hcl
resource "psm_rule_profile" "example" {
  name                 = "example-rule-profile"
  conn_track           = "enable"
  allow_session_reuse  = "disable"
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Rule Profile. This must be unique within the tenant.

* `conn_track` - (Optional) The connection tracking mode. Valid values are:
  * `"inherit"` (default)
  * `"enable"`
  * `"disable"`

* `allow_session_reuse` - (Optional) Whether to allow session reuse. Valid values are:
  * `"inherit"` (default)
  * `"enable"`
  * `"disable"`

### Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Rule Profile. This is the same as the `name`.

### Import

Rule Profiles can be imported using the `name`, e.g.,

```text
terraform import psm_rule_profile.example example-rule-profile
```
