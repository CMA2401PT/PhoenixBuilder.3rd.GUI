package trace

import (
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/block/cube"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/entity/physics"
	"phoenixbuilder_3rd_gui/fb/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"math"
)

// Result represents the result of a ray trace collision with a bounding box.
type Result interface {
	// AABB returns the bounding box collided with.
	AABB() physics.AABB
	// Position returns where the ray first collided with the bounding box.
	Position() mgl64.Vec3
	// Face returns the face of the bounding box that was collided on.
	Face() cube.Face
}

// Perform performs a ray trace between start and end, checking if any blocks or entities collided with the
// ray. The physics.AABB that's passed is used for checking if any entity within the bounding box collided
// with the ray.
func Perform(start, end mgl64.Vec3, w *world.World, aabb physics.AABB, ignoredEntities ...world.Entity) (hit Result, ok bool) {
	// Check if there's any blocks that we may collide with.
	TraverseBlocks(start, end, func(blockPos cube.Pos) (cont bool) {
		dist := math.MaxFloat64
		// Check if there's any entities that we may collide with.
		for _, entity := range w.CollidingEntities(aabb.Translate(blockPos.Vec3Centre()), ignoredEntities...) {
			// Check if we collide with the entities bounding box.
			result, ok := EntityIntercept(entity, start, end)
			if !ok {
				continue
			}

			if distance := start.Sub(result.Position()).LenSqr(); distance < dist {
				dist = distance
				hit = result
			}
		}

		if hit != nil {
			return false
		}

		b := w.Block(blockPos)

		// Check if we collide with the block's model.
		if result, ok := BlockIntercept(blockPos, w, b, start, end); ok {
			hit = result
			return false
		}
		return true
	})

	if hit == nil {
		return nil, false
	}

	return hit, true
}
