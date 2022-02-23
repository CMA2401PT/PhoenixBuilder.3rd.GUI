package trace

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/cube"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/entity/physics"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"math"
)

// BlockResult is the result of a ray trace collision with a block's model.
type BlockResult struct {
	bb   physics.AABB
	pos  mgl64.Vec3
	face cube.Face

	blockPos cube.Pos
}

// AABB returns the AABB that was collided within the block's model.
func (r BlockResult) AABB() physics.AABB {
	return r.bb
}

// Position ...
func (r BlockResult) Position() mgl64.Vec3 {
	return r.pos
}

// Face returns the hit block face.
func (r BlockResult) Face() cube.Face {
	return r.face
}

// BlockPosition returns the block that was collided with.
func (r BlockResult) BlockPosition() cube.Pos {
	return r.blockPos
}

// BlockIntercept performs a ray trace and calculates the point on the block model's edge nearest to the start position
// that the ray collided with.
// BlockIntercept returns a BlockResult with the block collided with and with the colliding vector closest to the start position,
// if no colliding point was found, a zero BlockResult is returned and ok is false.
func BlockIntercept(pos cube.Pos, w *world.World, b world.Block, start, end mgl64.Vec3) (result BlockResult, ok bool) {
	bbs := b.Model().AABB(pos, w)
	if len(bbs) == 0 {
		return
	}

	var (
		hit  Result
		dist = math.MaxFloat64
	)

	for _, bb := range bbs {
		next, ok := AABBIntercept(bb.Translate(pos.Vec3()), start, end)
		if !ok {
			continue
		}

		nextDist := next.Position().Sub(start).LenSqr()
		if nextDist < dist {
			hit = next
			dist = nextDist
		}
	}

	if hit == nil {
		return result, false
	}

	return BlockResult{bb: hit.AABB(), pos: hit.Position(), face: hit.Face(), blockPos: pos}, true
}
