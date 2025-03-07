package psm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVRF() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVRFCreate,
		ReadContext:   resourceVRFRead,
		UpdateContext: resourceVRFUpdate,
		DeleteContext: resourceVRFDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVRFImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ingress_security_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"egress_security_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"connection_tracking_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"allow_session_reuse": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_fragments_forwarding": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ingress_nat_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"egress_nat_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ipsec_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"flow_export_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"maximum_cps_per_network": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"maximum_sessions_per_network": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
		},
	}
}

type VRF struct {
	Kind       interface{} `json:"kind"`
	APIVersion interface{} `json:"api-version"`
	Meta       struct {
		Name            string      `json:"name"`
		Tenant          string      `json:"tenant"`
		Namespace       interface{} `json:"namespace"`
		GenerationID    interface{} `json:"generation-id"`
		ResourceVersion interface{} `json:"resource-version"`
		UUID            interface{} `json:"uuid"`
		Labels          interface{} `json:"labels"`
		SelfLink        interface{} `json:"self-link"`
		DisplayName     interface{} `json:"display-name"`
	} `json:"meta"`
	Spec struct {
		Type                                                  string      `json:"type"`
		RouterMacAddress                                      interface{} `json:"router-mac-address"`
		VxlanVni                                              interface{} `json:"vxlan-vni"`
		DefaultIpamPolicy                                     interface{} `json:"default-ipam-policy"`
		IngressSecurityPolicy                                 []string    `json:"ingress-security-policy"`
		EgressSecurityPolicy                                  []string    `json:"egress-security-policy"`
		MaximumCpsPerNetworkPerDistributedServicesEntity      int         `json:"maximum-cps-per-network-per-distributed-services-entity"`
		MaximumSessionsPerNetworkPerDistributedServicesEntity int         `json:"maximum-sessions-per-network-per-distributed-services-entity"`
		FlowExportPolicy                                      []string    `json:"flow-export-policy"`
		IngressNatPolicy                                      []string    `json:"ingress-nat-policy"`
		EgressNatPolicy                                       []string    `json:"egress-nat-policy"`
		IpsecPolicy                                           []string    `json:"ipsec-policy"`
		SelectCPS                                             int         `json:"selectCPS"`
		SelectSessions                                        int         `json:"selectSessions"`
		ConnectionTracking                                    string      `json:"connection-tracking-mode"`
		AllowSessionReuse                                     string      `json:"allow-session-reuse"`
		IpFragmentsForwarding                                 string      `json:"ip-fragments-forwarding"`
	} `json:"spec"`
}

func resourceVRFCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client()

	vrf := &VRF{}
	vrf.Meta.Name = d.Get("name").(string)
	vrf.Meta.Tenant = "default"
	vrf.Spec.Type = "unknown"
	vrf.Spec.ConnectionTracking = d.Get("connection_tracking_mode").(string)
	vrf.Spec.AllowSessionReuse = d.Get("allow_session_reuse").(string)
	vrf.Spec.IpFragmentsForwarding = d.Get("ip_fragments_forwarding").(string)
	vrfName := d.Get("name").(string)
	vrf.Spec.MaximumCpsPerNetworkPerDistributedServicesEntity = d.Get("maximum_cps_per_network").(int)
	vrf.Spec.MaximumSessionsPerNetworkPerDistributedServicesEntity = d.Get("maximum_sessions_per_network").(int)

	if v, ok := d.GetOk("ingress_security_policy"); ok {
		vrf.Spec.IngressSecurityPolicy = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("egress_security_policy"); ok {
		vrf.Spec.EgressSecurityPolicy = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("ingress_nat_policy"); ok {
		vrf.Spec.IngressNatPolicy = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("egress_nat_policy"); ok {
		vrf.Spec.EgressNatPolicy = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("ipsec_policy"); ok {
		vrf.Spec.IpsecPolicy = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("flow_export_policy"); ok {
		vrf.Spec.FlowExportPolicy = expandStringList(v.([]interface{}))
	}

	if vrfName == "default" {
		d.SetId("default")
		return nil
	}

	jsonBytes, err := json.Marshal(vrf)
	if err != nil {
		log.Printf("[ERROR] Error marshalling VRF: %s", err)
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Creating VRF with name: %s", vrf.Meta.Name)

	req, err := http.NewRequestWithContext(ctx, "POST", config.Server+"/configs/network/v1/tenant/default/virtualrouters", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return diag.FromErr(err)
	}

	req.AddCookie(&http.Cookie{Name: "sid", Value: config.SID})

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Error when creating VRF: %s", err)
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("failed to create VRF: HTTP %d %s: %s", resp.StatusCode, resp.Status, bodyBytes)
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "VRF creation failed",
				Detail:   errMsg,
			},
		}
	}

	responseBody := &VRF{}
	if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(responseBody.Meta.UUID.(string))

	log.Printf("[DEBUG] VRF created with UUID: %s", responseBody.Meta.UUID.(string))

	return append(diag.Diagnostics{}, resourceVRFRead(ctx, d, m)...)
}

func resourceVRFRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client()

	url := config.Server + "/configs/network/v1/tenant/default/virtualrouters/" + d.Get("name").(string)
	log.Printf("[DEBUG] Reading VRF with name: %s", d.Get("name").(string))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.AddCookie(&http.Cookie{Name: "sid", Value: config.SID})

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Error when reading VRF: %s", err)
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("failed to read VRF: HTTP %s", resp.Status)
	}

	vrf := &VRF{}
	if err := json.NewDecoder(resp.Body).Decode(vrf); err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", vrf.Meta.Name)
	d.Set("kind", vrf.Kind)
	d.Set("api_version", vrf.APIVersion)
	d.Set("ingress_security_policy", vrf.Spec.IngressSecurityPolicy)
	d.Set("egress_security_policy", vrf.Spec.EgressSecurityPolicy)
	d.Set("connection_tracking_mode", vrf.Spec.ConnectionTracking)
	d.Set("allow_session_reuse", vrf.Spec.AllowSessionReuse)
	d.Set("ip_fragments_forwarding", vrf.Spec.IpFragmentsForwarding)
	d.Set("ingress_nat_policy", vrf.Spec.IngressNatPolicy)
	d.Set("egress_nat_policy", vrf.Spec.EgressNatPolicy)
	d.Set("ipsec_policy", vrf.Spec.IpsecPolicy)
	d.Set("flow_export_policy", vrf.Spec.FlowExportPolicy)
	d.Set("maximum_cps_per_network", vrf.Spec.MaximumCpsPerNetworkPerDistributedServicesEntity)
	d.Set("maximum_sessions_per_network", vrf.Spec.MaximumSessionsPerNetworkPerDistributedServicesEntity)

	return nil
}

func resourceVRFDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client()
	vrfName := d.Get("name").(string)

	if vrfName == "default" {
		d.SetId("default")
		return nil
	}

	url := config.Server + "/configs/network/v1/tenant/default/virtualrouters/" + d.Get("name").(string)

	log.Printf("[DEBUG] Deleting VRF with name: %s", d.Get("name").(string))

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.AddCookie(&http.Cookie{Name: "sid", Value: config.SID})

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Error when deleting VRF: %s", err)
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("failed to delete VRF: HTTP %s", resp.Status)
	}

	d.SetId("")

	return nil
}

func resourceVRFUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client()

	vrf := &VRF{}
	vrf.Meta.Name = d.Get("name").(string)
	vrf.Meta.Tenant = "default"
	vrf.Spec.Type = "unknown"
	vrf.Spec.ConnectionTracking = d.Get("connection_tracking_mode").(string)
	vrf.Spec.AllowSessionReuse = d.Get("allow_session_reuse").(string)
	vrf.Spec.IpFragmentsForwarding = d.Get("ip_fragments_forwarding").(string)
	vrfName := d.Get("name").(string)
	vrf.Spec.MaximumCpsPerNetworkPerDistributedServicesEntity = d.Get("maximum_cps_per_network").(int)
	vrf.Spec.MaximumSessionsPerNetworkPerDistributedServicesEntity = d.Get("maximum_sessions_per_network").(int)

	if v, ok := d.GetOk("ingress_security_policy"); ok {
		vrf.Spec.IngressSecurityPolicy = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("egress_security_policy"); ok {
		vrf.Spec.EgressSecurityPolicy = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("ingress_nat_policy"); ok {
		vrf.Spec.IngressNatPolicy = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("egress_nat_policy"); ok {
		vrf.Spec.EgressNatPolicy = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("ipsec_policy"); ok {
		vrf.Spec.IpsecPolicy = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("flow_export_policy"); ok {
		vrf.Spec.FlowExportPolicy = expandStringList(v.([]interface{}))
	}

	if vrfName == "default" {
		d.SetId("default")
		return nil
	}

	jsonBytes, err := json.Marshal(vrf)
	if err != nil {
		log.Printf("[ERROR] Error marshalling VRF: %s", err)
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Updating VRF with name: %s", vrf.Meta.Name)

	req, err := http.NewRequestWithContext(ctx, "PUT", config.Server+"/configs/network/v1/tenant/default/virtualrouters/"+vrf.Meta.Name, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return diag.FromErr(err)
	}

	req.AddCookie(&http.Cookie{Name: "sid", Value: config.SID})

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Error when updating VRF: %s", err)
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("failed to update VRF: HTTP %d %s: %s", resp.StatusCode, resp.Status, bodyBytes)
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "VRF update failed",
				Detail:   errMsg,
			},
		}
	}

	responseBody := &VRF{}
	if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(responseBody.Meta.UUID.(string))

	log.Printf("[DEBUG] VRF updated with UUID: %s", responseBody.Meta.UUID.(string))

	return append(diag.Diagnostics{}, resourceVRFRead(ctx, d, m)...)
}

func resourceVRFImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	config := m.(*Config)

	name := d.Id()

	url := fmt.Sprintf("%s/configs/network/v1/tenant/default/virtualrouters/%s", config.Server, name)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request to import VRF: %v", err)
	}

	req.AddCookie(&http.Cookie{Name: "sid", Value: config.SID})

	client := config.Client()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error importing VRF: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to import VRF: HTTP %s", resp.Status)
	}

	var vrf VRF
	if err := json.NewDecoder(resp.Body).Decode(&vrf); err != nil {
		return nil, fmt.Errorf("error decoding VRF response: %v", err)
	}

	d.SetId(vrf.Meta.UUID.(string))
	d.Set("name", vrf.Meta.Name)
	d.Set("ingress_security_policy", vrf.Spec.IngressSecurityPolicy)
	d.Set("egress_security_policy", vrf.Spec.EgressSecurityPolicy)
	d.Set("connection_tracking_mode", vrf.Spec.ConnectionTracking)
	d.Set("allow_session_reuse", vrf.Spec.AllowSessionReuse)
	d.Set("ip_fragments_forwarding", vrf.Spec.IpFragmentsForwarding)
	d.Set("ingress_nat_policy", vrf.Spec.IngressNatPolicy)
	d.Set("egress_nat_policy", vrf.Spec.EgressNatPolicy)
	d.Set("ipsec_policy", vrf.Spec.IpsecPolicy)
	d.Set("flow_export_policy", vrf.Spec.FlowExportPolicy)

	return []*schema.ResourceData{d}, nil
}
