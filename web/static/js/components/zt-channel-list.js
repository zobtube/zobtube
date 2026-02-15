(function() {
"use strict";
function ZtChannelList() {
  var el = Reflect.construct(HTMLElement, [], ZtChannelList);
  return el;
}
ZtChannelList.prototype = Object.create(HTMLElement.prototype);
ZtChannelList.prototype.connectedCallback = function() {
  var self = this;
  var admin = (window.__USER__ && window.__USER__.admin);
  fetch("/api/channel", { credentials: "same-origin" })
    .then(function(r) { return r.json(); })
    .then(function(data) {
      var items = data.items || [];
      var html = '<div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-podcast"></i></span><h3>Channels' +
        (admin ? ' <a href="/channel/new"><i class="fas fa-plus-circle"></i></a>' : '') + '</h3><hr /></div>' +
        '<div class="row row-cols-1 row-cols-md-4 g-4">';
      items.forEach(function(ch) {
        var id = ch.ID || ch.id;
        var name = (ch.Name || ch.name || "").replace(/&/g,"&amp;").replace(/</g,"&lt;");
        var thumbUrl = "/api/channel/" + encodeURIComponent(id) + "/thumb";
        html += '<div class="col"><div class="card"><a href="/channel/' + id + '"><img data-src="' + thumbUrl + '" class="card-img-top lazy" alt=""></a><div class="card-body"><h5 class="card-title"><a class="stretched-link" href="/channel/' + id + '">' + name + '</a></h5></div></div></div>';
      });
      html += "</div>";
      self.innerHTML = html;
      if (window.zt && window.zt.loadLazyIn) window.zt.loadLazyIn(self);
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() { self.innerHTML = '<div class="alert alert-danger">Failed to load.</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); });
};
customElements.define("zt-channel-list", ZtChannelList);
})();
