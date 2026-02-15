(function() {
"use strict";
function thumbPreviewHtml(v) {
  var urlView = (v.Type||v.type)==="c" ? "/clip/"+(v.ID||v.id) : "/video/"+(v.ID||v.id);
  var urlThumb = "/api/video/"+(v.ID||v.id)+"/thumb_xs";
  var dur = (v.Duration||v.duration) ? (window.ztNiceDuration||function(ns){var s=Math.floor(ns/1e9),m=Math.floor(s/60);s%=60;var h=Math.floor(m/60);m%=60;return h>0?String(h).padStart(2,0)+":"+String(m).padStart(2,0)+":"+String(s).padStart(2,0):String(m).padStart(2,0)+":"+String(s).padStart(2,0)})(v.Duration||v.duration) : "";
  var name = (v.Name||v.name||"").replace(/&/g,"&amp;").replace(/</g,"&lt;");
  var actors = v.Actors||v.actors||[];
  var actorsHtml = actors.length ? '<span class="video-posts-author"><i>Featuring: </i></span>' + actors.map(function(a){
    var an = (a.Name||a.name||"").replace(/&/g,"&amp;").replace(/</g,"&lt;");
    return '<span class="video-posts-author"><a href="/actor/'+(a.ID||a.id)+'">'+an+'</a></span>';
  }).join(" ") : "";
  return '<div class="single-review"><div class="review-img"><a href="'+urlView+'"><img src="'+urlThumb+'" class="lazy"></a>'+(dur?'<span class="video-duration">'+dur+'</span>':'')+'</div><div class="review-content"><h4><a href="'+urlView+'" class="video-title">'+name+'</a></h4>'+actorsHtml+'</div></div>';
}
window.ztThumbPreviewHtml = thumbPreviewHtml;
})();
