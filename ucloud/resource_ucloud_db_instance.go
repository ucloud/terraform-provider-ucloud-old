package ucloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

func resourceUCloudDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceUCloudDBInstanceCreate,
		Read:   resourceUCloudDBInstanceRead,
		Update: resourceUCloudDBInstanceUpdate,
		Delete: resourceUCloudDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"password": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validateInstancePassword,
			},

			"engine": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateStringInChoices([]string{"mysql", "percona", "postgresql"}),
				ForceNew:     true,
				Required:     true,
			},

			"engine_version": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateStringInChoices([]string{"5.1", "5.5", "5.6", "5.7", "9.4", "9. 6", "10.4"}),
				ForceNew:     true,
				Required:     true,
			},

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateInstanceName,
			},

			"port": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(1, 65535),
			},

			"instance_storage": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateDataDiskSize(20, 3000),
			},

			"param_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"memory_limit": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntInChoices([]int{1000, 2000, 4000, 6000, 8000, 12000, 16000, 24000, 32000, 48000, 64000, 96000}),
			},

			"instance_charge_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Month",
				ValidateFunc: validateStringInChoices([]string{"Year", "Month", "Dynamic"}),
			},

			"instance_duration": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"username": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "root",
				ValidateFunc: validateInstanceName,
			},

			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
		},
	}
}

func resourceUCloudDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.udbconn

	req := conn.NewCreateUDBInstanceRequest()
	req.AdminPassword = ucloud.String(d.Get("password").(string))
	req.Zone = ucloud.String(d.Get("availability_zone").(string))
	engine := d.Get("engine").(string)
	engineVersion := d.Get("engine_version").(string)
	req.DBTypeId = ucloud.String(strings.Join([]string{engine, engineVersion}, "-"))
	req.Port = ucloud.Int(d.Get("port").(int))
	req.DiskSpace = ucloud.Int(d.Get("instance_storage").(int))
	req.ParamGroupId = ucloud.String(d.Get("param_group_id").(string))
	req.MemoryLimit = ucloud.Int(d.Get("memory_limit").(int))
	req.ChargeType = ucloud.String(d.Get("instance_charge_type").(string))
	req.Quantity = ucloud.Int(d.Get("instance_duration").(int))
	req.AdminUser = ucloud.String(d.Get("username").(string))

	if val, ok := d.GetOk("vpc_id"); ok {
		req.VPCId = ucloud.String(val.(string))
	}

	if val, ok := d.GetOk("subnet_id"); ok {
		req.SubnetId = ucloud.String(val.(string))
	}

	resp, err := conn.CreateUDBInstance(req)
	if err != nil {
		return fmt.Errorf("error in create db instance, %s", err)
	}

	d.SetId(resp.DBId)

	// after create db, we need to wait it initialized
	stateConf := dbWaitForState(client, d.Id(), "Running")

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("wait for db initialize failed in create db %s, %s", d.Id(), err)
	}

	return resourceUCloudDBInstanceUpdate(d, meta)
}

func resourceUCloudDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.udbconn

	d.Partial(true)

	if d.HasChange("name") && !d.IsNewResource() {
		d.SetPartial("name")
		req := conn.NewModifyUDBInstanceNameRequest()
		req.DBId = ucloud.String(d.Id())
		req.Name = ucloud.String(d.Get("name").(string))

		if _, err := conn.ModifyUDBInstanceName(req); err != nil {
			return fmt.Errorf("do %s failed in update db %s, %s", "ModifyUDBInstanceName", d.Id(), err)
		}

		// after update db name, we need to wait it completed
		stateConf := dbWaitForState(client, d.Id(), "Running")

		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("wait for update db name failed in update db %s, %s", d.Id(), err)
		}
	}

	if d.HasChange("password") && !d.IsNewResource() {
		d.SetPartial("password")
		req := conn.NewModifyUDBInstancePasswordRequest()
		req.DBId = ucloud.String(d.Id())
		req.Password = ucloud.String(d.Get("password").(string))

		if _, err := conn.ModifyUDBInstancePassword(req); err != nil {
			return fmt.Errorf("do %s failed in update db %s, %s", "ModifyUDBInstancePassword", d.Id(), err)
		}

		// after update db password, we need to wait it completed
		stateConf := dbWaitForState(client, d.Id(), "Running")

		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("wait for update db password failed in update db %s, %s", d.Id(), err)
		}
	}

	isChanged := false
	req := conn.NewResizeUDBInstanceRequest()
	req.DBId = ucloud.String(d.Id())

	if d.HasChange("memory_limit") && !d.IsNewResource() {
		d.SetPartial("memory_limit")
		req.MemoryLimit = ucloud.Int(d.Get("memory_limit").(int))
		isChanged = true
	}

	if d.HasChange("instance_storage") && !d.IsNewResource() {
		d.SetPartial("instance_storage")
		req.DiskSpace = ucloud.Int(d.Get("instance_storage").(int))
		isChanged = true
	}

	if isChanged {

		if _, err := conn.ResizeUDBInstance(req); err != nil {
			return fmt.Errorf("do %s failed in update db %s, %s", "ResizeUDBInstance", d.Id(), err)
		}

		// after resize db instance, we need to wait it completed
		stateConf := dbWaitForState(client, d.Id(), "Running")

		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("wait for resize db instance failed in update db %s, %s", d.Id(), err)
		}
	}

	d.Partial(false)

	return resourceUCloudDBInstanceRead(d, meta)
}

func resourceUCloudDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)

	db, err := client.describeDbInstanceById(d.Id())
	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("do %s failed in read db %s, %s", "DescribeUDBInstance", d.Id(), err)
	}

	arr := strings.Split(db.DBTypeId, "-")
	d.Set("name", db.Name)
	d.Set("engine", arr[0])
	d.Set("engine_version", arr[1])
	d.Set("param_group_id", db.ParamGroupId)
	d.Set("port", db.Port)
	d.Set("status", db.State)
	d.Set("create_time", timestampToString(db.CreateTime))
	d.Set("expire_time", timestampToString(db.ExpiredTime))
	d.Set("modify_time", timestampToString(db.ModifyTime))

	return nil
}

func resourceUCloudDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.udbconn

	req := conn.NewDeleteUDBInstanceRequest()
	req.DBId = ucloud.String(d.Id())
	stopReq := conn.NewStopUDBInstanceRequest()
	stopReq.DBId = ucloud.String(d.Id())

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		db, err := client.describeDbInstanceById(d.Id())
		if err != nil {
			if isNotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if db.State != "Shutoff" {
			if _, err := conn.StopUDBInstance(stopReq); err != nil {
				return resource.RetryableError(fmt.Errorf("do %s failed in delete db instance %s, %s", "StopUDBInstance", d.Id(), err))
			}

			// after instance stop, we need to wait it stoped
			stateConf := dbWaitForState(client, d.Id(), "Shutoff")

			if _, err := stateConf.WaitForState(); err != nil {
				return resource.RetryableError(fmt.Errorf("wait for db instance stop failed in delete db %s, %s", d.Id(), err))
			}
		}

		if _, err := conn.DeleteUDBInstance(req); err != nil {
			return resource.NonRetryableError(fmt.Errorf("error in delete db %s, %s", d.Id(), err))
		}

		if _, err := client.describeDbInstanceById(d.Id()); err != nil {
			if isNotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("do %s failed in delete db %s, %s", "DescribeUDBInstance", d.Id(), err))
		}

		return resource.RetryableError(fmt.Errorf("delete db but it still exists"))
	})
}

func dbWaitForState(client *UCloudClient, dbId, target string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{target},
		Timeout:    5 * time.Minute,
		Delay:      2 * time.Second,
		MinTimeout: 1 * time.Second,
		Refresh: func() (interface{}, string, error) {
			db, err := client.describeDbInstanceById(dbId)
			if err != nil {
				if isNotFoundError(err) {
					return nil, "pending", nil
				}
				return nil, "", err
			}

			state := db.State
			if state != target {
				state = "pending"
			}

			return db, state, nil
		},
	}
}
