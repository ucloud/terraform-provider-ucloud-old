package ucloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
)

func resourceUCloudEIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceUCloudEIPCreate,
		Read:   resourceUCloudEIPRead,
		Update: resourceUCloudEIPUpdate,
		Delete: resourceUCloudEIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bandwidth": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 800),
			},

			"internet_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "bgp",
				ValidateFunc: validateStringInChoices([]string{"bgp", "international"}),
			},

			"internet_charge_type": &schema.Schema{
				Type:         schema.TypeString,
				Default:      "month",
				Optional:     true,
				ValidateFunc: validateStringInChoices([]string{"month", "year", "dynamic"}),
			},

			"internet_charge_mode": &schema.Schema{
				Type:         schema.TypeString,
				Default:      "bandwidth",
				Optional:     true,
				ValidateFunc: validateStringInChoices([]string{"traffic", "bandwidth"}),
			},

			"eip_duration": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 9),
			},

			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateName,
			},

			"remark": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tag": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateName,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_set": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"internet_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"resource": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"resource_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"create_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"expire_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceUCloudEIPCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.unetconn

	req := conn.NewAllocateEIPRequest()
	req.Bandwidth = ucloud.Int(d.Get("bandwidth").(int))
	req.Quantity = ucloud.Int(d.Get("eip_duration").(int))
	req.ChargeType = ucloud.String(d.Get("internet_charge_type").(string))
	req.PayMode = ucloud.String(d.Get("internet_charge_mode").(string))
	req.OperatorName = ucloud.String(upperCamelCvt.mustUnconvert(d.Get("internet_type").(string)))

	if val, ok := d.GetOk("name"); ok {
		req.Name = ucloud.String(val.(string))
	}

	if val, ok := d.GetOk("tag"); ok {
		req.Tag = ucloud.String(val.(string))
	}

	if val, ok := d.GetOk("remark"); ok {
		req.Remark = ucloud.String(val.(string))
	}

	resp, err := conn.AllocateEIP(req)
	if err != nil {
		query, err := request.ToQueryMap(req)
		fmt.Printf("[DEBUG] %s %s", d.Get("internet_type").(string), upperCamelCvt.mustUnconvert(d.Get("internet_type").(string)))
		fmt.Printf("[DEBUG] allocate eip arguments %+v\n", query)
		return fmt.Errorf("error on creating eip, %s", err)
	}

	if len(resp.EIPSet) != 1 {
		return fmt.Errorf("error on creating eip, expected exactly one eip, got %v", len(resp.EIPSet))
	}

	eip := resp.EIPSet[0]
	d.SetId(eip.EIPId)

	// after create eip, we need to wait it initialized
	stateConf := eipWaitForState(client, d.Id())

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error on waiting for eip is initialized when creating %s, %s", d.Id(), err)
	}

	return resourceUCloudEIPUpdate(d, meta)
}

func resourceUCloudEIPUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.unetconn

	d.Partial(true)

	if d.HasChange("bandwidth") && !d.IsNewResource() {
		reqBand := conn.NewModifyEIPBandwidthRequest()
		reqBand.EIPId = ucloud.String(d.Id())
		reqBand.Bandwidth = ucloud.Int(d.Get("bandwidth").(int))

		_, err := conn.ModifyEIPBandwidth(reqBand)
		if err != nil {
			return fmt.Errorf("error on %s to eip %s, %s", "ModifyEIPBandwidth", d.Id(), err)
		}

		d.SetPartial("bandwidth")

		// after update eip bandwidth, we need to wait it completed
		stateConf := eipWaitForState(client, d.Id())

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("error on waiting for eip is completed when updating %s, %s", d.Id(), err)
		}
	}

	if d.HasChange("internet_charge_mode") && !d.IsNewResource() {
		reqCharge := conn.NewSetEIPPayModeRequest()
		reqCharge.EIPId = ucloud.String(d.Id())
		reqCharge.PayMode = ucloud.String(d.Get("internet_charge_mode").(string))
		reqCharge.Bandwidth = ucloud.Int(d.Get("bandwidth").(int))

		_, err := conn.SetEIPPayMode(reqCharge)
		if err != nil {
			return fmt.Errorf("error on %s to eip %s, %s", "SetEIPPayMode", d.Id(), err)
		}

		d.SetPartial("internet_charge_mode")

		// after update eip internet charge mode, we need to wait it completed
		stateConf := eipWaitForState(client, d.Id())

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("error on waiting for eip is completed when updating %s, %s", d.Id(), err)
		}
	}

	isChanged := false
	reqAttribute := conn.NewUpdateEIPAttributeRequest()
	reqAttribute.EIPId = ucloud.String(d.Id())

	if d.HasChange("name") && !d.IsNewResource() {
		reqAttribute.Name = ucloud.String(d.Get("name").(string))
		isChanged = true
	}

	if d.HasChange("tag") && !d.IsNewResource() {
		reqAttribute.Tag = ucloud.String(d.Get("tag").(string))
		isChanged = true
	}

	if d.HasChange("remark") && !d.IsNewResource() {
		reqAttribute.Remark = ucloud.String(d.Get("remark").(string))
		isChanged = true
	}

	if isChanged {
		_, err := conn.UpdateEIPAttribute(reqAttribute)
		if err != nil {
			return fmt.Errorf("error on %s to eip %s, %s", "UpdateEIPAttribute", d.Id(), err)
		}

		d.SetPartial("name")
		d.SetPartial("tag")
		d.SetPartial("remark")

		// after eip update eip attribute, we need to wait it completed
		stateConf := eipWaitForState(client, d.Id())

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("error on waiting for eip is completed when updating %s, %s", d.Id(), err)
		}
	}

	d.Partial(false)

	return resourceUCloudEIPRead(d, meta)
}

func resourceUCloudEIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)

	eip, err := client.describeEIPById(d.Id())
	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error on reading eip %s, %s", d.Id(), err)
	}

	d.Set("bandwidth", eip.Bandwidth)
	d.Set("internet_charge_type", eip.ChargeType)
	d.Set("internet_charge_mode", eip.PayMode)
	d.Set("name", eip.Name)
	d.Set("remark", eip.Remark)
	d.Set("tag", eip.Tag)
	d.Set("status", eip.Status)
	d.Set("create_time", timestampToString(eip.CreateTime))
	d.Set("expire_time", timestampToString(eip.ExpireTime))

	eipAddr := []map[string]interface{}{}
	for _, item := range eip.EIPAddr {
		eipAddr = append(eipAddr, map[string]interface{}{
			"ip":            item.IP,
			"internet_type": item.OperatorName,
		})
	}
	if err := d.Set("ip_set", eipAddr); err != nil {
		return fmt.Errorf("...")
	}

	if err := d.Set("resource", map[string]string{
		"resource_type": lowerCaseProdCvt.mustUnconvert(eip.Resource.ResourceType),
		"resource_id":   eip.Resource.ResourceId,
	}); err != nil {
		return fmt.Errorf("...")
	}

	return nil
}

func resourceUCloudEIPDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*UCloudClient)
	conn := client.unetconn

	req := conn.NewReleaseEIPRequest()
	req.EIPId = ucloud.String(d.Id())

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := conn.ReleaseEIP(req); err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on deleting eip %s, %s", d.Id(), err))
		}

		_, err := client.describeEIPById(d.Id())
		if err != nil {
			if isNotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("error on reading eip when deleting %s, %s", d.Id(), err))
		}

		return resource.RetryableError(fmt.Errorf("the specified eip %s has not been deleted due to unknown error", d.Id()))
	})
}

func eipWaitForState(client *UCloudClient, eipID string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"free"},
		Timeout:    5 * time.Minute,
		Delay:      2 * time.Second,
		MinTimeout: 1 * time.Second,
		Refresh: func() (interface{}, string, error) {
			eip, err := client.describeEIPById(eipID)
			if err != nil {
				if isNotFoundError(err) {
					return nil, "pending", nil
				}
				return nil, "", err
			}

			state := eip.Status
			if state != "free" {
				state = "pending"
			}

			return eip, state, nil
		},
	}
}
