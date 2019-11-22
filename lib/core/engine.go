// state is inspired by https://www.youtube.com/watch?v=pBK-lb-k-rs&list=PLEETnX-uPtBXP_B2yupUKlflXBznWIlL5&index=4
package core

import (
	"fmt"
	"time"

	"github.com/stojg/geode/lib/components"
	"github.com/stojg/geode/lib/debug"
	"github.com/stojg/geode/lib/input"
	"github.com/stojg/geode/lib/rendering"
)

func NewEngine(width, height int, title string, l components.Logger) (*Engine, error) {
	window, err := NewWindow(width, height, title, false)
	if err != nil {
		return nil, err
	}

	input.SetWindow(window.Instance())

	renderer := rendering.New(window.width, window.height, window.viewPortWidth, window.viewPortHeight, l)
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

func (e *Engine) Start() {
	e.run()
}

func (e *Engine) Width() int {
	return e.window.width
}

func (e *Engine) Height() int {
	return e.window.height
}

func (e *Engine) AddObject(object components.Object) {
	e.scene.AddObject(object)
}

func (e *Engine) run() {
	e.isRunning = true
	defer e.window.Close()

	const updateStep = 8 * time.Millisecond

	var total, accumulator, debugTimer time.Duration
	var updatedFrames, renderedFrames int

	currentTime := time.Now()

	for e.isRunning {
		newTime := time.Now()
		frameTime := newTime.Sub(currentTime)
		currentTime = newTime

		accumulator += frameTime
		debugTimer += frameTime

		if e.window.ShouldClose() {
			e.isRunning = false
		}

		// The renderer produces time and the simulation consumes it in discrete dt sized steps.
		for accumulator >= updateStep {
			// @todo why all this input update, scene input and then scene update?
			// need to figure out where this would fit in a "physics integrate loop" way
			input.Update()
			e.scene.Input(updateStep)
			e.scene.Update(updateStep)
			accumulator -= updateStep
			total += updateStep
			updatedFrames++
		}

		e.scene.Render(e.renderer)
		renderedFrames++

		e.window.Maintenance()

		if debugTimer >= time.Second*5 {
			fps := renderedFrames / 5.0
			ups := updatedFrames / 5.0
			secondsPerFrame := (time.Second * 5 / time.Duration(renderedFrames)).Seconds()
			msPerFrame := secondsPerFrame * 1000
			percent := (msPerFrame / (1000 / 60)) * 100
			dc := debug.GetDrawcalls() / uint64(renderedFrames)
			ss := debug.GetShaderSwitches() / uint64(renderedFrames)
			us := debug.GetUniformSet() / uint64(renderedFrames)
			vb := debug.GetVertexBind() / uint64(renderedFrames)
			logLine := fmt.Sprintf("%s | %0.1fms - %0.0f%% (%d fps, %d ups), %d draw calls, %d shader switches, %d uniform updates, %d vertex binds, %d particles", total, msPerFrame, percent, fps, ups, dc, ss, us, vb, debug.GetParticles())
			fmt.Println(logLine)
			e.logger.Println(logLine)
			renderedFrames = 0
			debugTimer = 0
			updatedFrames = 0
		}
	}
}
