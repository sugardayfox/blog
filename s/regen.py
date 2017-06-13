#!/usr/bin/env python

"""
Re-generate html files from markdown (.md) files.

pip install markdown web.py
"""

import os
import codecs
import markdown
import web

def read_file_utf8(path):
	with codecs.open(path, "r", "utf8") as fo:
		s = fo.read()
	return s

def write_file_utf8(path, s):
	with codecs.open(path, "w", "utf8") as fo:
		s = fo.write(s)

def ext(path):
    return os.path.splitext(path)[1].lower()

# returns full paths of files in a given directory, potentially recursively,
# potentially filtering file names by filter_func (which takes file path as
# an argument)
def list_files_g(d, filter_func=None, recur=False):
    to_visit = [d]
    while len(to_visit) > 0:
        d = to_visit.pop(0)
        for f in os.listdir(d):
            path = os.path.join(d, f)
            isdir = os.path.isdir(path)
            if isdir:
                if recur:
                    to_visit.append(path)
            else:
                if filter_func is not None:
                    if filter_func(path):
                        yield path
                else:
                    yield path

# generator => array
def list_files(d, filter_func=None, recur=False):
    return [path for path in list_files_g(d, filter_func, recur)]

def delete_file(path):
    if os.path.exists(path):
        os.remove(path)

def is_markdown_file(path):
	return ext(path) in [".md"]

class MdInfo(object):
	def __init__(self, meta_data, s):
		self.meta_data = meta_data
		self.s = s

# returns MdInfo from content of the .md file
def parse_md(s):
	lines = s.split("\n")
	# lines at the top that are in the format:
	# Key: value
	# are considered meta-data
	meta_data = {}
	while len(lines) > 0:
		l = lines[0]
		parts = l.split(":", 1)
		if len(parts) != 2:
			break
		key = parts[0].lower().strip()
		val = parts[1].strip()
		meta_data[key] = val
		lines.pop(0)
	s = "\n".join(lines)
	return MdInfo(meta_data, s)

def tmpl_for_src_path(src):
	dir = os.path.dirname(src)
	path = os.path.join(dir, "_md.tmpl.html")
	tmpl_data = open(path).read()
	return web.template.Template(tmpl_data, filename="md_tmpl.html")

def md_to_html(src, dst):
	s = read_file_utf8(src)
	md_info = parse_md(s)
	body = markdown.markdown(md_info.s)
	tmpl = tmpl_for_src_path(src)
	#print("Found template: %s" % mdtmpl)
	title = md_info.meta_data["title"]
	#print(vars.keys())
	html = str(tmpl(title, body))
	util.delete_file(dst)
	print("wrote %s" % dst)
	write_file_utf8(dst, html)

def main():
	md_files = list_files("www", is_markdown_file, recur=True)
	for md_file in md_files:
		html_file = md_file[:-2] + "html"
		md_to_html(md_file, html_file)

if __name__ == "__main__":
	main()
