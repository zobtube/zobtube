(function() {
"use strict";
function videoUrlView(v) {
  if (!v) return "#";
  var id = v.ID || v.id;
  var t = v.Type || v.type || "v";
  return (t === "c" ? "/clip/" : "/video/") + id;
}
function videoUrlThumbXS(v) {
  var id = (v && (v.ID || v.id));
  return id ? "/api/video/" + id + "/thumb_xs" : "";
}
function niceDuration(ns) {
  if (!ns) return "";
  var s = Math.floor(ns / 1e9);
  var m = Math.floor(s / 60); s %= 60;
  var h = Math.floor(m / 60); m %= 60;
  if (h > 0) return (h<10?"0":"")+h+":"+(m<10?"0":"")+m+":"+(s<10?"0":"")+s;
  return (m<10?"0":"")+m+":"+(s<10?"0":"")+s;
}
function ZtVideoTile() {
  var el = Reflect.construct(HTMLElement, [], ZtVideoTile);
  return el;
}
ZtVideoTile.prototype = Object.create(HTMLElement.prototype);
ZtVideoTile.prototype.connectedCallback = function() {
  var v = this._video || JSON.parse(this.getAttribute("data-video") || "null");
  if (!v) return;
  var urlView = videoUrlView(v);
  var urlThumb = videoUrlThumbXS(v);
  var d = v.Duration || v.duration;
  var dur = d ? niceDuration(d) : "";
  var name = (v.Name || v.name || v.Filename || v.filename || "Untitled").replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/"/g,"&quot;");
  this.innerHTML = '<div class="single-video"><div class="video-img">' +
    '<a href="' + urlView + '"><img class="lazy" data-src="' + urlThumb + '"></a>' +
    (dur ? '<span class="video-duration">' + dur + '</span>' : '') +
    '</div><div class="video-mini-title"><h4><a href="' + urlView + '" class="video-title">' + name + '</a></h4></div></div>';
};
customElements.define("zt-video-tile", ZtVideoTile);
window.ztVideoUrlView = videoUrlView;
window.ztVideoUrlThumbXS = videoUrlThumbXS;
window.ztNiceDuration = niceDuration;
})();
