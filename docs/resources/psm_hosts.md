# PSM Hosts Resource

The `psm_hosts` resource allows you to manage hosts in the PSM.  
A host represents a Baremetal or Hypervisor server.

## Example Usage DSS

```hcl
resource "psm_hosts" "example_dss" {
  name  = "dss-host-example"
  
  dscs {
    id  = "1111.2222.3333"
  }
  
  pnic_info {
    id  = "4444.5555.6666"
  }
}
```

## Example Usage PNIC
```hcl
resource "psm_hosts" "example_pnic" {
  name          = "pnic-host-example"

  pnic_info {
    mac_address = "ccdd.eeff.0010"
    name        = "PNIC-1"
  }

  pnic_info {
    mac_address = "ccdd.eeff.0011"
    name        = "PNIC-2"
  }

  pnic_info {
    mac_address = "ccdd.eeff.0012"
    name        = "PNIC-3"
  }

  pnic_info {
    mac_address = "ccdd.eeff.0013"
    name        = "PNIC-4"
  }
}
```

## Example Usage combined (DSS and PNIC)
```hcl
resource "psm_hosts" "example_combined_host" {
  name  = "combined-host-example"

  dscs {
    mac_address = "1122.3344.5566"
  }

  dscs {
    id = "DSS-002"
  }

  pnic_info {
    mac_address = "aabb.ccdd.eeff"
    name        = "PNIC-3"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the host.
* `dscs` - (Optional) A set of DSC (Distributed Services Card) configurations.
    * `id` - (Optional) The ID of the DSC.
    * `mac_address` - (Optional) The MAC address of the DSC.
* `pnic_info` - (Optional) A set of PNIC (Physical Network Interface Card) configurations.
    * `mac_address` - (Required) The MAC address of the PNIC.
    * `name` - (Required) The name of the PNIC.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `uuid` - The UUID of the host resource.

## Import

Hosts can be imported using the `name`, e.g.

```
$ terraform import psm_hosts.example <name>
```