var os = require('os'),
  http = require("http"),
  url = require("url"),
  qs = require("querystring"),
  path = require("path"),
  fs = require("fs"),
  mime = require('mime'),
  moment = require('moment'),
  ct = 0,
  port = parseInt(process.argv[2], 10) || 8888, // 启动端口 port
  tp = process.argv[3] || 'ftp', // 输出类型 type
  wp = process.argv[4] || process.cwd(); // 指定目录 dir

if ("-v" == process.argv[2] || "--version" == process.argv[2]) {
  console.log("v0.0.6");
  return;
}

// var contentType = {
//   ftp: 'application/octet-stream',
//   http: {
//     html: 'text/html',
//     xml: 'text/xml',
//     js: 'text/javascript',
//     css: 'text/css',
//     gif: 'image/gif',
//     jpg: 'image/jpg',
//     jpeg: 'image/jpeg',
//     svg: 'image/svg+xml',
//     png: 'image/png'
//   }
// };
// var contType = contentType[tp] ? contentType[tp] : false;

var html_str = "<!DOCTYPE html>\n" +
  "<html i18n-values=dir:textdirection;lang:language>\n" +
  "<head>\n" +
  "<meta charset=utf-8>\n" +
  "<script>function addRow(a,b,c,j,d){if(a==\".\"){return}var g=\"\"+document.location;if(g.substr(-1)!==\"/\"){g+=\"/\"}var h=document.getElementById(\"table\");var i=document.createElement(\"tr\");var e=document.createElement(\"td\");var f=document.createElement(\"a\");f.className=c?\"icon dir\":\"icon file\";if(a==\"..\"){f.href=g+\"..\";f.innerText=document.getElementById(\"parentDirText\").innerText;f.className=\"icon up\";j=\"\";d=\"\"}else{if(c){a=a+\"/\";b=b+\"/\";j=\"\"}else{f.draggable=\"true\";f.addEventListener(\"dragstart\",onDragStart,false)}f.innerText=a;f.href=g+b}e.appendChild(f);i.appendChild(e);i.appendChild(createCell(j));i.appendChild(createCell(d));h.appendChild(i)}function onDragStart(d){var c=d.srcElement;var b=c.innerText.replace(\":\",\"\");var a=\"application/octet-stream:\"+b+\":\"+c.href;d.dataTransfer.setData(\"DownloadURL\",a);d.dataTransfer.effectAllowed=\"copy\"}function createCell(b){var a=document.createElement(\"td\");a.setAttribute(\"class\",\"detailsColumn\");a.innerText=b;return a}function start(a){var b=document.getElementById(\"header\");b.innerText=b.innerText.replace(\"LOCATION\",a);document.getElementById(\"title\").innerText=b.innerText}function onListingParsingError(){var a=document.getElementById(\"listingParsingErrorBox\");a.innerHTML=a.innerHTML.replace(\"LOCATION\",encodeURI(document.location)+\"?raw\");a.style.display=\"block\"}</script>\n" +
  "<style>body{font-family:Consolas,Monaco,Lucida Console,monospace;}h1{border-bottom:1px solid #c0c0c0;margin-bottom:10px;padding-bottom:10px;white-space:nowrap}table{border-collapse:collapse}tr.header{font-weight:bold}td.detailsColumn{-webkit-padding-start:2em;text-align:end;white-space:nowrap}a.icon{-webkit-padding-start:1.5em;text-decoration:none}a.icon:hover{text-decoration:underline}a.file{background:url(\"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAIAAACQkWg2AAAABnRSTlMAAAAAAABupgeRAAABHUlEQVR42o2RMW7DIBiF3498iHRJD5JKHurL+CRVBp+i2T16tTynF2gO0KSb5ZrBBl4HHDBuK/WXACH4eO9/CAAAbdvijzLGNE1TVZXfZuHg6XCAQESAZXbOKaXO57eiKG6ft9PrKQIkCQqFoIiQFBGlFIB5nvM8t9aOX2Nd18oDzjnPgCDpn/BH4zh2XZdlWVmWiUK4IgCBoFMUz9eP6zRN75cLgEQhcmTQIbl72O0f9865qLAAsURAAgKBJKEtgLXWvyjLuFsThCSstb8rBCaAQhDYWgIZ7myM+TUBjDHrHlZcbMYYk34cN0YSLcgS+wL0fe9TXDMbY33fR2AYBvyQ8L0Gk8MwREBrTfKe4TpTzwhArXWi8HI84h/1DfwI5mhxJamFAAAAAElFTkSuQmCC \") left top no-repeat}a.dir{background:url(\"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAd5JREFUeNqMU79rFUEQ/vbuodFEEkzAImBpkUabFP4ldpaJhZXYm/RiZWsv/hkWFglBUyTIgyAIIfgIRjHv3r39MePM7N3LcbxAFvZ2b2bn22/mm3XMjF+HL3YW7q28YSIw8mBKoBihhhgCsoORot9d3/ywg3YowMXwNde/PzGnk2vn6PitrT+/PGeNaecg4+qNY3D43vy16A5wDDd4Aqg/ngmrjl/GoN0U5V1QquHQG3q+TPDVhVwyBffcmQGJmSVfyZk7R3SngI4JKfwDJ2+05zIg8gbiereTZRHhJ5KCMOwDFLjhoBTn2g0ghagfKeIYJDPFyibJVBtTREwq60SpYvh5++PpwatHsxSm9QRLSQpEVSd7/TYJUb49TX7gztpjjEffnoVw66+Ytovs14Yp7HaKmUXeX9rKUoMoLNW3srqI5fWn8JejrVkK0QcrkFLOgS39yoKUQe292WJ1guUHG8K2o8K00oO1BTvXoW4yasclUTgZYJY9aFNfAThX5CZRmczAV52oAPoupHhWRIUUAOoyUIlYVaAa/VbLbyiZUiyFbjQFNwiZQSGl4IDy9sO5Wrty0QLKhdZPxmgGcDo8ejn+c/6eiK9poz15Kw7Dr/vN/z6W7q++091/AQYA5mZ8GYJ9K0AAAAAASUVORK5CYII= \") left top no-repeat}a.up{background:url(\"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAmlJREFUeNpsU0toU0EUPfPysx/tTxuDH9SCWhUDooIbd7oRUUTMouqi2iIoCO6lceHWhegy4EJFinWjrlQUpVm0IIoFpVDEIthm0dpikpf3ZuZ6Z94nrXhhMjM3c8895977BBHB2PznK8WPtDgyWH5q77cPH8PpdXuhpQT4ifR9u5sfJb1bmw6VivahATDrxcRZ2njfoaMv+2j7mLDn93MPiNRMvGbL18L9IpF8h9/TN+EYkMffSiOXJ5+hkD+PdqcLpICWHOHc2CC+LEyA/K+cKQMnlQHJX8wqYG3MAJy88Wa4OLDvEqAEOpJd0LxHIMdHBziowSwVlF8D6QaicK01krw/JynwcKoEwZczewroTvZirlKJs5CqQ5CG8pb57FnJUA0LYCXMX5fibd+p8LWDDemcPZbzQyjvH+Ki1TlIciElA7ghwLKV4kRZstt2sANWRjYTAGzuP2hXZFpJ/GsxgGJ0ox1aoFWsDXyyxqCs26+ydmagFN/rRjymJ1898bzGzmQE0HCZpmk5A0RFIv8Pn0WYPsiu6t/Rsj6PauVTwffTSzGAGZhUG2F06hEc9ibS7OPMNp6ErYFlKavo7MkhmTqCxZ/jwzGA9Hx82H2BZSw1NTN9Gx8ycHkajU/7M+jInsDC7DiaEmo1bNl1AMr9ASFgqVu9MCTIzoGUimXVAnnaN0PdBBDCCYbEtMk6wkpQwIG0sn0PQIUF4GsTwLSIFKNqF6DVrQq+IWVrQDxAYQC/1SsYOI4pOxKZrfifiUSbDUisif7XlpGIPufXd/uvdvZm760M0no1FZcnrzUdjw7au3vu/BVgAFLXeuTxhTXVAAAAAElFTkSuQmCC \") left top no-repeat}html[dir=rtl] a{background-position-x:right}#listingParsingErrorBox{border:1px solid black;background:#fae691;padding:10px;display:none}</style>\n" +
  "<title id=title></title>\n" +
  "</head>\n" +
  "<body>\n" +
  "<span id=\"parentDirText\" style=\"display:none\" i18n-content=\"parentDirText\">[上级目录]</span>\n" +
  "<table id=\"table\">\n" +
  "<tbody><tr class=\"header\">\n" +
  "<td i18n-content=\"headerName\">名称</td>\n" +
  "<td class=\"detailsColumn\" i18n-content=\"headerSize\">大小</td>\n" +
  "<td class=\"detailsColumn\" i18n-content=\"headerDateModified\">修改日期</td>\n" +
  "</tr></tbody>\n" +
  "</table>\n" +
  "</body>\n" +
  "</html>\n" +
  "<script>var loadTimeData;function LoadTimeData(){}(function(){LoadTimeData.prototype={set datafunction(c){a(!this.data_,\"Re-setting data.\");this.data_=c},createJsEvalContext:function(){return new JsEvalContext(this.data_)},valueExists:function(c){return c in this.data_},getValue:function(d){a(this.data_,\"No data. Did you remember to include strings.js?\");var c=this.data_[d];a(typeof c!=\"undefined\",\"Could not find value for \"+d);return c},getString:function(d){var c=this.getValue(d);b(d,c,\"string\");return(c)},getStringF:function(f,d){var c=this.getString(f);if(!c){return\"\"}var e=arguments;return c.replace(/\\$[$1-9]/g,function(g){return g==\"$$\"?\"$\":e[g[1]]})},getBoolean:function(d){var c=this.getValue(d);b(d,c,\"boolean\");return(c)},getInteger:function(d){var c=this.getValue(d);b(d,c,\"number\");a(c==Math.floor(c),\"Number isn't integer: \"+c);return(c)},overrideValues:function(d){a(typeof d==\"object\",\"Replacements must be a dictionary object.\");for(var c in d){this.data_[c]=d[c]}}};function a(d,c){if(!d){console.error(\"Unexpected condition on \"+document.location.href+\": \"+c)}}function b(e,d,c){a(typeof d==c,\"[\"+d+\"] (\"+e+\") is not a \"+c)}a(!loadTimeData,\"should only include this file once\");loadTimeData=new LoadTimeData})();</script><script>loadTimeData.data={header:\"LOCATION 的索引\",headerDateModified:\"修改日期\",headerName:\"名称\",headerSize:\"大小\",listingParsingErrorBoxText:'糟糕！Google Chrome无法解读服务器所发送的数据。请\\u003Ca href=\"http://code.google.com/p/chromium/issues/entry\">报告错误\\u003C/a>，并附上\\u003Ca href=\"LOCATION\">原始列表\\u003C/a>。',parentDirText:\"[上级目录]\",textdirection:\"ltr\"};</script><script>var i18nTemplate=(function(){var b={\"i18n-content\":function(f,e,g){f.textContent=g.getString(e)},\"i18n-options\":function(e,g,h){var f=h.getValue(g);f.forEach(function(j){var i=typeof j==\"string\"?new Option(j):new Option(j[1],j[0]);e.appendChild(i)})},\"i18n-values\":function(e,f,h){var g=f.replace(/\\s/g,\"\").split(/;/);g.forEach(function(l){if(!l){return}var i=l.match(/^([^:]+):(.+)$/);if(!i){throw new Error(\"malformed i18n-values: \"+f)}var n=i[1];var j=i[2];var m=h.getValue(j);if(n[0]==\".\"){var o=n.slice(1).split(\".\");var k=e;while(k&&o.length>1){k=k[o.shift()]}if(k){k[o]=m;if(o==\"innerHTML\"){d(e,h)}}}else{e.setAttribute(n,(m))}})}};var c=Object.keys(b);var a=\"[\"+c.join(\"],[\")+\"]\";function d(l,n){var m=l.querySelectorAll(a);for(var h,g=0;h=m[g];g++){for(var f=0;f<c.length;f++){var e=c[f];var k=h.getAttribute(e);if(k!=null){b[e](h,k,n)}}}}return{process:d}}());i18nTemplate.process(document,loadTimeData);</script>\n";

