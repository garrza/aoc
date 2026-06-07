// graph.js — dependency-free 3D scatter plot rendered on a 2D canvas.
// Implements its own rotation + perspective projection, orbit-by-drag,
// wheel/pinch zoom, sweet-spot zones, and click/hover picking. No libraries.
import { scoreColor, overallScore, AXES, SWEET_SPOTS } from './rating.js';

const RANGE = 5; // world coords span -RANGE..RANGE on each axis
const CAM_DIST = 20; // camera distance for perspective
const AXIS_COLORS = { x: '#ff6b9d', y: '#6bc1ff', z: '#9dff6b' };
const DEFAULT_VIEW = { yaw: 0.7, pitch: 0.5, zoom: 1 };

export class Graph3D {
  constructor(container) {
    this.container = container;
    this.canvas = document.createElement('canvas');
    this.ctx = this.canvas.getContext('2d');
    this.container.appendChild(this.canvas);

    this.data = [];
    this.projected = [];
    this.onSelect = null;
    this.onHover = null;
    this.highlightId = null;
    this.showZones = true;

    this.yaw = DEFAULT_VIEW.yaw;
    this.pitch = DEFAULT_VIEW.pitch;
    this.zoom = DEFAULT_VIEW.zoom;

    this._initInput();
    this._resize();
    window.addEventListener('resize', () => this._resize());
    this._loop();
  }

  setData(data) {
    this.data = data.map((d) => ({ ...d, score: overallScore(d) }));
  }
  highlight(id) { this.highlightId = id; }
  toggleZones() { this.showZones = !this.showZones; return this.showZones; }
  zoomBy(f) { this.zoom = clampN(this.zoom * f, 0.4, 3); }
  resetView() { Object.assign(this, DEFAULT_VIEW); }

  // ---- view math ----
  _rotate(p) {
    const cy = Math.cos(this.yaw), sy = Math.sin(this.yaw);
    const cx = Math.cos(this.pitch), sx = Math.sin(this.pitch);
    const x1 = p.x * cy + p.z * sy;
    const z1 = -p.x * sy + p.z * cy;
    const y1 = p.y;
    const y2 = y1 * cx - z1 * sx;
    const z2 = y1 * sx + z1 * cx;
    return { x: x1, y: y2, z: z2 };
  }
  _project(p) {
    const r = this._rotate(p);
    const persp = CAM_DIST / (CAM_DIST - r.z);
    return { sx: this.cx + r.x * persp * this.unit, sy: this.cy - r.y * persp * this.unit, depth: r.z, persp };
  }

  // ---- input (mouse + touch) ----
  _initInput() {
    const el = this.canvas;
    const pointers = new Map();
    let lastX = 0, lastY = 0, moved = false, pinchDist = 0;

    el.addEventListener('pointerdown', (e) => {
      pointers.set(e.pointerId, { x: e.clientX, y: e.clientY });
      el.setPointerCapture(e.pointerId);
      moved = false; lastX = e.clientX; lastY = e.clientY;
      el.style.cursor = 'grabbing';
      if (pointers.size === 2) pinchDist = this._pinchDist(pointers);
    });
    el.addEventListener('pointermove', (e) => {
      if (pointers.has(e.pointerId)) pointers.set(e.pointerId, { x: e.clientX, y: e.clientY });

      if (pointers.size >= 2) {
        const d = this._pinchDist(pointers);
        if (pinchDist) this.zoomBy(d / pinchDist);
        pinchDist = d; moved = true;
        return;
      }
      if (pointers.size === 1) {
        const dx = e.clientX - lastX, dy = e.clientY - lastY;
        if (Math.abs(dx) + Math.abs(dy) > 2) moved = true;
        this.yaw += dx * 0.01;
        this.pitch = clampN(this.pitch + dy * 0.01, -1.45, 1.45);
        lastX = e.clientX; lastY = e.clientY;
      } else {
        this._hoverAt(e);
      }
    });
    const end = (e) => {
      const wasOne = pointers.size === 1;
      pointers.delete(e.pointerId);
      pinchDist = 0;
      el.style.cursor = 'grab';
      if (wasOne && !moved) this._clickAt(e);
    };
    el.addEventListener('pointerup', end);
    el.addEventListener('pointercancel', (e) => pointers.delete(e.pointerId));
    el.addEventListener('pointerleave', () => { if (this.onHover) this.onHover(null); });
    el.addEventListener('wheel', (e) => {
      e.preventDefault();
      this.zoomBy(e.deltaY < 0 ? 1.1 : 0.9);
    }, { passive: false });
  }
  _pinchDist(pointers) {
    const [a, b] = [...pointers.values()];
    return Math.hypot(a.x - b.x, a.y - b.y);
  }

