package ucloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
)

func resourceUCloudKVStoreInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceUCloudKVStoreInstanceCreate,
		Read:   resourceUCloudKVStoreInstanceRead,
		Update: resourceUCloudKVStoreInstanceUpdate,
		Delete: resourceUCloudKVStoreInstanceDelete,
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

			"engine": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringInChoices([]string{"memcache", "redis"}),
			},

			"instance_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateKVStoreInstanceType,
			},

			"engine_version": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateStringInChoices([]string{"3.0", "3.2", "4.0"}),
			},

			"instance_charge_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
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

			"password": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validateKVStoreInstancePassword,
			},

			"parameter_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"backup_begin_time": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(0, 23),
			},

			"backup_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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

func resourceUCloudKVStoreInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	// skip error, because it has been validated at schema
	t, _ := parseKVStoreInstanceType(d.Get("instance_type").(string))

	if engine := d.Get("engine").(string); t.Engine != engine {
		return fmt.Errorf("error in create kvstore instance, engine of instance type %s must be same as engine %s", t.Engine, engine)
	}

	if t.Engine == "redis" && t.Type == "master" {
		return createActiveStandbyRedisInstance(d, meta)
	}

	if t.Engine == "redis" && t.Type == "distributed" {
		return createDistributedRedisInstance(d, meta)
	}

	if t.Engine == "memcache" && t.Type == "master" {
		return createActiveStandbyMemcacheInstance(d, meta)
	}

	if t.Engine == "memcache" && t.Type == "distributed" {
		return fmt.Errorf("error in create kvstore instance, distributed memcache is not supported")
	}

	return fmt.Errorf("error in create kvstore instance, %s is not supported", t.Engine)
}

func resourceUCloudKVStoreInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	// skip error, because it has been validated at schema
	t, _ := parseKVStoreInstanceType(d.Get("instance_type").(string))

	if engine := d.Get("engine").(string); t.Engine != engine {
		return fmt.Errorf("error in update kvstore instance, engine of instance type %s must be same as engine %s", t.Engine, engine)
	}

	if t.Engine == "redis" && t.Type == "master" {
		return updateActiveStandbyRedisInstance(d, meta)
	}

	if t.Engine == "redis" && t.Type == "distributed" {
		return updateDistributedRedisInstance(d, meta)
	}

	if t.Engine == "memcache" && t.Type == "master" {
		return updateActiveStandbyMemcacheInstance(d, meta)
	}

	return fmt.Errorf("error in update kvstore instance, current engine is not supported")
}

func resourceUCloudKVStoreInstanceRead(d *schema.ResourceData, meta interface{}) error {
	t, err := parseKVStoreInstanceType(d.Get("instance_type").(string))
	if err != nil {
		return fmt.Errorf("error in read kvstore instance, %s", err)
	}

	if t.Engine == "redis" && t.Type == "master" {
		return readActiveStandbyRedisInstance(d, meta)
	}

	if t.Engine == "redis" && t.Type == "distributed" {
		return readDistributedRedisInstance(d, meta)
	}

	if t.Engine == "memcache" && t.Type == "master" {
		return readActiveStandbyMemcacheInstance(d, meta)
	}

	return fmt.Errorf("error in read kvstore instance, current engine is not supported")
}

func resourceUCloudKVStoreInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	// skip error, because it has been validated at schema
	t, _ := parseKVStoreInstanceType(d.Get("instance_type").(string))

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if t.Engine == "redis" && t.Type == "master" {
			return deleteActiveStandbyRedisInstance(d, meta)
		}

		if t.Engine == "redis" && t.Type == "distributed" {
			return deleteDistributedRedisInstance(d, meta)
		}

		if t.Engine == "memcache" && t.Type == "master" {
			return deleteActiveStandbyMemcacheInstance(d, meta)
		}

		return resource.NonRetryableError(fmt.Errorf("error in delete kvstore instance, current engine is not supported"))
	})
}

