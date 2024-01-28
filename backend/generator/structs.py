#!/usr/bin/env python3

import typing, json, caseconverter, io 
from pathlib import Path 
from dataclasses import dataclass
from typing import Any, Optional
from enum import Enum

SCRIPT_ROOT = Path(__file__).parent 

GO_MODULE = SCRIPT_ROOT / ".." / "game" / "structs.go"
JS_MODULE = SCRIPT_ROOT / ".." / ".." / "frontend" / "structs.js"
API_MODULE = SCRIPT_ROOT / ".." / "api.html"

class Graphics:
    pass 
    
GO_TYPES: dict[type,str] = {
    int: "int",
    str: "string",
    bool: "bool",
    float : "float32",
    Any: "interface{}",
    None | str : "*string",
    list[str]: "[]string",
    None | list[str]: "[]string",
    dict[str, bool]: "map[string]bool",
    Graphics: "Graphics",
}

assert Optional[str] == None | str 

class ApiDirection(Enum):
    event = "event"
    command = "command"
    struct = "struct"
    enum = "enum"

    def is_top_level(self) -> bool:
        return  (self == ApiDirection.event) or (self == ApiDirection.command)


@dataclass
class ApiType:
    dir: ApiDirection
    pytype: type

    name: str = ""
    json_tag: str = ""
    go_tag: str = ""

type_registry: dict[str,ApiType] = {}

def register_custom_type(name: str, cls: type) :

    # Add some useful aliases:
    GO_TYPES[cls] = f"{name}"
    GO_TYPES[None | cls] = f"*{name}"
    GO_TYPES[list[cls]] = f"[]{name}"
    GO_TYPES[None | list[cls]] = f"[]{name}"


def api_command(cls):
    assert cls.__name__ not in type_registry
    type_registry[cls.__name__] = ApiType(dir=ApiDirection.command, pytype=cls)
    return cls 

def api_event(cls):
    assert cls.__name__ not in type_registry
    type_registry[cls.__name__] = ApiType(dir=ApiDirection.event, pytype=cls)
    return cls 



def api_struct(cls: type):
    assert cls.__name__ not in type_registry
    type_registry[cls.__name__] = ApiType(dir=ApiDirection.struct, pytype=cls)
    register_custom_type(cls.__name__, cls )
    return cls 

def api_enum(cls: type):
    assert cls.__name__ not in type_registry
    type_registry[cls.__name__] = ApiType(dir=ApiDirection.enum, pytype=cls)
    register_custom_type(cls.__name__, cls ) # all enums are serialized as integers in go
    return cls

@api_struct
class Sticker:
    id: str
    x: float 
    y: float 

@api_enum
class GameView(Enum):
    title = "title"
    lobby = "lobby"
    promptselection = "promptselection" # (B)
    artstudioGeneric = "artstudio-generic" # (A) (generic voting or empty)
    artstudioActive = "artstudio-active" # (C)
    artstudioSticker = "artstudio-sticker" # (A) + Sticker Mode
    gallery = "gallery"
    announcer = "announcer"  # Shows text instead of canvas

@api_enum
class Effect(Enum):
    # none = ""
    flashlight = "flashlight"
    drunk = "drunk"
    flip = "flip"
    swap_tool = "swap_tool"
    lock_pencil = "lock_pencil"

@api_enum
class UserAction(Enum):
    setReady = "set-ready" # sets the player ready
    setNotReady = "set-not-ready" # sets the player to "not ready" in the lobby
    continueGame = "continue" # the player wants to continue/leave the current screen

@api_enum
class Backdrop(Enum):
	arctic  = "arctic"
	graveyard  = "graveyard"
	pirateShip  = "pirate_ship"
	theaterStage1  = "theater_stage1"
	desert  = "desert"

@api_command
class CreateSessionCommand:
    nickName: str


@api_command
class JoinSessionCommand:
    nickName: str 
    sessionId: str 


@api_command
class LeaveSessionCommand:
    pass 

@api_command
class UserCommand:
    action: UserAction

@api_command
class VoteCommand:
    option: str # user has voted for an option from ChangeGameViewEvent.voteOptions


