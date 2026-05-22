// LabNexus PDF template — Sprint 1.5.D
// Stack: pandoc (MD → Typst markup) → typst compile → PDF
// Brand AICertus: palette estratta da labnexus.app/ecosistema/aicertus.

#let labnexus-primary = rgb("#010C23")  // navy dark
#let labnexus-accent  = rgb("#006A9C")  // blu corporate
#let labnexus-light   = rgb("#0F8DBD")  // azzurro chiaro
#let labnexus-bg-pale = rgb("#F0F5FE")  // sfondo box
#let labnexus-text    = rgb("#333333")  // body
#let labnexus-muted   = rgb("#AAAAAA")  // footer

// META è popolato dal runner Go via interpolazione di variabili.
// Default vuoti per fallback se non passato.
#let meta-capability = sys.inputs.at("capability", default: "")
#let meta-profilo    = sys.inputs.at("profilo", default: "")
#let meta-modello    = sys.inputs.at("modello", default: "")
#let meta-provider   = sys.inputs.at("provider", default: "")
#let meta-data       = sys.inputs.at("data_esecuzione", default: "")
#let meta-durata     = sys.inputs.at("durata_secondi", default: "")
#let meta-token      = sys.inputs.at("token_stimati", default: "")
#let meta-fileinput  = sys.inputs.at("file_input", default: "")
#let meta-stato      = sys.inputs.at("stato", default: "completato")
#let meta-title      = if meta-capability != "" { meta-capability } else { meta-profilo }

// Page setup
#set page(
  paper: "a4",
  margin: (x: 25mm, top: 28mm, bottom: 22mm),
  header: [
    #set text(size: 9pt)
    #grid(
      columns: (1fr, auto),
      align: (left, right),
      [
        #text(fill: labnexus-primary, weight: "bold", size: 10pt)[LabNexus]
        #text(fill: labnexus-accent)[ | ]
        #text(fill: labnexus-text)[#meta-title]
      ],
      text(fill: labnexus-muted)[#meta-data],
    )
    #v(2mm)
    #line(length: 100%, stroke: 0.3pt + labnexus-accent)
  ],
  footer: [
    #line(length: 100%, stroke: 0.3pt + labnexus-accent)
    #v(1mm)
    #set text(size: 8pt, fill: labnexus-muted)
    #grid(
      columns: (1fr, auto),
      align: (left, right),
      [profilo: #meta-profilo  ·  #meta-modello],
      [pag #context counter(page).display() / #context counter(page).final().first()],
    )
  ],
)

// Typography
#set text(font: ("Helvetica Neue", "Helvetica", "Arial"), size: 10.5pt, fill: labnexus-text, lang: "it")
#set par(justify: true, leading: 0.65em, first-line-indent: 0pt)

#show heading.where(level: 1): h => block(below: 1em)[
  #set text(fill: labnexus-primary, weight: "bold", size: 20pt)
  #h.body
  #v(-0.3em)
  #line(length: 100%, stroke: 0.5pt + labnexus-accent)
]
#show heading.where(level: 2): h => block(above: 1.2em, below: 0.6em)[
  #set text(fill: labnexus-primary, weight: "bold", size: 14pt)
  #h.body
  #v(-0.3em)
  #line(length: 100%, stroke: 0.2pt + labnexus-accent)
]
#show heading.where(level: 3): h => block(above: 1em, below: 0.4em)[
  #set text(fill: labnexus-primary, weight: "bold", size: 12pt)
  #h.body
]
#show heading.where(level: 4): h => block(above: 0.8em, below: 0.3em)[
  #set text(fill: labnexus-primary, weight: "bold", size: 11pt)
  #h.body
]

// Tables: header bg pale + righe alternate
#show table: t => {
  set table(
    stroke: 0.3pt + labnexus-muted,
    fill: (col, row) => if row == 0 { labnexus-bg-pale } else if calc.odd(row) { rgb("#FAFCFF") } else { white },
    inset: 6pt,
  )
  t
}

// Code block / code span
#show raw.where(block: true): it => block(
  width: 100%,
  fill: labnexus-bg-pale,
  stroke: (left: 1.5pt + labnexus-accent),
  inset: 8pt,
  radius: 0pt,
  text(font: ("Courier", "Menlo", "Monaco"), size: 9pt, fill: labnexus-text)[#it.text],
)
#show raw.where(block: false): it => box(
  fill: labnexus-bg-pale,
  inset: (x: 3pt, y: 1pt),
  outset: (y: 1pt),
  text(font: ("Courier", "Menlo", "Monaco"), size: 9pt)[#it.text],
)

// Links: accent color, no underline
#show link: it => text(fill: labnexus-accent)[#it]

// Blockquote standard (non-callout)
#show quote.where(block: true): q => block(
  width: 100%,
  inset: (left: 8pt, right: 4pt, y: 4pt),
  stroke: (left: 1.5pt + labnexus-accent),
  text(style: "italic")[#q.body],
)

// Pandoc emette `#horizontalrule` per `---` (thematic break MD): typst non
// definisce questa funzione di default, quindi la mappiamo a una linea sottile.
#let horizontalrule = block(above: 0.8em, below: 0.8em)[
  #line(length: 100%, stroke: 0.3pt + labnexus-muted)
]

// Callout LabNexus (chiamato dai raw typst blocks generati dal pre-processor Go)
#let callout(kind: "", title: none, ..body) = block(
  width: 100%,
  fill: labnexus-bg-pale,
  stroke: (left: 2pt + labnexus-accent),
  inset: 10pt,
  above: 0.8em,
  below: 0.8em,
  {
    text(fill: labnexus-accent, weight: "bold", size: 9pt)[\[#upper(kind)\]]
    if title != none and title != "" [
      #linebreak()
      #text(fill: labnexus-primary, weight: "bold", size: 10.5pt)[#title]
    ]
    if body.pos().len() > 0 [
      #linebreak()
      #text(fill: labnexus-text)[#body.pos().join()]
    ]
  }
)

// Box metadati "Dati di esecuzione" — disegnato in cima alla 1° pagina
#let metadata-box() = {
  if meta-profilo == "" { return }
  block(
    width: 100%,
    fill: labnexus-bg-pale,
    stroke: (left: 2pt + labnexus-accent),
    inset: 10pt,
    above: 0pt,
    below: 1.2em,
    {
      text(fill: labnexus-primary, weight: "bold", size: 11pt)[Dati di esecuzione]
      v(0.4em)
      let rows = (
        if meta-capability != "" { ("Capability", meta-capability) } else { none },
        ("Profilo", meta-profilo),
        ("Modello", meta-modello),
        ("Provider", meta-provider),
        ("Data", meta-data),
        if meta-durata != "" { ("Durata", meta-durata + " s") } else { none },
        if meta-token != "" { ("Token stimati", meta-token) } else { none },
        if meta-fileinput != "" { ("File input", meta-fileinput) } else { none },
        ("Stato", meta-stato),
      ).filter(r => r != none)
      table(
        columns: (auto, 1fr),
        stroke: none,
        inset: (x: 0pt, y: 2pt),
        fill: none,
        ..rows.map(r => (
          text(fill: labnexus-primary, weight: "bold", size: 9pt)[#r.at(0)],
          text(fill: labnexus-text, size: 9pt)[#r.at(1)],
        )).flatten()
      )
    }
  )
}

#metadata-box()
