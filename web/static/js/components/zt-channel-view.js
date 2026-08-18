(function() {
"use strict";
function ZtChannelView() {
  var el = Reflect.construct(HTMLElement, [], ZtChannelView);
  return el;
}
ZtChannelView.prototype = Object.create(HTMLElement.prototype);
ZtChannelView.prototype.connectedCallback = function() {
  var self = this;
  var id = this.getAttribute("data-id");
  if (!id) { self.innerHTML = "Missing id"; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); return; }
  fetch("/api/channel/" + encodeURIComponent(id), { credentials: "same-origin" })
    .then(function(r) { if (!r.ok) throw new Error(r.status); return r.json(); })
    .then(function(data) {
      var ch = data.channel || data;
      var videos = data.videos || [];
      var admin = (window.__USER__ && window.__USER__.admin);
      var esc = window.ztEscHtml;
      var name = esc(ch.Name||ch.name||"");
      var idAttr = esc(id);
      var html = '<div style="display:flex"><div style="width:250px;margin-right:25px"><img class="img-rounded" src="/api/channel/'+idAttr+'/thumb" style="height:250px;width:250px"></div><div><h2 class="card-title actor_name">'+name+'</h2>'+(admin?' <a href="/channel/'+idAttr+'/edit"><i>Edit channel</i></a>':'')+'</div></div><hr /><br /><div class="row row-cols-1 row-cols-md-6 g-4">';
      videos.forEach(function(v){ html += (window.ztThumbPreviewHtml||function(){return"";})(v); });
      html += "</div>";
      self.innerHTML = html;
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() { self.innerHTML = '<div class="alert alert-danger">Not found.</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); });
};
customElements.define("zt-channel-view", ZtChannelView);
})();
