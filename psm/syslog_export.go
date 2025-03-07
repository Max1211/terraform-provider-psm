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

func resourceSyslogPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSyslogPolicyCreate,
		ReadContext:   resourceSyslogPolicyRead,
		UpdateContext: resourceSyslogPolicyUpdate,
		DeleteContext: resourceSyslogPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSyslogPolicyImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"format": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"filter": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"syslogconfig": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"facility": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
						"disable_batching": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: false,
						},
					},
				},
			},
			"psm_target": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"targets": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": {
							Type:     schema.TypeString,
							Required: true,
						},
						"transport": {
							Type:     schema.TypeString,
							Required: true,
						},
						"trusted_certs": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the trusted certificate for TLS connections",
						},
						"client_certificate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the client certificate for TLS connections",
						},
						"hostname_verification": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Hostname to verify against the server certificate",
						},
						"skip_cert_verification": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether to skip server certificate verification",
						},
					},
				},
			},
		},
	}
}

type SyslogPolicy struct {
	Meta struct {
		Name   string `json:"name"`
		Tenant string `json:"tenant"`
		UUID   string `json:"uuid"`
	} `json:"meta"`
	Spec struct {
		Format    string       `json:"format"`
		Filter    []string     `json:"filter"`
		Config    SyslogConfig `json:"config"`
		PsmTarget PsmTarget    `json:"psm-target"`
		Targets   []Target     `json:"targets"`
	} `json:"spec"`
}

type SyslogConfig struct {
	FacilityOverride string `json:"facility-override"`
	DisableBatching  bool   `json:"disable-batching,omitempty"`
}

type PsmTarget struct {
	Enable bool `json:"enable"`
}

type ServerCertVerification struct {
	Hostname                   string `json:"hostname,omitempty"`
	SkipServerCertVerification *bool  `json:"skip-server-cert-verification,omitempty"`
}

type Target struct {
	Destination                      string                  `json:"destination"`
	Transport                        string                  `json:"transport"`
	TrustedCerts                     string                  `json:"trusted-certs,omitempty"`
	ClientCertificate                string                  `json:"client-certificate,omitempty"`
	ServerCertificateVerificationOpt *ServerCertVerification `json:"server-certificate-verification-options,omitempty"`
}

func resourceSyslogPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client()

	syslogPolicy := &SyslogPolicy{}
	syslogPolicy.Meta.Name = d.Get("name").(string)
	syslogPolicy.Meta.Tenant = "default"
	syslogPolicy.Spec.Format = d.Get("format").(string)

	psmTarget := d.Get("psm_target").([]interface{})[0].(map[string]interface{})
	syslogPolicy.Spec.PsmTarget = PsmTarget{
		Enable: psmTarget["enable"].(bool),
	}

	filterInterface := d.Get("filter").([]interface{})
	filterStrings := make([]string, len(filterInterface))
	for i, v := range filterInterface {
		if v != nil {
			filterStrings[i] = v.(string)
		}
	}

	syslogConfig := d.Get("syslogconfig").([]interface{})[0].(map[string]interface{})
	syslogPolicy.Spec.Config = SyslogConfig{
		FacilityOverride: syslogConfig["facility"].(string),
		DisableBatching:  syslogConfig["disable_batching"].(bool),
	}

	syslogPolicy.Spec.Filter = filterStrings

	targets := d.Get("targets").([]interface{})
	for _, t := range targets {
		target := t.(map[string]interface{})
		newTarget := Target{
			Destination: target["destination"].(string),
			Transport:   target["transport"].(string),
		}

		// Add certificate fields if provided
		if trustedCerts, ok := target["trusted_certs"].(string); ok && trustedCerts != "" {
			newTarget.TrustedCerts = trustedCerts
		}

		if clientCert, ok := target["client_certificate"].(string); ok && clientCert != "" {
			newTarget.ClientCertificate = clientCert
		}

		// Check if we need to set up certificate verification options
		hostnameProvided := false
		skipVerificationProvided := false
		var hostname string
		var skipVerification bool

		if h, ok := target["hostname_verification"].(string); ok && h != "" {
			hostnameProvided = true
			hostname = h
		}

		if skip, ok := target["skip_cert_verification"].(bool); ok {
			skipVerificationProvided = true
			skipVerification = skip
		}

		// If either option is provided, create the verification options
		if hostnameProvided || skipVerificationProvided {
			verificationOpt := &ServerCertVerification{}

			if hostnameProvided {
				verificationOpt.Hostname = hostname
			}

			if skipVerificationProvided && skipVerification {
				verificationOpt.SkipServerCertVerification = &skipVerification
			}

			newTarget.ServerCertificateVerificationOpt = verificationOpt
		}

		syslogPolicy.Spec.Targets = append(syslogPolicy.Spec.Targets, newTarget)
	}

	jsonBytes, err := json.Marshal(syslogPolicy)
	if err != nil {
		log.Printf("[ERROR] Error marshalling Syslog Policy: %s", err)
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Creating Syslog Policy with name: %s", syslogPolicy.Meta.Name)
	log.Printf("[DEBUG] Request JSON: %s", jsonBytes)

	req, err := http.NewRequestWithContext(ctx, "POST", config.Server+"/configs/monitoring/v1/tenant/default/fwlogPolicy", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return diag.FromErr(err)
	}

	req.AddCookie(&http.Cookie{Name: "sid", Value: config.SID})

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Error when creating Syslog Policy: %s", err)
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("failed to create Syslog Policy: HTTP %d %s: %s", resp.StatusCode, resp.Status, bodyBytes)
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "Syslog Policy creation failed",
				Detail:   errMsg,
			},
		}
	}

	responseBody := &SyslogPolicy{}
	if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(responseBody.Meta.UUID)

	log.Printf("[DEBUG] Syslog Policy created with UUID: %s", responseBody.Meta.UUID)

	return resourceSyslogPolicyRead(ctx, d, m)
}

func resourceSyslogPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client()

	fwlogPolicyName := d.Get("name").(string)
	req, err := http.NewRequestWithContext(ctx, "GET", config.Server+"/configs/monitoring/v1/tenant/default/fwlogPolicy/"+fwlogPolicyName, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.AddCookie(&http.Cookie{Name: "sid", Value: config.SID})

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Error when reading Syslog Policy: %s", err)
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("failed to read Syslog Policy: HTTP %d %s: %s", resp.StatusCode, resp.Status, bodyBytes)
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "Syslog Policy read failed",
				Detail:   errMsg,
			},
		}
	}

	responseBody := &SyslogPolicy{}
	if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", responseBody.Meta.Name)
	d.Set("format", responseBody.Spec.Format)
	d.Set("filter", responseBody.Spec.Filter)
	d.Set("syslogconfig", []interface{}{map[string]interface{}{
		"facility":         responseBody.Spec.Config.FacilityOverride,
		"disable_batching": responseBody.Spec.Config.DisableBatching,
	}})
	d.Set("psm_target", []interface{}{map[string]interface{}{
		"enable": responseBody.Spec.PsmTarget.Enable,
	}})

	targets := make([]map[string]interface{}, len(responseBody.Spec.Targets))
	for i, target := range responseBody.Spec.Targets {
		targetMap := map[string]interface{}{
			"destination": target.Destination,
			"transport":   target.Transport,
		}

		// Handle certificate fields
		if target.TrustedCerts != "" {
			targetMap["trusted_certs"] = target.TrustedCerts
		}

		if target.ClientCertificate != "" {
			targetMap["client_certificate"] = target.ClientCertificate
		}

		if target.ServerCertificateVerificationOpt != nil {
			if target.ServerCertificateVerificationOpt.Hostname != "" {
				targetMap["hostname_verification"] = target.ServerCertificateVerificationOpt.Hostname
			}

			if target.ServerCertificateVerificationOpt.SkipServerCertVerification != nil {
				targetMap["skip_cert_verification"] = *target.ServerCertificateVerificationOpt.SkipServerCertVerification
			}
		}

		targets[i] = targetMap
	}
	d.Set("targets", targets)

	log.Printf("[DEBUG] Syslog Policy read with UUID: %s", responseBody.Meta.UUID)

	return nil
}

func resourceSyslogPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client()

	syslogPolicy := &SyslogPolicy{}
	syslogPolicy.Meta.Name = d.Get("name").(string)
	syslogPolicy.Meta.Tenant = "default"
	syslogPolicy.Spec.Format = d.Get("format").(string)

	psmTarget := d.Get("psm_target").([]interface{})[0].(map[string]interface{})
	syslogPolicy.Spec.PsmTarget = PsmTarget{
		Enable: psmTarget["enable"].(bool),
	}

	filterInterface := d.Get("filter").([]interface{})
	filterStrings := make([]string, len(filterInterface))
	for i, v := range filterInterface {
		if v != nil {
			filterStrings[i] = v.(string)
		}
	}

	syslogConfig := d.Get("syslogconfig").([]interface{})[0].(map[string]interface{})
	syslogPolicy.Spec.Config = SyslogConfig{
		FacilityOverride: syslogConfig["facility"].(string),
		DisableBatching:  syslogConfig["disable_batching"].(bool),
	}

	syslogPolicy.Spec.Filter = filterStrings

	targets := d.Get("targets").([]interface{})
	for _, t := range targets {
		target := t.(map[string]interface{})
		newTarget := Target{
			Destination: target["destination"].(string),
			Transport:   target["transport"].(string),
		}

		// Add certificate fields if provided
		if trustedCerts, ok := target["trusted_certs"].(string); ok && trustedCerts != "" {
			newTarget.TrustedCerts = trustedCerts
		}

		if clientCert, ok := target["client_certificate"].(string); ok && clientCert != "" {
			newTarget.ClientCertificate = clientCert
		}

		if hostname, ok := target["hostname_verification"].(string); ok && hostname != "" {
			newTarget.ServerCertificateVerificationOpt = &ServerCertVerification{
				Hostname: hostname,
			}
		}

		syslogPolicy.Spec.Targets = append(syslogPolicy.Spec.Targets, newTarget)
	}

	jsonBytes, err := json.Marshal(syslogPolicy)
	if err != nil {
		log.Printf("[ERROR] Error marshalling Syslog Policy: %s", err)
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Updating Syslog Policy with name: %s", syslogPolicy.Meta.Name)
	log.Printf("[DEBUG] Request JSON: %s", jsonBytes)

	fwlogPolicyName := d.Get("name").(string)
	req, err := http.NewRequestWithContext(ctx, "PUT", config.Server+"/configs/monitoring/v1/tenant/default/fwlogPolicy/"+fwlogPolicyName, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return diag.FromErr(err)
	}

	req.AddCookie(&http.Cookie{Name: "sid", Value: config.SID})

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Error when updating Syslog Policy: %s", err)
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("failed to update Syslog Policy: HTTP %d %s: %s", resp.StatusCode, resp.Status, bodyBytes)
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "Syslog Policy update failed",
				Detail:   errMsg,
			},
		}
	}

	responseBody := &SyslogPolicy{}
	if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(responseBody.Meta.UUID)

	log.Printf("[DEBUG] Syslog Policy updated with UUID: %s", responseBody.Meta.UUID)

	return resourceSyslogPolicyRead(ctx, d, m)
}

func resourceSyslogPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client()

	fwlogPolicyName := d.Get("name").(string)
	req, err := http.NewRequestWithContext(ctx, "DELETE", config.Server+"/configs/monitoring/v1/tenant/default/fwlogPolicy/"+fwlogPolicyName, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	req.AddCookie(&http.Cookie{Name: "sid", Value: config.SID})

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Error when deleting Syslog Policy: %s", err)
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		errMsg := fmt.Sprintf("failed to delete Syslog Policy: HTTP %d %s: %s", resp.StatusCode, resp.Status, bodyBytes)
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "Syslog Policy deletion failed",
				Detail:   errMsg,
			},
		}
	}

	log.Printf("[DEBUG] Syslog Policy deleted with UUID: %s", d.Id())

	d.SetId("")

	return nil
}

func resourceSyslogPolicyImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	config := m.(*Config)
	client := config.Client()

	// The ID is expected to be the name of the Syslog Policy
	fwlogPolicyName := d.Id()

	req, err := http.NewRequestWithContext(ctx, "GET", config.Server+"/configs/monitoring/v1/tenant/default/fwlogPolicy/"+fwlogPolicyName, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err)
	}

	req.AddCookie(&http.Cookie{Name: "sid", Value: config.SID})

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error reading Syslog Policy: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to read Syslog Policy: HTTP %d %s: %s", resp.StatusCode, resp.Status, bodyBytes)
	}

	responseBody := &SyslogPolicy{}
	if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
		return nil, fmt.Errorf("error decoding response: %s", err)
	}

	d.SetId(responseBody.Meta.UUID)
	d.Set("name", responseBody.Meta.Name)
	d.Set("format", responseBody.Spec.Format)
	d.Set("filter", responseBody.Spec.Filter)

	syslogConfig := []map[string]interface{}{
		{
			"facility":         responseBody.Spec.Config.FacilityOverride,
			"disable_batching": responseBody.Spec.Config.DisableBatching,
		},
	}
	d.Set("syslogconfig", syslogConfig)

	psmTarget := []map[string]interface{}{
		{
			"enable": responseBody.Spec.PsmTarget.Enable,
		},
	}
	d.Set("psm_target", psmTarget)

	targets := make([]map[string]interface{}, len(responseBody.Spec.Targets))
	for i, target := range responseBody.Spec.Targets {
		targetMap := map[string]interface{}{
			"destination": target.Destination,
			"transport":   target.Transport,
		}

		// Handle certificate fields
		if target.TrustedCerts != "" {
			targetMap["trusted_certs"] = target.TrustedCerts
		}

		if target.ClientCertificate != "" {
			targetMap["client_certificate"] = target.ClientCertificate
		}

		if target.ServerCertificateVerificationOpt != nil {
			if target.ServerCertificateVerificationOpt.Hostname != "" {
				targetMap["hostname_verification"] = target.ServerCertificateVerificationOpt.Hostname
			}

			if target.ServerCertificateVerificationOpt.SkipServerCertVerification != nil {
				targetMap["skip_cert_verification"] = *target.ServerCertificateVerificationOpt.SkipServerCertVerification
			}
		}

		targets[i] = targetMap
	}
	d.Set("targets", targets)

	return []*schema.ResourceData{d}, nil
}
