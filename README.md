# psm-provider

## Build Process 
Currently the provider is distributed and built from Github. Build the code locally and then reference this provider. 

```
sudo apt update
sudo apt install golang git
git clone git clone https://github.com/farsonic/psm-provider.git
cd psm-provider
go mod tidy
make install 
```

The code will now be placed into ~/.terraform.d/plugins/local/provider/psm/X.X.X/linux_amd64 where is the version number of the build. This should be referenced as local/provider/psm. Once you install the provider it will be hosted locally with the current **Hostname = local** and the **Namespace = provider**. The Name of the provider is **PSM**.

Within your Terraform infrastructure file (ie main.tf) specify the provider with the following syntax. You can specify the version if wanted but it will always use the latest. If you do specify the version you will need to use the -upgrade switch to force an upgrade. 

```
git clone https://github.com/farsonic/psm-provider.git 
cd psm-provider
make
```

Once you install the provider it will be hosted locally with the current **Hostname = local** and the **Namespace = provider**. The Name of the provider is **PSM**.

Within your Terraform infrastructure file (ie main.tf) specify the provider with the following syntax. 

```
terraform { 
  required_providers {
   psm = { 
      source = "local/provider/psm"
  }
 }
}
```

You will need to configure the provider to communicate directly with the PSM server either with or without SSL certificate validation using the following definition. 

```
provider "psm" { 
  user = "admin"
  server = "https://PSM_SERVER"
  password = "PSM_PASSWORD"
  insecure = true
}
```

## Usage examples

### VRF Instance
VRF's provide isolation of routing tables as well as networks within the platform. The following definition will configure a new VRF instance capability within the DPU, however the underlying switch will also need to have the VRF configuration in place also. 

```
resource "psm_vrf" "customerABC" { 
  name = "CustomerABC"
  allow_session_reuse      = "enable"
  connection_tracking_mode = "enable"
}
```


### Network 
Within PSM a network definition defines the name of the network and the VLAN that will be redirected to a DPU. The following resource definition will create a network called "Database Network" which redirect VLAN 123 traffic to the DPU. The VLAN will need to be configured on at least one switch in the network to be successfully propagated to the DPU. Not you still need to conform to PSM naming guidelines for the name of the network. 

```
resource "psm_network" "network" {
  name     = "DatabaseNetwork"
  vlan_id  = 123
  connection_tracking_mode = "disable"
  allow_session_reuse      = "disable"
  service_bypass           = true
  virtual_router           = "CustomerABC"
  tenant   = "default" 
  vlan_id  = 123
}
```

### IP Collections
PSM allows the user to create groups of IP Addresses called IP Collections. These are then used within Security Policies (and elsewhere) to define the source and destination IP Addresses used for matches. Addresses must be a list of strings, commar seperated if there is more than one subnet. No mask on the address is also acceptable and will result in an implicit /32 host mask. 

```
resource "psm_ipcollection" "ipcollections" {
  name     = "DatabaseServers"
  addresses = ["10.10.10.0/24"]
}
```

An IP collection can also be nested. A nested IP collection is one that references another IP Collection. Up to four levels of nesting is supported with IP Collections.

```
resource "psm_ipcollection" "db01" {
  name     = "DatabaseServer01"
  addresses = ["10.10.10.100/32"]
}
resource "psm_ipcollection" "db02" {
  name     = "DatabaseServer02"
  addresses = ["10.10.10.101/32"]
}
resource "psm_ipcollection" "dbsrvs" {
  name     = "DatabaseServers"
  ip_collections = ["DatabaseServer01", "DatabaseServer02"]
  addresses = ["10.10.10.0/24"] 
}
```

### Security Policies 
Security policies are attached to either an individual network or to the VRF. If attached to a VRF then the policy is inherited by networks associated with that particular VRF. If the tenant and/or the policy_distribution_target is not defined these will default to the default VRF. The from IP_Collections will need to be defined prior to them being mapped within the rule. When pushing a security policy you can configure a definition wihtout any rules and then add rules. Rules will be applied in order, so order matters. Rules will need at least a pair of to/from_ip_addresses and/or to/from_ip_collections. 

```
resource "psm_rules" "ApplicationA_Stack" {
  policy_name                = "ApplicationStack"
  tenant                     = "default"
  policy_distribution_target = "default"
  rule {
      rule_name = "AllowSSHTraffic"
      description = "This rule allows SSH traffic from public IPs"
      from_ip_addresses   = ["10.9.0.0/24"]
      to_ip_addresses     = ["10.10.0.0/23", "10.99.1.1/32"]
      from_ip_collections = ["network2","network1"]
      to_ip_collections   = ["network4"]
      from_workloadgroups = ["WorkloadGroup1","WorkloadGroup2"]
      to_workloadgroups   = ["WorkloadGroup3","WorkloadGroup4"]
      apps = ["SSH"]
      action = "permit"
      disable             = false
      labels = {
        "Application" : "SSH"
      }
      apps = ["SSH"]
      action = "permit"
    }
}
```

