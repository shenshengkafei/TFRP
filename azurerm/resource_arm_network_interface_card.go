package azurerm

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmNetworkInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetworkInterfaceCreate,
		Read:   resourceArmNetworkInterfaceRead,
		Update: resourceArmNetworkInterfaceCreate,
		Delete: resourceArmNetworkInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"network_security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mac_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"private_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"virtual_machine_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"ip_configuration": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"subnet_id": {
							Type:     schema.TypeString,
							Required: true,
						},

						"private_ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"private_ip_address_allocation": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateFunc:     validateNetworkInterfacePrivateIpAddressAllocation,
							StateFunc:        ignoreCaseStateFunc,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"public_ip_address_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"load_balancer_backend_address_pools_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"load_balancer_inbound_nat_rules_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
					},
				},
				Set: resourceArmNetworkInterfaceIpConfigurationHash,
			},

			"dns_servers": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"internal_dns_name_label": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"applied_dns_servers": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"internal_fqdn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"enable_ip_forwarding": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmNetworkInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	ifaceClient := client.ifaceClient

	log.Printf("[INFO] preparing arguments for Azure ARM Network Interface creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resGroup := d.Get("resource_group_name").(string)
	enableIpForwarding := d.Get("enable_ip_forwarding").(bool)
	tags := d.Get("tags").(map[string]interface{})

	properties := network.InterfacePropertiesFormat{
		EnableIPForwarding: &enableIpForwarding,
	}

	if v, ok := d.GetOk("network_security_group_id"); ok {
		nsgId := v.(string)
		properties.NetworkSecurityGroup = &network.SecurityGroup{
			ID: &nsgId,
		}

		networkSecurityGroupName, err := parseNetworkSecurityGroupName(nsgId)
		if err != nil {
			return err
		}

		azureRMLockByName(networkSecurityGroupName, networkSecurityGroupResourceName)
		defer azureRMUnlockByName(networkSecurityGroupName, networkSecurityGroupResourceName)
	}

	dns, hasDns := d.GetOk("dns_servers")
	nameLabel, hasNameLabel := d.GetOk("internal_dns_name_label")
	if hasDns || hasNameLabel {
		ifaceDnsSettings := network.InterfaceDNSSettings{}

		if hasDns {
			var dnsServers []string
			dns := dns.(*schema.Set).List()
			for _, v := range dns {
				str := v.(string)
				dnsServers = append(dnsServers, str)
			}
			ifaceDnsSettings.DNSServers = &dnsServers
		}

		if hasNameLabel {
			name_label := nameLabel.(string)
			ifaceDnsSettings.InternalDNSNameLabel = &name_label
		}

		properties.DNSSettings = &ifaceDnsSettings
	}

	ipConfigs, subnetnToLock, vnnToLock, sgErr := expandAzureRmNetworkInterfaceIpConfigurations(d)
	if sgErr != nil {
		return fmt.Errorf("Error Building list of Network Interface IP Configurations: %s", sgErr)
	}

	azureRMLockMultipleByName(subnetnToLock, subnetResourceName)
	defer azureRMUnlockMultipleByName(subnetnToLock, subnetResourceName)

	azureRMLockMultipleByName(vnnToLock, virtualNetworkResourceName)
	defer azureRMUnlockMultipleByName(vnnToLock, virtualNetworkResourceName)

	if len(ipConfigs) > 0 {
		properties.IPConfigurations = &ipConfigs
	}

	iface := network.Interface{
		Name:                      &name,
		Location:                  &location,
		InterfacePropertiesFormat: &properties,
		Tags: expandTags(tags),
	}

	_, error := ifaceClient.CreateOrUpdate(resGroup, name, iface, make(chan struct{}))
	err := <-error
	if err != nil {
		return err
	}

	read, err := ifaceClient.Get(resGroup, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read NIC %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmNetworkInterfaceRead(d, meta)
}

func resourceArmNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	ifaceClient := meta.(*ArmClient).ifaceClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["networkInterfaces"]

	resp, err := ifaceClient.Get(resGroup, name, "")
	if err != nil {
		if responseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Network Interface %s: %s", name, err)
	}

	iface := *resp.InterfacePropertiesFormat

	if iface.MacAddress != nil {
		if *iface.MacAddress != "" {
			d.Set("mac_address", iface.MacAddress)
		}
	}

	if iface.IPConfigurations != nil && len(*iface.IPConfigurations) > 0 {
		var privateIPAddress *string
		///TODO: Change this to a loop when https://github.com/Azure/azure-sdk-for-go/issues/259 is fixed
		if (*iface.IPConfigurations)[0].InterfaceIPConfigurationPropertiesFormat != nil {
			privateIPAddress = (*iface.IPConfigurations)[0].InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress
		}

		if *privateIPAddress != "" {
			d.Set("private_ip_address", *privateIPAddress)
		}
	}

	if iface.IPConfigurations != nil {
		d.Set("ip_configuration", schema.NewSet(resourceArmNetworkInterfaceIpConfigurationHash, flattenNetworkInterfaceIPConfigurations(iface.IPConfigurations)))
	}

	if iface.VirtualMachine != nil {
		if *iface.VirtualMachine.ID != "" {
			d.Set("virtual_machine_id", *iface.VirtualMachine.ID)
		}
	}

	var appliedDNSServers []string
	var dnsServers []string
	if iface.DNSSettings != nil {
		if iface.DNSSettings.AppliedDNSServers != nil && len(*iface.DNSSettings.AppliedDNSServers) > 0 {
			for _, applied := range *iface.DNSSettings.AppliedDNSServers {
				appliedDNSServers = append(appliedDNSServers, applied)
			}
		}

		if iface.DNSSettings.DNSServers != nil && len(*iface.DNSSettings.DNSServers) > 0 {
			for _, dns := range *iface.DNSSettings.DNSServers {
				dnsServers = append(dnsServers, dns)
			}
		}

		if iface.DNSSettings.InternalFqdn != nil && *iface.DNSSettings.InternalFqdn != "" {
			d.Set("internal_fqdn", iface.DNSSettings.InternalFqdn)
		}

		d.Set("internal_dns_name_label", iface.DNSSettings.InternalDNSNameLabel)
	}

	if iface.NetworkSecurityGroup != nil {
		d.Set("network_security_group_id", resp.NetworkSecurityGroup.ID)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))
	d.Set("applied_dns_servers", appliedDNSServers)
	d.Set("dns_servers", dnsServers)
	d.Set("enable_ip_forwarding", resp.EnableIPForwarding)

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmNetworkInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	ifaceClient := meta.(*ArmClient).ifaceClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["networkInterfaces"]

	if v, ok := d.GetOk("network_security_group_id"); ok {
		networkSecurityGroupId := v.(string)
		networkSecurityGroupName, err := parseNetworkSecurityGroupName(networkSecurityGroupId)
		if err != nil {
			return err
		}

		azureRMLockByName(networkSecurityGroupName, networkSecurityGroupResourceName)
		defer azureRMUnlockByName(networkSecurityGroupName, networkSecurityGroupResourceName)
	}

	configs := d.Get("ip_configuration").(*schema.Set).List()
	subnetNamesToLock := make([]string, 0)
	virtualNetworkNamesToLock := make([]string, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		subnet_id := data["subnet_id"].(string)
		subnetId, err := parseAzureResourceID(subnet_id)
		if err != nil {
			return err
		}
		subnetName := subnetId.Path["subnets"]
		subnetNamesToLock = append(subnetNamesToLock, subnetName)

		virtualNetworkName := subnetId.Path["virtualNetworks"]
		virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, virtualNetworkName)
	}

	azureRMLockMultipleByName(&subnetNamesToLock, subnetResourceName)
	defer azureRMUnlockMultipleByName(&subnetNamesToLock, subnetResourceName)

	azureRMLockMultipleByName(&virtualNetworkNamesToLock, virtualNetworkResourceName)
	defer azureRMUnlockMultipleByName(&virtualNetworkNamesToLock, virtualNetworkResourceName)

	_, error := ifaceClient.Delete(resGroup, name, make(chan struct{}))
	err = <-error

	return err
}

