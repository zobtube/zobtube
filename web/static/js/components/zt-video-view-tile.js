(function() {
"use strict";
function ZtVideoViewTile() {
  var el = Reflect.construct(HTMLElement, [], ZtVideoViewTile);
  return el;
}
ZtVideoViewTile.prototype = Object.create(HTMLElement.prototype);
ZtVideoViewTile.prototype.connectedCallback = function() {
  var d = JSON.parse(this.getAttribute("data-item") || "{}");
  var v = d.video || d;
  var count = d.count || 0;
  var id = v.ID || v.id;
  var t = v.Type || v.type || "v";
  var urlView = t === "c" ? "/clip/" + id : "/video/" + id;
  var urlThumb = "/api/video/" + id + "/thumb_xs";
  var dur = (v.Duration || v.duration) ? (window.ztNiceDuration || function(ns){var s=Math.floor(ns/1e9),m=Math.floor(s/60);s%=60;var h=Math.floor(m/60);m%=60;return h>0?String(h).padStart(2,0)+":"+String(m).padStart(2,0)+":"+String(s).padStart(2,0):String(m).padStart(2,0)+":"+String(s).padStart(2,0)})(v.Duration||v.duration) : "";
  var name = (v.Name||v.name||v.Filename||v.filename||"Untitled").replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/"/g,"&quot;");
  this.innerHTML = '<div class="single-video"><div class="video-img">' +
    '<a href="' + urlView + '"><img class="lazy" data-src="' + urlThumb + '"></a>' +
    (dur ? '<span class="video-duration">' + dur + '</span>' : '') +
    '</div><div class="video-mini-title"><h4><a href="' + urlView + '" class="video-title">' + name + '</a></h4>' +
    '<div class="video-counter"><div class="video-viewers"><span class="fa fa-eye view-icon"></span><span>' + count + '</span></div></div></div>';
};
customElements.define("zt-video-view-tile", ZtVideoViewTile);
})();
