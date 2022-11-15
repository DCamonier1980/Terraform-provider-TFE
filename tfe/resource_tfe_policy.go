package tfe

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-tfe"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTFEPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTFEPolicyCreate,
		Read:   resourceTFEPolicyRead,
		Update: resourceTFEPolicyUpdate,
		Delete: resourceTFEPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceTFEPolicyImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"organization": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"kind": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(tfe.Sentinel),
				ValidateFunc: validation.StringInSlice(
					[]string{
						string(tfe.OPA),
						string(tfe.Sentinel),
					}, false),
			},

			"query": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"policy": {
				Type:     schema.TypeString,
				Required: true,
			},

			"enforce_mode": {
				Type: schema.TypeString,
				Description: fmt.Sprintf(
					"The enforce_mode of the policy. For Sentinel, valid values are `%s`, `%s`, and `%s`. For OPA, Valid values are `%s`and `%s`",
					tfe.EnforcementHard, tfe.EnforcementSoft, tfe.EnforcementAdvisory,
					tfe.EnforcementMandatory, tfe.EnforcementAdvisory),
				Optional: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						string(tfe.EnforcementAdvisory),
						string(tfe.EnforcementHard),
						string(tfe.EnforcementSoft),
						string(tfe.EnforcementMandatory),
					},
					false,
				),
			},
		},
	}
}

func resourceTFEPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	tfeClient := meta.(*tfe.Client)

	// Get the name and organization.
	name := d.Get("name").(string)
	organization := d.Get("organization").(string)

	var kind string
	if vKind, ok := d.GetOk("kind"); ok {
		kind = vKind.(string)
	}

	// Setup common policy options
	options := &tfe.PolicyCreateOptions{
		Name: tfe.String(name),
		Kind: tfe.PolicyKind(kind),
	}

	if desc, ok := d.GetOk("description"); ok {
		options.Description = tfe.String(desc.(string))
	}

	//  Setup per-kind policy options
	switch tfe.PolicyKind(kind) {
	case tfe.Sentinel:
		options = createSentinelPolicyOptions(options, d, meta)
	case tfe.OPA:
		options = createOPAPolicyOptions(options, d, meta)
	default:
		return fmt.Errorf(
			"Unsupported policy kind %s: has to be one of [%s, %s]", kind, string(tfe.Sentinel), string(tfe.OPA))
	}

	log.Printf("[DEBUG] Create %s policy %s for organization: %s", kind, name, organization)
	policy, err := tfeClient.Policies.Create(ctx, organization, *options)
	if err != nil {
		return fmt.Errorf(
			"Error creating %s policy %s for organization %s: %w", kind, name, organization, err)
	}

	d.SetId(policy.ID)

	log.Printf("[DEBUG] Upload %s policy %s for organization: %s", kind, name, organization)
	err = tfeClient.Policies.Upload(ctx, policy.ID, []byte(d.Get("policy").(string)))
	if err != nil {
		return fmt.Errorf(
			"Error uploading %s policy %s for organization %s: %w", kind, name, organization, err)
	}

	return resourceTFEPolicyRead(d, meta)
}

func createOPAPolicyOptions(options *tfe.PolicyCreateOptions, d *schema.ResourceData, meta interface{}) *tfe.PolicyCreateOptions {
	name := d.Get("name").(string)
	path := name + ".rego"
	options.Enforce = []*tfe.EnforcementOptions{
		{
			Path: tfe.String(path),
			Mode: tfe.EnforcementMode(tfe.EnforcementLevel(d.Get("enforce_mode").(string))),
		},
	}
	if vQuery, ok := d.GetOk("query"); ok {
		options.Query = tfe.String(vQuery.(string))
	}
	return options
}

func createSentinelPolicyOptions(options *tfe.PolicyCreateOptions, d *schema.ResourceData, meta interface{}) *tfe.PolicyCreateOptions {
	name := d.Get("name").(string)
	path := name + ".rego"
	options.Enforce = []*tfe.EnforcementOptions{
		{
			Path: tfe.String(path),
			Mode: tfe.EnforcementMode(tfe.EnforcementLevel(d.Get("enforce_mode").(string))),
		},
	}
	return options
}

func resourceTFEPolicyRead(d *schema.ResourceData, meta interface{}) error {
	tfeClient := meta.(*tfe.Client)

	log.Printf("[DEBUG] Read policy: %s", d.Id())
	policy, err := tfeClient.Policies.Read(ctx, d.Id())
	if err != nil {
		if err == tfe.ErrResourceNotFound {
			log.Printf("[DEBUG] Policy %s does no longer exist", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Policy %s: %w", d.Id(), err)
	}

	// Update the config.
	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("kind", policy.Kind)

	if len(policy.Enforce) == 1 {
		d.Set("enforce_mode", string(policy.Enforce[0].Mode))
	}

	content, err := tfeClient.Policies.Download(ctx, policy.ID)
	if err != nil {
		return fmt.Errorf("Error downloading policy %s: %w", d.Id(), err)
	}
	d.Set("policy", string(content))

	return nil
}

func resourceTFEPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	tfeClient := meta.(*tfe.Client)

	if d.HasChange("description") || d.HasChange("enforce_mode") {
		// Create a new options struct.
		options := tfe.PolicyUpdateOptions{}

		if desc, ok := d.GetOk("description"); ok {
			options.Description = tfe.String(desc.(string))
		}

		path := d.Get("name").(string) + ".sentinel"
		vKind, ok := d.GetOk("kind")
		if ok {
			if vKind == tfe.OPA {
				path = d.Get("name").(string) + ".rego"
			}
		}
		if d.HasChange("enforce_mode") {
			options.Enforce = []*tfe.EnforcementOptions{
				{
					Path: tfe.String(path),
					Mode: tfe.EnforcementMode(tfe.EnforcementLevel(d.Get("enforce_mode").(string))),
				},
			}
		}

		log.Printf("[DEBUG] Update configuration for %s policy: %s", vKind, d.Id())
		_, err := tfeClient.Policies.Update(ctx, d.Id(), options)
		if err != nil {
			return fmt.Errorf(
				"Error updating configuration for %s policy %s: %w", vKind, d.Id(), err)
		}
	}

	if d.HasChange("policy") {
		vKind := d.Get("kind").(string)
		log.Printf("[DEBUG] Update %s policy: %s", vKind, d.Id())
		err := tfeClient.Policies.Upload(ctx, d.Id(), []byte(d.Get("policy").(string)))
		if err != nil {
			return fmt.Errorf("Error updating %s policy %s: %w", vKind, d.Id(), err)
		}
	}

	return resourceTFEPolicyRead(d, meta)
}

func resourceTFEPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	tfeClient := meta.(*tfe.Client)

	log.Printf("[DEBUG] Delete policy: %s", d.Id())
	err := tfeClient.Policies.Delete(ctx, d.Id())
	if err != nil {
		if err == tfe.ErrResourceNotFound {
			return nil
		}
		return fmt.Errorf("Error deleting policy %s: %w", d.Id(), err)
	}

	return nil
}

func resourceTFEPolicyImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := strings.SplitN(d.Id(), "/", 2)
	if len(s) != 2 {
		return nil, fmt.Errorf(
			"invalid policy import format: %s (expected <ORGANIZATION>/<POLICY ID>)",
			d.Id(),
		)
	}

	// Set the fields that are part of the import ID.
	d.Set("organization", s[0])
	d.SetId(s[1])

	return []*schema.ResourceData{d}, nil
}
