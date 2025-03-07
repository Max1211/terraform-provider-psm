package psm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func validateMirrorSessionName(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	// Name must be 2-64 characters, alphanumeric, hyphen, underscore
	nameRegex := regexp.MustCompile(`^[a-zA-Z0-9-_]{2,64}$`)
	if !nameRegex.MatchString(v) {
		errs = append(errs, fmt.Errorf("%q must be 2-64 characters, only alphanumeric, hyphen and underscore allowed", key))
	}
	return
}

func resourceMirrorSession() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMirrorSessionCreate,
		ReadContext:   resourceMirrorSessionRead,
		UpdateContext: resourceMirrorSessionUpdate,
		DeleteContext: resourceMirrorSessionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMirrorSessionImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateMirrorSessionName,
				Description:  "Name of the mirror session (2-64 characters, alphanumeric, hyphen, underscore)",
			},
			"span_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 1023),
				Description:  "SPAN ID (1-1023)",
			},
			"packet_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2048,
				ValidateFunc: validation.IntBetween(64, 2048),
				Description:  "Packet size (64-2048)",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the mirror session is disabled",
			},
			"policy_distribution_target": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "Policy distribution target",
			},
			"collector": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"erspan_type_3"}, false),
							Description:  "Type of collector (currently only 'erspan_type_3' is supported)",
						},
						"destination": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP address of the collector destination",
						},
						"virtual_router": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "default",
							Description: "Virtual router to use for the collector",
						},
					},
				},
			},
		},
	}
}

type MirrorSession struct {
	Kind       string `json:"kind,omitempty"`
	APIVersion string `json:"api-version,omitempty"`
	Meta       struct {
		Name            string `json:"name"`
		Tenant          string `json:"tenant"`
		Namespace       string `json:"namespace,omitempty"`
		GenerationID    string `json:"generation-id,omitempty"`
		ResourceVersion string `json:"resource-version,omitempty"`
		UUID            string `json:"uuid,omitempty"`
		SelfLink        string `json:"self-link,omitempty"`
	} `json:"meta"`
	Spec struct {
		PacketSize                int                    `json:"packet-size"`
		StartCondition            map[string]interface{} `json:"start-condition"`
		Collectors                []Collector            `json:"collectors"`
		MatchRules                interface{}            `json:"match-rules"`
		SpanID                    int                    `json:"span-id"`
		Disabled                  bool                   `json:"disabled"`
		PolicyDistributionTargets []string               `json:"policy-distribution-targets"`
	} `json:"spec"`
}

type Collector struct {
	Type         string       `json:"type"`
	ExportConfig ExportConfig `json:"export-config"`
}

type ExportConfig struct {
	Destination   string `json:"destination"`
	VirtualRouter string `json:"virtual-router"`
}

func resourceMirrorSessionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client()

	session := &MirrorSession{}
	session.Meta.Name = d.Get("name").(string)
	session.Meta.Tenant = "default"
	session.Meta.Namespace = "default"

	session.Spec.SpanID = d.Get("span_id").(int)
	session.Spec.PacketSize = d.Get("packet_size").(int)
	session.Spec.Disabled = d.Get("disabled").(bool)
	session.Spec.StartCondition = map[string]interface{}{}
	session.Spec.PolicyDistributionTargets = []string{d.Get("policy_distribution_target").(string)}

	// Add collectors
	collectors := d.Get("collector").([]interface{})
	for _, c := range collectors {
		collector := c.(map[string]interface{})
		newCollector := Collector{
			Type: collector["type"].(string),
			ExportConfig: ExportConfig{
				Destination:   collector["destination"].(string),
				VirtualRouter: collector["virtual_router"].(string),
			},
		}
		session.Spec.Collectors = append(session.Spec.Collectors, newCollector)
	}

	jsonBytes, err := json.Marshal(session)
	if err != nil {
		log.Printf("[ERROR] Error marshalling Mirror Session: %s", err)
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Creating Mirror Session with name: %s", session.Meta.Name)
	log.Printf("[DEBUG] Request JSON: %s", jsonBytes)

	req, err := http.NewRequestWithContext(ctx, "POST", config.Server+"/configs/monitoring/v1/tenant/default/MirrorSession", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return diag.FromErr(err)
	}

	req.AddCookie(&http.Cookie{Name: "sid", Value: config.SID})

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Error when creating Mirror Session: %s", err)
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("failed to create Mirror Session: HTTP %d %s: %s", resp.StatusCode, resp.Status, bodyBytes)
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "Mirror Session creation failed",
				Detail:   errMsg,
			},
		}
	}

	responseBody := &MirrorSession{}
	if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(responseBody.Meta.UUID)
	log.Printf("[DEBUG] Mirror Session created with UUID: %s", responseBody.Meta.UUID)

	return resourceMirrorSessionRead(ctx, d, m)
}

func resourceMirrorSessionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client()

	sessionName := d.Get("name").(string)
	req, err := http.NewRequestWithContext(ctx, "GET", config.Server+"/configs/monitoring/v1/tenant/default/MirrorSession/"+sessionName, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.AddCookie(&http.Cookie{Name: "sid", Value: config.SID})

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Error when reading Mirror Session: %s", err)
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		bodyBytes, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("failed to read Mirror Session: HTTP %d %s: %s", resp.StatusCode, resp.Status, bodyBytes)
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "Mirror Session read failed",
				Detail:   errMsg,
			},
		}
	}

	responseBody := &MirrorSession{}
	if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", responseBody.Meta.Name)
	d.Set("span_id", responseBody.Spec.SpanID)
	d.Set("packet_size", responseBody.Spec.PacketSize)
	d.Set("disabled", responseBody.Spec.Disabled)

	if len(responseBody.Spec.PolicyDistributionTargets) > 0 {
		d.Set("policy_distribution_target", responseBody.Spec.PolicyDistributionTargets[0])
	}

	collectors := make([]map[string]interface{}, len(responseBody.Spec.Collectors))
	for i, collector := range responseBody.Spec.Collectors {
		collectors[i] = map[string]interface{}{
			"type":           collector.Type,
			"destination":    collector.ExportConfig.Destination,
			"virtual_router": collector.ExportConfig.VirtualRouter,
		}
	}
	d.Set("collector", collectors)

	log.Printf("[DEBUG] Mirror Session read with UUID: %s", responseBody.Meta.UUID)

	return nil
}

func resourceMirrorSessionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client()

	session := &MirrorSession{}
	session.Meta.Name = d.Get("name").(string)
	session.Meta.Tenant = "default"
	session.Meta.Namespace = "default"

	session.Spec.SpanID = d.Get("span_id").(int)
	session.Spec.PacketSize = d.Get("packet_size").(int)
	session.Spec.Disabled = d.Get("disabled").(bool)
	session.Spec.StartCondition = map[string]interface{}{}
	session.Spec.PolicyDistributionTargets = []string{d.Get("policy_distribution_target").(string)}

	// Add collectors
	collectors := d.Get("collector").([]interface{})
	for _, c := range collectors {
		collector := c.(map[string]interface{})
		newCollector := Collector{
			Type: collector["type"].(string),
			ExportConfig: ExportConfig{
				Destination:   collector["destination"].(string),
				VirtualRouter: collector["virtual_router"].(string),
			},
		}
		session.Spec.Collectors = append(session.Spec.Collectors, newCollector)
	}

	jsonBytes, err := json.Marshal(session)
	if err != nil {
		log.Printf("[ERROR] Error marshalling Mirror Session: %s", err)
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Updating Mirror Session with name: %s", session.Meta.Name)
	log.Printf("[DEBUG] Request JSON: %s", jsonBytes)

	sessionName := d.Get("name").(string)
	req, err := http.NewRequestWithContext(ctx, "PUT", config.Server+"/configs/monitoring/v1/tenant/default/MirrorSession/"+sessionName, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return diag.FromErr(err)
	}

	req.AddCookie(&http.Cookie{Name: "sid", Value: config.SID})

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Error when updating Mirror Session: %s", err)
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("failed to update Mirror Session: HTTP %d %s: %s", resp.StatusCode, resp.Status, bodyBytes)
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "Mirror Session update failed",
				Detail:   errMsg,
			},
		}
	}

	responseBody := &MirrorSession{}
	if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Mirror Session updated with UUID: %s", responseBody.Meta.UUID)

	return resourceMirrorSessionRead(ctx, d, m)
}

func resourceMirrorSessionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client()

	sessionName := d.Get("name").(string)
	req, err := http.NewRequestWithContext(ctx, "DELETE", config.Server+"/configs/monitoring/v1/tenant/default/MirrorSession/"+sessionName, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.AddCookie(&http.Cookie{Name: "sid", Value: config.SID})

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Error when deleting Mirror Session: %s", err)
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("failed to delete Mirror Session: HTTP %d %s: %s", resp.StatusCode, resp.Status, bodyBytes)
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "Mirror Session deletion failed",
				Detail:   errMsg,
			},
		}
	}

	log.Printf("[DEBUG] Mirror Session deleted: %s", sessionName)
	d.SetId("")

	return nil
}

func resourceMirrorSessionImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	config := m.(*Config)
	client := config.Client()

	// The ID is expected to be the name of the Mirror Session
	sessionName := d.Id()

	req, err := http.NewRequestWithContext(ctx, "GET", config.Server+"/configs/monitoring/v1/tenant/default/MirrorSession/"+sessionName, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err)
	}

	req.AddCookie(&http.Cookie{Name: "sid", Value: config.SID})

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error reading Mirror Session: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to read Mirror Session: HTTP %d %s: %s", resp.StatusCode, resp.Status, bodyBytes)
	}

	responseBody := &MirrorSession{}
	if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
		return nil, fmt.Errorf("error decoding response: %s", err)
	}

	d.SetId(responseBody.Meta.UUID)
	d.Set("name", responseBody.Meta.Name)
	d.Set("span_id", responseBody.Spec.SpanID)
	d.Set("packet_size", responseBody.Spec.PacketSize)
	d.Set("disabled", responseBody.Spec.Disabled)

	if len(responseBody.Spec.PolicyDistributionTargets) > 0 {
		d.Set("policy_distribution_target", responseBody.Spec.PolicyDistributionTargets[0])
	}

	collectors := make([]map[string]interface{}, len(responseBody.Spec.Collectors))
	for i, collector := range responseBody.Spec.Collectors {
		collectors[i] = map[string]interface{}{
			"type":           collector.Type,
			"destination":    collector.ExportConfig.Destination,
			"virtual_router": collector.ExportConfig.VirtualRouter,
		}
	}
	d.Set("collector", collectors)

	return []*schema.ResourceData{d}, nil
}
