#!/usr/bin/env python
#coding:utf8

import sys, re

infile = "type.txt"
outfile = "../server/src/msg/proto_msg.go"

# protos = [ ["proto name", val], ...]
protos = []

def parse(f_in):
	# mods = [ {"name" = xx, "val" = xx, "protos" = []}, ...]
	mods = []
	pattern = re.compile(r'''(\w+)\s*=\s*(\d+)''')
	cur_mod = None
	for line in f_in:
		if line[0] == '\n' or line[0] == '\\':
			continue
		elif line[:1] == '#':
			cur_mod = {}
			mods.append(cur_mod)
			cur_mod["protos"] = []
			r = pattern.search(line[1:-1])
			if not r:
				print("err parsing %s" % line)
				sys.exit(-1)
			name, value = r.groups()
			cur_mod["name"] = name
			cur_mod["val"] = int(value)
		elif cur_mod != None:
			cur_mod["protos"].append(line[:-1])

	for item in mods:
		val = item["val"]
		i = 1
		for p in item["protos"]:
			pair = []
			pair.append(p)
			pair.append(val << 16 | i)
			i = i + 1
			protos.append(pair)

def create(f_out):
	print(protos)
	f_out.write('''\
package msg
const(
''')
	for item in protos:
		f_out.write("\t%s = %s\n" % (item[0].upper(), item[1]))
	f_out.write(')\n')

def main():
	f_in = file(infile, 'r')
	parse(f_in)
	f_in.close()
	f_out = file(outfile, 'w')
	create(f_out)
	f_out.close()

if __name__ == '__main__':
	main()
