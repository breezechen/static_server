#ifndef _HTML_H_
#define _HTML_H_

#include <string>
#include <vector>

struct file_info
{
	bool is_folder;
	std::string filename;
	std::string url;
	std::string date;
	std::string size;
};

std::string file_contents_to_html(const std::vector<file_info> &files);

#endif