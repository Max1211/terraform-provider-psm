package psm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tenant": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"virtual_router": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ingress_security_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"egress_security_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connection_tracking_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"allow_session_reuse": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"service_bypass": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
			},
			"ip_fragments_forwarding": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"ingress_mirror_session": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of mirror session to be applied on ingress traffic",
			},
			"egress_mirror_session": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of mirror session to be applied on egress traffic",
			},
		},
	}
}

type Network struct {
	Meta struct {
		Kind            interface{} `json:"kind" default:"null"`
		APIVersion      interface{} `json:"api-version" default:"null"`
		Name            string      `json:"name"`
		Tenant          string      `json:"tenant"`
		Namespace       interface{} `json:"namespace" default:"null"`
		GenerationID    interface{} `json:"generation-id" default:"null"`
		ResourceVersion interface{} `json:"resource-version" default:"null"`
		UUID            string      `json:"uuid" default:"null"`
		Labels          interface{} `json:"labels" default:"null"`
		SelfLink        interface{} `json:"self-link" default:"null"`
		DisplayName     interface{} `json:"display-name" default:"null"`
	}
	Spec struct {
		Ipv4Subnet            interface{}   `json:"ipv4-subnet" default:"null"`
		Ipv4Gateway           interface{}   `json:"ipv4-gateway" default:"null"`
		Ipv6Subnet            interface{}   `json:"ipv6-subnet" default:"null"`
		Ipv6Gateway           interface{}   `json:"ipv6-gateway" default:"null"`
		VxlanVni              interface{}   `json:"vxlan-vni" default:"null"`
		IpamPolicy            interface{}   `json:"ipam-policy" default:"null"`
		Orchestrators         []interface{} `json:"orchestrators"`
		IngressSecurityPolicy []interface{} `json:"ingress-security-policy" default:"null"`
		EgressSecurityPolicy  []interface{} `json:"egress-security-policy" default:"null"`
		IngressMirrorSession  []interface{} `json:"ingress-mirror-session,omitempty"`
		EgressMirrorSession   []interface{} `json:"egress-mirror-session,omitempty"`
		FirewallProfile       struct {
			MaximumCpsPerDistributedServicesEntity      int `json:"maximum-cps-per-distributed-services-entity" default:"-1"`
			MaximumSessionsPerDistributedServicesEntity int `json:"maximum-sessions-per-distributed-services-entity" default:"-1"`
		} `json:"firewall-profile"`
		RouteImportExport     interface{} `json:"route-import-export" default:"null"`
		Type                  string      `json:"type" default:"bridged"`
		VirtualRouter         string      `json:"virtual-router" default:"default"`
		VRF                   string      `json:"vrf" default:"default"`
		ConnectionTracking    string      `json:"connection-tracking-mode" default:"inherit from vrf"`
		AllowSessionReuse     string      `json:"allow-session-reuse" default:"inherit from vrf"`
		IpFragmentsForwarding string      `json:"ip-fragments-forwarding" default:"inherit from vrf"`
		VlanID                int         `json:"vlan-id"`
		SelectVlanOrIpv4      int         `json:"selectVlanOrIpv4" default:"1"`
		SelectCPS             int         `json:"selectCPS" default:"-1"`
		SelectSessions        int         `json:"selectSessions" default:"-1"`
		ServiceBypass         bool        `json:"service-bypass" default:"false"`
	}
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client()

	network := &Network{}
	network.Meta.Name = d.Get("name").(string)
	network.Meta.Tenant = d.Get("tenant").(string)
	network.Spec.VlanID = d.Get("vlan_id").(int)
	network.Spec.Type = "bridged"
	network.Meta.Namespace = "default"
	network.Spec.VirtualRouter = d.Get("virtual_router").(string)
	network.Spec.ConnectionTracking = d.Get("connection_tracking_mode").(string)
	network.Spec.AllowSessionReuse = d.Get("allow_session_reuse").(string)
	network.Spec.ServiceBypass = d.Get("service_bypass").(bool)
	network.Spec.IpFragmentsForwarding = d.Get("ip_fragments_forwarding").(string)

	// Check if the ingress_security_policy and egress_security_policy values are provided and set them
	if v, ok := d.GetOk("ingress_security_policy"); ok {
		network.Spec.IngressSecurityPolicy = []interface{}{v.(string)}
	}
	if v, ok := d.GetOk("egress_security_policy"); ok {
		network.Spec.EgressSecurityPolicy = []interface{}{v.(string)}
	}

	// Check if mirror sessions are provided
	if v, ok := d.GetOk("ingress_mirror_session"); ok {
		network.Spec.IngressMirrorSession = []interface{}{v.(string)}
	}
	if v, ok := d.GetOk("egress_mirror_session"); ok {
		network.Spec.EgressMirrorSession = []interface{}{v.(string)}
	}

	// Convert the Network struct to JSON.
	jsonBytes, err := json.Marshal(network)
	if err != nil {
		return diag.FromErr(err)
	}

	//req, err := http.NewRequestWithContext(ctx, "POST", "/configs/network/v1/tenant/default/networks", bytes.NewBuffer(jsonBytes))
	req, err := http.NewRequestWithContext(ctx, "POST", config.Server+"/configs/network/v1/tenant/default/networks", bytes.NewBuffer(jsonBytes))

	if err != nil {
		return diag.FromErr(err)
	}

	// Set SID cookie for authentication
	req.AddCookie(&http.Cookie{Name: "sid", Value: config.SID})

	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("failed to create network: HTTP %d %s: %s", resp.StatusCode, resp.Status, bodyBytes)
		// Added for additional debug if the JSON we send to the PSM server is invalid.
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "Network creation failed",
				Detail:   errMsg,
			},
		}
	}

	responseBody := &Network{}
	if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
		return diag.FromErr(err)
	}

	// Set the Terraform resource ID to the UUID returned by the API.
	d.SetId(responseBody.Meta.UUID)

	return append(diag.Diagnostics{}, resourceNetworkRead(ctx, d, m)...)

	//return diag.Diagnostics{}
}

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config) // Cast to Config instead of *http.Client
	client := config.Client()

	// Construct the URL for the network based on its name attached at the end
	url := config.Server + "/configs/network/v1/tenant/default/networks/" + d.Get("name").(string)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set SID cookie for authentication
	req.AddCookie(&http.Cookie{Name: "sid", Value: config.SID})

	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("failed to read network: HTTP %s", resp.Status)
	}

	// Decode the response body into the Network struct
	network := &Network{}

	if err := json.NewDecoder(resp.Body).Decode(network); err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", network.Meta.Name)
	d.Set("vlan_id", network.Spec.VlanID)

	// Read mirror session configurations if present
	if len(network.Spec.IngressMirrorSession) > 0 {
		if mirrorSession, ok := network.Spec.IngressMirrorSession[0].(string); ok {
			d.Set("ingress_mirror_session", mirrorSession)
		}
	}

	if len(network.Spec.EgressMirrorSession) > 0 {
		if mirrorSession, ok := network.Spec.EgressMirrorSession[0].(string); ok {
			d.Set("egress_mirror_session", mirrorSession)
		}
	}

	return nil
}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	isDebugEnabled := func() bool {
		return os.Getenv("TF_LOG") == "debug"
	}

	config := m.(*Config)
	client := config.Client()

	url := config.Server + "/configs/network/v1/tenant/default/networks/" + d.Get("name").(string)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	req.AddCookie(&http.Cookie{Name: "sid", Value: config.SID})

	resp, err := client.Do(req)
	if err != nil {
		if isDebugEnabled() {
			log.Printf("[DEBUG] Error getting current network state: %s", err)
		}
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if isDebugEnabled() {
			log.Printf("[DEBUG] Unexpected HTTP status when getting current network state: %s", resp.Status)
		}
		return diag.Errorf("failed to get current network state: HTTP %s", resp.Status)
	}

	networkCurrent := &Network{}
	if err := json.NewDecoder(resp.Body).Decode(networkCurrent); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("virtual_router") {
		if val, ok := d.GetOk("virtual_router"); ok {
			newVirtualRouter := val.(string)
			networkCurrent.Spec.VirtualRouter = newVirtualRouter
		} else {
			networkCurrent.Spec.VirtualRouter = "default" // default value
		}
	}

	if d.HasChange("ingress_security_policy") {
		if val, ok := d.GetOk("ingress_security_policy"); ok {
			newIngressPolicy := val.(string)
			networkCurrent.Spec.IngressSecurityPolicy = []interface{}{newIngressPolicy}
		} else {
			networkCurrent.Spec.IngressSecurityPolicy = nil
		}
	}

	if d.HasChange("egress_security_policy") {
		if val, ok := d.GetOk("egress_security_policy"); ok {
			newEgressPolicy := val.(string)
			networkCurrent.Spec.EgressSecurityPolicy = []interface{}{newEgressPolicy}
		} else {
			networkCurrent.Spec.EgressSecurityPolicy = nil
		}
	}

	if d.HasChange("connection_tracking_mode") {
		if val, ok := d.GetOk("connection_tracking_mode"); ok {
			newConnectionTrackingMode := val.(string)
			networkCurrent.Spec.ConnectionTracking = newConnectionTrackingMode
		} else {
			networkCurrent.Spec.ConnectionTracking = "inherit from vrf" // default value
		}
	}

	if d.HasChange("allow_session_reuse") {
		if val, ok := d.GetOk("allow_session_reuse"); ok {
			newAllowSessionReuse := val.(string)
			networkCurrent.Spec.AllowSessionReuse = newAllowSessionReuse
		} else {
			networkCurrent.Spec.AllowSessionReuse = "inherit from vrf" // default value
		}
	}

	if d.HasChange("ip_fragments_forwarding") {
		if val, ok := d.GetOk("ip_fragments_forwarding"); ok {
			newIpFragmentsForwarding := val.(string)
			networkCurrent.Spec.IpFragmentsForwarding = newIpFragmentsForwarding
		} else {
			networkCurrent.Spec.IpFragmentsForwarding = "inherit from vrf" // default value
		}
	}

	if d.HasChange("service_bypass") {
		networkCurrent.Spec.ServiceBypass = d.Get("service_bypass").(bool)
	}

	// Handle mirror session changes
	if d.HasChange("ingress_mirror_session") {
		if val, ok := d.GetOk("ingress_mirror_session"); ok {
			newIngressMirrorSession := val.(string)
			networkCurrent.Spec.IngressMirrorSession = []interface{}{newIngressMirrorSession}
		} else {
			networkCurrent.Spec.IngressMirrorSession = nil
		}
	}

	if d.HasChange("egress_mirror_session") {
		if val, ok := d.GetOk("egress_mirror_session"); ok {
			newEgressMirrorSession := val.(string)
			networkCurrent.Spec.EgressMirrorSession = []interface{}{newEgressMirrorSession}
		} else {
			networkCurrent.Spec.EgressMirrorSession = nil
		}
	}

	jsonBytes, err := json.Marshal(networkCurrent)
	if err != nil {
		return diag.FromErr(err)
	}

	reqUpdate, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return diag.FromErr(err)
	}

	reqUpdate.AddCookie(&http.Cookie{Name: "sid", Value: config.SID})

	respUpdate, err := client.Do(reqUpdate)
	if err != nil {
		return diag.FromErr(err)
	}
	defer respUpdate.Body.Close()

	if respUpdate.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(respUpdate.Body)
		errMsg := fmt.Sprintf("failed to update network: HTTP %d %s: %s", respUpdate.StatusCode, respUpdate.Status, bodyBytes)
		if isDebugEnabled() {
			log.Printf("[DEBUG] Network update failed with response: %s", bodyBytes)
		}
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "Network update failed",
				Detail:   errMsg,
			},
		}
	}

	if isDebugEnabled() {
		log.Printf("[DEBUG] Network updated successfully")
	}

	return resourceNetworkRead(ctx, d, m)
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client()

	// Construct the URL for the network based on its name

	url := config.Server + "/configs/network/v1/tenant/default/networks/" + d.Get("name").(string)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set SID cookie for authentication
	req.AddCookie(&http.Cookie{Name: "sid", Value: config.SID})

	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("failed to delete network: HTTP %s", resp.Status)
	}

	// Clear the resource ID as it's been deleted from the PSM server.
	d.SetId("")

	return nil
}
