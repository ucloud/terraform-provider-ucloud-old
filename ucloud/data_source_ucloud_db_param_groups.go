package ucloud

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ucloud/ucloud-sdk-go/services/udb"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

func dataSourceUCloudDBParamGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUCloudDBParamGroupsRead,

		Schema: map[string]*schema.Schema{
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ForceNew: true,
			},

			"class_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateStringInChoices([]string{"sql", "nosql", "postgresql"}),
			},

			"engine": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateStringInChoices([]string{"mysql", "percona", "postgresql"}),
			},

			"engine_version": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateStringInChoices([]string{"5.1", "5.5", "5.6", "5.7", "9.4", "9.6", "10.4"}),
			},

			"region_flag": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"total_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"param_groups": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"db_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"modifiable": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},

						"param_member": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},

									"value": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},

									"value_type": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},

									"allowed_value": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},

									"apply_type": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},

									"modifiable": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},

									"format_type": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceUCloudDBParamGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.udbconn
	var fetched []udb.UDBParamGroupSet
	var filtered []udb.UDBParamGroupSet
	var paramGroups []udb.UDBParamGroupSet
	var totalCount int
	var limit int = 100
	var offset int = 0

	if ids, ok := d.GetOk("ids"); ok && len(ids.([]interface{})) > 0 {
		for _, id := range ifaceToStringSlice(ids) {
			dbPg, err := client.describeDBParamGroupById(id)
			if err != nil {
				return fmt.Errorf("error in read db param group %s, %s", id, err)
			}

			totalCount++
			paramGroups = append(paramGroups, *dbPg)
		}
	} else {
		for {
			req := conn.NewDescribeUDBParamGroupRequest()
			req.Limit = ucloud.Int(limit)
			req.Offset = ucloud.Int(offset)
			if val, ok := d.GetOk("availability_zone"); ok {
				req.Zone = ucloud.String(val.(string))
			}

			if val, ok := d.GetOk("region_flag"); ok {
				req.RegionFlag = ucloud.Bool(val.(bool))
			}

			if val, ok := d.GetOk("class_type"); ok {
				req.ClassType = ucloud.String(val.(string))
			}

			resp, err := conn.DescribeUDBParamGroup(req)
			if err != nil {
				return fmt.Errorf("error in read db param groups, %s", err)
			}

			if resp == nil || len(resp.DataSet) < 1 {
				break
			}

			fetched = append(fetched, resp.DataSet...)
			totalCount += len(resp.DataSet)

			if len(resp.DataSet) < limit {
				break
			}

			offset = offset + limit
		}

		engine, eOk := d.GetOk("engine")
		for _, item := range fetched {
			if eOk && !strings.HasPrefix(item.DBTypeId, engine.(string)) {
				continue
			}

			filtered = append(filtered, item)
		}

		engineVersion, evOk := d.GetOk("engine_version")
		for _, item := range filtered {
			if evOk && !strings.HasSuffix(item.DBTypeId, engineVersion.(string)) {
				continue
			}

			paramGroups = append(paramGroups, item)
			totalCount++
		}
	}

	d.Set("total_count", totalCount)

	err := dataSourceUCloudDBParamGroupsSave(d, paramGroups)
	if err != nil {
		return fmt.Errorf("error in read param groups, %s", err)
	}

	return nil
}

func dataSourceUCloudDBParamGroupsSave(d *schema.ResourceData, paramGroups []udb.UDBParamGroupSet) error {
	ids := []string{}
	data := []map[string]interface{}{}

	for _, paramGroup := range paramGroups {
		ids = append(ids, strconv.Itoa(paramGroup.GroupId))
		paramMember := []map[string]interface{}{}
		for _, item := range paramGroup.ParamMember {
			paramMember = append(paramMember, map[string]interface{}{
				"key":           item.Key,
				"value":         item.Value,
				"value_type":    item.ValueType,
				"allowed_value": item.AllowedVal,
				"apply_type":    item.ApplyType,
				"modifiable":    item.Modifiable,
				"format_type":   item.FormatType,
			})
		}

		data = append(data, map[string]interface{}{
			"id":           paramGroup.GroupId,
			"name":         paramGroup.GroupName,
			"db_type":      paramGroup.DBTypeId,
			"description":  paramGroup.Description,
			"modifiable":   paramGroup.Modifiable,
			"param_member": paramMember,
		})
	}

	d.SetId(hashStringArray(ids))
	if err := d.Set("param_groups", data); err != nil {
		return err
	}

	if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
		writeToFile(outputFile.(string), data)
	}

	return nil
}
