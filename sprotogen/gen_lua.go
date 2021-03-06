package main

import (
	"github.com/lexical005/sproto/meta"
)

const luaCodeTemplate = `-- Generated by github.com/lexical005/sproto/sprotogen
-- DO NOT EDIT!

local sproto = require "sproto/sproto"

local {{.PackageName}} = {
	Schema = sproto.parse([[
	{{range .Structs}}
	.{{.Name}} {	{{range .StFields}}	
		{{.Name}} {{.TagNumber}} : {{.CompatibleTypeString}} {{end}}
	}
	{{end}}
	]]),

	MsgID = { {{range $a, $enumObj := .Enums}}{{range .Fields}}
		{{.Name}} = {{.TagNumber}}, {{end}}{{end}}
	},
	
	MsgName = { {{range $a, $enumObj := .Enums}}{{range .Fields}}
		"{{.Name}}", {{end}}{{end}}
	},
	
	Reset = { {{range .Structs}}
		["{{.Name}}"] = function( obj ) -- {{.Name}}
			if obj == nil then return end {{range .StFields}}
			obj.{{.Name}} = {{.LuaDefaultValueString}} {{end}}
		end, {{end}}
	},
}

local Schema = csProto.Schema
local MsgName = csProto.MsgName
function csProto.Encode(msgID, msgTable)
	return Schema:encode(MsgName[msgID], msgTable)
end

function csProto.Decode(msgID, msgBytes)
	return Schema:decode(MsgName[msgID], msgBytes)
end

return {{.PackageName}}
`

func (self *fieldModel) LuaDefaultValueString() string {

	if self.Repeatd {
		return "nil"
	}

	switch self.Type {
	case meta.FieldType_Bool:
		return "false"
	case meta.FieldType_Int32,
		meta.FieldType_Int64,
		meta.FieldType_UInt32,
		meta.FieldType_UInt64,
		meta.FieldType_Integer,
		meta.FieldType_Float32,
		meta.FieldType_Float64,
		meta.FieldType_Enum:
		return "0"
	case meta.FieldType_String:
		return "\"\""
	case meta.FieldType_Struct,
		meta.FieldType_Bytes:
		return "nil"
	}

	return "unknown type" + self.Type.String()
}

func gen_lua(fm *fileModel, filename string) {

	addData(fm, "lua")

	generateCode("sp->lua", luaCodeTemplate, filename, fm, nil)

}
