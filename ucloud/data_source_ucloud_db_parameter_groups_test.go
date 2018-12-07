package ucloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccUCloudDBParameterGroupsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDBParameterGroupsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ucloud_db_parameter_groups.foo"),
					resource.TestCheckResourceAttr("data.ucloud_db_parameter_groups.foo", "parameter_groups.#", "2"),
				),
			},
		},
	})
}

const testAccDataDBParameterGroupsConfig = `
data "ucloud_zones" "default" {
}

resource "ucloud_db_parameter_group" "foo" {
	count = 2

	availability_zone = "${data.ucloud_zones.default.zones.0.id}"

	name = "testAccDBParameterGroups"
	src_group_id = "18"
	description = "this is a test"
	engine = "mysql"
	engine_version = "5.7"
}

data "ucloud_db_parameter_groups" "foo" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	ids = ["${ucloud_db_parameter_group.foo.*.id}"]
}
`
