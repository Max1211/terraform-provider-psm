# Terraform Provider Documentation

## Resource: psm_user_preferences

The `psm_user_preferences` resource allows you to manage User Preferences in the PSM system.

### Example Usage

```hcl
resource "psm_user_preferences" "example" {
  name            = "admin"
  timezone_client = true
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the user preferences. This is typically "admin".

* `timezone_utc` - (Optional) Set to `true` to use UTC timezone. Conflicts with `timezone_name` and `timezone_client`.

* `timezone_name` - (Optional) The name of the timezone to use. Conflicts with `timezone_utc` and `timezone_client`.

* `timezone_client` - (Optional) Set to `true` to use the client's timezone. Conflicts with `timezone_utc` and `timezone_name`.

### Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the user preferences. This is always set to "admin".

### Import

User Preferences can be imported using the `name`, e.g.,

```text
terraform import psm_user_preferences.example admin
```
