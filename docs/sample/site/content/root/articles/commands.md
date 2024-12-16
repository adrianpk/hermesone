---
title: "Hermes Commands"
description: "A detailed guide to the commands available in Hermes."
date: "2024-11-01"
type: "article"
section: root
slug: "commands"
layout: "article"
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