Currently there is no ability to add individual protocol/port entries (watch this space) as well as ability to define custom application definitions. 


### Syslog Export

PSM can create Syslog Export Targets, so the DSM will be able to send telemetry data directly to up to four external receivers
Supported formats are: "syslog-rfc5424" and "syslog-bsd".
Filters can be applied to limit the applied action on sessions, supported filters: all, allow, deny


```
resource "psm_syslog_export_policy" "policy01" {
  name   = "policy01"
  format = "syslog-rfc5424" //syslog-rfc5424, syslog-bsd
  filter = ["deny"]         //all, allow, deny

  syslogconfig {
    facility         = "user"
    disable_batching = false //true, false
  }

  psm_target {
    enable = false
  }

  targets {
    destination = "1.1.1.1"
    transport   = "udp/514"
  }

  targets {
    destination = "2.2.2.2"
    transport   = "udp/514"
  }

  targets {
    destination = "3.3.3.3"
    transport   = "udp/514"
  }

  targets {
    destination = "4.4.4.4"
    transport   = "udp/514"
  }
}  
```

Binding the export policy is currently done manually via the DSS menu on PSM.

### Apps

Apps can be used inside Security Policies for easier handling. They can be nested up to five levels. 
They can be built simply based on tcp / udp ports:

```
resource "psm_app" "app" {      
  display_name = "example_app"      
  spec {      
    proto_ports {      
      protocol = "tcp"      
      ports    = "8080,9090"      
    }     
  }      
}    
```

Or more complex ones including other apps as nested objects:

```
resource "psm_app" "app" {
  display_name = "example_app"
  spec {
    proto_ports {
      protocol = "tcp"
      ports    = "7080,8080,8090-8092"
    }
    proto_ports {
      protocol = "udp"
      ports    = "5600-5800"
    }
    apps    = ["IMAP", "IMAPS"]
  }
}

```
ALGs (Application Layer Gateway) can be used to define Applications beside L3 / L4 based information

```
resource "psm_app" "dns_alg" {
  display_name = "dns_alg_test"
  spec {
    proto_ports {
      protocol = "udp"
      ports    = "53,5353"
    }
    alg {
      type = "dns"
      dns {
        drop_multi_question_packets    = true
        drop_long_label_packets        = true
        drop_large_domain_name_packets = false
        max_message_length             = 512
      }
    }
  }
}

resource "psm_app" "icmp_alg" {
  display_name = "icmp_alg_test"
  spec {
    alg {
      type = "icmp"
      icmp {
        type = "8"
        code = "0"
      }
    }
  }
}

resource "psm_app" "ftp_alg" {
  display_name = "ftp_alg_test"
  spec {
    proto_ports {
      protocol = "tcp"
      ports    = "21"
    }
    alg {
      type = "ftp"
      ftp {
        allow_mismatch_ip_address = true
      }
    }
  }
}

resource "psm_app" "sunrpc_alg" {
  display_name = "sunrpc_alg_test"
  spec {
    proto_ports {
      protocol = "tcp"
      ports    = "111"
    }
    timeout = "1h30m"
    alg {
      type = "sunrpc"
      sunrpc {
        program_id = "10024"
      }
    }
  }
}

resource "psm_app" "msrpc_alg" {
  display_name = "msrpc_alg_test"
  spec {
    proto_ports {
      protocol = "tcp"
      ports    = "135"
    }
    timeout = "1h30m"
    alg {
      type = "msrpc"
      msrpc {
        program_uuid = "a4f1db00-ca47-1067-b31f-00dd010662da"
      }
    }
  }
}

resource "psm_app" "tftp_alg" {
  display_name = "alg_tftp"
  spec {
    proto_ports {
      protocol = "udp"
      ports    = "69"
    }
    alg {
      type = "tftp"
    }
  }
}

resource "psm_app" "rtsp_alg" {
  display_name = "alg_rtsp"
  spec {
    proto_ports {
      protocol = "tcp"
      ports    = "554"
    }
    alg {
      type = "rtsp"
    }
  }
}
```

### Orchestrator / Hypervisor Integration

Building and working with Workload Groups in Security Policies relies on Workload Objects, which are basically VM Tags extracted from the vSphere environment.
A Read-Only user is needed for that.

To monitor only specific Data Centers inside the vSphere environment, namespaces can be configured optionally (default "all_namespaces")

Currently supported Hypervisors: vSphere 6.7, 7 and 8.


```
resource "psm_orchestrator" "vsphere01" {
  type      = "vcenter"
  name      = "vcenter01"
  uri       = "vcenter01.domain.name"
  auth_type = "username-password"
  username  = "psm@domain.name"
  password  = "Pensando0$!"

  namespaces {
    name = "dc01"
  }

  namespaces {
    name = "dc02"
  }
}
```