  _pickAt(e) {
    const rect = this.canvas.getBoundingClientRect();
    const mx = e.clientX - rect.left, my = e.clientY - rect.top;
    let best = null, bestD = 18;
    for (const pr of this.projected) {
      const d = Math.hypot(pr.sx - mx, pr.sy - my);
      if (d < pr.r + 8 && (!best || pr.depth > best.depth)) { best = pr; bestD = d; }
    }
    return best;
  }
  _clickAt(e) { const h = this._pickAt(e); if (h && this.onSelect) this.onSelect(h.id); }
  _hoverAt(e) {
    const h = this._pickAt(e);
    this.canvas.style.cursor = h ? 'pointer' : 'grab';
    if (this.onHover) this.onHover(h ? { ...h.data } : null, e);
  }

  _resize() {
    const w = this.container.clientWidth, h = this.container.clientHeight;
    const dpr = window.devicePixelRatio || 1;
    this.w = w; this.h = h;
    this.canvas.width = w * dpr; this.canvas.height = h * dpr;
    this.canvas.style.width = w + 'px'; this.canvas.style.height = h + 'px';
    this.ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
    this.cx = w / 2; this.cy = h / 2;
    this.canvas.style.cursor = 'grab';
  }

  _loop() {
    requestAnimationFrame(() => this._loop());
    this.unit = (Math.min(this.w, this.h) / 2 / RANGE) * 0.56 * this.zoom;
    this._draw();
  }

  _draw() {
    const ctx = this.ctx;
    ctx.clearRect(0, 0, this.w, this.h);
    ctx.fillStyle = '#0a0c11';
    ctx.fillRect(0, 0, this.w, this.h);

    if (this.showZones) this._drawZones();
    this._drawCube();
    this._drawAxes();
    this._drawPoints();
  }

  _cubeCorners() {
    const c = [];
    for (const x of [-RANGE, RANGE]) for (const y of [-RANGE, RANGE]) for (const z of [-RANGE, RANGE]) c.push({ x, y, z });
    return c;
  }

  _drawCube() {
    const ctx = this.ctx;
    const corners = this._cubeCorners().map((p) => this._project(p));
    const edges = [[0,1],[0,2],[0,4],[1,3],[1,5],[2,3],[2,6],[3,7],[4,5],[4,6],[5,7],[6,7]];
    ctx.strokeStyle = '#1f2632'; ctx.lineWidth = 1;
    for (const [a, b] of edges) {
      ctx.beginPath(); ctx.moveTo(corners[a].sx, corners[a].sy); ctx.lineTo(corners[b].sx, corners[b].sy); ctx.stroke();
    }
  }

  _drawAxes() {
    const ctx = this.ctx;
    const defs = [
      { key: 'x', a: { x: -RANGE, y: 0, z: 0 }, b: { x: RANGE, y: 0, z: 0 } },
      { key: 'y', a: { x: 0, y: -RANGE, z: 0 }, b: { x: 0, y: RANGE, z: 0 } },
      { key: 'z', a: { x: 0, y: 0, z: -RANGE }, b: { x: 0, y: 0, z: RANGE } },
    ];
    for (const ax of defs) {
      const pa = this._project(ax.a), pb = this._project(ax.b);
      ctx.strokeStyle = AXIS_COLORS[ax.key]; ctx.globalAlpha = 0.45; ctx.lineWidth = 1.5;
      ctx.beginPath(); ctx.moveTo(pa.sx, pa.sy); ctx.lineTo(pb.sx, pb.sy); ctx.stroke();
      ctx.globalAlpha = 1;
      this._poleLabel(`${AXES[ax.key].pos} +`, ax.b, AXIS_COLORS[ax.key]);
      this._poleLabel(`− ${AXES[ax.key].neg}`, ax.a, AXIS_COLORS[ax.key]);
    }
  }

  _poleLabel(text, worldPos, color) {
    const ctx = this.ctx;
    const out = { x: worldPos.x * 1.16, y: worldPos.y * 1.16, z: worldPos.z * 1.16 };
    const p = this._project(out);
    ctx.font = 'bold 12px ui-monospace, Menlo, monospace';
    ctx.fillStyle = color; ctx.textAlign = 'center'; ctx.textBaseline = 'middle';
    ctx.fillText(text, p.sx, p.sy);
    ctx.textAlign = 'left';
  }

