(function() {
"use strict";
function escapeAttr(s) {
  return String(s).replace(/&/g,"&amp;").replace(/"/g,"&quot;").replace(/</g,"&lt;").replace(/>/g,"&gt;");
}
function ZtHome() {
  var el = Reflect.construct(HTMLElement, [], ZtHome);
  return el;
}
ZtHome.prototype = Object.create(HTMLElement.prototype);
ZtHome.prototype.connectedCallback = function() {
  var self = this;
  fetch("/api/home", { credentials: "same-origin" })
    .then(function(r) { return r.json(); })
    .then(function(data) {
      var items = data.items || [];
      var html = '<div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-bolt"></i></span><h3>Latest Trending Videos</h3><hr /></div><div class="row">';
      if (items.length === 0) {
        html += '<div class="col-md-12"><div class="alert alert-warning" role="alert">No new video available yet!</div></div>';
      } else {
        items.forEach(function(v) {
          html += '<div class="col-md-3"><zt-video-tile data-video="' + escapeAttr(JSON.stringify(v)) + '"></zt-video-tile></div>';
        });
      }
      html += "</div>";
      self.innerHTML = html;
      if (window.zt && window.zt.loadLazyIn) window.zt.loadLazyIn(self);
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() {
      self.innerHTML = '<div class="alert alert-danger">Failed to load.</div>';
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    });
};
customElements.define("zt-home", ZtHome);
})();
