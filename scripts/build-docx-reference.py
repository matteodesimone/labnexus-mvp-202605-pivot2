#!/usr/bin/env python3
"""Build LabNexus reference.docx per pandoc --reference-doc.

Tool one-off Sprint 1.5.D. Lanciato manualmente quando serve rigenerare il
reference template (es. cambia palette, aggiungiamo nuovo style, etc.).
NON è invocato automaticamente da Makefile — il reference.docx prodotto è
committed nel repo come binary asset embedded via go:embed.

USO:
  cd <repo-root>
  python3 scripts/build-docx-reference.py

Prerequisiti:
  - pandoc presente in dist/labnexus-sprint1-darwin-arm64/bin/pandoc
    (esegui `make pdf-tools` prima se manca)
  - python3 >= 3.6

Output:
  internal/pdf/reference.docx (~12KB, sovrascritto)

Parte dal default pandoc reference.docx ed applica i design token AICertus:
- Heading colors: #010C23 (navy, vs default Word #0F4761)
- BlockText (blockquote): bordo sinistro accent #006A9C 2pt + indent + spacing 12pt
- VerbatimChar (code inline): background pale #F0F5FE
- Table firstRow: background pale #F0F5FE
- Heading1/2: bottom border accent
- Header pagina: "LabNexus | {{CAPABILITY}}" con separatore accent
- Footer pagina: "profilo: {{PROFILO}}  -  {{MODELLO}}     pag X / Y"

I placeholders {{CAPABILITY}}, {{PROFILO}}, {{MODELLO}} vengono sostituiti
dal pre-processor Go (internal/pdf/pdf.go:writeReferenceDocxWithMeta) prima
della compilazione pandoc.
"""
import subprocess
import shutil
import sys
import tempfile
import zipfile
import re
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
PANDOC = REPO_ROOT / "dist/labnexus-sprint1-darwin-arm64/bin/pandoc"
OUT = REPO_ROOT / "internal/pdf/reference.docx"

if not PANDOC.exists():
    sys.exit(f"errore: pandoc non trovato in {PANDOC}. Esegui 'make pdf-tools' prima.")

TMP = Path(tempfile.mkdtemp(prefix="labnexus-docx-build-"))
SRC = TMP / "ref-default.docx"
with SRC.open("wb") as _f:
    subprocess.run([str(PANDOC), "--print-default-data-file", "reference.docx"],
                   stdout=_f, check=True)

WORK = TMP / "ref-extract"

# Clean working dir
if WORK.exists():
    shutil.rmtree(WORK)
WORK.mkdir()

# Extract default
with zipfile.ZipFile(SRC) as zf:
    zf.extractall(WORK)

# ============================================================
# 1. word/styles.xml — apply brand tokens
# ============================================================
styles_path = WORK / "word/styles.xml"
styles = styles_path.read_text()

# Heading colors → navy #010C23 (sostituisce il default #0F4761)
styles = styles.replace('w:val="0F4761"', 'w:val="010C23"')
styles = styles.replace('w:themeColor="accent1"', 'w:themeColor="text1"')

# BlockText (blockquote): aggiungi border-left accent 2pt + padding
# Cerca: <w:style w:type="paragraph" w:styleId="BlockText"> ... </w:style>
def patch_block_text(m):
    inner = m.group(1)
    # Inserisci pBdr + indent + spacing before/after.
    # Spacing: w:before/w:after sono in 1/20 di punto (twips).
    # 240 twips = 12pt ≈ 4.2mm — spazio respiratorio prima/dopo il blockquote
    # così non si attacca ai paragrafi adiacenti.
    pbdr = (
        '<w:pBdr><w:left w:val="single" w:sz="18" w:space="6" w:color="006A9C"/></w:pBdr>'
        '<w:spacing w:before="240" w:after="240"/>'
        '<w:ind w:left="284"/>'
    )
    if "<w:pPr>" in inner:
        inner = inner.replace("<w:pPr>", "<w:pPr>" + pbdr, 1)
    else:
        inner = "<w:pPr>" + pbdr + "</w:pPr>" + inner
    return f'<w:style w:type="paragraph" w:styleId="BlockText">{inner}</w:style>'

styles = re.sub(
    r'<w:style w:type="paragraph" w:styleId="BlockText">(.*?)</w:style>',
    patch_block_text,
    styles,
    flags=re.DOTALL,
)

