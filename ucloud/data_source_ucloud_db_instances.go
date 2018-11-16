package ucloud

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ucloud/ucloud-sdk-go/services/udb"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

func dataSourceUCloudDBInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUCloudDBInstancesRead,

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

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"total_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"db_instances": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"zone": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"backup_zone": &schema.Schema{
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

						"param_group_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"src_db_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"instance_charge_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"memory": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},

						"instance_storage": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},

						"role": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"disk_used_size": &schema.Schema{
							Type:     schema.TypeFloat,
							Computed: true,
						},

						"data_file_size": &schema.Schema{
							Type:     schema.TypeFloat,
							Computed: true,
						},

						"log_file_size": &schema.Schema{
							Type:     schema.TypeFloat,
							Computed: true,
						},

						"backup_date": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"instance_mode": &schema.Schema{
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

						"expire_time": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"modify_time": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"ip_set": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},

									"port": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},

						"slave_instances": &schema.Schema{
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

									"param_group_id": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},

									"src_db_id": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},

									"instance_charge_type": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},

									"memory": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},

									"instance_storage": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},

									"role": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},

									"disk_used_size": &schema.Schema{
										Type:     schema.TypeFloat,
										Computed: true,
									},

									"data_file_size": &schema.Schema{
										Type:     schema.TypeFloat,
										Computed: true,
									},

									"log_file_size": &schema.Schema{
										Type:     schema.TypeFloat,
										Computed: true,
									},

									"backup_date": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},

									"instance_mode": &schema.Schema{
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

									"expire_time": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},

									"modify_time": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},

									"ip_set": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": &schema.Schema{
													Type:     schema.TypeString,
													Computed: true,
												},

												"port": &schema.Schema{
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
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

func dataSourceUCloudDBInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.udbconn

	var dbInstances []udb.UDBInstanceSet
	var totalCount int
	limit := 100
	offset := 0
	if ids, ok := d.GetOk("ids"); ok && len(ids.([]interface{})) > 0 {
		for _, id := range ifaceToStringSlice(ids) {
			db, err := client.describeDBInstanceById(id)
			if err != nil {
				return fmt.Errorf("error in read db instance %s, %s", id, err)
			}

			totalCount++
			dbInstances = append(dbInstances, *db)
		}
	} else {
		for {
			req := conn.NewDescribeUDBInstanceRequest()
			req.Limit = ucloud.Int(limit)
			req.Offset = ucloud.Int(offset)
			if val, ok := d.GetOk("availability_zone"); ok {
				req.Zone = ucloud.String(val.(string))
			}

			if val, ok := d.GetOk("class_type"); ok {
				req.ClassType = ucloud.String(val.(string))
			}

			resp, err := conn.DescribeUDBInstance(req)
			if err != nil {
				return fmt.Errorf("error in read db instances, %s", err)
			}

			if resp == nil || len(resp.DataSet) < 1 {
				break
			}

			dbInstances = append(dbInstances, resp.DataSet...)
			totalCount += len(resp.DataSet)

			if len(resp.DataSet) < limit {
				break
			}

			offset = offset + limit
		}
	}

	d.Set("total_count", totalCount)

	err := dataSourceUCloudDBInstancesSave(d, dbInstances)
	if err != nil {
		return fmt.Errorf("error in read instances, %s", err)
	}

	return nil
}

func dataSourceUCloudDBInstancesSave(d *schema.ResourceData, dbInstances []udb.UDBInstanceSet) error {
	ids := []string{}
	data := []map[string]interface{}{}

	for _, dbInstance := range dbInstances {
		ids = append(ids, dbInstance.DBId)
		slaveInstances := []map[string]interface{}{}
		for _, item := range dbInstance.DataSet {
			ids = append(ids, item.DBId)
			ipSetSlave := []map[string]interface{}{}
			ipSetSlave = append(ipSetSlave, map[string]interface{}{
				"ip":   item.VirtualIP,
				"port": item.Port,
			})
			arr := strings.Split(item.DBTypeId, "-")
			slaveInstances = append(slaveInstances, map[string]interface{}{
				"id":                   item.DBId,
				"name":                 item.Name,
				"engine":               arr[0],
				"engine_version":       arr[1],
				"param_group_id":       strconv.Itoa(item.ParamGroupId),
				"src_db_id":            item.SrcDBId,
				"instance_charge_type": item.ChargeType,
				"memory":               item.MemoryLimit,
				"instance_storage":     item.DiskSpace,
				"role":                 item.Role,
				"disk_used_size":       item.DiskUsedSize,
				"data_file_size":       item.DataFileSize,
				"log_file_size":        item.LogFileSize,
				"backup_date":          item.BackupDate,
				"instance_mode":        item.InstanceMode,
				"status":               item.State,
				"create_time":          timestampToString(item.CreateTime),
				"expire_time":          timestampToString(item.ExpiredTime),
				"modify_time":          timestampToString(item.ModifyTime),
				"ip_set":               ipSetSlave,
			})
		}

		ipSetMaster := []map[string]interface{}{}
		ipSetMaster = append(ipSetMaster, map[string]interface{}{
			"ip":   dbInstance.VirtualIP,
			"port": dbInstance.Port,
		})

		brr := strings.Split(dbInstance.DBTypeId, "-")
		data = append(data, map[string]interface{}{
			"id":                   dbInstance.DBId,
			"name":                 dbInstance.Name,
			"engine":               brr[0],
			"engine_version":       brr[1],
			"param_group_id":       strconv.Itoa(dbInstance.ParamGroupId),
			"src_db_id":            dbInstance.SrcDBId,
			"instance_charge_type": dbInstance.ChargeType,
			"memory":               dbInstance.MemoryLimit,
			"instance_storage":     dbInstance.DiskSpace,
			"role":                 dbInstance.Role,
			"disk_used_size":       dbInstance.DiskUsedSize,
			"data_file_size":       dbInstance.DataFileSize,
			"log_file_size":        dbInstance.LogFileSize,
			"backup_date":          dbInstance.BackupDate,
			"instance_mode":        dbInstance.InstanceMode,
			"status":               dbInstance.State,
			"create_time":          timestampToString(dbInstance.CreateTime),
			"expire_time":          timestampToString(dbInstance.ExpiredTime),
			"modify_time":          timestampToString(dbInstance.ModifyTime),
			"ip_set":               ipSetMaster,
			"slave_instances":      slaveInstances,
		})
	}

	d.SetId(hashStringArray(ids))
	if err := d.Set("db_instances", data); err != nil {
		return err
	}

	if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
		writeToFile(outputFile.(string), data)
	}

	return nil
}
