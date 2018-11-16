package ucloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccUCloudDBInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDBInstancesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ucloud_db_instances.foo"),
					resource.TestCheckResourceAttr("data.ucloud_db_instances.foo", "db_instances.#", "1"),
				),
			},
		},
	})
}

const testAccDataDBInstancesConfig = `
resource "ucloud_db_instance" "foo" {

	availability_zone = "cn-sh2-02"
	name = "testAccDBInstances"
	instance_storage = 20
	memory = 1
	engine = "mysql"
	engine_version = "5.7"
	password = "2018_UClou"
	port = 3306
	param_group_id = "18"
	instance_type = "SATA_SSD"
}

data "ucloud_db_instances" "foo" {
	ids = ["${ucloud_db_instance.foo.*.id}"]
}
`
