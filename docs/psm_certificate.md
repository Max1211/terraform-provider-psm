# Resource: psm_certificate

Manages SSL/TLS certificates in the PSM system.

## Example Usage

```hcl
resource "psm_certificate" "example" {
  name             = "example-cert"
  certificate_data = file("path/to/certificate.pem")
  private_key      = file("path/to/private_key.pem")
  description      = "Example SSL certificate"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the certificate. This must be unique within the PSM system.

* `certificate_data` - (Required) The certificate data in PEM format.

* `private_key` - (Optional) The private key associated with the certificate, in PEM format. This is sensitive information and will be stored securely.

* `description` - (Optional) A description of the certificate.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the certificate (same as `name`).

* `kind` - The kind of the resource (always "Certificate").

* `api_version` - The API version of the resource.

## Import

Certificates can be imported using the `name`, e.g.,

```
$ terraform import psm_certificate.example example-cert
```

This will import the certificate with the name "example-cert" into the Terraform resource named "example".

## Lifecycle Management

### Creation

When creating a new certificate, the resource will wait for up to 30 seconds for the certificate to become available in the PSM system. If the certificate is not available within this time, the creation will fail with a timeout error.

### Update

Certificates can be updated after creation. Note that updating a certificate will replace the existing certificate in the PSM system. Be cautious when updating certificates that are in use, as this may impact services relying on the certificate.

The `private_key` field is optional during updates. If not provided, the existing private key will be retained.

### Deletion

When a certificate is deleted through Terraform, it will be removed from the PSM system. If the certificate is still in use by any services, those services may be impacted.

## Security Considerations

* The `private_key` is treated as sensitive data. It is never output in logs or returned in API responses.
* Ensure that you manage the `private_key` securely in your Terraform configurations and state files.
* When possible, use secure methods to pass the certificate and private key data to Terraform, such as environment variables or secure secret management systems.

## Error Handling

* If a certificate with the specified name already exists during creation, the operation will fail.
* If a certificate is not found during read or delete operations, it will be removed from the Terraform state.
* Any unexpected HTTP status codes during API calls will result in an error being returned to Terraform.

## Best Practices

1. Use meaningful names for your certificates to easily identify them in both Terraform and the PSM system.
2. Always provide a description to document the purpose and usage of the certificate.
3. Regularly rotate your certificates to maintain security best practices.
4. Use Terraform's `sensitive` function when passing certificate data and private keys as variables to prevent accidental exposure in logs.