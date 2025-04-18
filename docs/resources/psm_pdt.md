# Resource: psm_pdt

Manages Policy Distribution Targets (PDTs) in the PSM system. PDTs define where policies should be distributed within the network.

## Example Usage

```hcl
resource "psm_pdt" "example" {
  name = "example-pdt"
  dses = ["1234.1234.1234", "5678.5678.5678"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Policy Distribution Target. This must be unique within the PSM system.

* `dses` - (Optional) A set of Distributed Services Engines (DSEs) where the policies should be distributed.  
  
  -> Consider using variables with meaningful names for each DSE.


## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Policy Distribution Target (UUID).

## Import

Policy Distribution Targets can be imported using the `name`, e.g.,

```
$ terraform import psm_pdt.example example-pdt
```