func createActiveStandbyRedisInstance(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.umemconn

	req := conn.NewCreateURedisGroupRequest()
	req.Name = ucloud.String(d.Get("name").(string))
	req.Zone = ucloud.String(d.Get("availability_zone").(string))
	req.Size = ucloud.Int(getKVStoreCapability(d.Get("instance_type").(string)))
	req.Quantity = ucloud.Int(d.Get("instance_duration").(int))
	req.ChargeType = ucloud.String(d.Get("instance_charge_type").(string))
	req.HighAvailability = ucloud.String("enable")

	if v, ok := d.GetOk("engine_version"); ok {
		req.Version = ucloud.String(v.(string))
	} else {
		return fmt.Errorf("error in create kvstore instance, engine version is required")
	}

	if v, ok := d.GetOk("password"); ok {
		req.Password = ucloud.String(v.(string))
	}

	if v, ok := d.GetOk("tag"); ok {
		req.Tag = ucloud.String(v.(string))
	}

	if v, ok := d.GetOk("parameter_group_id"); ok {
		req.ConfigId = ucloud.String(v.(string))
	}

	if v, ok := d.GetOk("backup_begin_time"); ok {
		req.AutoBackup = ucloud.String("enable")
		req.BackupTime = ucloud.Int(v.(int))
	}

	if v, ok := d.GetOk("backup_id"); ok {
		req.BackupId = ucloud.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		req.VPCId = ucloud.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		req.SubnetId = ucloud.String(v.(string))
	}

	resp, err := conn.CreateURedisGroup(req)
	if err != nil {
		return fmt.Errorf("error in create kvstore instance, %s", err)
	}

	d.SetId(resp.GroupId)

	if err := client.waitActiveStandbyRedisRunning(d.Id()); err != nil {
		return fmt.Errorf("error in create momory instance, %s", err)
	}

	return updateActiveStandbyRedisInstance(d, meta)
}

func createDistributedRedisInstance(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.umemconn

	req := conn.NewCreateUMemSpaceRequest()
	req.Name = ucloud.String(d.Get("name").(string))
	req.Zone = ucloud.String(d.Get("availability_zone").(string))
	req.Size = ucloud.Int(getKVStoreCapability(d.Get("instance_type").(string)))
	req.Quantity = ucloud.Int(d.Get("instance_duration").(int))
	req.ChargeType = ucloud.String(d.Get("instance_charge_type").(string))
	req.Protocol = ucloud.String("redis")

	if v, ok := d.GetOk("vpc_id"); ok {
		req.VPCId = ucloud.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		req.SubnetId = ucloud.String(v.(string))
	}

	if v, ok := d.GetOk("tag"); ok {
		req.Tag = ucloud.String(v.(string))
	}

	resp, err := conn.CreateUMemSpace(req)
	if err != nil {
		return fmt.Errorf("error in create kvstore instance, %s", err)
	}

	d.SetId(resp.SpaceId)

	if err := client.waitDistributedRedisRunning(d.Id()); err != nil {
		return fmt.Errorf("error in create momory instance, %s", err)
	}

	return updateDistributedRedisInstance(d, meta)
}

func createActiveStandbyMemcacheInstance(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.umemconn

	req := conn.NewCreateUMemcacheGroupRequest()
	req.Name = ucloud.String(d.Get("name").(string))
	req.Zone = ucloud.String(d.Get("availability_zone").(string))
	req.Size = ucloud.Int(getKVStoreCapability(d.Get("instance_type").(string)))
	req.Quantity = ucloud.Int(d.Get("instance_duration").(int))
	req.ChargeType = ucloud.String(d.Get("instance_charge_type").(string))
	req.Protocol = ucloud.String("memcache")

	if v, ok := d.GetOk("vpc_id"); ok {
		req.VPCId = ucloud.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		req.SubnetId = ucloud.String(v.(string))
	}

	if v, ok := d.GetOk("tag"); ok {
		req.Tag = ucloud.String(v.(string))
	}

	resp, err := conn.CreateUMemcacheGroup(req)
	if err != nil {
		return fmt.Errorf("error in create kvstore instance, %s", err)
	}

	d.SetId(resp.GroupId)

	if err := client.waitActiveStandbyMemcacheRunning(d.Id()); err != nil {
		return fmt.Errorf("error in create momory instance, %s", err)
	}

	return updateActiveStandbyMemcacheInstance(d, meta)
}

func updateActiveStandbyRedisInstance(d *schema.ResourceData, meta interface{}) error {
	if err := updateActiveStandbyRedisInstanceWithoutRead(d, meta); err != nil {
		return err
	}
	return readActiveStandbyRedisInstance(d, meta)
}

func updateActiveStandbyRedisInstanceWithoutRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.umemconn

	d.Partial(true)

	if d.HasChange("name") && !d.IsNewResource() {
		req := conn.NewModifyURedisGroupNameRequest()
		req.GroupId = ucloud.String(d.Id())
		req.Name = ucloud.String(d.Get("name").(string))

		_, err := conn.ModifyURedisGroupName(req)
		if err != nil {
			return fmt.Errorf("do ModifyURedisGroupName failed in update kvstore instance %s, %s", d.Id(), err)
		}

		if err := client.waitActiveStandbyRedisRunning(d.Id()); err != nil {
			return fmt.Errorf("do wait for ModifyURedisGroupName failed in update kvstore instance %s, %s", d.Id(), err)
		}

		d.SetPartial("name")
	}

	if d.HasChange("instance_type") && !d.IsNewResource() {
		o, n := d.GetChange("instance_type")
		newCapability := getKVStoreCapability(n.(string))
		if newCapability == getKVStoreCapability(o.(string)) {
			return fmt.Errorf("instance_type is invalid in update kvstore instance %s, intance type changed but memory capability is not changed", d.Id())
		}

		req := conn.NewResizeURedisGroupRequest()
		req.GroupId = ucloud.String(d.Id())
		req.Size = ucloud.Int(newCapability)

		_, err := conn.ResizeURedisGroup(req)
		if err != nil {
			return fmt.Errorf("do ResizeURedisGroup failed in update kvstore instance %s, %s", d.Id(), err)
		}

		if err := client.waitActiveStandbyRedisRunning(d.Id()); err != nil {
			return fmt.Errorf("do wait for ResizeURedisGroup failed in update kvstore instance %s, %s", d.Id(), err)
		}

		d.SetPartial("instance_type")
	}

	if d.HasChange("password") && !d.IsNewResource() {
		password := d.Get("password").(string)

		req := client.pumemconn.NewModifyURedisGroupPasswordRequest()
		req.GroupId = ucloud.String(d.Id())
		req.Password = ucloud.String(password)

		_, err := client.pumemconn.ModifyURedisGroupPassword(req)
		if err != nil {
			return fmt.Errorf("do ModifyURedisGroupPassword failed in update kvstore instance %s, %s", d.Id(), err)
		}

		if err := client.waitActiveStandbyRedisRunning(d.Id()); err != nil {
			return fmt.Errorf("do wait for ModifyURedisGroupPassword failed in update kvstore instance %s, %s", d.Id(), err)
		}

		d.SetPartial("password")
	}

	if d.HasChange("parameter_group_id") && !d.IsNewResource() {
		configId := d.Get("parameter_group_id").(string)

		req := client.pumemconn.NewChangeURedisConfigRequest()
		req.GroupId = ucloud.String(d.Id())
		req.ConfigId = ucloud.String(configId)

		_, err := client.pumemconn.ChangeURedisConfig(req)
		if err != nil {
			return fmt.Errorf("do ChangeURedisConfig failed in update kvstore instance %s, %s", d.Id(), err)
		}

		if err := client.waitActiveStandbyRedisRunning(d.Id()); err != nil {
			return fmt.Errorf("do wait for ChangeURedisConfig failed in update kvstore instance %s, %s", d.Id(), err)
		}

		d.SetPartial("parameter_group_id")
	}

	d.Partial(false)
	return nil
}

func updateDistributedRedisInstance(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.umemconn

	d.Partial(true)

	if d.HasChange("name") && !d.IsNewResource() {
		req := conn.NewModifyUMemSpaceNameRequest()
		req.SpaceId = ucloud.String(d.Id())
		req.Name = ucloud.String(d.Get("name").(string))

		_, err := conn.ModifyUMemSpaceName(req)
		if err != nil {
			return fmt.Errorf("do ModifyUMemSpaceName failed in update kvstore instance %s, %s", d.Id(), err)
		}

		if err := client.waitDistributedRedisRunning(d.Id()); err != nil {
			return fmt.Errorf("do wait for modify name failed in update kvstore instance %s, %s", d.Id(), err)
		}

		d.SetPartial("name")
	}

	if d.HasChange("instance_type") && !d.IsNewResource() {
		o, n := d.GetChange("instance_type")
		newCapability := getKVStoreCapability(n.(string))
		if newCapability == getKVStoreCapability(o.(string)) {
			return fmt.Errorf("instance_type is invalid in update kvstore instance %s, intance type changed but memory capability is not changed", d.Id())
		}

		req := conn.NewResizeUMemSpaceRequest()
		req.SpaceId = ucloud.String(d.Id())
		req.Size = ucloud.Int(newCapability)

		_, err := conn.ResizeUMemSpace(req)
		if err != nil {
			return fmt.Errorf("do ResizeUMemSpace failed in update kvstore instance %s, %s", d.Id(), err)
		}

		if err := client.waitDistributedRedisRunning(d.Id()); err != nil {
			return fmt.Errorf("do wait for resize failed in update kvstore instance %s, %s", d.Id(), err)
		}

		d.SetPartial("instance_type")
	}

	d.Partial(false)

	return readDistributedRedisInstance(d, meta)
}

