package strategies

import (
	"github.com/jimschubert/ignore/internal/util"
	"github.com/jimschubert/ignore/parser"
	"github.com/jimschubert/ignore/rules"
	"github.com/jimschubert/ignore/strategy"
)

func DefaultStrategy() strategy.Strategy {
	return GitignoreStrategy()
}

func GitignoreStrategy() strategy.Strategy {
	p := parser.NewGitignoreParser()
	return simpleStrategy{
		fullPath: ".gitignore",
		parser:   p,
		ruleBuilder: strategy.BuildRulesFrom(func(parts []parser.TokenValue) (rules.Rule, error) {
			if len(parts) == 1 && parts[0] == parser.NewLine {
				return rules.NewEmptyRule(util.StringValue(parts[0].Line), parts)
			}

			switch len(parts) {
			case 0:
				return rules.NewEmptyRule(util.StringValue(parts[0].Line), parts)
			case 1:
				part := parts[0]
				if part.Token == parser.MatchAny {
					return rules.NewRootedFileRule(util.StringValue(parts[0].Line), parts)
				}

				return rules.NewFileRule(util.StringValue(parts[0].Line), parts)
			}

			head := parts[0]
			tail := parts[len(parts)-1]

			if tail.Token == parser.DirectoryMarker {
				return rules.NewDirectoryRule(util.StringValue(parts[0].Line), parts)
			}

			if head.Token == parser.RootedMarker {
				return rules.NewRootedFileRule(util.StringValue(parts[0].Line), parts)
			}

			return rules.NewFileRule(util.StringValue(parts[0].Line), parts)
		}),
	}
}
