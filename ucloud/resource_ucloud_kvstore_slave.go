package ucloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

func resourceUCloudKVStoreSlave() *schema.Resource {
	return &schema.Resource{
		Create: resourceUCloudKVStoreSlaveCreate,
		Read:   resourceUCloudKVStoreSlaveRead,
		Update: resourceUCloudKVStoreSlaveUpdate,
		Delete: resourceUCloudKVStoreSlaveDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateKVStoreInstanceName,
			},

			"instance_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateKVStoreInstanceType,
			},

			"master_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"password": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validateKVStoreInstancePassword,
			},

			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"parameter_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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

			"instance_charge_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"create_time": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"update_time": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"expire_time": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceUCloudKVStoreSlaveCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.umemconn

	req := conn.NewCreateURedisGroupRequest()
	req.Name = ucloud.String(d.Get("name").(string))
	req.Zone = ucloud.String(d.Get("availability_zone").(string))
	req.Size = ucloud.Int(getKVStoreCapability(d.Get("instance_type").(string)))
	req.HighAvailability = ucloud.String("enable")

	if v, ok := d.GetOk("parameter_group_id"); ok {
		req.ConfigId = ucloud.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		req.Password = ucloud.String(v.(string))
	}

	if v, ok := d.GetOk("tag"); ok {
		req.Tag = ucloud.String(v.(string))
	}

	if v, ok := d.GetOk("master_id"); ok {
		req.MasterGroupId = ucloud.String(v.(string))
	}

	resp, err := conn.CreateURedisGroup(req)
	if err != nil {
		return fmt.Errorf("error in create memory instance, %s", err)
	}

	d.SetId(resp.GroupId)

	if err := client.waitActiveStandbyRedisRunning(d.Id()); err != nil {
		return fmt.Errorf("error in create momory instance, %s", err)
	}
	return resourceUCloudKVStoreSlaveUpdate(d, meta)
}

func resourceUCloudKVStoreSlaveUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := updateActiveStandbyRedisInstanceWithoutRead(d, meta); err != nil {
		return err
	}
	return resourceUCloudKVStoreSlaveRead(d, meta)
}

func resourceUCloudKVStoreSlaveRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)

	inst, err := client.describeActiveStandbyRedisById(d.Id())
	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("do %s failed in read memory instance %s, %s", "DescribeURedisGroup", d.Id(), err)
	}

	d.Set("name", inst.Name)
	d.Set("tag", inst.Tag)
	d.Set("instance_charge_type", inst.ChargeType)
	d.Set("instance_type", fmt.Sprintf("redis-master-%v", inst.Size))
	d.Set("parameter_group_id", inst.ConfigId)
	d.Set("ip_set", []map[string]interface{}{{
		"ip":   inst.VirtualIP,
		"port": inst.Port,
	}})

	d.Set("create_time", timestampToString(inst.CreateTime))
	d.Set("update_time", timestampToString(inst.ModifyTime))
	d.Set("expire_time", timestampToString(inst.ExpireTime))
	d.Set("status", inst.State)
	return nil
}

func resourceUCloudKVStoreSlaveDelete(d *schema.ResourceData, meta interface{}) error {
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		return deleteActiveStandbyRedisInstance(d, meta)
	})
}
