---
layout: "tfe"
page_title: "Terraform Enterprise: tfe_workspace"
sidebar_current: "docs-datasource-tfe-workspace-x"
description: |-
  Get information on a workspace.
---

# Data Source: tfe_workspace

Use this data source to get information about a workspace.

## Example Usage

```hcl
data "tfe_workspace" "test" {
  name         = "my-workspace-name"
  organization = "my-org-name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the workspace.
* `organization` - (Required) Name of the organization.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The workspace ID.
* `auto_apply` - Indicates whether to automatically apply changes when a
  Terraform plan is successful.
* `file_triggers_enabled` - Indicates whether runs are triggered based on the changed files in a VCS push (if `true`) or always triggered on every push (if `false`).
* `operations` - Indicates whether the workspace is using remote execution mode.
* `queue_all_runs` - Indicates whether all runs should be queued.
* `ssh_key_id` - The ID of an SSH key assigned to the workspace.
* `terraform_version` - The version of Terraform used for this workspace.
* `trigger_prefixes` - List of repository-root-relative paths which describe all locations to be tracked for changes.
* `vcs_repo` - Settings for the workspace's VCS repository.
* `working_directory` - A relative path that Terraform will execute within.

The `vcs_repo` block contains:

* `identifier` - A reference to your VCS repository in the format `:org/:repo`
  where `:org` and `:repo` refer to the organization and repository in your VCS
  provider.
* `ingress_submodules` - Indicates whether submodules should be fetched when
  cloning the VCS repository.
* `oauth_token_id` - OAuth token ID of the configured VCS connection.
