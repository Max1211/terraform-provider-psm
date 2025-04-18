# Resource: psm_flow_export_policy

Manages Flow Export Policies in the PSM system. These policies define how network flow data is exported from the system.

## Example Usage

```hcl
resource "psm_flow_export_policy" "example" {
  name     = "example-policy"
  interval = "60s"
  format   = "ipfix"

  target {
    destination = "192.168.1.100:4739"
    transport   = "udp"
  }

  target {
    destination = "192.168.1.101:4739"
    transport   = "tcp"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Flow Export Policy. This must be unique within the PSM system.

* `interval` - (Required) The interval at which flow data is exported. Typically specified in seconds (e.g., "60s").

* `format` - (Required) The format of the exported flow data. Typically "ipfix" (IP Flow Information Export).

* `target` - (Required) One or more export targets. Each target is a block that supports:
  * `destination` - (Required) The destination address and port for the exported data (e.g., "192.168.1.100:4739").
  * `transport` - (Required) The transport protocol used for exporting. Typically "udp" or "tcp".

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Flow Export Policy (UUID).

## Import

Flow Export Policies can be imported using the `name`, e.g.,

```
$ terraform import psm_flow_export_policy.example example-policy
```

This will import the Flow Export Policy with the name "example-policy" into the Terraform resource named "example".

## Lifecycle Management

### Creation

When creating a new Flow Export Policy, the system will assign a UUID to the policy. This UUID becomes the `id` of the Terraform resource.

### Read

The resource can be read at any time to sync the Terraform state with the actual configuration in the PSM system.

### Update

Flow Export Policies can be updated after creation. All fields can be modified.

### Deletion

When a Flow Export Policy is deleted through Terraform, it will be removed from the PSM system.

## Notes

1. The `interval` should be specified in a format the PSM system understands, typically in seconds (e.g., "60s").

2. The `format` field is typically set to "ipfix", but confirm with your PSM system documentation for other possible values.

3. Multiple `target` blocks can be specified to export flow data to multiple destinations.

4. The `transport` field in the `target` block is typically either "udp" or "tcp". Ensure you choose the correct protocol as supported by your collector.

5. Changes to a Flow Export Policy may impact ongoing flow data collection and export. Plan updates carefully.

## Best Practices

1. Use meaningful names for your Flow Export Policies to easily identify their purpose.

2. Consider the impact of the export interval on your network and collector performance. Too short an interval may generate excessive traffic, while too long an interval may delay analysis.

3. Ensure that your flow collectors are properly configured to receive data in the format and transport protocol you specify.

4. Regularly review your Flow Export Policies to ensure they align with your current network monitoring needs.

5. When possible, use variable inputs for destinations to make your Terraform configurations more flexible and reusable across different environments.