# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Communication style

When making code changes, briefly describe what changed and why — do not list removed/added/edited lines.

## Code style

Prefer `var` declarations over `:=` everywhere.

## Build & Run

```bash
go build ./...          # build
go run .                # run
./build-linux-to-linux.bash   # release build
```

No tests exist in this project.

## Architecture

This is a terminal pixel-art renderer built on `termbox-go`. Tiles from a PNG atlas are blitted directly into the termbox cell buffer using Unicode block characters as sub-cell pixels.

**Why 12×12 tiles?** 12 divides evenly by 1, 2, 3, and 4 — so every supported cell grid fits without partial cells: 2×1 (half blocks), 2×2 (quarter blocks), 2×3 (sextants), 2×4 (braille).

### Rendering pipeline

```
SetTile(col, row, tileID, color, tcw, tch, blit, atlas, atlasW)
  → extracts 12×12 pixels from atlas into scratch buffer, colorized
  → calls blit(pixels, 12, 12, col*tcw, row*tch)
    → groups pixels into cell-sized blocks
    → computes bitmask of lit pixels per cell
    → maps bitmask → Unicode rune via *Rune function
    → termbox.SetCell(...)
```

The terminal cell buffer doubles as the tile state store — `GetTile` reads color back from `termbox.CellBuffer()` with no separate data structure.

### Blit modes (all in `blit.go`)

| BlitFunc | Grid | Rune func | Unicode range |
|---|---|---|---|
| `BlitHalf` | 1×2 | `halfblockRune` | U+2580–U+2584 |
| `BlitQuarter` | 2×2 | `quarterblockRune` | U+2596–U+259F |
| `Blit` | 2×3 | `sextantRune` | U+1FB00–U+1FB3B |
| `BlitBraille` | 2×4 | `brailleRune` | U+2800–U+28FF |

The active mode and matching `tcw`/`tch` (terminal cells per tile) are chosen together in `main()`.

### Key encoding details

- **Sextants** (`sextantRune`): bitmask 0–63, skips values 21, 42, 63 which duplicate existing half/full block chars.
- **Braille** (`brailleRune`): row-major input bits are remapped to braille dot order (`{0,3,1,4,2,5,6,7}`) before indexing into U+2800.

### Files

- `main.go` — tile API (`SetTile`, `GetTile`, `ClearTile`) and `main()`
- `blit.go` — all `*Rune` mask functions and `Blit*` renderers
- `png.go` — `LoadPNG`: decodes a PNG, maps opaque pixels → `ColorWhite`, transparent → `ColorDefault`
