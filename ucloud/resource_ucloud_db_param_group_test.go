package ucloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/ucloud/ucloud-sdk-go/services/udb"
)

func TestAccUCloudDBParamGroup_basic(t *testing.T) {
	var dbPg udb.UDBParamGroupSet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ucloud_db_param_group.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDBParamGroupDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDBParamGroupConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBParamGroupExists("ucloud_db_param_group.foo", &dbPg),
					testAccCheckDBParamGroupAttributes(&dbPg),
					resource.TestCheckResourceAttr("ucloud_db_param_group.foo", "name", "tf-testDBParamGroup-basic"),
					resource.TestCheckResourceAttr("ucloud_db_param_group.foo", "description", "this is a test"),
					resource.TestCheckResourceAttr("ucloud_db_param_group.foo", "src_group_id", "18"),
					resource.TestCheckResourceAttr("ucloud_db_param_group.foo", "engine", "mysql"),
					resource.TestCheckResourceAttr("ucloud_db_param_group.foo", "engine_version", "5.7"),
				),
			},

			resource.TestStep{
				Config: testAccDBParamGroupConfigTwo,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckDBParamGroupExists("ucloud_db_param_group.foo", &dbPg),
					testAccCheckDBParamGroupAttributes(&dbPg),
					resource.TestCheckResourceAttr("ucloud_db_param_group.foo", "name", "tf-testDBParamGroup-basic"),
					resource.TestCheckResourceAttr("ucloud_db_param_group.foo", "description", "this is a test"),
					resource.TestCheckResourceAttr("ucloud_db_param_group.foo", "src_group_id", "18"),
					resource.TestCheckResourceAttr("ucloud_db_param_group.foo", "engine", "mysql"),
					resource.TestCheckResourceAttr("ucloud_db_param_group.foo", "engine_version", "5.7"),
				),
			},
		},
	})

}

func testAccCheckDBParamGroupExists(n string, dbPg *udb.UDBParamGroupSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("db param group id is empty")
		}

		client := testAccProvider.Meta().(*UCloudClient)
		zone := rs.Primary.Attributes["availability_zone"]
		ptr, err := client.describeDBParamGroupById(rs.Primary.ID, zone)

		log.Printf("[INFO] db param group id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		*dbPg = *ptr
		return nil
	}
}

func testAccCheckDBParamGroupAttributes(dbPg *udb.UDBParamGroupSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if dbPg.GroupId == 0 {
			return fmt.Errorf("db param group id is empty")
		}
		return nil
	}
}

func testAccCheckDBParamGroupDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ucloud_db_param_group" {
			continue
		}

		client := testAccProvider.Meta().(*UCloudClient)
		zone := rs.Primary.Attributes["availability_zone"]
		d, err := client.describeDBParamGroupById(rs.Primary.ID, zone)

		// Verify the error is what we want
		if err != nil {
			if isNotFoundError(err) {
				continue
			}
			return err
		}

		if d.GroupId != 0 {
			return fmt.Errorf("udb param group still exist")
		}
	}

	return nil
}

const testAccDBParamGroupConfig = `
data "ucloud_zones" "default" {
}

resource "ucloud_db_param_group" "foo" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	name = "tf-testDBParamGroup-basic"
	src_group_id = "18"
	description = "this is a test"
	engine = "mysql"
	engine_version = "5.7"
} 
`

const testAccDBParamGroupConfigTwo = `
data "ucloud_zones" "default" {
}

resource "ucloud_db_param_group" "foo" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	name = "tf-testDBParamGroup-basicUpdate"
	description = "this is a test"
	src_group_id = "18"
	engine = "mysql"
	engine_version = "5.7"
	parameter_input {
		key = max_connections
		value = 3000
	}
	parameter_input {
		key = slow_query_log
		value = 1
	}
}
`

func Test_replaceConfigContent(t *testing.T) {
	type args struct {
		content string
		data    map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceConfigContent(tt.args.content, tt.args.data); got != tt.want {
				t.Errorf("replaceConfigContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
