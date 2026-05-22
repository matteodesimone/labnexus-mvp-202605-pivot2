package config

// ResolvePDFEnabled risolve l'attivazione del PDF accoppiato all'output MD
// secondo la gerarchia "più specifico vince":
//
//	1. cli      — flag --pdf/--no-pdf (priorità massima, per-run)
//	2. job      — sezione [pdf] in _labnexus.toml (per-capability)
//	3. master   — sezione [pdf] in labnexus.config.toml (default globale)
//	4. default  — false (conservativo: PDF off salvo opt-in esplicito)
//
// Ogni parametro è *bool: nil = non setted, &true = on, &false = off.
// Il primo livello non-nil vince. Se tutti sono nil, ritorna false.
func ResolvePDFEnabled(cli, job, master *bool) bool {
	if cli != nil {
		return *cli
	}
	if job != nil {
		return *job
	}
	if master != nil {
		return *master
	}
	return false
}
