#include "html.h"
#include <sstream>
#include <boost/filesystem.hpp>

const char* html_js = 
"function addRow(a, b, c, j, d) {"
"    if (a == \".\") {"
"        return;"
"    }"
"    var g = \"\" + document.location;"
"    if (g.substr(-1) !== \"/\") {"
"        g += \"/\";"
"    }"
"    var h = document.getElementById(\"table\");"
"    var i = document.createElement(\"tr\");"
"    var e = document.createElement(\"td\");"
"    var f = document.createElement(\"a\");"
"    f.className = c ? \"icon dir\" : \"icon file\";"
"    if (a == \"..\") {"
"        f.href = g + \"..\";"
"        f.innerText = document.getElementById(\"parentDirText\").innerText;"
"        f.className = \"icon up\";"
"        j = \"\";"
"        d = \"\";"
"    } else {"
"        if (c) {"
"            a = a + \"/\";"
"            b = b + \"/\";"
"            j = \"\";"
"        } else {"
"            f.draggable = \"true\";"
"            f.addEventListener(\"dragstart\", onDragStart, false);"
"        }"
"        f.innerText = a;"
"        f.href = g + b;"
"    }"
"    e.appendChild(f);"
"    i.appendChild(e);"
"    i.appendChild(createCell(j));"
"    i.appendChild(createCell(d));"
"    h.appendChild(i);"
"}"
""
"function onDragStart(d) {"
"    var c = d.srcElement;"
"    var b = c.innerText.replace(\":\", \"\");"
"    var a = \"application/octet-stream:\" + b + \":\" + c.href;"
"    d.dataTransfer.setData(\"DownloadURL\", a);"
"    d.dataTransfer.effectAllowed = \"copy\";"
"}"
""
"function createCell(b) {"
"    var a = document.createElement(\"td\");"
"    a.setAttribute(\"class\", \"detailsColumn\");"
"    a.innerText = b;"
"    return a;"
"}"
""
"function start(a) {"
"    var b = document.getElementById(\"header\");"
"    b.innerText = b.innerText.replace(\"LOCATION\", a);"
"    document.getElementById(\"title\").innerText = b.innerText;"
"}"
""
"function onListingParsingError() {"
"    var a = document.getElementById(\"listingParsingErrorBox\");"
"    a.innerHTML = a.innerHTML.replace(\"LOCATION\", encodeURI(document.location) + \"?raw\");"
"    a.style.display = \"block\";"
"}";

