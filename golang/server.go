package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var confPort string
var confRoot string
var confIp string

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func sizeString(size int64) string {
	us := "BKMGT"
	r := float64(size)
	i := 0
	for ; r > 1000 && i < 5; i++ {
		r = r / 1024
	}
	return fmt.Sprintf("%.2f %c", r, us[i])
}

func listDir(rw http.ResponseWriter, req *http.Request, fpath string) {
	fmt.Fprintln(rw, html)
	fmt.Fprintf(rw, "<script>start(\"【%s】\");</script>\n", req.Host+req.URL.Path)
	root, _ := filepath.Abs(confRoot)
	flist, _ := ioutil.ReadDir(fpath)
	var files, dirs []os.FileInfo
	for _, f := range flist {
		if f.IsDir() {
			dirs = append(dirs, f)
		} else {
			files = append(files, f)
		}
	}
	if root != fpath {
		fmt.Fprintf(rw, "<script>addRow(\"..\",\"..\",1,0,\"0 B\", 0,\"\");</script>\n")
	}
	for _, item := range dirs {
		encoded := strings.Replace(url.QueryEscape(item.Name()), "+", "%20", -1)
		fmt.Fprintf(rw, "<script>addRow(\"%s\",\"%s\",1,0,\"0 B\", %d,\"%s\");</script>\n",
			item.Name(), encoded, item.ModTime().Unix(), item.ModTime().Format("2006-01-02 15:04:05"))
	}
	for _, item := range files {
		encoded := strings.Replace(url.QueryEscape(item.Name()), "+", "%20", -1)
		fmt.Fprintf(rw, "<script>addRow(\"%s\",\"%s\",0,%d,\"%s\", %d,\"%s\");</script>\n",
			item.Name(), encoded, item.Size(), sizeString(item.Size()), item.ModTime().Unix(), item.ModTime().Format("2006-01-02 15:04:05"))
	}
}

func rootHandler(rw http.ResponseWriter, req *http.Request) {
	fpath, _ := filepath.Abs(path.Join(confRoot, req.URL.Path))
	finfo, err := os.Stat(fpath)
	log.Println(err)
	if err != nil && os.IsNotExist(err) {
		http.Error(rw, "404", http.StatusNotFound)
	} else {
		if finfo.IsDir() {
			listDir(rw, req, fpath)
		} else {
			http.ServeFile(rw, req, fpath)
		}
	}
}

func main() {
	flag.StringVar(&confIp, "ip", "0.0.0.0", "listening ip")
	flag.StringVar(&confPort, "port", "80", "listening port")
	flag.StringVar(&confRoot, "root", ".", "root directory")
	flag.Parse()

	fmt.Printf("Serving HTTP on %s port %s, root: %s", confIp, confPort, confRoot)
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(confIp+":"+confPort, nil))
}

