'use strict';

// ---------------------------------------------------------------------------
// Band color CSS values
// ---------------------------------------------------------------------------

const BAND_CSS = {
  black:  '#111111',
  brown:  '#6b3a2a',
  red:    '#c0392b',
  orange: '#e67e22',
  yellow: '#f1c40f',
  green:  '#27ae60',
  blue:   '#2980b9',
  violet: '#8e44ad',
  grey:   '#7f8c8d',
  white:  '#f0f0f0',
  gold:   '#d4a017',
  silver: '#bdc3c7',
  none:   'transparent',
};

// ---------------------------------------------------------------------------
// DOM helpers — no innerHTML anywhere
// ---------------------------------------------------------------------------

function el(tag, attrs) {
  const e = document.createElement(tag);
  if (attrs) {
    for (const k of Object.keys(attrs)) {
      if (k === 'textContent') e.textContent = attrs[k];
      else if (k === 'className') e.className = attrs[k];
      else e.setAttribute(k, attrs[k]);
    }
  }
  return e;
}

function append(parent, ...children) {
  for (const c of children) {
    if (c == null) continue;
    if (typeof c === 'string') parent.appendChild(document.createTextNode(c));
    else parent.appendChild(c);
  }
  return parent;
}

function clear(node) {
  while (node.firstChild) node.removeChild(node.firstChild);
}

function setResult(node, content) {
  clear(node);
  if (typeof content === 'string') {
    node.textContent = content;
  } else {
    node.appendChild(content);
  }
}

// Build a key→value table fragment
function resultTable(entries) {
  const frag = document.createDocumentFragment();
  for (const [k, v, cls] of entries) {
    const row = el('div', {className: 'r-row'});
    append(row,
      el('span', {className: 'r-key', textContent: k + ':'}),
      el('span', {className: cls ? 'r-val ' + cls : 'r-val', textContent: String(v)})
    );
    frag.appendChild(row);
  }
  return frag;
}

// Build resistor body diagram from array of color name strings
function buildDiagram(bands) {
  const wrap = el('div', {className: 'r-diagram'});
  append(wrap, el('div', {className: 'r-lead'}));
  const body = el('div', {className: 'r-body'});
  for (const c of bands) {
    const band = el('div', {className: 'r-band'});
    band.style.background = BAND_CSS[c] || '#888';
    body.appendChild(band);
  }
  append(wrap, body, el('div', {className: 'r-lead'}));
  return wrap;
}

// Build band chips (colored label chips)
function buildBandChips(bands) {
  const list = el('div', {className: 'band-list'});
  for (const c of bands) {
    const chip = el('div', {className: 'band-chip'});
    const dot = el('span', {className: 'band-dot'});
    dot.style.background = BAND_CSS[c] || '#888';
    append(chip, dot, c);
    list.appendChild(chip);
  }
  return list;
}

// ---------------------------------------------------------------------------
// Format helpers
// ---------------------------------------------------------------------------

function formatR(ohms) {
  if (ohms == null || isNaN(ohms)) return '—';
  if (ohms >= 1e6) return (ohms / 1e6).toPrecision(4) + ' MΩ';
  if (ohms >= 1e3) return (ohms / 1e3).toPrecision(4) + ' kΩ';
  return ohms.toPrecision(4) + ' Ω';
}

function formatSeries(n) {
  if (!n) return '—';
  return 'E' + n;
}

// ---------------------------------------------------------------------------
// WASM bootstrap
// ---------------------------------------------------------------------------

let wasmReady = false;

async function initWASM() {
  const statusEl = document.getElementById('wasm-status');
  try {
    const go = new Go();
    const result = await WebAssembly.instantiateStreaming(
      fetch('/static/resistor.wasm'),
      go.importObject
    );
    go.run(result.instance);
    wasmReady = true;
    statusEl.dataset.state = 'ready';
    statusEl.textContent = 'WASM ready';
    document.querySelectorAll('button.action').forEach(b => b.removeAttribute('disabled'));
  } catch (err) {
    statusEl.dataset.state = 'failed';
    statusEl.textContent = 'WASM load failed: ' + err.message;
  }
}