func resourceArmNetworkInterfaceIpConfigurationHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["subnet_id"].(string)))
	if m["private_ip_address"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["private_ip_address"].(string)))
	}
	buf.WriteString(fmt.Sprintf("%s-", m["private_ip_address_allocation"].(string)))
	if m["public_ip_address_id"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["public_ip_address_id"].(string)))
	}
	if m["load_balancer_backend_address_pools_ids"] != nil {
		ids := m["load_balancer_backend_address_pools_ids"].(*schema.Set).List()
		for _, id := range ids {
			buf.WriteString(fmt.Sprintf("%d-", schema.HashString(id.(string))))
		}
	}
	if m["load_balancer_inbound_nat_rules_ids"] != nil {
		ids := m["load_balancer_inbound_nat_rules_ids"].(*schema.Set).List()
		for _, id := range ids {
			buf.WriteString(fmt.Sprintf("%d-", schema.HashString(id.(string))))
		}
	}

	return hashcode.String(buf.String())
}

func validateNetworkInterfacePrivateIpAddressAllocation(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	allocations := map[string]bool{
		"static":  true,
		"dynamic": true,
	}

	if !allocations[value] {
		errors = append(errors, fmt.Errorf("Network Interface Allocations can only be Static or Dynamic"))
	}
	return
}

