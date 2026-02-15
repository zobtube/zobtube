(function() {
"use strict";
function clipUrlView(v) {
  if (!v) return "#";
  var id = v.ID || v.id;
  return "/clip/" + id;
}
function clipUrlThumbXS(v) {
  var id = (v && (v.ID || v.id));
  return id ? "/api/video/" + id + "/thumb_xs" : "";
}
function ZtClipList() {
  var el = Reflect.construct(HTMLElement, [], ZtClipList);
  return el;
}
ZtClipList.prototype = Object.create(HTMLElement.prototype);
ZtClipList.prototype.connectedCallback = function() {
  var self = this;
  fetch("/api/clip", { credentials: "same-origin" })
    .then(function(r) { return r.json(); })
    .then(function(data) {
      var items = data.items || [];
      var css = ".zt-clip-grid{column-count:4;column-gap:5px}.zt-clip-grid .zt-clip-column{width:100%;margin-top:0;margin-bottom:5px;break-inside:avoid}";
      var html = '<style>' + css + '</style>';
      html += '<div class="row"><div class="col-md-12">';
      html += '<div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-bolt"></i></span><h3>All clips</h3><hr /></div>';
      if (items.length === 0) {
        html += '<div class="alert alert-warning" role="alert">No clip available so far!</div>';
      } else {
        html += '<div class="zt-clip-grid">';
        items.forEach(function(v) {
          var urlView = clipUrlView(v);
          var urlThumb = clipUrlThumbXS(v);
          html += '<div class="zt-clip-column"><a href="' + urlView + '"><img class="lazy" alt="Clip" data-src="' + urlThumb + '"></a></div>';
        });
        html += '</div>';
      }
      html += '</div></div>';
      self.innerHTML = html;
      if (window.zt && window.zt.loadLazyIn) window.zt.loadLazyIn(self);
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() { self.innerHTML = '<div class="alert alert-danger">Failed to load.</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); });
};
customElements.define("zt-clip-list", ZtClipList);
})();