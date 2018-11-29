package ucloud

import (
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/ucloud/ucloud-sdk-go/services/umem"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

func (c *UCloudClient) describeActiveStandbyRedisById(id string) (*umem.URedisGroupSet, error) {
	conn := c.umemconn

	req := conn.NewDescribeURedisGroupRequest()
	req.GroupId = ucloud.String(id)

	resp, err := conn.DescribeURedisGroup(req)
	if err != nil {
		return nil, err
	}

	if resp == nil || len(resp.DataSet) < 1 {
		return nil, newNotFoundError(getNotFoundMessage("redis group", id))
	}

	return &resp.DataSet[0], nil
}

func (c *UCloudClient) describeDistributedRedisById(id string) (*umem.UMemSpaceSet, error) {
	conn := c.umemconn

	req := conn.NewDescribeUMemSpaceRequest()
	req.SpaceId = ucloud.String(id)

	resp, err := conn.DescribeUMemSpace(req)
	if err != nil {
		return nil, err
	}

	if resp == nil || len(resp.DataSet) < 1 {
		return nil, newNotFoundError(getNotFoundMessage("redis space", id))
	}

	return &resp.DataSet[0], nil
}

func (c *UCloudClient) describeActiveStandbyMemcacheById(id string) (*umem.UMemcacheGroupSet, error) {
	conn := c.umemconn

	req := conn.NewDescribeUMemcacheGroupRequest()
	req.GroupId = ucloud.String(id)

	resp, err := conn.DescribeUMemcacheGroup(req)
	if err != nil {
		return nil, err
	}

	if resp == nil || len(resp.DataSet) < 1 {
		return nil, newNotFoundError(getNotFoundMessage("memcache", id))
	}

	return &resp.DataSet[0], nil
}

func (c *UCloudClient) waitActiveStandbyRedisRunning(id string) error {
	refresh := func() (interface{}, string, error) {
		resp, err := c.describeActiveStandbyRedisById(id)
		if err != nil {
			if isNotFoundError(err) {
				return nil, "pending", nil
			}
			return nil, "", err
		}

		if resp.State != "Running" {
			return nil, "pending", nil
		}
		return resp, "ok", nil
	}

	return waitForMemoryInstance(refresh)
}

func (c *UCloudClient) waitDistributedRedisRunning(id string) error {
	refresh := func() (interface{}, string, error) {
		resp, err := c.describeDistributedRedisById(id)
		if err != nil {
			if isNotFoundError(err) {
				return nil, "pending", nil
			}
			return nil, "", err
		}

		if resp.State != "Running" {
			return nil, "pending", nil
		}
		return resp, "ok", nil
	}

	return waitForMemoryInstance(refresh)
}

func (c *UCloudClient) waitActiveStandbyMemcacheRunning(id string) error {
	refresh := func() (interface{}, string, error) {
		resp, err := c.describeActiveStandbyMemcacheById(id)
		if err != nil {
			if isNotFoundError(err) {
				return nil, "pending", nil
			}
			return nil, "", err
		}

		if resp.State != "Running" {
			return nil, "pending", nil
		}
		return resp, "ok", nil
	}

	return waitForMemoryInstance(refresh)
}

func waitForMemoryInstance(refresh func() (interface{}, string, error)) error {
	conf := resource.StateChangeConf{
		Timeout: 10 * time.Minute,
		Target:  []string{"ok"},
		Pending: []string{"pending"},
		Refresh: refresh,
	}

	_, err := conf.WaitForState()
	if err != nil {
		return err
	}

	return nil
}
