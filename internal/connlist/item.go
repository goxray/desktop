package connlist

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	vpn "github.com/goxray/tun/pkg/client"
	"github.com/lilendian0x00/xray-knife/v2/xray"

	"github.com/goxray/ui/internal/netchart"
)

type Client interface {
	Connect(string) error
	Disconnect(context.Context) error
	BytesRead() int
	BytesWritten() int
}

// Item is a combine that is passed (via interface segregation) throughout the system to apply
// centralized changes to connections with the smallest overhead as possible.
type Item struct {
	label      string
	link       string
	xconfigMap map[string]string
	active     bool

	parent   *Collection
	client   Client
	recorder NetworkRecorder
}

func newItem(label, link string, parent *Collection) (*Item, error) {
	itm := &Item{
		label: label,
		link:  link,
	}
	proto, err := xray.ParseXrayConfig(link)
	if err != nil {
		return nil, fmt.Errorf("invalid xray link: %s", err)
	}

	itm.xconfigMap = itm.xrayBaseConfigToMap(proto.ConvertToGeneralConfig())

	// TODO: too hardcoded, ok for now
	cl, err := vpn.NewClientWithOpts(vpn.Config{
		Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	})
	if err != nil {
		panic(fmt.Errorf("create vpn client: %v", err))
	}
	itm.client = cl

	itm.recorder = netchart.NewRecorder(itm.client)
	itm.recorder.Start()
	itm.parent = parent

	return itm, nil
}

func (c *Item) init() error {
	proto, err := xray.ParseXrayConfig(c.Link())
	if err != nil {
		return fmt.Errorf("invalid xray link: %s", err)
	}

	c.xconfigMap = c.xrayBaseConfigToMap(proto.ConvertToGeneralConfig())

	cl, err := vpn.NewClientWithOpts(vpn.Config{
		Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	})
	if err != nil {
		return fmt.Errorf("create vpn client: %v", err)
	}
	c.client = cl

	c.recorder = netchart.NewRecorder(c.client)
	c.recorder.Start()

	return nil
}

func (c *Item) Update(link, label string) error {
	c.link = link
	c.label = label
	if err := c.init(); err != nil {
		return err
	}
	c.parent.onChange()
	return nil
}

func (c *Item) Active() bool {
	return c.active
}

func (c *Item) SetActive(active bool) {
	c.active = active
	c.parent.onChange()
}

func (c *Item) Connect() error {
	return c.client.Connect(c.Link())
}

func (c *Item) Disconnect() error {
	return c.client.Disconnect(context.Background())
}

func (c *Item) Label() string {
	return c.label
}

func (c *Item) Link() string {
	return c.link
}

func (c *Item) XRayConfig() map[string]string {
	return c.xconfigMap
}

func (c *Item) SetXRayConfig(cfg xray.GeneralConfig) {
	c.xconfigMap = c.xrayBaseConfigToMap(cfg)
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
