// rating.js — transparent scoring model for the 3D ranking.
//
// Each axis is a spectrum from -5 (negative pole) to +5 (positive pole):
//   x: Naca  (-5)  <->  Fresa (+5)
//   y: Basic (-5)  <->  Artsy (+5)
//   z: Culera(-5)  <->  Chida (+5)
//
// The "ideal corner" is (+5, +5, +5). Overall score blends the three axes
// with weights that favor personality (Chida) the most, since being chida
// matters more than aesthetics.

export const AXES = {
  x: { neg: 'Naca', pos: 'Fresa', weight: 0.30 },
  y: { neg: 'Basic', pos: 'Artsy', weight: 0.25 },
  z: { neg: 'Culera', pos: 'Chida', weight: 0.45 },
};

// Map an axis value [-5, 5] to a 0..1 share.
const norm = (v) => (clamp(v, -5, 5) + 5) / 10;

export function clamp(v, lo, hi) {
  return Math.max(lo, Math.min(hi, v));
}

// Overall rating on a 0..10 scale (one decimal of resolution).
export function overallScore({ x, y, z }) {
  const blend =
    norm(x) * AXES.x.weight +
    norm(y) * AXES.y.weight +
    norm(z) * AXES.z.weight;
  return Math.round(blend * 1000) / 100; // 0..10, 2-decimal precision
}

// Sweet spots — named cuboid regions of the space that mean something.
// Each has axis-aligned bounds [min, max] per axis, a color and a label.
export const SWEET_SPOTS = [
  {
    key: 'elite',
    label: 'ELITE',
    desc: 'Fresa · Artsy · Chida — the dream corner',
    color: '#2dd4bf',
    badge: '◆',
    bounds: { x: [2.5, 5], y: [2.5, 5], z: [2.5, 5] },
  },
  {
    key: 'altbae',
    label: 'ALT BAE',
    desc: 'Artsy & Chida, naca-leaning — the cool alt type',
    color: '#a78bfa',
    badge: '✦',
    bounds: { x: [-5, 1], y: [3, 5], z: [2.5, 5] },
  },
  {
    key: 'avoid',
    label: 'RED FLAG',
    desc: 'Culera zone — proceed at your own risk',
    color: '#f4506b',
    badge: '⚑',
    bounds: { x: [-5, 5], y: [-5, 5], z: [-5, -2.5] },
  },
];

// Return the first sweet spot a point falls inside, or null.
export function zoneFor({ x, y, z }) {
  for (const s of SWEET_SPOTS) {
    const b = s.bounds;
    if (x >= b.x[0] && x <= b.x[1] && y >= b.y[0] && y <= b.y[1] && z >= b.z[0] && z <= b.z[1]) {
      return s;
    }
  }
  return null;
}

// Tier label + a stable key used for coloring.
export function tierFor(score) {
  if (score >= 8.5) return { label: 'S', key: 's' };
  if (score >= 7.0) return { label: 'A', key: 'a' };
  if (score >= 5.5) return { label: 'B', key: 'b' };
  if (score >= 4.0) return { label: 'C', key: 'c' };
  return { label: 'D', key: 'd' };
}

// A color along a red -> amber -> green ramp based on a 0..10 score.
// Returns a CSS hex string.
export function scoreColor(score) {
  const t = clamp(score / 10, 0, 1);
  // hue 0 (red) -> 130 (green)
  const hue = t * 130;
  return hslToHex(hue, 70, 52);
}

function hslToHex(h, s, l) {
  s /= 100;
  l /= 100;
  const k = (n) => (n + h / 30) % 12;
  const a = s * Math.min(l, 1 - l);
  const f = (n) => {
    const color = l - a * Math.max(-1, Math.min(k(n) - 3, Math.min(9 - k(n), 1)));
    return Math.round(255 * color)
      .toString(16)
      .padStart(2, '0');
  };
  return `#${f(0)}${f(8)}${f(4)}`;
}
