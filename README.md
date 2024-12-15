# Hermes - Static Site Generator (SSG)

## Work In Progress (WIP)

This document is a work in progress. It serves as a way to capture ideas, behaviors, features, and other relevant information about Hermes. The content here is subject to change and will be polished and organized in a better way later.

## Overview

Hermes is a static site generator (SSG) designed to help you create fast, secure, and easily maintainable websites. It takes your content written in Markdown and transforms it into a static website using predefined layouts and templates.

## Commands

This section provides a quick introduction to the available commands in Hermes. We will go deeper into each command with more detail later.

- **`init`**: Initializes a Hermes site project.

- **`gen`**: Generates the HTML files in the `output` directory using the Markdown and layout files. It also takes care of copying the images and rewriting the paths accordingly.

- **`upgrade`**: If you upgrade Hermes and the new version includes updates to the layout files, this command lets you replicate those updates in your site. A timestamped backup of the current layout files is created in place to prevent accidental overwrites.

- **`new`**: Generates new empty files for different sections, content types, etc. While you can still create the files manually, this command provides an easy way to create a normalized Markdown file with front matter metadata in the proper location.

- **`publish`**: Publishes the rendered content to GitHub Pages.

- **`backup`**: Versions the content in GitHub.

Worth mentioning that content is versioned in a different repository than the final HTML content.

## Content Structure

Under the `content` directory, we store the Markdown files that will later be converted to HTML.

- **Root Directory**: The root directory defines the content for the root section of the site. This content is rendered directly under the domain root (e.g., `domain.tld`) and not under `domain.tld/root`.
- **Non-Root Sections**: Any other directory under `content` defines a non-root section. These sections are rendered under `domain.tld/section-name`.

### Example

- `content/root/index.md` will be rendered as `domain.tld/index.html`.
- `content/guides/article.md` will be rendered as `domain.tld/guides/article.html`.


## Content Types

We have different content types to categorize the content on the site:

- **Page**: Used for site pages, usually for standard static pages like about-us, contact, etc. These are not considered "real" content but more like standard static pages.
- **Article**: Used for real content. Some articles are not intended to be chronologically accessible as in a blog, for example, an essay or a paper. This type makes sense for this kind of posts/content.
- **Blog**: Self-explanatory, content ordered chronologically.
- **Series**: Sometimes we want to group articles in series, such as a guide, a tutorial, or a course. A series of connected and self-contained posts that make sense to organize and make accessible in a series of connected posts. We use the article type for this.

## Content Type Paths

The paths for different content types are as follows:

- **Pages and Articles**: These are rendered at the section level. For the root section, they are rendered directly under the domain. If they belong to another non-root section, they are rendered under the section path.
  - Example: `content/root/about-us.md` will be rendered as `domain.tld/about-us.html`.
  - Example: `content/guides/article.md` will be rendered as `domain.tld/guides/article.html`.

- **Blog**: The blog content type is accessible differently.
  - For the root blog, it is accessible under `domain.tld/blog`.
  - For sections, it is accessible under `domain.tld/section/blog`.
  - Example: `content/root/blog/post.md` will be rendered as `domain.tld/blog/post.html`.
  - Example: `content/guides/blog/post.md` will be rendered as `domain.tld/guides/blog/post.html`.

- **Series**: For now, there can be some misalignment between what we say here and the actual behavior because the series is a work in progress and we are experimenting with the best way to do it. The idea is that they are accessible:
  - For the root, under `domain.tld/series/series-name/article-name-nn` where `nn` is the delivery number of each article of the series.
  - For sections, it is the same but also includes the section name.
  - Example: `content/root/series/guide/article-01.md` will be rendered as `domain.tld/series/guide/article-01.html`.
  - Example: `content/guides/series/guide/article-01.md` will be rendered as `domain.tld/guides/series/guide/article-01.html`.

## Layout Selection Logic

### Default Layouts

By default, Hermes uses a simple and straightforward approach for selecting layouts:

- **Global Default Layout**: All posts use the default layout stored at `layout/default.html`.
- **Content-Type Specific Default Layouts**: You can also have default layouts for specific content types (page, article, blog, series) stored at `layout/page/default.html`, `layout/article/default.html`, `layout/blog/default.html`, and `layout/series/default.html`.

By default, all these default layouts are the same, providing a convenient way for users to customize each one if required. If these specific defaults are not present, Hermes will fall back to using `layout/default.html`.

### Custom Layouts for Specific Content

For most users, the default layouts are sufficient. However, if you need more granular control over the layout used for specific content, Hermes provides a flexible mechanism:

- **Specific Content Layouts**: You can create a layout for a specific content item by naming the layout file after the content item. For example, `layout/default/article/article-name.html` will be used to render an article with the name `article-name.md`. This applies to other content types as well.

### Section-Specific Layouts

In addition to global and content-type specific layouts, you can also define layouts for specific sections:

- **Section Default Layouts**: Create a layout for a specific section by placing a default layout in the section's directory. For example, `layout/section-name/default.html` will be used for all content in the `section-name` section.
- **Section-Specific Content-Type Layouts**: Similar to global content-type layouts, you can have section-specific layouts for each content type. For example, `layout/section-name/article/default.html` will be used for all articles in the `section-name` section.
- **Section-Specific Content Layouts**: You can also create layouts for specific content items within a section. For example, `layout/section-name/article/article-name.html` will be used to render an article named `article-name.md` within the `section-name` section.

### Summary

- **Global Default**: `layout/default.html`
- **Content-Type Defaults**: `layout/page/default.html`, `layout/article/default.html`, `layout/blog/default.html`, `layout/series/default.html`
- **Specific Content Layouts**: `layout/default/article/article-name.html`
- **Section Default**: `layout/section-name/default.html`
- **Section-Specific Content-Type Layouts**: `layout/section-name/article/default.html`
- **Section-Specific Content Layouts**: `layout/section-name/article/article-name.html`

This flexible layout selection logic allows you to customize the appearance of your site at various levels, from global defaults to specific content items within sections.

**This may sound overwhelming, but remember, you can rely on the defaults to start and customize later as needed. Most of the time you won't need this granular control.**

### Image Storage

Regardless of where your `.md` file is located under the `content` directory (root, section, or content type), all images for a particular post should be stored in a sibling directory with the same name as the content file. For example, if your content file is `content-name.md`, then its images should be stored in a directory called `content-name`.

### Referencing Images

You can reference images using the standard Markdown syntax. For example:

```
![An example image](index/example.png "Example Image")
```

In this example, `index.md` is the content file, and `example.png` is an image stored in the `index` directory, which is a sibling to `index.md`.

### Image Path Handling

Hermes will automatically copy the images and rewrite the paths accordingly during the build process. This ensures that your images are correctly linked in the generated HTML files.

### Header Images

For now, if you need to insert a header image for the post, you need to rely on the same method of storing and referencing images as described above. However, it is planned that Hermes will support a default header image picked using a convention-over-configuration approach. Additionally, you will be able to override this default using front matter metadata in your Markdown files.

### Example

Assume you have the following structure:

```
content/
└── guides/
    ├── article.md
    └── article/
        └── example.png
```

In `article.md`, you can reference `example.png` as follows:

```
![An example image](article/example.png "Example Image")
```

Hermes will handle copying the image and updating the path in the generated HTML.

This approach helps keep your images organized and ensures that they are easily accessible and correctly linked in your site.

## Future Plans

- Improve and polish the documentation.
- Add more detailed examples and use cases.

Stay tuned for updates and feel free to contribute or provide feedback!

---

This document is a living document and will be updated regularly as Hermes evolves.
