package goaviatrix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const Phase1AuthDefault = "SHA-256"
const Phase1DhGroupDefault = "14"
const Phase1EncryptionDefault = "AES-256-CBC"
const Phase2AuthDefault = "HMAC-SHA-256"
const Phase2DhGroupDefault = "14"
const Phase2EncryptionDefault = "AES-256-CBC"
const SslServerPoolDefault = "192.168.44.0/24"

// Site2Cloud simple struct to hold site2cloud details
type Site2Cloud struct {
	Action                        string   `form:"action,omitempty"`
	CID                           string   `form:"CID,omitempty"`
	VpcID                         string   `form:"vpc_id,omitempty" json:"vpc_id,omitempty"`
	TunnelName                    string   `form:"connection_name" json:"name,omitempty"`
	RemoteGwType                  string   `form:"remote_gateway_type,omitempty"`
	ConnType                      string   `form:"connection_type,omitempty" json:"type,omitempty"`
	TunnelType                    string   `form:"tunnel_type,omitempty" json:"tunnel_type,omitempty"`
	GwName                        string   `form:"primary_cloud_gateway_name,omitempty" json:"gw_name,omitempty"`
	BackupGwName                  string   `form:"backup_gateway_name,omitempty"`
	RemoteGwIP                    string   `form:"remote_gateway_ip,omitempty" json:"peer_ip,omitempty"`
	RemoteGwIP2                   string   `form:"backup_remote_gateway_ip,omitempty"`
	PreSharedKey                  string   `form:"pre_shared_key,omitempty"`
	BackupPreSharedKey            string   `form:"backup_pre_shared_key,omitempty"`
	RemoteSubnet                  string   `form:"remote_subnet_cidr,omitempty" json:"remote_cidr,omitempty"`
	LocalSubnet                   string   `form:"local_subnet_cidr,omitempty" json:"local_cidr,omitempty"`
	HAEnabled                     string   `form:"ha_enabled,omitempty" json:"ha_status,omitempty"`
	PeerType                      string   `form:"peer_type,omitempty"`
	SslServerPool                 string   `form:"ssl_server_pool,omitempty"`
	NetworkType                   string   `form:"network_type,omitempty"`
	CloudSubnetCidr               string   `form:"cloud_subnet_cidr,omitempty"`
	RemoteCidr                    string   `form:"remote_cidr,omitempty"`
	RemoteSubnetVirtual           string   `form:"virtual_remote_subnet_cidr,omitempty" json:"virtual_remote_subnet_cidr,omitempty"`
	LocalSubnetVirtual            string   `form:"virtual_local_subnet_cidr,omitempty" json:"virtual_local_subnet_cidr,omitempty"`
	Phase1Auth                    string   `form:"phase1_auth,omitempty"`
	Phase1DhGroups                string   `form:"phase1_dh_group,omitempty"`
	Phase1Encryption              string   `form:"phase1_encryption,omitempty"`
	Phase2Auth                    string   `form:"phase2_auth,omitempty"`
	Phase2DhGroups                string   `form:"phase2_dh_group,omitempty"`
	Phase2Encryption              string   `form:"phase2_encryption,omitempty"`
	EnableIKEv2                   string   `form:"enable_ikev2,omitempty"`
	PrivateRouteEncryption        string   `form:"private_route_encryption,omitempty"`
	RemoteGwLatitude              float64  `form:"remote_gateway_latitude,omitempty"`
	RemoteGwLongitude             float64  `form:"remote_gateway_longitude,omitempty"`
	BackupRemoteGwLatitude        float64  `form:"backup_remote_gateway_latitude,omitempty"`
	BackupRemoteGwLongitude       float64  `form:"backup_remote_gateway_longitude,omitempty"`
	RouteTableList                []string `form:"route_table_list,omitempty"`
	CustomAlgorithms              bool
	DeadPeerDetection             bool
	EnableActiveActive            bool
	ForwardToTransit              bool
	EventTriggeredHA              bool
	CustomMap                     bool   `form:"custom_map,omitempty"`
	RemoteSourceRealCIDRs         string `form:"remote_src_real_cidrs,omitempty"`
	RemoteSourceVirtualCIDRs      string `form:"remote_src_virt_cidrs,omitempty"`
	RemoteDestinationRealCIDRs    string `form:"remote_dst_real_cidrs,omitempty"`
	RemoteDestinationVirtualCIDRs string `form:"remote_dst_virt_cidrs,omitempty"`
	LocalSourceRealCIDRs          string `form:"local_src_real_cidrs,omitempty"`
	LocalSourceVirtualCIDRs       string `form:"local_src_virt_cidrs,omitempty"`
	LocalDestinationRealCIDRs     string `form:"local_dst_real_cidrs,omitempty"`
	LocalDestinationVirtualCIDRs  string `form:"local_dst_virt_cidrs,omitempty"`
	LocalTunnelIp                 string `form:"local_tunnel_ip,omitempty"`
	RemoteTunnelIp                string `form:"remote_tunnel_ip,omitempty"`
	BackupLocalTunnelIp           string `form:"backup_local_tunnel_ip,omitempty"`
	BackupRemoteTunnelIp          string `form:"backup_remote_tunnel_ip,omitempty"`
	EnableSingleIpHA              bool
	Phase1RemoteIdentifier        string
}

