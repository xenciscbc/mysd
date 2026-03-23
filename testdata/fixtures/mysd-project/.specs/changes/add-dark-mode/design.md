## Architecture

Theme support is implemented using CSS custom properties (variables) at the root level, toggled via a data attribute on the html element.

## Key Decisions

- CSS variables over class-based theming for simpler component code
- localStorage for preference persistence (no backend needed for v1)
- prefers-color-scheme media query for OS detection
