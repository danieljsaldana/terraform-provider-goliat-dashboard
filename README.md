# Terraform Provider for Goliat Dashboard

The Terraform provider for [Goliat Dashboard](https://github.com/danieljsaldana/goliat-dashboard) allows for the management and automation of resources within Goliat Dashboard, facilitating Infrastructure as Code (IaC) for your project.

## Overview

Goliat Dashboard is a comprehensive solution for monitoring and managing services. This provider enables users to create and manage organizations and projects in Goliat Dashboard using Terraform.

## Features

- Manage organizations (`goliatdashboard_organization`).
- Manage projects (`goliatdashboard_project`).

## Prerequisites

- [Terraform](https://www.terraform.io/downloads.html) v1.0.0 or higher.
- An authentication token for Goliat Dashboard.
- Access to a running instance of Goliat Dashboard.

## Installation

1. Install the provider from the Terraform Registry or clone this repository and build it:
   ```bash
   go install github.com/your-username/terraform-provider-goliat-dashboard@latest
   ```

2. Configure the provider in your `main.tf` file:
   ```hcl
   terraform {
     required_providers {
       goliatdashboard = {
         source  = "your-username/goliatdashboard"
         version = "1.0.0"
       }
     }
   }
   ```

## Configuration

Include the provider configuration in your `main.tf` file:
```hcl
provider "goliatdashboard" {
  backend_url = "https://your-goliat-dashboard.com"
  token       = "your-authentication-token"
}
```

## Usage Examples

### Create an Organization

```hcl
resource "goliatdashboard_organization" "example" {
  name = "New Provider Organization"
  type = "providerOrganizations"
}
```

### Create a Project

```hcl
resource "goliatdashboard_project" "example" {
  organization = "example_organization_id"
  name         = "Example Project"
  description  = "This is an example project"
}
```

## Contributing

Contributions are welcome. To contribute, please open an issue or pull request on [Goliat Dashboard](https://github.com/danieljsaldana/goliat-dashboard).

## License

This project is licensed under the [MPL-2.0](LICENSE).