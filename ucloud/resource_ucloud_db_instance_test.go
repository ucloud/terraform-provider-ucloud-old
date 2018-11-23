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
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "instance_type", "mysql-ha-1"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "engine", "mysql"),
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
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "instance_type", "mysql-ha-2"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "engine", "mysql"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "engine_version", "5.7"),
				),
			},
		},
	})
}

func TestAccUCloudDBInstance_pgsql(t *testing.T) {
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
				Config: testAccDBInstanceConfigPgsql,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists("ucloud_db_instance.foo", &db),
					testAccCheckDBInstanceAttributes(&db),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "name", "tf-testDBInstance-pgsql"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "instance_storage", "20"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "instance_type", "postgresql-basic-1"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "engine", "postgresql"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "engine_version", "9.6"),
				),
			},

			resource.TestStep{
				Config: testAccDBInstanceConfigPgsqlTwo,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists("ucloud_db_instance.foo", &db),
					testAccCheckDBInstanceAttributes(&db),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "name", "tf-testDBInstance-pgsqlUpdate"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "instance_storage", "30"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "instance_type", "postgresql-basic-2"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "engine", "postgresql"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "engine_version", "9.6"),
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
			return fmt.Errorf("db instance still exist")
		}
	}

	return nil
}

const testAccDBInstanceConfig = `
data "ucloud_zones" "default" {
}

data "ucloud_db_param_groups" "default" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	region_flag = "false"
	engine = "mysql"
	engine_version = "5.7"
}

resource "ucloud_db_instance" "foo" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	name = "tf-testDBInstance-basic"
	instance_storage = 20
	instance_type = "mysql-ha-1"
	engine = "mysql"
	engine_version = "5.7"
	password = "2018_UClou"
	parameter_group_id = "${data.ucloud_db_param_groups.default.param_groups.0.id}"
}
`
const testAccDBInstanceConfigTwo = `
data "ucloud_zones" "default" {
}

data "ucloud_db_param_groups" "default" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	region_flag = "false"
	engine = "mysql"
	engine_version = "5.7"
}

resource "ucloud_db_instance" "foo" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	name = "tf-testDBInstance-basicUpdate"
	instance_storage = 30
	instance_type = "mysql-ha-2"
	engine = "mysql"
	engine_version = "5.7"
	password = "2018_UClou"
	parameter_group_id = "${data.ucloud_db_param_groups.default.param_groups.0.id}"
}
`
const testAccDBInstanceConfigPgsql = `
data "ucloud_zones" "default" {
}

data "ucloud_db_param_groups" "default" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	region_flag = "false"
	engine = "postgresql"
	engine_version = "9.6"
}

resource "ucloud_db_instance" "foo" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	name = "tf-testDBInstance-pgsql"
	instance_storage = 20
	instance_type = "postgresql-basic-1"
	engine = "postgresql"
	engine_version = "9.6"
	password = "2018_UClou"
	parameter_group_id = "${data.ucloud_db_param_groups.default.param_groups.0.id}"
}
`
const testAccDBInstanceConfigPgsqlTwo = `
data "ucloud_zones" "default" {
}

data "ucloud_db_param_groups" "default" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	region_flag = "false"
	engine = "postgresql"
	engine_version = "9.6"
}

resource "ucloud_db_instance" "foo" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	name = "tf-testDBInstance-pgsqlUpdate"
	instance_storage = 30
	instance_type = "postgresql-basic-2"
	engine = "postgresql"
	engine_version = "9.6"
	password = "2018_UClou"
	parameter_group_id = "${data.ucloud_db_param_groups.default.param_groups.0.id}"
}
`
