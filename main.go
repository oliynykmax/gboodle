package main

import (
	//"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var (
	White  = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} // W
	Red    = color.NRGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF} // R
	Green  = color.NRGBA{R: 0x00, G: 0x80, B: 0x00, A: 0xFF} // G
	Blue   = color.NRGBA{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF} // B
	Orange = color.NRGBA{R: 0xFF, G: 0xA5, B: 0x00, A: 0xFF} // O
	Yellow = color.NRGBA{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF} // Y
)

type CubeState [54]color.NRGBA

func rotateLeft[T any](s []T) {

	firstElement := s[0]

	for i := 0; i < len(s)-1; i++ {
		s[i] = s[i+1]
	}
	s[len(s)-1] = firstElement
}

func doU(s CubeState) CubeState {
	U := 0
	L := 9
	F := 18
	R := 27
	B := 36
	//D := 45
	//R1 := 2
	//R2 := 5
	//R3 := 8

	s = rotateSide(U, s)
	stack := make([]color.NRGBA, 12)
	stack[0] = s[F]
	stack[1] = s[F+1]
	stack[2] = s[F+2]
	stack[3] = s[R]
	stack[4] = s[R+1]
	stack[5] = s[R+2]
	stack[6] = s[B]
	stack[7] = s[B+1]
	stack[8] = s[B+2]
	stack[9] = s[L]
	stack[10] = s[L+1]
	stack[11] = s[L+2]
	rotateLeft(stack)
	rotateLeft(stack)
	rotateLeft(stack)
	s[F] = stack[0]
	s[F+1] = stack[1]
	s[F+2] = stack[2]
	s[R] = stack[3]
	s[R+1] = stack[4]
	s[R+2] = stack[5]
	s[B] = stack[6]
	s[B+1] = stack[7]
	s[B+2] = stack[8]
	s[L] = stack[9]
	s[L+1] = stack[10]
	s[L+2] = stack[11]
	return s
}

func rotateSide(baseIndex int, s CubeState) CubeState {
	tempC := s[baseIndex]
	s[baseIndex] = s[baseIndex+6]   // 33
	s[baseIndex+6] = s[baseIndex+8] // 35
	s[baseIndex+8] = s[baseIndex+2] // 29
	s[baseIndex+2] = tempC
	tempE := s[baseIndex+1]
	s[baseIndex+1] = s[baseIndex+3] // 30
	s[baseIndex+3] = s[baseIndex+7] // 34
	s[baseIndex+7] = s[baseIndex+5] // 32
	s[baseIndex+5] = tempE
	return s
}

func doR(s CubeState) CubeState {
	U := 0
	//L := 9
	F := 18
	R := 27
	B := 36
	D := 45
	R1 := 2
	R2 := 5
	R3 := 8

	temp1 := s[U+2]
	temp2 := s[U+2+3]
	temp3 := s[U+2+3+3]
	faces := []int{U, F, D, B}
	for i, face := range faces {
		if i == 2 {
			break
		}
		s[face+R1] = s[faces[i+1]+R1]
		s[face+R2] = s[faces[i+1]+R2]
		s[face+R3] = s[faces[i+1]+R3]
	}
	s[D+R1] = s[B-2+R1]
	s[D+R2] = s[B-2+R2]
	s[D+R3] = s[B-2+R3]
	s[B] = temp1
	s[B+3] = temp2
	s[B+3+3] = temp3
	s = rotateSide(R, s)
	return s
}

func NewSolvedCube(scramble string) CubeState {
	moves := strings.Split(scramble, " ")
	//U := 0
	////L := 9
	//F := 18
	////R := 27
	//B := 36
	//D := 45
	//R1 := 2
	//R2 := 5
	//R3 := 8
	var s CubeState
	// U Face
	for i := 0; i < 9; i++ {
		s[i] = White
	}
	// L Face
	for i := 9; i < 18; i++ {
		s[i] = Orange
	}
	// F Face
	for i := 18; i < 27; i++ {
		s[i] = Green
	}
	// R Face
	for i := 27; i < 36; i++ {
		s[i] = Red
	}
	// B Face
	for i := 36; i < 45; i++ {
		s[i] = Blue
	}
	// D Face
	for i := 45; i < 54; i++ {
		s[i] = Yellow
	}
	for _, move := range moves {
		switch move {
		case "R":
			s = doR(s)
		case "R'":
			s = doR(s)
			s = doR(s)
			s = doR(s)
		case "R2":
			s = doR(s)
			s = doR(s)
			//case "U":
			//	s = doU(s)

		}
	}
	return s
}

