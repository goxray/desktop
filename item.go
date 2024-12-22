package main

import (
	"github.com/lilendian0x00/xray-knife/xray"
)

type Item struct {
	LabelVal     string            `json:"Label"`
	LinkVal      string            `json:"Link"`
	XRyConfigVal map[string]string `json:"xray_config"`

	active bool
}

func NewItem(label, link string, cfg xray.GeneralConfig) *Item {
	itm := &Item{
		LabelVal: label,
		LinkVal:  link,
	}
	itm.XRyConfigVal = itm.xrayBaseConfigToMap(cfg)

	return itm
}

func (c *Item) Active() bool {
	return c.active
}

func (c *Item) SetActive(active bool) {
	c.active = active
}

func (c *Item) Label() string {
	return c.LabelVal
}

func (c *Item) Link() string {
	return c.LinkVal
}

func (c *Item) XRayConfig() map[string]string {
	return c.XRyConfigVal
}

func (c *Item) SetXRayConfig(cfg xray.GeneralConfig) {
	c.XRyConfigVal = c.xrayBaseConfigToMap(cfg)
}

func (c *Item) xrayBaseConfigToMap(x xray.GeneralConfig) map[string]string {
	return map[string]string{
		"Protocol": x.Protocol, "Address": x.Address,
		"Security": x.Security, "Aid": x.Aid, "Host": x.Host,
		"ID": x.ID, "Network": x.Network, "Path": x.Path,
		"Port": x.Port, "Remark": x.Remark, "TLS": x.TLS,
		"SNI": x.SNI, "ALPN": x.ALPN, "TlsFingerprint": x.TlsFingerprint,
		"Authority": x.Authority, "ServiceName": x.ServiceName,
		"Mode": x.Mode, "Type": x.Type,
		"OrigLink": x.OrigLink,
	}
}