# VerbatimChar (code inline): aggiungi background pale
def patch_verbatim(m):
    inner = m.group(1)
    shd = '<w:shd w:val="clear" w:color="auto" w:fill="F0F5FE"/>'
    if "<w:rPr>" in inner:
        inner = inner.replace("<w:rPr>", "<w:rPr>" + shd, 1)
    else:
        inner = inner + "<w:rPr>" + shd + "</w:rPr>"
    return f'<w:style w:type="character" w:customStyle="1" w:styleId="VerbatimChar">{inner}</w:style>'

styles = re.sub(
    r'<w:style w:type="character" w:customStyle="1" w:styleId="VerbatimChar">(.*?)</w:style>',
    patch_verbatim,
    styles,
    flags=re.DOTALL,
)

# Table firstRow: shading pale per header row
def patch_table_firstrow(m):
    inner = m.group(1)
    shd = '<w:shd w:val="clear" w:color="auto" w:fill="F0F5FE"/>'
    if "<w:tcPr>" in inner:
        inner = inner.replace("<w:tcPr>", "<w:tcPr>" + shd, 1)
    return f'<w:tblStylePr w:type="firstRow">{inner}</w:tblStylePr>'

styles = re.sub(
    r'<w:tblStylePr w:type="firstRow">(.*?)</w:tblStylePr>',
    patch_table_firstrow,
    styles,
    flags=re.DOTALL,
)

# Heading2: aggiungi border-bottom accent (sottolinea sezione)
def patch_heading(style_id, color):
    def fn(m):
        inner = m.group(1)
        pbdr = f'<w:pBdr><w:bottom w:val="single" w:sz="4" w:space="2" w:color="{color}"/></w:pBdr>'
        if "<w:pPr>" in inner:
            inner = inner.replace("<w:pPr>", "<w:pPr>" + pbdr, 1)
        else:
            inner = "<w:pPr>" + pbdr + "</w:pPr>" + inner
        return f'<w:style w:type="paragraph" w:styleId="{style_id}">{inner}</w:style>'
    return fn

styles = re.sub(
    r'<w:style w:type="paragraph" w:styleId="Heading1">(.*?)</w:style>',
    patch_heading("Heading1", "006A9C"),
    styles,
    flags=re.DOTALL,
)
styles = re.sub(
    r'<w:style w:type="paragraph" w:styleId="Heading2">(.*?)</w:style>',
    patch_heading("Heading2", "006A9C"),
    styles,
    flags=re.DOTALL,
)

styles_path.write_text(styles)

# ============================================================
# 2. Header / Footer XML files
# ============================================================
HEADER_XML = '''<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:hdr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <w:p>
    <w:pPr>
      <w:pBdr><w:bottom w:val="single" w:sz="4" w:space="2" w:color="006A9C"/></w:pBdr>
      <w:spacing w:after="0"/>
    </w:pPr>
    <w:r>
      <w:rPr><w:rFonts w:ascii="Helvetica Neue" w:hAnsi="Helvetica Neue" w:cs="Helvetica Neue"/><w:b/><w:color w:val="010C23"/><w:sz w:val="22"/></w:rPr>
      <w:t>LabNexus</w:t>
    </w:r>
    <w:r>
      <w:rPr><w:color w:val="006A9C"/><w:sz w:val="22"/></w:rPr>
      <w:t xml:space="preserve">  |  </w:t>
    </w:r>
    <w:r>
      <w:rPr><w:rFonts w:ascii="Helvetica Neue" w:hAnsi="Helvetica Neue" w:cs="Helvetica Neue"/><w:color w:val="333333"/><w:sz w:val="20"/></w:rPr>
      <w:t>{{CAPABILITY}}</w:t>
    </w:r>
  </w:p>
</w:hdr>
'''

FOOTER_XML = '''<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:ftr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <w:p>
    <w:pPr>
      <w:pBdr><w:top w:val="single" w:sz="4" w:space="2" w:color="006A9C"/></w:pBdr>
      <w:tabs>
        <w:tab w:val="right" w:pos="9072"/>
      </w:tabs>
      <w:spacing w:before="0"/>
    </w:pPr>
    <w:r>
      <w:rPr><w:color w:val="AAAAAA"/><w:sz w:val="16"/></w:rPr>
      <w:t xml:space="preserve">profilo: {{PROFILO}}  -  {{MODELLO}}</w:t>
    </w:r>
    <w:r>
      <w:rPr><w:color w:val="AAAAAA"/><w:sz w:val="16"/></w:rPr>
      <w:tab/>
      <w:t xml:space="preserve">pag </w:t>
    </w:r>
    <w:fldSimple w:instr="PAGE">
      <w:r><w:rPr><w:color w:val="AAAAAA"/><w:sz w:val="16"/></w:rPr><w:t>1</w:t></w:r>
    </w:fldSimple>
    <w:r>
      <w:rPr><w:color w:val="AAAAAA"/><w:sz w:val="16"/></w:rPr>
      <w:t xml:space="preserve"> / </w:t>
    </w:r>
    <w:fldSimple w:instr="NUMPAGES">
      <w:r><w:rPr><w:color w:val="AAAAAA"/><w:sz w:val="16"/></w:rPr><w:t>1</w:t></w:r>
    </w:fldSimple>
  </w:p>
</w:ftr>
'''

