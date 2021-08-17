package allure

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/formatters"
	"github.com/google/uuid"
)

var formatterRegister sync.Once

func RegisterAllureFormatter(resultsPath string) {
	if resultsPath == "" {
		resultsPath = "./allure-results"
	}

	formatterRegister.Do(func() {
		err := os.MkdirAll(resultsPath, 0o700)
		if err != nil {
			panic(err)
		}

		godog.Format("allure", "Allure formatter.",
			func(s string, writer io.Writer) formatters.Formatter {
				return &Formatter{
					BaseFmt: godog.NewBaseFmt(s, writer),
				}
			})
	})
}

type Formatter struct {
	res         *TestCase
	lastTime    TimestampMs
	ResultsPath string

	*godog.BaseFmt
}

func (f *Formatter) writeResult(r *TestCase) error {
	f.lastTime = getTimestampMs()

	r.Stage = "finished"
	r.Stop = f.lastTime

	j, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fmt.Sprintf("%s/%s-result.json", f.resultsPath(), r.UUID), j, 0o600)
}

func (f *Formatter) resultsPath() string {
	if f.ResultsPath != "" {
		return f.ResultsPath
	}

	return "./allure-results"
}

// TestRunStarted is triggered on test start.
// func (f *Formatter) TestRunStarted() {}

// Feature receives gherkin document.
//func (f *Formatter) Feature(feature *godog.Feature, string, []byte) {
//	feature.Feature.Name
//}

// Pickle receives scenario.
func (f *Formatter) Pickle(scenario *godog.Scenario) {
	if f.res != nil {
		if err := f.writeResult(f.res); err != nil {
			panic(err)
		}
	}

	f.lastTime = getTimestampMs()

	feature := f.Storage.MustGetFeature(scenario.Uri)

	f.res = &TestCase{
		UUID:        uuid.New().String(),
		Name:        scenario.Name,
		FullName:    scenario.Uri + ":" + scenario.Name,
		Description: scenario.Uri,
		Start:       f.lastTime,
		Labels: []Label{
			{Name: "feature", Value: feature.Feature.Name},
			{Name: "suite", Value: "features"},
			{Name: "framework", Value: "godog"},
			{Name: "language", Value: "Go"},
		},
	}
}

func getTimestampMs() TimestampMs {
	return TimestampMs(time.Now().UnixNano() / int64(time.Millisecond))
}

// Defined receives step definition.
//func (f *Formatter) Defined(*godog.Scenario, *godog.Step, *godog.StepDefinition) {
//}

func (f *Formatter) step(st *godog.Step, sd *godog.StepDefinition) Step {
	step := Step{
		Name:  st.Text,
		Stage: "finished",
		Start: f.lastTime,
		//Parameters: []Parameter{{
		//	//Name:  "text",
		//	Value: st.Text,
		//}},
	}

	if st.Argument != nil {
		if st.Argument.DocString != nil {
			mt := "text/plain"
			if st.Argument.DocString.MediaType == "json" {
				mt = "application/json"
			}

			att, err := NewAttachment("Doc", mt, f.resultsPath(), []byte(st.Argument.DocString.Content))
			if err != nil {
				panic(err)
			}

			step.Attachments = append(step.Attachments, *att)
		} else if st.Argument.DataTable != nil {
			mt := "text/csv"
			buf := bytes.NewBuffer(nil)
			c := csv.NewWriter(buf)

			for _, r := range st.Argument.DataTable.Rows {
				var rec []string
				for _, cell := range r.Cells {
					rec = append(rec, cell.Value)
				}
				c.Write(rec)
			}
			c.Flush()

			att, err := NewAttachment("Table", mt, f.resultsPath(), buf.Bytes())
			if err != nil {
				panic(err)
			}

			step.Attachments = append(step.Attachments, *att)
		}

		//if st.Argument.DocString != nil {
		//	step.Parameters = []Parameter{{
		//		//Name:  st.Argument.DocString.MediaType,
		//		Value: st.Argument.DocString.Content,
		//	}}
		//} else if st.Argument.DataTable != nil {
		//
		//}
	}

	f.lastTime = getTimestampMs()
	step.Stop = f.lastTime

	return step
}

// Passed captures passed step.
func (f *Formatter) Passed(sc *godog.Scenario, st *godog.Step, sd *godog.StepDefinition) {
	step := f.step(st, sd)
	step.Status = Passed
	f.res.Steps = append(f.res.Steps, step)
	f.res.Status = Passed
}

// Skipped captures skipped step.
func (f *Formatter) Skipped(sc *godog.Scenario, st *godog.Step, sd *godog.StepDefinition) {
	step := f.step(st, sd)
	step.Status = Skipped
	f.res.Steps = append(f.res.Steps, step)
}

// Undefined captures undefined step.
func (f *Formatter) Undefined(sc *godog.Scenario, st *godog.Step, sd *godog.StepDefinition) {
	step := f.step(st, sd)
	step.Status = Broken

	f.res.Steps = append(f.res.Steps, step)
}

// Failed captures failed step.
func (f *Formatter) Failed(sc *godog.Scenario, st *godog.Step, sd *godog.StepDefinition, err error) {
	details := &StatusDetails{
		Message: err.Error(),
	}

	step := f.step(st, sd)
	step.Status = Failed
	step.StatusDetails = details

	f.res.Steps = append(f.res.Steps, step)
	f.res.Status = Failed
	f.res.StatusDetails = details
}

// Pending captures pending step.
func (f *Formatter) Pending(*godog.Scenario, *godog.Step, *godog.StepDefinition) {
}

func (f *Formatter) writeJSON(name string, v interface{}) error {
	j, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(f.resultsPath()+"/"+name, j, 0o600)
}

func (f *Formatter) Summary() {
	if f.res != nil {
		if err := f.writeResult(f.res); err != nil {
			panic(err)
		}
	}

	e := Executor{
		Name: "GitHub Actions",
		Type: "github",
	}

	f.writeJSON("executor.json", e)
}
