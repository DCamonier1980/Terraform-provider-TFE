---
layout: "tfe"
page_title: "Terraform Enterprise: tfe_team_member"
sidebar_current: "docs-tfe-team-member"
description: |-
  Add or remove a user from a team.
---

# tfe_team_member

Add or remove a user from a team.

~> NOTE on managing team memberships: Terraform currently provides two resources
for managing team memberships. The [tfe_team_member](team_member.html) resource
can be used multiple times as it manages the team membership for a single user.
The [tfe_team_members](team_members.html) resource, on the other hand, is used
to manage all team memberships for a specific team and can only be used once.
Both resources cannot be used for the same team simultaneously.

## Example Usage

Basic usage:

```hcl
resource "tfe_team" "team" {
  name = "my-team-name"
  organization = "my-org-name"
}

resource "tfe_team_member" "member" {
	team_id = "${tfe_team.team.id}"
  username = "sander"
}
```

## Argument Reference

The following arguments are supported:

* `team_id` - (Required) ID of the team.
* `username` - (Required) Name of the user to add.