(WORK / "word/header1.xml").write_text(HEADER_XML)
(WORK / "word/footer1.xml").write_text(FOOTER_XML)

# ============================================================
# 3. [Content_Types].xml — registra header1.xml e footer1.xml
# ============================================================
ct_path = WORK / "[Content_Types].xml"
ct = ct_path.read_text()
header_ct = '<Override PartName="/word/header1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.header+xml"/>'
footer_ct = '<Override PartName="/word/footer1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.footer+xml"/>'
if header_ct not in ct:
    ct = ct.replace("</Types>", header_ct + footer_ct + "</Types>")
ct_path.write_text(ct)

# ============================================================
# 4. word/_rels/document.xml.rels — registra relationship per header/footer
# ============================================================
rels_path = WORK / "word/_rels/document.xml.rels"
rels = rels_path.read_text()
# Trova un rId libero
existing_ids = re.findall(r'Id="rId(\d+)"', rels)
max_id = max(int(x) for x in existing_ids) if existing_ids else 0
header_rid = f"rId{max_id + 1}"
footer_rid = f"rId{max_id + 2}"

header_rel = f'<Relationship Id="{header_rid}" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/header" Target="header1.xml"/>'
footer_rel = f'<Relationship Id="{footer_rid}" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/footer" Target="footer1.xml"/>'
rels = rels.replace("</Relationships>", header_rel + footer_rel + "</Relationships>")
rels_path.write_text(rels)

# ============================================================
# 5. word/document.xml — aggiungi sectPr con headerReference + footerReference
# ============================================================
doc_path = WORK / "word/document.xml"
doc = doc_path.read_text()

header_ref = f'<w:headerReference w:type="default" r:id="{header_rid}"/>'
footer_ref = f'<w:footerReference w:type="default" r:id="{footer_rid}"/>'

# Cerca sectPr esistente
if "<w:sectPr>" in doc:
    doc = doc.replace("<w:sectPr>", "<w:sectPr>" + header_ref + footer_ref, 1)
elif "<w:sectPr " in doc:
    doc = re.sub(r"<w:sectPr([^>]*)>", lambda m: f"<w:sectPr{m.group(1)}>" + header_ref + footer_ref, doc, count=1)
else:
    # Crea sectPr prima della chiusura body
    sect_pr = (
        '<w:sectPr>' + header_ref + footer_ref +
        '<w:pgSz w:w="12240" w:h="15840"/>'
        '<w:pgMar w:top="1440" w:right="1440" w:bottom="1440" w:left="1440" w:header="720" w:footer="720" w:gutter="0"/>'
        '</w:sectPr>'
    )
    doc = doc.replace("</w:body>", sect_pr + "</w:body>")

doc_path.write_text(doc)

# ============================================================
# 6. Re-zip — deterministico
# ============================================================
# Rigenerazioni successive devono produrre il MEDESIMO byte-stream così il
# git diff binario è zero quando il content semantico è identico. Otteniamo
# determinismo via:
#   - ordering alfabetico dei file (sorted)
#   - timestamp fisso 1980-01-01 (epoch zip minimum)
#   - external_attr fisso (0o644 << 16)
EPOCH = (1980, 1, 1, 0, 0, 0)
if OUT.exists():
    OUT.unlink()
with zipfile.ZipFile(OUT, "w", zipfile.ZIP_DEFLATED) as zf:
    files = sorted(p for p in WORK.rglob("*") if p.is_file())
    for p in files:
        info = zipfile.ZipInfo(str(p.relative_to(WORK)))
        info.date_time = EPOCH
        info.compress_type = zipfile.ZIP_DEFLATED
        info.external_attr = 0o644 << 16
        with p.open("rb") as fh:
            zf.writestr(info, fh.read())

print(f"OK reference.docx scritto in {OUT} ({OUT.stat().st_size} bytes)")
