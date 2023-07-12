// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccMonitoringMonitoredProject_monitoringMonitoredProjectBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"project_id":    envvar.GetTestProjectFromEnv(),
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringMonitoredProjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringMonitoredProject_monitoringMonitoredProjectBasicExample(context),
			},
			{
				ResourceName:            "google_monitoring_monitored_project.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metrics_scope"},
			},
		},
	})
}

func testAccMonitoringMonitoredProject_monitoringMonitoredProjectBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_monitoring_monitored_project" "primary" {
  metrics_scope = "%{project_id}"
  name          = google_project.basic.project_id
}

resource "google_project" "basic" {
  project_id = "tf-test-m-id%{random_suffix}"
  name       = "tf-test-m-id%{random_suffix}-display"
  org_id     = "%{org_id}"
}
`, context)
}

func testAccCheckMonitoringMonitoredProjectDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_monitoring_monitored_project" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{MonitoringBasePath}}v1/locations/global/metricsScopes/{{metrics_scope}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})

			rName := tpgresource.GetResourceNameFromSelfLink(rs.Primary.Attributes["name"])
			project, err := config.NewResourceManagerClient(config.UserAgent).Projects.Get(rName).Do()
			rName = strconv.FormatInt(project.ProjectNumber, 10)

			for _, monitoredProject := range res["monitoredProjects"].([]any) {
				if strings.HasSuffix(monitoredProject.(map[string]any)["name"].(string), rName) {
					return fmt.Errorf("MonitoringMonitoredProject still exists at %s", url)
				}
			}
		}

		return nil
	}
}
