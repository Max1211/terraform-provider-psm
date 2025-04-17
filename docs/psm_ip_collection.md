# Resource: psm_ip_collection

Manages IP Collections in the PSM system. IP Collections are groups of IP addresses or other IP Collections that can be used in various network configurations and policies.

## Example Usage

```hcl
resource "psm_ip_collection" "example" {
  name = "example-ip-collection"
  
  addresses = [
    "192.168.1.0/24",
    "10.0.0.1",
    "2001:db8::/32"
  ]
  
  ip_collections = [
    "other-ip-collection-1",
    "other-ip-collection-2"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the IP Collection. This must be unique within the PSM system.

* `addresses` - (Optional) A list of IP addresses, CIDR blocks, or IP ranges to include in the collection.

* `ip_collections` - (Optional) A list of other IP Collection names to include in this collection.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the IP Collection (same as `name`).

## Import

IP Collections can be imported using the `name`, e.g.,

```
$ terraform import psm_ip_collection.example example-ip-collection
```

This will import the IP Collection with the name "example-ip-collection" into the Terraform resource named "example".

## Lifecycle Management

### Creation

When creating a new IP Collection, you must provide a unique `name`. You can optionally specify `addresses` and `ip_collections` to populate the collection.

### Read

The resource can be read at any time to sync the Terraform state with the actual configuration in the PSM system.

### Update

IP Collections can be updated after creation. You can modify the `addresses` and `ip_collections` lists.

### Deletion

When an IP Collection is deleted through Terraform, it will be removed from the PSM system.

## Notes

1. The `name` field is used as the identifier for the IP Collection and cannot be changed after creation. If you need to rename an IP Collection, you must create a new resource and delete the old one.

2. IP addresses in the `addresses` list can be specified in various formats:
   - Individual IP addresses (e.g., "192.168.1.1")
   - CIDR notation (e.g., "10.0.0.0/24")
   - IP ranges (e.g., "192.168.1.1-192.168.1.10")
   - IPv6 addresses and networks are also supported

3. When specifying other IP Collections in the `ip_collections` list, use their names, not their IDs.

4. Be cautious when updating IP Collections that are in use by other resources or policies, as changes may have cascading effects.

## Best Practices

1. Use meaningful names for your IP Collections to easily identify their purpose.

2. Consider using variables for IP addresses and other IP Collection names to make your Terraform configurations more flexible and reusable across different environments.

3. Group related IP addresses into collections for easier management and policy application.

4. Regularly review and update your IP Collections to ensure they reflect your current network architecture.

5. When possible, use CIDR notation for network ranges to make your configurations more concise.

6. Be mindful of the hierarchy when nesting IP Collections (using the `ip_collections` attribute) to avoid circular references.

7. Document the purpose and contents of each IP Collection in your Terraform configurations or external documentation.