func updateActiveStandbyMemcacheInstance(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.pumemconn

	d.Partial(true)

	if d.HasChange("name") && !d.IsNewResource() {
		req := conn.NewModifyUMemcacheGroupNameRequest()
		req.GroupId = ucloud.String(d.Id())
		req.Name = ucloud.String(d.Get("name").(string))

		_, err := conn.ModifyUMemcacheGroupName(req)
		if err != nil {
			return fmt.Errorf("do ModifyUMemcacheGroupName failed in update kvstore instance %s, %s", d.Id(), err)
		}

		if err := client.waitActiveStandbyMemcacheRunning(d.Id()); err != nil {
			return fmt.Errorf("do wait for ModifyUMemcacheGroupName failed in update kvstore instance %s, %s", d.Id(), err)
		}
	}

	if d.HasChange("instance_type") && !d.IsNewResource() {
		o, n := d.GetChange("instance_type")
		newCapability := getKVStoreCapability(n.(string))
		if newCapability == getKVStoreCapability(o.(string)) {
			return fmt.Errorf("instance_type is invalid in update kvstore instance %s, intance type changed but memory capability is not changed", d.Id())
		}

		req := conn.NewResizeUMemcacheGroupRequest()
		req.GroupId = ucloud.String(d.Id())
		req.Size = ucloud.Int(newCapability)

		_, err := conn.ResizeUMemcacheGroup(req)
		if err != nil {
			return fmt.Errorf("do ResizeUMemcacheGroup failed in update kvstore instance %s, %s", d.Id(), err)
		}

		if err := client.waitActiveStandbyMemcacheRunning(d.Id()); err != nil {
			return fmt.Errorf("do wait for ResizeUMemcacheGroup failed in update kvstore instance %s, %s", d.Id(), err)
		}
	}

	d.Partial(false)

	return readActiveStandbyMemcacheInstance(d, meta)
}

func readActiveStandbyRedisInstance(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)

	inst, err := client.describeActiveStandbyRedisById(d.Id())
	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("do %s failed in read kvstore instance %s, %s", "DescribeURedisGroup", d.Id(), err)
	}

	d.Set("name", inst.Name)
	d.Set("tag", inst.Tag)
	d.Set("instance_charge_type", inst.ChargeType)
	d.Set("instance_type", fmt.Sprintf("redis-master-%v", inst.Size))
	d.Set("parameter_group_id", inst.ConfigId)
	d.Set("vpc_id", inst.VPCId)
	d.Set("subnet_id", inst.SubnetId)

	if inst.AutoBackup == "enable" {
		d.Set("backup_begin_time", inst.BackupTime)
	}

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

func readDistributedRedisInstance(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)

	inst, err := client.describeDistributedRedisById(d.Id())
	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("do %s failed in read kvstore instance %s, %s", "DescribeUMemSpace", d.Id(), err)
	}

	d.Set("name", inst.Name)
	d.Set("tag", inst.Tag)
	d.Set("instance_charge_type", inst.ChargeType)
	d.Set("instance_type", fmt.Sprintf("redis-distributed-%v", inst.Size))
	d.Set("vpc_id", inst.VPCId)
	d.Set("subnet_id", inst.SubnetId)

	addresses := []map[string]interface{}{}
	for _, addr := range inst.Address {
		ipItem := map[string]interface{}{
			"ip":   addr.IP,
			"port": addr.Port,
		}
		addresses = append(addresses, ipItem)
	}

	if len(addresses) == 0 {
		return fmt.Errorf("do %s failed in read kvstore instance %s, no availability ip address", "DescribeUMemSpace", d.Id())
	}

	d.Set("ip_set", addresses)
	d.Set("create_time", timestampToString(inst.CreateTime))
	d.Set("expire_time", timestampToString(inst.ExpireTime))
	d.Set("status", inst.State)
	return nil
}