@api_command
class PlaceStickerCommand:
    sticker: str
    x: float 
    y: float  

@api_command
class SetPaintingCommand:
    graphics: Graphics 


@api_event
class EnterSessionEvent:
    sessionId: str 

@api_event
class JoinSessionFailedEvent:
    reason: str 

@api_event
class KickedEvent:
    reason: str

@api_struct
class Painting:
    prompt: str # shows the current drawing prompt
    graphics: Graphics # the current painting data
    backdrop: Backdrop # the ID of the backdrop 
    stickers: list[Sticker] # the current list of stickers that should be shown
    winner: bool # the painting is the winner

@api_event
class ChangeGameViewEvent:
    view: GameView # what view the frontend should show

    painting: Painting # artstudio*: the current painting data

    results: list[Painting] # gallery: contains all winning images

    votePrompt: str # artstudioGeneric: the prompt that is shown when 
    voteOptions: list[str] # promptselection, artstudioGeneric: list of options that the player can vote for.

    announcer: str # the text shown on the announcer screen

@api_event
class TimerChangedEvent:
    secondsLeft: int 

@api_event
class ChangeToolModifierEvent:
    modifier: Effect
    duration: int # milliseconds

@api_event
class PaintingChangedEvent:
    graphics: Graphics # the new painting

@api_event
class PlayersChangedEvent:
    players: list[str] # new list of present player
    addedPlayer: None | str # player that joined
    removedPlayer: None | str # player that left

@api_event
class PlayerReadyChangedEvent:
    players: dict[str,bool] # contains a map from nick to "is ready"

@api_event
class PopUpEvent:
    message: str

@api_event
class DebugMessageEvent:
    message: str # Show this text as a debug overlay somewhere

###############################################################################

def generate_go_file(file: io.IOBase):

    def lineout(*args):

        file.write("".join(str(a) for a in args)+"\n")

    lineout('package game')
    lineout('')
    lineout('import (')
    lineout('	"encoding/json"')
    lineout('	"errors"')
    lineout(')')
    lineout()
    lineout("const (")
    for atype in type_registry.values():
        if not atype.dir.is_top_level():
            continue 
        lineout("\t", atype.go_tag, ' = "', atype.json_tag, '"')
    lineout(")")

    lineout(
"""
func DeserializeMessage(data []byte) (Message, error) {

	var raw_map map[string]interface{} // must be an object

	err := json.Unmarshal(data, &raw_map)
	if err != nil {
		return nil, err
	}

	type_tag, ok := raw_map["type"]
	if !ok {
		return nil, errors.New("Invalid json")
	}

	var out Message

	switch type_tag {
""")
    for atype in type_registry.values():
        if not atype.dir.is_top_level():
            continue 
        lineout("\tcase ",atype.go_tag,":")
        lineout("\t\tout = &",atype.name,"{}")	

    lineout("""
	default:
		return nil, errors.New("Invalid type")
	}

	err = json.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
"""
    )

    for atype in type_registry.values():

        if atype.dir == ApiDirection.enum:
            
            lineout("type ", atype.name, " string")
            lineout("const (")

            for item in atype.pytype:
                lineout("\t", caseconverter.macrocase(atype.name), "_", caseconverter.macrocase(item.name), ' ', atype.name, ' = "', item.value, '"')

            lineout(")")
            lineout("var ALL_", caseconverter.macrocase(atype.name), "_ITEMS = []",atype.name,"{")
            for item in atype.pytype:
                lineout("\t", '"', item.value, '",')
            lineout("}")

        else:

            lineout("type ", atype.name, " struct {")

            # lineout("\treflect.TypeOf(&",atype.name,"{}): ",atype.go_tag,",")

            for field, hint in typing.get_type_hints(atype.pytype).items():
                
                if hint not in GO_TYPES:
                    print("Could not find mapping for type ", hint)
                    print("available mappings are:")
                    for ktype, gtype in GO_TYPES.items():
                        print(ktype, "=>", gtype)
                    exit(1)

                go_name = caseconverter.pascalcase(field)
                go_type = GO_TYPES[hint]

                lineout("\t", go_name, " ", go_type, ' `json:"', field, '"`')

            lineout("}")
        
        lineout()

    lineout()
    for atype in type_registry.values():
        if not atype.dir.is_top_level():
            continue 
        lineout("func (item *", atype.name, ") GetJsonType() string {")
        lineout('\treturn "',atype.json_tag, '"')
        lineout("}")

        lineout("func (item *", atype.name, ") FixNils() Message {")
        lineout("\t", "copy := *item")

        for field, hint in typing.get_type_hints(atype.pytype).items():
            go_name = caseconverter.pascalcase(field)
            go_type = GO_TYPES[hint]
            if go_type.startswith("map[") or go_type.startswith("["):
                lineout("\t", "if copy.", go_name, " == nil {")
                lineout("\t\t", "copy.", go_name, " = ", go_type, "{}")
                lineout("\t", "}")

        lineout("\t", "return &copy")
        lineout("}")

        lineout()



