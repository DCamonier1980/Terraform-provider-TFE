package tfe

import (
	"fmt"
	"log"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTFETeamMembers() *schema.Resource {
	return &schema.Resource{
		Create: resourceTFETeamMembersCreate,
		Read:   resourceTFETeamMembersRead,
		Update: resourceTFETeamMembersUpdate,
		Delete: resourceTFETeamMembersDelete,

		Schema: map[string]*schema.Schema{
			"team_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"usernames": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceTFETeamMembersCreate(d *schema.ResourceData, meta interface{}) error {
	tfeClient := meta.(*tfe.Client)

	// Get the team ID and username..
	teamID := d.Get("team_id").(string)

	// Create a new options struct.
	options := tfe.TeamMemberAddOptions{}

	// Add all the users that need to be added.
	for _, username := range d.Get("usernames").(*schema.Set).List() {
		options.Usernames = append(options.Usernames, username.(string))
	}

	log.Printf("[DEBUG] Add users to team: %s", teamID)
	err := tfeClient.TeamMembers.Add(ctx, teamID, options)
	if err != nil {
		return fmt.Errorf("Error adding users to team %s: %v", teamID, err)
	}

	d.SetId(teamID)

	return nil
}

func resourceTFETeamMembersRead(d *schema.ResourceData, meta interface{}) error {
	tfeClient := meta.(*tfe.Client)

	// Get the team ID and username..
	teamID := d.Get("team_id").(string)

	log.Printf("[DEBUG] Read users from team: %s", teamID)
	users, err := tfeClient.TeamMembers.List(ctx, teamID)
	if err != nil {
		if err == tfe.ErrResourceNotFound {
			log.Printf("[DEBUG] Users do no longer exist")
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading users from team %s: %v", teamID, err)
	}

	usernames := resourceTFETeamMembers().Schema["usernames"].ZeroValue().(*schema.Set)
	for _, user := range users {
		usernames.Add(user.Username)
	}

	if usernames.Len() > 0 {
		d.Set("usernames", usernames)
	} else {
		log.Printf("[DEBUG] Users do no longer exist")
		d.SetId("")
	}

	return nil
}

func resourceTFETeamMembersUpdate(d *schema.ResourceData, meta interface{}) error {
	tfeClient := meta.(*tfe.Client)

	if d.HasChange("usernames") {
		old, new := d.GetChange("usernames")
		oldUsers := old.(*schema.Set).Difference(new.(*schema.Set))
		newUsers := new.(*schema.Set).Difference(old.(*schema.Set))

		// Get the team ID and username..
		teamID := d.Get("team_id").(string)

		// First add the new users.
		if newUsers.Len() > 0 {
			// Create a new options struct.
			options := tfe.TeamMemberAddOptions{}

			// Add all the users that need to be added.
			for _, username := range newUsers.List() {
				options.Usernames = append(options.Usernames, username.(string))
			}

			log.Printf("[DEBUG] Add users to team: %s", teamID)
			err := tfeClient.TeamMembers.Add(ctx, teamID, options)
			if err != nil {
				return fmt.Errorf("Error adding users to team %s: %v", teamID, err)
			}
		}

		// Then delete all the old users.
		if oldUsers.Len() > 0 {
			// Create a new options struct.
			options := tfe.TeamMemberRemoveOptions{}

			// Add all the users that need to be added.
			for _, username := range oldUsers.List() {
				options.Usernames = append(options.Usernames, username.(string))
			}

			log.Printf("[DEBUG] Remove users from team: %s", teamID)
			err := tfeClient.TeamMembers.Remove(ctx, teamID, options)
			if err != nil {
				return fmt.Errorf("Error removing users to team %s: %v", teamID, err)
			}
		}
	}

	return nil
}

func resourceTFETeamMembersDelete(d *schema.ResourceData, meta interface{}) error {
	tfeClient := meta.(*tfe.Client)

	log.Printf("[DEBUG] Retrieve users to remove from team: %s", d.Id())
	users, err := tfeClient.TeamMembers.List(ctx, d.Id())
	if err != nil {
		if err == tfe.ErrResourceNotFound {
			return nil
		}
		return fmt.Errorf("Error retrieving users to remove from team %s: %v", d.Id(), err)
	}

	// Create a new options struct.
	options := tfe.TeamMemberRemoveOptions{}

	// Add all the users that need to be removed.
	for _, user := range users {
		options.Usernames = append(options.Usernames, user.Username)
	}

	log.Printf("[DEBUG] Remove users from team: %s", d.Id())
	err = tfeClient.TeamMembers.Remove(ctx, d.Id(), options)
	if err != nil {
		return fmt.Errorf("Error removing users to team %s: %v", d.Id(), err)
	}

	return nil
}
