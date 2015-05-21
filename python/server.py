#coding=utf-8
import sys
import os
import socket
import datetime

from mimetypes import guess_type
from gevent.server import StreamServer
from gevent.monkey import patch_all
from urllib import quote, unquote

html_str = '''
<!DOCTYPE html>
<html i18n-values=dir:textdirection;lang:language>
<head>
<meta charset=utf-8>
<script>function addRow(a,b,c,j,d){if(a=="."){return}var g=""+document.location;if(g.substr(-1)!=="/"){g+="/"}var h=document.getElementById("table");var i=document.createElement("tr");var e=document.createElement("td");var f=document.createElement("a");f.className=c?"icon dir":"icon file";if(a==".."){f.href=g+"..";f.innerText=document.getElementById("parentDirText").innerText;f.className="icon up";j="";d=""}else{if(c){a=a+"/";b=b+"/";j=""}else{f.draggable="true";f.addEventListener("dragstart",onDragStart,false)}f.innerText=a;f.href=g+b}e.appendChild(f);i.appendChild(e);i.appendChild(createCell(j));i.appendChild(createCell(d));h.appendChild(i)}function onDragStart(d){var c=d.srcElement;var b=c.innerText.replace(":","");var a="application/octet-stream:"+b+":"+c.href;d.dataTransfer.setData("DownloadURL",a);d.dataTransfer.effectAllowed="copy"}function createCell(b){var a=document.createElement("td");a.setAttribute("class","detailsColumn");a.innerText=b;return a}function start(a){var b=document.getElementById("header");b.innerText=b.innerText.replace("LOCATION",a);document.getElementById("title").innerText=b.innerText}function onListingParsingError(){var a=document.getElementById("listingParsingErrorBox");a.innerHTML=a.innerHTML.replace("LOCATION",encodeURI(document.location)+"?raw");a.style.display="block"}</script>
<style>body{font-family:Consolas,Monaco,Lucida Console,monospace;}h1{border-bottom:1px solid #c0c0c0;margin-bottom:10px;padding-bottom:10px;white-space:nowrap}table{border-collapse:collapse}tr.header{font-weight:bold}td.detailsColumn{-webkit-padding-start:2em;text-align:end;white-space:nowrap}a.icon{-webkit-padding-start:1.5em;text-decoration:none}a.icon:hover{text-decoration:underline}a.file{background:url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAIAAACQkWg2AAAABnRSTlMAAAAAAABupgeRAAABHUlEQVR42o2RMW7DIBiF3498iHRJD5JKHurL+CRVBp+i2T16tTynF2gO0KSb5ZrBBl4HHDBuK/WXACH4eO9/CAAAbdvijzLGNE1TVZXfZuHg6XCAQESAZXbOKaXO57eiKG6ft9PrKQIkCQqFoIiQFBGlFIB5nvM8t9aOX2Nd18oDzjnPgCDpn/BH4zh2XZdlWVmWiUK4IgCBoFMUz9eP6zRN75cLgEQhcmTQIbl72O0f9865qLAAsURAAgKBJKEtgLXWvyjLuFsThCSstb8rBCaAQhDYWgIZ7myM+TUBjDHrHlZcbMYYk34cN0YSLcgS+wL0fe9TXDMbY33fR2AYBvyQ8L0Gk8MwREBrTfKe4TpTzwhArXWi8HI84h/1DfwI5mhxJamFAAAAAElFTkSuQmCC ") left top no-repeat}a.dir{background:url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAd5JREFUeNqMU79rFUEQ/vbuodFEEkzAImBpkUabFP4ldpaJhZXYm/RiZWsv/hkWFglBUyTIgyAIIfgIRjHv3r39MePM7N3LcbxAFvZ2b2bn22/mm3XMjF+HL3YW7q28YSIw8mBKoBihhhgCsoORot9d3/ywg3YowMXwNde/PzGnk2vn6PitrT+/PGeNaecg4+qNY3D43vy16A5wDDd4Aqg/ngmrjl/GoN0U5V1QquHQG3q+TPDVhVwyBffcmQGJmSVfyZk7R3SngI4JKfwDJ2+05zIg8gbiereTZRHhJ5KCMOwDFLjhoBTn2g0ghagfKeIYJDPFyibJVBtTREwq60SpYvh5++PpwatHsxSm9QRLSQpEVSd7/TYJUb49TX7gztpjjEffnoVw66+Ytovs14Yp7HaKmUXeX9rKUoMoLNW3srqI5fWn8JejrVkK0QcrkFLOgS39yoKUQe292WJ1guUHG8K2o8K00oO1BTvXoW4yasclUTgZYJY9aFNfAThX5CZRmczAV52oAPoupHhWRIUUAOoyUIlYVaAa/VbLbyiZUiyFbjQFNwiZQSGl4IDy9sO5Wrty0QLKhdZPxmgGcDo8ejn+c/6eiK9poz15Kw7Dr/vN/z6W7q++091/AQYA5mZ8GYJ9K0AAAAAASUVORK5CYII= ") left top no-repeat}a.up{background:url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAmlJREFUeNpsU0toU0EUPfPysx/tTxuDH9SCWhUDooIbd7oRUUTMouqi2iIoCO6lceHWhegy4EJFinWjrlQUpVm0IIoFpVDEIthm0dpikpf3ZuZ6Z94nrXhhMjM3c8895977BBHB2PznK8WPtDgyWH5q77cPH8PpdXuhpQT4ifR9u5sfJb1bmw6VivahATDrxcRZ2njfoaMv+2j7mLDn93MPiNRMvGbL18L9IpF8h9/TN+EYkMffSiOXJ5+hkD+PdqcLpICWHOHc2CC+LEyA/K+cKQMnlQHJX8wqYG3MAJy88Wa4OLDvEqAEOpJd0LxHIMdHBziowSwVlF8D6QaicK01krw/JynwcKoEwZczewroTvZirlKJs5CqQ5CG8pb57FnJUA0LYCXMX5fibd+p8LWDDemcPZbzQyjvH+Ki1TlIciElA7ghwLKV4kRZstt2sANWRjYTAGzuP2hXZFpJ/GsxgGJ0ox1aoFWsDXyyxqCs26+ydmagFN/rRjymJ1898bzGzmQE0HCZpmk5A0RFIv8Pn0WYPsiu6t/Rsj6PauVTwffTSzGAGZhUG2F06hEc9ibS7OPMNp6ErYFlKavo7MkhmTqCxZ/jwzGA9Hx82H2BZSw1NTN9Gx8ycHkajU/7M+jInsDC7DiaEmo1bNl1AMr9ASFgqVu9MCTIzoGUimXVAnnaN0PdBBDCCYbEtMk6wkpQwIG0sn0PQIUF4GsTwLSIFKNqF6DVrQq+IWVrQDxAYQC/1SsYOI4pOxKZrfifiUSbDUisif7XlpGIPufXd/uvdvZm760M0no1FZcnrzUdjw7au3vu/BVgAFLXeuTxhTXVAAAAAElFTkSuQmCC ") left top no-repeat}html[dir=rtl] a{background-position-x:right}#listingParsingErrorBox{border:1px solid black;background:#fae691;padding:10px;display:none}</style>
<title id=title></title>
</head>
<body>
<span id="parentDirText" style="display:none" i18n-content="parentDirText">[上级目录]</span>
<table id="table">
<tbody><tr class="header">
<td i18n-content="headerName">名称</td>
<td class="detailsColumn" i18n-content="headerSize">大小</td>
<td class="detailsColumn" i18n-content="headerDateModified">修改日期</td>
</tr></tbody>
</table>
</body>
</html>
<script>var loadTimeData;function LoadTimeData(){}(function(){LoadTimeData.prototype={set datafunction(c){a(!this.data_,"Re-setting data.");this.data_=c},createJsEvalContext:function(){return new JsEvalContext(this.data_)},valueExists:function(c){return c in this.data_},getValue:function(d){a(this.data_,"No data. Did you remember to include strings.js?");var c=this.data_[d];a(typeof c!="undefined","Could not find value for "+d);return c},getString:function(d){var c=this.getValue(d);b(d,c,"string");return(c)},getStringF:function(f,d){var c=this.getString(f);if(!c){return""}var e=arguments;return c.replace(/\$[$1-9]/g,function(g){return g=="$$"?"$":e[g[1]]})},getBoolean:function(d){var c=this.getValue(d);b(d,c,"boolean");return(c)},getInteger:function(d){var c=this.getValue(d);b(d,c,"number");a(c==Math.floor(c),"Number isn't integer: "+c);return(c)},overrideValues:function(d){a(typeof d=="object","Replacements must be a dictionary object.");for(var c in d){this.data_[c]=d[c]}}};function a(d,c){if(!d){console.error("Unexpected condition on "+document.location.href+": "+c)}}function b(e,d,c){a(typeof d==c,"["+d+"] ("+e+") is not a "+c)}a(!loadTimeData,"should only include this file once");loadTimeData=new LoadTimeData})();</script><script>loadTimeData.data={header:"LOCATION 的索引",headerDateModified:"修改日期",headerName:"名称",headerSize:"大小",listingParsingErrorBoxText:'糟糕！Google Chrome无法解读服务器所发送的数据。请\u003Ca href="http://code.google.com/p/chromium/issues/entry">报告错误\u003C/a>，并附上\u003Ca href="LOCATION">原始列表\u003C/a>。',parentDirText:"[上级目录]",textdirection:"ltr"};</script><script>var i18nTemplate=(function(){var b={"i18n-content":function(f,e,g){f.textContent=g.getString(e)},"i18n-options":function(e,g,h){var f=h.getValue(g);f.forEach(function(j){var i=typeof j=="string"?new Option(j):new Option(j[1],j[0]);e.appendChild(i)})},"i18n-values":function(e,f,h){var g=f.replace(/\s/g,"").split(/;/);g.forEach(function(l){if(!l){return}var i=l.match(/^([^:]+):(.+)$/);if(!i){throw new Error("malformed i18n-values: "+f)}var n=i[1];var j=i[2];var m=h.getValue(j);if(n[0]=="."){var o=n.slice(1).split(".");var k=e;while(k&&o.length>1){k=k[o.shift()]}if(k){k[o]=m;if(o=="innerHTML"){d(e,h)}}}else{e.setAttribute(n,(m))}})}};var c=Object.keys(b);var a="["+c.join("],[")+"]";function d(l,n){var m=l.querySelectorAll(a);for(var h,g=0;h=m[g];g++){for(var f=0;f<c.length;f++){var e=c[f];var k=h.getAttribute(e);if(k!=null){b[e](h,k,n)}}}}return{process:d}}());i18nTemplate.process(document,loadTimeData);</script>
'''

