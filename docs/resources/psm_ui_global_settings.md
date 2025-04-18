# Resource: psm_ui_global_settings

Manages the global UI settings in the PSM system.

## Example Usage

```hcl
resource "psm_ui_global_settings" "example" {
  duration               = "90m"
  warning_time           = "15s"
  enable_object_renaming = false
}
```

## Argument Reference

The following arguments are supported:

* `duration` - (Optional) The duration of the idle timeout. Defaults to "60m".

* `warning_time` - (Optional) The warning time before the idle timeout occurs. Defaults to "10s".

* `enable_object_renaming` - (Optional) Whether to enable object renaming in the UI. Defaults to true.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `name` - The name of the UI global settings. This is always set to "default-ui-global-settings".

## Import

The UI global settings can be imported using a placeholder ID, e.g.,

```
$ terraform import psm_ui_global_settings.example default-ui-global-settings
```

## Notes

1. This resource manages a singleton configuration. Only one instance of this resource should be defined in your Terraform configuration.

2. The `name` attribute is computed and will always be "default-ui-global-settings".

3. When specifying durations, use Go duration string format (e.g., "60m" for 60 minutes, "10s" for 10 seconds).

4. The `delete` operation for this resource does not actually delete the configuration but resets it to default values.

## Lifecycle Management

### Creation and Update

When creating or updating the UI global settings, the resource will set the specified values or use defaults if not provided.

### Read

The resource can be read at any time to sync the Terraform state with the actual configuration in the PSM system.

### Delete

As mentioned in the notes, the delete operation resets the configuration to default values:
- `duration`: "60m"
- `warning_time`: "10s"
- `enable_object_renaming`: true

## Best Practices

1. Use variables for the duration and warning time to make your Terraform configurations more flexible and reusable across different environments.

2. Consider the implications of changing the idle timeout settings on user experience and security.

3. Be cautious when disabling object renaming, as it may impact user workflows in the PSM UI.

4. Always use meaningful names for your resource instances, even though this resource manages a singleton configuration.

5. Regularly review and update your UI global settings to ensure they align with your organization's policies and user needs.

6. When importing existing settings, verify that the imported values match your expectations and adjust your Terraform configuration accordingly.

7. Remember that changes to these settings will affect all users of the PSM UI, so plan and communicate changes appropriately.
