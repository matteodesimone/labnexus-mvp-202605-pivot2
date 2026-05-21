// xls.go — Parser per Excel 97-2003 OLE binary (.xls). FR-18.
//
// Lib usata: github.com/extrame/xls (decisione 2026-05-21: testato su file
// reale Denis "M132.00_Rapporto Proficiency Testing.xls" → estrae dati ok,
// con nil-guard sui Row indici inesistenti — la lib panica senza guard).
package input

import (
	"fmt"
	"strings"

	"github.com/extrame/xls"
)

// parseXls estrae plain text da un file Excel 97-2003 OLE binary. Tutte le
// sheet vengono concatenate con separatore `--- SHEET: <name> ---`.
// Celle nil/vuote skippate.
func parseXls(path string) (s string, err error) {
	// Defense: la lib extrame/xls può panicare su row indici nil. Recovery → error.
	defer func() {
		if r := recover(); r != nil {
			s = ""
			err = fmt.Errorf("input: parse .xls %q panicked: %v", path, r)
		}
	}()
	wb, err := xls.Open(path, "utf-8")
	if err != nil {
		return "", fmt.Errorf("input: open .xls %q: %w", path, err)
	}
	var out strings.Builder
	for i := 0; i < wb.NumSheets(); i++ {
		sh := wb.GetSheet(i)
		if sh == nil {
			continue
		}
		fmt.Fprintf(&out, "--- SHEET: %s ---\n", sh.Name)
		appendSheetText(&out, sh)
		out.WriteByte('\n')
	}
	return out.String(), nil
}

// appendSheetText itera righe e celle non-nil scrivendo le celle separate da tab.
// Per ogni riga, recover da eventuali panic (extrame/xls può panicare su
// indici interni inconsistenti) e continua con la riga successiva.
func appendSheetText(out *strings.Builder, sh *xls.WorkSheet) {
	max := int(sh.MaxRow)
	for r := 0; r <= max; r++ {
		appendRowSafely(out, sh, r)
	}
}

func appendRowSafely(out *strings.Builder, sh *xls.WorkSheet, r int) {
	defer func() { _ = recover() }()
	row := sh.Row(r)
	if row == nil {
		return
	}
	first, last := row.FirstCol(), row.LastCol()
	writeRowCells(out, row, first, last)
}

func writeRowCells(out *strings.Builder, row *xls.Row, first, last int) {
	wroteCell := false
	for c := first; c < last; c++ {
		v := row.Col(c)
		if v == "" {
			continue
		}
		if wroteCell {
			out.WriteByte('\t')
		}
		out.WriteString(v)
		wroteCell = true
	}
	if wroteCell {
		out.WriteByte('\n')
	}
}
