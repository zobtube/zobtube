(function() {
"use strict";
function ZtAdmObjectList() {
  var el = Reflect.construct(HTMLElement, [], ZtAdmObjectList);
  return el;
}
ZtAdmObjectList.prototype = Object.create(HTMLElement.prototype);
ZtAdmObjectList.prototype.connectedCallback = function() {
  var self = this;
  var obj = this.getAttribute("data-object") || "Video";
  var api = obj === "Video" ? "/api/adm/video" : obj === "Actor" ? "/api/adm/actor" : obj === "Channel" ? "/api/adm/channel" : "/api/adm/user";
  var tab = "overview";
  fetch(api, { credentials: "same-origin" })
    .then(function(r) { return r.json(); })
    .then(function(d) {
      var items = d.items || [];
      var base = obj === "Video" ? "/video" : obj === "Actor" ? "/actor" : obj === "Channel" ? "/channel" : "/adm/user";
      var html = '<zt-adm-tabs data-active="'+tab+'"></zt-adm-tabs><div class="row"><div class="col-md-12"><h4>'+obj+'s</h4><ul class="list-group">';
      items.forEach(function(it) {
        var id = it.ID || it.id;
        var name = (it.Name || it.name || it.Username || it.username || "").replace(/&/g,"&amp;");
        html += '<li class="list-group-item"><a href="'+base+'/'+id+(obj==="Video"?'/edit':'')+'">'+name+'</a></li>';
      });
      html += '</ul></div></div>';
      self.innerHTML = html;
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() { self.innerHTML = '<div class="alert alert-danger">Failed.</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); });
};
customElements.define("zt-adm-object-list", ZtAdmObjectList);
})();
