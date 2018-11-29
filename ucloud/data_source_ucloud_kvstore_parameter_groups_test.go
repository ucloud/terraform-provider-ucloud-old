package ucloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccUCloudKVStoreParameterGroupsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreParameterGroupsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ucloud_kvstore_parameter_groups.foo"),
				),
			},
		},
	})
}

const testAccKVStoreParameterGroupsConfig = `
data "ucloud_kvstore_parameter_groups" "foo" {
}
`
