package ignore

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/jimschubert/ignore/internal/strategies"
	"github.com/jimschubert/ignore/parser"
	"github.com/jimschubert/ignore/rules"
	"github.com/jimschubert/ignore/strategy"
)

type Processor struct {
	strategy    strategy.Strategy
	ruleList    []rules.Rule
	initialized bool
}

func (p *Processor) processIgnoreFile() error {
	if p.initialized {
		return nil
	}

	p.initialized = true
	ignoreFile := p.strategy.DefinitionPath()
	file, err := os.Open(ignoreFile)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	parts, err := p.strategy.Parser().ParseAll(file)
	if err != nil {
		return err
	}

	maxLen := len(parts)
	for i := 0; i < maxLen; i++ {
		if parts[i].Token == parser.LineFeed {
			continue
		}

		width := 0
		for _, value := range parts[i:] {
			if value.Token == parser.LineFeed {
				break
			}
			width += 1
		}

		if width == 0 {
			continue
		}

		var rule rules.Rule

		// TODO: Decide if definition is really need here
		if rule, err = p.strategy.RuleBuilder().RuleFor(parts[i : i+width]); err != nil {
			return err
		}

		p.ruleList = append(p.ruleList, rule)

		i += width
	}

	return nil
}

func (p *Processor) AllowsFile(path string) (bool, error) {
	if err := p.processIgnoreFile(); err != nil {
		return true, err
	}

	exclude := false
	hasIncludes := false
	for _, rule := range p.ruleList {
		switch r := rule.(type) {
		case rules.EvaluatingRule:
			op, err := r.Evaluate(path)
			if err != nil {
				return false, err
			}

			if op == rules.ExcludeAndTerminate {
				return false, err
			}

			// invalid rules will not impact include/exclude analysis.
			if op != rules.Invalid {
				exclude = exclude || op == rules.Exclude
				hasIncludes = hasIncludes || op == rules.Include
			}
		}
	}

	// regardless of any combination… if user explicitly includes anywhere it will include.
	if exclude && hasIncludes {
		return true, nil
	}

	return !exclude, nil
}

type ProcessorOption func(*Processor) error

// WithGitignoreStrategy is a functional option which applies the strategy for parsing .gitignore files
func WithGitignoreStrategy() ProcessorOption {
	return func(processor *Processor) error {
		processor.strategy = strategies.GitignoreStrategy()
		return nil
	}
}

// WithIgnoreFilePath is a functional option which allows the user to parse a non-standard filename for a given strategy
func WithIgnoreFilePath(filePath string) ProcessorOption {
	return func(processor *Processor) error {
		m, err := strategies.AsMutable(processor.strategy)
		if err != nil {
			return err
		}

		var fullPath string
		fullPath, _ = filepath.Abs(filePath)
		fileInfo, err := os.Stat(fullPath)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			return errors.New("ignore file path must be a regular file, not a directory")
		}

		if err := m.SetDefinitionPath(fullPath); err != nil {
			return err
		}

		processor.strategy = m
		return nil
	}
}

func NewProcessor(opts ...ProcessorOption) (*Processor, error) {
	// TODO: supporting other strategies would mean inferring strategy from ignore filenames
	ruleList := make([]rules.Rule, 0)
	processor := &Processor{
		strategy: strategies.DefaultStrategy(),
		ruleList: ruleList,
	}

	for _, opt := range opts {
		if err := opt(processor); err != nil {
			return nil, err
		}
	}

	return processor, nil
}
