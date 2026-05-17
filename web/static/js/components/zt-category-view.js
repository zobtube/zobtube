(function() {
"use strict";
function ZtCategoryView() {
  var el = Reflect.construct(HTMLElement, [], ZtCategoryView);
  return el;
}
ZtCategoryView.prototype = Object.create(HTMLElement.prototype);
ZtCategoryView.prototype.connectedCallback = function() {
  var self = this;
  var id = this.getAttribute("data-id");
  if (!id) { self.innerHTML = "Missing id"; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); return; }
  fetch("/api/category/" + encodeURIComponent(id), { credentials: "same-origin" })
    .then(function(r) { if (!r.ok) throw new Error(r.status); return r.json(); })
    .then(function(data) {
      var sub = data;
      var name = (sub.Name||sub.name||"").replace(/&/g,"&amp;").replace(/</g,"&lt;");
      var videos = sub.Videos || sub.videos || [];
      var actors = sub.Actors || sub.actors || [];
      var byId = {};
      videos.forEach(function(v) {
        var vid = v.ID || v.id;
        if (vid) byId[vid] = v;
      });
      actors.forEach(function(act) {
        (act.Videos || act.videos || []).forEach(function(v) {
          var vid = v.ID || v.id;
          if (vid && !byId[vid]) byId[vid] = v;
        });
      });
      var uniqueVideos = Object.keys(byId).map(function(k) { return byId[k]; });
      var html = '<h2 class="card-title actor_name">Category: '+name+'</h2><hr /><div class="row row-cols-1 row-cols-md-6 g-4">';
      uniqueVideos.forEach(function(v) { html += (window.ztThumbPreviewHtml || function() { return ""; })(v); });
      html += "</div>";
      self.innerHTML = html;
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() { self.innerHTML = '<div class="alert alert-danger">Not found.</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); });
};
customElements.define("zt-category-view", ZtCategoryView);
})();
