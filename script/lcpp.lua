#!/usr/bin/env lua
-- from https://github.com/facebookresearch/CParser
-- Copyright (c) Facebook, Inc. and its affiliates.
-- This source code is licensed under the MIT license found in the
-- LICENSE file in the root directory of this source tree.
--
-- Front-end for the Lua C preprocessor.


local function smartrequire(s)
   local savedpath=package.path
   local thisfile=debug.getinfo(2).source
   if type(thisfile)=='string' and thisfile:find("^@") then
      local thisdir=thisfile:gsub("^@",""):gsub("[^/]*$","")
      package.path=thisdir.."?.lua;"..package.path
   end
   local r = require(s)
   package.path = savedpath
   return r
end

cparser = smartrequire 'cparser'
io = require 'io'

local function die(s,...)
   if s then print(string.format(s,...)) end
   error(nil)
end   

local function usage()
   die([[
Usage: lcpp [options] inputfile.c [-o outputfile.c]
Main options:
  -Idirname        : Add dirname to the include search path
  -I-              : Marks beginning of system include search path
  -Dsym[=value]    : Define a macro symbol.
  -Usym            : Undefine a macro symbol.
  -Zcppdef         : Extracts initial macro definitions from cpp.
  -Znopass         : Do not copy unrecognized preprocessor directives.
  -dM              : Dump final macro definitions.
]])
end

local options={}
local outputfile
local inputfile


local outputi
for i,v in ipairs{...} do
   if outputi == i then
      outputfile = v
   elseif v == '-o' and not outputfile then
      outputi = i+1
   elseif v:find("^-") then
      options[1+#options] = v
   elseif not inputfile then
      inputfile = v
   else
      usage()
   end
end

local function exists(fname)
   local fd = io.open(fname,"r")
   if fd then fd:close() return true end
   return false
end

if not inputfile then
   usage()
elseif not exists(inputfile) then
   die("cparser: cannot read file '%s'", inputfile)
end

local outputfd
if outputfile then
   outputfd = io.open(outputfile,"w")
   if not outputfd then
      die("cparser: cannot open '%s' for writing", outputfile)
   end
end

local s,m = pcall(cparser.cpp, inputfile, outputfd, options)
if io.type(outputfd)=='file' then outputfd:close() end
if not s then die( not m:find("aborted") and m ) end

