package internal

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CreateBackup(skills []SkillEntry, targetPath string) error {
	file, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("create archive: %w", err)
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	manifest := Manifest{
		Version:   ToolVersion,
		CreatedAt: time.Now(),
		Skills:    skills,
		Total:     len(skills),
	}

	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal manifest: %w", err)
	}

	if err := writeTarEntry(tw, "metadata.json", manifestData); err != nil {
		return err
	}

	for _, skill := range skills {
		src := filepath.Join(skill.Path, "SKILL.md")
		data, err := os.ReadFile(src)
		if err != nil {
			continue
		}

		archivePath := "skills/" + skill.Name + "/SKILL.md"
		if err := writeTarEntry(tw, archivePath, data); err != nil {
			return err
		}
	}

	return nil
}

func RestoreBackup(archivePath string, targetBase string) ([]string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return nil, fmt.Errorf("open archive: %w", err)
	}
	defer file.Close()

	gr, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("read gzip: %w", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	var restored []string

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read tar entry: %w", err)
		}

		if header.Typeflag != tar.TypeReg {
			continue
		}

		if !strings.HasPrefix(header.Name, "skills/") {
			continue
		}

		parts := strings.SplitN(strings.TrimPrefix(header.Name, "skills/"), "/", 2)
		if len(parts) < 2 {
			continue
		}

		skillName := parts[0]
		destDir := filepath.Join(targetBase, skillName)
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return nil, fmt.Errorf("create dir %s: %w", destDir, err)
		}

		destFile := filepath.Join(destDir, "SKILL.md")
		out, err := os.Create(destFile)
		if err != nil {
			return nil, fmt.Errorf("create file %s: %w", destFile, err)
		}

		if _, err := io.Copy(out, tr); err != nil {
			out.Close()
			return nil, fmt.Errorf("write %s: %w", destFile, err)
		}
		out.Close()

		restored = append(restored, skillName)
	}

	return restored, nil
}

func writeTarEntry(tw *tar.Writer, name string, body []byte) error {
	header := &tar.Header{
		Name: name,
		Size: int64(len(body)),
		Mode: 0644,
	}

	if err := tw.WriteHeader(header); err != nil {
		return fmt.Errorf("write tar header %s: %w", name, err)
	}

	if _, err := tw.Write(body); err != nil {
		return fmt.Errorf("write tar body %s: %w", name, err)
	}

	return nil
}