function size_str(size) {

  var c =  ["B", "KB", "MB", "GB"];
  var u = [ 1, 1 << 10, 1 << 20, 1 << 30 ];
  var i;
  var test = size;
  for (i = 0; i < 4; i++)
  {
    test = test / 1024;
    if (test <= 1) {
      break;
    }
  }

  if (i == 4) {
    i--;
  }

  if (i > 0) {
    return '' + (size * 1.00 / u[i]).toFixed(2) + ' ' + c[i];
  } else {
    return '' + size + ' B';
  }

}

var hftp = function(request, response) {

  request.setEncoding('utf8'); //请求编码

  var uri = qs.unescape(url.parse(request.url).pathname),
    filename = path.join(wp, uri);

  fs.exists(filename, function(exists) {
    if (!exists) {
      response.writeHead(404, {
        "Content-Type": "text/plain"
      });
      response.write("404 Not Found\n");
      response.end();
      return;
    }

    fs.stat(filename, function(err, stats) {
      if (err) {
        response.writeHead(500, {
          "Content-Type": "text/plain"
        });
        response.write(err + "\n");
        response.end();
        return;

      } else {
        if (stats.isDirectory()) {
          var files = fs.readdirSync(filename);
          if (files) {
            if (uri.lastIndexOf('/') == uri.length - 1) {
              uri = uri.substring(0, uri.length - 1);
            }
            var dir = '/';
            if (uri.lastIndexOf('/') > 0) {
              dir = uri.substring(0, uri.lastIndexOf('/'));
            }
            var all = html_str;
            all += '<script>start("' + uri + '");</script>';
            if (uri.length > 0) {
              all += '<script>addRow("..","..",1,"0 B","");</script>';
            }
            uri += '/';
            for (var i in files) {
              try {
                var is = fs.lstatSync(filename + files[i]);
                if (is.isDirectory()) {
                  all += '<script>addRow("' + files[i] + '","' + qs.escape(files[i]) + '", 1, "0 B", "' + moment(is.mtime).format("YYYY/MM/DD hh:mm:ss") + '");</script>';
                }
              } catch (ex) {}
            }
            for (var i in files) {
              try {
                var is = fs.lstatSync(filename + files[i]);
                if (is.isFile()) {
                  all += '<script>addRow("' + files[i] + '","' + qs.escape(files[i]) + '", 0, "' + size_str(is.size) + '", "' + moment(is.mtime).format("YYYY/MM/DD hh:mm:ss") + '");</script>';
                }
              } catch (ex) {}
            }
            response.writeHead(202, {
              "Content-Type": "text/html; charset=UTF-8"
            });
            response.write(all);
            response.end();
          }
        } else {
          ct ++;
          // var t = contType ? contType[filename.replace(/.*\.(\w+)/,'$1').toLowerCase()] : false;
          var hd = {
            'Content-Type': tp === 'http' ? mime.lookup(filename) : 'application/octet-stream',
            'Content-Length': stats.size
          };
          response.writeHead(200, hd);

          fs.createReadStream(filename).pipe(response);
        }

      }
    });

  });
};

var server = http.createServer(hftp);
server.on('error', function(e) {
  if (e.code == 'EADDRINUSE') {
    console.log('Port has been used. 端口已被使用 :' + port);
    server.close();
  } else {
    throw e;
  }
});
server.listen(port);

console.log("Static file server running at\n  => http://localhost:" + port + " type=" + tp + ", workspace=" + wp + "/\nCTRL + C to shutdown");