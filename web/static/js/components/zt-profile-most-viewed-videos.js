(function() {
"use strict";
function escapeAttr(s) { return String(s).replace(/&/g,"&amp;").replace(/"/g,"&quot;").replace(/</g,"&lt;").replace(/>/g,"&gt;"); }
function ZtProfileMostViewedVideos() {
  var el = Reflect.construct(HTMLElement, [], ZtProfileMostViewedVideos);
  return el;
}
ZtProfileMostViewedVideos.prototype = Object.create(HTMLElement.prototype);
ZtProfileMostViewedVideos.prototype.connectedCallback = function() {
  var self = this;
  fetch("/api/profile", { credentials: "same-origin" })
    .then(function(r) {
      if (!r.ok) throw new Error(r.status);
      return r.json();
    })
    .then(function(data) {
      var videoViews = data.video_views || [];
      var html = '<div class="row"><div class="col-md-3 col-lg-3"><zt-profile-tabs data-active="most-viewed-videos"></zt-profile-tabs></div><div class="col-md-9 col-lg-9">';
      html += '<div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-video"></i></span><h3>Most viewed videos</h3><hr /></div>';
      html += '<div class="row">';
      if (videoViews.length === 0) {
        html += '<div class="col-md-12"><div class="alert alert-warning" role="alert">No trending videos so far!</div></div>';
      } else {
        videoViews.forEach(function(vv) {
          var v = vv.Video || vv.video;
          var videoID = vv.VideoID || vv.video_id || vv.videoId;
          if (v && !v.ID && !v.id && videoID) {
            v = Object.assign({}, v, { ID: videoID });
          }
          if (!v || (!v.ID && !v.id)) return;
          html += '<div class="col-md-3"><zt-video-view-tile data-item="'+escapeAttr(JSON.stringify({video:v,count:vv.Count||vv.count||0}))+'"></zt-video-view-tile></div>';
        });
      }
      html += '</div></div></div>';
      self.innerHTML = html;
      if (window.zt && window.zt.loadLazyIn) window.zt.loadLazyIn(self);
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function(e) {
      if (e && e.message === "401") {
        window.location.href = "/auth/login?next=" + encodeURIComponent(window.location.pathname);
        return;
      }
      self.innerHTML = '<div class="alert alert-danger">Failed to load.</div>';
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    });
};
customElements.define("zt-profile-most-viewed-videos", ZtProfileMostViewedVideos);
})();