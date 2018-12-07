package ucloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/ucloud/ucloud-sdk-go/services/umem"
)

func TestAccUCloudActiveStandbyRedis_basic(t *testing.T) {
	var inst umem.URedisGroupSet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ucloud_kvstore_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckActiveStandbyRedisDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccActiveStandbyRedisConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckActiveStandbyRedisExists("ucloud_kvstore_instance.foo", &inst),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "tag", "tf-acc"),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "name", "tf-acc-redis"),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "instance_type", "redis-master-1"),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "engine", "redis"),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "engine_version", "4.0"),
				),
			},

			resource.TestStep{
				Config: testAccActiveStandbyRedisConfigUpdate,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckActiveStandbyRedisExists("ucloud_kvstore_instance.foo", &inst),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "tag", "tf-acc"),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "name", "tf-acc-redis-renamed"),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "instance_type", "redis-master-2"),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "engine", "redis"),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "engine_version", "4.0"),
				),
			},
		},
	})
}

func TestAccUCloudActiveStandbyMemcache_basic(t *testing.T) {
	var inst umem.UMemcacheGroupSet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ucloud_kvstore_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckActiveStandbyMemcacheDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccActiveStandbyMemcacheConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckActiveStandbyMemcacheExists("ucloud_kvstore_instance.foo", &inst),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "name", "tf-acc-memcache"),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "engine", "memcache"),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "engine_version", "4.0"),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "instance_type", "memcache-master-1"),
				),
			},
		},
	})
}

func TestAccUCloudDistributedRedis_basic(t *testing.T) {
	var inst umem.UMemSpaceSet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "ucloud_kvstore_instance.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDistributedRedisDestroy,

		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDistributedRedisConfig,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckDistributedRedisExists("ucloud_kvstore_instance.foo", &inst),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "instance_type", "redis-distributed-16"),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "name", "tf-acc-redis"),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "engine", "redis"),
				),
			},

			resource.TestStep{
				Config: testAccDistributedRedisConfigUpdate,

				Check: resource.ComposeTestCheckFunc(
					testAccCheckDistributedRedisExists("ucloud_kvstore_instance.foo", &inst),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "instance_type", "redis-distributed-20"),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "name", "tf-acc-redis-renamed"),
					resource.TestCheckResourceAttr("ucloud_kvstore_instance.foo", "engine", "redis"),
				),
			},
		},
	})
}

func testAccCheckActiveStandbyRedisExists(n string, inst *umem.URedisGroupSet) resource.TestCheckFunc {
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

func testAccCheckActiveStandbyRedisDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ucloud_kvstore_instance" {
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

func testAccCheckDistributedRedisExists(n string, inst *umem.UMemSpaceSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("active standby redis id is empty")
		}

		client := testAccProvider.Meta().(*UCloudClient)
		ptr, err := client.describeDistributedRedisById(rs.Primary.ID)

		log.Printf("[INFO] distributed redis id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		*inst = *ptr
		return nil
	}
}

func testAccCheckDistributedRedisDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ucloud_kvstore_instance" {
			continue
		}

		client := testAccProvider.Meta().(*UCloudClient)
		d, err := client.describeDistributedRedisById(rs.Primary.ID)

		if err != nil {
			if isNotFoundError(err) {
				continue
			}
			return err
		}

		if d.SpaceId != "" {
			return fmt.Errorf("distributed redis still exist")
		}
	}

	return nil
}

func testAccCheckActiveStandbyMemcacheExists(n string, inst *umem.UMemcacheGroupSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("active standby memcache id is empty")
		}

		client := testAccProvider.Meta().(*UCloudClient)
		ptr, err := client.describeActiveStandbyMemcacheById(rs.Primary.ID)

		log.Printf("[INFO] active standby memcache id %#v", rs.Primary.ID)

		if err != nil {
			return err
		}

		*inst = *ptr
		return nil
	}
}

func testAccCheckActiveStandbyMemcacheDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ucloud_kvstore_instance" {
			continue
		}

		client := testAccProvider.Meta().(*UCloudClient)
		d, err := client.describeActiveStandbyMemcacheById(rs.Primary.ID)

		if err != nil {
			if isNotFoundError(err) {
				continue
			}
			return err
		}

		if d.GroupId != "" {
			return fmt.Errorf("active standby memcache still exist")
		}
	}

	return nil
}

const testAccActiveStandbyRedisConfig = `
data "ucloud_zones" "default" {
}

data "ucloud_kvstore_parameter_groups" "default" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	engine_version = "4.0"
}

resource "ucloud_kvstore_instance" "foo" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	engine = "redis"
	engine_version = "4.0"
	instance_type = "redis-master-1"
	password = "2018_tfacc"
	name = "tf-acc-redis"
	tag = "tf-acc"
	parameter_group_id = "${data.ucloud_kvstore_parameter_groups.default.parameter_groups.0.id}"
	backup_begin_time = 3
}
`

const testAccActiveStandbyRedisConfigUpdate = `
resource "ucloud_kvstore_instance" "foo" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	engine = "redis"
	engine_version = "4.0"
	instance_type = "redis-master-2"
	password = "2018_tfacc"
	name = "tf-acc-redis-renamed"
	tag = "tf-acc"
}
`

const testAccActiveStandbyMemcacheConfig = `
resource "ucloud_kvstore_instance" "foo" {
	availability_zone = "${data.ucloud_zones.default.zones.0.id}"
	name = "tf-acc-memcache"
	engine = "memcache"
	engine_version = "4.0"
	instance_type = "memcache-master-1"
}
`

const testAccDistributedRedisConfig = `
resource "ucloud_kvstore_instance" "foo" {
	availability_zone = "cn-bj2-04"
	name = "tf-acc-redis"
	engine = "redis"
	instance_type = "redis-distributed-16"
}
`

const testAccDistributedRedisConfigUpdate = `
resource "ucloud_kvstore_instance" "foo" {
	availability_zone = "cn-bj2-04"
	name = "tf-acc-redis-renamed"
	engine = "redis"
	instance_type = "redis-distributed-20"
}
`
