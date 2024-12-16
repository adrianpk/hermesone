---
title: "Layout Selection Logic in Hermes"
description: "An explanation of the layout selection logic used in Hermes."
date: "2024-11-01"
type: "article"
section: root
slug: "layout-selection-logic"
layout: "article"
---

# Layout Selection Logic

## Default Layouts

By default, Hermes uses a simple and straightforward approach for selecting layouts:

- **Global Default Layout**: All posts use the default layout stored at `layout/default.html`.
- **Content-Type Specific Default Layouts**: You can also have default layouts for specific content types (page, article, blog, series) stored at `layout/page/default.html`, `layout/article/default.html`, `layout/blog/default.html`, and `layout/series/default.html`.

By default, all these default layouts are the same, providing a convenient way for users to customize each one if required. If these specific defaults are not present, Hermes will fall back to using `layout/default.html`.

## Custom Layouts for Specific Content

For most users, the default layouts are sufficient. However, if you need more granular control over the layout used for specific content, Hermes provides a flexible mechanism:

- **Specific Content Layouts**: You can create a layout for a specific content item by naming the layout file after the content item. For example, `layout/default/article/article-name.html` will be used to render an article with the name `article-name.md`. This applies to other content types as well.

## Section-Specific Layouts

In addition to global and content-type specific layouts, you can also define layouts for specific sections:

- **Section Default Layouts**: Create a layout for a specific section by placing a default layout in the section's directory. For example, `layout/section-name/default.html` will be used for all content in the `section-name` section.
- **Section-Specific Content-Type Layouts**: Similar to global content-type layouts, you can have section-specific layouts for each content type. For example, `layout/section-name/article/default.html` will be used for all articles in the `section-name` section.
- **Section-Specific Content Layouts**: You can also create layouts for specific content items within a section. For example, `layout/section-name/article/article-name.html` will be used to render an article named `article-name.md` within the `section-name` section.

## Summary

- **Global Default**: `layout/default.html`
- **Content-Type Defaults**: `layout/page/default.html`, `layout/article/default.html`, `layout/blog/default.html`, `layout/series/default.html`
- **Specific Content Layouts**: `layout/default/article/article-name.html`
- **Section Default**: `layout/section-name/default.html`
- **Section-Specific Content-Type Layouts**: `layout/section-name/article/default.html`
- **Section-Specific Content Layouts**: `layout/section-name/article/article-name.html`

This flexible layout selection logic allows you to customize the appearance of your site at various levels, from global defaults to specific content items within sections.

**This may sound overwhelming, but remember, you can rely on the defaults to start and customize later as needed. Most of the time you won't need this granular control.**

---

This document explains the layout selection logic used in Hermes.

---
