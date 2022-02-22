package enchantment

import (
	"phoenixbuilder/dragonfly/server/item"
	"phoenixbuilder/dragonfly/server/item/tool"
)

// SilkTouch is an enchantment that allows many blocks to drop themselves instead of their usual items when mined.
type SilkTouch struct{ enchantment }

// Name ...
func (e SilkTouch) Name() string {
	return "Silk Touch"
}

// MaxLevel ...
func (e SilkTouch) MaxLevel() int {
	return 1
}

// WithLevel ...
func (e SilkTouch) WithLevel(level int) item.Enchantment {
	return SilkTouch{e.withLevel(level, e)}
}

// CompatibleWith ...
func (e SilkTouch) CompatibleWith(s item.Stack) bool {
	t, ok := s.Item().(tool.Tool)
	//TODO: Fortune
	return ok && (t.ToolType() != tool.TypeSword && t.ToolType() != tool.TypeNone)
}