type EditSite2Cloud struct {
	Action                        string `form:"action,omitempty"`
	CID                           string `form:"CID,omitempty"`
	VpcID                         string `form:"vpc_id,omitempty"`
	ConnName                      string `form:"conn_name"`
	GwName                        string `form:"primary_cloud_gateway_name,omitempty"`
	NetworkType                   string `form:"network_type,omitempty"`
	CloudSubnetCidr               string `form:"cloud_subnet_cidr,omitempty"`
	CloudSubnetVirtual            string `form:"cloud_subnet_virtual,omitempty"`
	RemoteSourceRealCIDRs         string `form:"remote_src_real_cidrs,omitempty"`
	RemoteSourceVirtualCIDRs      string `form:"remote_src_virt_cidrs,omitempty"`
	RemoteDestinationRealCIDRs    string `form:"remote_dst_real_cidrs,omitempty"`
	RemoteDestinationVirtualCIDRs string `form:"remote_dst_virt_cidrs,omitempty"`
	LocalSourceRealCIDRs          string `form:"local_src_real_cidrs,omitempty"`
	LocalSourceVirtualCIDRs       string `form:"local_src_virt_cidrs,omitempty"`
	LocalDestinationRealCIDRs     string `form:"local_dst_real_cidrs,omitempty"`
	LocalDestinationVirtualCIDRs  string `form:"local_dst_virt_cidrs,omitempty"`
	Phase1RemoteIdentifier        string `form:"phase1_remote_identifier,omitempty"`
}

type Site2CloudResp struct {
	Return  bool               `json:"return"`
	Results Site2CloudConnList `json:"results"`
	Reason  string             `json:"reason"`
}

type Site2CloudConnList struct {
	Connections []Site2Cloud `json:"connections"`
}

