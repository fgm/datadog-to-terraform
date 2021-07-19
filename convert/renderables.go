package convert

import (
	"fmt"
	"sort"
	"strings"
)

// Renderable implementations transforms a value a Terraform HCL definition,
// which can be either a property assignement or a - possibly repeated - block
// definition, e.g.:
//
//   - title: "some title"
//   - widget { ... }
type Renderable interface {
	Render(level int, name string) (string, error)
}

type boolValue bool

func (s boolValue) Render(level int, name string) (string, error) {
	return fmt.Sprintf("%s%s = %t", indent(level), name, s), nil
}

type hiddenValue string

func (s hiddenValue) Render(level int, name string) (string, error) {
	return fmt.Sprintf(`%s// %s is hidden"`, indent(level), name), nil
}

type intValue int

func (iv intValue) Render(level int, name string) (string, error) {
	return fmt.Sprintf("%s%s = %d,", indent(level), name, iv), nil
}


type intListValue []int

func (il intListValue) Render(level int, name string) (string, error) {
	build := []string{
		fmt.Sprintf("%s%s = [", indent(level), name),
	}
	for _, s := range il {
		build = append(build, fmt.Sprintf(`%s%d,`, indent(level+1), s))
	}
	build = append(build, indent(level) + "]")
	return strings.Join(build, "\n"), nil
}


type stringValue string

func (s stringValue) Render(level int, name string) (string, error) {
	return fmt.Sprintf(`%s%s = "%s"`, indent(level), name, s), nil
}

type stringListValue []string

func (sl stringListValue) Render(level int, name string) (string, error) {
	build := []string{
		fmt.Sprintf("%s%s = [", indent(level), name),
	}
	if len(sl) == 0 {
		build = append(build, "]")
		return strings.Join(build, ""), nil
	}
	for _, s := range sl {
		build = append(build, fmt.Sprintf(`%s"%s",`, indent(level+1), s))
	}
	build = append(build, indent(level) + "]")
	return strings.Join(build, "\n"), nil
}

func (b *TFBlock) buildProps(level int) ([]string, error) {
	build := make([]string, 0, len(b.props))
	keys := make([]string, 0, len(b.props))
	for k := range b.props {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		prop := b.props[k]
		row, err := prop.Render(level, k)
		if err != nil {
			return nil, fmt.Errorf("rendering %s in %s %s: %w", k, b.RootTag, b.JSONTag, err)
		}
		build = append(build, row)
	}
	build = append(build, "}")
	return build, nil
}

func (b *TFBlock) buildBlocks(level int) ([]string, error) {
	return nil, nil
}

func (b *TFBlock) Render(level int, name string) (string, error) {
	build := make([]string, 0, 1 + len(b.props) + len(b.blockLists) + 1)
	build = append(build, fmt.Sprintf("%s%s {", indent(level), name))
	props, err := b.buildProps(level+1)
	if err != nil {
		return "", fmt.Errorf("rendering %s props: %w", name, err)
	}
	build = append(build, props...)

	blocks, err := b.buildBlocks(level+1)
	if err != nil {
		return "", fmt.Errorf("rendering %s blocks: %w", name, err)
	}
	build = append(build, blocks...)

	return strings.Join(build, "\n") + "\n", nil
}

func (bl TFBlockList) Render(level int, name string) (string, error) {
	build := []string{
		"\n",
	}
	for i, b := range bl {
		row, err := b.Render(level, name)
		if err != nil {
			return "", fmt.Errorf("failed rendering block list entry %d under key %s in %s context %s: %w",
				i, name, b.RootTag, b.JSONTag, err)
		}
		build = append(build, row)
	}
	out := strings.Join(build, "\n")
	return out, nil
}