def generate_js_file(file: io.IOBase):

    def lineout(*args):
        file.write("".join(str(a) for a in args)+"\n")

    def scoped_kv(name: str, items: dict [str,Any] ):

        lineout("const ", name, " = {")
        for key, value in items.items():
            lineout("    ", key, " : '", value, "',")
        lineout("};")

    scoped_kv("CommandId", {
        atype.name.removesuffix("Command"): atype.json_tag
        for atype in type_registry.values()
        if atype.dir == ApiDirection.command
    })
    lineout()
    scoped_kv("EventId", {
        atype.name.removesuffix("Event"): atype.json_tag
        for atype in type_registry.values()
        if atype.dir == ApiDirection.event
    })
    lineout()

    for atype in type_registry.values():
        if atype.dir == ApiDirection.enum:
            
            lineout("// Enum:")
            scoped_kv(atype.name, {
                item.name: item.value 
                for item in atype.pytype
            })

        elif atype.dir == ApiDirection.command:
           
            lineout("// Command:")
            lineout("function send", atype.name, "(", ", ".join(typing.get_type_hints(atype.pytype).keys()), ")")
            lineout("{")
            lineout("    socket.send(JSON.stringify({")
            lineout("        type : CommandId.", atype.name.removesuffix("Command"), ",")
            
            for field, hint in typing.get_type_hints(atype.pytype).items():
                lineout("        ", field, " : ", field, ", // ", hint.__name__)

            lineout("    }));")
            lineout("}")

            # lineout("type ", atype.name, " struct {")

            # # lineout("\treflect.TypeOf(&",atype.name,"{}): ",atype.go_tag,",")

                
            #     if hint not in GO_TYPES:
            #         print("Could not find mapping for type ", hint)
            #         print("available mappings are:")
            #         for ktype, gtype in GO_TYPES.items():
            #             print(ktype, "=>", gtype)
            #         exit(1)

            #     go_name = caseconverter.pascalcase(field)
            #     go_type = GO_TYPES[hint]

            #     lineout("\t", go_name, " ", go_type, ' `json:"', field, '"`')

            # lineout("}")
        else:
            # skip all unsupported types
            continue 

        lineout()


