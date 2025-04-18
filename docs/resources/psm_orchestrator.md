# Resource: psm_orchestrator

Manages Orchestrator integrations in the PSM system.

## Example Usage

```hcl
resource "psm_orchestrator" "example" {
  name     = "example-orchestrator"
  type     = "kubernetes"
  uri      = "https://kubernetes.example.com"
  username = "admin"
  password = "password123"

  ca_data = file("path/to/ca.crt")
  # OR
  # disable_server_authentication = true

  namespaces {
    name = "default"
    mode = "smartservicemonitored"
  }

  namespaces {
    name = "production"
    mode = "smartserviceunmonitored"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Orchestrator integration.

* `type` - (Required) The type of the Orchestrator. Typically "kubernetes".

* `uri` - (Required) The URI of the Orchestrator.

* `username` - (Required) The username for authentication.

* `password` - (Required) The password for authentication.

* `ca_data` - (Optional) The CA certificate data for server authentication.

* `disable_server_authentication` - (Optional) Whether to disable server authentication. Defaults to `true` if `ca_data` is not provided.

* `namespaces` - (Optional) A list of namespace configurations. If not provided, a default namespace "all_namespaces" with mode "smartservicemonitored" will be used. Each namespace block supports:
  * `name` - (Required) The name of the namespace.
  * `mode` - (Required) The mode of the namespace. Can be "smartservicemonitored" or "smartserviceunmonitored".

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Orchestrator integration (UUID).

## Import

Orchestrator integrations can be imported using the `name`, e.g.,

```
$ terraform import psm_orchestrator.example example-orchestrator
```

Note: The `password` field is not returned by the API during import. You will need to set this value manually in your Terraform configuration after import.

## Notes

1. Either `ca_data` or `disable_server_authentication` should be set. If neither is explicitly set, `disable_server_authentication` defaults to `true`.

2. The `password` is sensitive and will not be displayed in logs or console output.

3. If no `namespaces` are specified, a default namespace "all_namespaces" with mode "smartservicemonitored" will be used.

4. During import, if the Orchestrator has no namespaces specified, the default "all_namespaces" will be set in the Terraform state.

## Best Practices

1. Use meaningful names for your Orchestrator integrations to easily identify them.

2. Store sensitive information like passwords using Terraform's encrypted state or external secret management systems.

3. When possible, use `ca_data` for server authentication instead of disabling server authentication.

4. Regularly rotate credentials used for Orchestrator integrations.

5. Use variables for URIs, usernames, and other configuration elements to make your Terraform configurations more flexible and reusable across different environments.

6. Be cautious when changing the `type` or `name` of an existing Orchestrator integration, as these are `ForceNew` fields and will result in the destruction and recreation of the resource.

7. Review and adjust namespace configurations regularly to ensure they align with your current orchestration needs.
