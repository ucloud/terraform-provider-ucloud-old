package ucloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccUCloudDBParamGroupsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDBParamGroupsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ucloud_db_param_groups.foo"),
					resource.TestCheckResourceAttr("data.ucloud_db_param_groups.foo", "instances.#", "2"),
				),
			},
		},
	})
}

const testAccDataDBParamGroupsConfig = `
resource "ucloud_db_param_group" "foo" {
	count = 2

	availability_zone = "cn-sh2-02"

	name = "testAccDBParamGroups"
	src_group_id = "18"
	description = "this is a test"
	engine = "mysql"
	engine_version = "5.7"
}

data "ucloud_db_param_groups" "foo" {
	ids = ["${ucloud_db_param_group.foo.*.id}"]
}
`
