// app.js — UI wiring, persistence, leaderboard and controls.
import { Graph3D } from './graph.js';
import { overallScore, tierFor, scoreColor, zoneFor, SWEET_SPOTS } from './rating.js';

const STORE_KEY = 'herbgraph:v1';

let entries = load();
let selectedId = null;

const $ = (id) => document.getElementById(id);
const form = $('entry-form');
const nameEl = $('name');
const igEl = $('ig');
const editIdEl = $('edit-id');
const formTitle = $('form-title');
const submitBtn = $('submit-btn');
const cancelBtn = $('cancel-btn');
const sliders = { x: $('ax-x'), y: $('ax-y'), z: $('ax-z') };
const outs = { x: $('out-x'), y: $('out-y'), z: $('out-z') };
const previewScore = $('preview-score');
const previewTier = $('preview-tier');
const previewZone = $('preview-zone');
const rankingList = $('ranking-list');
const emptyNote = $('empty-note');
const tooltip = $('tooltip');

const graph = new Graph3D($('graph'));
graph.onSelect = (id) => selectEntry(id);
graph.onHover = (data, event) => showTooltip(data, event);

buildLegend();
wireControls();

// ---- persistence ----
function load() {
  try {
    const raw = localStorage.getItem(STORE_KEY);
    return raw ? JSON.parse(raw) : [];
  } catch { return []; }
}
function save() { localStorage.setItem(STORE_KEY, JSON.stringify(entries)); }

