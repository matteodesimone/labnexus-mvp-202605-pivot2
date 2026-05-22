-- callout-filter.lua — pandoc Lua filter (Sprint 1.5.D)
--
-- Intercetta i blockquote Obsidian `> [!TYPE] title\n> body...` e li trasforma
-- in raw blocks Typst che invocano la funzione `#callout(kind, title)[body]`
-- definita nel template LabNexus.
--
-- Esempio markdown input:
--   > [!ATTENZIONE] Scadenza ACCREDIA
--   > La RT-23 Rev.05 entra in vigore il 2026-09-30.
--   > Le modifiche devono essere recepite prima.
--
-- Output typst raw (renderizzato come box stilizzato dal template):
--   #callout(kind: "ATTENZIONE", title: [Scadenza ACCREDIA])[
--   La RT-23 Rev.05 entra in vigore il 2026-09-30. Le modifiche devono essere recepite prima.
--   ]

function BlockQuote(blockquote)
  if #blockquote.content == 0 then return nil end
  local first_block = blockquote.content[1]
  if first_block.t ~= "Para" then return nil end

  -- La prima inline deve essere uno Str che matcha "[!TYPE]"
  local first_inline = first_block.content[1]
  if not first_inline or first_inline.t ~= "Str" then return nil end
  local kind = first_inline.text:match("^%[%!([%w_-]+)%]$")
  if not kind then return nil end

  -- Estrai title (inline dopo "[!TYPE]" fino al primo Break/SoftBreak)
  -- e resto-della-prima-para (dopo il break)
  local title_inlines = {}
  local rest_first_para = {}
  local saw_break = false
  local skipping_leading_space = true
  for i = 2, #first_block.content do
    local item = first_block.content[i]
    if not saw_break then
      if item.t == "SoftBreak" or item.t == "LineBreak" then
        saw_break = true
        skipping_leading_space = true
      elseif skipping_leading_space and item.t == "Space" then
        -- skip
      else
        skipping_leading_space = false
        table.insert(title_inlines, item)
      end
    else
      if skipping_leading_space and item.t == "Space" then
        -- skip
      else
        skipping_leading_space = false
        table.insert(rest_first_para, item)
      end
    end
  end

  -- Body blocks: resto della prima para (se c'è) + paragrafi successivi del blockquote
  local body_blocks = {}
  if #rest_first_para > 0 then
    table.insert(body_blocks, pandoc.Para(rest_first_para))
  end
  for i = 2, #blockquote.content do
    table.insert(body_blocks, blockquote.content[i])
  end

  -- Output dispatch: typst usa #callout(...) come raw block (template definisce
  -- la funzione); docx riscrive come BlockQuote con header "[KIND] title" bold
  -- in run separati (no styling custom — perde il box visuale del PDF, ma
  -- preserva semantica e leggibilità con prefix testuale).
  if FORMAT == "typst" then
    -- Converti title/body in typst markup
    local title_typst = ""
    if #title_inlines > 0 then
      local title_doc = pandoc.Pandoc({pandoc.Plain(title_inlines)}, pandoc.Meta({}))
      title_typst = pandoc.write(title_doc, "typst")
      title_typst = title_typst:gsub("\n+$", ""):gsub("\n", " ")
    end
    local body_typst = ""
    if #body_blocks > 0 then
      local body_doc = pandoc.Pandoc(body_blocks, pandoc.Meta({}))
      body_typst = pandoc.write(body_doc, "typst")
      body_typst = body_typst:gsub("\n+$", "")
    end
    local raw
    if title_typst ~= "" and body_typst ~= "" then
      raw = string.format('#callout(kind: "%s", title: [%s])[\n%s\n]', kind, title_typst, body_typst)
    elseif title_typst ~= "" then
      raw = string.format('#callout(kind: "%s", title: [%s])[]', kind, title_typst)
    else
      raw = string.format('#callout(kind: "%s")[\n%s\n]', kind, body_typst)
    end
    return pandoc.RawBlock("typst", raw)
  end

  -- Default branch (es. docx): ricostruisci BlockQuote con header bold
  -- "[KIND] title" + body blocks. Pandoc DOCX writer rende come Quote style
  -- (bordo sinistro indented), e il prefix testuale "[KIND]" preserva la
  -- "tipizzazione" del callout anche senza box visuale custom.
  local header_inlines = {pandoc.Strong({pandoc.Str("[" .. kind .. "]")})}
  if #title_inlines > 0 then
    table.insert(header_inlines, pandoc.Space())
    table.insert(header_inlines, pandoc.Strong(title_inlines))
  end
  local new_blocks = {pandoc.Para(header_inlines)}
  for _, b in ipairs(body_blocks) do
    table.insert(new_blocks, b)
  end
  return pandoc.BlockQuote(new_blocks)
end
