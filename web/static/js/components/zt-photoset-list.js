(function() {
"use strict";
function esc(s) { return String(s).replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/"/g,"&quot;"); }
function ZtPhotosetList() {
  return Reflect.construct(HTMLElement, [], ZtPhotosetList);
}
ZtPhotosetList.prototype = Object.create(HTMLElement.prototype);
ZtPhotosetList.prototype.connectedCallback = function() {
  var self = this;
  fetch("/api/photoset", { credentials: "same-origin" })
    .then(function(r) { return r.json(); })
    .then(function(data) {
      var items = data.items || [];
      var html = '<div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-images"></i></span><h3>Photosets</h3><hr /></div><div class="row">';
      if (items.length === 0) {
        html += '<div class="col-md-12"><div class="alert alert-warning">No photosets yet.</div></div>';
      } else {
        items.forEach(function(ps) {
          var id = ps.ID || ps.id;
          var name = esc(ps.Name || ps.name || "Untitled");
          var status = ps.Status || ps.status || "";
          var cover = "/api/photoset/" + id + "/cover";
          html += '<div class="col-md-3 col-sm-6 mb-4"><a href="/photoset/' + id + '" class="text-decoration-none text-dark">';
          html += '<div class="card h-100"><div class="ratio ratio-4x3 bg-light">';
          html += '<img class="lazy" data-src="' + cover + '" alt="' + name + '" style="object-fit:cover;width:100%;height:100%">';
          html += '</div><div class="card-body p-2"><h6 class="card-title mb-0">' + name + '</h6>';
          if (status && status !== "ready") html += '<small class="text-muted">' + esc(status) + '</small>';
          html += '</div></div></a></div>';
        });
      }
      html += "</div>";
      self.innerHTML = html;
      if (window.zt && window.zt.loadLazyIn) window.zt.loadLazyIn(self);
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() {
      self.innerHTML = '<div class="alert alert-danger">Failed to load photosets.</div>';
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    });
};
customElements.define("zt-photoset-list", ZtPhotosetList);
})();