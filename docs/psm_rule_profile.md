# Terraform Provider Documentation

## Resource: psm_rule_profile

The `psm_rule_profile` resource allows you to manage Rule Profiles in PSM (Pensando Service Mesh).

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
  - `"inherit"` (default)
  - `"enable"`
  - `"disable"`

* `allow_session_reuse` - (Optional) Whether to allow session reuse. Valid values are:
  - `"inherit"` (default)
  - `"enable"`
  - `"disable"`

### Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Rule Profile. This is the same as the `name`.

### Import

Rule Profiles can be imported using the `name`, e.g.,

```
$ terraform import psm_rule_profile.example example-rule-profile
```

### Notes

* The `Kind` and `APIVersion` fields are automatically set to "RuleProfile" and "v1" respectively.
* The resource supports full CRUD operations (Create, Read, Update, Delete).
* All operations are performed in the context of the "default" tenant.

### Error Handling

* If the API returns a non-200 status code during any operation, the provider will return an error with details about the failed operation, including the HTTP status code and, where available, the response body.
* Network errors or issues with JSON marshalling/unmarshalling will also result in an error being returned.

This resource allows for comprehensive management of Rule Profiles in PSM. It provides options to configure connection tracking and session reuse behavior, which can be applied to security rules for fine-grained control over network traffic.