// codes from chrome local file browser
var html string = `<!DOCTYPE html>

<html i18n-values="dir:textdirection;lang:language">

<head>
<meta charset="utf-8">
<meta name="google" value="notranslate">

<script>
function addRow(name, url, isdir,
    size, size_string, date_modified, date_modified_string) {
  if (name == ".")
    return;

  var root = document.location.pathname;
  if (root.substr(-1) !== "/")
    root += "/";

  var tbody = document.getElementById("tbody");
  var row = document.createElement("tr");
  var file_cell = document.createElement("td");
  var link = document.createElement("a");

  link.className = isdir ? "icon dir" : "icon file";

  if (name == "..") {
    link.href = root + "..";
    link.innerText = document.getElementById("parentDirText").innerText;
    link.className = "icon up";
    size = 0;
    size_string = "";
    date_modified = 0;
    date_modified_string = "";
  } else {
    if (isdir) {
      name = name + "/";
      url = url + "/";
      size = 0;
      size_string = "";
    } else {
      link.draggable = "true";
      link.addEventListener("dragstart", onDragStart, false);
    }
    link.innerText = name;
    link.href = root + url;
  }
  file_cell.dataset.value = name;
  file_cell.appendChild(link);

  row.appendChild(file_cell);
  row.appendChild(createCell(size, size_string));
  row.appendChild(createCell(date_modified, date_modified_string));

  tbody.appendChild(row);
}

function onDragStart(e) {
  var el = e.srcElement;
  var name = el.innerText.replace(":", "");
  var download_url_data = "application/octet-stream:" + name + ":" + el.href;
  e.dataTransfer.setData("DownloadURL", download_url_data);
  e.dataTransfer.effectAllowed = "copy";
}

function createCell(value, text) {
  var cell = document.createElement("td");
  cell.setAttribute("class", "detailsColumn");
  cell.dataset.value = value;
  cell.innerText = text;
  return cell;
}

function start(location) {
  var header = document.getElementById("header");
  header.innerText = header.innerText.replace("LOCATION", location);

  document.getElementById("title").innerText = header.innerText;
}

function onListingParsingError() {
  var box = document.getElementById("listingParsingErrorBox");
  box.innerHTML = box.innerHTML.replace("LOCATION", encodeURI(document.location)
      + "?raw");
  box.style.display = "block";
}

function sortTable(column) {
  var theader = document.getElementById("theader");
  var oldOrder = theader.cells[column].dataset.order || '1';
  oldOrder = parseInt(oldOrder, 10)
  var newOrder = 0 - oldOrder;
  theader.cells[column].dataset.order = newOrder;

  var tbody = document.getElementById("tbody");
  var rows = tbody.rows;
  var list = [], i;
  for (i = 0; i < rows.length; i++) {
    list.push(rows[i]);
  }

  list.sort(function(row1, row2) {
    var a = row1.cells[column].dataset.value;
    var b = row2.cells[column].dataset.value;
    if (column) {
      a = parseInt(a, 10);
      b = parseInt(b, 10);
      return a > b ? newOrder : a < b ? oldOrder : 0;
    }

    // Column 0 is text.
    // Also the parent directory should always be sorted at one of the ends.
    if (b == ".." | a > b) {
      return newOrder;
    } else if (a == ".." | a < b) {
      return oldOrder;
    } else {
      return 0;
    }
  });

  // Appending an existing child again just moves it.
  for (i = 0; i < list.length; i++) {
    tbody.appendChild(list[i]);
  }
}
</script>

<style>
  body {
    font-family:Consolas,Monaco,Lucida Console,monospace;
  }
  h1 {
    border-bottom: 1px solid #c0c0c0;
    margin-bottom: 10px;
    padding-bottom: 10px;
    white-space: nowrap;
  }

  table {
    border-collapse: collapse;
  }

  th {
    cursor: pointer;
  }

  td.detailsColumn {
    -webkit-padding-start: 2em;
    text-align: end;
    white-space: nowrap;
  }

  a.icon {
    -webkit-padding-start: 1.5em;
    text-decoration: none;
  }

  a.icon:hover {
    text-decoration: underline;
  }

  a.file {
    background : url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAIAAACQkWg2AAAABnRSTlMAAAAAAABupgeRAAABHUlEQVR42o2RMW7DIBiF3498iHRJD5JKHurL+CRVBp+i2T16tTynF2gO0KSb5ZrBBl4HHDBuK/WXACH4eO9/CAAAbdvijzLGNE1TVZXfZuHg6XCAQESAZXbOKaXO57eiKG6ft9PrKQIkCQqFoIiQFBGlFIB5nvM8t9aOX2Nd18oDzjnPgCDpn/BH4zh2XZdlWVmWiUK4IgCBoFMUz9eP6zRN75cLgEQhcmTQIbl72O0f9865qLAAsURAAgKBJKEtgLXWvyjLuFsThCSstb8rBCaAQhDYWgIZ7myM+TUBjDHrHlZcbMYYk34cN0YSLcgS+wL0fe9TXDMbY33fR2AYBvyQ8L0Gk8MwREBrTfKe4TpTzwhArXWi8HI84h/1DfwI5mhxJamFAAAAAElFTkSuQmCC ") left top no-repeat;
  }

  a.dir {
    background : url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAd5JREFUeNqMU79rFUEQ/vbuodFEEkzAImBpkUabFP4ldpaJhZXYm/RiZWsv/hkWFglBUyTIgyAIIfgIRjHv3r39MePM7N3LcbxAFvZ2b2bn22/mm3XMjF+HL3YW7q28YSIw8mBKoBihhhgCsoORot9d3/ywg3YowMXwNde/PzGnk2vn6PitrT+/PGeNaecg4+qNY3D43vy16A5wDDd4Aqg/ngmrjl/GoN0U5V1QquHQG3q+TPDVhVwyBffcmQGJmSVfyZk7R3SngI4JKfwDJ2+05zIg8gbiereTZRHhJ5KCMOwDFLjhoBTn2g0ghagfKeIYJDPFyibJVBtTREwq60SpYvh5++PpwatHsxSm9QRLSQpEVSd7/TYJUb49TX7gztpjjEffnoVw66+Ytovs14Yp7HaKmUXeX9rKUoMoLNW3srqI5fWn8JejrVkK0QcrkFLOgS39yoKUQe292WJ1guUHG8K2o8K00oO1BTvXoW4yasclUTgZYJY9aFNfAThX5CZRmczAV52oAPoupHhWRIUUAOoyUIlYVaAa/VbLbyiZUiyFbjQFNwiZQSGl4IDy9sO5Wrty0QLKhdZPxmgGcDo8ejn+c/6eiK9poz15Kw7Dr/vN/z6W7q++091/AQYA5mZ8GYJ9K0AAAAAASUVORK5CYII= ") left top no-repeat;
  }

  a.up {
    background : url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAmlJREFUeNpsU0toU0EUPfPysx/tTxuDH9SCWhUDooIbd7oRUUTMouqi2iIoCO6lceHWhegy4EJFinWjrlQUpVm0IIoFpVDEIthm0dpikpf3ZuZ6Z94nrXhhMjM3c8895977BBHB2PznK8WPtDgyWH5q77cPH8PpdXuhpQT4ifR9u5sfJb1bmw6VivahATDrxcRZ2njfoaMv+2j7mLDn93MPiNRMvGbL18L9IpF8h9/TN+EYkMffSiOXJ5+hkD+PdqcLpICWHOHc2CC+LEyA/K+cKQMnlQHJX8wqYG3MAJy88Wa4OLDvEqAEOpJd0LxHIMdHBziowSwVlF8D6QaicK01krw/JynwcKoEwZczewroTvZirlKJs5CqQ5CG8pb57FnJUA0LYCXMX5fibd+p8LWDDemcPZbzQyjvH+Ki1TlIciElA7ghwLKV4kRZstt2sANWRjYTAGzuP2hXZFpJ/GsxgGJ0ox1aoFWsDXyyxqCs26+ydmagFN/rRjymJ1898bzGzmQE0HCZpmk5A0RFIv8Pn0WYPsiu6t/Rsj6PauVTwffTSzGAGZhUG2F06hEc9ibS7OPMNp6ErYFlKavo7MkhmTqCxZ/jwzGA9Hx82H2BZSw1NTN9Gx8ycHkajU/7M+jInsDC7DiaEmo1bNl1AMr9ASFgqVu9MCTIzoGUimXVAnnaN0PdBBDCCYbEtMk6wkpQwIG0sn0PQIUF4GsTwLSIFKNqF6DVrQq+IWVrQDxAYQC/1SsYOI4pOxKZrfifiUSbDUisif7XlpGIPufXd/uvdvZm760M0no1FZcnrzUdjw7au3vu/BVgAFLXeuTxhTXVAAAAAElFTkSuQmCC ") left top no-repeat;
  }

  html[dir=rtl] a {
    background-position-x: right;
  }

  #listingParsingErrorBox {
    border: 1px solid black;
    background: #fae691;
    padding: 10px;
    display: none;
  }
</style>

<title id="title"></title>

</head>

<body>

<div id="listingParsingErrorBox" i18n-values=".innerHTML:listingParsingErrorBoxText"></div>

<span id="parentDirText" style="display:none" i18n-content="parentDirText"></span>

<h1 id="header" i18n-content="header"></h1>

<table>
  <thead>
    <tr class="header" id="theader">
      <th i18n-content="headerName" onclick="javascript:sortTable(0);"></th>
      <th class="detailsColumn" i18n-content="headerSize" onclick="javascript:sortTable(1);"></th>
      <th class="detailsColumn" i18n-content="headerDateModified" onclick="javascript:sortTable(2);"></th>
    </tr>
  </thead>
  <tbody id="tbody">
  </tbody>
</table>

</body>

</html>
<script>// Copyright (c) 2012 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

/**
 * @fileoverview This file defines a singleton which provides access to all data
 * that is available as soon as the page's resources are loaded (before DOM
 * content has finished loading). This data includes both localized strings and
 * any data that is important to have ready from a very early stage (e.g. things
 * that must be displayed right away).
 */

/** @type {!LoadTimeData} */ var loadTimeData;

// Expose this type globally as a temporary work around until
// https://github.com/google/closure-compiler/issues/544 is fixed.
/** @constructor */
function LoadTimeData() {}

(function() {
  'use strict';

  LoadTimeData.prototype = {
    /**
     * Sets the backing object.
     *
     * Note that there is no getter for |data_| to discourage abuse of the form:
     *
     *     var value = loadTimeData.data()['key'];
     *
     * @param {Object} value The de-serialized page data.
     */
    set data(value) {
      expect(!this.data_, 'Re-setting data.');
      this.data_ = value;
    },

    /**
     * Returns a JsEvalContext for |data_|.
     * @returns {JsEvalContext}
     */
    createJsEvalContext: function() {
      return new JsEvalContext(this.data_);
    },

    /**
     * @param {string} id An ID of a value that might exist.
     * @return {boolean} True if |id| is a key in the dictionary.
     */
    valueExists: function(id) {
      return id in this.data_;
    },

    /**
     * Fetches a value, expecting that it exists.
     * @param {string} id The key that identifies the desired value.
     * @return {*} The corresponding value.
     */
    getValue: function(id) {
      expect(this.data_, 'No data. Did you remember to include strings.js?');
      var value = this.data_[id];
      expect(typeof value != 'undefined', 'Could not find value for ' + id);
      return value;
    },

    /**
     * As above, but also makes sure that the value is a string.
     * @param {string} id The key that identifies the desired string.
     * @return {string} The corresponding string value.
     */
    getString: function(id) {
      var value = this.getValue(id);
      expectIsType(id, value, 'string');
      return /** @type {string} */ (value);
    },

    /**
     * Returns a formatted localized string where $1 to $9 are replaced by the
     * second to the tenth argument.
     * @param {string} id The ID of the string we want.
     * @param {...(string|number)} var_args The extra values to include in the
     *     formatted output.
     * @return {string} The formatted string.
     */
    getStringF: function(id, var_args) {
      var value = this.getString(id);
      if (!value)
        return '';

      var varArgs = arguments;
      return value.replace(/\$[$1-9]/g, function(m) {
        return m == '$$' ? '$' : varArgs[m[1]];
      });
    },

    /**
     * As above, but also makes sure that the value is a boolean.
     * @param {string} id The key that identifies the desired boolean.
     * @return {boolean} The corresponding boolean value.
     */
    getBoolean: function(id) {
      var value = this.getValue(id);
      expectIsType(id, value, 'boolean');
      return /** @type {boolean} */ (value);
    },

    /**
     * As above, but also makes sure that the value is an integer.
     * @param {string} id The key that identifies the desired number.
     * @return {number} The corresponding number value.
     */
    getInteger: function(id) {
      var value = this.getValue(id);
      expectIsType(id, value, 'number');
      expect(value == Math.floor(value), 'Number isn\'t integer: ' + value);
      return /** @type {number} */ (value);
    },

    /**
     * Override values in loadTimeData with the values found in |replacements|.
     * @param {Object} replacements The dictionary object of keys to replace.
     */
    overrideValues: function(replacements) {
      expect(typeof replacements == 'object',
             'Replacements must be a dictionary object.');
      for (var key in replacements) {
        this.data_[key] = replacements[key];
      }
    }
  };

  /**
   * Checks condition, displays error message if expectation fails.
   * @param {*} condition The condition to check for truthiness.
   * @param {string} message The message to display if the check fails.
   */
  function expect(condition, message) {
    if (!condition) {
      console.error('Unexpected condition on ' + document.location.href + ': ' +
                    message);
    }
  }

  /**
   * Checks that the given value has the given type.
   * @param {string} id The id of the value (only used for error message).
   * @param {*} value The value to check the type on.
   * @param {string} type The type we expect |value| to be.
   */
  function expectIsType(id, value, type) {
    expect(typeof value == type, '[' + value + '] (' + id +
                                 ') is not a ' + type);
  }

  expect(!loadTimeData, 'should only include this file once');
  loadTimeData = new LoadTimeData;
})();
</script><script>loadTimeData.data = {"header":"LOCATION 的索引","headerDateModified":"修改日期","headerName":"名称","headerSize":"大小","listingParsingErrorBoxText":"糟糕！Google Chrome无法解读服务器所发送的数据。请\u003Ca href=\"http://code.google.com/p/chromium/issues/entry\">报告错误\u003C/a>，并附上\u003Ca href=\"LOCATION\">原始列表\u003C/a>。","parentDirText":"[上级目录]","textdirection":"ltr"};</script><script>// Copyright (c) 2012 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// // Copyright (c) 2012 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

/** @typedef {Document|DocumentFragment|Element} */
var ProcessingRoot;

/**
 * @fileoverview This is a simple template engine inspired by JsTemplates
 * optimized for i18n.
 *
 * It currently supports three handlers:
 *
 *   * i18n-content which sets the textContent of the element.
 *
 *     <span i18n-content="myContent"></span>
 *
 *   * i18n-options which generates <option> elements for a <select>.
 *
 *     <select i18n-options="myOptionList"></select>
 *
 *   * i18n-values is a list of attribute-value or property-value pairs.
 *     Properties are prefixed with a '.' and can contain nested properties.
 *
 *     <span i18n-values="title:myTitle;.style.fontSize:fontSize"></span>
 *
 * This file is a copy of i18n_template.js, with minor tweaks to support using
 * load_time_data.js. It should replace i18n_template.js eventually.
 */

var i18nTemplate = (function() {
  /**
   * This provides the handlers for the templating engine. The key is used as
   * the attribute name and the value is the function that gets called for every
   * single node that has this attribute.
   * @type {!Object}
   */
  var handlers = {
    /**
     * This handler sets the textContent of the element.
     * @param {!HTMLElement} element The node to modify.
     * @param {string} key The name of the value in |data|.
     * @param {!LoadTimeData} data The data source to draw from.
     * @param {!Set<ProcessingRoot>} visited
     */
    'i18n-content': function(element, key, data, visited) {
      element.textContent = data.getString(key);
    },

    /**
     * This handler adds options to a <select> element.
     * @param {!HTMLElement} select The node to modify.
     * @param {string} key The name of the value in |data|. It should
     *     identify an array of values to initialize an <option>. Each value,
     *     if a pair, represents [content, value]. Otherwise, it should be a
     *     content string with no value.
     * @param {!LoadTimeData} data The data source to draw from.
     * @param {!Set<ProcessingRoot>} visited
     */
    'i18n-options': function(select, key, data, visited) {
      var options = data.getValue(key);
      options.forEach(function(optionData) {
        var option = typeof optionData == 'string' ?
            new Option(optionData) :
            new Option(optionData[1], optionData[0]);
        select.appendChild(option);
      });
    },

    /**
     * This is used to set HTML attributes and DOM properties. The syntax is:
     *   attributename:key;
     *   .domProperty:key;
     *   .nested.dom.property:key
     * @param {!HTMLElement} element The node to modify.
     * @param {string} attributeAndKeys The path of the attribute to modify
     *     followed by a colon, and the name of the value in |data|.
     *     Multiple attribute/key pairs may be separated by semicolons.
     * @param {!LoadTimeData} data The data source to draw from.
     * @param {!Set<ProcessingRoot>} visited
     */
    'i18n-values': function(element, attributeAndKeys, data, visited) {
      var parts = attributeAndKeys.replace(/\s/g, '').split(/;/);
      parts.forEach(function(part) {
        if (!part)
          return;

        var attributeAndKeyPair = part.match(/^([^:]+):(.+)$/);
        if (!attributeAndKeyPair)
          throw new Error('malformed i18n-values: ' + attributeAndKeys);

        var propName = attributeAndKeyPair[1];
        var propExpr = attributeAndKeyPair[2];

        var value = data.getValue(propExpr);

        // Allow a property of the form '.foo.bar' to assign a value into
        // element.foo.bar.
        if (propName[0] == '.') {
          var path = propName.slice(1).split('.');
          var targetObject = element;
          while (targetObject && path.length > 1) {
            targetObject = targetObject[path.shift()];
          }
          if (targetObject) {
            targetObject[path] = value;
            // In case we set innerHTML (ignoring others) we need to recursively
            // check the content.
            if (path == 'innerHTML') {
              for (var i = 0; i < element.children.length; ++i) {
                processWithoutCycles(element.children[i], data, visited, false);
              }
            }
          }
        } else {
          element.setAttribute(propName, /** @type {string} */(value));
        }
      });
    }
  };

  var prefixes = [''];

  // Only look through shadow DOM when it's supported. As of April 2015, iOS
  // Chrome doesn't support shadow DOM.
  if (Element.prototype.createShadowRoot)
    prefixes.push('* /deep/ ');

  var attributeNames = Object.keys(handlers);
  var selector = prefixes.map(function(prefix) {
    return prefix + '[' + attributeNames.join('], ' + prefix + '[') + ']';
  }).join(', ');

  /**
   * Processes a DOM tree using a |data| source to populate template values.
   * @param {!ProcessingRoot} root The root of the DOM tree to process.
   * @param {!LoadTimeData} data The data to draw from.
   */
  function process(root, data) {
    processWithoutCycles(root, data, new Set(), true);
  }

  /**
   * Internal process() method that stops cycles while processing.
   * @param {!ProcessingRoot} root
   * @param {!LoadTimeData} data
   * @param {!Set<ProcessingRoot>} visited Already visited roots.
   * @param {boolean} mark Whether nodes should be marked processed.
   */
  function processWithoutCycles(root, data, visited, mark) {
    if (visited.has(root)) {
      // Found a cycle. Stop it.
      return;
    }

    // Mark the node as visited before recursing.
    visited.add(root);

    var importLinks = root.querySelectorAll('link[rel=import]');
    for (var i = 0; i < importLinks.length; ++i) {
      var importLink = /** @type {!HTMLLinkElement} */(importLinks[i]);
      if (!importLink.import) {
        // Happens when a <link rel=import> is inside a <template>.
        // TODO(dbeam): should we log an error if we detect that here?
        continue;
      }
      processWithoutCycles(importLink.import, data, visited, mark);
    }

    var templates = root.querySelectorAll('template');
    for (var i = 0; i < templates.length; ++i) {
      var template = /** @type {HTMLTemplateElement} */(templates[i]);
      if (!template.content)
        continue;
      processWithoutCycles(template.content, data, visited, mark);
    }

    var isElement = root instanceof Element;
    if (isElement && root.webkitMatchesSelector(selector))
      processElement(/** @type {!Element} */(root), data, visited);

    var elements = root.querySelectorAll(selector);
    for (var i = 0; i < elements.length; ++i) {
      processElement(elements[i], data, visited);
    }

    if (mark) {
      var processed = isElement ? [root] : root.children;
      if (processed) {
        for (var i = 0; i < processed.length; ++i) {
          processed[i].setAttribute('i18n-processed', '');
        }
      }
    }
  }

  /**
   * Run through various [i18n-*] attributes and populate.
   * @param {!Element} element
   * @param {!LoadTimeData} data
   * @param {!Set<ProcessingRoot>} visited
   */
  function processElement(element, data, visited) {
    for (var i = 0; i < attributeNames.length; i++) {
      var name = attributeNames[i];
      var attribute = element.getAttribute(name);
      if (attribute != null)
        handlers[name](element, attribute, data, visited);
    }
  }

  return {
    process: process
  };
}());


i18nTemplate.process(document, loadTimeData);
</script>`
