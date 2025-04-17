package psm

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// convertToStringSlice converts various input types to a slice of strings
func convertToStringSlice(input interface{}) []string {
	var result []string
	inputSlice, ok := input.([]interface{})
	if !ok {
		return result
	}
	for _, v := range inputSlice {
		str, ok := v.(string)
		if !ok {
			continue
		}
		result = append(result, str)
	}
	return result
}

func getStringOrEmpty(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func getBoolOrDefault(m map[string]interface{}, key string, defaultValue bool) bool {
	if v, ok := m[key]; ok {
		b, ok := v.(bool)
		if ok {
			return b
		}
	}
	return defaultValue
}

func getStringSlice(m map[string]interface{}, key string) []string {
	if v, ok := m[key]; ok {
		if slice, ok := v.([]interface{}); ok {
			result := make([]string, len(slice))
			for i, item := range slice {
				result[i] = fmt.Sprintf("%v", item)
			}
			return result
		}
	}
	return nil
}

func expandStringList(list []interface{}) []string {
	result := make([]string, len(list))
	for i, v := range list {
		result[i] = v.(string)
	}
	return result
}

func ExpandStringSet(set *schema.Set) []string {
	list := set.List()
	result := make([]string, len(list))
	for i, v := range list {
		result[i] = v.(string)
	}
	return result
}

func FlattenStringList(list []string) *schema.Set {
	set := schema.NewSet(schema.HashString, []interface{}{})
	for _, v := range list {
		set.Add(v)
	}
	return set
}

func expandSpec(d *schema.ResourceData) TunnelSpec {
	tunnelData := d.Get("tunnel").([]interface{})[0].(map[string]interface{})

	spec := TunnelSpec{
		HAMode:                    getStringOrEmpty(tunnelData, "ha_mode"),
		TunnelEndpoints:           expandTunnelEndpoints(tunnelData["tunnel_endpoints"].([]interface{})),
		PolicyDistributionTargets: convertToStringSlice(tunnelData["policy_distribution_targets"].([]interface{})),
		DisableTCPMSSAdjust:       getBoolOrDefault(tunnelData, "disable_tcp_mss_adjust", false),
	}

	return spec
}

func expandTunnelEndpoints(endpoints []interface{}) []TunnelEndpoint {
	var tunnelEndpoints []TunnelEndpoint

	for _, endpoint := range endpoints {
		endpointData := endpoint.(map[string]interface{})
		tunnelEndpoint := TunnelEndpoint{
			InterfaceName: getStringOrEmpty(endpointData, "interface_name"),
			DSE:           getStringOrEmpty(endpointData, "dse"),
			IKEVersion:    getStringOrEmpty(endpointData, "ike_version"),
		}

		if ikeSA, ok := endpointData["ike_sa"].([]interface{}); ok && len(ikeSA) > 0 {
			tunnelEndpoint.IKESA = expandIKESA(ikeSA[0].(map[string]interface{}))
		}

		if ipsecSA, ok := endpointData["ipsec_sa"].([]interface{}); ok && len(ipsecSA) > 0 {
			tunnelEndpoint.IPSECSA = expandIPSECSA(ipsecSA[0].(map[string]interface{}))
		}

		if localID, ok := endpointData["local_identifier"].([]interface{}); ok && len(localID) > 0 {
			tunnelEndpoint.LocalIdentifier = expandIdentifier(localID[0].(map[string]interface{}))
		}

		if remoteID, ok := endpointData["remote_identifier"].([]interface{}); ok && len(remoteID) > 0 {
			tunnelEndpoint.RemoteIdentifier = expandIdentifier(remoteID[0].(map[string]interface{}))
		}

		// Expand lifetime for this endpoint
		if lifetime, ok := endpointData["lifetime"].([]interface{}); ok && len(lifetime) > 0 {
			lifetimeData := lifetime[0].(map[string]interface{})
			tunnelEndpoint.Lifetime = &Lifetime{
				SALifetime:  getStringOrEmpty(lifetimeData, "sa_lifetime"),
				IKELifetime: getStringOrEmpty(lifetimeData, "ike_lifetime"),
			}
		}

		tunnelEndpoints = append(tunnelEndpoints, tunnelEndpoint)
	}

	return tunnelEndpoints
}

func expandIKESA(data map[string]interface{}) *IKESA {
	ikesa := &IKESA{
		EncryptionAlgorithms: convertToStringSlice(data["encryption_algorithms"].([]interface{})),
		HashAlgorithms:       convertToStringSlice(data["hash_algorithms"].([]interface{})),
		DHGroups:             convertToStringSlice(data["dh_groups"].([]interface{})),
		RekeyLifetime:        getStringOrEmpty(data, "rekey_lifetime"),
		PreSharedKey:         getStringOrEmpty(data, "pre_shared_key"),
		ReauthLifetime:       getStringOrEmpty(data, "reauth_lifetime"),
		DPDDelay:             getStringOrEmpty(data, "dpd_delay"),
		IKEV1DPDTimeout:      getStringOrEmpty(data, "ikev1_dpd_timeout"),
		IKEInitiator:         getBoolOrDefault(data, "ike_initiator", true),
		AuthType:             getStringOrEmpty(data, "auth_type"),
	}

	if ikesa.AuthType == "certificates" {
		ikesa.LocalIdentityCertificates = getStringOrEmpty(data, "local_identity_certificates")
		if certs, ok := data["remote_ca_certificates"].([]interface{}); ok {
			ikesa.RemoteCACertificates = convertToStringSlice(certs)
		}
	}

	return ikesa
}

func expandIPSECSA(data map[string]interface{}) *IPSECSA {
	return &IPSECSA{
		EncryptionAlgorithms: convertToStringSlice(data["encryption_algorithms"].([]interface{})),
		DHGroups:             convertToStringSlice(data["dh_groups"].([]interface{})),
		RekeyLifetime:        getStringOrEmpty(data, "rekey_lifetime"),
	}
}

func expandIdentifier(data map[string]interface{}) Identifier {
	return Identifier{
		Type:  getStringOrEmpty(data, "type"),
		Value: getStringOrEmpty(data, "value"),
	}
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func flattenSpec(spec *TunnelSpec, d *schema.ResourceData) error {
	if spec == nil {
		return fmt.Errorf("TunnelSpec is nil")
	}

	tunnelMap := make(map[string]interface{})
	tunnelMap["ha_mode"] = spec.HAMode
	tunnelMap["policy_distribution_targets"] = spec.PolicyDistributionTargets
	tunnelMap["disable_tcp_mss_adjust"] = spec.DisableTCPMSSAdjust

	tunnelEndpoints := flattenTunnelEndpoints(spec.TunnelEndpoints, d)
	tunnelMap["tunnel_endpoints"] = tunnelEndpoints

	if spec.Config != nil {
		lifetime := make(map[string]interface{})

		// Preserve existing values if new values are empty
		oldTunnel, ok := d.GetOk("tunnel")
		if ok && len(oldTunnel.([]interface{})) > 0 {
			oldTunnelMap := oldTunnel.([]interface{})[0].(map[string]interface{})
			if oldLifetime, ok := oldTunnelMap["lifetime"].([]interface{}); ok && len(oldLifetime) > 0 {
				oldLifetimeMap := oldLifetime[0].(map[string]interface{})

				if spec.Config.SALifetime != "" {
					lifetime["sa_lifetime"] = spec.Config.SALifetime
				} else if saLifetime, ok := oldLifetimeMap["sa_lifetime"]; ok {
					lifetime["sa_lifetime"] = saLifetime
				}

				if spec.Config.IKELifetime != "" {
					lifetime["ike_lifetime"] = spec.Config.IKELifetime
				} else if ikeLifetime, ok := oldLifetimeMap["ike_lifetime"]; ok {
					lifetime["ike_lifetime"] = ikeLifetime
				}
			}
		} else {
			// If there are no old values, only set the new values if they're not empty
			if spec.Config.SALifetime != "" {
				lifetime["sa_lifetime"] = spec.Config.SALifetime
			}
			if spec.Config.IKELifetime != "" {
				lifetime["ike_lifetime"] = spec.Config.IKELifetime
			}
		}

		// Only add the lifetime block if it's not empty
		if len(lifetime) > 0 {
			tunnelMap["lifetime"] = []interface{}{lifetime}
		}
	}

	if err := d.Set("tunnel", []interface{}{tunnelMap}); err != nil {
		return fmt.Errorf("error setting tunnel: %s", err)
	}

	return nil
}

func flattenTunnelEndpoints(endpoints []TunnelEndpoint, d *schema.ResourceData) []interface{} {
	tunnelEndpoints := make([]interface{}, len(endpoints))
	for i, endpoint := range endpoints {
		endpointMap := make(map[string]interface{})
		endpointMap["interface_name"] = endpoint.InterfaceName
		endpointMap["dse"] = endpoint.DSE
		endpointMap["ike_version"] = endpoint.IKEVersion

		if endpoint.IKESA != nil {
			ikesa := make(map[string]interface{})
			ikesa["encryption_algorithms"] = endpoint.IKESA.EncryptionAlgorithms
			ikesa["hash_algorithms"] = endpoint.IKESA.HashAlgorithms
			ikesa["dh_groups"] = endpoint.IKESA.DHGroups
			ikesa["rekey_lifetime"] = endpoint.IKESA.RekeyLifetime
			ikesa["reauth_lifetime"] = endpoint.IKESA.ReauthLifetime
			ikesa["dpd_delay"] = endpoint.IKESA.DPDDelay
			ikesa["ikev1_dpd_timeout"] = endpoint.IKESA.IKEV1DPDTimeout
			ikesa["ike_initiator"] = endpoint.IKESA.IKEInitiator
			ikesa["auth_type"] = endpoint.IKESA.AuthType
			if endpoint.IKESA.LocalIdentityCertificates != "" {
				ikesa["local_identity_certificates"] = endpoint.IKESA.LocalIdentityCertificates
			}
			if len(endpoint.IKESA.RemoteCACertificates) > 0 {
				ikesa["remote_ca_certificates"] = endpoint.IKESA.RemoteCACertificates
			}

			// Preserve pre_shared_key if it exists in the current state
			if v, ok := d.GetOk(fmt.Sprintf("tunnel.0.tunnel_endpoints.%d.ike_sa.0.pre_shared_key", i)); ok {
				ikesa["pre_shared_key"] = v.(string)
			}

			endpointMap["ike_sa"] = []interface{}{ikesa}
		}

		if endpoint.IPSECSA != nil {
			ipsecsa := make(map[string]interface{})
			ipsecsa["encryption_algorithms"] = endpoint.IPSECSA.EncryptionAlgorithms
			ipsecsa["dh_groups"] = endpoint.IPSECSA.DHGroups
			ipsecsa["rekey_lifetime"] = endpoint.IPSECSA.RekeyLifetime
			endpointMap["ipsec_sa"] = []interface{}{ipsecsa}
		}

		endpointMap["local_identifier"] = []interface{}{
			map[string]interface{}{
				"type":  endpoint.LocalIdentifier.Type,
				"value": endpoint.LocalIdentifier.Value,
			},
		}

		endpointMap["remote_identifier"] = []interface{}{
			map[string]interface{}{
				"type":  endpoint.RemoteIdentifier.Type,
				"value": endpoint.RemoteIdentifier.Value,
			},
		}

		// Add lifetime configuration if it exists
		if endpoint.Lifetime != nil {
			endpointMap["lifetime"] = []interface{}{
				map[string]interface{}{
					"sa_lifetime":  endpoint.Lifetime.SALifetime,
					"ike_lifetime": endpoint.Lifetime.IKELifetime,
				},
			}
		} else {
			// Check if lifetime exists in the current state and preserve it
			if v, ok := d.GetOk(fmt.Sprintf("tunnel.0.tunnel_endpoints.%d.lifetime", i)); ok {
				endpointMap["lifetime"] = v.([]interface{})
			}
		}

		tunnelEndpoints[i] = endpointMap
	}

	return tunnelEndpoints
}

func suppressMissingLifetimeValues(k, old, new string, d *schema.ResourceData) bool {
	// If the new value is empty, suppress the diff
	return new == ""
}

func validateNATRules(rules []interface{}) error {
	for i, r := range rules {
		rule := r.(map[string]interface{})

		// Validate name
		name, ok := rule["name"].(string)
		if !ok || name == "" {
			return fmt.Errorf("rule %d: name is required", i)
		}

		// Validate type
		ruleType, ok := rule["type"].(string)
		if !ok || ruleType == "" {
			return fmt.Errorf("rule %d: type is required", i)
		}

		// Source is optional, no validation needed

		// Destination is optional, no validation needed

		// Validate destination_proto_port (optional)
		if destProtoPort, ok := rule["destination_proto_port"].([]interface{}); ok && len(destProtoPort) > 0 {
			dpp := destProtoPort[0].(map[string]interface{})
			protocol, ok := dpp["protocol"].(string)
			if !ok || protocol == "" {
				return fmt.Errorf("rule %d: protocol in destination_proto_port is required if destination_proto_port is specified", i)
			}
			// Ports are optional, no need to validate
		}

		// Validate that at least one of translated_source or translated_destination is specified
		hasTranslatedSource := rule["translated_source"] != nil && len(rule["translated_source"].([]interface{})) > 0
		hasTranslatedDest := rule["translated_destination"] != nil && len(rule["translated_destination"].([]interface{})) > 0

		if !hasTranslatedSource && !hasTranslatedDest {
			return fmt.Errorf("rule %d: either translated_source or translated_destination must be specified", i)
		}

		// Validate translated_destination_port (optional)
		if translatedDestPort, ok := rule["translated_destination_port"].(string); ok && translatedDestPort != "" {
			if !hasTranslatedDest {
				return fmt.Errorf("rule %d: translated_destination is required when translated_destination_port is specified", i)
			}
		}
	}
	return nil
}

func createNatRule(rule map[string]interface{}) NatRule {
	natRule := NatRule{
		Name:    rule["name"].(string),
		Disable: rule["disable"].(bool),
		Type:    rule["type"].(string),
	}

	// Handle Source (optional)
	if source, ok := rule["source"].([]interface{}); ok && len(source) > 0 {
		sourceMap := source[0].(map[string]interface{})
		src := createAddressCollection(sourceMap)
		if len(src.Addresses) > 0 || len(src.IPCollections) > 0 {
			natRule.Source = &src
		}
	}

	// Handle Destination (optional)
	if destination, ok := rule["destination"].([]interface{}); ok && len(destination) > 0 {
		destMap := destination[0].(map[string]interface{})
		dest := createAddressCollection(destMap)
		if len(dest.Addresses) > 0 || len(dest.IPCollections) > 0 {
			natRule.Destination = &dest
		}
	}

	// Handle DestinationProtoPort (optional)
	if protoPort, ok := rule["destination_proto_port"].([]interface{}); ok && len(protoPort) > 0 {
		pp := protoPort[0].(map[string]interface{})
		natRule.DestinationProtoPort.Protocol = pp["protocol"].(string)
		if ports, ok := pp["ports"]; ok {
			natRule.DestinationProtoPort.Ports = ports.(string)
		}
	}

	// Handle TranslatedSource (optional)
	if translatedSource, ok := rule["translated_source"].([]interface{}); ok && len(translatedSource) > 0 {
		tsMap := translatedSource[0].(map[string]interface{})
		ts := createAddressCollection(tsMap)
		if len(ts.Addresses) > 0 || len(ts.IPCollections) > 0 {
			natRule.TranslatedSource = &ts
		}
	}

	// Handle TranslatedDestination (optional)
	if translatedDest, ok := rule["translated_destination"].([]interface{}); ok && len(translatedDest) > 0 {
		tdMap := translatedDest[0].(map[string]interface{})
		td := createAddressCollection(tdMap)
		if len(td.Addresses) > 0 || len(td.IPCollections) > 0 {
			natRule.TranslatedDestination = &td
		}
	}

	// Handle TranslatedDestinationPort (optional)
	if translatedDestPort, ok := rule["translated_destination_port"].(string); ok && translatedDestPort != "" {
		natRule.TranslatedDestinationPort = translatedDestPort
	}

	return natRule
}

// Helper function to create AddressCollection
func createAddressCollection(data map[string]interface{}) AddressCollection {
	return AddressCollection{
		Addresses:     expandStringList(data["addresses"].([]interface{})),
		IPCollections: expandStringList(data["ipcollections"].([]interface{})),
	}
}

func expandPermissions(permissions []interface{}) []struct {
	ResourceGroup     string   `json:"resource-group"`
	ResourceKind      string   `json:"resource-kind"`
	ResourceNamespace string   `json:"resource-namespace"`
	Actions           []string `json:"actions"`
} {
	result := make([]struct {
		ResourceGroup     string   `json:"resource-group"`
		ResourceKind      string   `json:"resource-kind"`
		ResourceNamespace string   `json:"resource-namespace"`
		Actions           []string `json:"actions"`
	}, len(permissions))

	for i, perm := range permissions {
		p := perm.(map[string]interface{})
		result[i] = struct {
			ResourceGroup     string   `json:"resource-group"`
			ResourceKind      string   `json:"resource-kind"`
			ResourceNamespace string   `json:"resource-namespace"`
			Actions           []string `json:"actions"`
		}{
			ResourceGroup:     p["resource_group"].(string),
			ResourceKind:      p["resource_kind"].(string),
			ResourceNamespace: "*_ALL_*",
			Actions:           expandStringList(p["actions"].([]interface{})),
		}
	}

	return result
}

func flattenPermissions(permissions []struct {
	ResourceGroup     string   `json:"resource-group"`
	ResourceKind      string   `json:"resource-kind"`
	ResourceNamespace string   `json:"resource-namespace"`
	Actions           []string `json:"actions"`
}) []interface{} {
	result := make([]interface{}, len(permissions))

	for i, perm := range permissions {
		result[i] = map[string]interface{}{
			"resource_group": perm.ResourceGroup,
			"resource_kind":  perm.ResourceKind,
			"actions":        perm.Actions,
		}
	}

	return result
}

func dscHash(v interface{}) int {
	m := v.(map[string]interface{})
	var buf bytes.Buffer
	if id, ok := m["id"].(string); ok {
		buf.WriteString(fmt.Sprintf("%s-", id))
	}
	if mac, ok := m["mac_address"].(string); ok {
		buf.WriteString(fmt.Sprintf("%s-", mac))
	}
	return schema.HashString(buf.String())
}

func pnicHash(v interface{}) int {
	m := v.(map[string]interface{})
	return schema.HashString(m["mac_address"].(string) + m["name"].(string))
}

func determineHostType(d *schema.ResourceData) string {
	hasDSCs := d.Get("dscs").(*schema.Set).Len() > 0
	hasPNICs := d.Get("pnic_info").(*schema.Set).Len() > 0

	if hasDSCs && hasPNICs {
		return "both"
	} else if hasDSCs {
		return "dss"
	} else if hasPNICs {
		return "pnic"
	}
	return ""
}

func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func validateAddressFamily(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	switch v {
	case "IPv4", "IPv6":
		// valid
	default:
		errs = append(errs, fmt.Errorf("%q must be either 'IPv4' or 'IPv6', got: %s", key, v))
	}
	return
}