doc_root = os.getcwd()
port = 80
addr = '0.0.0.0'
fsencoding = sys.getfilesystemencoding()
decoding = (fsencoding.lower() != 'utf-8')

def size_str(bytes):
    bytes = float(bytes)
    if bytes >= 1099511627776:
        terabytes = bytes / 1099511627776
        size = '%.2f T' % terabytes
    elif bytes >= 1073741824:
        gigabytes = bytes / 1073741824
        size = '%.2f G' % gigabytes
    elif bytes >= 1048576:
        megabytes = bytes / 1048576
        size = '%.2f M' % megabytes
    elif bytes >= 1024:
        kilobytes = bytes / 1024
        size = '%.2f K' % kilobytes
    else:
        size = '%d B' % int(bytes)
    return size

def time_str(mtime):
    return datetime.datetime.fromtimestamp(mtime).strftime('%Y-%m-%d %H:%M:%S')

def http_server(sock, addr):
    """Handler for http server requests."""
    try:
        msg = receive_message(sock)
        uri = unquote(parse_request(msg))

        filepath = os.path.join(doc_root, uri)
        if os.path.isfile(filepath):
            size = os.path.getsize(filepath)
            sock.send(build_header(guess_type(filepath)[0], os.path.getsize(filepath)))
            with open(filepath, 'rb') as infile:
                data = infile.read(4096)
                while data:
                    sock.sendall(data)
                    data = infile.read(4096)

        elif os.path.isdir(filepath):
            contents = os.listdir(filepath)
            rep = html_str
            rep += '<script>start("%s");</script>' % uri
            if uri != '/' :
                rep += '<script>addRow("..","..",1,"0 B","");</script>'
            for i in range(len(contents)):
                try:
                    p = os.path.join(filepath, contents[i])
                    n = contents[i]
                    if decoding:
                        n = contents[i].decode(fsencoding).encode('utf-8')
                    if os.path.isdir(p):
                        rep += '<script>addRow("%s","%s",1,"0 B","%s");</script>\n' % (n, quote(contents[i]), time_str(os.path.getmtime(p)))
                except:
                    pass
            for i in range(len(contents)):
                try:
                    p = os.path.join(filepath, contents[i])
                    n = contents[i]
                    if decoding:
                        n = contents[i].decode(fsencoding).encode('utf-8')
                    if os.path.isfile(p):
                        rep += '<script>addRow("%s","%s",0,"%s","%s");</script>\n' % (n, quote(contents[i]), size_str(os.path.getsize(p)), time_str(os.path.getmtime(p)))
                except:
                    pass
            sock.sendall(build_header('text/html', len(rep)) + rep)
        else:
            #If what we received was not a file or a directory, raise an Error404.
            raise Error404("404: File not found.")

    except (Error404, Error405) as e:
        rep = e.message
        sock.sendall(build_header('text/plain', len(rep), e.code) + rep)

    except:
        rep = "500 Internal Server Error"
        sock.sendall(build_header('text/plain', len(rep), '500') + rep)
        pass
    finally:
       sock.shutdown(socket.SHUT_WR)
       sock.close()


