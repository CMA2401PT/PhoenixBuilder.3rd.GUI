package nbtconv

import (
	"bytes"
	"encoding/gob"
	"phoenixbuilder/dragonfly/server/item"
	"phoenixbuilder/dragonfly/server/world"
)

// ReadItem decodes the data of an item into an item stack.
func ReadItem(data map[string]interface{}, s *item.Stack) item.Stack {
	disk := s == nil
	if disk {
		a := readItemStack(data)
		s = &a
	}
	readDamage(data, s, disk)
	readDisplay(data, s)
	readEnchantments(data, s)
	readDragonflyData(data, s)
	return *s
}

// ReadBlock decodes the data of a block into a world.Block.
func ReadBlock(m map[string]interface{}) world.Block {
	//lint:ignore S1005 Double assignment is done explicitly to prevent panics.
	nameVal, _ := m["name"]
	name, _ := nameVal.(string)
	//lint:ignore S1005 Double assignment is done explicitly to prevent panics.
	statesVal, _ := m["states"]
	properties, _ := statesVal.(map[string]interface{})

	b, _ := world.BlockByName(name, properties)
	return b
}

// readItemStack reads an item.Stack from the NBT in the map passed.
func readItemStack(m map[string]interface{}) item.Stack {
	var it world.Item
	if blockItem, ok := MapBlock(m, "Block").(world.Item); ok {
		it = blockItem
	}
	if v, ok := world.ItemByName(MapString(m, "Name"), MapInt16(m, "Damage")); ok {
		it = v
	}
	if it == nil {
		return item.Stack{}
	}
	if n, ok := it.(world.NBTer); ok {
		it = n.DecodeNBT(m).(world.Item)
	}
	return item.NewStack(it, int(MapByte(m, "Count")))
}

// readDamage reads the damage value stored in the NBT with the Damage tag and saves it to the item.Stack passed.
func readDamage(m map[string]interface{}, s *item.Stack, disk bool) {
	if disk {
		*s = s.Damage(int(MapInt16(m, "Damage")))
		return
	}
	*s = s.Damage(int(MapInt32(m, "Damage")))
}

// readEnchantments reads the enchantments stored in the ench tag of the NBT passed and stores it into an item.Stack.
func readEnchantments(m map[string]interface{}, s *item.Stack) {
	if enchantmentList, ok := m["ench"]; ok {
		enchantments, ok := enchantmentList.([]map[string]interface{})
		if !ok {
			for _, e := range MapSlice(m, "ench") {
				if v, ok := e.(map[string]interface{}); ok {
					enchantments = append(enchantments, v)
				}
			}
		}
		for _, ench := range enchantments {
			if e, ok := item.EnchantmentByID(int(MapInt16(ench, "id"))); ok {
				*s = s.WithEnchantment(e.WithLevel(int(MapInt16(ench, "lvl"))))
			}
		}
	}
}

// readDisplay reads the display data present in the display field in the NBT. It includes a custom name of the item
// and the lore.
func readDisplay(m map[string]interface{}, s *item.Stack) {
	if displayInterface, ok := m["display"]; ok {
		if display, ok := displayInterface.(map[string]interface{}); ok {
			if _, ok := display["Name"]; ok {
				// Only add the custom name if actually set.
				*s = s.WithCustomName(MapString(display, "Name"))
			}
			if loreInterface, ok := display["Lore"]; ok {
				if lore, ok := loreInterface.([]string); ok {
					*s = s.WithLore(lore...)
				} else if lore, ok := loreInterface.([]interface{}); ok {
					loreLines := make([]string, 0, len(lore))
					for _, l := range lore {
						loreLines = append(loreLines, l.(string))
					}
					*s = s.WithLore(loreLines...)
				}
			}
		}
	}
}

// readDragonflyData reads data written to the dragonflyData field in the NBT of an item and adds it to the item.Stack
// passed.
func readDragonflyData(m map[string]interface{}, s *item.Stack) {
	if customData, ok := m["dragonflyData"]; ok {
		d, ok := customData.([]byte)
		if !ok {
			if itf, ok := customData.([]interface{}); ok {
				for _, v := range itf {
					b, _ := v.(byte)
					d = append(d, b)
				}
			}
		}
		var m map[string]interface{}
		if err := gob.NewDecoder(bytes.NewBuffer(d)).Decode(&m); err != nil {
			panic("error decoding item user data: " + err.Error())
		}
		for k, v := range m {
			*s = s.WithValue(k, v)
		}
	}
}
