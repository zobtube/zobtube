(function() {
"use strict";
function ZtActorView() {
  var el = Reflect.construct(HTMLElement, [], ZtActorView);
  return el;
}
ZtActorView.prototype = Object.create(HTMLElement.prototype);
ZtActorView.prototype.connectedCallback = function() {
  var self = this;
  var id = this.getAttribute("data-id");
  if (!id) { self.innerHTML = "Missing id"; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); return; }
  fetch("/api/actor/" + encodeURIComponent(id), { credentials: "same-origin" })
    .then(function(r) { if (!r.ok) throw new Error(r.status); return r.json(); })
    .then(function(a) {
      var admin = (window.__USER__ && window.__USER__.admin);
      var aliases = (a.Aliases || a.aliases || []).map(function(x){ return x.Name || x.name; }).filter(Boolean).join(" / ");
      var catsHtml = (a.Categories || a.categories || []).map(function(s){ var n=(s.Name||s.name||"").replace(/&/g,"&amp;").replace(/</g,"&lt;"); return '<div class="category">'+n+'</div>'; }).join("");
      var links = (a.Links || a.links || []);
      var videos = a.Videos || a.videos || [];
      var name = (a.Name||a.name||"").replace(/&/g,"&amp;").replace(/</g,"&lt;");
      var desc = (a.Description||a.description||"").replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/\n/g,"<br>");
      var linksHtml = links.map(function(l){return '<a href="'+(l.URL||l.url)+'" target="_blank" rel="noopener noreferrer"><img class="img-rounded" src="/static/images/provider-'+(l.Provider||l.provider)+'.png" style="height:80px;width:80px;margin-top:5px"></a>';}).join('');
      var html = '<style>.actor_name{font-size:3rem}.bio_detail_label{color:#6b6b6b}</style><div style="display:flex"><div style="width:250px;margin-right:25px"><img class="img-rounded" src="/api/actor/'+id+'/thumb" style="height:250px;width:250px"></div><div id="bio" style="flex-grow:1;margin:0"><h2 class="card-title actor_name">'+name+'</h2>'+(aliases?'<h4>aka '+aliases+'</h4>':'')+(admin?' <a href="/actor/'+id+'/edit"><i>Edit profile</i></a>':'')+'<div class="categories mb-4" style="padding-top:25px">'+catsHtml+'</div><div>'+desc+'</div></div><div style="display:flex"><div style="width:170px"><div style="margin-top:20px;float:right">'+linksHtml+'</div></div></div></div><hr /><br /><div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-play"></i></span><h3>Videos</h3></div><div class="row row-cols-1 row-cols-md-6 g-4">';
      videos.forEach(function(v){ html += (window.ztThumbPreviewHtml || function(){ return ""; })(v); });
      html += "</div>";
      self.innerHTML = html;
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() { self.innerHTML = '<div class="alert alert-danger">Not found.</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); });
};
customElements.define("zt-actor-view", ZtActorView);
})();
