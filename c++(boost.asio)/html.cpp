#include "html.h"
#include "webix_css.h"
#include "webix_js.h"
#include <fstream>
#include <sstream>
#include <boost/filesystem.hpp>

const unsigned char* get_webix_css()
{
    boost::filesystem::path p("webix.css");
    if (!boost::filesystem::exists(p)) {
        return webix_css_str;
    }

    static std::time_t lastmtime = 0;
    static std::string* static_webix_css = new std::string;

    std::time_t mtime = boost::filesystem::last_write_time(p);
    if (mtime != lastmtime) {
        std::ifstream is(p.string().c_str(), std::ios::in | std::ios::binary);
        std::stringstream ss;
        ss << is.rdbuf();
        *static_webix_css = ss.str();
    }
    return (const unsigned char*)static_webix_css->c_str();
}

const unsigned char* get_webix_js()
{
    boost::filesystem::path p("webix.js");
    if (!boost::filesystem::exists(p)) {
        return webix_js_str;
    }

    static std::time_t lastmtime = 0;
    static std::string* static_webix_js = new std::string;

    std::time_t mtime = boost::filesystem::last_write_time(p);
    if (mtime != lastmtime) {
        std::ifstream is(p.string().c_str(), std::ios::in | std::ios::binary);
        std::stringstream ss;
        ss << is.rdbuf();
        *static_webix_js = ss.str();
    }
    return (const unsigned char*)static_webix_js->c_str();
}

std::string file_contents_to_html(const std::vector<file_info> &files)
{
    std::stringstream ss;

    ss << "<!DOCTYPE html>" << std::endl;
    ss << "<html><head><meta charset=\"utf-8\"><style type=\"text/css\">" << std::endl;
    ss << get_webix_css() << std::endl;
    ss << "</style><script type=\"text/javascript\">" << std::endl;
    ss << get_webix_js() << std::endl;
    ss << "</script>" << std::endl;
    ss << "<title>index of " << files[0].filename << "</title>" << std::endl;
    ss << "</head><body>" << std::endl;
    ss << "<script type=\"text/javascript\">" << std::endl;
    ss << "var file_contents = [" << std::endl;

    std::string curl = files[0].url;
    std::string::size_type len = curl.length();
    std::string purl = curl.substr(curl.rfind("%2f") + 3);
    if (len > 3) {
        ss << "{filename:\"[上级目录]\",url:\"..\",size:\"\",date:\"\"}," << std::endl;
    }
    for (int i = 1; i < files.size(); ++i)
    {
        ss << "{filename:\"" ;
        if (files[i].is_folder) {
            ss << "<div class=\\\"webix_tree_folder\\\"></div>";
        } else {
            ss << "<div class=\\\"webix_tree_file\\\"></div>";
        }
        ss << files[i].filename;
        ss <<  "\",";
        ss << "url:\"";
        if (!purl.empty()) {
            ss << purl << "/" ;
        }
        ss << files[i].url << "\",";
        ss << "size: \"" << files[i].size << "\",";
        ss << "date: \"" << files[i].date << "\"}," << std::endl;
    }
    ss << "];" << std::endl;
    ss << "webix.ready(function(){webix.ui({view:\"scrollview\",body:{type:\"space\",autoheight:!0,cols:[{type:\"clean\",id:\"findex\",view:\"datatable\",columns:[{id:\"filename\",header:\"文件名\",fillspace:!0,template:\"<a href='#url#'>#filename#</a>\"},{id:\"size\",header:\"大小\",width:80,css:{\"text-align\":\"right\"}},{id:\"date\",header:\"修改日期\",width:180}],data:file_contents},{width:400,multi:!0,rows:[{header:\"远程下载\",body:{},collapsed:!0},{header:\"上传\",collapsed:!0,body:{}},{}]}]}})});" << std::endl;
    ss << "</script></body></html>" << std::endl;
    return ss.str();
}