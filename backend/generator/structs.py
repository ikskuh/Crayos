#!/usr/bin/env python3

import typing, json, caseconverter, io 
from pathlib import Path 
from dataclasses import dataclass
from typing import Any, Optional
from enum import Enum

SCRIPT_ROOT = Path(__file__).parent 

GO_MODULE = SCRIPT_ROOT / ".." / "game" / "structs.go"

GO_TYPES: dict[type,str] = {
    int: "int",
    str: "string",
    float : "float32",
    Any: "interface{}",
    None | str : "*string",
    list[str]: "[]string",
    None | list[str]: "[]string",
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
    register_custom_type("string", cls ) # all enums are serialized as integers in go
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
    promptselection = "promptselection"
    artstudioEmpty = "artstudio-empty"
    artstudioActive = "artstudio-active"
    exhibition = "exhibition"
    exhibitionVoting = "exhibition-voting"
    exhibitionStickering = "exhibition-stickering"
    showcase = "showcase"
    gallery = "gallery"

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
    action: str 

@api_command
class VoteCommand:
    option: str


@api_command
class PlaceStickerCommand:
    sticker: str
    x: float 
    y: float  

@api_command
class SetPaintingCommand:
    path: Any 


@api_event
class EnterSessionEvent:
    sessionId: str 

@api_event
class JoinSessionFailedEvent:
    reason: str 

@api_event
class KickedEvent:
    reason: str 

@api_event
class ChangeGameViewEvent:
    view: GameView

    painting: Any 
    paintingPrompt: None | str
    paintingBackdrop: None | str
    paintingStickers: None | list[Sticker]

    availableStickers: None | list[str]

    # VotePrompt  *string  `json:"vote-prompt"`
    # VoteOptions []string `json:"vote-options"`
    pass 

@api_event
class ChangeToolModifierEvent:
    modifier: str # TODO(fqu): Add enum here

@api_event
class PaintingChangedEvent:
    path: Any

@api_event
class PlayersChangedEvent:
    players: list[str] 
    joinedPlayer: None | str
    pass 



def generate_go_file(file: io.IOBase):

    def lineout(*args):

        file.write("".join(str(a) for a in args)+"\n")

    lineout('package game')
    lineout('')
    lineout('import (')
    lineout('	"encoding/json"')
    lineout('	"errors"')
    lineout('	"reflect"')
    lineout(')')
    lineout()
    lineout("const (")
    for atype in type_registry.values():
        if not atype.dir.is_top_level():
            continue 
        lineout("\t", atype.go_tag, ' = "', atype.json_tag, '"')
    lineout(")")
    lineout()
    lineout("var JSON_TYPE_ID = map[reflect.Type]string{")
    for atype in type_registry.values():
        if not atype.dir.is_top_level():
            continue 
        lineout("\treflect.TypeOf(&",atype.name,"{}): ",atype.go_tag, ",")
    lineout("}")
    lineout()

    lineout(
"""
type Message interface {
}

func SerializeMessage(msg Message) ([]byte, error) {

	temp, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	var dummy map[string]interface{}

	err = json.Unmarshal(temp, &dummy)
	if err != nil {
		return nil, err
	}

	dummy["type"] = JSON_TYPE_ID[reflect.TypeOf(msg)]

	return json.Marshal(dummy)
}

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

	var out interface{}

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
            
            lineout("const (")

            for item in atype.pytype:

                lineout("\t", caseconverter.macrocase(atype.name), "_", caseconverter.macrocase(item.name), ' = "', item.value, '"')

            lineout(")")

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


def main():

    # preprocess the classes
    for name, atype in type_registry.items():
        
        atype.name = name 
        atype.go_tag = caseconverter.macrocase(name) + "_TAG"
        atype.json_tag = caseconverter.kebabcase(name) 

    # print(type_registry)

    with GO_MODULE.open("w") as f:
        generate_go_file(f)



if __name__ == "__main__":
    main()