package block

import (
	"fmt"
	"phoenixbuilder/dragonfly/server/block/cube"
	"phoenixbuilder/dragonfly/server/internal/nbtconv"
	"phoenixbuilder/dragonfly/server/item"
	"phoenixbuilder/dragonfly/server/item/inventory"
	"phoenixbuilder/dragonfly/server/world"
	"phoenixbuilder/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"strings"
	"sync"
)

// Barrel is a fisherman's job site block, used to store items. It functions like a single chest although
// it requires no air space above it to be opened.
type Barrel struct {
	solid
	bass

	// Facing is the direction that the barrel is facing.
	Facing cube.Face
	// Open is whether the barrel is open or not.
	Open bool
	// CustomName is the custom name of the barrel. This name is displayed when the barrel is opened, and may
	// include colour codes.
	CustomName string

	inventory *inventory.Inventory
	viewerMu  *sync.RWMutex
	viewers   map[ContainerViewer]struct{}
}

// NewBarrel creates a new initialised barrel. The inventory is properly initialised.
func NewBarrel() Barrel {
	m := new(sync.RWMutex)
	v := make(map[ContainerViewer]struct{}, 1)
	return Barrel{
		inventory: inventory.New(27, func(slot int, item item.Stack) {
			m.RLock()
			defer m.RUnlock()
			for viewer := range v {
				viewer.ViewSlotChange(slot, item)
			}
		}),
		viewerMu: m,
		viewers:  v,
	}
}

// Inventory returns the inventory of the barrel. The size of the inventory will be 27.
func (b Barrel) Inventory() *inventory.Inventory {
	return b.inventory
}

// WithName returns the barrel after applying a specific name to the block.
func (b Barrel) WithName(a ...interface{}) world.Item {
	b.CustomName = strings.TrimSuffix(fmt.Sprintln(a...), "\n")
	return b
}

// open opens the barrel, displaying the animation and playing a sound.
func (b Barrel) open(w *world.World, pos cube.Pos) {
	b.Open = true
	w.PlaySound(pos.Vec3Centre(), sound.BarrelOpen{})
	w.SetBlock(pos, b)
}

// close closes the barrel, displaying the animation and playing a sound.
func (b Barrel) close(w *world.World, pos cube.Pos) {
	b.Open = false
	w.PlaySound(pos.Vec3Centre(), sound.BarrelClose{})
	w.SetBlock(pos, b)
}

// AddViewer adds a viewer to the barrel, so that it is updated whenever the inventory of the barrel is changed.
func (b Barrel) AddViewer(v ContainerViewer, w *world.World, pos cube.Pos) {
	b.viewerMu.Lock()
	defer b.viewerMu.Unlock()
	if len(b.viewers) == 0 {
		b.open(w, pos)
	}
	b.viewers[v] = struct{}{}
}

// RemoveViewer removes a viewer from the barrel, so that slot updates in the inventory are no longer sent to
// it.
func (b Barrel) RemoveViewer(v ContainerViewer, w *world.World, pos cube.Pos) {
	b.viewerMu.Lock()
	defer b.viewerMu.Unlock()
	if len(b.viewers) == 0 {
		return
	}
	delete(b.viewers, v)
	if len(b.viewers) == 0 {
		b.close(w, pos)
	}
}

// Activate ...
func (b Barrel) Activate(pos cube.Pos, _ cube.Face, _ *world.World, u item.User) {
	if opener, ok := u.(ContainerOpener); ok {
		opener.OpenBlockContainer(pos)
	}
}

// UseOnBlock ...
func (b Barrel) UseOnBlock(pos cube.Pos, face cube.Face, clickPos mgl64.Vec3, w *world.World, user item.User, ctx *item.UseContext) (used bool) {
	pos, _, used = firstReplaceable(w, pos, face, b)
	if !used {
		return
	}
	//noinspection GoAssignmentToReceiver
	b = NewBarrel()
	b.Facing = calculateFace(user, pos)

	place(w, pos, b, user, ctx)
	return placed(ctx)
}

// BreakInfo ...
func (b Barrel) BreakInfo() BreakInfo {
	return newBreakInfo(2.5, alwaysHarvestable, axeEffective, simpleDrops(append(b.inventory.Contents(), item.NewStack(b, 1))...))
}

// FlammabilityInfo ...
func (b Barrel) FlammabilityInfo() FlammabilityInfo {
	return newFlammabilityInfo(0, 0, true)
}

// DecodeNBT ...
func (b Barrel) DecodeNBT(data map[string]interface{}) interface{} {
	facing := b.Facing
	//noinspection GoAssignmentToReceiver
	b = NewBarrel()
	b.Facing = facing
	b.CustomName = nbtconv.MapString(data, "CustomName")
	nbtconv.InvFromNBT(b.inventory, nbtconv.MapSlice(data, "Items"))
	return b
}

// EncodeNBT ...
func (b Barrel) EncodeNBT() map[string]interface{} {
	if b.inventory == nil {
		facing, customName := b.Facing, b.CustomName
		//noinspection GoAssignmentToReceiver
		b = NewBarrel()
		b.Facing, b.CustomName = facing, customName
	}
	m := map[string]interface{}{
		"Items": nbtconv.InvToNBT(b.inventory),
		"id":    "Barrel",
	}
	if b.CustomName != "" {
		m["CustomName"] = b.CustomName
	}
	return m
}

// EncodeBlock ...
func (b Barrel) EncodeBlock() (string, map[string]interface{}) {
	return "minecraft:barrel", map[string]interface{}{"open_bit": boolByte(b.Open), "facing_direction": int32(b.Facing)}
}

// EncodeItem ...
func (b Barrel) EncodeItem() (name string, meta int16) {
	return "minecraft:barrel", 0
}

func allBarrels() (b []world.Block) {
	for i := cube.Face(0); i < 6; i++ {
		b = append(b, Barrel{Facing: i})
		b = append(b, Barrel{Facing: i, Open: true})
	}
	return
}
