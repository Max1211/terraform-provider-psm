# Resource: psm_policy_distribution_target

Manages Policy Distribution Targets (PDTs) in the PSM system. PDTs define where policies should be distributed within the network.

## Example Usage

```hcl
resource "psm_policy_distribution_target" "example" {
  name = "example-pdt"
  dses = ["dse1", "dse2", "dse3"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Policy Distribution Target. This must be unique within the PSM system.

* `dses` - (Optional) A set of Distributed Services Engines (DSEs) where the policies should be distributed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Policy Distribution Target (UUID).

## Import

Policy Distribution Targets can be imported using the `name`, e.g.,

```
$ terraform import psm_policy_distribution_target.example example-pdt
```

## Lifecycle Management

### Creation

When creating a new Policy Distribution Target, you must provide a unique `name`. If `dses` are specified during creation, they will be set after the initial creation of the PDT.

### Read

The resource can be read at any time to sync the Terraform state with the actual configuration in the PSM system.

### Update

Policy Distribution Targets can be updated after creation. You can modify the `dses` set to add or remove DSEs from the distribution target.

### Deletion

When a Policy Distribution Target is deleted through Terraform, it will be removed from the PSM system.

## Notes

1. The `name` field is used as the identifier for the Policy Distribution Target and cannot be changed after creation. If you need to rename a PDT, you must create a new resource and delete the old one.

2. The `dses` field is a set, which means the order of the DSEs doesn't matter and duplicates are automatically removed.

3. If no `dses` are specified, the Policy Distribution Target will be created without any associated DSEs.

4. Changes to the Policy Distribution Target may affect where policies are distributed in your network, potentially impacting network behavior and security.

## Best Practices

1. Use meaningful names for your Policy Distribution Targets to easily identify their purpose and scope.

2. Consider using variables for DSE names to make your Terraform configurations more flexible and reusable across different environments.

3. Regularly review and update your Policy Distribution Targets to ensure they align with your current network topology and policy distribution needs.

4. Be cautious when modifying existing Policy Distribution Targets, especially in production environments. Consider the impact on policy distribution and network behavior.

5. Document the purpose and configuration of each Policy Distribution Target in your Terraform code comments or external documentation.

6. Use consistent naming conventions for your Policy Distribution Targets across your infrastructure.

7. When possible, group related DSEs into a single Policy Distribution Target for easier management and clearer policy distribution strategies.

8. Regularly audit your Policy Distribution Targets to ensure that policies are being distributed to the correct DSEs and that there are no unused or misconfigured targets.
