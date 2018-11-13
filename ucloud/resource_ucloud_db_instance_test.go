package ucloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/ucloud/ucloud-sdk-go/services/udb"
)

func TestAccUCloudDBInstance_basic(t *testing.T) {
	var db udb.UDBInstanceSet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ucloud_db_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDBInstanceDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBInstanceConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists("ucloud_db_instance.foo", &db),
					testAccCheckDBInstanceAttributes(&db),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "name", "tf-testDBInstance-basic"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "instance_storage", "20"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "memory", "1"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "engine", "mysql"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "param_group_id", "18"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "engine_version", "5.7"),
				),
			},

			resource.TestStep{
				Config: testAccDBInstanceConfigTwo,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists("ucloud_db_instance.foo", &db),
					testAccCheckDBInstanceAttributes(&db),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "name", "tf-testDBInstance-basicUpdate"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "instance_storage", "30"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "param_group_id", "18"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "memory", "2"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "engine", "mysql"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "engine_version", "5.7"),
				),
			},
		},
	})
}

func TestAccUCloudDBInstance_slave(t *testing.T) {
	var db udb.UDBInstanceSet
	var dbTwo udb.UDBInstanceSet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ucloud_db_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDBInstanceDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBInstanceConfigSlave,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists("ucloud_db_instance.foo", &dbTwo),
					testAccCheckDBInstanceAttributes(&dbTwo),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "name", "tf-testDBInstance-master"),
					testAccCheckDBInstanceExists("ucloud_db_instance.bar", &db),
					testAccCheckDBInstanceAttributes(&db),
					resource.TestCheckResourceAttr("ucloud_db_instance.bar", "name", "tf-testDBInstance-slave"),
					resource.TestCheckResourceAttr("ucloud_db_instance.bar", "instance_storage", "20"),
					resource.TestCheckResourceAttr("ucloud_db_instance.bar", "memory", "1"),
					resource.TestCheckResourceAttr("ucloud_db_instance.bar", "engine", "mysql"),
					resource.TestCheckResourceAttr("ucloud_db_instance.bar", "param_group_id", "18"),
					resource.TestCheckResourceAttr("ucloud_db_instance.bar", "engine_version", "5.7"),
				),
			},

			resource.TestStep{
				Config: testAccDBInstanceConfigSlavePromote,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists("ucloud_db_instance.foo", &dbTwo),
					testAccCheckDBInstanceAttributes(&dbTwo),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "name", "tf-testDBInstance-master"),
					testAccCheckDBInstanceExists("ucloud_db_instance.bar", &db),
					testAccCheckDBInstanceAttributes(&db),
					resource.TestCheckResourceAttr("ucloud_db_instance.bar", "name", "tf-testDBInstance-promote"),
					resource.TestCheckResourceAttr("ucloud_db_instance.bar", "instance_storage", "20"),
					resource.TestCheckResourceAttr("ucloud_db_instance.bar", "memory", "1"),
					resource.TestCheckResourceAttr("ucloud_db_instance.bar", "engine", "mysql"),
					resource.TestCheckResourceAttr("ucloud_db_instance.bar", "param_group_id", "18"),
					resource.TestCheckResourceAttr("ucloud_db_instance.bar", "engine_version", "5.7"),
				),
			},
		},
	})

}

func testAccCheckDBInstanceExists(n string, db *udb.UDBInstanceSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("db instance id is empty")
		}

		client := testAccProvider.Meta().(*UCloudClient)
		ptr, err := client.describeDBInstanceById(rs.Primary.ID)

		log.Printf("[INFO] db instance id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		*db = *ptr
		return nil
	}
}

func testAccCheckDBInstanceAttributes(db *udb.UDBInstanceSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if db.DBId == "" {
			return fmt.Errorf("db instance id is empty")
		}
		return nil
	}
}

func testAccCheckDBInstanceDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ucloud_db_instance" {
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
			return fmt.Errorf("udb instance still exist")
		}
	}

	return nil
}

const testAccDBInstanceConfig = `
data "ucloud_zones" "default" {
}

resource "ucloud_db_instance" "foo" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	name = "tf-testDBInstance-basic"
	instance_storage = 20
	memory = 1
	engine = "mysql"
	engine_version = "5.7"
	password = "2018_UClou"
	port = 3306
	param_group_id = "18"
	instance_type = "SATA_SSD"
}
`
const testAccDBInstanceConfigTwo = `
data "ucloud_zones" "default" {
}

resource "ucloud_db_instance" "foo" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	name = "tf-testDBInstance-basicUpdate"
	instance_storage = 30
	memory = 2
	engine = "mysql"
	engine_version = "5.7"
	password = "2018_UClou"
	port = 3306
	param_group_id = "18"
	instance_type = "SATA_SSD"
}
`
const testAccDBInstanceConfigSlave = `
data "ucloud_zones" "default" {
}

resource "ucloud_db_instance" "foo" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	name = "tf-testDBInstance-master"
	instance_storage = 20
	memory = 1
	engine = "mysql"
	engine_version = "5.7"
	password = "2018_UClou"
	port = 3306
	param_group_id = "18"
	instance_type = "SATA_SSD"
}

resource "ucloud_db_instance" "bar" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	name = "tf-testDBInstance-slave"
	instance_storage = 20
	memory = 1
	engine = "mysql"
	engine_version = "5.7"
	password = "2018_UClou"
	port = 3306
	param_group_id = "18"
	instance_type = "SATA_SSD"
	master_id = "${ucloud_db_instance.foo.id}"
}
`
const testAccDBInstanceConfigSlavePromote = `
data "ucloud_zones" "default" {
}
resource "ucloud_db_instance" "foo" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	name = "tf-testDBInstance-master"
	instance_storage = 20
	memory = 1
	engine = "mysql"
	engine_version = "5.7"
	password = "2018_UClou"
	port = 3306
	param_group_id = "18"
	instance_type = "SATA_SSD"
}

resource "ucloud_db_instance" "bar" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	name = "tf-testDBInstance-promote"
	instance_storage = 20
	memory = 1
	engine = "mysql"
	engine_version = "5.7"
	password = "2018_UClou"
	port = 3306
	param_group_id = "18"
	instance_type = "SATA_SSD"
}
`
