package game

import (
	"fmt"
	tl "github.com/ilackarms/termloop"
)

type Inventory struct {
	StrText    *tl.Text
	DexText    *tl.Text
	IntText    *tl.Text
	HPText     *tl.Text
	ManaText   *tl.Text
	ACText     *tl.Text
	DamageText *tl.Text

	ArmorText *tl.Text
	HandsText *tl.Text

	BackpackText *tl.Text

	Background *tl.Rectangle
}

func NewInventory() *Inventory {
	return &Inventory{
		Background: tl.NewRectangle(0, 0, 0, 0, tl.RgbTo256Color(0, 0, 0)),

		StrText:      newText(),
		DexText:      newText(),
		IntText:      newText(),
		HPText:       newText(),
		ManaText:     newText(),
		ACText:       newText(),
		DamageText:   newText(),
		ArmorText:    newText(),
		HandsText:    newText(),
		BackpackText: newText(),
	}
}

func (i *Inventory) Draw(screen *tl.Screen, offsetX, offsetY int, stats Stats) {
	w, h := screen.Size()
	i.Background.SetPosition(-1*offsetX, -1*offsetY)
	i.Background.SetSize(w, h)
	i.Background.Draw(screen)

	drawText(i.HPText, 0-offsetX, 0-offsetY, fmt.Sprintf("HP: %v/%v", stats.CurrHP, stats.MaxHP), screen)
	drawText(i.ManaText, 0-offsetX, 1-offsetY, fmt.Sprintf("MANA: %v/%v", stats.CurrMana, stats.MaxMana), screen)

	drawText(i.StrText, 0-offsetX, 3-offsetY, fmt.Sprintf("STR: %v", stats.Str), screen)
	drawText(i.DexText, 0-offsetX, 4-offsetY, fmt.Sprintf("DEX: %v", stats.Dex), screen)
	drawText(i.IntText, 0-offsetX, 5-offsetY, fmt.Sprintf("INT: %v", stats.Int), screen)

	drawText(i.ACText, 0-offsetX, 7-offsetY, fmt.Sprintf("AC: %v", stats.AC), screen)
	drawText(i.DamageText, 0-offsetX, 8-offsetY, fmt.Sprintf("DAMAGE: %v", stats.Damage), screen)

	drawText(i.ArmorText, 0-offsetX, 10-offsetY, fmt.Sprintf("ARMOR: %s", stats.Armor), screen)
	drawText(i.HandsText, 0-offsetX, 11-offsetY, fmt.Sprintf("HANDS: %s", stats.Armor), screen)

	drawText(i.BackpackText, 0-offsetX, 13-offsetY, fmt.Sprintf("BACKPACK: "), screen)
	for i, item := range stats.Backpack {
		drawText(newText(), 2-offsetX, 14+i-offsetY, fmt.Sprintf("%v", item), screen)
	}
}

func drawText(text *tl.Text, x, y int, val string, screen *tl.Screen) {
	text.SetPosition(x, y)
	text.SetText(val)
	text.Draw(screen)
}

func newText() *tl.Text {
	return tl.NewText(0, 0, "", tl.ColorWhite, tl.RgbTo256Color(0, 0, 0))
}
