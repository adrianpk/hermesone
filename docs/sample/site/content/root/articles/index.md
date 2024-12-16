---
title: "Hermes - Static Site Generator (SSG)"
description: "An introduction to Hermes, the ultimate content creation tool."
summary: "An overview of Hermes, a static site generator designed to create fast, secure, and maintainable websites."
date: "2024-11-01"
published-at: "2024-11-01T10:58:21+01:00"
created-at: "2024-11-15T10:58:21+01:00"
updated-at: "2024-11-16T02:04:50+01:00"
type: "article"
section: root
slug: "hermes-ssg"
image: "https://example.com/logo.png"
social-image: "https://example.com/social-logo.png"
layout: "article"
canonical-url: "https://example.com/hermes-ssg"
locale: "en"
robots: "index, follow"
excerpt: "An overview of Hermes, a static site generator designed to create fast, secure, and maintainable websites."
permalink: "/hermes-ssg/"
draft: false
table-of-contents: true
share: true
featured: true
comments: true
author: ["Adrian PK"]
categories: ["Content Creation", "Tools"]
tags: ["Hermes", "Content Creation", "SSG"]
keywords: ["Hermes", "Content Creation", "SSG"]
sitemap:
    priority: 0.8
    changefreq: "monthly"
debug: false
---

# Hermes - Static Site Generator (SSG)

Welcome to the **Hermes** documentation site. This site provides comprehensive information about Hermes, a static site generator designed to help you create fast, secure, and easily maintainable websites.

## Sections

Below are the main sections of the documentation. Click on each link to navigate to the respective page.

- [Overview](overview.md)
- [Commands](commands.md)
- [Content Structure](content-structure.md)
- [Content Types](content-types.md)
- [Content Type Paths](content-type-paths.md)
- [Layout Selection Logic](layout-selection-logic.md)
- [Image Storage](image-storage.md)
- [Index Generation](index-generation.md)
- [Future Plans](future-plans.md)

## Overview

Hermes is a static site generator (SSG) designed to help you create fast, secure, and easily maintainable websites. It takes your content written in Markdown and transforms it into a static website using predefined layouts and templates.

For more details, visit the [Overview](overview.md) page.

## Commands

Hermes provides several commands to help you manage your site. These include:

- **`init`**: Initializes a Hermes site project.
- **`gen`**: Generates the HTML files in the `output` directory.
- **`upgrade`**: Updates layout files to the latest version.
- **`new`**: Generates new empty files for different sections.
- **`publish`**: Publishes the rendered content to GitHub Pages.
- **`backup`**: Versions the content in GitHub.

For more details, visit the [Commands](commands.md) page.

## Content Structure

The content structure in Hermes is organized under the `content` directory. This directory contains Markdown files that will be converted to HTML.

For more details, visit the [Content Structure](content-structure.md) page.

## Content Types

Hermes supports different content types to categorize the content on your site. These include Page, Article, Blog, and Series.

For more details, visit the [Content Types](content-types.md) page.

## Content Type Paths

The paths for different content types are structured to ensure proper rendering under the domain.

For more details, visit the [Content Type Paths](content-type-paths.md) page.

## Layout Selection Logic

Hermes uses a flexible layout selection logic to customize the appearance of your site at various levels.

For more details, visit the [Layout Selection Logic](layout-selection-logic.md) page.

## Image Storage

Hermes handles image storage and path rewriting during the build process to ensure correct linking in the generated HTML files.

For more details, visit the [Image Storage](image-storage.md) page.

## Index Generation

Hermes provides a mechanism for generating index pages for the root and sections. If an `index.md` file is present and published, it will be used as the index page. Otherwise, Hermes generates an index page with a paginated list of content.

For more details, visit the [Index Generation](index-generation.md) page.

## Future Plans

Hermes is continuously evolving. Stay tuned for updates and feel free to contribute or provide feedback.

For more details, visit the [Future Plans](future-plans.md) page.

---

This document is a living document and will be updated regularly as Hermes evolves.
