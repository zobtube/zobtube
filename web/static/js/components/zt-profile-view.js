(function() {
"use strict";
function escapeAttr(s) { return String(s).replace(/&/g,"&amp;").replace(/"/g,"&quot;").replace(/</g,"&lt;").replace(/>/g,"&gt;"); }
function ZtProfileView() {
  var el = Reflect.construct(HTMLElement, [], ZtProfileView);
  return el;
}
ZtProfileView.prototype = Object.create(HTMLElement.prototype);
ZtProfileView.prototype.connectedCallback = function() {
  var self = this;
  fetch("/api/profile", { credentials: "same-origin" })
    .then(function(r) {
      if (!r.ok) throw new Error(r.status);
      return r.json();
    })
    .then(function(data) {
      var videoViews = data.video_views || [];
      var actorViews = data.actor_views || [];
      var html = '<div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-user"></i></span><h3>Your account</h3><hr /></div>';
      html += '<div class="row"><div class="col-md-12"><h3>Most viewed videos</h3><div class="row">';
      if (videoViews.length === 0) {
        html += '<div class="col-md-12"><div class="alert alert-warning" role="alert">No trending videos so far!</div></div>';
      } else {
        videoViews.forEach(function(vv) {
          var v = vv.Video || vv.video;
          if (!v) return;
          html += '<div class="col-md-3"><zt-video-view-tile data-item="'+escapeAttr(JSON.stringify({video:v,count:vv.Count||vv.count||0}))+'"></zt-video-view-tile></div>';
        });
      }
      html += '</div></div><div class="col-md-12"><h3 style="padding-top:50px;">Most viewed actors</h3><div class="row">';
      if (actorViews.length === 0) {
        html += '<div class="col-md-12"><div class="alert alert-warning" role="alert">No video viewed so far!</div></div>';
      } else {
        actorViews.forEach(function(av) {
          var a = av.actor || av.Actor;
          var c = av.Count || av.count || 0;
          if (!a) return;
          var aid = a.ID || a.id;
          var an = (a.Name||a.name||"").replace(/&/g,"&amp;").replace(/</g,"&lt;");
          var thumb = (a.Thumbnail||a.thumbnail) ? '<img class="lazy" data-src="/api/actor/'+encodeURIComponent(aid)+'/thumb" class="card-img-top lazy" alt="">' : "";
          html += '<div class="col-md-2"><div class="video-img"><a href="/actor/'+encodeURIComponent(aid)+'">'+thumb+'</a></div><div class="video-content"><h4><a href="/actor/'+encodeURIComponent(aid)+'" class="video-title">'+an+'</a></h4><div class="video-counter"><div class="video-viewers"><span class="fa fa-eye view-icon"></span><span>'+c+'</span></div></div></div></div>';
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
customElements.define("zt-profile-view", ZtProfileView);
})();
