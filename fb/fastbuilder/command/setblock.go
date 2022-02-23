package command

import (
	"fmt"
	"phoenixbuilder_3rd_gui/fb/fastbuilder/types"
)

func SetBlockRequest(module *types.Module, config *types.MainConfig) string {
	Block := module.Block
	Point := module.Point
	Method := config.Method
	if Block != nil {
		return fmt.Sprintf("setblock %v %v %v %v %v %v", Point.X, Point.Y, Point.Z, *Block.Name, Block.Data, Method)
	} else {
		return fmt.Sprintf("setblock %v %v %v %v %v %v", Point.X, Point.Y, Point.Z, config.Block.Name, config.Block.Data, Method)
	}
}
