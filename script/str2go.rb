#!/usr/bin/env ruby
#
# str2go <package_name> <constant_name> <path_to_file>
#

package_name  = ARGV[0]
constant_name = ARGV[1]
path_to_file  = ARGV[2]

data = File.read path_to_file

out = $stdout

def line_to_go(line)
  '"' + line.gsub("\n", '\n').gsub('"', '\"').gsub(/[^[:ascii:]]/) do |char|
    "\\#{char.ord}"
  end + '"'
end

def lines_to_go(lines)
  lines.map do |line|
    "\t#{line_to_go(line)}"
  end.join(" +\n").sub("\t", '')
end

out.puts "package #{package_name}"
out.puts

out.puts "const #{constant_name} = #{lines_to_go(data.lines)}"