// getScramble generates a new Rubik's cube scramble string.
func getScramble() string {
	movesCube := [18]string{"R", "R'", "R2", "L", "L'", "L2",
		"U", "U'", "U2", "D", "D'", "D2",
		"F", "F'", "F2", "B", "B'", "B2"}
	var selected []string
	previousIndex := -1
	for len(selected) < rand.Intn(24-20+1)+20 {
		randomIndex := rand.Intn(18)
		if previousIndex/3 != randomIndex/3 {
			selected = append(selected, movesCube[randomIndex])
		}
		previousIndex = randomIndex
	}
	return strings.Join(selected, " ")
}

type CubeNetWidget struct {
	State       CubeState
	StickerSize unit.Dp
	StickerGap  unit.Dp
}

func (c *CubeNetWidget) Layout(gtx layout.Context) layout.Dimensions {
	stickerSizePx := gtx.Dp(c.StickerSize)
	stickerGapPx := gtx.Dp(c.StickerGap)
	faceSizePx := stickerSizePx*3 + stickerGapPx*2

	faceOffsets := map[string]image.Point{
		"U": {X: faceSizePx + stickerGapPx, Y: 0},
		"L": {X: 0, Y: faceSizePx + stickerGapPx},
		"F": {X: faceSizePx + stickerGapPx, Y: faceSizePx + stickerGapPx},
		"R": {X: 2*(faceSizePx+stickerGapPx) + stickerGapPx, Y: faceSizePx + stickerGapPx},
		"B": {X: 3*(faceSizePx+stickerGapPx) + 2*stickerGapPx, Y: faceSizePx + stickerGapPx},
		"D": {X: faceSizePx + stickerGapPx, Y: 2*(faceSizePx+stickerGapPx) + 2*stickerGapPx},
	}

	maxWidth := 0
	maxHeight := 0
	for _, offset := range faceOffsets {
		if offset.X+faceSizePx > maxWidth {
			maxWidth = offset.X + faceSizePx
		}
		if offset.Y+faceSizePx > maxHeight {
			maxHeight = offset.Y + faceSizePx
		}
	}
	faceStickerRanges := map[string][2]int{
		"U": {0, 9},   // Indices 0-8
		"L": {9, 18},  // Indices 9-17
		"F": {18, 27}, // Indices 18-26
		"R": {27, 36}, // Indices 27-35
		"B": {36, 45}, // Indices 36-44
		"D": {45, 54}, // Indices 45-53
	}
	for faceName, faceOffset := range faceOffsets {
		startIndex := faceStickerRanges[faceName][0]
		for i := 0; i < 9; i++ {
			row := i / 3
			col := i % 3

			stickerColor := c.State[startIndex+i]

			x := faceOffset.X + col*(stickerSizePx+stickerGapPx)
			y := faceOffset.Y + row*(stickerSizePx+stickerGapPx)

			stack := clip.Rect{Min: image.Pt(x, y), Max: image.Pt(x+stickerSizePx, y+stickerSizePx)}.Push(gtx.Ops)
			paint.ColorOp{Color: stickerColor}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			stack.Pop()
		}
	}
	return layout.Dimensions{
		Size: image.Pt(maxWidth, maxHeight),
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	go func() {
		var w app.Window
		w.Option(app.Title("Rubik's cube scrambler generator"))
		if err := run(&w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	currentScramble := getScramble()
	var regenerateButton widget.Clickable
	var scrambleDisplayEditor widget.Editor
	scrambleDisplayEditor.SetText(currentScramble) // Set initial text
	scrambleDisplayEditor.SingleLine = true
	scrambleDisplayEditor.Alignment = text.Middle
	scrambleDisplayEditor.ReadOnly = true
	cubeState := NewSolvedCube(currentScramble)
	cubeNet := CubeNetWidget{
		State:       cubeState,
		StickerSize: unit.Dp(30),
		StickerGap:  unit.Dp(3),
	}

	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			fillBackground(gtx, color.NRGBA{A: 0xFF})

			if regenerateButton.Clicked(gtx) {
				currentScramble = getScramble()
				scrambleDisplayEditor.SetText(currentScramble)
				cubeState = NewSolvedCube(currentScramble)
				cubeNet.State = cubeState
				w.Invalidate()
			}
			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceAround,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top:    unit.Dp(20),
						Bottom: unit.Dp(10),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						editorStyle := material.Editor(theme, &scrambleDisplayEditor, "") // No hint text directly here
						editorStyle.TextSize = unit.Sp(40)
						editorStyle.Color = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
						return editorStyle.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, cubeNet.Layout)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top: unit.Dp(20),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Axis:      layout.Horizontal,
								Alignment: layout.Middle,
								Spacing:   layout.SpaceAround,
							}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return material.Button(theme, &regenerateButton, "Regenerate Scramble").Layout(gtx)
								}),
							)
						})
					})
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}

func fillBackground(gtx layout.Context, col color.NRGBA) {
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
