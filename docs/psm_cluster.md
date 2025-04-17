# Resource: psm_cluster

Manages the cluster configuration in the PSM system.

## Example Usage

```hcl
resource "psm_cluster" "example" {
  name             = "example-cluster"
  quorum_nodes     = ["node1", "node2", "node3"]
  virtual_ip       = "192.168.1.100"
  ntp_servers      = ["ntp1.example.com", "ntp2.example.com"]
  auto_admit_dscs  = true
  certs            = file("path/to/certs.pem")
  key              = file("path/to/key.pem")
  bootstrap_ipam_policy = "default-ipam-policy"
  sites            = ["site1", "site2"]
  
  labels = {
    environment = "production"
    team        = "networking"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the cluster.

* `quorum_nodes` - (Optional) A list of quorum node addresses.

* `virtual_ip` - (Optional) The virtual IP address for the cluster.

* `ntp_servers` - (Optional) A list of NTP server addresses.

* `auto_admit_dscs` - (Optional) Whether to automatically admit DSCs (Distributed Services Cards). Defaults to `false`.

* `certs` - (Optional) The certificates for the cluster in PEM format.

* `key` - (Optional, Sensitive) The private key for the cluster.

* `bootstrap_ipam_policy` - (Optional) The bootstrap IPAM policy for the cluster.

* `sites` - (Optional) A list of sites associated with the cluster.

* `labels` - (Optional) A map of labels to assign to the cluster.

* `certificate` - (Optional) The certificate for the cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the cluster (UUID).

## Import

Cluster can be imported using a placeholder ID, e.g.,

```
$ terraform import psm_cluster.example cluster
```

This will import the existing cluster configuration into your Terraform state.

## Lifecycle Management

### Creation

When creating a new cluster configuration, if a configuration already exists, the operation will fail. In such cases, use `terraform import` to manage the existing configuration.

### Update

The cluster configuration can be updated after creation. All fields can be modified.

### Deletion

Cluster resources cannot be deleted. Attempting to delete will result in an error.

## Notes

1. This resource manages a singleton cluster configuration. Only one instance of this resource should be defined in your Terraform configuration.

2. The `sites` argument is used to set the `system.multisite` label internally. The sites are joined with `|||` as a separator.

3. The `key` field is treated as sensitive data and will not be displayed in logs or console output.

4. When updating the cluster configuration, be cautious as changes may impact the entire PSM system.

5. If you need to manage an existing cluster configuration, use the `terraform import` command to bring it under Terraform management.

6. The `certificate` field is both settable and readable, while the `key` field is write-only for security reasons.

## Best Practices

1. Always use meaningful names for your cluster to easily identify it in both Terraform and the PSM system.

2. Regularly review and update your NTP servers to ensure time synchronization across the cluster.

3. Be cautious when modifying quorum nodes, as this can impact cluster stability.

4. Use Terraform variables or secure secret management solutions to handle sensitive data like certificates and keys.

5. Implement proper access controls and review processes before applying changes to the cluster configuration, as these changes can have system-wide impacts.