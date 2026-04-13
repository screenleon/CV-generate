---
name: design-to-code
description: Use when the user provides screenshots, mockups, or design assets and wants implementation that stays visually close to the source.
---

# Design To Code

Use this skill for image-led UI implementation.

## Pre-implementation checklist

Before writing any code:

1. Read the existing component library or UI pattern files in the codebase
2. Check `DECISIONS.md` for prior decisions on UI framework, token system, or component conventions
3. Check `docs/operating-rules.md` → Project-specific constraints for any CSS/styling rules
4. Confirm with `demand-triage` whether this task also involves logic or flow changes — if so, route through `feature-planner` first

## Design breakdown (required before coding)

Break down the design into these layers before writing a single line of code:

1. **Layout** — identify the grid, flexbox, or other layout system. Note container widths, alignment, and nesting depth. Map to existing layout components or patterns in the codebase.
2. **Typography** — font family, size, weight, line-height, letter-spacing for each text element. Map to the project's type scale or design tokens if they exist.
3. **Color** — exact color values for backgrounds, text, borders, and interactive elements. Map to the project's color tokens or variables. Flag any color not in the existing palette.
4. **Spacing** — margins, paddings, gaps between elements. Map to the project's spacing scale if it exists.
5. **Radius, border, shadow** — border-radius, border widths and styles, box-shadows. Note which elements have them and which do not.
6. **Interaction states** — default, hover, focus, active, disabled, loading, error, and empty states. If the design does not show all states, infer reasonable defaults and call them out.
7. **Responsive or device differences** — if the design shows multiple breakpoints, note the differences. If only one size is shown, state the assumed breakpoint and note what might need adaptation.

## Matching vs. inferring

At the end of the breakdown, explicitly separate:

- **Exact matches** — values and layouts clearly visible in the design
- **Inferences** — values you estimated because the design was incomplete or ambiguous

This lets the reviewer know what to double-check. Never silently infer values — always surface them.

## Scale adaptation

Use `demand-triage` to classify the task before starting:

| Scale | Typical scenario | Additional steps |
|-------|-----------------|------------------|
| **Small** | Single state change, copy update, or color adjustment on an existing component | Inline preamble; skip planning agent |
| **Medium** | New component following existing patterns; 2–5 files | Follow standard Medium workflow |
| **Large** | New screen or flow, design system changes, or redesign affecting many components | Planning agent first; critic review; risk-reviewer |

If the design requires **new state management, routing, or API calls**, reclassify to Medium or Large and route through `feature-planner` before implementation.

## Implementation guidelines

- **Reuse existing components** — before creating a new component, check if one already exists that can be extended or composed.
- **Use project tokens** — prefer design tokens, CSS variables, or theme values over hard-coded colors, sizes, and fonts.
- **Mobile-first is default** — unless the design clearly targets desktop-first, build for narrow viewports and scale up.
- **Accessibility minimums** — ensure sufficient color contrast (WCAG AA), semantic HTML elements, keyboard navigability, and alt text for images.
- **Do not redesign** — if the design is inconsistent with the project's existing patterns, implement the design as-is and flag the inconsistency in your output. Do not silently "fix" the design.

## Anti-patterns

| Anti-pattern | Why it is wrong | What to do instead |
|-------------|----------------|--------------------|
| Replacing the provided design with a generic template component | Loses visual fidelity; the design is the spec | Implement what is shown; reuse structure only where the design matches |
| Abstracting too early | Creates a "reusable" component from a one-off screen before the pattern is proven | Implement as a single-use component; promote to shared only after second use |
| Hard-coding pixel values | Breaks when the design system changes; diverges from the project's spacing scale | Map to the nearest existing token or CSS variable |
| Silently inferring interaction states | Leads to missing hover/focus/disabled states the reviewer will not notice until QA | Always list inferred states in the Matching vs. Inferring section |
| Skipping responsive behavior | The component breaks on mobile even if the design was desktop-only | State the assumed breakpoint and at minimum add a mobile layout note |
| Implementing new data fetching or state | Mixes UI fidelity work with logic/API work; increases blast radius | Stop and reclassify with `demand-triage`; route through feature-planner |

## Use this skill when

- the user provides screenshots, mockups, or design files
- the task is primarily about visual fidelity
- the deliverable is a UI component or page that must match a visual reference

## Do not use this skill when

- the request does not include a visual reference (use `application-implementer` instead)
- the design also requires new API endpoints, schema changes, or auth changes — plan those with `feature-planner` first, then use this skill for the UI phase only
- the primary deliverable is a logic change that happens to touch a component

## Conformance self-check

Before marking implementation complete, verify:

- [ ] Design breakdown (all 7 layers) was produced before writing code
- [ ] Exact matches and inferences are explicitly listed and separated
- [ ] Demand triage was run and scale classification is stated
- [ ] Existing components were checked before creating a new one
- [ ] Project design tokens or CSS variables are used, not hard-coded values
- [ ] All interaction states (hover, focus, active, disabled, loading, error, empty) are handled
- [ ] Accessibility minimums checked: contrast, semantic HTML, keyboard nav, alt text
- [ ] If design is inconsistent with existing patterns, inconsistency is flagged (not silently fixed)
- [ ] Responsive behavior is addressed or explicitly noted as out of scope
- [ ] If new logic or API work was discovered, reclassification was triggered and task scope was re-approved