type EditSite2CloudConnDetail struct {
	VpcID                         []string      `json:"vpc_id,omitempty"`
	TunnelName                    []string      `json:"name,omitempty"`
	ConnType                      string        `json:"type,omitempty"`
	TunnelType                    string        `json:"tunnel_type,omitempty"`
	GwName                        string        `json:"gw_name,omitempty"`
	Tunnels                       []TunnelInfo  `json:"tunnels,omitempty"`
	RemoteSubnet                  string        `json:"real_remote_cidr,omitempty"`
	LocalSubnet                   string        `json:"real_local_cidr,omitempty"`
	RemoteCidr                    string        `json:"remote_cidr,omitempty"`
	LocalCidr                     string        `json:"local_cidr,omitempty"`
	HAEnabled                     string        `json:"ha_status,omitempty"`
	PeerType                      string        `json:"peer_type,omitempty"`
	RemoteSubnetVirtual           string        `json:"virt_remote_cidr,omitempty"`
	LocalSubnetVirtual            string        `json:"virt_local_cidr,omitempty"`
	Algorithm                     AlgorithmInfo `json:"algorithm,omitempty"`
	RouteTableList                []string      `json:"rtbls,omitempty"`
	SslServerPool                 []string      `json:"ssl_server_pool,omitempty"`
	DeadPeerDetectionConfig       string        `json:"dpd_config,omitempty"`
	EnableActiveActive            string        `json:"active_active_ha,omitempty"`
	EnableIKEv2                   string        `json:"ike_ver,omitempty"`
	BgpLocalASN                   string        `json:"bgp_local_asn_number,omitempty"`
	BgpLocalIP                    string        `json:"bgp_local_ip,omitempty"`
	BgpBackupLocalIP              string        `json:"bgp_backup_local_ip,omitempty"`
	BgpRemoteASN                  string        `json:"bgp_remote_asn_number,omitempty"`
	BgpRemoteIP                   string        `json:"bgp_remote_ip,omitempty"`
	BgpBackupRemoteIP             string        `json:"bgp_backup_remote_ip,omitempty"`
	EnableGlobalAccelerator       bool          `json:"globalaccel,omitempty"`
	AzureAccountName              string        `json:"arm_account_name,omitempty"`
	AzureResourceGroup            string        `json:"arm_resource_group,omitempty"`
	AzureVhubName                 string        `json:"arm_vhub_name,omitempty"`
	ForwardToTransit              string        `json:"forward_to_transit"`
	RemoteSourceRealCIDRs         string        `json:"remote_src_real_cidrs"`
	RemoteSourceVirtualCIDRs      string        `json:"remote_src_virt_cidrs"`
	RemoteDestinationRealCIDRs    string        `json:"remote_dst_real_cidrs"`
	RemoteDestinationVirtualCIDRs string        `json:"remote_dst_virt_cidrs"`
	LocalSourceRealCIDRs          string        `json:"local_src_real_cidrs"`
	LocalSourceVirtualCIDRs       string        `json:"local_src_virt_cidrs"`
	LocalDestinationRealCIDRs     string        `json:"local_dst_real_cidrs"`
	LocalDestinationVirtualCIDRs  string        `json:"local_dst_virt_cidrs"`
	ManualBGPCidrs                []string      `json:"conn_bgp_manual_advertise_cidrs"`
	EventTriggeredHA              string        `json:"event_triggered_ha"`
	EnableSingleIpHA              string        `json:"single_ip_ha,omitempty"`
	Phase1RemoteIdentifier        string        `json:"phase1_remote_id"`
}

type Site2CloudConnDetailResp struct {
	Return  bool                     `json:"return"`
	Results Site2CloudConnDetailList `json:"results"`
	Reason  string                   `json:"reason"`
}

type Site2CloudConnDetailList struct {
	Connections EditSite2CloudConnDetail `json:"connections"`
}

type TunnelInfo struct {
	Status         string `json:"status"`
	IPAddr         string `json:"ip_addr"`
	Name           string `json:"name"`
	PeerIP         string `json:"peer_ip"`
	GwName         string `json:"gw_name"`
	TunnelStatus   string `json:"tunnel_status"`
	TunnelProtocol string `json:"tunnel_protocol"`
}

type AlgorithmInfo struct {
	Phase1Auth      []string `json:"ph1_auth,omitempty"`
	Phase1DhGroups  []string `json:"ph1_dh,omitempty"`
	Phase1Encrption []string `json:"ph1_encr,omitempty"`
	Phase2Auth      []string `json:"ph2_auth,omitempty"`
	Phase2DhGroups  []string `json:"ph2_dh,omitempty"`
	Phase2Encrption []string `json:"ph2_encr,omitempty"`
}

