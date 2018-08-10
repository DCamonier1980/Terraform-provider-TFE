---
layout: "tfe"
page_title: "Terraform Enterprise: tfe_organization_token"
sidebar_current: "docs-tfe-organization-token"
description: |-
Generates a new organization token, replacing any existing token.
---

# tfe_organization_token

Generates a new organization token, replacing any existing token. This token
can be used to act as the organization service account.

## Example Usage

Basic usage:

```hcl
resource "tfe_organization_token" "token" {
	organization = "my-org-name"
}
```

## Argument Reference

The following arguments are supported:

* `organization` - (Required) Name of the organization.

## Attributes Reference

* `id` - The ID of the token.
* `token` - The generated token.
