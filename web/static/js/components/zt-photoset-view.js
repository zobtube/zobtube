(function() {
"use strict";
function esc(s) { return String(s).replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/"/g,"&quot;"); }
function ZtPhotosetView() {
  return Reflect.construct(HTMLElement, [], ZtPhotosetView);
}
ZtPhotosetView.prototype = Object.create(HTMLElement.prototype);
ZtPhotosetView.prototype.connectedCallback = function() {
  var self = this;
  var id = this.getAttribute("data-id");
  if (!id) { self.innerHTML = "Missing id"; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); return; }
  fetch("/api/photoset/" + encodeURIComponent(id), { credentials: "same-origin" })
    .then(function(r) { if (!r.ok) throw new Error(r.status); return r.json(); })
    .then(function(data) {
      var ps = data.photoset || data;
      var photos = data.photos || [];
      var cats = data.categories || {};
      var name = esc(ps.Name || ps.name || "Untitled");
      var admin = window.__USER__ && window.__USER__.admin;
      var channel = ps.Channel || ps.channel;
      var channelHtml = channel ? '<a class="btn btn-sm btn-dark" href="/channel/'+(channel.ID||channel.id)+'">' + esc(channel.Name||channel.name||"") + '</a>' : "";
      var actors = ps.Actors || ps.actors || [];
      var actorsHtml = actors.map(function(a) {
        return '<a class="btn btn-sm btn-danger me-1" href="/actor/'+(a.ID||a.id)+'">' + esc(a.Name||a.name||"") + '</a>';
      }).join("");
      var catsHtml = Object.keys(cats).map(function(cid) {
        return '<a class="btn btn-sm btn-secondary me-1" href="/category/'+cid+'">' + esc(cats[cid]) + '</a>';
      }).join("");
      var html = '<div><h3 class="mt-3">' + name + '</h3>';
      html += '<div class="mb-3 d-flex flex-wrap gap-1">' + channelHtml + actorsHtml + catsHtml;
      if (admin) html += ' <a class="btn btn-sm btn-outline-dark" href="/photoset/'+id+'/edit"><i class="fa fa-edit"></i> Edit</a>';
      html += '</div><div class="row g-2" id="zt-photoset-grid">';
      photos.forEach(function(pv, idx) {
        var p = pv.Photo || pv;
        var pid = p.ID || p.id;
        var thumb = "/api/photo/" + pid + "/thumb_mini";
        html += '<div class="col-6 col-md-3 col-lg-2"><a href="#" class="zt-photo-thumb d-block" data-index="'+idx+'">';
        html += '<img class="img-fluid rounded lazy" data-src="'+thumb+'" alt="'+esc(p.Filename||p.filename||"")+'" style="width:100%;aspect-ratio:1;object-fit:cover">';
        html += '</a></div>';
      });
      html += '</div></div>';
      self.innerHTML = html;
      self._photos = photos;
      self._photosetId = id;
      self.querySelectorAll(".zt-photo-thumb").forEach(function(a) {
        a.addEventListener("click", function(e) {
          e.preventDefault();
          var idx = parseInt(a.getAttribute("data-index"), 10) || 0;
          var lb = document.createElement("zt-photo-lightbox");
          lb.setAttribute("data-photos", JSON.stringify(photos));
          lb.setAttribute("data-index", String(idx));
          lb.setAttribute("data-photoset-id", id);
          document.body.appendChild(lb);
        });
      });
      if (window.zt && window.zt.loadLazyIn) window.zt.loadLazyIn(self);
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() {
      self.innerHTML = '<div class="alert alert-danger">Not found.</div>';
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    });
};
customElements.define("zt-photoset-view", ZtPhotosetView);
})();