func (c *Client) CreateSite2Cloud(site2cloud *Site2Cloud) error {
	// TODO: use PostAPI - long form
	Url, err := url.Parse(c.baseURL)
	if err != nil {
		return errors.New("url Parsing failed for add_site2cloud " + err.Error())
	}
	addSite2cloud := url.Values{}
	addSite2cloud.Add("CID", c.CID)
	addSite2cloud.Add("action", "add_site2cloud")
	addSite2cloud.Add("vpc_id", site2cloud.VpcID)
	addSite2cloud.Add("connection_name", site2cloud.TunnelName)
	addSite2cloud.Add("connection_type", site2cloud.ConnType)
	addSite2cloud.Add("remote_gateway_type", site2cloud.RemoteGwType)
	addSite2cloud.Add("tunnel_type", site2cloud.TunnelType)

	addSite2cloud.Add("ha_enabled", site2cloud.HAEnabled)
	addSite2cloud.Add("backup_gateway_name", site2cloud.BackupGwName)
	addSite2cloud.Add("backup_remote_gateway_ip", site2cloud.RemoteGwIP2)

	addSite2cloud.Add("phase1_auth", site2cloud.Phase1Auth)
	addSite2cloud.Add("phase1_dh_group", site2cloud.Phase1DhGroups)
	addSite2cloud.Add("phase1_encryption", site2cloud.Phase1Encryption)
	addSite2cloud.Add("phase2_auth", site2cloud.Phase2Auth)
	addSite2cloud.Add("phase2_dh_group", site2cloud.Phase2DhGroups)
	addSite2cloud.Add("phase2_encryption", site2cloud.Phase2Encryption)

	if site2cloud.TunnelType == "tcp" {
		addSite2cloud.Add("ssl_server_pool", site2cloud.SslServerPool)
	}

	if site2cloud.EnableIKEv2 == "true" {
		addSite2cloud.Add("enable_ikev2", "true")
	}

	if site2cloud.PrivateRouteEncryption == "true" {
		addSite2cloud.Add("private_route_encryption", site2cloud.PrivateRouteEncryption)
		if len(site2cloud.RouteTableList) != 0 {
			for i := range site2cloud.RouteTableList {
				addSite2cloud.Add("route_table_list["+strconv.Itoa(i)+"]", site2cloud.RouteTableList[i])
			}
		}
		latitude := fmt.Sprintf("%f", site2cloud.RemoteGwLatitude)
		longitude := fmt.Sprintf("%f", site2cloud.RemoteGwLongitude)
		addSite2cloud.Add("remote_gateway_latitude", latitude)
		addSite2cloud.Add("remote_gateway_longitude", longitude)
		if site2cloud.HAEnabled == "yes" {
			backupLatitude := fmt.Sprintf("%f", site2cloud.BackupRemoteGwLatitude)
			backupLongitude := fmt.Sprintf("%f", site2cloud.BackupRemoteGwLongitude)
			addSite2cloud.Add("remote_gateway_latitude", backupLatitude)
			addSite2cloud.Add("remote_gateway_longitude", backupLongitude)
		}
	}

	addSite2cloud.Add("primary_cloud_gateway_name", site2cloud.GwName)
	addSite2cloud.Add("remote_gateway_ip", site2cloud.RemoteGwIP)
	addSite2cloud.Add("remote_subnet_cidr", site2cloud.RemoteSubnet)
	addSite2cloud.Add("local_subnet_cidr", site2cloud.LocalSubnet)
	addSite2cloud.Add("virtual_remote_subnet_cidr", site2cloud.RemoteSubnetVirtual)
	addSite2cloud.Add("virtual_local_subnet_cidr", site2cloud.LocalSubnetVirtual)

	addSite2cloud.Add("pre_shared_key", site2cloud.PreSharedKey)
	addSite2cloud.Add("backup_pre_shared_key", site2cloud.BackupPreSharedKey)

	if site2cloud.CustomMap {
		addSite2cloud.Add("custom_map", strconv.FormatBool(site2cloud.CustomMap))
		addSite2cloud.Add("remote_src_real_cidrs", site2cloud.RemoteSourceRealCIDRs)
		addSite2cloud.Add("remote_src_virt_cidrs", site2cloud.RemoteSourceVirtualCIDRs)
		addSite2cloud.Add("remote_dst_real_cidrs", site2cloud.RemoteDestinationRealCIDRs)
		addSite2cloud.Add("remote_dst_virt_cidrs", site2cloud.RemoteDestinationVirtualCIDRs)
		addSite2cloud.Add("local_src_real_cidrs", site2cloud.LocalSourceRealCIDRs)
		addSite2cloud.Add("local_src_virt_cidrs", site2cloud.LocalSourceVirtualCIDRs)
		addSite2cloud.Add("local_dst_real_cidrs", site2cloud.LocalDestinationRealCIDRs)
		addSite2cloud.Add("local_dst_virt_cidrs", site2cloud.LocalDestinationVirtualCIDRs)
	}

	addSite2cloud.Add("local_tunnel_ip", site2cloud.LocalTunnelIp)
	addSite2cloud.Add("remote_tunnel_ip", site2cloud.RemoteTunnelIp)
	addSite2cloud.Add("backup_local_tunnel_ip", site2cloud.BackupLocalTunnelIp)
	addSite2cloud.Add("backup_remote_tunnel_ip", site2cloud.BackupRemoteTunnelIp)
	if site2cloud.EnableSingleIpHA {
		addSite2cloud.Add("enable_single_ip_ha", "true")
	}

	Url.RawQuery = addSite2cloud.Encode()
	resp, err := c.Get(Url.String(), nil)
	if err != nil {
		return errors.New("HTTP Get add_site2cloud failed: " + err.Error())
	}
	var data APIResp
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	bodyIoCopy := strings.NewReader(bodyString)
	if err = json.NewDecoder(bodyIoCopy).Decode(&data); err != nil {
		return errors.New("Json Decode add_site2cloud failed: " + err.Error() + "\n Body: " + bodyString)
	}
	if !data.Return {
		return errors.New("Rest API add_site2cloud Get failed: " + data.Reason)
	}
	return nil
}

