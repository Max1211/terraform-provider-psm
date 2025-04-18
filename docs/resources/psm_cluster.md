# Resource: psm_cluster

Manages the cluster configuration in the PSM system.

## Example Usage

```hcl
resource "psm_cluster" "example" {
  name             = "example-cluster"
  ntp_servers      = ["ntp1.example.com", "ntp2.example.com"]
  auto_admit_dscs  = true
  certs            = file("path/to/cert.pem")
  key              = file("path/to/key.pem")
  sites            = ["site1", "site2"]
  
  labels = {
    environment = "production"
    team        = "networking"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the PSM cluster.

* `ntp_servers` - (Optional) A list of NTP server addresses.

* `auto_admit_dscs` - (Optional) Whether to automatically admit DSCs (Distributed Services Cards). Defaults to `false`.

* `certs` - (Optional) The certificates for the cluster in PEM format.

* `key` - (Optional, Sensitive) The private key for the cluster.

* `sites` - (Optional) A list of sites associated with the cluster.

* `labels` - (Optional) A map of labels to assign to the cluster.

* `certificate` - (Optional) The certificate for the cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the cluster (UUID).

* `quorum_nodes` - (Optional) A list of quorum node addresses.

* `virtual_ip` - (Optional) The virtual IP address for the cluster.

* `bootstrap_ipam_policy` - (Optional) The bootstrap IPAM policy for the cluster.

## Import

Cluster can be imported using a placeholder ID, e.g.,

```
$ terraform import psm_cluster.example cluster
```

This will import the existing cluster configuration into your Terraform state.