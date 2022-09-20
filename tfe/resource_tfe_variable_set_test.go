package tfe

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTFEVariableSet_full(t *testing.T) {
	variableSet := &tfe.VariableSet{}
	rInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTFEVariableSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTFEVariableSet_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTFEVariableSetExists(
						"tfe_variable_set.foobar", variableSet),
					testAccCheckTFEVariableSetAttributes(variableSet),
					resource.TestCheckResourceAttr(
						"tfe_variable_set.foobar", "name", "variable_set_test"),
					resource.TestCheckResourceAttr(
						"tfe_variable_set.foobar", "description", "a test variable set"),
					resource.TestCheckResourceAttr(
						"tfe_variable_set.foobar", "global", "false"),
				),
			},
			{
				Config: testAccTFEVariableSet_update(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTFEVariableSetExists(
						"tfe_variable_set.foobar", variableSet),
					testAccCheckTFEVariableSetAttributesUpdate(variableSet),
					resource.TestCheckResourceAttr(
						"tfe_variable_set.foobar", "name", "variable_set_test_updated"),
					resource.TestCheckResourceAttr(
						"tfe_variable_set.foobar", "description", "another description"),
					resource.TestCheckResourceAttr(
						"tfe_variable_set.foobar", "global", "true"),
				),
			},
		},
	})
}

func TestAccTFEVariableSet_import(t *testing.T) {
	rInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTFEVariableSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTFEVariableSet_basic(rInt),
			},

			{
				ResourceName:        "tfe_variable_set.foobar",
				ImportState:         true,
				ImportStateIdPrefix: "",
				ImportStateVerify:   true,
			},
		},
	})
}

func testAccCheckTFEVariableSetExists(
	n string, variableSet *tfe.VariableSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		tfeClient := testAccProvider.Meta().(*tfe.Client)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		vs, err := tfeClient.VariableSets.Read(
			ctx,
			rs.Primary.ID,
			&tfe.VariableSetReadOptions{Include: &[]tfe.VariableSetIncludeOpt{tfe.VariableSetWorkspaces}},
		)
		if err != nil {
			return err
		}

		*variableSet = *vs

		return nil
	}
}

func testAccCheckTFEVariableSetAttributes(
	variableSet *tfe.VariableSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if variableSet.Name != "variable_set_test" {
			return fmt.Errorf("Bad name: %s", variableSet.Name)
		}
		if variableSet.Description != "a test variable set" {
			return fmt.Errorf("Bad description: %s", variableSet.Description)
		}
		if variableSet.Global != false {
			return fmt.Errorf("Bad global: %t", variableSet.Global)
		}

		return nil
	}
}

func testAccCheckTFEVariableSetAttributesUpdate(
	variableSet *tfe.VariableSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if variableSet.Name != "variable_set_test_updated" {
			return fmt.Errorf("Bad name: %s", variableSet.Name)
		}
		if variableSet.Description != "another description" {
			return fmt.Errorf("Bad description: %s", variableSet.Description)
		}
		if variableSet.Global != true {
			return fmt.Errorf("Bad global: %t", variableSet.Global)
		}

		return nil
	}
}

func testAccCheckTFEVariableSetDestroy(s *terraform.State) error {
	tfeClient := testAccProvider.Meta().(*tfe.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tfe_variable_set" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		_, err := tfeClient.VariableSets.Read(ctx, rs.Primary.ID, nil)
		if err == nil {
			return fmt.Errorf("Variable Set %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccTFEVariableSet_basic(rInt int) string {
	return fmt.Sprintf(`
resource "tfe_organization" "foobar" {
  name = "tst-terraform-%d"
  email = "admin@company.com"
}

resource "tfe_variable_set" "foobar" {
  name         = "variable_set_test"
  description  = "a test variable set"
  global       = false
  organization = tfe_organization.foobar.id
}`, rInt)
}

func testAccTFEVariableSet_update(rInt int) string {
	return fmt.Sprintf(`
resource "tfe_organization" "foobar" {
  name = "tst-terraform-%d"
  email = "admin@company.com"
}

resource "tfe_variable_set" "foobar" {
  name         = "variable_set_test_updated"
  description  = "another description"
  global       = true
  organization = tfe_organization.foobar.id
}`, rInt)
}
