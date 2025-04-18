# Terraform Provider Documentation

## Resource: psm_user_preferences

The `psm_user_preferences` resource allows you to manage User Preferences in PSM (Pensando Service Mesh).

### Example Usage

```hcl
resource "psm_user_preferences" "example" {
  name            = "admin"
  timezone_utc    = true
  service_cards   = ["NetworkCard", "SecurityCard"]
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the user preferences. This is typically "admin".

* `timezone_utc` - (Optional) Set to `true` to use UTC timezone. Conflicts with `timezone_name` and `timezone_client`.

* `timezone_name` - (Optional) The name of the timezone to use. Conflicts with `timezone_utc` and `timezone_client`.

* `timezone_client` - (Optional) Set to `true` to use the client's timezone. Conflicts with `timezone_utc` and `timezone_name`.

* `service_cards` - (Optional) A set of service cards to be active on the dashboard.

### Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the user preferences. This is always set to "admin".

### Import

User Preferences can be imported using the `name`, e.g.,

```
$ terraform import psm_user_preferences.example admin
```

### Notes

* The `Kind` and `APIVersion` fields are automatically set to "UserPreference" and "v1" respectively.
* The tenant and namespace are always set to "default".
* This resource doesn't support true deletion. When deleted, it resets preferences to default values.
* Only one of `timezone_utc`, `timezone_name`, or `timezone_client` can be set at a time.

### Limitations

* This resource is designed to work with a single user (typically "admin"). It may not support multiple user preferences.
* The API interactions are always performed for the "admin" user in the "default" tenant and namespace.

### Error Handling

* If the API returns a non-200 status code during any operation, the provider will return an error with details about the failed operation.
* Network errors or issues with JSON marshalling/unmarshalling will also result in an error being returned.

This resource allows for management of user preferences in PSM, including timezone settings and active service cards on the dashboard. It's particularly useful for customizing the PSM user interface experience.
