# Yuki Buy Log — Client

Mobile-first SPA for tracking purchases and expenses.

## Tech Stack

- **Svelte 5** (runes, NOT SvelteKit) with TypeScript
- **Vite 7** as build tool
- **melt** (v0.44+) — headless UI component library for Svelte 5 (import from `melt/builders`)
- **lucide-svelte** — icon library

## Project Structure

```
client/
├── src/
│   ├── main.ts          # App entry point, mounts App.svelte to #app
│   ├── app.css           # Global styles (reset, dark/light theme)
│   ├── App.svelte        # Root component with bottom navigation (Tabs from melt)
│   └── pages/            # Screen components, one per menu tab
│       ├── HomePage.svelte
│       ├── SearchPage.svelte
│       ├── AddPage.svelte
│       ├── ListPage.svelte
│       └── ProfilePage.svelte
├── index.html            # HTML shell with mobile viewport meta tags
├── vite.config.ts
├── svelte.config.js
├── tsconfig.json
└── package.json
```

## Architecture

- App behaves like a native mobile app: fixed full-screen layout, no page scroll, bottom nav always visible.
- Navigation uses melt `Tabs` builder — each menu item is a tab trigger, each page is tab content.
- `App.svelte` exposes `navigateTo(tab: string)` function for programmatic navigation between tabs.
- Pages are loaded eagerly (no lazy loading) as plain Svelte components.
- Layout uses `100dvh` + `position: fixed` to prevent browser chrome issues on mobile.
- Safe area insets (`env(safe-area-inset-bottom)`) are respected for notched devices.

## Conventions

- Use Svelte 5 runes (`$state`, `$derived`, `$effect`) — no legacy reactive statements.
- Keep all page-level components in `src/pages/`.
- Reusable UI widgets go in `src/lib/` (create as needed).
- Avoid adding dependencies beyond `melt` and `lucide-svelte`.
- Use `<script lang="ts">` in all `.svelte` files.
- All interface should use english language.

## Design Philosophy

The visual language is **white canvas with intentional accent**:

- Surfaces (`--color-surface`, `--color-bg`) form the neutral foundation — clean white/off-white.
- Accent colors appear **only where they carry meaning** — selected state, required action, status signal, interactive focus.
- At any moment no more than ~20% of the visible UI should use an accent color.
- Depth and hierarchy are communicated through **shadows and surface layering**, not heavy colors.

### Color Semantics

Each accent has a single semantic role. Do not use accent colors decoratively or interchangeably.

| Token | Role | Use cases |
|---|---|---|
| `--color-blue` | Primary interactive | Focus rings, selected/active state, primary buttons, tags/pills, links |
| `--color-green` | Positive / success | Save confirmations, completed status, positive metrics |
| `--color-red` | Destructive / error | Delete buttons, error validation, warning badges |
| `--color-yellow` | Attention / caution | Pending state, draft status, uncertain/missing data |
| `--color-dark-blue` | Strong emphasis | Page/section headers, brand marks — use sparingly |

### Tinting Rule

Never fill large areas with a solid accent color. Use `color-mix` to create tinted versions:

```css
/* 8–12% — subtle selected/hovered background */
background: color-mix(in srgb, var(--color-blue) 10%, transparent);

/* 20–30% — visible accent border */
border-color: color-mix(in srgb, var(--color-blue) 25%, transparent);

/* 100% — only for small icons, text labels, or tiny elements */
color: var(--color-blue);
```

### Interactive State Conventions

Every interactive element must visually respond to all relevant states:

- **Default** — `--color-border` border, `--color-surface` background, `--color-text` label.
- **Hover** — lighten background to `--color-bg`; no border change needed.
- **Focus** — set `border-color: var(--color-blue)` and add `box-shadow: var(--focus-ring)`.
- **Active / selected** — accent-colored text or a left-side 2px accent border stroke; tinted background (8–12%).
- **Disabled** — `--color-disabled` text and border; no hover or focus effects.

### Elevation Model

Use shadows to show which layer a surface belongs to — not for decoration.

- **No shadow** — flat inline elements (labels, dividers, form fields).
- `--shadow-sm` — subtly raised inline cards or highlighted items.
- `--shadow-md` — dropdowns, popovers, floating cards.
- `--shadow-lg` — modals, bottom sheets, overlays.

### Typography Conventions

- Page titles: `--text-lg`, `font-weight: 600`, `--color-text`.
- Section/card headers: `--text-md`, `font-weight: 600`, `--color-text`.
- Body / input text: `--text-base`, `font-weight: 400`, `--color-text`.
- Labels above inputs: `--text-sm`, `font-weight: 500`, `--color-text-secondary`.
- Secondary/helper text, pills: `--text-sm`, `--color-text-secondary`.
- Captions, timestamps, meta: `--text-xs`, `--color-disabled`.

---

## Styling Rules

- **Use CSS variables from `:root`** — never hardcode colors or magic numbers. All tokens are in `app.css`:
  - Colors: `--color-bg`, `--color-surface`, `--color-text`, `--color-text-secondary`, `--color-disabled`, `--color-border`
  - Accent colors: `--color-red`, `--color-yellow`, `--color-blue`, `--color-dark-blue`, `--color-green`
  - Radii: `--radius-xs`, `--radius-sm`, `--radius-md`, `--radius-lg`, `--radius-full`
  - Shadows: `--shadow-sm`, `--shadow-md`, `--shadow-lg`
  - Transitions: `--transition-fast`, `--transition-base`
  - Focus ring: `--focus-ring`
  - Font sizes: `--text-xs`, `--text-sm`, `--text-base`, `--text-md`, `--text-lg`
  - Spacing scale: `--space-1` (2px) · `--space-2` (4px) · `--space-3` (6px) · `--space-4` (8px) · `--space-5` (12px) · `--space-6` (16px) · `--space-7` (20px) · `--space-8` (24px)
  - Component sizes: `--page-padding` (4px) · `--card-padding` (12px) · `--section-gap` (12px) · `--form-gap` (8px) · `--label-gap` (4px) · `--input-height` (38px) · `--input-padding` (8px 12px) · `--nav-height` (52px)
- **Apply radius tokens by element type:**
  - Inputs, selects, textareas, small cards → `var(--radius-md)`
  - Pills, chips, small badges, icon buttons → `var(--radius-sm)` or `var(--radius-xs)`
  - Large cards, modals, bottom sheets, dropdowns → `var(--radius-lg)`
  - Avatars, toggles → `var(--radius-full)`
- **Apply transitions consistently** — use `var(--transition-fast)` for color/bg/opacity, `var(--transition-base)` for border/shadow/transform. Never use ad-hoc duration values.
- **Never show scrollbars** — this is a SPA with a native-app feel. Every scrollable element must hide its scrollbar: `scrollbar-width: none` (Firefox) and `::-webkit-scrollbar { display: none }` (Chrome/Safari). No scrollbar should ever be visible anywhere in the app.
- **Widgets must not control their own size.** A widget does not set `width` or `height` on itself. The parent is responsible for sizing via grid, flex, or explicit dimensions.
- **Keep padding small** — widgets target small mobile screens, use `padding: 4px` as the default inside cards/widgets.

## Commands

```bash
npm run dev       # Start dev server
npm run build     # Production build
npm run preview   # Preview production build
npm run check     # Type-check with svelte-check + tsc
```