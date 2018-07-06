// state is inspired by https://www.youtube.com/watch?v=pBK-lb-k-rs&list=PLEETnX-uPtBXP_B2yupUKlflXBznWIlL5&index=4
package core

import (
	"fmt"
	"time"

	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/debug"
	"github.com/stojg/graphics/lib/input"
	"github.com/stojg/graphics/lib/rendering"
)

func NewEngine(width, height int, title string, l components.Logger) (*Engine, error) {

	window, err := NewWindow(width, height, title, false)
	if err != nil {
		return nil, err
	}

	input.SetWindow(window.Instance())

	renderer := rendering.New(window.viewPortWidth, window.viewPortHeight, l)
	scene := NewScene()

	engine := &Engine{
		scene:    scene,
		window:   window,
		renderer: renderer,
		logger:   l,
	}
	engine.scene.SetState(renderer.State())

	return engine, nil

}

type Engine struct {
	window    *Window
	scene     *Scene
	renderer  *rendering.Renderer
	isRunning bool
	logger    components.Logger
}

func (e *Engine) Renderer() components.Renderer {
	return e.renderer
}

func (m *Engine) Start() {
	m.run()
}

func (m *Engine) Width() int {
	return m.window.viewPortWidth
}

func (m *Engine) Height() int {
	return m.window.viewPortHeight
}

func (m *Engine) AddObject(object *GameObject) {
	m.scene.AddObject(object)
}

func (m *Engine) AddTerrain(object *GameObject) {
	m.scene.AddTerrain(object)
}

func (m *Engine) run() {
	m.isRunning = true

	var renderFrames int
	var frameCounter time.Duration

	var t time.Duration
	var dt = time.Millisecond

	currentTime := time.Now()
	var accumulator time.Duration

	defer m.window.Close()

	for m.isRunning {

		newTime := time.Now()
		frameTime := newTime.Sub(currentTime)
		currentTime = newTime

		accumulator += frameTime

		for accumulator >= frameTime {
			if m.window.ShouldClose() {
				m.isRunning = false
			}
			input.Update()
			m.scene.Input(dt)
			m.scene.Update(dt)
			accumulator -= dt
			t += dt
			frameCounter += dt
		}

		m.scene.Render(m.renderer)
		m.window.Maintenance()

		renderFrames++

		if frameCounter >= time.Second*5 {
			//fps := renderFrames/5)
			secondsPerFrame := (time.Second * 5 / time.Duration(renderFrames)).Seconds()
			msPerFrame := secondsPerFrame * 1000
			percent := (msPerFrame / (1000 / 60)) * 100
			dc := debug.GetDrawcalls() / uint64(renderFrames)
			ss := debug.GetShaderSwitches() / uint64(renderFrames)
			us := debug.GetUniformSet() / uint64(renderFrames)
			logLine := fmt.Sprintf("%0.1fms - %0.0f%%, %d draw calls, %d shader switches, %d uniform updates", msPerFrame, percent, dc, ss, us)
			fmt.Println(logLine)
			m.logger.Println(logLine)
			renderFrames = 0
			frameCounter = 0
		}
	}
}
