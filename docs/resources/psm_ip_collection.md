# Resource: psm_ip_collection

Manages IP Collections in the PSM system. IP Collections are groups of IP addresses or other IP Collections that can be used in various network configurations and policies.

## Example Usage

```hcl
resource "psm_ip_collection" "example" {
  name = "example-ip-collection"
  
  addresses = [
    "192.168.1.0/24",
    "10.0.0.1",
    "10.1.1.1-10.1.1.100"
  ]
}
```

## Example Usage nested IP collection

```hcl

resource "psm_ip_collection" "example-1" {
  display_name = "example-ip-collection-1"
  
  addresses = [
    "10.1.1.1"
  ]
}

resource "psm_ip_collection" "example-2" {
  display_name = "example-ip-collection-2"
  
  addresses = [
    "10.1.1.2"
  ]
}

resource "psm_ip_collection" "example" {
  display_name = "example-ip-collection"
  
  ip_collections = [
    "example-ip-collection-1",
    "example-ip-collection-2"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) The name of the IP Collection. This must be unique within the PSM system.

* `addresses` - (Optional) A list of IP addresses, CIDR blocks, or IP ranges to include in the collection.

* `ip_collections` - (Optional) A list of other IP Collection names to include in this collection.  
  
* `address_family` - (Optional) Address Family.  
  Default value: IPv4
  Possible values: `IPv4`, `IPv6`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Flow Export Policy (UUID).

## Import

IP Collections can be imported using the `name`, e.g.,

```text
terraform import psm_ip_collection.example example-ip-collection
```
