package godogx

import (
	"bytes"
	"io"
	"sync"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/formatters"
)

var prettyFailedRegister sync.Once

// RegisterPrettyFailedFormatter adds `pretty-failed` formatter that skips output of successful scenarios
// and shows failed.
func RegisterPrettyFailedFormatter() {
	prettyFailedRegister.Do(func() {
		godog.Format("pretty-failed", "Pretty failed formatter skips successful scenarios.",
			func(s string, writer io.Writer) formatters.Formatter {
				buf := bytes.NewBuffer(nil)

				return &prettyFailedFormatter{
					PrettyFmt: godog.NewPrettyFmt(s, buf),
					w:         writer,
					buf:       buf,
				}
			})
	})
}

type prettyFailedFormatter struct {
	*godog.PrettyFmt

	buf *bytes.Buffer
	w   io.Writer
}

func (p *prettyFailedFormatter) Failed(pickle *godog.Scenario, step *godog.Step, match *formatters.StepDefinition, err error) {
	p.PrettyFmt.Failed(pickle, step, match, err)

	p.Lock.Lock()
	defer p.Lock.Unlock()
	if _, err := io.Copy(p.w, p.buf); err != nil {
		panic(err)
	}
}

func (p *prettyFailedFormatter) Pickle(scenario *godog.Scenario) {
	p.Lock.Lock()
	p.buf.Reset()
	p.Lock.Unlock()

	p.PrettyFmt.Pickle(scenario)
}

func (p *prettyFailedFormatter) Feature(f *godog.GherkinDocument, ps string, c []byte) {
	p.Lock.Lock()
	p.buf.Reset()
	p.Lock.Unlock()

	p.PrettyFmt.Feature(f, ps, c)

	p.Lock.Lock()
	defer p.Lock.Unlock()
	if _, err := io.Copy(p.w, p.buf); err != nil {
		panic(err)
	}
}

func (p *prettyFailedFormatter) Summary() {
	p.Lock.Lock()
	p.buf.Reset()
	p.Lock.Unlock()

	p.PrettyFmt.Summary()

	p.Lock.Lock()
	defer p.Lock.Unlock()
	if _, err := io.Copy(p.w, p.buf); err != nil {
		panic(err)
	}
}
