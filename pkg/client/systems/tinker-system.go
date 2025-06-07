package systems

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	ark "github.com/mlange-42/ark/ecs"

	"github.com/ouijan/ingenuity/pkg/core/ecs"
	"github.com/ouijan/ingenuity/pkg/core/ecs/components"
	"github.com/ouijan/ingenuity/pkg/core/utils"
)

type TinkerSystem struct {
	textEntities  *ark.Filter2[components.Metadata, components.Text]
	netEntities   *ark.Filter1[components.NetworkedEntity]
	netComponents *ark.Map3[components.Metadata, components.Transform2D, components.Text]
}

func (s *TinkerSystem) OnCreate(ea *ecs.EntityAdmin) {
	s.textEntities = ark.NewFilter2[components.Metadata, components.Text](&ea.World)
	s.netEntities = ark.NewFilter1[components.NetworkedEntity](&ea.World)
	s.netComponents = ark.NewMap3[components.Metadata, components.Transform2D, components.Text](
		&ea.World,
	)
}

func (s *TinkerSystem) Update(dt float32) error {
	fpsMsg := fmt.Sprintf("FPS: %d, DT: %f", rl.GetFPS(), rl.GetFrameTime())

	// TODO: this is only grabbing 1 message per frame

	// data := &NetMsg{}
	// msg, hasMsg := utils.ChanSelect(ca.client.MsgCh)
	// if hasMsg {
	// 	// log.Info("<- Received message %s", msg)
	// 	err := json.Unmarshal(msg.Payload, data)
	// 	if err != nil {
	// 		log.Error("Failed to unmarshal message: %v \n %v", err, msg.Payload)
	// 		return err
	// 	}
	// }

	query := s.textEntities.Query()
	for query.Next() {
		meta, text := query.Get()

		switch meta.Name {
		case utils.FPSDisplayName:
			text.Content = fpsMsg
			// case MessageDisplayName:
			// 	if hasMsg && data.Debug != nil {
			// 		text.Content = data.Debug.Msg
			// 	}
		}
	}

	// if hasMsg && data.EntityUpdate != nil {
	// 	query := s.netEntities.Query()
	// 	for query.Next() {
	// 		net := query.Get()
	// 		if net.Id == data.EntityUpdate.Id {
	// 			s.updateEntity(data.EntityUpdate, query.Entity())
	// 		}
	// 	}
	// }

	return nil
}

// func (s *TestClientSystem) updateEntity(update *NetMsg_EntityUpdate, entity ark.Entity) {
// 	meta, trans, text := s.netComponents.Get(entity)
// 	if update.Metadata != nil {
// 		meta.Name = update.Metadata.Name
// 		meta.Tags = update.Metadata.Tags
// 	}
// 	if update.Transform != nil {
// 		trans.X = update.Transform.X
// 		trans.Y = update.Transform.Y
// 	}
// 	if update.Text != nil {
// 		text.Content = update.Text.Content
// 		text.FontSize = update.Text.FontSize
// 		text.Colour = update.Text.Colour
// 	}
// }

func (s *TinkerSystem) OnDestroy() {
	// Cleanup if needed
}

var _ ecs.System = &TinkerSystem{}
