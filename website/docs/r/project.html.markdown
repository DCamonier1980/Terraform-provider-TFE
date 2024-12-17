---
layout: "tfe"
page_title: "Terraform Enterprise: tfe_project"
description: |-
Manages projects.
---

# tfe_project

Provides a project resource.

## Example Usage

Basic usage:

```hcl
resource "tfe_organization" "test-organization" {
  name  = "my-org-name"
  email = "admin@company.com"
}

resource "tfe_project" "test" {
  organization = tfe_organization.test-organization.name
  name = "projectname"
}
```

With tag bindings:

```hcl
resource "tfe_organization" "test-organization" {
  name  = "my-org-name"
  email = "admin@company.com"
}

resource "tfe_project" "test" {
  organization = tfe_organization.test-organization.name
  name         = "projectname"
  tag_bindings = {
      cost_center = "infrastructure"
      team = "platform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the project.
    *  TFE versions v202404-2 and earlier support between 3-36 characters
    *  TFE versions v202405-1 and later support between 3-40 characters
* `organization` - (Optional) Name of the organization. If omitted, organization must be defined in the provider config.
* `description` - (Optional) A description for the project.
* `tag_bindings` - (Optional) A map of key-value tags to add to the project.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The project ID.

## Import

Projects can be imported; use `<PROJECT ID>` as the import ID. For example:

```shell
terraform import tfe_project.test prj-niVoeESBXT8ZREhr
```
