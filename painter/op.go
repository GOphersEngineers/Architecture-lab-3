package painter

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/exp/shiny/screen"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool { return true }

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

// WhiteFill зафарбовує тестуру у білий колір. Може бути викоистана як Operation через OperationFunc(WhiteFill).
func WhiteFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.White, screen.Src)
}

// GreenFill зафарбовує тестуру у зелений колір. Може бути викоистана як Operation через OperationFunc(GreenFill).
func GreenFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.RGBA{G: 0xff, A: 0xff}, screen.Src)
}

type BgRect struct {
	X1 int
	Y1 int
	X2 int
	Y2 int
}

type TFigure struct {
	X int
	Y int
}

type MoveFigures struct {
	X       int
	Y       int
	Figures []*TFigure
}

func (op *BgRect) Do(texture screen.Texture) bool {
	texture.Fill(image.Rect(op.X1, op.Y1, op.X2, op.Y2), color.Black, screen.Src)
	return false
}

func (op *TFigure) Do(texture screen.Texture) bool {
	blueColor := color.RGBA{B: 255}
	texture.Fill(image.Rect(op.X-100, op.Y-100, op.X+100, op.Y-50), blueColor, draw.Src)
	texture.Fill(image.Rect(op.X-25, op.Y-50, op.X+25, op.Y+100), blueColor, draw.Src)
	return false
}

func (op *MoveFigures) Do(texture screen.Texture) bool {
	for i := range op.Figures {
		op.Figures[i].X += op.X
		op.Figures[i].Y += op.Y
	}
	return false
}

func Reset(texture screen.Texture) {
	texture.Fill(texture.Bounds(), color.Black, screen.Src)
}
