package plugin

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"phoenixbuilder/fastbuilder/plugin_structs"
	"phoenixbuilder/minecraft"
	bridge_fmt "phoenixbuilder/session/bridge/fmt"
	"plugin"
	"unsafe"
)

func StartPluginSystem(conn *minecraft.Conn) {
	plugins := loadConfigPath()
	files, _ := ioutil.ReadDir(plugins)
	pluginbridge := plugin_structs.PluginBridge(&PluginBridgeImpl{
		sessionConnection: conn,
	})
	for _, file := range files {
		path := fmt.Sprintf("%s/%s", plugins, file.Name())
		if filepath.Ext(path) != ".so" {
			continue
		}
		go func() {
			RunPlugin(conn, path, pluginbridge)
		}()
	}
}

func RunPlugin(conn *minecraft.Conn, path string, bridge plugin_structs.PluginBridge) {
	plugin, err := plugin.Open(path)
	if err != nil {
		bridge_fmt.Printf("Failed to load plugin: %s\n%v\n", path, err)
		return
	}
	mainfunc, err := plugin.Lookup("PluginInit")
	if err != nil {
		bridge_fmt.Printf("Failed to find initial entry point for plugin %s.\n", path)
		return
	}
	mainref, err := plugin.Lookup("Main")
	if err != nil {
		bridge_fmt.Printf("Failed to find entry point for plugin %s.\n", path)
		return
	}

	name := mainfunc.(func(unsafe.Pointer, interface{}) string)(unsafe.Pointer(&bridge), mainref)
	bridge_fmt.Printf("Plugin %s(%s) loaded!\n", name, path)
}

func loadConfigPath() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		bridge_fmt.Println("[PLUGIN] WARNING - Failed to obtain the user's home directory. made homedir=\".\";\n")
		homedir = "."
	}
	fbconfigdir := filepath.Join(homedir, ".config/fastbuilder/plugins")
	os.MkdirAll(fbconfigdir, 0755)
	return fbconfigdir
}
