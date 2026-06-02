package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labnexus/labnexus/internal/profile"
)

// Bug: il trigger override a livello job (_labnexus.toml) non veniva mai usato
// dal motore. resolveTrigger guardava solo i campi del PROFILO, ignorando
// l'override del job. Effetto: i Prompt_INPUT_Rev.00.rtf nelle cartelle lavori/
// erano file morti.
//
// Design (deciso col CTO): il trigger è XOR (testo O file) a DUE livelli.
// Il profilo propone un default (inline o file); il _labnexus.toml locale fa
// override (inline o file). Precedenza:
//   override job (trigger_prompt OPPURE trigger_prompt_file) › default profilo.
func TestResolveTrigger_Precedence(t *testing.T) {
	const profileDefaultInline = "DEFAULT inline del profilo, lungo abbastanza per il minLen di 50 caratteri."
	const profileDefaultFileBody = "DEFAULT del profilo letto da file"
	const jobOverrideInline = "OVERRIDE inline dal job via _labnexus.toml, sopra i 50 caratteri richiesti."
	const jobOverrideFileBody = "OVERRIDE dal job letto da file"

	writeFile := func(t *testing.T, dir, name, body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	t.Run("job_file_override_vince_su_default_inline_profilo", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "ovr.txt", jobOverrideFileBody)
		p := &profile.Profile{TriggerPrompt: profileDefaultInline}
		ov := triggerOverride{File: "ovr.txt"}

		got, source, err := resolveTrigger(p, ov, dir)
		if err != nil {
			t.Fatalf("resolveTrigger: %v", err)
		}
		if got != jobOverrideFileBody {
			t.Errorf("trigger: got %q, want override del job %q", got, jobOverrideFileBody)
		}
		if !strings.Contains(strings.ToLower(source), "override") || !strings.Contains(strings.ToLower(source), "job") {
			t.Errorf("source audit: got %q, deve dichiarare l'override del job", source)
		}
	})

	t.Run("job_inline_override_vince_su_default_file_profilo", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "profile_default.txt", profileDefaultFileBody)
		p := &profile.Profile{TriggerPromptFile: "profile_default.txt"}
		ov := triggerOverride{Inline: jobOverrideInline}

		got, source, err := resolveTrigger(p, ov, dir)
		if err != nil {
			t.Fatalf("resolveTrigger: %v", err)
		}
		if got != jobOverrideInline {
			t.Errorf("trigger: got %q, want override inline del job %q", got, jobOverrideInline)
		}
		if !strings.Contains(strings.ToLower(source), "override") {
			t.Errorf("source audit: got %q, deve dichiarare l'override del job", source)
		}
	})

	t.Run("nessun_override_usa_default_file_profilo", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "profile_default.txt", profileDefaultFileBody)
		p := &profile.Profile{TriggerPromptFile: "profile_default.txt"}

		got, source, err := resolveTrigger(p, triggerOverride{}, dir)
		if err != nil {
			t.Fatalf("resolveTrigger: %v", err)
		}
		if got != profileDefaultFileBody {
			t.Errorf("trigger: got %q, want default file del profilo %q", got, profileDefaultFileBody)
		}
		if !strings.Contains(strings.ToLower(source), "profilo") {
			t.Errorf("source audit: got %q, deve dichiarare il default del profilo", source)
		}
	})

	t.Run("nessun_override_usa_default_inline_profilo", func(t *testing.T) {
		dir := t.TempDir()
		p := &profile.Profile{TriggerPrompt: profileDefaultInline}

		got, source, err := resolveTrigger(p, triggerOverride{}, dir)
		if err != nil {
			t.Fatalf("resolveTrigger: %v", err)
		}
		if got != profileDefaultInline {
			t.Errorf("trigger: got %q, want default inline del profilo %q", got, profileDefaultInline)
		}
		if !strings.Contains(strings.ToLower(source), "profilo") {
			t.Errorf("source audit: got %q, deve dichiarare il default del profilo", source)
		}
	})

	t.Run("effective_trigger_file_precedenza", func(t *testing.T) {
		cases := []struct {
			name string
			p    *profile.Profile
			ov   triggerOverride
			want string
		}{
			{"job_file_vince", &profile.Profile{TriggerPromptFile: "prof.txt"}, triggerOverride{File: "job.rtf"}, "job.rtf"},
			{"job_inline_nessun_file", &profile.Profile{TriggerPromptFile: "prof.txt"}, triggerOverride{Inline: "x"}, ""},
			{"profilo_file", &profile.Profile{TriggerPromptFile: "prof.txt"}, triggerOverride{}, "prof.txt"},
			{"profilo_inline_nessun_file", &profile.Profile{TriggerPrompt: "x"}, triggerOverride{}, ""},
		}
		for _, tc := range cases {
			if got := effectiveTriggerFile(tc.p, tc.ov); got != tc.want {
				t.Errorf("%s: effectiveTriggerFile = %q, want %q", tc.name, got, tc.want)
			}
		}
	})

	// Robustness: un override whitespace-only (es. trigger_prompt_file = "   ")
	// passa la validazione XOR come "assente" (TrimSpace), quindi a runtime deve
	// fare fallback al default del profilo, NON tentare di leggere un path vuoto.
	t.Run("override_whitespace_fa_fallback_al_profilo", func(t *testing.T) {
		dir := t.TempDir()
		p := &profile.Profile{TriggerPrompt: profileDefaultInline}

		for _, ov := range []triggerOverride{{File: "   "}, {Inline: "  \n\t "}} {
			got, source, err := resolveTrigger(p, ov, dir)
			if err != nil {
				t.Fatalf("resolveTrigger(%+v): atteso fallback senza errore, got %v", ov, err)
			}
			if got != profileDefaultInline {
				t.Errorf("override %+v: got %q, want fallback al default profilo", ov, got)
			}
			if !strings.Contains(strings.ToLower(source), "profilo") {
				t.Errorf("override %+v: source %q deve dichiarare il default del profilo", ov, source)
			}
		}
	})
}
