# HerbGraph

A Palantir-styled, dependency-free **3D ranking app**. Plot people across
three slang axes, see them in a rotatable 3D space, and get an overall
rating with a tier label and "sweet spot" zones.

No build step, no frameworks, no external libraries — the 3D scene is a
custom canvas renderer, so it works offline and on any modern browser
(including iPhone Safari).

## The three dimensions

| Axis | – pole | + pole |
|------|--------|--------|
| X (pink)  | Naca   | Fresa |
| Y (blue)  | Basic  | Artsy |
| Z (green) | Culera | Chida |

Each axis is a slider from **-5** to **+5**. The ideal corner is
`(Fresa, Artsy, Chida)`.

## Overall rating

The score (0–10) is a weighted blend of the three axes, normalized so that
`-5 → 0` and `+5 → 10` on each axis:

```
overall = 10 * ( 0.30*norm(x) + 0.25*norm(y) + 0.45*norm(z) )
norm(v) = (v + 5) / 10
```

**Chida–Culera** carries the most weight (0.45) — personality over aesthetics.

| Tier | Score |
|------|-------|
| S | ≥ 8.5 |
| A | ≥ 7.0 |
| B | ≥ 5.5 |
| C | ≥ 4.0 |
| D | < 4.0 |

## Sweet spots

Named regions of the space, rendered as translucent zones in the 3D view
and flagged on each subject:

| Zone | Region | Meaning |
|------|--------|---------|
| ◆ ELITE   | Fresa+, Artsy+, Chida+ | the dream corner |
| ✦ ALT BAE | Artsy+, Chida+, naca-leaning | the cool alt type |
| ⚑ RED FLAG | very Culera | proceed at your own risk |

Weights, thresholds and zones all live in `rating.js` if you want to tune them.

## Features

- Add / edit / delete subjects with name + optional **Instagram handle**
  (click ◎ in the ranking to open the profile) + 3 sliders, with a live
  score / tier / zone preview.
- Custom **3D scatter plot**: rotate (drag), zoom (pinch / scroll / buttons),
  reset view, toggle sweet-spot zones. Points are color-coded by score with
  labels and floor stems for depth.
- Tap a point (or a row) to select; hover for a tooltip.
- Sorted **leaderboard** + a top stats bar (count, average, top subject).
- **Fullscreen** button — great for AirPlay screen-mirroring to Apple TV.
- Responsive layout for phone screens.
- Everything persists locally via `localStorage`.

## Running locally

It uses ES modules, so serve it over HTTP (don't open the file directly):

```bash
cd apps/herbgraph
python3 -m http.server 8080
# open http://localhost:8080
```

## Viewing on your iPhone / Apple TV

The repo ships a GitHub Pages workflow (`.github/workflows/deploy-herbgraph.yml`)
that publishes this folder as a website.

**One-time setup:** in the GitHub repo → **Settings → Pages**, set
*Build and deployment → Source* to **GitHub Actions**. The workflow then runs
on every push that touches `apps/herbgraph/` and gives you a URL like:

```
https://<your-user>.github.io/aoc/
```

Open that on your iPhone, tap the ⤢ fullscreen button, then use AirPlay
screen mirroring to cast to your Apple TV.
