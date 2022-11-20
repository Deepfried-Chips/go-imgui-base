package main

import (
	"fmt"
	"github.com/Deepfried-Chips/go-imgui-base/platforms"
	"github.com/Deepfried-Chips/go-imgui-base/renderers"
	imgui "github.com/inkyblackness/imgui-go/v4"
	"os"
	"time"
)

const sleepDuration = time.Millisecond * 25

type Platform interface {
	// ShouldStop is regularly called as the abort condition for the program loop.
	ShouldStop() bool
	// ProcessEvents is called once per render loop to dispatch any pending events.
	ProcessEvents()
	// DisplaySize returns the dimension of the display.
	DisplaySize() [2]float32
	// FramebufferSize returns the dimension of the framebuffer.
	FramebufferSize() [2]float32
	// NewFrame marks the begin of a render pass. It must update the imgui IO state according to user input (mouse, keyboard, ...)
	NewFrame()
	// PostRender marks the completion of one render pass. Typically this causes the display buffer to be swapped.
	PostRender()
	// ClipboardText returns the current text of the clipboard, if available.
	ClipboardText() (string, error)
	// SetClipboardText sets the text as the current text of the clipboard.
	SetClipboardText(text string)
}

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
	platform.Window.SetTitle("SDL2 + OpenGL3 Base")
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
		imgui.SetNextWindowPosV(imgui.Vec2{X: demoX, Y: demoY}, imgui.ConditionFirstUseEver, imgui.Vec2{})

		showDemoWindow := true
		imgui.ShowDemoWindow(&showDemoWindow)

		imgui.Render()

		renderer.PreRender([3]float32{0.0, 0.0, 0.0})

		renderer.Render(platform.DisplaySize(), platform.FramebufferSize(), imgui.RenderedDrawData())
		platform.PostRender()

		<-time.After(sleepDuration)
	}
}
