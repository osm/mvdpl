package mvdparser

import (
	"strings"

	"github.com/osm/quake/common/args"
	"github.com/osm/quake/common/ascii"
	"github.com/osm/quake/common/context"
	"github.com/osm/quake/common/infostring"
	"github.com/osm/quake/demo/mvd"
	"github.com/osm/quake/packet/command/stufftext"
	"github.com/osm/quake/packet/command/updateping"
	"github.com/osm/quake/packet/command/updatepl"
	"github.com/osm/quake/packet/command/updateuserinfo"
	"github.com/osm/quake/packet/svc"
)

type Parser struct {
	isStarted bool
	players   map[byte]*Player
	elapsed   float64

	Events []Event
}

type Event interface {
	Name() string
	Timestamp() float64
	Value() int16
	Suffix() string
}

type PacketLoss struct {
	name      string
	timestamp float64
	lossage   int16
}

func (pl *PacketLoss) Name() string       { return pl.name }
func (pl *PacketLoss) Timestamp() float64 { return pl.timestamp }
func (pl *PacketLoss) Value() int16       { return pl.lossage }
func (pl *PacketLoss) Suffix() string     { return "% pl" }

type Ping struct {
	name      string
	timestamp float64
	ping      int16
}

func (p *Ping) Name() string       { return p.name }
func (p *Ping) Timestamp() float64 { return p.timestamp }
func (p *Ping) Value() int16       { return p.ping }
func (p *Ping) Suffix() string     { return " ms" }

type Player struct {
	Team string
	Name string
}

func New() *Parser {
	return &Parser{
		players: make(map[byte]*Player),
	}
}

func (p *Parser) Parse(data []byte) ([]Event, error) {
	demo, err := mvd.Parse(context.New(), data)
	if err != nil {
		return nil, err
	}

	for _, d := range demo.Data {
		if p.isStarted {
			p.elapsed += float64(d.Timestamp) * 0.001
		}

		if d.Read == nil {
			continue
		}

		gd, ok := d.Read.Packet.(*svc.GameData)
		if !ok {
			continue
		}

		for _, cmd := range gd.Commands {
			switch c := cmd.(type) {
			case *updateuserinfo.Command:
				p.handleUpdateUserinfo(c)
			case *stufftext.Command:
				p.handleStufftext(c)
			case *updatepl.Command:
				p.handleUpdatePL(c)
			case *updateping.Command:
				p.handleUpdatePing(c)
			}
		}
	}

	return p.Events, nil
}

func (p *Parser) handleUpdateUserinfo(cmd *updateuserinfo.Command) {
	is := infostring.Parse(cmd.UserInfo)
	name := ascii.Parse(is.Get("name"))
	if name == "" {
		return
	}

	if _, exists := p.players[cmd.PlayerIndex]; !exists {
		p.players[cmd.PlayerIndex] = &Player{}
	}

	pl := p.players[cmd.PlayerIndex]
	pl.Name = name
	pl.Team = strings.TrimSpace(ascii.Parse(is.Get("team")))
}

func (p *Parser) handleStufftext(cmd *stufftext.Command) {
	for _, a := range args.Parse(cmd.String) {
		if a.Cmd == "//ktx" && len(a.Args) == 1 && a.Args[0] == "matchstart" {
			p.isStarted = true
		}
	}
}

func (p *Parser) handleUpdatePL(cmd *updatepl.Command) {
	pl, exists := p.players[cmd.PlayerIndex]
	if !exists {
		return
	}

	p.Events = append(p.Events,
		&PacketLoss{name: pl.Name, timestamp: p.elapsed, lossage: int16(cmd.PL)})
}

func (p *Parser) handleUpdatePing(cmd *updateping.Command) {
	pl, exists := p.players[cmd.PlayerIndex]
	if !exists {
		return
	}

	p.Events = append(p.Events,
		&Ping{name: pl.Name, timestamp: p.elapsed, ping: cmd.Ping})
}
