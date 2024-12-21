---
title: Image Storage
description: An explanation of how image storage and path rewriting are handled in Hermes.
filepath: ""
summary: ""
date: "2024-11-01"
published-at: "2024-11-01T10:58:21+01:00"
created-at: "2024-12-15T14:16:19+01:00"
updated-at: "2024-12-22T15:33:45+01:00"
type: article
section: root
slug: image-storage
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

Hermes handles image storage and path rewriting during the build process to ensure correct linking in the generated HTML files.

## Image Storage

Regardless of where your `.md` file is located under the `content` directory (root, section, or content type), all images for a particular post should be stored in a sibling directory with the same name as the content file. For example, if your content file is `content-name.md`, then its images should be stored in a directory called `content-name`.

## Referencing Images

You can reference images using the standard Markdown syntax. For example:

```markdown
![An example image](index/example.png "Example Image")
```

In this example, `index.md` is the content file, and `example.png` is an image stored in the `index` directory, which is a sibling to `index.md`.

## Image Path Handling

Hermes will automatically copy the images and rewrite the paths accordingly during the build process. This ensures that your images are correctly linked in the generated HTML files.

## Header Images

For now, if you need to insert a header image for the post, you need to rely on the same method of storing and referencing images as described above. However, it is planned that Hermes will support a default header image picked using a convention-over-configuration approach. Additionally, you will be able to override this default using front matter metadata in your Markdown files.

## Example

Assume you have the following structure:

```
content/
└── guides/
    ├── article.md
    └── article/
        └── example.png
```

In `article.md`, you can reference `example.png` as follows:

```markdown
![An example image](article/example.png "Example Image")
```

Hermes will handle copying the image and updating the path in the generated HTML.

This approach helps keep your images organized and ensures that they are easily accessible and correctly linked in your site.

---

This document explains how image storage and path rewriting are handled in Hermes.

---