func readActiveStandbyMemcacheInstance(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)

	inst, err := client.describeActiveStandbyMemcacheById(d.Id())
	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("do %s failed in read kvstore instance %s, %s", "DescribeURedisGroup", d.Id(), err)
	}

	d.Set("name", inst.Name)
	d.Set("tag", inst.Tag)
	d.Set("instance_charge_type", inst.ChargeType)
	d.Set("instance_type", fmt.Sprintf("memcache-master-%v", inst.Size))
	d.Set("vpc_id", inst.VPCId)
	d.Set("subnet_id", inst.SubnetId)
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

func deleteActiveStandbyRedisInstance(d *schema.ResourceData, meta interface{}) *resource.RetryError {
	client := meta.(*UCloudClient)
	conn := client.umemconn

	_, err := client.describeActiveStandbyRedisById(d.Id())
	if err != nil {
		if isNotFoundError(err) {
			return nil
		}
		return resource.NonRetryableError(fmt.Errorf("do %s failed in delete kvstore instance %s, %s", "DescribeRedisGroup", d.Id(), err))
	}

	req := conn.NewDeleteURedisGroupRequest()
	req.GroupId = ucloud.String(d.Id())
	if _, err := conn.DeleteURedisGroup(req); err != nil {
		return resource.NonRetryableError(fmt.Errorf("error in delete kvstore instance %s, %s", d.Id(), err))
	}

	_, err = client.describeActiveStandbyRedisById(d.Id())
	if err != nil {
		if isNotFoundError(err) {
			return nil
		}
		return resource.NonRetryableError(fmt.Errorf("do %s failed in delete kvstore instance %s, %s", "DescribeRedisGroup", d.Id(), err))
	}

	return resource.RetryableError(fmt.Errorf("delete kvstore instance but it still exists"))
}

func deleteDistributedRedisInstance(d *schema.ResourceData, meta interface{}) *resource.RetryError {
	client := meta.(*UCloudClient)
	conn := client.umemconn

	_, err := client.describeDistributedRedisById(d.Id())
	if err != nil {
		if isNotFoundError(err) {
			return nil
		}
		return resource.NonRetryableError(fmt.Errorf("do %s failed in delete kvstore instance %s, %s", "DescribeUMemSpace", d.Id(), err))
	}

	req := conn.NewDeleteUMemSpaceRequest()
	req.SpaceId = ucloud.String(d.Id())
	if _, err := conn.DeleteUMemSpace(req); err != nil {
		return resource.NonRetryableError(fmt.Errorf("error in delete kvstore instance %s, %s", d.Id(), err))
	}

	_, err = client.describeDistributedRedisById(d.Id())
	if err != nil {
		if isNotFoundError(err) {
			return nil
		}
		return resource.NonRetryableError(fmt.Errorf("do %s failed in delete kvstore instance %s, %s", "DescribeUMemSpace", d.Id(), err))
	}
	return resource.RetryableError(fmt.Errorf("delete kvstore instance but it still exists"))
}

func deleteActiveStandbyMemcacheInstance(d *schema.ResourceData, meta interface{}) *resource.RetryError {
	client := meta.(*UCloudClient)
	conn := client.umemconn

	_, err := client.describeActiveStandbyMemcacheById(d.Id())
	if err != nil {
		if isNotFoundError(err) {
			return nil
		}
		return resource.NonRetryableError(fmt.Errorf("do %s failed in delete kvstore instance %s, %s", "DescribeUMemcacheGroup", d.Id(), err))
	}

	req := conn.NewDeleteUMemcacheGroupRequest()
	req.GroupId = ucloud.String(d.Id())
	if _, err := conn.DeleteUMemcacheGroup(req); err != nil {
		return resource.NonRetryableError(fmt.Errorf("error in delete kvstore instance %s, %s", d.Id(), err))
	}

	_, err = client.describeActiveStandbyMemcacheById(d.Id())
	if err != nil {
		if isNotFoundError(err) {
			return nil
		}
		return resource.NonRetryableError(fmt.Errorf("do %s failed in delete kvstore instance %s, %s", "DescribeUMemcacheGroup", d.Id(), err))
	}

	return resource.RetryableError(fmt.Errorf("delete kvstore instance but it still exists"))
}

func getKVStoreCapability(instType string) int {
	// skip error, because it has been validated at schema
	t, _ := parseKVStoreInstanceType(instType)
	return t.Memory
}
