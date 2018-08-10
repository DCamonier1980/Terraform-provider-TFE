---
layout: "tfe"
page_title: "Terraform Enterprise: tfe_workspoce"
sidebar_current: "docs-tfe-workspace"
description: |-
  Workspaces represent running infrastructure managed by Terraform.
---

# tfe_workspoce

Provides a workspace resource.

## Example Usage

Basic usage:

```hcl
resource "tfe_workspoce" "my-workspace" {
	name = "my-workspace"
	organization = "my-organization"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the workspace.
* `organization` - (Required) Name of the organization.
* `auto_apply` - (Optional) Whether to automatically apply changes when a
  Terraform plan is successful. Defaults to `false`.
* `terraform_version` - (Optional) The version of Terraform to use for this
  workspace. Defaults to the latest available version.
* `working_directory` - (Optional) A relative path that Terraform will execute
  within.  Defaults to the root of your repository.
* `vcs_repo` - (Optional) Settings for the workspace's VCS repository.

The `vcs_repo` block supports:

* `identifier` - (Required) A reference to your VCS repository in the format
  `:org/:repo` where `:org` and `:repo` refer to the organization and repository
  in your VCS provider.
* `branch` - (Optional) The repository branch that Terraform will execute from.
  Default to `master`.
* `ingress_submodules` - (Optional) Whether submodules should be fetched when
  cloning the VCS repository. Defaults to `false`.
* `oauth_token_id` - (Required) Token ID of the VCS Connection (OAuth Conection
  + Token) to use.