func (c *Client) GetSite2Cloud(site2cloud *Site2Cloud) (*Site2Cloud, error) {
	form := map[string]string{
		"CID":             c.CID,
		"action":          "list_site2cloud_conn",
		"connection_name": site2cloud.TunnelName,
	}

	var data Site2CloudResp

	err := c.GetAPI(&data, form["action"], form, BasicCheck)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(data.Results.Connections); i++ {
		conn := data.Results.Connections[i]
		if site2cloud.VpcID == conn.VpcID && site2cloud.TunnelName == conn.TunnelName {
			return &conn, nil
		}
	}
	return nil, ErrNotFound
}

func (c *Client) GetSite2CloudConnDetail(site2cloud *Site2Cloud) (*Site2Cloud, error) {
	form := map[string]string{
		"CID":       c.CID,
		"action":    "get_site2cloud_conn_detail",
		"conn_name": site2cloud.TunnelName,
		"vpc_id":    site2cloud.VpcID,
	}
	check := func(action, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "does not exist") {
				return ErrNotFound
			}
			return fmt.Errorf("rest API %s %s failed: %s", action, method, reason)
		}
		return nil
	}
	var data Site2CloudConnDetailResp
	err := c.GetAPI(&data, form["action"], form, check)
	if err != nil {
		return nil, err
	}

	s2cConnDetail := data.Results.Connections
	if len(s2cConnDetail.TunnelName) != 0 {
		site2cloud.GwName = s2cConnDetail.GwName
		site2cloud.ConnType = s2cConnDetail.ConnType
		if s2cConnDetail.TunnelType == "policy" || s2cConnDetail.TunnelType == "Policy" || s2cConnDetail.TunnelType == "Site2Cloud_Policy" {
			site2cloud.TunnelType = "policy"
		} else if s2cConnDetail.TunnelType == "route" || s2cConnDetail.TunnelType == "Route" || s2cConnDetail.TunnelType == "Site2Cloud_Routed" {
			site2cloud.TunnelType = "route"
		}
		site2cloud.RemoteGwType = s2cConnDetail.PeerType
		if site2cloud.ConnType == "mapped" {
			site2cloud.RemoteSubnet = s2cConnDetail.RemoteSubnet
			site2cloud.LocalSubnet = s2cConnDetail.LocalSubnet
			site2cloud.RemoteSubnetVirtual = s2cConnDetail.RemoteSubnetVirtual
			site2cloud.LocalSubnetVirtual = s2cConnDetail.LocalSubnetVirtual
		} else {
			site2cloud.RemoteSubnet = s2cConnDetail.RemoteCidr
			site2cloud.LocalSubnet = s2cConnDetail.LocalCidr
		}
		site2cloud.HAEnabled = s2cConnDetail.HAEnabled
		for i := range s2cConnDetail.Tunnels {
			if s2cConnDetail.Tunnels[i].GwName == site2cloud.GwName {
				site2cloud.RemoteGwIP = s2cConnDetail.Tunnels[i].PeerIP
			} else {
				site2cloud.BackupGwName = s2cConnDetail.Tunnels[i].GwName
				site2cloud.RemoteGwIP2 = s2cConnDetail.Tunnels[i].PeerIP
			}
		}
		if s2cConnDetail.Algorithm.Phase1Auth[0] == Phase1AuthDefault &&
			s2cConnDetail.Algorithm.Phase2Auth[0] == Phase2AuthDefault &&
			s2cConnDetail.Algorithm.Phase1DhGroups[0] == Phase1DhGroupDefault &&
			s2cConnDetail.Algorithm.Phase2DhGroups[0] == Phase2DhGroupDefault &&
			s2cConnDetail.Algorithm.Phase1Encrption[0] == Phase1EncryptionDefault &&
			s2cConnDetail.Algorithm.Phase2Encrption[0] == Phase2EncryptionDefault {
			site2cloud.CustomAlgorithms = false
			site2cloud.Phase1Auth = ""
			site2cloud.Phase2Auth = ""
			site2cloud.Phase1DhGroups = ""
			site2cloud.Phase2DhGroups = ""
			site2cloud.Phase1Encryption = ""
			site2cloud.Phase2Encryption = ""
		} else {
			site2cloud.CustomAlgorithms = true
			site2cloud.Phase1Auth = s2cConnDetail.Algorithm.Phase1Auth[0]
			site2cloud.Phase2Auth = s2cConnDetail.Algorithm.Phase2Auth[0]
			site2cloud.Phase1DhGroups = s2cConnDetail.Algorithm.Phase1DhGroups[0]
			site2cloud.Phase2DhGroups = s2cConnDetail.Algorithm.Phase2DhGroups[0]
			site2cloud.Phase1Encryption = s2cConnDetail.Algorithm.Phase1Encrption[0]
			site2cloud.Phase2Encryption = s2cConnDetail.Algorithm.Phase2Encrption[0]
		}
		if len(s2cConnDetail.RouteTableList) > 0 {
			site2cloud.RouteTableList = s2cConnDetail.RouteTableList
			site2cloud.PrivateRouteEncryption = "true"
		} else {
			site2cloud.PrivateRouteEncryption = "false"
		}
		if s2cConnDetail.SslServerPool[0] != "192.168.44.0/24" {
			site2cloud.SslServerPool = s2cConnDetail.SslServerPool[0]
		}
		if s2cConnDetail.DeadPeerDetectionConfig == "enable" {
			site2cloud.DeadPeerDetection = true
		} else if s2cConnDetail.DeadPeerDetectionConfig == "disable" {
			site2cloud.DeadPeerDetection = false
		}
		if s2cConnDetail.EnableActiveActive == "enable" || s2cConnDetail.EnableActiveActive == "Enable" {
			site2cloud.EnableActiveActive = true
		} else {
			site2cloud.EnableActiveActive = false
		}
		if s2cConnDetail.ForwardToTransit == "enable" {
			site2cloud.ForwardToTransit = true
		} else {
			site2cloud.ForwardToTransit = false
		}
		if s2cConnDetail.EnableIKEv2 == "2" {
			site2cloud.EnableIKEv2 = "true"
		}
		site2cloud.EventTriggeredHA = s2cConnDetail.EventTriggeredHA == "enabled"
		site2cloud.RemoteSourceRealCIDRs = s2cConnDetail.RemoteSourceRealCIDRs
		site2cloud.RemoteSourceVirtualCIDRs = s2cConnDetail.RemoteSourceVirtualCIDRs
		site2cloud.RemoteDestinationRealCIDRs = s2cConnDetail.RemoteDestinationRealCIDRs
		site2cloud.RemoteDestinationVirtualCIDRs = s2cConnDetail.RemoteDestinationVirtualCIDRs
		site2cloud.LocalSourceRealCIDRs = s2cConnDetail.LocalSourceRealCIDRs
		site2cloud.LocalSourceVirtualCIDRs = s2cConnDetail.LocalSourceVirtualCIDRs
		site2cloud.LocalDestinationRealCIDRs = s2cConnDetail.LocalDestinationRealCIDRs
		site2cloud.LocalDestinationVirtualCIDRs = s2cConnDetail.LocalDestinationVirtualCIDRs
		site2cloud.LocalTunnelIp = s2cConnDetail.BgpLocalIP
		site2cloud.RemoteTunnelIp = s2cConnDetail.BgpRemoteIP
		if site2cloud.HAEnabled == "enabled" {
			site2cloud.BackupLocalTunnelIp = s2cConnDetail.BgpBackupLocalIP
			site2cloud.BackupRemoteTunnelIp = s2cConnDetail.BgpBackupRemoteIP
		}
		site2cloud.EnableSingleIpHA = s2cConnDetail.EnableSingleIpHA == "enabled"
		site2cloud.Phase1RemoteIdentifier = s2cConnDetail.Phase1RemoteIdentifier
		return site2cloud, nil
	}

	return nil, ErrNotFound
}

