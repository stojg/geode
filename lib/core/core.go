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
	game := &Core{
		game: NewGame(),
		input: input,
		win: window,
	}
	game.Start()
	//if err := win.Open(); err != nil {
	//	return err
	//}
	//defer win.Close()
	//
	//if err := gl.Init(); err != nil {
	//	return err
	//}
	//version := gl.GoStr(gl.GetString(gl.VERSION))
	//log.Printf("OpenGL Version %s\n", version)
	//
	//gl.Disable(gl.MULTISAMPLE)
	return nil
}

type Core struct {
	game *Game
	input *Input
	isRunning bool
	win *Window
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
	ClearScreen()
	m.game.Render()
	m.win.Render()
}

func (m *Core) cleanup() {
	m.win.Close()
}