// ---------------------------------------------------------------------------
// Tab management
// ---------------------------------------------------------------------------

function initTabs() {
  const tabs   = Array.from(document.querySelectorAll('[role="tab"]'));
  const panels = Array.from(document.querySelectorAll('[role="tabpanel"]'));

  function activate(tab) {
    tabs.forEach(t => {
      t.setAttribute('aria-selected', 'false');
      t.setAttribute('tabindex', '-1');
    });
    panels.forEach(p => p.setAttribute('hidden', ''));

    tab.setAttribute('aria-selected', 'true');
    tab.setAttribute('tabindex', '0');
    const panel = document.getElementById(tab.getAttribute('aria-controls'));
    if (panel) panel.removeAttribute('hidden');
  }

  tabs.forEach((tab, i) => {
    tab.addEventListener('click', () => activate(tab));
    tab.addEventListener('keydown', e => {
      let next = i;
      if (e.key === 'ArrowRight') next = (i + 1) % tabs.length;
      else if (e.key === 'ArrowLeft') next = (i - 1 + tabs.length) % tabs.length;
      else if (e.key === 'Home') next = 0;
      else if (e.key === 'End') next = tabs.length - 1;
      else return;
      e.preventDefault();
      activate(tabs[next]);
      tabs[next].focus();
    });
  });
}

// ---------------------------------------------------------------------------
// Band visibility (show/hide band rows by count)
// ---------------------------------------------------------------------------

function syncBandRows(prefix, countSelectId, max) {
  const count = parseInt(document.getElementById(countSelectId).value, 10);
  for (let i = 1; i <= max; i++) {
    const row = document.getElementById(prefix + i);
    if (!row) continue;
    if (i <= count) row.removeAttribute('hidden');
    else row.setAttribute('hidden', '');
  }
}

// ---------------------------------------------------------------------------
// Tab: Select Standard Resistor
// ---------------------------------------------------------------------------

function runSelect() {
  const out = document.getElementById('sel-out');
  if (!wasmReady) { out.textContent = 'WASM not ready'; return; }

  const r = parseFloat(document.getElementById('sel-r').value);
  if (!r || r <= 0) { out.textContent = 'Enter a positive resistance value'; return; }

  const series   = document.getElementById('sel-series').value;
  const rounding = document.getElementById('sel-round').value;

  const res = resistor.selectStandardResistor(JSON.stringify({
    resistance: r, series, rounding
  }));

  if (!res.ok) { out.textContent = 'Error: ' + res.error; return; }

  const v = res.value;
  const entries = [
    ['Requested',    formatR(v.RequestedResistance)],
    ['Selected',     formatR(v.SelectedResistance)],
    ['Series',       formatSeries(v.Series)],
    ['Tolerance',    v.TolerancePct + ' %'],
    ['Rounding',     v.Rounding],
  ];
  if (v.Bands && v.Bands.length) {
    entries.push(['Bands', v.Bands.join(', ')]);
  }
  if (v.Assumptions && v.Assumptions.length) {
    entries.push(['Assumptions', v.Assumptions.join('; ')]);
  }

  const frag = document.createDocumentFragment();
  frag.appendChild(resultTable(entries));
  if (v.Bands && v.Bands.length) {
    frag.appendChild(buildDiagram(v.Bands));
    frag.appendChild(buildBandChips(v.Bands));
  }
  setResult(out, frag);
}

// ---------------------------------------------------------------------------
// Tab: Decode Bands
// ---------------------------------------------------------------------------

