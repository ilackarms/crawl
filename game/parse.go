package game

import (
	tl "github.com/ilackarms/termloop"
	"encoding/json"
	"github.com/emc-advanced-dev/pkg/errors"
)

func SerializeLevel(level Level) (levelData, error) {
	drawables := make([]drawableData, len(level.Entities))
	for _, entity := range level.Entities {
		var drawable drawableData
		switch entity.(type) {
		case *tl.Entity:
			drawable.Type = tl.DrawableType_Entity
		case *tl.Text:
			drawable.Type = tl.DrawableType_Text
		case *tl.Rectangle:
			drawable.Type = tl.DrawableType_Rectangle
		case *PlayerRep:
			drawable.Type = DrawableType_PlayerRep
		default:
			drawable.Type = tl.DrawableType_Custom
		}
		data, err := json.Marshal(entity)
		if err != nil {
			return levelData{}, errors.New("could not convert entity to json", err)
		}
		drawable.Data = data
		drawables = append(drawables, drawable)
	}
	ld := levelData{
		UUID: level.UUID,
		Drawables: drawables,
		Bg: level.Bg,
		Offsetx: level.Offsetx,
		Offsety: level.Offsety,
	}
	return ld, nil
}

//returned level has no callback
func DeserializeLevel(ld levelData) (*Level, error) {
	level := &Level{
		BaseLevel: tl.NewBaseLevel(ld.Bg),
	}
	level.Offsetx = ld.Offsetx
	level.Offsety = ld.Offsety
	level.UUID = ld.UUID
	for _, drawable := range ld.Drawables {
		var d tl.Drawable
		switch drawable.Type {
		case tl.DrawableType_Entity:
			var e *tl.Entity
			if err := json.Unmarshal(drawable.Data, e); err != nil {
				return nil, errors.New("unmarshalling "+string(drawable.Data)+" to entity", err)
			}
			d = e
		case tl.DrawableType_Rectangle:
			var r *tl.Rectangle
			if err := json.Unmarshal(drawable.Data, r); err != nil {
				return nil, errors.New("unmarshalling "+string(drawable.Data)+" to rectangle", err)
			}
			d = r
		case tl.DrawableType_Text:
			var t *tl.Text
			if err := json.Unmarshal(drawable.Data, t); err != nil {
				return nil, errors.New("unmarshalling "+string(drawable.Data)+" to text", err)
			}
			d = t
		case DrawableType_PlayerRep:
			playerRep, err := DeserializePlayerRep(drawable.Data)
			if err != nil {
				return nil, err
			}
			d = playerRep
		default:
			return nil, errors.New("unsupported drawable type: "+string(drawable.Data), nil)
		}
		level.AddEntity(d)
	}
	return level, nil
}