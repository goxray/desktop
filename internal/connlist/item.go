package connlist

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	vpn "github.com/goxray/tun/pkg/client"
	"github.com/lilendian0x00/xray-knife/v2/xray"

	"github.com/goxray/desktop/internal/netchart"
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
	if err := itm.init(); err != nil {
		return nil, err
	}
	itm.parent = parent

	return itm, nil
}

func (c *Item) init() error {
	proto, err := xray.ParseXrayConfig(c.Link())
	if err != nil {
		return fmt.Errorf("invalid xray link: %s", err)
	}

	c.xconfigMap, err = c.xrayBaseConfigToMap(proto)
	if err != nil {
		return fmt.Errorf("parse xray config to map: %s", err)
	}

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

func (c *Item) xrayBaseConfigToMap(proto xray.Protocol) (map[string]string, error) {
	x := proto.ConvertToGeneralConfig()
	xmap := map[string]string{
		"Protocol": x.Protocol, "Address": x.Address,
		"Security": x.Security, "Aid": x.Aid, "Host": x.Host,
		"ID": x.ID, "Network": x.Network, "Path": x.Path,
		"Port": x.Port, "Remark": x.Remark, "TLS": x.TLS,
		"SNI": x.SNI, "ALPN": x.ALPN, "TlsFingerprint": x.TlsFingerprint,
		"Authority": x.Authority, "ServiceName": x.ServiceName,
		"Mode": x.Mode, "Type": x.Type,
		"OrigLink": x.OrigLink,
	}

	// Marshalling will marshall the actual protocol, like proto.(*xray.Vmess)
	b, err := json.Marshal(proto)
	if err != nil {
		return nil, fmt.Errorf("marshal xray protocol: %w", err)
	}

	// Unmarshalling it will add protocol-specific values to the map.
	if err := json.Unmarshal(b, &xmap); err != nil {
		return nil, fmt.Errorf("unmarshal xray protocol: %w", err)
	}

	// Keys that duplicate base protocol values.
	removeDupKeys := []string{"add", "ps", "sni", "fp", "id", "OrigLink"}

	// Make all keys in map start with uppercase letter and remove duplicate keys.
	for k, _ := range xmap {
		if slices.Contains(removeDupKeys, k) {
			delete(xmap, k)
			continue
		}

		if strings.Title(k) != k {
			xmap[strings.Title(k)] = xmap[k]
			delete(xmap, k)
		}
	}

	return xmap, nil
}
