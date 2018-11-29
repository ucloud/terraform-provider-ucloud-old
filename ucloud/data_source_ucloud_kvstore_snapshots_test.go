package ucloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccUCloudKVStoreSnapshotsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreSnapshotsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIDExists("data.ucloud_kvstore_snapshots.foo"),
				),
			},
		},
	})
}

const testAccKVStoreSnapshotsConfig = `
data "ucloud_kvstore_snapshots" "foo" {
}
`