def receive_message(conn, buffsize=4096):
    """When a connection is received by the http_server, this function
    pieces together the message received and returns it.
    """

    msg = ''
    while True:
        msg_part = conn.recv(buffsize)
        msg += msg_part
        if len(msg_part) < buffsize:
            break

    conn.shutdown(socket.SHUT_RD)
    return msg


def parse_request(request):
    first_rn = request.find('\r\n')
    first_line = request[:first_rn]
    if first_line.split()[0] == 'GET':
        uri = first_line.split()[1]
        return uri
    else:
        raise Error405("405: Method not allowed. Only GET is allowed.")


def build_header(mimetype, bytelen, code="200 OK"):
    """Build a response with the specified code and content."""

    resp_list = []
    resp_list.append('HTTP/1.1 %s' % code)
    resp_list.append('Content-Type: %s; char=utf-8' % mimetype)
    resp_list.append('Content-Length: %s' % str(bytelen))
    resp_list.append('\r\n')
    resp = '\r\n'.join(resp_list)
    return resp



class Error404(BaseException):
    """Exception raised when a file specified by a URI does not exist."""
    code = '404'


class Error405(BaseException):
    """Exception raised when a request is made with a method other than GET."""
    code = '405'

def parse_cmd():
    global doc_root, port, addr
    arvc = len(sys.argv)
    if arvc > 4 or arvc == 1 or (arvc == 2 and (sys.argv[1] == '-h' or sys.argv[1] == '--help')):
        print 'python static_server.py address port doc_root'
        print '  For IPv4, try:\n'
        print '    python static_server.py 0.0.0.0 80 .\n'
        print '  For IPv6, try:\n'
        print '    python static_server.py 0::0 80 .\n'
        exit(1)
    if arvc > 3:
        doc_root = sys.argv[3]
    if arvc > 2:
        port = int(sys.argv[2])
    if arvc > 1:
        addr = sys.argv[1]
    print doc_root

def main():
    parse_cmd()
    patch_all()
    server = StreamServer((addr, port), http_server)
    print('Starting server on %s:%d, doc_root:%s' % (addr, port, doc_root))
    server.serve_forever()

if __name__ == '__main__':
    main()
