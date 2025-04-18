# Resource: psm_ui_global_settings

Manages the global UI settings in the PSM system.

## Example Usage

```hcl
resource "psm_ui_global_settings" "example" {
  duration               = "30m"
  warning_time           = "30s"
  enable_object_renaming = false
}
```

## Argument Reference

The following arguments are supported:

* `duration` - (Optional) The duration of the idle timeout. Defaults to "60m".  
  Possible values: "30m", "60m".

* `warning_time` - (Optional) The warning time before the idle timeout occurs. Defaults to "10s".  
    Possible values: "10s", "20s", "40s", "60s", "2m".

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

### Delete

As mentioned in the notes, the delete operation resets the configuration to default values:
- `duration`: "60m"
- `warning_time`: "10s"
- `enable_object_renaming`: true