function runDecodeBands() {
  const out = document.getElementById('decode-bands-out');
  if (!wasmReady) { out.textContent = 'WASM not ready'; return; }

  const count = parseInt(document.getElementById('db-count').value, 10);
  const colors = [];
  for (let i = 1; i <= count; i++) {
    const sel = document.getElementById('db-band' + i);
    if (sel) colors.push(sel.value);
  }

  const res = resistor.decodeBands(JSON.stringify(colors));
  if (!res.ok) { out.textContent = 'Error: ' + res.error; return; }

  const v = res.value;
  const entries = [
    ['Resistance',  formatR(v.ResistanceOhms)],
    ['Tolerance',   v.TolerancePct ? v.TolerancePct + ' %' : '±20 % (no band)'],
    ['Power',       v.PowerWatts   ? v.PowerWatts + ' W'  : '—'],
    ['Type',        v.Type    || '—'],
    ['Package',     v.Package || '—'],
  ];

  const frag = document.createDocumentFragment();
  frag.appendChild(resultTable(entries));
  frag.appendChild(buildDiagram(colors));
  frag.appendChild(buildBandChips(colors));
  setResult(out, frag);
}

// ---------------------------------------------------------------------------
// Tab: Encode Bands
// ---------------------------------------------------------------------------

function runEncodeBands() {
  const out = document.getElementById('encode-bands-out');
  if (!wasmReady) { out.textContent = 'WASM not ready'; return; }

  const r   = parseFloat(document.getElementById('eb-r').value);
  const tol = parseFloat(document.getElementById('eb-tol').value) || 5;
  if (!r || r <= 0) { out.textContent = 'Enter a positive resistance value'; return; }

  const res = resistor.encodeBands(JSON.stringify({resistanceOhms: r, tolerancePct: tol}));
  if (!res.ok) { out.textContent = 'Error: ' + res.error; return; }

  const bands = res.value;
  const frag = document.createDocumentFragment();
  frag.appendChild(resultTable([['Bands', bands.join(', ')]]));
  frag.appendChild(buildDiagram(bands));
  frag.appendChild(buildBandChips(bands));
  setResult(out, frag);
}

// ---------------------------------------------------------------------------
// Tab: SMD decode
// ---------------------------------------------------------------------------

function runDecodeSMD() {
  const out = document.getElementById('smd-decode-out');
  if (!wasmReady) { out.textContent = 'WASM not ready'; return; }

  const marking = document.getElementById('smd-marking').value.trim();
  if (!marking) { out.textContent = 'Enter an SMD marking'; return; }

  const res = resistor.decodeSMD(marking);
  if (!res.ok) { out.textContent = 'Error: ' + res.error; return; }

  const v = res.value;
  setResult(out, resultTable([
    ['Resistance',  formatR(v.ResistanceOhms)],
    ['Tolerance',   v.TolerancePct ? v.TolerancePct + ' %' : '—'],
  ]));
}

// ---------------------------------------------------------------------------
// Tab: SMD encode
// ---------------------------------------------------------------------------

function runEncodeSMD() {
  const out = document.getElementById('smd-encode-out');
  if (!wasmReady) { out.textContent = 'WASM not ready'; return; }

  const r    = parseFloat(document.getElementById('smd-enc-r').value);
  const mode = document.getElementById('smd-mode').value;
  if (!r || r <= 0) { out.textContent = 'Enter a positive resistance value'; return; }

  const res = resistor.encodeSMD(JSON.stringify({resistance: r, mode}));
  if (!res.ok) { out.textContent = 'Error: ' + res.error; return; }

  setResult(out, resultTable([['Marking', res.value]]));
}

// ---------------------------------------------------------------------------
// Tab: Infer
// ---------------------------------------------------------------------------

