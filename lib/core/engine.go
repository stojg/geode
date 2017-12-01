// engine is inspired by https://www.youtube.com/watch?v=pBK-lb-k-rs&list=PLEETnX-uPtBXP_B2yupUKlflXBznWIlL5&index=4
package core

import (
	"fmt"
	"time"

	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/input"
	"github.com/stojg/graphics/lib/rendering"
)

func NewEngine(width, height int, title string) (*Engine, error) {

	window, err := NewWindow(width, height, title)
	if err != nil {
		return nil, err
	}

	input.SetWindow(window.Instance())

	engine := &Engine{
		game:            NewGame(),
		window:          window,
		renderingEngine: rendering.NewEngine(width, height),
	}
	engine.game.SetEngine(engine)

	return engine, nil

}

type Engine struct {
	window          *Window
	game            *Game
	renderingEngine *rendering.Engine
	isRunning       bool
}

func (m *Engine) Start() {
	if m.isRunning {
		return
	}
	m.run()
}

func (m *Engine) AddObject(object *GameObject) {
	m.game.AddObject(object)
}

func (m *Engine) Stop() {
	if !m.isRunning {
		return
	}
	m.isRunning = false
}

func (m *Engine) run() {
	m.isRunning = true

	var renderFrames int
	var frameCounter time.Duration

	var t time.Duration
	var dt = time.Millisecond

	currentTime := time.Now()
	var accumulator time.Duration

	defer m.cleanup()

	for m.isRunning {

		newTime := time.Now()
		frameTime := newTime.Sub(currentTime)
		currentTime = newTime

		accumulator += frameTime

		for accumulator >= frameTime {
			if m.window.ShouldClose() {
				m.Stop()
			}
			input.Update()
			m.game.Input(dt)
			m.game.Update(dt)
			accumulator -= dt
			t += dt
			frameCounter += dt
		}

		m.render()
		renderFrames++

		if frameCounter >= time.Second*5 {
			fmt.Printf("%0.4f\t%d\n", (time.Second * 5 / time.Duration(renderFrames)).Seconds(), renderFrames/5)
			renderFrames = 0
			frameCounter = 0
		}
	}
	m.cleanup()
}

func (m *Engine) render() {
	m.game.Render(m.renderingEngine)
	m.window.Render()
}

func (m *Engine) cleanup() {
	m.window.Close()
}

func (m *Engine) RenderingEngine() components.RenderingEngine {
	return m.renderingEngine
}
