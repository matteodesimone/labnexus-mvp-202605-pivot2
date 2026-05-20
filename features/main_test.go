package features

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/cucumber/godog"
)

// binaryPath è il path al binario `labnexus` costruito una volta in TestMain.
var binaryPath string

// repoRoot è la radice del repository (la dir parent di `features/`).
var repoRoot string

var godogOpts = godog.Options{
	Format: "pretty",
	Paths:  []string{"."},
	Tags:   "~@manual", // esclude gli scenari hardware-only (Finder/Gatekeeper/drag&drop)
	Strict: false,
}

func init() {
	godog.BindCommandLineFlags("godog.", &godogOpts)
}

// TestMain: costruisce il binario, poi avvia la suite BDD.
func TestMain(m *testing.M) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("getwd: %v", err)
	}
	repoRoot = filepath.Dir(wd)

	tmp, err := os.MkdirTemp("", "labnexus-bin-")
	if err != nil {
		log.Fatalf("mkdir tmp: %v", err)
	}
	defer os.RemoveAll(tmp)

	binaryPath = filepath.Join(tmp, "labnexus")
	build := exec.Command("go", "build", "-o", binaryPath, "./cmd/labnexus")
	build.Dir = repoRoot
	build.Env = os.Environ()
	if out, err := build.CombinedOutput(); err != nil {
		log.Fatalf("build failed: %v\n%s", err, string(out))
	}

	godogStatus := godog.TestSuite{
		Name:                "labnexus-fetta1",
		ScenarioInitializer: InitializeScenario,
		Options:             &godogOpts,
	}.Run()
	if godogStatus != 0 {
		fmt.Fprintf(os.Stderr, "godog: %d failures\n", godogStatus)
	}
	// Esegue anche i test Go standard nel pacchetto features (es. unit test
	// sulle assertion estratte in assertions_fetta2_test.go — review loop 1
	// H4 fix). Prima questo `m.Run()` mancava e i test venivano saltati.
	testStatus := m.Run()
	if godogStatus != 0 {
		os.Exit(godogStatus)
	}
	os.Exit(testStatus)
}

// InitializeScenario registra hook before/after + tutti gli step.
func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(c context.Context, sc *godog.Scenario) (context.Context, error) {
		s, err := newScenarioState()
		if err != nil {
			return c, err
		}
		return context.WithValue(c, stateKey{}, s), nil
	})
	ctx.After(func(c context.Context, sc *godog.Scenario, scErr error) (context.Context, error) {
		if s, ok := c.Value(stateKey{}).(*scenarioState); ok {
			s.teardown()
		}
		return c, nil
	})
	registerSteps(ctx)
}
