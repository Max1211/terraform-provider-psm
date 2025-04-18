---
page_title: "PSM Provider"
subcategory: ""
description: |-
  Terraform provider for interacting with AMD Policy and Services Manager (PSM).
---

# PSM Provider

This provider allows you to interact with the AMD Policy and Services Manager (PSM), enabling infrastructure automation and configuration management of AMD PSM resources.

## Example Usage

```terraform
terraform {
  required_providers {
    psm = {
      source = "Max1211/psm"
      version = "0.5.16" # Use the latest version available
    }
  }
}

provider "psm" {
  user     = "admin"           # Username for PSM
  password = "your-password"   # Password for authentication
  server   = "psm.example.com" # PSM server address
  insecure = false             # Set to true to skip TLS verification
}

# Example resource
# resource "psm_resource_type" "example" {
#   # ...
# }
```

## Authentication

The PSM provider requires authentication credentials to interact with your PSM instance.

## Schema

### Required

- `user` (String) - Username used to authenticate with PSM
- `password` (String) - Password used for authentication
- `server` (String) - Hostname/IP address of the PSM server to connect to

### Optional

- `insecure` (Boolean) - Whether to skip TLS verification when connecting to the server (default: `false`)
