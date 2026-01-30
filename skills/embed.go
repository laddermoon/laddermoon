package skills

import "embed"

//go:embed laddermoon-feed/SKILL.md
//go:embed laddermoon-sync/SKILL.md
//go:embed laddermoon-audit/SKILL.md
//go:embed laddermoon-propose/SKILL.md
//go:embed laddermoon-review/SKILL.md
//go:embed laddermoon-criticize/SKILL.md
//go:embed laddermoon-clarify/SKILL.md
//go:embed laddermoon-code/SKILL.md
//go:embed laddermoon-apply/SKILL.md
var SkillsFS embed.FS

// SkillNames lists all available LadderMoon skills
var SkillNames = []string{
	"laddermoon-feed",
	"laddermoon-sync",
	"laddermoon-audit",
	"laddermoon-propose",
	"laddermoon-review",
	"laddermoon-criticize",
	"laddermoon-clarify",
	"laddermoon-code",
	"laddermoon-apply",
}
