package ucloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/ucloud/ucloud-sdk-go/services/udb"
)

func TestAccUCloudDBParameterGroup_basic(t *testing.T) {
	var dbPg udb.UDBParamGroupSet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ucloud_db_parameter_group.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDBParameterGroupDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBParameterGroupConfigBasic,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBParameterGroupExists("ucloud_db_parameter_group.foo", &dbPg),
					testAccCheckDBParameterGroupAttributes(&dbPg),
					resource.TestCheckResourceAttr("ucloud_db_parameter_group.foo", "name", "tf-testDBParameterGroup-basic"),
					resource.TestCheckResourceAttr("ucloud_db_parameter_group.foo", "description", "this is a test"),
					resource.TestCheckResourceAttr("ucloud_db_parameter_group.foo", "engine", "mysql"),
					resource.TestCheckResourceAttr("ucloud_db_parameter_group.foo", "engine_version", "5.7"),
				),
			},
		},
	})
}

func TestAccUCloudDBParameterGroup_key(t *testing.T) {
	var dbPg udb.UDBParamGroupSet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ucloud_db_parameter_group.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDBParameterGroupDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBParameterGroupConfigKey,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBParameterGroupExists("ucloud_db_parameter_group.foo", &dbPg),
					testAccCheckDBParameterGroupAttributes(&dbPg),
					resource.TestCheckResourceAttr("ucloud_db_parameter_group.foo", "name", "tf-testDBParameterGroup-key"),
					resource.TestCheckResourceAttr("ucloud_db_parameter_group.foo", "description", "this is a test"),
					resource.TestCheckResourceAttr("ucloud_db_parameter_group.foo", "engine", "mysql"),
					resource.TestCheckResourceAttr("ucloud_db_parameter_group.foo", "engine_version", "5.7"),
				),
			},
		},
	})

}

func testAccCheckDBParameterGroupExists(n string, dbPg *udb.UDBParamGroupSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("db parameter group id is empty")
		}

		client := testAccProvider.Meta().(*UCloudClient)
		zone := rs.Primary.Attributes["availability_zone"]
		ptr, err := client.describeDBParameterGroupByIdAndZone(rs.Primary.ID, zone)

		log.Printf("[INFO] db parameter group id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		*dbPg = *ptr
		return nil
	}
}

func testAccCheckDBParameterGroupAttributes(dbPg *udb.UDBParamGroupSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if dbPg.GroupId == 0 {
			return fmt.Errorf("db parameter group id is empty")
		}
		return nil
	}
}

func testAccCheckDBParameterGroupDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ucloud_db_param_group" {
			continue
		}

		client := testAccProvider.Meta().(*UCloudClient)
		zone := rs.Primary.Attributes["availability_zone"]
		d, err := client.describeDBParameterGroupByIdAndZone(rs.Primary.ID, zone)

		// Verify the error is what we want
		if err != nil {
			if isNotFoundError(err) {
				continue
			}
			return err
		}

		if d.GroupId != 0 {
			return fmt.Errorf("db parameter group still exist")
		}
	}

	return nil
}

const testAccDBParameterGroupConfigBasic = `
data "ucloud_zones" "default" {
}

data "ucloud_db_parameter_groups" "default" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	multi_az = "false"
	engine = "mysql"
	engine_version = "5.7"
}

resource "ucloud_db_parameter_group" "foo" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	name = "tf-testDBParameterGroup-basic"
	src_group_id = "${data.ucloud_db_parameter_groups.default.parameter_groups.0.id}"
	description = "this is a test"
	engine = "mysql"
	engine_version = "5.7"
} 
`

const testAccDBParameterGroupConfigKey = `
data "ucloud_zones" "default" {
}

data "ucloud_db_parameter_groups" "default" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	multi_az = "false"
	engine = "mysql"
	engine_version = "5.7"
}

resource "ucloud_db_parameter_group" "foo" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	name = "tf-testDBParameterGroup-key"
	description = "this is a test"
	src_group_id = "${data.ucloud_db_parameter_groups.default.parameter_groups.0.id}"
	engine = "mysql"
	engine_version = "5.7"
	parameter_input {
		key = "max_connections"
		value = "3000"
	}
	parameter_input {
		key = "slow_query_log"
		value = "1"
	}
}
`