  _drawZones() {
    const ctx = this.ctx;
    // gather zone faces with depth so they sort correctly against each other
    for (const s of SWEET_SPOTS) {
      const b = s.bounds;
      const corners = [];
      for (const x of [b.x[0], b.x[1]]) for (const y of [b.y[0], b.y[1]]) for (const z of [b.z[0], b.z[1]])
        corners.push(this._project({ x, y, z }));
      // 6 faces by corner index (order matches x,y,z nesting above)
      const faces = [
        [0,1,3,2],[4,5,7,6], // x- , x+
        [0,1,5,4],[2,3,7,6], // y- , y+
        [0,2,6,4],[1,3,7,5], // z- , z+
      ];
      const polys = faces.map((f) => ({
        pts: f.map((i) => corners[i]),
        depth: f.reduce((s2, i) => s2 + corners[i].depth, 0) / 4,
      })).sort((a, c) => a.depth - c.depth);

      for (const poly of polys) {
        ctx.beginPath();
        poly.pts.forEach((p, i) => (i ? ctx.lineTo(p.sx, p.sy) : ctx.moveTo(p.sx, p.sy)));
        ctx.closePath();
        ctx.fillStyle = hexA(s.color, 0.06); ctx.fill();
        ctx.strokeStyle = hexA(s.color, 0.5); ctx.lineWidth = 1; ctx.stroke();
      }
      // label at zone center
      const cx = (b.x[0] + b.x[1]) / 2, cy = (b.y[0] + b.y[1]) / 2, cz = (b.z[0] + b.z[1]) / 2;
      const c = this._project({ x: cx, y: cy, z: cz });
      ctx.font = 'bold 11px ui-monospace, Menlo, monospace';
      ctx.fillStyle = s.color; ctx.textAlign = 'center'; ctx.textBaseline = 'middle';
      ctx.fillText(`${s.badge} ${s.label}`, c.sx, c.sy);
      ctx.textAlign = 'left';
    }
  }

  _drawPoints() {
    const ctx = this.ctx;
    this.projected = [];
    const pts = this.data.map((d) => ({ d, pr: this._project(d) })).sort((a, b) => a.pr.depth - b.pr.depth);

    for (const { d, pr } of pts) {
      const isHi = d.id === this.highlightId;
      const r = Math.max(4, 7 * pr.persp) * (isHi ? 1.7 : 1);
      const color = scoreColor(d.score);

      const foot = this._project({ x: d.x, y: -RANGE, z: d.z });
      ctx.strokeStyle = color; ctx.globalAlpha = 0.2; ctx.lineWidth = 1;
      ctx.beginPath(); ctx.moveTo(pr.sx, pr.sy); ctx.lineTo(foot.sx, foot.sy); ctx.stroke();
      ctx.globalAlpha = 1;

      const grad = ctx.createRadialGradient(pr.sx - r * 0.3, pr.sy - r * 0.3, r * 0.1, pr.sx, pr.sy, r);
      grad.addColorStop(0, '#ffffff'); grad.addColorStop(0.25, color); grad.addColorStop(1, shade(color, -0.4));
      ctx.fillStyle = grad;
      ctx.beginPath(); ctx.arc(pr.sx, pr.sy, r, 0, Math.PI * 2); ctx.fill();
      if (isHi) { ctx.strokeStyle = '#fff'; ctx.lineWidth = 2; ctx.stroke(); }

      ctx.font = `${isHi ? 'bold ' : ''}13px Inter, system-ui, sans-serif`;
      const tw = ctx.measureText(d.name).width;
      ctx.fillStyle = 'rgba(10,12,17,0.7)';
      roundRect(ctx, pr.sx - tw / 2 - 5, pr.sy - r - 23, tw + 10, 18, 4); ctx.fill();
      ctx.fillStyle = '#e6edf3'; ctx.textAlign = 'center'; ctx.textBaseline = 'middle';
      ctx.fillText(d.name, pr.sx, pr.sy - r - 14);
      ctx.textAlign = 'left';

      this.projected.push({ id: d.id, sx: pr.sx, sy: pr.sy, r, depth: pr.depth, data: d });
    }
  }
}

// ---- helpers ----
function clampN(v, lo, hi) { return Math.max(lo, Math.min(hi, v)); }

function roundRect(ctx, x, y, w, h, r) {
  ctx.beginPath();
  ctx.moveTo(x + r, y);
  ctx.arcTo(x + w, y, x + w, y + h, r);
  ctx.arcTo(x + w, y + h, x, y + h, r);
  ctx.arcTo(x, y + h, x, y, r);
  ctx.arcTo(x, y, x + w, y, r);
  ctx.closePath();
}
function shade(hex, amt) {
  const n = parseInt(hex.slice(1), 16);
  let r = (n >> 16) & 255, g = (n >> 8) & 255, b = n & 255;
  const f = amt < 0 ? 1 + amt : 1, t = amt < 0 ? 0 : 255 * amt;
  r = Math.round(r * f + t); g = Math.round(g * f + t); b = Math.round(b * f + t);
  return `rgb(${r},${g},${b})`;
}
function hexA(hex, a) {
  const n = parseInt(hex.slice(1), 16);
  return `rgba(${(n >> 16) & 255},${(n >> 8) & 255},${n & 255},${a})`;
}
