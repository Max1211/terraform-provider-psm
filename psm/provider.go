package psm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type FabricName string
type SwitchName string

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"psm_network":              resourceNetwork(),
			"psm_workload":             resourceWorkload(),
			"psm_rules":                resourceRules(),
			"psm_vrf":                  resourceVRF(),
			"psm_ipcollection":         resourceIPCollection(),
			"psm_orchestrator":         resourceOrchestrator(),
			"psm_workloadgroup":        resourceWorkloadGroup(),
			"psm_flow_export_policy":   resourceFlowExportPolicy(),
			"psm_syslog_export_policy": resourceSyslogPolicy(),
			"psm_app":                  resourceApps(),
			"psm_cluster":              resourceCluster(),
			"psm_rule_profile":         resourceRuleProfile(),
			"psm_dss":                  resourceDistributedServiceCard(),
			"psm_pdt":                  resourcePolicyDistributionTarget(),
			"psm_uiglobalsettings":     resourcePSMUIGlobalSettings(),
			"psm_authpolicy":           resourceAuthnPolicy(),
			"psm_user":                 resourceUser(),
			"psm_user_role":            resourceRole(),
			"psm_role_binding":         resourceRoleBinding(),
			"psm_ipsec_policy":         resourceIPSecPolicy(),
			"psm_certificate":          resourceCertificate(),
			"psm_nat_policy":           resourceNATPolicy(),
			"psm_user_preferences":     resourcePSMUserPreferences(),
			"psm_hosts":                resourceHosts(),
			"psm_mirror_session":       resourceMirrorSession(),
		},
		Schema: map[string]*schema.Schema{
			"user": {
				Description: "The username for the PSM Server",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("API_USER", nil),
			},
			"password": {
				Description: "The users password for the PSM Server",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("API_PASSWORD", nil),
			},
			"server": {
				Description: "The PSM server IP address or URL",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("API_SERVER", nil),
			},
			"insecure": {
				Description: "Skip SSL certificate verification.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := &Config{
		User:     d.Get("user").(string),
		Password: d.Get("password").(string),
		Server:   d.Get("server").(string),
		Insecure: d.Get("insecure").(bool),
	}

	err := config.Authenticate()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return config, nil
}
