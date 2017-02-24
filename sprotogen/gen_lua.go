package main

import "github.com/davyxu/gosproto/meta"

const luaCodeTemplate = `-- Generated by github.com/davyxu/gosproto/sprotogen
-- DO NOT EDIT!

-- Enum:
--[[
{{range $a, $enumObj := .Enums}}
-- {{$enumObj.Name}} {{range .Fields}}	
local {{$enumObj.Name}}_{{.Name}} = {{.TagNumber}} {{end}}
{{end}}
]]

local sproto = {
	Schema = [[
{{range .Structs}}
.{{.Name}} {	{{range .StFields}}	
	{{.Name}} {{.TagNumber}} : {{.CompatibleTypeString}} {{end}}
}
{{end}}
	]],

	NameByID = { {{range .Structs}}
		[{{.MsgID}}] = "{{.Name}}",{{end}}
	},
	
	IDByName = {},
}

local t = sproto.IDByName
for k, v in pairs(sproto.NameByID) do
	t[v] = k
end

return sproto

`

func gen_lua(fileD *meta.FileDescriptor, packageName, filename string) {

	fm := &fileModel{
		FileDescriptor: fileD,
		PackageName:    packageName,
	}

	addData(fm, fileD)

	generateCode("sp->lua", luaCodeTemplate, filename, fm, nil)

}
