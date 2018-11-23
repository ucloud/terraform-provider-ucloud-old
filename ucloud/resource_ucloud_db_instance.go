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
				ValidateFunc: validateStringInChoices([]string{"5.5", "5.6", "5.7", "9.4", "9.6"}),
				ForceNew:     true,
				Required:     true,
			},

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDBInstanceName,
			},

			"instance_storage": &schema.Schema{
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateDataDiskSize(20, 3000),
			},

			"parameter_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"instance_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDBInstanceType,
			},

			"port": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(3306, 65535),
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

			"backup_begin_time": &schema.Schema{
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDBInstanceBlackList,
			},

			"tag": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceName,
				Computed:     true,
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

	engine := d.Get("engine").(string)
	// skip error because it has been validated by schema
	dbType, _ := parseDBInstanceType(d.Get("instance_type").(string))
	if dbType.Engine != engine {
		return fmt.Errorf("error in create db instance, engine of instance type %s must be same as engine %s", dbType.Engine, engine)
	}

	if dbType.Engine == "postgresql" && dbType.Type == "ha" {
		return fmt.Errorf("error in create db instance, high availability postgresql is not supported at this time")
	}

	req := conn.NewCreateUDBInstanceRequest()
	req.Name = ucloud.String(d.Get("name").(string))
	req.AdminPassword = ucloud.String(d.Get("password").(string))
	zone := d.Get("availability_zone").(string)
	req.Zone = ucloud.String(zone)
	req.DiskSpace = ucloud.Int(d.Get("instance_storage").(int))
	req.ChargeType = ucloud.String(d.Get("instance_charge_type").(string))
	req.Quantity = ucloud.Int(d.Get("instance_duration").(int))
	req.AdminUser = ucloud.String("root")
	req.InstanceType = ucloud.String("SATA_SSD")
	req.MemoryLimit = ucloud.Int(dbType.Memory * 1000)
	req.InstanceMode = ucloud.String(dbMap.convert(dbType.Type))
	engineVersion := d.Get("engine_version").(string)
	if engine == "mysql" || engine == "percona" {
		if err := checkStringIn(engineVersion, []string{"5.5", "5.6", "5.7"}); err != nil {
			return fmt.Errorf("The current engine version is not supported, %s", err)
		}
	} else {
		if err := checkStringIn(engineVersion, []string{"9.4", "9.6"}); err != nil {
			return fmt.Errorf("The current engine version is not supported, %s", err)
		}
	}
	req.DBTypeId = ucloud.String(strings.Join([]string{engine, engineVersion}, "-"))

	// if val, ok := d.GetOk("tag"); ok {
	// 	req.Tag = ucloud.String(val.(string))
	// }

	if val, ok := d.GetOk("port"); ok {
		req.Port = ucloud.Int(val.(int))
	} else {
		if engine == "mysql" || engine == "percona" {
			req.Port = ucloud.Int(3306)
		}
		if engine == "postgresql" {
			req.Port = ucloud.Int(5432)
		}
	}

	if val, ok := d.GetOk("backup_zone"); ok {
		if val.(string) != zone {
			req.BackupZone = ucloud.String(val.(string))
		}
	}

	if val, ok := d.GetOk("backup_count"); ok {
		req.BackupCount = ucloud.Int(val.(int))
	}

	if val, ok := d.GetOk("backup_begin_time"); ok {
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

	if val, ok := d.GetOk("vpc_id"); ok {
		req.VPCId = ucloud.String(val.(string))
	}

	if val, ok := d.GetOk("subnet_id"); ok {
		req.SubnetId = ucloud.String(val.(string))
	}

	pgId, err := strconv.Atoi(d.Get("parameter_group_id").(string))
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
	stateConf := client.dbWaitForState(d.Id(), "Running")

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
		stateConf := client.dbWaitForState(d.Id(), "Running")

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
		stateConf := client.dbWaitForState(d.Id(), "Running")

		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("wait for update db password failed in update db %s, %s", d.Id(), err)
		}
	}

	isChanged := false
	req := conn.NewResizeUDBInstanceRequest()
	req.DBId = ucloud.String(d.Id())

	if d.HasChange("instance_type") && !d.IsNewResource() {
		d.SetPartial("instance_type")
		engine := d.Get("engine").(string)
		old, new := d.GetChange("instance_type")

		oldType, _ := parseDBInstanceType(old.(string))

		newType, _ := parseDBInstanceType(new.(string))

		if newType.Engine != engine {
			return fmt.Errorf("error in update db instance, engine of instance type %s must be same as engine %s", newType.Engine, engine)
		}

		if newType.Engine == "postgresql" && newType.Type == "ha" {
			return fmt.Errorf("error in update db instance, high availability of postgresql is not supported at this time")
		}

		if oldType.Memory != newType.Memory {
			req.MemoryLimit = ucloud.Int(newType.Memory * 1000)
		}

		isChanged = true
	}

	if d.HasChange("instance_storage") && !d.IsNewResource() {
		d.SetPartial("instance_storage")
		req.DiskSpace = ucloud.Int(d.Get("instance_storage").(int))
		isChanged = true
	}

	if isChanged {
		db, err := client.describeDBInstanceById(d.Id())
		if err != nil {
			if isNotFoundError(err) {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("do %s failed in update db instance %s, %s", "DescribeUDBInstance", d.Id(), err)
		}

		if db.InstanceMode == "Normal" {
			//update these attributes of db instance, we need to wait it stopped
			stopReq := conn.NewStopUDBInstanceRequest()
			stopReq.DBId = ucloud.String(d.Id())
			stopReq.Zone = ucloud.String(d.Get("availability_zone").(string))
			if db.State != "Shutoff" {
				_, err := conn.StopUDBInstance(stopReq)

				if err != nil {
					return fmt.Errorf("do %s failed in update db instance %s, %s", "StopUDBInstance", d.Id(), err)
				}

				// after stop db instance, we need to wait it stopped
				stateConf := client.dbWaitForState(d.Id(), "Shutoff")

				if _, err := stateConf.WaitForState(); err != nil {
					return fmt.Errorf("wait for stop db instance failed in update db instance %s, %s", d.Id(), err)
				}
			}

			if _, err := conn.ResizeUDBInstance(req); err != nil {
				return fmt.Errorf("do %s failed in update db instance %s, %s", "ResizeUDBInstance", d.Id(), err)
			}

			// after resize db instance, we need to wait it completed
			stateConf := client.dbWaitForState(d.Id(), "Shutoff")

			if _, err := stateConf.WaitForState(); err != nil {
				return fmt.Errorf("wait for resize db instance failed in update db %s, %s", d.Id(), err)
			}

			// after update these attributes of db instance, we need to wait it started
			startReq := conn.NewStartUDBInstanceRequest()
			startReq.DBId = ucloud.String(d.Id())
			startReq.Zone = ucloud.String(d.Get("availability_zone").(string))

			_, err = conn.StartUDBInstance(startReq)

			if err != nil {
				return fmt.Errorf("do %s failed in update db instance %s, %s", "StartUDBInstance", d.Id(), err)
			}

			// after db instance update, we need to wait it started
			stateConf = client.dbWaitForState(d.Id(), "Running")

			if _, err := stateConf.WaitForState(); err != nil {
				return fmt.Errorf("wait for start db instance failed in update db instance %s, %s", d.Id(), err)
			}
		}

		if _, err := conn.ResizeUDBInstance(req); err != nil {
			return fmt.Errorf("do %s failed in update db %s, %s", "ResizeUDBInstance", d.Id(), err)
		}

		// after resize db instance, we need to wait it completed
		stateConf := client.dbWaitForState(d.Id(), "Running")

		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("wait for resize db instance failed in update db %s, %s", d.Id(), err)
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

	if d.HasChange("backup_begin_time") && !d.IsNewResource() {
		d.SetPartial("backup_begin_time")
		buReq.BackupTime = ucloud.Int(d.Get("backup_begin_time").(int))
		backupChanged = true
	}

	if backupChanged {
		if _, err := conn.UpdateUDBInstanceBackupStrategy(buReq); err != nil {
			return fmt.Errorf("do %s failed in update db %s, %s", "UpdateUDBInstanceBackupStrategy", d.Id(), err)
		}

		// after update db backup strategy, we need to wait it completed
		stateConf := client.dbWaitForState(d.Id(), "Running")

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
	d.Set("parameter_group_id", strconv.Itoa(db.ParamGroupId))
	d.Set("port", db.Port)
	d.Set("status", db.State)
	d.Set("instance_charge_type", db.ChargeType)
	d.Set("instance_storage", db.DiskSpace)
	d.Set("backup_zone", db.BackupZone)
	d.Set("availability_zone", db.Zone)
	d.Set("backup_count", db.BackupCount)
	d.Set("backup_duration", db.BackupDuration)
	d.Set("backup_begin_time", db.BackupBeginTime)
	d.Set("backup_date", db.BackupDate)
	d.Set("black_list", db.BackupBlacklist)
	d.Set("backup_zone", db.BackupZone)
	// d.Set("vpc_id", db.VPCId)
	// d.Set("subnet_id", db.SubnetId)

	//d.Set("tag", db.Tag)
	d.Set("create_time", timestampToString(db.CreateTime))
	d.Set("expire_time", timestampToString(db.ExpiredTime))
	d.Set("modify_time", timestampToString(db.ModifyTime))
	var dbType dbInstanceType
	dbType.Memory = db.MemoryLimit / 1000
	dbType.Engine = arr[0]
	dbType.Type = dbMap.unconvert(db.InstanceMode)
	d.Set("instance_type", fmt.Sprintf("%s-%s-%d", dbType.Engine, dbType.Type, dbType.Memory))

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
			stateConf := client.dbWaitForState(d.Id(), "Shutoff")

			if _, err := stateConf.WaitForState(); err != nil {
				return resource.RetryableError(fmt.Errorf("wait for db instance stop failed in delete db instance %s, %s", d.Id(), err))
			}
		}

		if _, err := conn.DeleteUDBInstance(req); err != nil {
			return resource.NonRetryableError(fmt.Errorf("error in delete db instance %s, %s", d.Id(), err))
		}

		if _, err := client.describeDBInstanceById(d.Id()); err != nil {
			if isNotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("do %s failed in delete db instance %s, %s", "DescribeUDBInstance", d.Id(), err))
		}

		return resource.RetryableError(fmt.Errorf("delete db instance but it still exists"))
	})
}