// ---- helpers ----
function cleanHandle(s) {
  return (s || '').trim().replace(/^@+/, '').replace(/^https?:\/\/(www\.)?instagram\.com\//i, '').replace(/\/.*$/, '').trim();
}
function readForm() {
  return { x: parseFloat(sliders.x.value), y: parseFloat(sliders.y.value), z: parseFloat(sliders.z.value) };
}

// ---- form preview ----
function updatePreview() {
  const vals = readForm();
  outs.x.textContent = vals.x.toFixed(1);
  outs.y.textContent = vals.y.toFixed(1);
  outs.z.textContent = vals.z.toFixed(1);
  const score = overallScore(vals);
  const tier = tierFor(score);
  previewScore.textContent = score.toFixed(1);
  previewTier.textContent = tier.label;
  previewTier.style.background = scoreColor(score);
  const zone = zoneFor(vals);
  if (zone) {
    previewZone.hidden = false;
    previewZone.textContent = `${zone.badge} ${zone.label}`;
    previewZone.style.color = zone.color;
  } else {
    previewZone.hidden = true;
  }
}
for (const k of ['x', 'y', 'z']) sliders[k].addEventListener('input', updatePreview);

// ---- add / edit ----
form.addEventListener('submit', (e) => {
  e.preventDefault();
  const name = nameEl.value.trim();
  if (!name) return;
  const vals = readForm();
  const ig = cleanHandle(igEl.value);
  const id = editIdEl.value;
  if (id) {
    const entry = entries.find((x) => x.id === id);
    if (entry) Object.assign(entry, { name, ig, ...vals });
  } else {
    entries.push({ id: crypto.randomUUID(), name, ig, ...vals });
  }
  save();
  resetForm();
  render();
});
cancelBtn.addEventListener('click', resetForm);

function resetForm() {
  form.reset();
  editIdEl.value = '';
  formTitle.textContent = '// ADD SUBJECT';
  submitBtn.textContent = '+ ADD';
  cancelBtn.hidden = true;
  updatePreview();
}
function editEntry(id) {
  const entry = entries.find((x) => x.id === id);
  if (!entry) return;
  nameEl.value = entry.name;
  igEl.value = entry.ig || '';
  sliders.x.value = entry.x; sliders.y.value = entry.y; sliders.z.value = entry.z;
  editIdEl.value = id;
  formTitle.textContent = '// EDIT SUBJECT';
  submitBtn.textContent = 'SAVE';
  cancelBtn.hidden = false;
  updatePreview();
  nameEl.focus();
}
function deleteEntry(id) {
  entries = entries.filter((x) => x.id !== id);
  if (selectedId === id) selectedId = null;
  if (editIdEl.value === id) resetForm();
  save();
  render();
}
$('clear-all').addEventListener('click', () => {
  if (!entries.length) return;
  if (confirm('Clear all subjects? This cannot be undone.')) {
    entries = []; selectedId = null; save(); resetForm(); render();
  }
});

function selectEntry(id) {
  selectedId = id;
  graph.highlight(id);
  renderRanking();
  const row = rankingList.querySelector(`[data-id="${id}"]`);
  if (row) row.scrollIntoView({ block: 'nearest', behavior: 'smooth' });
}

// ---- tooltip ----
function showTooltip(data, event) {
  if (!data) { tooltip.hidden = true; return; }
  const score = overallScore(data);
  const zone = zoneFor(data);
  tooltip.hidden = false;
  tooltip.innerHTML =
    `<div class="tt-name">${escapeHtml(data.name)}</div>` +
    `<div class="tt-row">SCORE <b>${score.toFixed(1)}/10</b></div>` +
    (zone ? `<div class="tt-row" style="color:${zone.color}">${zone.badge} ${zone.label}</div>` : '') +
    (data.ig ? `<div class="tt-row">@${escapeHtml(data.ig)}</div>` : '');
  tooltip.style.left = `${event.clientX + 14}px`;
  tooltip.style.top = `${event.clientY + 14}px`;
}

// ---- legend ----
function buildLegend() {
  const el = $('legend');
  el.innerHTML = '<div class="legend-title">SWEET SPOTS</div>' +
    SWEET_SPOTS.map((s) =>
      `<div class="legend-item"><span class="legend-swatch" style="background:${s.color}"></span>${s.badge} ${s.label}</div>`
    ).join('');
}

// ---- controls ----
function wireControls() {
  $('zoom-in').addEventListener('click', () => graph.zoomBy(1.2));
  $('zoom-out').addEventListener('click', () => graph.zoomBy(0.8));
  $('reset-view').addEventListener('click', () => graph.resetView());
  $('toggle-zones').addEventListener('click', (e) => {
    const on = graph.toggleZones();
    e.currentTarget.dataset.on = on ? '1' : '0';
    $('legend').style.opacity = on ? '1' : '0.4';
  });
  $('fullscreen').addEventListener('click', () => {
    const root = document.documentElement;
    if (document.fullscreenElement) document.exitFullscreen();
    else (root.requestFullscreen || root.webkitRequestFullscreen)?.call(root);
  });
}

// ---- render ----
function render() {
  graph.setData(entries);
  if (selectedId) graph.highlight(selectedId);
  renderRanking();
  renderStats();
}

function renderStats() {
  const scores = entries.map((e) => overallScore(e));
  $('stat-count').textContent = entries.length;
  if (!scores.length) { $('stat-avg').textContent = '—'; $('stat-top').textContent = '—'; return; }
  const avg = scores.reduce((a, b) => a + b, 0) / scores.length;
  $('stat-avg').textContent = avg.toFixed(1);
  const top = entries[scores.indexOf(Math.max(...scores))];
  $('stat-top').textContent = top ? top.name : '—';
}

function renderRanking() {
  const ranked = entries
    .map((e) => ({ ...e, score: overallScore(e), tier: tierFor(overallScore(e)), zone: zoneFor(e) }))
    .sort((a, b) => b.score - a.score);

  emptyNote.hidden = ranked.length > 0;
  rankingList.innerHTML = '';

  ranked.forEach((e, i) => {
    const li = document.createElement('li');
    li.className = 'rank-row' + (e.id === selectedId ? ' selected' : '');
    li.dataset.id = e.id;
    const igLink = e.ig
      ? `<a class="ig-link" href="https://instagram.com/${encodeURIComponent(e.ig)}" target="_blank" rel="noopener" title="@${escapeHtml(e.ig)}">◎</a>`
      : '';
    const zoneTag = e.zone ? `<span class="rank-zone" style="color:${e.zone.color}">${e.zone.badge}</span>` : '';
    li.innerHTML = `
      <span class="rank-pos">${i + 1}</span>
      <span class="rank-tier" style="background:${scoreColor(e.score)}">${e.tier.label}</span>
      <span class="rank-main"><span class="rank-name">${escapeHtml(e.name)} ${zoneTag} ${igLink}</span></span>
      <span class="rank-score">${e.score.toFixed(1)}</span>
      <span class="rank-actions">
        <button class="icon-btn" data-act="edit" title="Edit">✎</button>
        <button class="icon-btn" data-act="del" title="Delete">✕</button>
      </span>`;
    li.querySelector('.rank-name').addEventListener('click', () => selectEntry(e.id));
    li.querySelector('.rank-pos').addEventListener('click', () => selectEntry(e.id));
    li.querySelector('[data-act="edit"]').addEventListener('click', () => editEntry(e.id));
    li.querySelector('[data-act="del"]').addEventListener('click', () => deleteEntry(e.id));
    rankingList.appendChild(li);
  });
}

function escapeHtml(s) {
  return String(s).replace(/[&<>"']/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]));
}

// ---- seed examples on first run ----
if (entries.length === 0 && !localStorage.getItem(STORE_KEY)) {
  entries = [
    { id: crypto.randomUUID(), name: 'Sample A', ig: '', x: 3.5, y: 4, z: 4.5 },
    { id: crypto.randomUUID(), name: 'Sample B', ig: '', x: -2, y: 3.5, z: 3 },
    { id: crypto.randomUUID(), name: 'Sample C', ig: '', x: 1, y: -1, z: -3.5 },
  ];
  save();
}

// ---- boot ----
updatePreview();
render();
