---
title: Hermes Commands
description: A detailed guide to the commands available in Hermes.
filepath: ""
summary: ""
date: "2024-11-01"
published-at: "2024-11-01T10:58:21+01:00"
created-at: "2024-12-15T13:57:08+01:00"
updated-at: "2024-12-22T09:27:12+01:00"
type: article
section: root
slug: commands
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

# Commands

This section provides a quick introduction to the available commands in Hermes. We will go deeper into each command with more detail later.

## `init`

Initializes a Hermes site project. This command sets up the necessary directory structure and configuration files for your site.

```bash
$ hermes init
```

## `gen`

Generates the HTML files in the `output` directory using the Markdown and layout files. It also takes care of copying the images and rewriting the paths accordingly.

```bash
$ hermes gen
```

## `upgrade`

If you upgrade Hermes and the new version includes updates to the layout files, this command lets you replicate those updates in your site. A timestamped backup of the current layout files is created in place to prevent accidental overwrites.

```bash
$ hermes upgrade
```

## `new`

Generates new empty files for different sections, content types, etc. While you can still create the files manually, this command provides an easy way to create a normalized Markdown file with front matter metadata in the proper location.

```bash
$ hermes new <type> <name>
```

## `publish`

Publishes the rendered content to GitHub Pages. This command automates the process of deploying your static site to GitHub Pages.

```bash
$ hermes publish
```

## `backup`

Versions the content in GitHub. This command helps you keep track of changes to your content by creating backups in a separate repository.

```bash
$ hermes backup
```

Worth mentioning that content is versioned in a different repository than the final HTML content.

---

This document provides a detailed guide to the commands available in Hermes.

---