### IPFIX Export

PSM can create IPFIX Export Targets, so the DSM will be able to send telemetry data directly to external receivers

```
resource "psm_flow_export_policy" "ipfix" {
  name     = "IPFIXv2"
  interval = "10s"
  format   = "ipfix"

  target {
    destination = "1.1.1.1"
    transport   = "udp/9995"
  }
}
```

Binding the export policy is currently done manually via the DSS menu on PSM.

### Advanced usage 

Combine this all together and define your networks, subnets and firewall policies into a single definition within terraform. There is currently constraints around the order of execution, so ensure your networks and IP Collections are defined before you atempt to assign them to a security policy. 

```
locals {
  networks = {
    network1 = {
      name   = "network1"
      description = "My network"
      department = "Production"
      vlan   = 123
      subnet = "10.45.45.0/28"
      vrf = "default"
    },
    network2 = {
      name   = "network2"
      description = "My network"
      department = "Development"
      vlan   = 456
      subnet = "10.56.78.0/24"
      vrf = "default"
    }
    network4 = {
      name   = "network4"
      description = "My network"
      department = "Development"
      vlan   = 479
      subnet = "10.58.178.0/24"
      vrf = "default"
    }
    network7 = {
      name   = "network7"
      description = "My network"
      department = "Development"
      vlan   = 477
      subnet = "10.8.177.0/24"
      vrf = "CustomerABC"
    }
    HomeLAN = {
      name   = "HomeLAN"
      description = "Home Network"
      department = "Production"
      vlan   = 199
      subnet = "192.168.0.0/24"
      vrf = "default"
    }
    DEMO = {
      name   = "DEMONETWORK"
      description = "DEMO Network"
      department = "Development"
      vlan   = 1119
      subnet = "192.168.19.0/24"
      vrf = "default"
    }
    nosubnet = {
      name   = "nosubnet"
      description = "nosubnet Network"
      department = "Development"
      vlan   = 1120
      subnet = ""
      vrf = "default"
    }
  }
  firewall_rules = [
    {
      name = "webservertraffic"
      description = "All the webserver traffic"
      from_collection = ["network1", "network2"]
      to_collection   = ["network4", "network7"]
      apps = ["AH", "GRE"]
      action = "permit"
    },
    {
      name = "telnet"
      description = "All the Telnet traffic"
      from_collection = ["network1", "network2"]
      to_collection   = ["network4", "network7", "network1"]
      apps = ["TELNET"]
      action = "permit"
    },
    {
      name = "sshtraffic"
      description = "SSH traffic"
      from_collection = ["network1", "network2"]
      to_collection   = ["network4", "network7"]
      apps = ["SSH"]
      action = "permit"
    },
    {
      name = "bgptraffic"
      description = "All the BGP traffic"
      from_collection = ["network1", "network2"]
      to_collection   = ["network4", "network7"]
      apps = ["BGP"]
      action = "permit"
    },
    {
      name = "webtraffic"
      description = "All the WEB traffic"
      from_collection = ["network1", "network2","DEMONETWORK"]
      to_collection   = ["network4", "network7"]
      apps = ["HTTPS"]
      action = "permit"
    },
    {
      name = "DefaultPolicy"
      description = "Default "
      from_addresses = ["any"]
      to_addresses   = ["any"]
      apps = ["ALL_TCP","ALL_UDP"]
      action = "deny"
    }
  ]
 filtered_networks = {for key, value in local.networks : key => value if value.subnet != ""}
}

resource "psm_vrf" "vrfs" {
  for_each = local.networks
  name     = each.value.vrf
}

resource "psm_network" "network" {
  for_each = local.networks
  name     = each.value.name
  tenant   = each.value.vrf
  vlan_id  = each.value.vlan
  depends_on = [psm_vrf.vrfs]
}

resource "psm_ipcollection" "ipcollections" {
  for_each = local.filtered_networks
  name     = each.value.name
  addresses = [each.value.subnet]
  depends_on = [psm_network.network]
}

resource "psm_rules" "default_vrf_policy" {
  policy_name                = "test"
  tenant                     = "default"
  policy_distribution_target = "default"

  dynamic "rule" {
    for_each = local.firewall_rules
    content {
      rule_name           = rule.value.name
      description         = rule.value.description
      from_ip_collections = lookup(rule.value, "from_collection", [])
      to_ip_collections   = lookup(rule.value, "to_collection", [])
      from_ip_addresses   = lookup(rule.value, "from_addresses", [])
      to_ip_addresses     = lookup(rule.value, "to_addresses", [])
      apps                = rule.value.apps
      action              = rule.value.action
    }
  }
  depends_on = [psm_ipcollection.ipcollections]
}
```
