# Asona Design System

## Overview
A clear, intuitive, and highly functional interface tailored for a project management and real-time chat application. The design emphasizes clean lines, balanced whitespace, and high information density to minimize visual noise while maximizing productivity.

## Colors
- **Primary** (`#10b981`): Main branding, calls to action (CTAs), primary buttons, and active interactive elements.
- **Secondary** (`#3b82f6`): Supporting UI elements, secondary actions, informational badges, and links.
- **Surface** (`#ffffff`): Default page backgrounds and main content areas (Light mode). Dark mode equivalent: `#0f172a`.
- **Surface Variant** (`#f3f4f6`): Card backgrounds, subtle delineations, and secondary panels (Light mode). Dark mode equivalent: `#1e293b`.
- **On-surface** (`#111827`): Primary text color on default surface backgrounds (Light mode). Dark mode equivalent: `#f8fafc`.
- **On-surface Variant** (`#6b7280`): Secondary text, disabled states, and subtle icons (Light mode). Dark mode equivalent: `#94a3b8`.
- **Error** (`#ef4444`): Destructive actions, validation errors, and critical alerts.

## Typography
- **Font Families**: Inter for both Headings and Body text.
- **Headlines**: Inter, Semi-bold (600), with tight tracking for an authoritative yet modern look.
- **Body**: Inter, Regular (400), typically 14px-16px, optimized for long-form readability.
- **Labels**: Inter, Medium (500), 12px, often uppercase with generous letter spacing for section headers, eyebrows, and tags.

## Components & Style Patterns
- **Rounding**: 
  - Standard elements (Cards, Modals): Rounded (`8px` to `12px`).
  - Small elements (Buttons, Inputs): Rounded (`6px` to `8px`).
  - Badges/Avatars: Fully rounded (`9999px`).
- **Borders**: Minimal use of `1px` subtle borders (`#e5e7eb` in Light mode) to define edges without heavy box-shadows.
- **Elevation**: Flat design preference. Use shadows sparingly (e.g., `0 4px 6px -1px rgb(0 0 0 / 0.1)`) only for floating elements like dropdowns, modals, or focused interactive states.
- **Interaction States**: 
  - Hover: Subtly darken or lighten the background color. 
  - Focus: Clear visible focus rings (`2px` outline with `2px` offset) on all interactive elements for accessibility.

## Do's and Don'ts
- **Do** maintain a strict contrast ratio (minimum 4.5:1) for all text elements.
- **Do** use the Primary color sparingly to draw the user's eye to the most important action on a page.
- **Don't** mix vastly different corner radii (e.g., mixing sharp `0px` corners with `16px` rounded corners) within the same contiguous view.
- **Don't** use more than three font weights across the platform to maintain load performance and visual consistency.
