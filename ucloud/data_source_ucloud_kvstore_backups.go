package ucloud

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ucloud/ucloud-sdk-go/services/umem"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

func dataSourceUCloudKVStoreBackups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUCloudKVStoreBackupsRead,
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

			"kvstore_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
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

			"backups": {
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

						"kvstore_instance_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"kvstore_instance_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"size": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"type": &schema.Schema{
							Type:     schema.TypeInt,
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
					},
				},
			},
		},
	}
}

func dataSourceUCloudKVStoreBackupsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*UCloudClient).umemconn

	req := conn.NewDescribeURedisBackupRequest()
	if v, ok := d.GetOk("kvstore_instance_id"); ok {
		req.GroupId = ucloud.String(v.(string))
	}

	var backups []umem.URedisBackupSet
	limit := 100
	offset := 0
	for {
		req.Limit = ucloud.Int(limit)
		req.Offset = ucloud.Int(offset)
		resp, err := conn.DescribeURedisBackup(req)
		if err != nil {
			return fmt.Errorf("error in read kvstore backup list, %s", err)
		}

		if resp == nil || len(resp.DataSet) < 1 {
			break
		}

		backups = append(backups, resp.DataSet...)
		if len(resp.DataSet) < limit {
			break
		}

		offset = offset + limit
	}

	if v, ok := d.GetOk("name_regex"); ok {
		r := regexp.MustCompile(v.(string))
		backups = filterBackups(backups, func(item *umem.URedisBackupSet) bool {
			return r.MatchString(item.BackupName)
		})
	}

	if v, ok := d.GetOk("ids"); ok {
		backups = filterBackups(backups, func(item *umem.URedisBackupSet) bool {
			err := checkStringIn(item.BackupId, ifaceToStringSlice(v))
			return err == nil
		})
	}

	d.Set("total_count", len(backups))
	err := dataSourceUCloudKVStoreBackupsSave(d, backups)
	if err != nil {
		return fmt.Errorf("error in read kvstore backup list, %s", err)
	}

	return nil
}

func filterBackups(backups []umem.URedisBackupSet, fn func(*umem.URedisBackupSet) bool) []umem.URedisBackupSet {
	var vL []umem.URedisBackupSet
	for _, v := range backups {
		if fn(&v) {
			vL = append(vL, v)
		}
	}
	return vL
}

func dataSourceUCloudKVStoreBackupsSave(d *schema.ResourceData, backups []umem.URedisBackupSet) error {
	ids := []string{}
	data := []map[string]interface{}{}

	for _, item := range backups {
		ids = append(ids, item.BackupId)
		data = append(data, map[string]interface{}{
			"availability_zone":     item.Zone,
			"id":                    item.BackupId,
			"name":                  item.BackupName,
			"kvstore_instance_id":   item.GroupId,
			"kvstore_instance_name": item.GroupName,
			"size":                  item.BackupSize,
			"type":                  item.BackupType,
			"status":                item.State,
			"create_time":           item.BackupTime,
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
