---
page_title: "Resource: psm_mirror_session"
description: |-
  Manages Mirror Sessions in AMD Policy and Services Manager.
---

# Resource: psm_mirror_session

Manages Mirror Sessions to export raw traffic for monitoring and troubleshooting in AMD Policy and Services Manager.

## Example Usage

```terraform
resource "psm_mirror_session" "example" {
  name                      = "traffic-monitoring"
  span_id                   = 1
  packet_size               = 2048
  disabled                  = false
  policy_distribution_target = "default"
  
  collector {
    type           = "erspan_type_3"
    destination    = "1.1.1.1"
  }
  
  collector {
    type           = "erspan_type_3"
    destination    = "2.2.2.2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the mirror session. Must be 2-64 characters, containing only alphanumeric characters, hyphens, and underscores.
* `span_id` - (Required) SPAN ID for the mirror session. Valid values are from 1 to 1023.
* `packet_size` - (Optional) Maximum packet size for mirrored traffic. Valid values are from 64 to 2048 bytes. Defaults to 2048.
* `disabled` - (Optional) Whether the mirror session is disabled. Defaults to false.
* `policy_distribution_target` - (Optional) Policy distribution target for the mirror session. Defaults to "default".
* `collector` - (Required) Up to 2 collector blocks defining where to send mirrored traffic. Each block supports:
  * `type` - (Required) Type of collector. Currently only "erspan_type_3" is supported.
  * `destination` - (Required) IP address of the collector destination.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the Mirror Session.

## Import

Mirror Sessions can be imported using the name, e.g.,

```
$ terraform import psm_mirror_session.example traffic-monitoring
```