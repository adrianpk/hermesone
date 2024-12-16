---
title: "Index Generation in Hermes"
description: "An explanation of how index generation works in Hermes."
date: "2024-11-01"
type: "article"
section: root
slug: "index-generation"
layout: "article"
---

# Index Generation

This document explains how the index generation works in Hermes, including the current implementation and future plans.

## Current Implementation

### Creating an Index Page

For both the root and sections, you can create an article or page called `index.md`. If its publishing date is in the past and it is not marked as a draft, it will be used as the index (`index.html`) page for that section.

### Automatic Index Generation

If there is no `index.md` file, or if it is not published, Hermes will automatically generate an index page. This generated index page includes a paginated chronological list from newest to oldest of all content published under that section (articles, blogs, and series).

## Future Plans

### Customizable Layout

The layout of the automatically generated index will be customizable. This will allow users to define how the index page should look and feel.

### Pagination Parameters

Hermes will support customizable pagination parameters. This will give users more control over how content is paginated on the index page.

## Summary

- **Manual Index**: Create an `index.md` file with a past publishing date and not marked as a draft to use it as the index page.
- **Automatic Index**: If no `index.md` is available or published, Hermes generates an index page with a paginated list of content.
- **Future Enhancements**: Customizable layout and pagination parameters are planned for future updates.

---

This document provides an overview of how index generation works in Hermes.

---
