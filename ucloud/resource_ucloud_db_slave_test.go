package ucloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/ucloud/ucloud-sdk-go/services/udb"
)

func TestAccUCloudDBSlave_basic(t *testing.T) {
	var dbSlave udb.UDBInstanceSet
	var dbInstance udb.UDBInstanceSet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ucloud_db_slave.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDBSlaveDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBSlaveConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists("ucloud_db_instance.foo", &dbInstance),
					testAccCheckDBSlaveExists("ucloud_db_slave.foo", &dbSlave),
					testAccCheckDBSlaveAttributes(&dbSlave),
					resource.TestCheckResourceAttr("ucloud_db_slave.foo", "name", "tf-testDBInstance-slave"),
					resource.TestCheckResourceAttr("ucloud_db_slave.foo", "instance_storage", "20"),
					resource.TestCheckResourceAttr("ucloud_db_slave.foo", "instance_type", "mysql-basic-1"),
				),
			},

			resource.TestStep{
				Config: testAccDBSlaveConfigTwo,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBSlaveExists("ucloud_db_instance.foo", &dbInstance),
					testAccCheckDBSlaveExists("ucloud_db_slave.foo", &dbSlave),
					testAccCheckDBSlaveAttributes(&dbSlave),
					resource.TestCheckResourceAttr("ucloud_db_slave.foo", "name", "tf-testDBInstance-slave-update"),
					resource.TestCheckResourceAttr("ucloud_db_slave.foo", "instance_storage", "30"),
					resource.TestCheckResourceAttr("ucloud_db_slave.foo", "instance_type", "mysql-basic-2"),
				),
			},
		},
	})
}

func testAccCheckDBSlaveExists(n string, dbSlave *udb.UDBInstanceSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("db slave id is empty")
		}

		client := testAccProvider.Meta().(*UCloudClient)
		ptr, err := client.describeDBInstanceById(rs.Primary.ID)

		log.Printf("[INFO] db slave id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		*dbSlave = *ptr
		return nil
	}
}

func testAccCheckDBSlaveAttributes(dbSlave *udb.UDBInstanceSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if dbSlave.DBId == "" {
			return fmt.Errorf("db slave id is empty")
		}
		return nil
	}
}

func testAccCheckDBSlaveDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ucloud_db_slave" {
			continue
		}

		client := testAccProvider.Meta().(*UCloudClient)
		d, err := client.describeDBInstanceById(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil {
			if isNotFoundError(err) {
				continue
			}
			return err
		}

		if d.DBId != "" {
			return fmt.Errorf("db slave still exist")
		}
	}

	return nil
}

const testAccDBSlaveConfig = `
data "ucloud_zones" "default" {
}

data "ucloud_db_parameter_groups" "default" {
	multi_az = "false"
	engine = "mysql"
	engine_version = "5.7"
}

resource "ucloud_db_instance" "foo" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	name = "tf-testDBInstance-master"
	instance_storage = 20
	instance_type = "mysql-ha-1"
	engine = "mysql"
	engine_version = "5.7"
	password = "2018_UClou"
	parameter_group_id = "${data.ucloud_db_parameter_groups.default.parameter_groups.0.id}"
}

resource "ucloud_db_slave" "foo" {
	name = "tf-testDBInstance-slave"
	instance_storage = 20
	instance_type = "mysql-basic-1"
	password = "2018_UClou"
	parameter_group_id = "${data.ucloud_db_parameter_groups.default.parameter_groups.0.id}"
	master_id = "${ucloud_db_instance.foo.id}"
}
`
const testAccDBSlaveConfigTwo = `
data "ucloud_zones" "default" {
}

data "ucloud_db_parameter_groups" "default" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	multi_az = "false"
	engine = "mysql"
	engine_version = "5.7"
}

resource "ucloud_db_instance" "foo" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	name = "tf-testDBInstance-master-update"
	instance_storage = 30
	instance_type = "mysql-ha-2"
	engine = "mysql"
	engine_version = "5.7"
	password = "2018_UClou"
	parameter_group_id = "${data.ucloud_db_parameter_groups.default.parameter_groups.0.id}"
}

resource "ucloud_db_slave" "foo" {
	name = "tf-testDBInstance-slave-update"
	instance_storage = 30
	instance_type = "mysql-basic-2"
	password = "2018_UClou"
	parameter_group_id = "${data.ucloud_db_parameter_groups.default.parameter_groups.0.id}"
	master_id = "${ucloud_db_instance.foo.id}"
}
`
