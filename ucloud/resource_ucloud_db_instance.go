package ucloud

import (
	"fmt"
	"strconv"
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

			"backup_zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"master_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"is_lock": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},

			"is_force": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"password": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
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
				ValidateFunc: validateStringInChoices([]string{"5.1", "5.5", "5.6", "5.7", "9.4", "9.6", "10.4"}),
				ForceNew:     true,
				Required:     true,
			},

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDBInstanceName,
			},

			"port": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(3306, 65535),
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

			"memory": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntInChoices([]int{1, 2, 4, 6, 8, 12, 16, 24, 32, 48, 64, 96, 128}),
			},

			"instance_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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

			"backup_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  7,
				ForceNew: true,
			},

			"backup_duration": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  24,
				ForceNew: true,
			},

			"backup_time": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"backup_date": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"backup_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"backup_black_list": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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

			"role": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceUCloudDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.udbconn

	masterId, slaveOk := d.GetOk("master_id")
	if !slaveOk {
		req := conn.NewCreateUDBInstanceRequest()
		req.InstanceMode = ucloud.String("HA")
		req.Name = ucloud.String(d.Get("name").(string))
		req.AdminPassword = ucloud.String(d.Get("password").(string))
		req.Zone = ucloud.String(d.Get("availability_zone").(string))
		engine := d.Get("engine").(string)
		engineVersion := d.Get("engine_version").(string)
		req.DBTypeId = ucloud.String(strings.Join([]string{engine, engineVersion}, "-"))
		req.DiskSpace = ucloud.Int(d.Get("instance_storage").(int))
		req.MemoryLimit = ucloud.Int(d.Get("memory").(int) * 1000)
		req.ChargeType = ucloud.String(d.Get("instance_charge_type").(string))
		req.Quantity = ucloud.Int(d.Get("instance_duration").(int))
		req.AdminUser = ucloud.String("root")

		if val, ok := d.GetOk("port"); ok {
			req.Port = ucloud.Int(val.(int))
		} else {
			if engine == "mysql" {
				req.Port = ucloud.Int(3306)
			}
			if engine == "postgresql" {
				req.Port = ucloud.Int(5432)
			}
		}

		if val, ok := d.GetOk("backup_count"); ok {
			req.BackupCount = ucloud.Int(val.(int))
		}

		if val, ok := d.GetOk("backup_time"); ok {
			req.BackupTime = ucloud.Int(val.(int))
		}

		if val, ok := d.GetOk("backup_duration"); ok {
			req.BackupDuration = ucloud.Int(val.(int))
		}

		if val, ok := d.GetOk("backup_id"); ok {
			backupId, err := strconv.Atoi(val.(string))
			if err != nil {
				return err
			}
			req.BackupId = ucloud.Int(backupId)
		}

		if val, ok := d.GetOk("instance_type"); ok {
			req.InstanceType = ucloud.String(val.(string))
		}

		if val, ok := d.GetOk("vpc_id"); ok {
			req.VPCId = ucloud.String(val.(string))
		}

		if val, ok := d.GetOk("subnet_id"); ok {
			req.SubnetId = ucloud.String(val.(string))
		}

		pgId, err := strconv.Atoi(d.Get("param_group_id").(string))
		if err != nil {
			return err
		}
		req.ParamGroupId = ucloud.Int(pgId)

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
	} else {
		req := conn.NewCreateUDBSlaveRequest()

		req.InstanceMode = ucloud.String("HA")
		req.SrcId = ucloud.String(masterId.(string))
		req.Name = ucloud.String(d.Get("name").(string))
		req.Port = ucloud.Int(d.Get("port").(int))
		req.DiskSpace = ucloud.Int(d.Get("instance_storage").(int))
		req.MemoryLimit = ucloud.Int(d.Get("memory").(int) * 1000)

		if val, ok := d.GetOk("instance_type"); ok {
			req.InstanceType = ucloud.String(val.(string))
		}

		if val, ok := d.GetOk("is_lock"); ok {
			req.IsLock = ucloud.Bool(val.(bool))
		}

		resp, err := conn.CreateUDBSlave(req)

		if err != nil {
			return fmt.Errorf("error in create slave db, %s", err)
		}

		d.SetId(resp.DBId)

		// after create db, we need to wait it initialized
		stateConf := dbWaitForState(client, d.Id(), "Running")

		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("wait for slave db initialize failed in create slave db %s, %s", d.Id(), err)
		}

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

	if d.HasChange("memory") && !d.IsNewResource() {
		d.SetPartial("memory")
		req.MemoryLimit = ucloud.Int(d.Get("memory").(int) * 1000)
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

	if d.HasChange("master_id") && !d.IsNewResource() {
		d.SetPartial("master_id")
		old, new := d.GetChange("master_id")
		if old.(string) == "" {
			return fmt.Errorf("the master db cannot be reduced to the slave db")
		}

		if new.(string) != "" {
			return fmt.Errorf("the master id can only be updated to %s, got %s", "", new.(string))
		}

		req := conn.NewPromoteUDBSlaveRequest()
		req.DBId = ucloud.String(d.Id())
		req.IsForce = ucloud.Bool(d.Get("is_force").(bool))

		if _, err := conn.PromoteUDBSlave(req); err != nil {
			return fmt.Errorf("do %s failed in update db %s, %s", "PromoteUDBSlave", d.Id(), err)
		}

		// after promote slave db to master db, we need to wait it completed
		stateConf := dbWaitForState(client, d.Id(), "Running")

		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("wait for promote slave db failed in update db %s, %s", d.Id(), err)
		}
	}

	backupChanged := false
	buReq := conn.NewUpdateUDBInstanceBackupStrategyRequest()
	buReq.DBId = ucloud.String(d.Id())

	if d.HasChange("backup_date") {
		d.SetPartial("backup_count")
		buReq.BackupDate = ucloud.String(d.Get("backup_count").(string))
		backupChanged = true
	}

	if d.HasChange("backup_time") && !d.IsNewResource() {
		d.SetPartial("backup_time")
		buReq.BackupTime = ucloud.Int(d.Get("backup_time").(int))
		backupChanged = true
	}

	if backupChanged {
		if _, err := conn.UpdateUDBInstanceBackupStrategy(buReq); err != nil {
			return fmt.Errorf("do %s failed in update db %s, %s", "UpdateUDBInstanceBackupStrategy", d.Id(), err)
		}

		// after update db backup strategy, we need to wait it completed
		stateConf := dbWaitForState(client, d.Id(), "Running")

		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("wait for update db backup strategy failed in update db %s, %s", d.Id(), err)
		}
	}

	if d.HasChange("backup_black_list") {
		d.SetPartial("backup_black_list")
		req := conn.NewEditUDBBackupBlacklistRequest()
		req.Blacklist = ucloud.String(d.Get("backup_black_list").(string))
		req.DBId = ucloud.String(d.Id())

		if _, err := conn.EditUDBBackupBlacklist(req); err != nil {
			return fmt.Errorf("do %s failed in update db %s, %s", "EditUDBBackupBlacklist", d.Id(), err)
		}
	}

	d.Partial(false)

	return resourceUCloudDBInstanceRead(d, meta)
}

func resourceUCloudDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)

	db, err := client.describeDBInstanceById(d.Id())
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
	d.Set("param_group_id", strconv.Itoa(db.ParamGroupId))
	d.Set("port", db.Port)
	d.Set("status", db.State)
	d.Set("instance_charge_type", db.ChargeType)
	d.Set("memory", db.MemoryLimit/1000)
	d.Set("instance_storage", db.DiskSpace)
	d.Set("role", db.Role)
	d.Set("backup_zone", db.BackupZone)
	d.Set("availability_zone", db.Zone)
	d.Set("instance_charge_type", db.ChargeType)

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
		db, err := client.describeDBInstanceById(d.Id())
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

		if _, err := client.describeDBInstanceById(d.Id()); err != nil {
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
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
		Refresh: func() (interface{}, string, error) {
			db, err := client.describeDBInstanceById(dbId)
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
