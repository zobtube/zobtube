(function() {
"use strict";
function ZtVideoList() {
  var el = Reflect.construct(HTMLElement, [], ZtVideoList);
  return el;
}
ZtVideoList.prototype = Object.create(HTMLElement.prototype);
ZtVideoList.prototype.connectedCallback = function() {
  var self = this;
  var type = this.getAttribute("data-type") || "video";
  var api = type === "clip" ? "/api/clip" : type === "movie" ? "/api/movie" : "/api/video";
  var title = type === "clip" ? "Clips" : type === "movie" ? "Movies" : "Videos";
  fetch(api, { credentials: "same-origin" })
    .then(function(r) { return r.json(); })
    .then(function(data) {
      var items = data.items || [];
      var html = '<div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-film"></i></span><h3>' + title + '</h3><hr /></div><div class="row">';
      if (items.length === 0) {
        html += '<div class="col-md-12"><div class="alert alert-warning">No ' + title.toLowerCase() + ' yet.</div></div>';
      } else {
        items.forEach(function(v) {
          html += '<div class="col-md-3"><zt-video-tile data-video="' + String(JSON.stringify(v)).replace(/&/g,"&amp;").replace(/"/g,"&quot;").replace(/</g,"&lt;") + '"></zt-video-tile></div>';
        });
      }
      html += "</div>";
      self.innerHTML = html;
      if (window.zt && window.zt.loadLazyIn) window.zt.loadLazyIn(self);
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() { self.innerHTML = '<div class="alert alert-danger">Failed to load.</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); });
};
customElements.define("zt-video-list", ZtVideoList);
})();
