package repo

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Skill struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Path        string
}

func GetSkillDetails(skillPath string) (Skill, error) {
	skillFile := filepath.Join(skillPath, "SKILL.md")
	f, err := os.Open(skillFile)
	if err != nil {
		return Skill{}, err
	}
	defer f.Close()

	var yamlLines []string
	scanner := bufio.NewScanner(f)
	inFrontmatter := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			} else {
				break
			}
		}
		if inFrontmatter {
			yamlLines = append(yamlLines, line)
		}
	}

	var skill Skill
	err = yaml.Unmarshal([]byte(strings.Join(yamlLines, "\n")), &skill)
	if err != nil {
		return Skill{}, err
	}

	skill.Path = skillPath
	if skill.Name == "" {
		skill.Name = filepath.Base(skillPath)
	}

	return skill, nil
}

func ListSkills(repoPath string) ([]Skill, error) {
	entries, err := os.ReadDir(repoPath)
	if err != nil {
		return nil, err
	}

	var skills []Skill
	for _, entry := range entries {
		if entry.IsDir() {
			skillPath := filepath.Join(repoPath, entry.Name())
			if _, err := os.Stat(filepath.Join(skillPath, "SKILL.md")); err == nil {
				skill, err := GetSkillDetails(skillPath)
				if err == nil {
					skills = append(skills, skill)
				}
			}
		}
	}
	return skills, nil
}
