package session

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"phoenixbuilder/fastbuilder/args"
	"phoenixbuilder/fastbuilder/command"
	"phoenixbuilder/fastbuilder/configuration"
	fbauth "phoenixbuilder/fastbuilder/cv4/auth"
	"phoenixbuilder/fastbuilder/function"
	I18n "phoenixbuilder/fastbuilder/i18n"
	"phoenixbuilder/fastbuilder/menu"
	"phoenixbuilder/fastbuilder/move"
	"phoenixbuilder/fastbuilder/nbtconstructor"
	"phoenixbuilder/fastbuilder/plugin"
	fbtask "phoenixbuilder/fastbuilder/task"
	"phoenixbuilder/fastbuilder/types"
	"phoenixbuilder/fastbuilder/utils"
	"phoenixbuilder/fastbuilder/world_provider"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
	bridge_fmt "phoenixbuilder/session/bridge/fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pterm/pterm"
)

type FBPlainToken struct {
	EncryptToken bool   `json:"encrypt_token"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

type SessionConfig struct {
	Lang          string `yaml:"lang" json:"lang"`
	FBUserName    string `yaml:"fb_username" json:"fb_username"`
	FBPassword    string `yaml:"fb_password" json:"fb_password"`
	FBToken       string `yaml:"fb_token" json:"fb_token"`
	ServerCode    string `yaml:"server_code" json:"server_code"`
	ServerPasswd  string `yaml:"server_passwd" json:"server_passwd"`
	RespondUser   string `yaml:"respond_user" json:"respond_user"`
	MuteWorldChat bool   `yaml:"mute_world_chat" json:"mute_world_chat"`
	iamDeveloper  bool
	// when "iamDeveloper" is true, the following fields are used,
	// otherwise, the fields are ignored (restore to default)
	NoPyRPC               bool   `yaml:"no_py_rpc" json:"no_py_rpc"`
	NBTConstructorEnabled bool   `yaml:"nbt_constructor_enable" json:"nbt_constructor_enable"`
	FBVersion             string `yaml:"fb_version" json:"fb_version"`
	FBHash                string `yaml:"fb_hash" json:"fb_hash"`
	FBCodeName            string `yaml:"fb_codename" json:"fb_codename"`
}

func NewConfig() *SessionConfig {
	return &SessionConfig{
		Lang:                  "zh_CN", // "en_US"
		FBUserName:            "",
		FBPassword:            "",
		FBToken:               "",
		ServerCode:            "",
		ServerPasswd:          "",
		RespondUser:           "",
		iamDeveloper:          false,
		MuteWorldChat:         false,
		NoPyRPC:               false,
		NBTConstructorEnabled: true,
		FBVersion:             DefaultFBVersion,
		FBHash:                DefaultFBHash,
		FBCodeName:            DefaultFBCodeName,
	}
}

type Session struct {
	// can use this to terminate the session
	stopChan chan struct{}

	// can use this to send command
	cmdChan          chan string
	closeFns         []func()
	worldChatChannel chan []string
	fbClinet         *fbauth.Client
	mcConn           *minecraft.Conn
	botRuntimeID     string
	Config           *SessionConfig
	// set/ set end callback
	CmdSetCbFn    func(X, Y, Z int)
	CmdSetEndCbFn func(X, Y, Z int)
}

//  what a pity, something global cannot be shared between different sessions
// so we need a flag to avoid multiple
// session running at the same time
var isStart bool

func init() {
	I18n.Init()
	isStart = false
}

func (config *SessionConfig) IsDeveloper() bool {
	return config.iamDeveloper
}

func NewSession(config *SessionConfig) *Session {
	// it's weird that we need to do this, because actually we can only hold one session
	// but maybe in the future we can support multiple sessions
	if isStart {
		return nil
	}

	config.iamDeveloper = false

	if !config.iamDeveloper {
		defaultConfig := NewConfig()
		config.NoPyRPC = defaultConfig.NoPyRPC
		config.NBTConstructorEnabled = defaultConfig.NBTConstructorEnabled
		config.FBVersion = defaultConfig.FBVersion
		config.FBHash = defaultConfig.FBHash
		config.FBCodeName = defaultConfig.FBCodeName
	}

	session := &Session{
		stopChan:      make(chan struct{}),
		cmdChan:       make(chan string),
		closeFns:      make([]func(), 0),
		Config:        config,
		CmdSetCbFn:    func(X, Y, Z int) {},
		CmdSetEndCbFn: func(X, Y, Z int) {},
	}
	I18n.SelectedLanguage = config.Lang
	I18n.UpdateLanguage()
	return session
}

func (s *Session) Start() (terminateChan chan string, startErr error) {
	// we need to make sure no multiple session is running
	if isStart {
		return nil, fmt.Errorf("Session is already started")
	}

	// before we start, we need to make sure that the session is valid
	// if not, we need to return an error

	err := s.beforeStart()
	if err != nil {
		return nil, err
	}

	// after we start, we need to return a channel that we can use to
	// notify reciver of this chan that the session is terminated
	// and the reason for termination

	isStart = true
	// when the session is terminated, we need to notify the caller
	c := s.afterStart()
	return c, nil
}

func (s *Session) afterStart() chan string {
	c := make(chan string)
	go s.routine(c)
	return c
}

func (s *Session) beforeStart() (err error) {
	// in this function, we need to make sure that the session is valid
	// first, we need to connect to the fb auth server and get the token
	// then, we try connecting to netease mc server

	// but first, we need to deal with the panic hidden in the code
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("Session Start Fail, because a panic occoured: \n%v", r)
		}
	}()

	// print copyright
	bridge_fmt.Print("The Following are copyright of Phoenix Builder Inside:")
	bridge_fmt.Print(I18n.T(I18n.Copyright_Notice_Headline))
	bridge_fmt.Print(I18n.T(I18n.Copyright_Notice_Line_1))
	bridge_fmt.Print(I18n.T(I18n.Copyright_Notice_Line_2))
	bridge_fmt.Print(I18n.T(I18n.Copyright_Notice_Line_3))
	bridge_fmt.Print("https://github.com/Sandertv/gophertunnel")
	bridge_fmt.Print("ファスト　ビルダー")
	bridge_fmt.Print("F A S T  B U I L D E R")
	bridge_fmt.Print("Contributors: Ruphane, CAIMEO")
	bridge_fmt.Print("Copyright (c) FastBuilder DevGroup, Bouldev 2022")
	bridge_fmt.Print("FastBuilder Phoenix " + s.Config.FBVersion)
	if I18n.ShouldDisplaySpecial() {
		bridge_fmt.Printf("%s", I18n.T(I18n.Special_Startup))
	}

	// check credentials
	if (s.Config.FBUserName == "" || s.Config.FBPassword == "") && s.Config.FBToken == "" {
		return fmt.Errorf("no credientials provided")
	}

	// check server configuration
	if s.Config.ServerCode == "" {
		return fmt.Errorf("no server code provided")
	}

	// do what phoenix builder does
	worldChatChannel := make(chan []string)
	s.worldChatChannel = worldChatChannel
	client := fbauth.CreateClient(worldChatChannel)
	s.fbClinet = client
	if s.Config.FBToken == "" {
		// we need to get the token
		tokenReq := &FBPlainToken{
			EncryptToken: true,
			Username:     s.Config.FBUserName,
			Password:     s.Config.FBPassword,
		}
		tokenReqStr, err := json.Marshal(tokenReq)
		if err != nil {
			return fmt.Errorf("cannot marshal token request to json: \n%v", err)
		}
		token := client.GetToken("", string(tokenReqStr))
		if token == "" {
			return fmt.Errorf("cannot get token: \n" + I18n.T(I18n.FBUC_LoginFailed))
		}
		s.Config.FBToken = token
	}
	bridge_fmt.Println(fmt.Sprintf("%s: %s", I18n.T(I18n.ServerCodeTrans), s.Config.ServerCode))
	dialer := minecraft.Dialer{
		ServerCode: s.Config.ServerCode,
		Password:   s.Config.ServerPasswd,
		Version:    s.Config.FBHash,
		Token:      s.Config.FBToken,
		Client:     client,
	}
	conn, err := dialer.Dial("raknet", "")
	if err != nil {
		return fmt.Errorf("cannot dial to netease mc server: (%v)", err)
	}
	s.mcConn = conn
	s.closeFns = append(s.closeFns, func() {
		conn.Close()
	})
	// override the default respond user
	if s.Config.RespondUser == "" {
		s.Config.RespondUser = client.ShouldRespondUser()
	}

	// TODO: don't make this global
	configuration.RespondUser = s.Config.RespondUser

	// set bot runtimeID
	s.botRuntimeID = fmt.Sprintf("%d", conn.GameData().EntityUniqueID)

	// send pyRPC
	if !s.Config.NoPyRPC {
		conn.WritePacket(&packet.PyRpc{
			Content: []byte{0x82, 0xc4, 0x8, 0x5f, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x5f, 0xc4, 0x5, 0x74, 0x75, 0x70, 0x6c, 0x65, 0xc4, 0x5, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x93, 0xc4, 0xc, 0x53, 0x79, 0x6e, 0x63, 0x55, 0x73, 0x69, 0x6e, 0x67, 0x4d, 0x6f, 0x64, 0x82, 0xc4, 0x8, 0x5f, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x5f, 0xc4, 0x5, 0x74, 0x75, 0x70, 0x6c, 0x65, 0xc4, 0x5, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x91, 0x90, 0xc0},
		})
		conn.WritePacket(&packet.PyRpc{
			Content: []byte{0x82, 0xc4, 0x8, 0x5f, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x5f, 0xc4, 0x5, 0x74, 0x75, 0x70, 0x6c, 0x65, 0xc4, 0x5, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x93, 0xc4, 0xf, 0x53, 0x79, 0x6e, 0x63, 0x56, 0x69, 0x70, 0x53, 0x6b, 0x69, 0x6e, 0x55, 0x75, 0x69, 0x64, 0x82, 0xc4, 0x8, 0x5f, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x5f, 0xc4, 0x5, 0x74, 0x75, 0x70, 0x6c, 0x65, 0xc4, 0x5, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x91, 0xc0, 0xc0},
		})
		conn.WritePacket(&packet.PyRpc{
			Content: []byte{0x82, 0xc4, 0x8, 0x5f, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x5f, 0xc4, 0x5, 0x74, 0x75, 0x70, 0x6c, 0x65, 0xc4, 0x5, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x93, 0xc4, 0x1f, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4c, 0x6f, 0x61, 0x64, 0x41, 0x64, 0x64, 0x6f, 0x6e, 0x73, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x46, 0x72, 0x6f, 0x6d, 0x47, 0x61, 0x63, 0x82, 0xc4, 0x8, 0x5f, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x5f, 0xc4, 0x5, 0x74, 0x75, 0x70, 0x6c, 0x65, 0xc4, 0x5, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x90, 0xc0},
		})
		conn.WritePacket(&packet.PyRpc{
			Content: bytes.Join([][]byte{[]byte{0x82, 0xc4, 0x8, 0x5f, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x5f, 0xc4, 0x5, 0x74, 0x75, 0x70, 0x6c, 0x65, 0xc4, 0x5, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x93, 0xc4, 0xb, 0x4d, 0x6f, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x43, 0x32, 0x53, 0x82, 0xc4, 0x8, 0x5f, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x5f, 0xc4, 0x5, 0x74, 0x75, 0x70, 0x6c, 0x65, 0xc4, 0x5, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x94, 0xc4, 0x9, 0x4d, 0x69, 0x6e, 0x65, 0x63, 0x72, 0x61, 0x66, 0x74, 0xc4, 0x6, 0x70, 0x72, 0x65, 0x73, 0x65, 0x74, 0xc4, 0x12, 0x47, 0x65, 0x74, 0x4c, 0x6f, 0x61, 0x64, 0x65, 0x64, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x81, 0xc4, 0x8, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0xc4},
				[]byte{byte(len(s.botRuntimeID))},
				[]byte(s.botRuntimeID),
				[]byte{0xc0},
			}, []byte{}),
		})
		conn.WritePacket(&packet.PyRpc{
			Content: []byte{0x82, 0xc4, 0x8, 0x5f, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x5f, 0xc4, 0x5, 0x74, 0x75, 0x70, 0x6c, 0x65, 0xc4, 0x5, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x93, 0xc4, 0x19, 0x61, 0x72, 0x65, 0x6e, 0x61, 0x47, 0x61, 0x6d, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x4c, 0x6f, 0x61, 0x64, 0x82, 0xc4, 0x8, 0x5f, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x5f, 0xc4, 0x5, 0x74, 0x75, 0x70, 0x6c, 0x65, 0xc4, 0x5, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x90, 0xc0},
		})
		conn.WritePacket(&packet.PyRpc{
			Content: bytes.Join([][]byte{[]byte{0x82, 0xc4, 0x8, 0x5f, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x5f, 0xc4, 0x5, 0x74, 0x75, 0x70, 0x6c, 0x65, 0xc4, 0x5, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x93, 0xc4, 0xb, 0x4d, 0x6f, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x43, 0x32, 0x53, 0x82, 0xc4, 0x8, 0x5f, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x5f, 0xc4, 0x5, 0x74, 0x75, 0x70, 0x6c, 0x65, 0xc4, 0x5, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x94, 0xc4, 0x9, 0x4d, 0x69, 0x6e, 0x65, 0x63, 0x72, 0x61, 0x66, 0x74, 0xc4, 0xe, 0x76, 0x69, 0x70, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0xc4, 0xc, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x55, 0x69, 0x49, 0x6e, 0x69, 0x74, 0xc4},
				[]byte{byte(len(s.botRuntimeID))},
				[]byte(s.botRuntimeID),
				[]byte{0xc0},
			}, []byte{}),
		})
	}

	// send the ClientCacheStatus packet
	conn.WritePacket(&packet.ClientCacheStatus{
		Enabled: false,
	})

	// init FB Functions
	function.InitInternalFunctions()

	// override the default nbt state
	if !s.Config.iamDeveloper {
		s.Config.NBTConstructorEnabled = !fbauth.ShouldDisableNBTConstructor
	}
	if s.Config.NBTConstructorEnabled {
		nbtconstructor.InitNBTConstructor()
	}

	// TODO: don't make those global
	fbtask.InitTaskStatusDisplay(conn)

	// init the bot movement sync
	move.ConnectTime = conn.GameData().ConnectTime
	move.Position = conn.GameData().PlayerPosition
	move.Pitch = conn.GameData().Pitch
	move.Yaw = conn.GameData().Yaw
	move.Connection = conn
	move.RuntimeID = conn.GameData().EntityRuntimeID

	// no necessary here
	// signalhandler.Init(conn)

	// TODO: it should have a better design
	zeroId, _ := uuid.NewUUID()
	oneId, _ := uuid.NewUUID()
	configuration.ZeroId = zeroId
	configuration.OneId = oneId
	types.ForwardedBrokSender = fbtask.BrokSender

	return nil
}

func (s *Session) routine(c chan string) {
	terminateReason := "Session terminated by user"
	defer func() {
		// we don't want the whole program to exit when there is a panic
		// hidden in the code
		r := recover()
		if r != nil {
			terminateReason = fmt.Sprintf("Session terminated, because a panic occoured: \n%v", r)
		}
		s.close()
		c <- terminateReason
	}()

	go func() {
		if s.Config.MuteWorldChat {
			return
		}
		for {
			select {
			case csmsg := <-s.worldChatChannel:
				command.WorldChatTellraw(s.mcConn, csmsg[0], csmsg[1])
			case <-s.stopChan:
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case cmd := <-s.cmdChan:
				if len(cmd) == 0 {
					continue
				}
				if cmd[0] == '.' {
					ud, _ := uuid.NewUUID()
					chann := make(chan *packet.CommandOutput)
					command.UUIDMap.Store(ud.String(), chann)
					command.SendCommand(cmd[1:], ud, s.mcConn)
					resp := <-chann
					bridge_fmt.Printf("%+v\n", resp)
				} else if cmd[0] == '!' {
					ud, _ := uuid.NewUUID()
					chann := make(chan *packet.CommandOutput)
					command.UUIDMap.Store(ud.String(), chann)
					command.SendWSCommand(cmd[1:], ud, s.mcConn)
					resp := <-chann
					bridge_fmt.Printf("%+v\n", resp)
				}
				if cmd == "menu" {
					menu.OpenMenu(s.mcConn)
					bridge_fmt.Printf("OK\n")
					continue
				}
				function.Process(s.mcConn, cmd)
			case <-s.stopChan:
				return
			}
		}
	}()

	// A loop that reads packets from the connection until it is closed.
	conn := s.mcConn
	user := s.Config.RespondUser
	zeroId := configuration.ZeroId
	client := s.fbClinet
	for {
		// Read a packet from the connection: ReadPacket returns an error if the connection is closed or if
		// a read timeout is set. You will generally want to return or break if this happens.
		pk, err := conn.ReadPacket()
		if err != nil {
			panic(err)
		}

		switch p := pk.(type) {
		case *packet.PyRpc:
			if args.NoPyRpc() {
				break
			}
			//fmt.Printf("PyRpc!\n")
			if strings.Contains(string(p.Content), "GetLoadingTime") {
				//fmt.Printf("GetLoadingTime!!\n")
				uid := s.mcConn.IdentityData().Uid
				num := uid&255 ^ (uid&65280)>>8
				curTime := time.Now().Unix()
				num = curTime&3 ^ (num&7)<<2 ^ (curTime&252)<<3 ^ (num&248)<<8
				numcont := make([]byte, 2)
				binary.BigEndian.PutUint16(numcont, uint16(num))
				s.mcConn.WritePacket(&packet.PyRpc{
					Content: []byte{0x82, 0xc4, 0x8, 0x5f, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x5f, 0xc4, 0x5, 0x74, 0x75, 0x70, 0x6c, 0x65, 0xc4, 0x5, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x93, 0xc4, 0x12, 0x53, 0x65, 0x74, 0x6c, 0x6f, 0x61, 0x64, 0x4c, 0x6f, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x82, 0xc4, 0x8, 0x5f, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x5f, 0xc4, 0x5, 0x74, 0x75, 0x70, 0x6c, 0x65, 0xc4, 0x5, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x91, 0xcd, numcont[0], numcont[1], 0xc0},
				})
				// Good job, netease, you wasted 3 days of my idle time
				// (-Ruphane)

				// See analyze/nemcfix/final.py for its python version
				// and see analyze/ for how I did it.
				s.tellraw("Welcome to FastBuilder!")
				s.tellraw(fmt.Sprintf("Operator: %s", user))
				s.sendCommand("testforblock ~ ~ ~ air", zeroId)
			} else if strings.Contains(string(p.Content), "check_server_contain_pet") {
				//fmt.Printf("Checkpet!!\n")

				// Pet req
				/*conn.WritePacket(&packet.PyRpc {
					Content: bytes.Join([][]byte{[]byte{0x82,0xc4,0x8,0x5f,0x5f,0x74,0x79,0x70,0x65,0x5f,0x5f,0xc4,0x5,0x74,0x75,0x70,0x6c,0x65,0xc4,0x5,0x76,0x61,0x6c,0x75,0x65,0x93,0xc4,0xb,0x4d,0x6f,0x64,0x45,0x76,0x65,0x6e,0x74,0x43,0x32,0x53,0x82,0xc4,0x8,0x5f,0x5f,0x74,0x79,0x70,0x65,0x5f,0x5f,0xc4,0x5,0x74,0x75,0x70,0x6c,0x65,0xc4,0x5,0x76,0x61,0x6c,0x75,0x65,0x94,0xc4,0x9,0x4d,0x69,0x6e,0x65,0x63,0x72,0x61,0x66,0x74,0xc4,0x3,0x70,0x65,0x74,0xc4,0x12,0x73,0x75,0x6d,0x6d,0x6f,0x6e,0x5f,0x70,0x65,0x74,0x5f,0x72,0x65,0x71,0x75,0x65,0x73,0x74,0x87,0xc4,0x13,0x61,0x6c,0x6c,0x6f,0x77,0x5f,0x73,0x74,0x65,0x70,0x5f,0x6f,0x6e,0x5f,0x62,0x6c,0x6f,0x63,0x6b,0xc2,0xc4,0xb,0x61,0x76,0x6f,0x69,0x64,0x5f,0x6f,0x77,0x6e,0x65,0x72,0xc3,0xc4,0x7,0x73,0x6b,0x69,0x6e,0x5f,0x69,0x64,0xcd,0x27,0x11,0xc4,0x9,0x70,0x6c,0x61,0x79,0x65,0x72,0x5f,0x69,0x64,0xc4},
							[]byte{byte(len(runtimeid))},
								[]byte(runtimeid),
								[]byte{0xc4,0x6,0x70,0x65,0x74,0x5f,0x69,0x64,0x1,0xc4,0xa,0x6d,0x6f,0x64,0x65,0x6c,0x5f,0x6e,0x61,0x6d,0x65,0xc4,0x14,0x74,0x79,0x5f,0x79,0x75,0x61,0x6e,0x73,0x68,0x65,0x6e,0x67,0x68,0x75,0x6c,0x69,0x5f,0x30,0x5f,0x30,0xc4,0x4,0x6e,0x61,0x6d,0x65,0xc4,0xc,0xe6,0x88,0x91,0xe7,0x9a,0x84,0xe4,0xbc,0x99,0xe4,0xbc,0xb4,0xc0},
						},[]byte{}),
				})*/

			} else if strings.Contains(string(p.Content), "summon_pet_response") {
				//fmt.Printf("summonpetres\n")
				/*conn.WritePacket(&packet.PyRpc {
					Content: []byte{0x82,0xc4,0x8,0x5f,0x5f,0x74,0x79,0x70,0x65,0x5f,0x5f,0xc4,0x5,0x74,0x75,0x70,0x6c,0x65,0xc4,0x5,0x76,0x61,0x6c,0x75,0x65,0x93,0xc4,0x19,0x61,0x72,0x65,0x6e,0x61,0x47,0x61,0x6d,0x65,0x50,0x6c,0x61,0x79,0x65,0x72,0x46,0x69,0x6e,0x69,0x73,0x68,0x4c,0x6f,0x61,0x64,0x82,0xc4,0x8,0x5f,0x5f,0x74,0x79,0x70,0x65,0x5f,0x5f,0xc4,0x5,0x74,0x75,0x70,0x6c,0x65,0xc4,0x5,0x76,0x61,0x6c,0x75,0x65,0x90,0xc0},
				})
				conn.WritePacket(&packet.PyRpc {
					Content: bytes.Join([][]byte{[]byte{0x82,0xc4,0x8,0x5f,0x5f,0x74,0x79,0x70,0x65,0x5f,0x5f,0xc4,0x5,0x74,0x75,0x70,0x6c,0x65,0xc4,0x5,0x76,0x61,0x6c,0x75,0x65,0x93,0xc4,0xb,0x4d,0x6f,0x64,0x45,0x76,0x65,0x6e,0x74,0x43,0x32,0x53,0x82,0xc4,0x8,0x5f,0x5f,0x74,0x79,0x70,0x65,0x5f,0x5f,0xc4,0x5,0x74,0x75,0x70,0x6c,0x65,0xc4,0x5,0x76,0x61,0x6c,0x75,0x65,0x94,0xc4,0x9,0x4d,0x69,0x6e,0x65,0x63,0x72,0x61,0x66,0x74,0xc4,0xe,0x76,0x69,0x70,0x45,0x76,0x65,0x6e,0x74,0x53,0x79,0x73,0x74,0x65,0x6d,0xc4,0xc,0x50,0x6c,0x61,0x79,0x65,0x72,0x55,0x69,0x49,0x6e,0x69,0x74,0xc4},
							[]byte{byte(len(runtimeid))},
								[]byte(runtimeid),
								[]byte{0xc0},
							},[]byte{}),
				})*/
				/*conn.WritePacket(&packet.PyRpc {
					Content: []byte{0x82,0xc4,0x8,0x5f,0x5f,0x74,0x79,0x70,0x65,0x5f,0x5f,0xc4,0x5,0x74,0x75,0x70,0x6c,0x65,0xc4,0x5,0x76,0x61,0x6c,0x75,0x65,0x93,0xc4,0x1f,0x43,0x6c,0x69,0x65,0x6e,0x74,0x4c,0x6f,0x61,0x64,0x41,0x64,0x64,0x6f,0x6e,0x73,0x46,0x69,0x6e,0x69,0x73,0x68,0x65,0x64,0x46,0x72,0x6f,0x6d,0x47,0x61,0x63,0x82,0xc4,0x8,0x5f,0x5f,0x74,0x79,0x70,0x65,0x5f,0x5f,0xc4,0x5,0x74,0x75,0x70,0x6c,0x65,0xc4,0x5,0x76,0x61,0x6c,0x75,0x65,0x90,0xc0},
				})*/
			} else if strings.Contains(string(p.Content), "GetStartType") {
				// 2021-12-22 10:51~11:55
				// Thank netease for wasting my time again ;)
				encData := p.Content[68 : len(p.Content)-1]
				response := client.TransferData(string(encData), fmt.Sprintf("%d", conn.IdentityData().Uid))
				conn.WritePacket(&packet.PyRpc{
					Content: bytes.Join([][]byte{[]byte{0x82, 0xc4, 0x8, 0x5f, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x5f, 0xc4, 0x5, 0x74, 0x75, 0x70, 0x6c, 0x65, 0xc4, 0x5, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x93, 0xc4, 0xc, 0x53, 0x65, 0x74, 0x53, 0x74, 0x61, 0x72, 0x74, 0x54, 0x79, 0x70, 0x65, 0x82, 0xc4, 0x8, 0x5f, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x5f, 0xc4, 0x5, 0x74, 0x75, 0x70, 0x6c, 0x65, 0xc4, 0x5, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x91, 0xc4},
						[]byte{byte(len(response))},
						[]byte(response),
						[]byte{0xc0},
					}, []byte{}),
				})
			}
			break
		case *packet.StructureTemplateDataResponse:
			//fmt.Printf("RESPONSE %+v\n",p.StructureTemplate)
			fbtask.ExportWaiter <- p.StructureTemplate
			break
		/*case *packet.InventoryContent:
		for _, item := range p.Content {
			fmt.Printf("InventorySlot %+v\n",item.Stack.NBTData["dataField"])
		}
		break*/
		/*case *packet.InventorySlot:
		fmt.Printf("Slot %d:%+v",p.Slot,p.NewItem.Stack)*/
		case *packet.Text:
			if p.TextType == packet.TextTypeChat {
				for _, item := range plugin.ChatEventListeners {
					item(p.SourceName, p.Message)
				}
				if user == p.SourceName {
					if p.Message[0] == '>' && len(p.Message) > 1 {
						umsg := p.Message[1:]
						if !client.CanSendMessage() {
							command.WorldChatTellraw(conn, "FasｔBuildeｒ", "Lose connection to the authentication server.")
							break
						}
						client.WorldChat(umsg)
					}
					break
					pterm.Println(pterm.Yellow(fmt.Sprintf("<%s>", user)), pterm.LightCyan(p.Message))
					if p.Message[0] == '>' {
						//umsg:=p.Message[1:]
						//
					}
					function.Process(conn, p.Message)
					break
				}
			}
		case *packet.CommandOutput:
			//if p.SuccessCount > 0 {
			if p.CommandOrigin.UUID.String() == configuration.ZeroId.String() {
				pos, _ := utils.SliceAtoi(p.OutputMessages[0].Parameters)
				if !(p.OutputMessages[0].Message == "commands.generic.unknown") {
					configuration.IsOp = true
				}
				if len(pos) == 0 {
					s.tellraw(I18n.T(I18n.InvalidPosition))
					break
				}
				configuration.GlobalFullConfig().Main().Position = types.Position{
					X: pos[0],
					Y: pos[1],
					Z: pos[2],
				}
				s.CmdSetCbFn(pos[0], pos[1], pos[2])
				s.tellraw(fmt.Sprintf("%s: %v", I18n.T(I18n.PositionGot), pos))
				break
			} else if p.CommandOrigin.UUID.String() == configuration.OneId.String() {
				pos, _ := utils.SliceAtoi(p.OutputMessages[0].Parameters)
				if len(pos) == 0 {
					s.tellraw(I18n.T(I18n.InvalidPosition))
					break
				}
				configuration.GlobalFullConfig().Main().End = types.Position{
					X: pos[0],
					Y: pos[1],
					Z: pos[2],
				}
				s.CmdSetEndCbFn(pos[0], pos[1], pos[2])
				s.tellraw(fmt.Sprintf("%s: %v", I18n.T(I18n.PositionGot_End), pos))
				break
			}
			//}
			pr, ok := command.UUIDMap.LoadAndDelete(p.CommandOrigin.UUID.String())
			if ok {
				pu := pr.(chan *packet.CommandOutput)
				pu <- p
			}
		case *packet.ActorEvent:
			if p.EventType == packet.ActorEventDeath && p.EntityRuntimeID == conn.GameData().EntityRuntimeID {
				conn.WritePacket(&packet.PlayerAction{
					EntityRuntimeID: conn.GameData().EntityRuntimeID,
					ActionType:      protocol.PlayerActionRespawn,
				})
			}
		case *packet.LevelChunk:
			if world_provider.ChunkInput != nil {
				world_provider.ChunkInput <- p
			} else {
				world_provider.DoCache(p)
			}
		case *packet.UpdateBlock:
			channel, h := command.BlockUpdateSubscribeMap.LoadAndDelete(p.Position)
			if h {
				ch := channel.(chan bool)
				ch <- true
			}
		case *packet.AddActor:
			if p.EntityType == "minecraft:villager_v2" {
				if nbtconstructor.AddVillagerChannel != nil {
					nbtconstructor.AddVillagerChannel <- p
				}
			}
		case *packet.InventoryContent:
			if p.WindowID == 0 {
				if len(p.Content) == 0 {
					break
				}
				if nbtconstructor.InventoryContentChannel != nil {
					nbtconstructor.InventoryContentChannel <- p
				}
			}
		case *packet.ItemStackResponse:
			if nbtconstructor.ItemStackResponseChannel != nil {
				nbtconstructor.ItemStackResponseChannel <- p
			}
		case *packet.UpdateTrade:
			if nbtconstructor.IsWorking {
				nbtconstructor.TradeWindowID = p.WindowID
			}
		case *packet.Respawn:
			if p.EntityRuntimeID == conn.GameData().EntityRuntimeID {
				move.Position = p.Position
			}
		case *packet.MovePlayer:
			if p.EntityRuntimeID == conn.GameData().EntityRuntimeID {
				move.Position = p.Position
			} else if p.EntityRuntimeID == move.TargetRuntimeID {
				move.Target = p.Position
			}
		case *packet.CorrectPlayerMovePrediction:
			//fmt.Printf("correct %v\n",time.Now())
			move.MoveP += 10
			if move.MoveP > 100 {
				move.MoveP = 0
			}
			move.Position = p.Position
			move.Jump()
		case *packet.AddPlayer:
			if move.TargetRuntimeID == 0 && p.EntityRuntimeID != conn.GameData().EntityRuntimeID {
				move.Target = p.Position
				move.TargetRuntimeID = p.EntityRuntimeID
				//fmt.Printf("Got target: %s\n",p.Username)
			}
		}
		select {
		case <-s.stopChan:
			terminateReason = "session terminated by user"
			return
		default:
		}
	}
}

func (s *Session) GetPos() (x, y, z int) {
	return configuration.GlobalFullConfig().Main().Position.X, configuration.GlobalFullConfig().Main().Position.Y, configuration.GlobalFullConfig().Main().Position.Z
}

func (s *Session) GetEndPos() (x, y, z int) {
	return configuration.GlobalFullConfig().Main().End.X, configuration.GlobalFullConfig().Main().End.Y, configuration.GlobalFullConfig().Main().End.Z
}

func (s *Session) sendCommand(commands string, UUID uuid.UUID) error {
	return command.SendCommand(commands, UUID, s.mcConn)
}

func (s *Session) tellraw(lines ...string) error {
	command.AdditionalChatCb(lines[0])
	return command.Tellraw(s.mcConn, lines[0])
}

func (s *Session) close() {
	isStart = false
	for _, fn := range s.closeFns {
		fn()
	}
	// let GC do the work
	s.fbClinet = nil
	s.mcConn = nil
}

func (s *Session) Execute(cmd string) {
	s.cmdChan <- cmd
}

func (s *Session) Stop() {
	// close the stopChan to nofitify the routine to stop session
	close(s.stopChan)
}
