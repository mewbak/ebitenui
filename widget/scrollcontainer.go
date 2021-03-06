package widget

import (
	img "image"
	"math"

	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/input"

	"github.com/hajimehoshi/ebiten"
)

type ScrollContainer struct {
	ScrollLeft float64
	ScrollTop  float64

	image               *ScrollContainerImage
	content             HasWidget
	padding             Insets
	stretchContentWidth bool

	widget    *Widget
	renderBuf *image.BufferedImage
	maskedBuf *image.BufferedImage
}

type ScrollContainerOpt func(s *ScrollContainer)

type ScrollContainerImage struct {
	Idle     *image.NineSlice
	Disabled *image.NineSlice
	Mask     *image.NineSlice
}

const ScrollContainerOpts = scrollContainerOpts(true)

type scrollContainerOpts bool

func NewScrollContainer(opts ...ScrollContainerOpt) *ScrollContainer {
	s := &ScrollContainer{
		widget: NewWidget(),

		renderBuf: &image.BufferedImage{},
		maskedBuf: &image.BufferedImage{},
	}

	for _, o := range opts {
		o(s)
	}

	return s
}

func (o scrollContainerOpts) WithLayoutData(ld interface{}) ScrollContainerOpt {
	return func(s *ScrollContainer) {
		WidgetOpts.WithLayoutData(ld)(s.widget)
	}
}

func (o scrollContainerOpts) WithImage(i *ScrollContainerImage) ScrollContainerOpt {
	return func(s *ScrollContainer) {
		s.image = i
	}
}

func (o scrollContainerOpts) WithContent(c HasWidget) ScrollContainerOpt {
	return func(s *ScrollContainer) {
		s.content = c
	}
}

func (o scrollContainerOpts) WithPadding(p Insets) ScrollContainerOpt {
	return func(s *ScrollContainer) {
		s.padding = p
	}
}

func (o scrollContainerOpts) WithStretchContentWidth() ScrollContainerOpt {
	return func(s *ScrollContainer) {
		s.stretchContentWidth = true
	}
}

func (s *ScrollContainer) GetWidget() *Widget {
	return s.widget
}

func (s *ScrollContainer) SetLocation(rect img.Rectangle) {
	s.widget.Rect = rect
}

func (s *ScrollContainer) PreferredSize() (int, int) {
	if s.content == nil {
		return 50, 50
	}

	p, ok := s.content.(PreferredSizer)
	if !ok {
		return 50, 50
	}

	w, h := p.PreferredSize()
	return w + s.padding.Dx(), h + s.padding.Dy()
}

func (s *ScrollContainer) SetupInputLayer(def input.DeferredSetupInputLayerFunc) {
	s.content.GetWidget().ElevateToNewInputLayer(&input.Layer{
		DebugLabel: "scroll container content",
		EventTypes: input.LayerEventTypeAll ^ input.LayerEventTypeWheel,
		BlockLower: true,
		FullScreen: false,
		RectFunc:   s.ContentRect,
	})

	if il, ok := s.content.(input.InputLayerer); ok {
		il.SetupInputLayer(def)
	}
}

func (s *ScrollContainer) Render(screen *ebiten.Image, def DeferredRenderFunc) {
	s.clampScroll()

	s.widget.Render(screen, def)

	s.draw(screen)

	s.renderContent(screen, def)
}

func (s *ScrollContainer) draw(screen *ebiten.Image) {
	i := s.image.Idle
	if s.widget.Disabled {
		if s.image.Disabled != nil {
			i = s.image.Disabled
		}
	}

	if i != nil {
		i.Draw(screen, s.widget.Rect.Dx(), s.widget.Rect.Dy(), func(opts *ebiten.DrawImageOptions) {
			s.widget.drawImageOptions(opts)
			s.drawImageOptions(opts)
		})
	}
}

func (s *ScrollContainer) drawImageOptions(opts *ebiten.DrawImageOptions) {
	if s.widget.Disabled && s.image.Disabled == nil {
		opts.ColorM.Scale(1, 1, 1, 0.35)
	}
}

func (s *ScrollContainer) renderContent(screen *ebiten.Image, def DeferredRenderFunc) {
	if s.content == nil {
		return
	}

	r, ok := s.content.(Renderer)
	if !ok {
		return
	}

	if l, ok := s.content.(Locateable); ok {
		var cw int
		var ch int
		if p, ok := s.content.(PreferredSizer); ok {
			cw, ch = p.PreferredSize()
		} else {
			cw, ch = 50, 50
		}

		crect := s.ContentRect()
		if s.stretchContentWidth {
			if cw < crect.Dx() {
				cw = crect.Dx()
			}
		}

		rect := img.Rect(0, 0, cw, ch)
		rect = rect.Add(s.widget.Rect.Min)
		rect = rect.Add(img.Point{s.padding.Left, s.padding.Top})

		rect = rect.Sub(img.Point{int(math.Round(float64(cw-crect.Dx()) * s.ScrollLeft)), int(math.Round(float64(ch-crect.Dy()) * s.ScrollTop))})

		l.SetLocation(rect)
		if r, ok := s.content.(Relayoutable); ok {
			r.RequestRelayout()
		}
	}

	w, h := screen.Size()

	s.renderBuf.Width, s.renderBuf.Height = w, h
	renderBuf := s.renderBuf.Image()
	_ = renderBuf.Clear()

	s.maskedBuf.Width, s.maskedBuf.Height = w, h
	maskedBuf := s.maskedBuf.Image()
	_ = maskedBuf.Clear()

	r.Render(renderBuf, def)

	s.image.Mask.Draw(maskedBuf, s.widget.Rect.Dx()-s.padding.Dx(), s.widget.Rect.Dy()-s.padding.Dy(), func(opts *ebiten.DrawImageOptions) {
		opts.GeoM.Translate(float64(s.widget.Rect.Min.X+s.padding.Left), float64(s.widget.Rect.Min.Y+s.padding.Top))
		opts.CompositeMode = ebiten.CompositeModeCopy
	})

	_ = maskedBuf.DrawImage(renderBuf, &ebiten.DrawImageOptions{
		CompositeMode: ebiten.CompositeModeSourceIn,
	})

	_ = screen.DrawImage(maskedBuf, nil)
}

func (s *ScrollContainer) ContentRect() img.Rectangle {
	return s.padding.Apply(s.widget.Rect)
}

func (s *ScrollContainer) clampScroll() {
	if s.ScrollTop < 0 {
		s.ScrollTop = 0
	} else if s.ScrollTop > 1 {
		s.ScrollTop = 1
	}

	if s.ScrollLeft < 0 {
		s.ScrollLeft = 0
	} else if s.ScrollLeft > 1 {
		s.ScrollLeft = 1
	}
}
