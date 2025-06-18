package mvdparser

import (
	"strings"

	"github.com/osm/quake/common/args"
	"github.com/osm/quake/common/ascii"
	"github.com/osm/quake/common/context"
	"github.com/osm/quake/common/infostring"
	"github.com/osm/quake/demo/mvd"
	"github.com/osm/quake/packet/command/stufftext"
	"github.com/osm/quake/packet/command/updatepl"
	"github.com/osm/quake/packet/command/updateuserinfo"
	"github.com/osm/quake/packet/svc"
)

type Parser struct {
	isStarted bool
	players   map[byte]*Player
	elapsed   float64

	PacketLoss []PacketLoss
}

type PacketLoss struct {
	Name      string
	Timestamp int
	Lossage   byte
}

type Player struct {
	Team string
	Name string
}

func New() *Parser {
	return &Parser{
		players: make(map[byte]*Player),
	}
}

func (p *Parser) Parse(data []byte) ([]PacketLoss, error) {
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
			}
		}
	}

	return p.PacketLoss, nil
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

	p.PacketLoss = append(p.PacketLoss,
		PacketLoss{Name: pl.Name, Timestamp: int(p.elapsed), Lossage: cmd.PL})
}