function runInfer() {
  const out = document.getElementById('infer-out');
  if (!wasmReady) { out.textContent = 'WASM not ready'; return; }

  const bandCount = parseInt(document.getElementById('infer-count').value, 10);
  const bands = [];
  for (let i = 1; i <= bandCount; i++) {
    const sel = document.getElementById('infer-band' + i);
    if (sel && sel.value !== 'none') bands.push(sel.value);
  }

  const bodyColor = document.getElementById('infer-body').value;
  const length    = parseFloat(document.getElementById('infer-len').value) || 0;
  const pkg       = document.getElementById('infer-pkg').value;
  const marking   = document.getElementById('infer-marking').value.trim();

  const input = {};
  if (bands.length)          input.bands     = bands;
  if (bodyColor !== 'none')  input.bodyColor  = bodyColor;
  if (length > 0)            input.lengthMM   = length;
  if (pkg)                   input.package    = pkg;
  if (marking)               input.marking    = marking;

  const res = resistor.inferResistor(JSON.stringify(input));
  if (!res.ok) { out.textContent = 'Error: ' + res.error; return; }

  const v    = res.value;
  const spec = v.Spec  || {};
  const meta = v.Meta  || {};

  const entries = [
    ['Resistance',  formatR(spec.ResistanceOhms)],
    ['Confidence',  ((meta.Confidence || 0) * 100).toFixed(1) + ' %'],
  ];
  if (v.VoltageRating) entries.push(['Voltage Rating', v.VoltageRating + ' V']);
  if (spec.PowerWatts) entries.push(['Power',          spec.PowerWatts + ' W']);
  if (spec.Type)       entries.push(['Type',           spec.Type]);
  if (spec.Package)    entries.push(['Package',        spec.Package]);
  if (meta.Assumptions && meta.Assumptions.length) {
    for (const a of meta.Assumptions) {
      entries.push(['Assumption', a]);
    }
  }

  setResult(out, resultTable(entries));
}

// ---------------------------------------------------------------------------
// Tab: Analyze
// ---------------------------------------------------------------------------

function runAnalyze() {
  const out = document.getElementById('analyze-out');
  if (!wasmReady) { out.textContent = 'WASM not ready'; return; }

  const r   = parseFloat(document.getElementById('an-r').value);
  const pwr = parseFloat(document.getElementById('an-pwr').value) || 0;
  const tol = parseFloat(document.getElementById('an-tol').value) || 0;
  const v   = parseFloat(document.getElementById('an-v').value)   || 0;
  const i   = parseFloat(document.getElementById('an-i').value)   || 0;

  if (!r || r <= 0)  { out.textContent = 'Enter a positive resistance value'; return; }
  if (!v && !i)      { out.textContent = 'Enter applied voltage or current';  return; }

  const res = resistor.analyzeResistor(JSON.stringify({
    spec: {resistanceOhms: r, powerWatts: pwr, tolerancePct: tol},
    appliedVoltage: v,
    appliedCurrent: i,
  }));

  if (!res.ok) { out.textContent = 'Error: ' + res.error; return; }

  const d = res.value;
  const entries = [
    ['Power Dissipation', d.PowerDissipation.toFixed(6) + ' W'],
    ['Voltage Drop',      d.VoltageDrop.toFixed(6)      + ' V'],
    ['Current',           (d.Current * 1000).toFixed(6) + ' mA'],
  ];
  if (d.DeratedSafePower != null) {
    entries.push(['Derated Safe Power', d.DeratedSafePower.toFixed(6) + ' W']);
  }
  if (d.WorstCaseResistanceMin != null) {
    entries.push(['Worst Case Min', formatR(d.WorstCaseResistanceMin)]);
    entries.push(['Worst Case Max', formatR(d.WorstCaseResistanceMax)]);
  }

  const frag = document.createDocumentFragment();
  frag.appendChild(resultTable(entries));

  if (d.Warnings && d.Warnings.length) {
    for (const w of d.Warnings) {
      const row = el('div', {className: 'r-row'});
      append(row,
        el('span', {className: 'r-key lvl-' + w.Level, textContent: '[' + w.Level + ']'}),
        el('span', {className: 'r-val', textContent: w.Message})
      );
      frag.appendChild(row);
    }
  }

  setResult(out, frag);
}

// ---------------------------------------------------------------------------
// Init
// ---------------------------------------------------------------------------

document.addEventListener('DOMContentLoaded', () => {
  initTabs();
  initWASM();

  const dbCount = document.getElementById('db-count');
  if (dbCount) {
    dbCount.addEventListener('change', () => syncBandRows('db-row', 'db-count', 6));
    syncBandRows('db-row', 'db-count', 6);
  }

  const inferCount = document.getElementById('infer-count');
  if (inferCount) {
    inferCount.addEventListener('change', () => syncBandRows('infer-row', 'infer-count', 6));
    syncBandRows('infer-row', 'infer-count', 6);
  }
});
