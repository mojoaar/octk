# octk — OpenCode ToolKit

Backup and restore your OpenCode skills, rules, agents, and config.

## Install

### From source

```bash
go install github.com/mojoaar/octk@latest
```

### Download binary

Grab the latest binary for your platform from the [releases page](https://github.com/mojoaar/octk/releases).

## Usage

```
octk help              Show help and logo
octk version           Show version information
octk list              List all discoverable skills
octk backup            Backup all skills to octk-skills-YYYY-MM-DD.tar.gz
octk restore <file>    Restore skills from a tar.gz archive
```

### Backup

Scans OpenCode skill discovery paths (project `.agents/skills/`, global `~/.agents/skills/`, and superpowers packages) and packages everything into a dated archive.

```bash
octk backup
# → octk-skills-2026-06-28.tar.gz
```

### Restore

Extracts all skills to `~/.agents/skills/` so OpenCode discovers them on the next session.

```bash
octk restore octk-skills-2026-06-28.tar.gz
```

## License

MIT

## Support

[Buy me a coffee](https://ko-fi.com/mojoaar)
