package skills

import (
	"fmt"
	"os"
	"path/filepath"
)

const skillName = "youve-got-marketing"

const skillContent = `---
name: youve-got-marketing
description: >
  Get brand DNA, marketing tasks, and content guidelines from You've Got Marketing.
  Use when writing user-facing copy, landing pages, docs, READMEs, error messages,
  microcopy, release notes, changelogs, or any marketing-related content.
---

# You've Got Marketing

This project uses [You've Got Marketing](https://youvegotmarketing.com) for brand
context and marketing task management.

## Getting Context

Run ` + "`ygm context`" + ` to get the full marketing context including brand voice, visual
guidelines, and actionable tasks. The output is JSON suitable for LLM consumption.

## Available Commands

- ` + "`ygm brand --json`" + ` - Brand DNA (colors, fonts, voice guidelines)
- ` + "`ygm context`" + ` - Full context dump (brand + plan + tasks)

### Task Management

- ` + "`ygm tasks --json`" + ` - List marketing tasks (filter with --status, --platform)
- ` + "`ygm tasks create --title \"...\" [--platform X] [--description \"...\"] [--date YYYY-MM-DD] --json`" + ` - Create a task
- ` + "`ygm tasks update <id> [--title \"...\"] [--description \"...\"] [--status pending|in_progress|completed] --json`" + ` - Update a task
- ` + "`ygm tasks discard <id> --json`" + ` - Discard (soft-delete) a task

## When to Use

**Always run ` + "`ygm context`" + ` before**:
- Writing user-facing copy (landing pages, docs, READMEs)
- Creating UI text, error messages, or microcopy
- Writing release notes or changelogs
- Any marketing or brand-related content

**Pattern**: Run ` + "`ygm brand --json`" + ` first to get voice guidelines, then write on-brand.

## Configuration

- Global config: ` + "`~/.config/ygm/config.yml`" + ` (auth tokens)
- Local config: ` + "`.ygm.yml`" + ` (which org to use in this project)
- Link a project: ` + "`ygm link [org-slug]`" + `
- Unlink: ` + "`ygm unlink`" + `

## Not Installed?

If ` + "`ygm`" + ` is not on PATH, install it:
` + "```" + `
curl -fsSL https://youvegotmarketing.com/install.sh | sh
` + "```" + `
`

// skillPaths returns the canonical directory and symlink directories for the
// skill under the given base directory.
func skillPaths(base string) (canonDir string, symlinkDirs []string) {
	canonDir = filepath.Join(base, ".agents", "skills", skillName)
	symlinkDirs = []string{
		filepath.Join(base, ".github", "skills", skillName),
		filepath.Join(base, ".claude", "skills", skillName),
	}
	return canonDir, symlinkDirs
}

// InstallGlobal creates the skill at ~/.agents/skills/youve-got-marketing/SKILL.md
// and symlinks from ~/.github/skills/ and ~/.claude/skills/.
func InstallGlobal() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot determine home directory: %w", err)
	}
	return install(home)
}

// InstallLocal creates the skill at .agents/skills/youve-got-marketing/SKILL.md
// and symlinks from .github/skills/ and .claude/skills/ in the current directory.
func InstallLocal() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot determine working directory: %w", err)
	}
	return install(cwd)
}

// RemoveLocal removes .agents/skills/youve-got-marketing and its symlinks
// from the current directory.
func RemoveLocal() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot determine working directory: %w", err)
	}

	canonDir, symlinkDirs := skillPaths(cwd)

	// Remove symlinks (only if they are actually symlinks)
	for _, dir := range symlinkDirs {
		fi, err := os.Lstat(dir)
		if err != nil {
			continue // doesn't exist or can't stat — skip
		}
		if fi.Mode()&os.ModeSymlink == 0 {
			fmt.Fprintf(os.Stderr, "Warning: %s is not a symlink, skipping removal\n", dir)
			continue
		}
		if err := os.Remove(dir); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not remove symlink %s: %v\n", dir, err)
		}
	}

	// Remove canonical dir
	if err := os.RemoveAll(canonDir); err != nil {
		return fmt.Errorf("failed to remove local skill: %w", err)
	}

	return nil
}

// install writes SKILL.md and creates symlinks under the given base directory.
func install(base string) error {
	canonDir, symlinkDirs := skillPaths(base)

	if err := writeSkill(canonDir); err != nil {
		return fmt.Errorf("failed to write skill: %w", err)
	}

	for _, dir := range symlinkDirs {
		if err := ensureSymlink(canonDir, dir); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not create symlink %s: %v\n", dir, err)
		}
	}

	return nil
}

// writeSkill writes SKILL.md into the given directory, creating it if needed.
func writeSkill(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(skillContent), 0644)
}

// ensureSymlink creates a symlink at linkPath pointing to target.
// If linkPath already exists and points to target, it's a no-op.
// If it's a symlink pointing elsewhere, it's replaced.
// If it's a real file or directory, it's left alone (returns error).
func ensureSymlink(target, linkPath string) error {
	if err := os.MkdirAll(filepath.Dir(linkPath), 0755); err != nil {
		return err
	}

	fi, err := os.Lstat(linkPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("cannot check %s: %w", linkPath, err)
		}
		// Nothing exists, create the symlink
		return os.Symlink(target, linkPath)
	}

	// Something exists — only replace if it's a symlink
	if fi.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("%s already exists and is not a symlink", linkPath)
	}

	existing, err := os.Readlink(linkPath)
	if err == nil && existing == target {
		return nil // Already correct
	}

	// Symlink points elsewhere — replace it
	if err := os.Remove(linkPath); err != nil {
		return err
	}
	return os.Symlink(target, linkPath)
}
