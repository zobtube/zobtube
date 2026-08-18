(function() {
"use strict";

// Shared HTML escaping helpers.
//
// Components build markup by string concatenation and assign it to innerHTML,
// so every interpolated value must be escaped first. Each component used to
// carry its own esc() copy and they disagreed on which characters to cover
// (some only &, some missed >, none covered '), which left holes wherever a
// value landed in a single-quoted attribute or an inline handler.
//
// escHtml covers all five characters, which is safe in both text and attribute
// contexts, so a single function is enough. escAttr is kept as an alias so the
// existing call sites keep reading naturally.

var ENTITIES = {
  "&": "&amp;",
  "<": "&lt;",
  ">": "&gt;",
  '"': "&quot;",
  "'": "&#39;"
};

function escHtml(s) {
  return String(s == null ? "" : s).replace(/[&<>"']/g, function(ch) {
    return ENTITIES[ch];
  });
}

// safeUrl returns a URL that is safe to place in href/src, or "#" when the
// scheme is one that can execute script (javascript:, vbscript:, data: ...).
// Relative and same-origin URLs pass through unchanged.
function safeUrl(u) {
  var raw = String(u == null ? "" : u);
  // Browsers ignore control characters and whitespace while parsing a scheme,
  // so "java\tscript:alert(1)" runs. Strip them before testing.
  var probe = raw.replace(/[\x00-\x20]/g, "").toLowerCase();
  var colon = probe.indexOf(":");
  if (colon > -1) {
    var scheme = probe.slice(0, colon);
    // A "scheme" containing / ? or # is not a scheme at all: it is a path such
    // as "/foo/bar:baz" or "?a=b:c", which is relative and therefore fine.
    if (!/[/?#]/.test(scheme) && scheme !== "http" && scheme !== "https" && scheme !== "mailto") {
      return "#";
    }
  }
  return raw;
}

window.zt = window.zt || {};
window.zt.escHtml = escHtml;
window.zt.escAttr = escHtml;
window.zt.safeUrl = safeUrl;

// Bare globals: components are plain IIFEs loaded via <script>, so this is how
// they reach shared code (see window.ztVideoUrlView and friends).
window.ztEscHtml = escHtml;
window.ztEscAttr = escHtml;
window.ztSafeUrl = safeUrl;
})();
