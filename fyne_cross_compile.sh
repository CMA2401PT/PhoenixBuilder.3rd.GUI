ICON="-icon unbundled_assets/Icon.png"
VERSION="-app-version 0.0.4"
APPID='-app-id "phoenixbuilder.third.gui"'
APPBUILD='-app-build 169'
ARGS="$ICON $VERSION $APPID $APPBUILD"

fyne-cross linux -arch=amd64 $ARGS -name "FastBuilder_3rd_Gui"
fyne-cross darwin -arch=amd64 $ARGS  -name "FastBuilder_3rd_Gui"
fyne-cross windows -arch=amd64 $ARGS -name "FastBuilder_3rd_Gui.exe"
fyne-cross android -arch=arm64 $ARGS -name "FastBuilder_3rd_Gui"
fyne-cross ios $ARGS -name "FastBuilder-3rd-Gui"