def generate_debug_file(file):
    def lineout(*args):
        file.write("".join(str(a) for a in args)+"\n")


    STATUS_FIELDS = {
        "sessionId": "Session ID",
        "view": "Current View",
        "players": "Players",
        "backdrop": "Backdrop",
        "timer": "Timer",
    }

    lineout("""<!DOCTYPE html>
<html lang="en">
  <head>
    <title>API</title>
    <style>
* {
    box-sizing: border-box;
}

body {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    position: absolute;
    margin: 0;
    padding: 1rem;
    
    width: 100%;
    height: 100%;
}

table#status {
    border-collapse: collapse;
    border: 1px solid black;
}

table#status tr {
    border: 1px solid black;    
}

table#status tr td {
    border: 1px solid black;
    padding: 5px;
}

table#status tr:nth-child(1) {
    background-color: #DDD;
}

table#status tr:nth-child(2) td {
    font-family: monospace;
}

.commands {
    display: flex;
    flex-direction: column;
    gap: 5px;
}

.commands .command {
    display: flex;
    flex-direction: row;
    gap: 5px;
}

.commands .command button {
    width: 12rem;
}

#log {
    border: 1px solid black;
    padding: 5px;
    flex: 1;
    resize: none;
    overflow: scroll;
    font-family: monospace;
}

#log p {
    margin: 0;
    padding: 0;
    white-space: pre;
}

#log button {
    margin: 2px;
    padding: 3px;
}

    </style>
    <script type="text/javascript">
        var socket;
    """)

    # just include the full API
    generate_js_file(file)

    lineout("""
        var log_area;
        var sessionId = null;

        const STATUS_FIELDS = {};

        function setStatus(field, value) {
            STATUS_FIELDS[field].innerText = String(value);
        }

        function logElement(element)
        {
            if(!log_area) {
                return;
            }
            log_area.appendChild(element);
            element.scrollIntoView();
        }

        function log(...text)
        {
            const text_node = document.createElement("p");
            text_node.innerText = text.join("");
            logElement(text_node);
        }

        function logButton(text, handler)
        {
            const btn = document.createElement("button");
            btn.innerText = text;
            btn.addEventListener("click", handler)
            logElement(btn);
        }

        function handleEnterSession(evt) {
            sessionId = evt.sessionId;
            setStatus("sessionId", sessionId);
        }

        function handleJoinSessionFailed(evt) {

        }

        function handleTimerChanged(evt) {
            setStatus("timer", evt.secondsLeft);
            return true;
        }

        function handleKicked(evt) {
            sessionId = null;
            setStatus("sessionId", "-");
        }

        function handleChangeGameView(evt) {
            setStatus("view", evt.view);
            setStatus("backdrop", evt.paintingBackdrop);

            log('ChangeGameViewEvent to ', JSON.stringify(evt.view));
            log('  painting: ', JSON.stringify(evt.painting));
            if (evt.voteOptions && evt.voteOptions.length > 0) {
                log('  vote:');
                for(const option of evt.voteOptions) {
                    logButton(option, function() {
                        sendVoteCommand(option);
                    });
                }
            } else {
                log('  vote: none');
            }


            switch(evt.view) {
                case GameView.lobby: 
                    log("Entered lobby, select if you're ready:")
                    logButton("ready", function() {
                        sendUserCommand(UserAction.setReady);
                    });
                    logButton("not ready", function() {
                        sendUserCommand(UserAction.setNotReady);
                    });
                    break;
            }
            
            return true;
        }

        function handleChangeToolModifier(evt) {

        }

        function handlePaintingChanged(evt) {

        }

        function handlePlayersChanged(evt) {
            setStatus("players", evt.players.join(", "));

            if (evt.addedPlayer) {
                log(evt.addedPlayer, " joined the game");
            }
            else if (evt.addedPlayer) {
                log(evt.removedPlayer, " joined the game");
            }
            else {
                log("Players changed to: ", evt.players.join(", "));
            }

            return true;
        }

        function handlePlayerReadyChanged(evt) {
            log(
                "READY: ", 
                Object.keys(evt.players).filter(k => evt.players[k]).join(", ") || "-",
                "\tNOT READY:",
                Object.keys(evt.players).filter(k => !evt.players[k]).join(", " || "-")
            )
            return true;
        }

""")

    for atype in type_registry.values():
        if atype.dir == ApiDirection.command:
           
            lineout("function autoSend", atype.name, "()")
            lineout("{")

            for field, hint in typing.get_type_hints(atype.pytype).items():
                lineout("    let ", field, ' = document.getElementById("', f"{atype.name}-arg-{field}", '").value;')
                
                if hint == float:
                    lineout("    ", field, " = Number(", field, ");")
                elif hint == str:
                    pass
                elif hint == Any or hint == Graphics:
                    pass 
                elif issubclass(hint, Enum):
                    pass  # enums are strings
                else:
                    print("Unsupported command type:", hint)
                    exit(1)


            lineout("    let cmd_struct = JSON.stringify({")
            lineout("        type : '", atype.json_tag, "',")
            
            for field, hint in typing.get_type_hints(atype.pytype).items():
                lineout("        ", field, " : ", field, ", // ", hint.__name__)

            lineout("    });")
            lineout("    console.log('Sending', cmd_struct);")
            lineout("    socket.send(cmd_struct);")
            lineout("}")

    

    lineout("function deserialize(msg)")
    lineout("{")
    lineout("    const obj = JSON.parse(msg);")
    lineout("    switch(obj.type) {")
    for atype in type_registry.values():
        if atype.dir == ApiDirection.event:
            lineout("    case '", atype.json_tag, "':")
            lineout("        if(handle", atype.name.removesuffix("Event"), "(obj)) {")
            lineout("            return;")
            lineout("        }")
            lineout("        log('event: ", atype.name ,"');")

            for field, hint in typing.get_type_hints(atype.pytype).items():
                lineout("        log('  ", field, ": ', JSON.stringify(obj.", field, "))")
            lineout("          log();")
            
            lineout("        break;")

    lineout("    default:")
    lineout("        log('received unknown object of type ', obj.type);")
    lineout("        break;")
    lineout("    }")
    lineout("}")

    lineout("""

        function reconnect() {
            if(socket) {
                socket.close();
            }
            socket = new WebSocket("ws://" + document.location.host + "/ws");
            socket.onclose = function (evt) {
              log("Connection closed.");
            };
            socket.onmessage = function (evt) {
                console.log("Recieved: " + evt.data);
                deserialize(evt.data)
            };
        }
    
        window.addEventListener("DOMContentLoaded", () => {
            log_area = document.getElementById("log");
            """)

    for field_name in  STATUS_FIELDS.keys():
        lineout("            STATUS_FIELDS['",field_name,"'] = document.getElementById('status-",field_name,"');");

    lineout("""
            const nick_names = [
                "xq", "manello", "captainhorst","philippwendel", "dionymoth", "Alm4nditte" 
            ];

            const nick = nick_names[Math.floor(Math.random()*nick_names.length)];

            document.getElementById("CreateSessionCommand-arg-nickName").value = nick;
            document.getElementById("JoinSessionCommand-arg-nickName").value = nick;


            reconnect();
        });

    </script>
  </head>
  <body>
    <table id="status">""")
    lineout("<tr>")
    for field_label in  STATUS_FIELDS.values():
        lineout("<td>",field_label,"</td>")
    lineout("</tr>")
    for field_label in  STATUS_FIELDS.keys():
        lineout("<td id=\"status-",field_label,"\">-</td>")
    
    lineout("""</table>
    <div class="commands">
    <div class="command">
        <button onClick="reconnect()">Reconnect ws</button>
    </div>
    """)

    for atype in type_registry.values():
        if atype.dir == ApiDirection.command:
            lineout("<div class=\"command\">")
            lineout(
                '<button onClick="autoSend', atype.name, '()">',atype.name,   '</button>'
            )

            for field, hint in typing.get_type_hints(atype.pytype).items():
                lineout("<span>", field, ":</span>")

                js_type = "text"
                if hint == float:
                    js_type = "number"

                elif issubclass(hint, Enum):
                    lineout('<select id="', f"{atype.name}-arg-{field}" ,'">')
                    for item in hint:
                        lineout('<option value="',item.value,'">',item.name,"</option>")

                    lineout("</select>")
                    continue

                lineout('<input id="', f"{atype.name}-arg-{field}" ,'" type="',js_type,'">')

            lineout("</div>")

    lineout("""
    </div>
    <div id="log"></div>
  </body>
</html>
""")

    pass 


def main():

    # preprocess the classes
    for name, atype in type_registry.items():
        
        atype.name = name 
        atype.go_tag = caseconverter.macrocase(name) + "_TAG"
        atype.json_tag = caseconverter.kebabcase(name) 

    # print(type_registry)

    with GO_MODULE.open("w") as f:
        generate_go_file(f)

    with JS_MODULE.open("w") as f:
        generate_js_file(f)

    with API_MODULE.open("w") as f:
        generate_debug_file(f)


if __name__ == "__main__":
    main()