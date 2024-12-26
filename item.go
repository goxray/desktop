package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/goxray/tun/pkg/client"
	vpn "github.com/goxray/tun/pkg/client"
	"github.com/lilendian0x00/xray-knife/xray"

	"github.com/goxray/ui/internal/netchart"
	"github.com/goxray/ui/window"
)

type Item struct {
	LabelVal     string            `json:"Label"`
	LinkVal      string            `json:"Link"`
	XRyConfigVal map[string]string `json:"xray_config"`

	active bool

	client   *client.Client
	recorder *netchart.Recorder
}

func NewItem(label, link string, cfg xray.GeneralConfig) *Item {
	itm := &Item{
		LabelVal: label,
		LinkVal:  link,
	}
	itm.XRyConfigVal = itm.xrayBaseConfigToMap(cfg)

	itm.Init()

	return itm
}

func (c *Item) Init() {
	cl, err := vpn.NewClientWithOpts(vpn.Config{
		Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	})
	if err != nil {
		panic(fmt.Errorf("create vpn client: %v", err))
	}
	c.client = cl

	c.recorder = netchart.NewRecorder(c.client)
	c.recorder.Start()
}

func (c *Item) Active() bool {
	return c.active
}

func (c *Item) SetActive(active bool) {
	c.active = active
}

func (c *Item) Connect() error {
	return c.client.Connect(c.Link())
}

func (c *Item) Disconnect() error {
	return c.client.Disconnect(context.Background())
}

func (c *Item) Recorder() window.NetworkRecorder {
	if c.client == nil {
		c.Init()
	}

	return c.recorder
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
