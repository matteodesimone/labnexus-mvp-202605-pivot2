package config_test

import (
	"testing"

	"github.com/labnexus/labnexus/internal/config"
)

// ResolvePDFEnabled implementa la gerarchia di attivazione PDF:
//
//	CLI flag --pdf/--no-pdf   →  override per-run (priorità massima)
//	_labnexus.toml [pdf]      →  override per-capability
//	labnexus.config.toml [pdf]→  default globale
//	default hardcoded         →  false (conservativo)
//
// Tutti i parametri sono *bool: nil = non setted, &true = on, &false = off.

func boolPtr(b bool) *bool { return &b }

func TestResolvePDFEnabled_CLIWinsOverJobAndMaster(t *testing.T) {
	got := config.ResolvePDFEnabled(boolPtr(false), boolPtr(true), boolPtr(true))
	if got != false {
		t.Errorf("CLI=false deve vincere su job=true e master=true, got %v", got)
	}
}

func TestResolvePDFEnabled_CLITrueWinsOverJobFalse(t *testing.T) {
	got := config.ResolvePDFEnabled(boolPtr(true), boolPtr(false), boolPtr(false))
	if got != true {
		t.Errorf("CLI=true deve vincere su job=false e master=false, got %v", got)
	}
}

func TestResolvePDFEnabled_JobWinsOverMasterWhenCLIAbsent(t *testing.T) {
	got := config.ResolvePDFEnabled(nil, boolPtr(false), boolPtr(true))
	if got != false {
		t.Errorf("job=false deve vincere su master=true quando CLI assente, got %v", got)
	}
}

func TestResolvePDFEnabled_MasterUsedWhenCLIAndJobAbsent(t *testing.T) {
	got := config.ResolvePDFEnabled(nil, nil, boolPtr(true))
	if got != true {
		t.Errorf("master=true deve essere usato quando CLI e job assenti, got %v", got)
	}
}

func TestResolvePDFEnabled_DefaultFalseWhenAllAbsent(t *testing.T) {
	got := config.ResolvePDFEnabled(nil, nil, nil)
	if got != false {
		t.Errorf("default deve essere false quando tutto assente, got %v", got)
	}
}
