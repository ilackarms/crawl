package game

import (
	tl "github.com/ilackarms/termloop"
	"encoding/json"
	"github.com/emc-advanced-dev/pkg/errors"
)

func SerializeLevel(level Level) ([]byte, error) {
	drawables := make([]drawableData, len(level.Entities))
	for _, entity := range level.Entities {
		var drawable drawableData
		switch entity.(type) {
		case tl.Entity:
			drawable.Type = tl.DrawableType_Entity
		case tl.Text:
			drawable.Type = tl.DrawableType_Text
		case tl.Rectangle:
			drawable.Type = tl.DrawableType_Rectangle
		default:
			drawable.Type = tl.DrawableType_Custom
		}
		data, err := json.Marshal(entity)
		if err != nil {
			return nil, errors.New("could not convert entity to json", err)
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
	data, err := json.Marshal(ld)
	if err != nil {
		return nil, errors.New("converting level data to json", err)
	}
	return data, nil
}

//returned level has no callback
func DeserializeLevel(data []byte) (Level, error) {
	var ld levelData
	if err := json.Unmarshal(data, &ld); err != nil {
		return nil, errors.New("unmarshalling "+string(data)+" to level data", err)
	}
	level := Level{
		BaseLevel: 	tl.NewBaseLevel(ld.Bg),
	}
	level.Offsetx = ld.Offsetx
	level.Offsety = ld.Offsety
	level.UUID = ld.UUID
	for _, drawable := range ld.Drawables {
		var entity tl.Drawable
		switch drawable.Type {
		case tl.DrawableType_Entity:
			var e tl.Entity
			if err := json.Unmarshal(drawable.Data, &e); err != nil {
				return nil, errors.New("unmarshalling "+string(drawable.Data)+" to entity", err)
			}
			entity = e
		case tl.DrawableType_Rectangle:
			var r tl.Rectangle
			if err := json.Unmarshal(drawable.Data, &r); err != nil {
				return nil, errors.New("unmarshalling "+string(drawable.Data)+" to rectangle", err)
			}
			entity = r
		case tl.DrawableType_Text:
			var t tl.Text
			if err := json.Unmarshal(drawable.Data, &t); err != nil {
				return nil, errors.New("unmarshalling "+string(drawable.Data)+" to text", err)
			}
			entity = t
		case tl.DrawableType_Custom:
			return nil, errors.New("unsupported drawable type: "+string(drawable.Data), nil)
		}
		level.AddEntity(entity)
	}
	return level, nil
}