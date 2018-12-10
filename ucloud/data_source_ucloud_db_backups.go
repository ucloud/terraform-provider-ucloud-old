package ucloud

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/ucloud/ucloud-sdk-go/services/udb"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

func dataSourceUCloudDBBackups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUCloudDBBackupsRead,

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

			"db_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"backup_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"automatic", "manual"}, false),
			},

			"begin_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRFC3339TimeString,
			},

			"end_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRFC3339TimeString,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"total_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"backups": &schema.Schema{
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

						"standby_zone": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"size": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"db_instance_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"db_instance_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"begin_time": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"end_time": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceUCloudDBBackupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.udbconn
	var backups []udb.UDBBackupSet
	var totalCount int
	limit := 100
	offset := 0
	if ids, ok := d.GetOk("ids"); ok && len(ids.([]interface{})) > 0 {
		var zone string
		if val, ok := d.GetOk("availability_zone"); ok {
			zone = val.(string)
		} else {
			return fmt.Errorf("availability zone must be set when look up backups by ids")
		}
		for _, id := range ifaceToStringSlice(ids) {
			backup, err := client.describeDBBackupByIdAndZone(id, zone)
			if err != nil {
				return fmt.Errorf("error in read db backup %s, %s", id, err)
			}

			totalCount++
			backups = append(backups, *backup)
		}
	} else if dbId, ok := d.GetOk("db_instance_id"); ok {
		for {
			req := conn.NewDescribeUDBBackupRequest()
			req.Limit = ucloud.Int(limit)
			req.Offset = ucloud.Int(offset)
			req.DBId = ucloud.String(dbId.(string))
			resp, err := conn.DescribeUDBBackup(req)
			if err != nil {
				return fmt.Errorf("error in read db backups, %s", err)
			}

			if resp == nil || len(resp.DataSet) < 1 {
				break
			}

			backups = append(backups, resp.DataSet...)
			totalCount += len(resp.DataSet)

			if len(resp.DataSet) < limit {
				break
			}

			offset = offset + limit
		}
	} else {
		for {
			req := conn.NewDescribeUDBBackupRequest()
			req.Limit = ucloud.Int(limit)
			req.Offset = ucloud.Int(offset)
			if val, ok := d.GetOk("availability_zone"); ok {
				req.Zone = ucloud.String(val.(string))
			}

			if val, ok := d.GetOk("end_time"); ok {
				// skip error because it has been validated by schema
				endTime, _ := stringToTimestamp(val.(string))
				req.EndTime = ucloud.Int(endTime)
			}

			if val, ok := d.GetOk("begin_time"); ok {
				// skip error because it has been validated by schema
				beginTime, _ := stringToTimestamp(val.(string))
				req.BeginTime = ucloud.Int(beginTime)
			}

			if val, ok := d.GetOk("backup_type"); ok {
				req.BackupType = ucloud.Int(backupTypeCvt.mustUnconvert(val.(string)))
			}

			resp, err := conn.DescribeUDBBackup(req)
			if err != nil {
				return fmt.Errorf("error in read db backups, %s", err)
			}

			if resp == nil || len(resp.DataSet) < 1 {
				break
			}

			backups = append(backups, resp.DataSet...)
			totalCount += len(resp.DataSet)

			if len(resp.DataSet) < limit {
				break
			}

			offset = offset + limit
		}
	}

	d.Set("total_count", totalCount)

	err := dataSourceUCloudDBBackupsSave(d, backups)
	if err != nil {
		return fmt.Errorf("error in read param groups, %s", err)
	}

	return nil
}

func dataSourceUCloudDBBackupsSave(d *schema.ResourceData, backups []udb.UDBBackupSet) error {
	ids := []string{}
	data := []map[string]interface{}{}
	valueType := make(map[int]string)
	valueType[0] = "unknown"
	valueType[10] = "int"
	valueType[20] = "string"
	valueType[30] = "bool"
	for _, backup := range backups {
		ids = append(ids, strconv.Itoa(backup.BackupId))
		data = append(data, map[string]interface{}{
			"id":               backup.BackupId,
			"name":             backup.BackupName,
			"size":             backup.BackupSize,
			"type":             backup.BackupType,
			"status":           backup.State,
			"db_instance_id":   backup.DBId,
			"db_instance_name": backup.DBName,
			"zone":             backup.Zone,
			"standby_zone":     backup.BackupZone,
			"begin_time":       timestampToString(backup.BackupTime),
			"end_time":         timestampToString(backup.BackupEndTime),
		})
	}

	d.SetId(hashStringArray(ids))
	if err := d.Set("backups", data); err != nil {
		return err
	}

	if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
		writeToFile(outputFile.(string), data)
	}

	return nil
}
