package skills

import "embed"

//go:embed laddermoon-feed/SKILL.md
//go:embed laddermoon-sync/SKILL.md
//go:embed laddermoon-audit/SKILL.md
//go:embed laddermoon-propose/SKILL.md
//go:embed laddermoon-solve/SKILL.md
var SkillsFS embed.FS

// SkillNames lists all available LadderMoon skills
var SkillNames = []string{
	"laddermoon-feed",
	"laddermoon-sync",
	"laddermoon-audit",
	"laddermoon-propose",
	"laddermoon-solve",
}
