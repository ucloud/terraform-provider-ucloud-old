package ucloud

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ucloud/ucloud-sdk-go/services/udb"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

func dataSourceUCloudDBParameterGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUCloudDBParameterGroupsRead,

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

			"multi_az": &schema.Schema{
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

			"parameter_groups": &schema.Schema{
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

						"engine": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"engine_version": &schema.Schema{
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

						"parameter_set": &schema.Schema{
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
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceUCloudDBParameterGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.udbconn
	var fetched []udb.UDBParamGroupSet
	var filtered []udb.UDBParamGroupSet
	var parameterGroups []udb.UDBParamGroupSet
	var totalCount int
	limit := 100
	offset := 0
	if ids, ok := d.GetOk("ids"); ok && len(ids.([]interface{})) > 0 {
		var zone string
		if val, ok := d.GetOk("availability_zone"); ok {
			zone = val.(string)
		} else {
			return fmt.Errorf("availability zone must be set when look up parameter groups  by ids")
		}
		for _, id := range ifaceToStringSlice(ids) {
			dbPg, err := client.describeDBParameterGroupByIdAndZone(id, zone)
			if err != nil {
				return fmt.Errorf("error in read db param group %s, %s", id, err)
			}

			totalCount++
			parameterGroups = append(parameterGroups, *dbPg)
		}
	} else {
		for {
			req := conn.NewDescribeUDBParamGroupRequest()
			req.Limit = ucloud.Int(limit)
			req.Offset = ucloud.Int(offset)
			if val, ok := d.GetOk("availability_zone"); ok {
				req.Zone = ucloud.String(val.(string))
			}

			if val, ok := d.GetOk("multi_az"); ok {
				req.RegionFlag = ucloud.Bool(val.(bool))
			}

			resp, err := conn.DescribeUDBParamGroup(req)
			if err != nil {
				return fmt.Errorf("error in read db parameter groups, %s", err)
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

			parameterGroups = append(parameterGroups, item)
			totalCount++
		}
	}

	d.Set("total_count", totalCount)

	err := dataSourceUCloudDBParameterGroupsSave(d, parameterGroups)
	if err != nil {
		return fmt.Errorf("error in read parameter groups, %s", err)
	}

	return nil
}

func dataSourceUCloudDBParameterGroupsSave(d *schema.ResourceData, parameterGroups []udb.UDBParamGroupSet) error {
	ids := []string{}
	data := []map[string]interface{}{}
	for _, parameterGroup := range parameterGroups {
		ids = append(ids, strconv.Itoa(parameterGroup.GroupId))
		parameterMember := []map[string]interface{}{}
		for _, item := range parameterGroup.ParamMember {
			parameterMember = append(parameterMember, map[string]interface{}{
				"key":           item.Key,
				"value":         item.Value,
				"value_type":    pgValueTypeCvt.mustConvert(item.ValueType),
				"allowed_value": item.AllowedVal,
			})
		}

		arr := strings.Split(parameterGroup.DBTypeId, "-")
		data = append(data, map[string]interface{}{
			"id":             strconv.Itoa(parameterGroup.GroupId),
			"name":           parameterGroup.GroupName,
			"engine":         arr[0],
			"engine_version": arr[1],
			"description":    parameterGroup.Description,
			"modifiable":     parameterGroup.Modifiable,
			"parameter_set":  parameterMember,
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
