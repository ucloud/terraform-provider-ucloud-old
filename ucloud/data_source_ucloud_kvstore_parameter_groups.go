package ucloud

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	pumem "github.com/ucloud/ucloud-sdk-go/private/services/umem"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

func dataSourceUCloudKVStoreParameterGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUCloudKVStoreParameterGroupsRead,
		Schema: map[string]*schema.Schema{
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateNameRegex,
			},

			"engine_version": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateStringInChoices([]string{"3.0", "3.2", "4.0"}),
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"total_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"parameter_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"engine_version": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"status": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"create_time": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"update_time": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceUCloudKVStoreParameterGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.pumemconn

	req := conn.NewDescribeURedisConfigRequest()
	if v, ok := d.GetOk("availability_zone"); ok {
		req.Zone = ucloud.String(v.(string))
	}

	if v, ok := d.GetOk("engine_version"); ok {
		req.Version = ucloud.String(v.(string))
	}

	var groups []pumem.URedisConfigSet
	var err error
	if v, ok := d.GetOk("ids"); ok {
		groups, err = describeURedisConfigBatch(client, req, ifaceToStringSlice(v))
	} else {
		groups, err = describeURedisConfigAll(client, req)
	}

	if err != nil {
		return fmt.Errorf("error in read kvstore parameter group list, %s", err)
	}

	if v, ok := d.GetOk("name_regex"); ok {
		r := regexp.MustCompile(v.(string))
		groups = filterParameterGroups(groups, func(item *pumem.URedisConfigSet) bool {
			return r.MatchString(item.Name)
		})
	}

	if v, ok := d.GetOk("ids"); ok {
		groups = filterParameterGroups(groups, func(item *pumem.URedisConfigSet) bool {
			return checkStringIn(item.ConfigId, ifaceToStringSlice(v)) == nil
		})
	}

	d.Set("total_count", len(groups))
	err = dataSourceUCloudKVStoreParameterGroupsSave(d, groups)
	if err != nil {
		return fmt.Errorf("error in read kvstore parameter group list, %s", err)
	}

	return nil
}

func filterParameterGroups(items []pumem.URedisConfigSet, fn func(*pumem.URedisConfigSet) bool) []pumem.URedisConfigSet {
	var vL []pumem.URedisConfigSet
	for _, v := range items {
		if fn(&v) {
			vL = append(vL, v)
		}
	}
	return vL
}

func dataSourceUCloudKVStoreParameterGroupsSave(d *schema.ResourceData, items []pumem.URedisConfigSet) error {
	ids := []string{}
	data := []map[string]interface{}{}

	for _, item := range items {
		ids = append(ids, item.ConfigId)
		data = append(data, map[string]interface{}{
			"availability_zone": item.Zone,
			"id":                item.ConfigId,
			"name":              item.Name,
			"description":       item.Description,
			"engine_version":    item.Version,
			"status":            item.State,
			"create_time":       timestampToString(item.CreateTime),
			"update_time":       timestampToString(item.ModifyTime),
		})
	}

	d.SetId(hashStringArray(ids))
	if err := d.Set("parameter_groups", data); err != nil {
		return err
	}

	if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
		writeToFile(outputFile.(string), data)
	}

	return nil
}

func describeURedisConfigAll(client *UCloudClient, req *pumem.DescribeURedisConfigRequest) ([]pumem.URedisConfigSet, error) {
	groups := []pumem.URedisConfigSet{}
	limit := 100
	offset := 0
	for {
		req.Limit = ucloud.Int(limit)
		req.Offset = ucloud.Int(offset)
		resp, err := client.pumemconn.DescribeURedisConfig(req)
		if err != nil {
			return nil, err
		}

		if resp == nil || len(resp.DataSet) < 1 {
			break
		}

		groups = append(groups, resp.DataSet...)
		if len(resp.DataSet) < limit {
			break
		}

		offset = offset + limit
	}
	return groups, nil
}

func describeURedisConfigBatch(client *UCloudClient, req *pumem.DescribeURedisConfigRequest, ids []string) ([]pumem.URedisConfigSet, error) {
	groups := []pumem.URedisConfigSet{}

	for _, id := range ids {
		req.ConfigId = ucloud.String(id)
		resp, err := client.pumemconn.DescribeURedisConfig(req)
		if err != nil {
			return nil, err
		}

		if resp == nil || len(resp.DataSet) < 1 {
			continue
		}

		groups = append(groups, resp.DataSet[0])
	}
	return groups, nil
}
