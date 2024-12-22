---
title: Content Types
description: An overview of the different content types supported by Hermes.
filepath: ""
summary: ""
date: "2024-11-01"
published-at: "2024-11-01T10:58:21+01:00"
created-at: "2024-12-15T14:04:02+01:00"
updated-at: "2024-12-22T09:27:12+01:00"
type: article
section: root
slug: content-types
image: ""
social-image: ""
layout: article
canonical-url: ""
locale: en_US
robots: ""
excerpt: ""
permalink: ""
draft: false
table-of-contents: false
share: false
featured: false
comments: false
author: []
categories: []
tags: []
keywords: []
sitemap:
    priority: 0
    changefreq: ""
debug: false
---

We have different content types to categorize the content on the site:

## Page

Used for site pages, usually for standard static pages like about-us, contact, etc. These are not considered "real" content but more like standard static pages.

## Article

Used for real content. Some articles are not intended to be chronologically accessible as in a blog, for example, an essay or a paper. This type makes sense for this kind of posts/content.

## Blog

Self-explanatory, content ordered chronologically.

## Series

Sometimes we want to group articles in series, such as a guide, a tutorial, or a course. A series of connected and self-contained posts that make sense to organize and make accessible in a series of connected posts. We use the article type for this.

## Example

- `content/root/about-us.md` will be rendered as `domain.tld/about-us.html`.
- `content/guides/article.md` will be rendered as `domain.tld/guides/article.html`.
