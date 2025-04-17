# Terraform Provider Documentation

## Resource: psm_workload

The `psm_workload` resource allows you to manage Workloads in PSM (Pensando Service Mesh).

### Example Usage

```hcl
resource "psm_workload" "example" {
  name       = "example-workload"
  ip_address = "192.168.1.100"
  vlan_id    = 10
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Workload. This must be unique within the tenant.

* `ip_address` - (Optional) The IP address assigned to the Workload. Defaults to "default".

* `vlan_id` - (Optional) The VLAN ID for the Workload. Defaults to 0.

### Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the Workload.

### Import

Workloads can be imported using the `name`, e.g.,

```
$ terraform import psm_workload.example example-workload
```

### Notes

* The tenant for the Workload is always set to "default" in the current implementation.
* The `Read` operation populates the Terraform state with the `name`, `ip_address`, and `vlan_id` attributes.
* This resource does not support updates. To change any attribute, you need to recreate the resource.

### Limitations

* The current implementation only supports a single interface per Workload.
* Many Workload properties (such as `host-name`, `mac-address`, `micro-seg-vlan`, `network`, `vni`, and `migration-timeout`) are not exposed through this resource.
* Error handling is implemented for non-200 HTTP status codes, but specific error messages from the API are not parsed or exposed in detail.

### Error Handling

* If the API returns a non-200 status code during create, read, or delete operations, the provider will return an error with details about the failed operation.
* Any network errors or issues with marshalling/unmarshalling JSON will also result in an error being returned.

This resource allows for basic management of Workloads. For more advanced configurations or to access additional Workload properties, you may need to extend this resource or use other means to interact with the PSM API.