func flattenNetworkInterfaceIPConfigurations(ipConfigs *[]network.InterfaceIPConfiguration) []interface{} {
	result := make([]interface{}, 0, len(*ipConfigs))
	for _, ipConfig := range *ipConfigs {
		niIPConfig := make(map[string]interface{})
		niIPConfig["name"] = *ipConfig.Name
		niIPConfig["subnet_id"] = *ipConfig.InterfaceIPConfigurationPropertiesFormat.Subnet.ID
		niIPConfig["private_ip_address_allocation"] = strings.ToLower(string(ipConfig.InterfaceIPConfigurationPropertiesFormat.PrivateIPAllocationMethod))

		if ipConfig.InterfaceIPConfigurationPropertiesFormat.PrivateIPAllocationMethod == network.Static {
			niIPConfig["private_ip_address"] = *ipConfig.InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress
		}

		if ipConfig.InterfaceIPConfigurationPropertiesFormat.PublicIPAddress != nil {
			niIPConfig["public_ip_address_id"] = *ipConfig.InterfaceIPConfigurationPropertiesFormat.PublicIPAddress.ID
		}

		var pools []interface{}
		if ipConfig.InterfaceIPConfigurationPropertiesFormat.LoadBalancerBackendAddressPools != nil {
			for _, pool := range *ipConfig.InterfaceIPConfigurationPropertiesFormat.LoadBalancerBackendAddressPools {
				pools = append(pools, *pool.ID)
			}
		}
		niIPConfig["load_balancer_backend_address_pools_ids"] = schema.NewSet(schema.HashString, pools)

		var rules []interface{}
		if ipConfig.InterfaceIPConfigurationPropertiesFormat.LoadBalancerInboundNatRules != nil {
			for _, rule := range *ipConfig.InterfaceIPConfigurationPropertiesFormat.LoadBalancerInboundNatRules {
				rules = append(rules, *rule.ID)
			}
		}
		niIPConfig["load_balancer_inbound_nat_rules_ids"] = schema.NewSet(schema.HashString, rules)

		result = append(result, niIPConfig)
	}
	return result
}

func expandAzureRmNetworkInterfaceIpConfigurations(d *schema.ResourceData) ([]network.InterfaceIPConfiguration, *[]string, *[]string, error) {
	configs := d.Get("ip_configuration").(*schema.Set).List()
	ipConfigs := make([]network.InterfaceIPConfiguration, 0, len(configs))
	subnetNamesToLock := make([]string, 0)
	virtualNetworkNamesToLock := make([]string, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		subnet_id := data["subnet_id"].(string)
		private_ip_allocation_method := data["private_ip_address_allocation"].(string)

		var allocationMethod network.IPAllocationMethod
		switch strings.ToLower(private_ip_allocation_method) {
		case "dynamic":
			allocationMethod = network.Dynamic
		case "static":
			allocationMethod = network.Static
		default:
			return []network.InterfaceIPConfiguration{}, nil, nil, fmt.Errorf(
				"valid values for private_ip_allocation_method are 'dynamic' and 'static' - got '%s'",
				private_ip_allocation_method)
		}

		properties := network.InterfaceIPConfigurationPropertiesFormat{
			Subnet: &network.Subnet{
				ID: &subnet_id,
			},
			PrivateIPAllocationMethod: allocationMethod,
		}

		subnetId, err := parseAzureResourceID(subnet_id)
		if err != nil {
			return []network.InterfaceIPConfiguration{}, nil, nil, err
		}
		subnetName := subnetId.Path["subnets"]
		virtualNetworkName := subnetId.Path["virtualNetworks"]
		subnetNamesToLock = append(subnetNamesToLock, subnetName)
		virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, virtualNetworkName)

		if v := data["private_ip_address"].(string); v != "" {
			properties.PrivateIPAddress = &v
		}

		if v := data["public_ip_address_id"].(string); v != "" {
			properties.PublicIPAddress = &network.PublicIPAddress{
				ID: &v,
			}
		}

		if v, ok := data["load_balancer_backend_address_pools_ids"]; ok {
			var ids []network.BackendAddressPool
			pools := v.(*schema.Set).List()
			for _, p := range pools {
				pool_id := p.(string)
				id := network.BackendAddressPool{
					ID: &pool_id,
				}

				ids = append(ids, id)
			}

			properties.LoadBalancerBackendAddressPools = &ids
		}

		if v, ok := data["load_balancer_inbound_nat_rules_ids"]; ok {
			var natRules []network.InboundNatRule
			rules := v.(*schema.Set).List()
			for _, r := range rules {
				rule_id := r.(string)
				rule := network.InboundNatRule{
					ID: &rule_id,
				}

				natRules = append(natRules, rule)
			}

			properties.LoadBalancerInboundNatRules = &natRules
		}

		name := data["name"].(string)
		ipConfig := network.InterfaceIPConfiguration{
			Name: &name,
			InterfaceIPConfigurationPropertiesFormat: &properties,
		}

		ipConfigs = append(ipConfigs, ipConfig)
	}

	return ipConfigs, &subnetNamesToLock, &virtualNetworkNamesToLock, nil
}