const char* html_css = 
"body {"
"    font-family: Consolas, Monaco, Lucida Console, monospace;"
"}"
"h1 {"
"    border-bottom: 1px solid #c0c0c0;"
"    margin-bottom: 10px;"
"    padding-bottom: 10px;"
"    white-space: nowrap;"
"}"
"table {"
"    border-collapse: collapse;"
"}"
"tr.header {"
"    font-weight: bold;"
"}"
"td.detailsColumn {"
"    -webkit-padding-start: 2em;"
"    text-align: end;"
"    white-space: nowrap;"
"}"
"a.icon {"
"    -webkit-padding-start: 1.5em;"
"    text-decoration: none;"
"}"
"a.icon:hover {"
"    text-decoration: underline;"
"}"
"a.file {"
"    background: url(\"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAIAAACQkWg2AAAABnRSTlMAAAAAAABupgeRAAABHUlEQVR42o2RMW7DIBiF3498iHRJD5JKHurL+CRVBp+i2T16tTynF2gO0KSb5ZrBBl4HHDBuK/WXACH4eO9/CAAAbdvijzLGNE1TVZXfZuHg6XCAQESAZXbOKaXO57eiKG6ft9PrKQIkCQqFoIiQFBGlFIB5nvM8t9aOX2Nd18oDzjnPgCDpn/BH4zh2XZdlWVmWiUK4IgCBoFMUz9eP6zRN75cLgEQhcmTQIbl72O0f9865qLAAsURAAgKBJKEtgLXWvyjLuFsThCSstb8rBCaAQhDYWgIZ7myM+TUBjDHrHlZcbMYYk34cN0YSLcgS+wL0fe9TXDMbY33fR2AYBvyQ8L0Gk8MwREBrTfKe4TpTzwhArXWi8HI84h/1DfwI5mhxJamFAAAAAElFTkSuQmCC \") left top no-repeat;"
"}"
"a.dir {"
"    background: url(\"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAd5JREFUeNqMU79rFUEQ/vbuodFEEkzAImBpkUabFP4ldpaJhZXYm/RiZWsv/hkWFglBUyTIgyAIIfgIRjHv3r39MePM7N3LcbxAFvZ2b2bn22/mm3XMjF+HL3YW7q28YSIw8mBKoBihhhgCsoORot9d3/ywg3YowMXwNde/PzGnk2vn6PitrT+/PGeNaecg4+qNY3D43vy16A5wDDd4Aqg/ngmrjl/GoN0U5V1QquHQG3q+TPDVhVwyBffcmQGJmSVfyZk7R3SngI4JKfwDJ2+05zIg8gbiereTZRHhJ5KCMOwDFLjhoBTn2g0ghagfKeIYJDPFyibJVBtTREwq60SpYvh5++PpwatHsxSm9QRLSQpEVSd7/TYJUb49TX7gztpjjEffnoVw66+Ytovs14Yp7HaKmUXeX9rKUoMoLNW3srqI5fWn8JejrVkK0QcrkFLOgS39yoKUQe292WJ1guUHG8K2o8K00oO1BTvXoW4yasclUTgZYJY9aFNfAThX5CZRmczAV52oAPoupHhWRIUUAOoyUIlYVaAa/VbLbyiZUiyFbjQFNwiZQSGl4IDy9sO5Wrty0QLKhdZPxmgGcDo8ejn+c/6eiK9poz15Kw7Dr/vN/z6W7q++091/AQYA5mZ8GYJ9K0AAAAAASUVORK5CYII= \") left top no-repeat;"
"}"
"a.up {"
"    background: url(\"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAmlJREFUeNpsU0toU0EUPfPysx/tTxuDH9SCWhUDooIbd7oRUUTMouqi2iIoCO6lceHWhegy4EJFinWjrlQUpVm0IIoFpVDEIthm0dpikpf3ZuZ6Z94nrXhhMjM3c8895977BBHB2PznK8WPtDgyWH5q77cPH8PpdXuhpQT4ifR9u5sfJb1bmw6VivahATDrxcRZ2njfoaMv+2j7mLDn93MPiNRMvGbL18L9IpF8h9/TN+EYkMffSiOXJ5+hkD+PdqcLpICWHOHc2CC+LEyA/K+cKQMnlQHJX8wqYG3MAJy88Wa4OLDvEqAEOpJd0LxHIMdHBziowSwVlF8D6QaicK01krw/JynwcKoEwZczewroTvZirlKJs5CqQ5CG8pb57FnJUA0LYCXMX5fibd+p8LWDDemcPZbzQyjvH+Ki1TlIciElA7ghwLKV4kRZstt2sANWRjYTAGzuP2hXZFpJ/GsxgGJ0ox1aoFWsDXyyxqCs26+ydmagFN/rRjymJ1898bzGzmQE0HCZpmk5A0RFIv8Pn0WYPsiu6t/Rsj6PauVTwffTSzGAGZhUG2F06hEc9ibS7OPMNp6ErYFlKavo7MkhmTqCxZ/jwzGA9Hx82H2BZSw1NTN9Gx8ycHkajU/7M+jInsDC7DiaEmo1bNl1AMr9ASFgqVu9MCTIzoGUimXVAnnaN0PdBBDCCYbEtMk6wkpQwIG0sn0PQIUF4GsTwLSIFKNqF6DVrQq+IWVrQDxAYQC/1SsYOI4pOxKZrfifiUSbDUisif7XlpGIPufXd/uvdvZm760M0no1FZcnrzUdjw7au3vu/BVgAFLXeuTxhTXVAAAAAElFTkSuQmCC \") left top no-repeat;"
"}"
"html[dir=rtl] a {"
"    background-position-x: right;"
"}"
"#listingParsingErrorBox {"
"    border: 1px solid black;"
"    background: #fae691;"
"    padding: 10px;"
"    display: none;"
"}";

const char* get_js()
{
    return html_js;
}

const char* get_css()
{
    return html_css;
}

std::string file_contents_to_html(const std::vector<file_info> &files)
{
    std::stringstream ss;
    std::string::size_type len = files[0].url.length();

    ss << "<!DOCTYPE html>" << std::endl;
    ss << "<html i18n-values=dir:textdirection;lang:language><head><meta charset=utf-8>" << std::endl;
    ss << "<script>" << get_js() << "</script>" << std::endl;
    ss << "<style>" << get_css() << "</style>" << std::endl;
    ss << "<title>index of " << files[0].filename << "</title>" << std::endl;
    ss << "</head><body>" << std::endl;
    ss << "<span id=\"parentDirText\" style=\"display:none\" i18n-content=\"parentDirText\">[上级目录]</span>" << std::endl;
    ss << "<table id=\"table\">" << std::endl;
    ss << "<tbody><tr class=\"header\">" << std::endl;
    ss << "<td i18n-content=\"headerName\">名称</td>" << std::endl;
    ss << "<td class=\"detailsColumn\" i18n-content=\"headerSize\">大小</td>" << std::endl;
    ss << "<td class=\"detailsColumn\" i18n-content=\"headerDateModified\">修改日期</td>" << std::endl;
    ss << "</tr></tbody>" << std::endl;
    ss << "</table></body></html>" << std::endl;
    ss << "<script>start(\"" << files[0].filename << "\");</script>" << std::endl;


    if (len > 3) {
        ss << "<script>addRow(\"..\", \"..\", 1, \"0 B\", \"\"); </script>" << std::endl;
    }

    for (int i = 1; i < files.size(); ++i)
    {
        ss << "<script>addRow(\"" << files[i].filename << "\",\"" << files[i].url << "\",";
        if (files[i].is_folder) {
            ss << "1, \"0 B\", ";
        } else {
            ss << "0, \"" << files[i].size << "\", ";
        }
        ss << "\"" << files[i].date << "\"" << ");</script>" << std::endl;
    }
    return ss.str();
}