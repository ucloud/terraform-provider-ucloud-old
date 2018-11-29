package ucloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/ucloud/ucloud-sdk-go/services/umem"
)

func TestAccUCloudKVStoreSlave_basic(t *testing.T) {
	var inst umem.URedisGroupSet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ucloud_kvstore_slave.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKVStoreSlaveDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKVStoreSlaveConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckKVStoreSlaveExists("ucloud_kvstore_instance.foo", &inst),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "instance_type", "redis-master-1"),
					resource.TestCheckResourceAttr("ucloud_kvstore_slave.foo", "instance_type", "redis-master-1"),
					resource.TestCheckResourceAttrSet("ucloud_kvstore_slave.foo", "master_id"),
				),
			},
		},
	})
}

func testAccCheckKVStoreSlaveExists(n string, inst *umem.URedisGroupSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("active standby redis id is empty")
		}

		client := testAccProvider.Meta().(*UCloudClient)
		ptr, err := client.describeActiveStandbyRedisById(rs.Primary.ID)

		log.Printf("[INFO] active standby redis id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		*inst = *ptr
		return nil
	}
}

func testAccCheckKVStoreSlaveDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ucloud_kvstore_slave" {
			continue
		}

		client := testAccProvider.Meta().(*UCloudClient)
		d, err := client.describeActiveStandbyRedisById(rs.Primary.ID)

		if err != nil {
			if isNotFoundError(err) {
				continue
			}
			return err
		}

		if d.GroupId != "" {
			return fmt.Errorf("active standby redis still exist")
		}
	}

	return nil
}

const testAccKVStoreSlaveConfig = `
data "ucloud_kvstore_parameter_groups" "default" {
	availability_zone = "cn-sh2-02"
	engine_version = "4.0"
}

resource "ucloud_kvstore_instance" "foo" {
	availability_zone = "cn-sh2-02"
	engine = "redis"
	engine_version = "3.2"
	instance_type = "redis-master-1"
	password = "2018_tfacc"
	name = "tf-acc-redis-master"
	tag = "tf-acc"
	parameter_group_id = "${data.ucloud_kvstore_parameter_groups.default.parameter_groups.0.id}"
	backup_begin_time = 3
}

resource "ucloud_kvstore_slave" "foo" {
	availability_zone = "cn-sh2-02"
	name = "tf-acc-redis-read-only-slave"
	instance_type = "redis-master-1"
	master_id = "${ucloud_kvstore_instance.foo.id}"
	parameter_group_id = "${data.ucloud_kvstore_parameter_groups.default.parameter_groups.0.id}"
}
`