func (c *Client) UpdateSite2Cloud(site2cloud *EditSite2Cloud) error {
	site2cloud.CID = c.CID
	site2cloud.Action = "edit_site2cloud_conn"

	return c.PostAPI(site2cloud.Action, site2cloud, BasicCheck)
}

func (c *Client) DeleteSite2Cloud(site2cloud *Site2Cloud) error {
	site2cloud.CID = c.CID
	site2cloud.Action = "delete_site2cloud_connection"

	return c.PostAPI(site2cloud.Action, site2cloud, BasicCheck)
}

func (c *Client) EnableDeadPeerDetection(site2cloud *Site2Cloud) error {
	form := map[string]string{
		"CID":             c.CID,
		"action":          "enable_dpd_config",
		"vpc_id":          site2cloud.VpcID,
		"connection_name": site2cloud.TunnelName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) DisableDeadPeerDetection(site2cloud *Site2Cloud) error {
	form := map[string]string{
		"CID":             c.CID,
		"action":          "disable_dpd_config",
		"vpc_id":          site2cloud.VpcID,
		"connection_name": site2cloud.TunnelName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) EnableSite2cloudActiveActive(site2cloud *Site2Cloud) error {
	form := map[string]string{
		"CID":             c.CID,
		"action":          "enable_site2cloud_active_active_ha",
		"vpc_id":          site2cloud.VpcID,
		"connection_name": site2cloud.TunnelName,
	}

	checkFunc := func(act, method, reason string, ret bool) error {
		if !ret {
			if strings.Contains(reason, "already enabled") {
				return nil
			}
			return fmt.Errorf("rest API %s %s failed: %s", act, method, reason)
		}
		return nil
	}

	return c.PostAPI(form["action"], form, checkFunc)
}

func (c *Client) DisableSite2cloudActiveActive(site2cloud *Site2Cloud) error {
	form := map[string]string{
		"CID":             c.CID,
		"action":          "disable_site2cloud_active_active_ha",
		"vpc_id":          site2cloud.VpcID,
		"connection_name": site2cloud.TunnelName,
	}

	return c.PostAPI(form["action"], form, BasicCheck)
}

func (c *Client) EnableSpokeMappedSite2CloudForwarding(site2cloud *Site2Cloud) error {
	data := map[string]string{
		"CID":             c.CID,
		"action":          "enable_spoke_mapped_site2cloud_forwarding",
		"vpc_id":          site2cloud.VpcID,
		"connection_name": site2cloud.TunnelName,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) DisableSpokeMappedSite2CloudForwarding(site2cloud *Site2Cloud) error {
	data := map[string]string{
		"CID":             c.CID,
		"action":          "disable_spoke_mapped_site2cloud_forwarding",
		"vpc_id":          site2cloud.VpcID,
		"connection_name": site2cloud.TunnelName,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) EnableSite2CloudEventTriggeredHA(vpcID, connectionName string) error {
	data := map[string]string{
		"CID":             c.CID,
		"action":          "enable_site2cloud_event_triggered_ha",
		"vpc_id":          vpcID,
		"connection_name": connectionName,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func (c *Client) DisableSite2CloudEventTriggeredHA(vpcID, connectionName string) error {
	data := map[string]string{
		"CID":             c.CID,
		"action":          "disable_site2cloud_event_triggered_ha",
		"vpc_id":          vpcID,
		"connection_name": connectionName,
	}
	return c.PostAPI(data["action"], data, BasicCheck)
}

func S2CPh1RemoteIdDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.HasChange("ha_enabled") || d.HasChange("enable_single_ip_ha") {
		return false
	}

	ip := d.Get("remote_gateway_ip").(string)
	haip := d.Get("backup_remote_gateway_ip").(string)
	o, n := d.GetChange("phase1_remote_identifier")
	haEnabled := d.Get("ha_enabled").(bool)
	singleIpHA := d.Get("enable_single_ip_ha").(bool)

	ph1RemoteIdListOld := ExpandStringList(o.([]interface{}))
	ph1RemoteIdListNew := ExpandStringList(n.([]interface{}))

	if len(ph1RemoteIdListOld) != 0 && len(ph1RemoteIdListNew) != 0 {
		if haEnabled && !singleIpHA {
			if len(ph1RemoteIdListNew) != 2 || len(ph1RemoteIdListOld) != 2 {
				return false
			}
			return ph1RemoteIdListOld[0] == ip && ph1RemoteIdListNew[0] == ip &&
				strings.TrimSpace(ph1RemoteIdListOld[1]) == haip && strings.TrimSpace(ph1RemoteIdListNew[1]) == haip
		} else if !haEnabled || singleIpHA {
			if len(ph1RemoteIdListNew) != 1 {
				return false
			}
			return ph1RemoteIdListOld[0] == ip && ph1RemoteIdListNew[0] == ip
		}
	}

	if !haEnabled || singleIpHA {
		if len(ph1RemoteIdListOld) == 1 && ph1RemoteIdListOld[0] == ip && len(ph1RemoteIdListNew) == 0 {
			return true
		}
	}

	if haEnabled && !singleIpHA {
		if len(ph1RemoteIdListOld) == 2 && ph1RemoteIdListOld[0] == ip && strings.TrimSpace(ph1RemoteIdListOld[1]) == haip && len(ph1RemoteIdListNew) == 0 {
			return true
		}
	}

	return false
}
