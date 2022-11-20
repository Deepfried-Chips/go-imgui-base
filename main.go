

package main

import (
	"fmt"
	"github.com/Deepfried-Chips/go-imgui-base/platforms"
	"github.com/Deepfried-Chips/go-imgui-base/renderers"
	"github.com/inkyblackness/imgui-go/v4"
	"os"
	"time"
)

const sleepDuration   = time.Millisecond * 25

type clipboard struct {
	platform Platform
}

func (board clipboard) Text() (string, error) {
	return board.platform.ClipboardText()
}

func (board clipboard) SetText(text string) {
	board.platform.SetClipboardText(text)
}

func main() {
	context := imgui.CreateContext(nil)
	defer context.Destroy()
	io := imgui.CurrentIO()

	platform, err := platforms.NewSDL(io, platforms.SDLClientAPIOpenGL3)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer platform.Dispose()

	renderer, err := renderers.NewOpenGL3(io)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer renderer.Dispose()

	for !platform.ShouldStop() {
		platform.ProcessEvents()

		platform.NewFrame()
		imgui.NewFrame()

		const demoX = 650
		const demoY = 20
		imgui.SetNextWindowPosV(imgui.Vec2(X: demoX, Y: demoY}, imgui.ConditionFirstUseEver, imgui.Vec2{})

		imgui.ShowDemoWindow(true)

		imgui.Render()

		renderer.PreRender([3]float32{0.0, 0.0, 0.0})

		renderer.Render(platform.DisplaySize(), platform.FramebufferSize(), imgui.RenderedDrawData())
		platform.PostRender()

		<-time.After(sleepDuration)
	}
}
