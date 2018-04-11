package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/dns"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmDnsZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDnsZoneCreate,
		Read:   resourceArmDnsZoneRead,
		Update: resourceArmDnsZoneCreate,
		Delete: resourceArmDnsZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: resourceAzurermResourceGroupNameDiffSuppress,
			},

			"number_of_record_sets": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"max_number_of_record_sets": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name_servers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmDnsZoneCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).zonesClient

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := "global"

	tags := d.Get("tags").(map[string]interface{})

	parameters := dns.Zone{
		Location: &location,
		Tags:     expandTags(tags),
	}

	etag := ""
	ifNoneMatch := "" // set to empty to allow updates to records after creation
	resp, err := client.CreateOrUpdate(resGroup, name, parameters, etag, ifNoneMatch)
	if err != nil {
		return err
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read DNS zone %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmDnsZoneRead(d, meta)
}

func resourceArmDnsZoneRead(d *schema.ResourceData, meta interface{}) error {
	zonesClient := meta.(*ArmClient).zonesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["dnszones"]

	resp, err := zonesClient.Get(resGroup, name)
	if err != nil {
		if responseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading DNS zone %s (resource group %s): %+v", name, resGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("number_of_record_sets", resp.NumberOfRecordSets)
	d.Set("max_number_of_record_sets", resp.MaxNumberOfRecordSets)

	nameServers := make([]string, 0, len(*resp.NameServers))
	for _, ns := range *resp.NameServers {
		nameServers = append(nameServers, ns)
	}
	if err := d.Set("name_servers", nameServers); err != nil {
		return err
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmDnsZoneDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).zonesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["dnszones"]

	etag := ""
	_, error := client.Delete(resGroup, name, etag, make(chan struct{}))
	err = <-error

	if err != nil {
		return fmt.Errorf("Error deleting DNS zone %s (resource group %s): %+v", name, resGroup, err)
	}

	return nil
}
