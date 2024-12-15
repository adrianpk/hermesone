package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrianpk/gohermes/internal/hermes"
)

// GenHTML generates the HTML files from the markdown files.
func GenHTML() error {
	pp, err := startPreProcessor(hermes.ContentDir)
	if err != nil {
		return err
	}

	err = filepath.Walk(hermes.ContentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".md" {
			relativePath, err := filepath.Rel(hermes.ContentDir, path)
			if err != nil {
				log.Printf("error getting relative path for %s: %v", path, err)
				return nil
			}

			fileData, exists := pp.FindFileData(relativePath)
			if !exists {
				log.Printf("file data not found for %s", relativePath)
				return nil
			}

			if !fileData.Published {
				return nil
			}

			outputPath := determineOutputPath(relativePath)

			if shouldRender(path, outputPath) {
				fileContent, err := os.ReadFile(path)
				if err != nil {
					log.Printf("error reading file %s: %v", path, err)
					return nil
				}

				content, err := hermes.Parse(fileContent, path)
				if err != nil {
					log.Printf("error parsing file %s: %v", path, err)
					return nil
				}

				layoutPath := findLayout(path)
				if layoutPath != "" {
					tmpl, err := template.New("webpage").Funcs(template.FuncMap{
						"safeHTML": safeHTML,
					}).ParseFiles(layoutPath)

					if err != nil {
						log.Printf("error parsing template files for %s: %v", path, err)
						return nil
					}

					var tmplBuf bytes.Buffer

					err = tmpl.Execute(&tmplBuf, content)
					if err != nil {
						log.Printf("error executing template for %s: %v", path, err)
						return nil
					}

					err = os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
					if err != nil {
						log.Printf("error creating directories for %s: %v", outputPath, err)
						return nil
					}

					outputFile, err := os.Create(outputPath)
					if err != nil {
						log.Printf("error creating output file %s: %v", outputPath, err)
						return nil
					}
					defer outputFile.Close()

					_, err = tmplBuf.WriteTo(outputFile)
					if err != nil {
						log.Printf("error writing to output file %s: %v", outputPath, err)
						return nil
					}

					err = copyImages(path, outputPath)
					if err != nil {
						log.Printf("error copying images for %s: %v", path, err)
						return nil
					}

				} else {
					log.Printf("no layout found for %s", path)
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	err = processRootSection(hermes.ContentDir, hermes.OutputDir, hermes.LayoutDir, pp)
	if err != nil {
		return fmt.Errorf("error processing root section: %w", err)
	}

	cfg, err := hermes.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	sections := cfg.Sections

	err = processSections(hermes.ContentDir, hermes.OutputDir, hermes.LayoutDir, pp, sections)
	if err != nil {
		return fmt.Errorf("error processing sections: %w", err)
	}

	err = addNoJekyll()
	if err != nil {
		return fmt.Errorf("error adding .nojekyll file: %w", err)
	}

	log.Println("content generated!")
	return nil
}

// genDefaultIndex generates the default index.html file.
// If no index.md is provided for the section then this one is ussed as default. It renders a list of all published
// content for this section.
// This is a WIP, a lot of logging is present to help debug the process. It will be removed as soon the loggic is
// stable. Also, there are a lot of hardcoded values that will be replaced by dinamic ones.
// Finally, this is generating the default index for the root section, it includes references to all the content in the
// site.
// We will also need a similar logic to generate the section index when content is not provided for it.
// This will show all the content for the specific section.
// Worth mentioning that the partial used to render the content should also be improved to show a nice presentation of
// the content (image, title, excerpt, etc).
func genDefaultIndex(pp *hermes.PreProcessor, rootIndexPath, outputPath string, fd []hermes.FileData) error {
	partial := "layout/default/partials/_index.html"
	log.Printf("using partial template: %s\n", partial)

	partialTmpl, err := template.New("_index.html").ParseFiles(partial)
	if err != nil {
		return err
	}

	var partialBuf bytes.Buffer

	err = partialTmpl.Execute(&partialBuf, fd)
	if err != nil {
		return err
	}

	content := map[string]interface{}{
		"HTML": partialBuf.String(),
	}

	layoutPath := findLayout(rootIndexPath)
	if layoutPath == "" {
		return fmt.Errorf("no layout found for %s", rootIndexPath)
	}

	layoutTmpl, err := template.New("webpage").Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	}).ParseFiles(layoutPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	var finalBuf bytes.Buffer

	err = layoutTmpl.Execute(&finalBuf, content)
	if err != nil {
		return err
	}

	_, err = finalBuf.WriteTo(outputFile)
	if err != nil {
		log.Printf("error writing to output file: %v\n", err)
		return err
	}

	return nil
}

func processRootSection(contentDir, outputDir, layoutDir string, pp *hermes.PreProcessor) error {
	log.Println("processing root section")

	indexPath := filepath.Join(hermes.ContentDir, "root", "index.md")
	outputPath := filepath.Join(hermes.OutputDir, "index.html")

	if !isValidIndex(indexPath, pp) {
		err := genDefaultIndex(pp, indexPath, outputPath, pp.GetAllPublished())
		if err != nil {
			return err
		}
	}

	return nil
}

func processSections(contentDir, outputDir, layoutDir string, pp *hermes.PreProcessor, sections []hermes.Section) error {
	for _, section := range sections {
		if section.Name == "root" {
			continue
		}

		indexPath := filepath.Join(hermes.ContentDir, section.Name, hermes.IndexMdFile)
		outputPath := filepath.Join(hermes.OutputDir, section.Name, hermes.IndexFile)
		layoutPath := filepath.Join(layoutDir, "index.html")

		log.Printf("layout path: %s\n", layoutPath)

		if !isValidIndex(indexPath, pp) {
			log.Printf("section index is not valid for section: %s, generating fallback index", section.Name)

			err := genDefaultIndex(pp, indexPath, outputPath, pp.GetPublishedBySection(section.Name))
			if err != nil {
				log.Printf("error generating fallback index for section: %s, error: %v\n", section.Name, err)
				continue
			}
		} else {
			log.Printf("section index is valid for section: %s, no need to generate fallback index", section.Name)
		}
	}

	log.Println("finished processSections")
	return nil
}

func isValidIndex(indexPath string, pp *hermes.PreProcessor) bool {
	relPath := strings.TrimPrefix(indexPath, "content/")

	fileData, _ := pp.FindFileData(relPath)

	return fileData.IsIndex()
}

func determineOutputPath(relativePath string) string {
	parts := strings.Split(relativePath, string(os.PathSeparator))

	if len(parts) < 2 {
		return outputPath(parts...)
	}

	section := parts[0]
	subdir := parts[1]
	needDir := needsCustomDir(subdir)

	switch section {
	case hermes.DefSection:
		if needDir {
			p := outputPath(append([]string{subdir}, parts[2:]...)...)
			return p
		} else {
			p := outputPath(append([]string{}, parts[2:]...)...)
			return p
		}

	case ct.Page, ct.Article:
		p := outputPath(parts[1:]...)
		return p

	default:
		if needDir {
			p := outputPath(append([]string{section, subdir}, parts[2:]...)...)
			return p
		} else {
			p := outputPath(append([]string{section}, parts[2:]...)...)
			return p
		}
	}
}

func needsCustomDir(dir string) bool {
	return dir != ct.Article && dir != ct.Page
}

func outputPath(parts ...string) string {
	trimmedPath := strings.TrimSuffix(strings.Join(parts, string(os.PathSeparator)), filepath.Ext(parts[len(parts)-1])) + ".html"
	return filepath.Join(hermes.OutputDir, trimmedPath)
}

// shouldRender checks if the markdown file is newer than the html file
// to determine if the html should be re-rendered.
func shouldRender(mdPath, htmlPath string) bool {
	htmlInfo, err := os.Stat(htmlPath)
	if os.IsNotExist(err) {
		return true
	}

	markdownInfo, err := os.Stat(mdPath)
	if err != nil {
		return false
	}

	return markdownInfo.ModTime().After(htmlInfo.ModTime())
}

// findLayout tries to find the layout file for the given markdown file.
func findLayout(path string) string {
	fmt.Println("")
	fmt.Println("=== findLayout start ===")
	defer func() {
		fmt.Println("=== findLayout end ===")
		fmt.Println("")
	}()

	fmt.Printf("Input path: %s\n", path)
	path = strings.TrimPrefix(path, "content/")
	fmt.Printf("Trimmed path: %s\n", path)

	base := filepath.Base(path)
	fmt.Printf("Base: %s\n", base)
	base = strings.TrimSuffix(base, filepath.Ext(base))
	fmt.Printf("Base without extension: %s\n", base)

	dir := filepath.Dir(path)
	fmt.Printf("Directory: %s\n", dir)
	section, _ := sectionTypeSegments(path)
	fmt.Printf("Section: %s\n", section)

	secTypeLayoutDir := filepath.Join(hermes.LayoutDir, dir)
	fmt.Printf("Section type layout directory: %s\n", secTypeLayoutDir)

	secLayoutDir := filepath.Join(hermes.LayoutDir, section)
	fmt.Printf("Section layout directory: %s\n", secLayoutDir)

	layoutPaths := []string{
		filepath.Join(secTypeLayoutDir, base+".html"),
		filepath.Join(secTypeLayoutDir, hermes.DefLayout),
		filepath.Join(secLayoutDir, hermes.DefLayout),
		filepath.Join(hermes.DefLayoutPath, filepath.Base(dir), hermes.DefLayout),
		filepath.Join(hermes.DefLayoutPath, base+".html"),
		filepath.Join(hermes.DefLayoutPath, hermes.DefLayout),
	}
	fmt.Printf("Layout paths: %v\n", layoutPaths)

	for _, layoutPath := range layoutPaths {
		fmt.Printf("Checking layout path: %s\n", layoutPath)
		if _, err := os.Stat(layoutPath); err == nil {
			fmt.Printf("Found layout path: %s\n", layoutPath)
			return layoutPath
		}
	}

	fmt.Println("No layout path found")
	return ""
}

// findLayout tries to find the layout file for the given markdown file.
func findLayout2(path string) string {
	path = strings.TrimPrefix(path, "content/")

	base := filepath.Base(path)
	base = strings.TrimSuffix(base, filepath.Ext(base))
	dir := filepath.Dir(path)
	section, _ := sectionTypeSegments(path)

	secTypeLayoutDir := filepath.Join(hermes.LayoutDir, dir)
	secLayoutDir := filepath.Join(hermes.LayoutDir, section)

	layoutPaths := []string{
		filepath.Join(secTypeLayoutDir, base+".html"),
		filepath.Join(secTypeLayoutDir, hermes.DefLayout),
		filepath.Join(secLayoutDir, hermes.DefLayout),
		filepath.Join(hermes.DefLayoutPath, filepath.Base(dir), hermes.DefLayout),
		filepath.Join(hermes.DefLayoutPath, base+".html"),
		filepath.Join(hermes.DefLayoutPath, hermes.DefLayout),
	}

	for _, layoutPath := range layoutPaths {
		if _, err := os.Stat(layoutPath); err == nil {
			return layoutPath
		}
	}

	return ""
}

func sectionTypeSegments(path string) (string, string) {
	dir := filepath.Dir(path)
	segments := strings.Split(dir, osFileSep)

	var sectionSegment, typeSegment string
	if len(segments) > 0 {
		sectionSegment = segments[0]
	}
	if len(segments) > 1 {
		typeSegment = segments[1]
	}

	return sectionSegment, typeSegment
}

// copyImages copies the images from the markdown directory to the output directory.
func copyImages(mdPath, htmlPath string) error {
	rootPrefix := "root/"
	imageDir := strings.TrimSuffix(mdPath, filepath.Ext(mdPath))
	relativeImageDir := strings.TrimPrefix(imageDir, hermes.ContentDir+"/")

	if strings.HasPrefix(relativeImageDir, rootPrefix) {
		relativeImageDir = strings.TrimPrefix(relativeImageDir, rootPrefix)

		parts := strings.Split(relativeImageDir, string(osFileSep))

		if len(parts) == 1 {
			relativeImageDir = filepath.Join(hermes.ImgDir, parts[0])
		} else if len(parts) > 1 {
			switch parts[0] {
			case ct.Blog, ct.Series:
				relativeImageDir = filepath.Join(hermes.ImgDir, parts[0], parts[1])
			case ct.Article, ct.Page:
				relativeImageDir = filepath.Join(hermes.ImgDir, parts[1])
			default:
				relativeImageDir = filepath.Join(hermes.ImgDir, strings.Join(parts, string(osFileSep)))
			}
		}
	} else {
		parts := strings.Split(relativeImageDir, string(osFileSep))

		if len(parts) > 2 && (parts[1] == ct.Article || parts[1] == ct.Page) {
			relativeImageDir = filepath.Join(hermes.ImgDir, parts[0], parts[2])
		} else if len(parts) > 2 && (parts[1] == ct.Blog || parts[1] == ct.Series) {
			relativeImageDir = filepath.Join(hermes.ImgDir, parts[0], parts[1], parts[2])
		} else {
			relativeImageDir = filepath.Join(hermes.ImgDir, strings.Join(parts, string(os.PathSeparator)))
		}
	}

	outputImageDir := filepath.Join(hermes.OutputDir, relativeImageDir)

	err := filepath.Walk(imageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relativePath, err := filepath.Rel(imageDir, path)
			if err != nil {
				return err
			}

			destPath := filepath.Join(outputImageDir, relativePath)

			err = os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
			if err != nil {
				return err
			}

			srcFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			destFile, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, srcFile)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// addNoJekyll adds a .nojekyll file to the output directory.
func addNoJekyll() error {
	noJekyllPath := filepath.Join(hermes.OutputDir, hermes.NoJekyllFile)
	if _, err := os.Stat(noJekyllPath); os.IsNotExist(err) {
		file, err := os.Create(noJekyllPath)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

func startPreProcessor(root string) (*hermes.PreProcessor, error) {
	pp := hermes.NewPreProcessor(root)
	err := pp.Build()
	if err != nil {
		log.Printf("error building pp: %v", err)
		return nil, err
	}

	err = pp.Sync()
	if err != nil {
		log.Printf("error syncing pp: %v", err)
		return nil, err
	}

	//pp.Debug()

	return pp, nil
}

// safeHTML function to replace the anonymous function
func safeHTML(s string) template.HTML {
	return template.HTML(s)
}
