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
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "instance_storage", "30"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "memory_limit", "1000"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "engine", "mysql"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "engine_version", "5.5"),
				),
			},

			resource.TestStep{
				Config: testAccDBInstanceConfigTwo,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBInstanceExists("ucloud_db_instance.foo", &db),
					testAccCheckDBInstanceAttributes(&db),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "name", "tf-testDBInstance-basicUpdate"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "instance_storage", "50"),
					resource.TestCheckResourceAttr("ucloud_db_instance.foo", "memory_limit", "2000"),
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
	instance_storage = 30
	memory_limit = 1000
	engine = "mysql"
	engine_version = "5.5"
	password = "2018_UClou"
	port = 3306
	param_group_id = 2
}
`
const testAccDBInstanceConfigTwo = `
data "ucloud_zones" "default" {
}

resource "ucloud_db_instance" "foo" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	name = "tf-testDBInstance-basicUpdate"
	instance_storage = 50
	memory_limit = 2000
	engine = "mysql"
	engine_version = "5.5"
	password = "2018_UClou"
	port = 3306
	param_group_id = 2
}
`
