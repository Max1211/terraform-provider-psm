# Resource: psm_certificate

Manages SSL/TLS certificates in the PSM system.

## Example Usage (with CA certificate)

```hcl
resource "psm_certificate" "example" {
  name             = "example-cert"
  certificate_data = file("path/to/certificate.pem")
  description      = "Example SSL certificate"
}
```

## Example Usage (with certificate and key)

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

->
 [`private_key`](#private_key) must be in PKCS#1 format.


## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the certificate.

* `kind` - The kind of the resource.

* `api_version` - The API version of the resource.

## Import

Certificates can be imported using the `name`, e.g.,

```hcl
$ terraform import psm_certificate.example example-cert
```