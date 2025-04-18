# Distributed Service Switch Resource

Manages a Distributed Service Switch (DSS) in the system.

## Example Usage

```hcl
resource "psm_dss" "example" {
  name                    = "example-dss"
  fwlog_policy_name       = "example-fwlog-policy"
  flow_export_policy_name = "example-flow-export-policy"
  
  labels = {
    environment = "production"
    department  = "networking"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Distributed Service Switch. This must be unique within the system. Changing this forces a new resource to be created.
  * The `name` is the DSS ID (Base MAC Address), e.g. 1234.1234.1234.

* `labels` - (Optional) A map of key/value labels to attach to the DSS.  
  * `system.psm.site` - (Optional) This is a reserved label which is used to assign the DSS to a site.

* `fwlog_policy_name` - (Optional) The name of the firewall log policy to attach to the DSS.

* `flow_export_policy_name` - (Optional) The name of the flow export policy to attach to the DSS.  

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `serial_num` - The serial number of the DSS.

* `primary_mac` - The primary MAC address of the DSS.

* `DSS_version` - The version of the DSS software.

* `DSS_sku` - The SKU of the DSS.

* `ip_address` - The IP address assigned to the DSS.

* `default_gw` - The default gateway for the DSS.

* `dns_servers` - A list of DNS servers configured for the DSS.

* `is_connected_to_psm` - Boolean indicating whether the DSS is connected to the PSM (Pensando Systems Manager).

* `host_name` - The hostname of the DSS.

* `dss_version` - The version of the Distributed Services Switch (DSS) associated with this DSS.

* `forwarding_profile` - The forwarding profile of the DSS.

* `security_policy_rule_scale_profile` - The security policy rule scale profile of the DSS.

* `dsms` - A list of Distributed Service Modules (DSMs) associated with this DSS. Each DSM has the following attributes:
  * `unit_id` - The unit ID of the DSM.
  * `mac_address` - The MAC address of the DSM.

## Import

Distributed Service Switchs can be imported using the `name`, e.g.,

```text
terraform import psm_distributed_service_Switch.example example-DSS
```

## Notes

* The Distributed Service Switch resource cannot be fully deleted from the system. When `terraform destroy` is called, it will only clear the labels associated with the DSS.
* Updating the `ip_address`, `default_gw`, or `dns_servers` fields will trigger an update of the DSS's IP configuration.
* The `fwlog_policy_name` and `flow_export_policy_name` fields can be updated after creation, which will change the associated policies for the DSS.
