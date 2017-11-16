// core is inspired by https://www.youtube.com/watch?v=pBK-lb-k-rs&list=PLEETnX-uPtBXP_B2yupUKlflXBznWIlL5&index=4
package core

import (
	"fmt"
	"time"
)

const frameCap time.Duration = 5000

func Main(log Logger) error {
	window := NewWindow(800, 600, "games")

	InitGraphics()
	input := NewInput()

	s := NewShader("simple")

	game := &Core{
		game:  NewGame(s),
		input: input,
		win:   window,
	}

	CheckForError("core.Main [before game.Start]")

	game.Start()
	return nil
}

type Core struct {
	game      *Game
	input     *Input
	isRunning bool
	win       *Window
}

func (m *Core) Start() {
	if m.isRunning {
		return
	}
	m.run()
}

func (m *Core) Stop() {
	if !m.isRunning {
		return
	}
	m.isRunning = false
}

func (m *Core) run() {
	m.isRunning = true

	var frames int
	var frameCount time.Duration
	frameTime := time.Second / frameCap

	lastTime := time.Now()
	var unProcessedTime time.Duration

	for m.isRunning {

		render := false

		startTime := time.Now()
		passedTime := startTime.Sub(lastTime)
		lastTime = startTime

		unProcessedTime += passedTime
		frameCount += passedTime

		for unProcessedTime > frameTime {

			render = true

			unProcessedTime -= frameTime

			if m.win.ShouldClose() {
				m.Stop()
			}
			m.input.Update()
			m.game.Input(m.input)
			m.game.Update()

			if frameCount >= time.Second {
				fmt.Println(frames)
				frames = 0
				frameCount = 0
			}
		}

		if render {
			m.render()
			frames++
		} else {
			time.Sleep(time.Millisecond)
		}
	}
	m.cleanup()
}

func (m *Core) render() {
	CheckForError("Core.render [start]")
	ClearScreen()
	m.game.Render()
	m.win.Render()
	CheckForError("Core.render [end]")
}

func (m *Core) cleanup() {
	m.win.Close()
}
