(function() {
"use strict";
function ZtChannelEdit() {
  var el = Reflect.construct(HTMLElement, [], ZtChannelEdit);
  return el;
}
ZtChannelEdit.prototype = Object.create(HTMLElement.prototype);
ZtChannelEdit.prototype.connectedCallback = function() {
  var self = this;
  var id = this.getAttribute("data-id");
  if (!id) { self.innerHTML = "Missing id"; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); return; }
  fetch("/api/channel/" + encodeURIComponent(id), { credentials: "same-origin" })
    .then(function(r) { return r.json(); })
    .then(function(data) {
      var ch = data.channel || data;
      var name = (ch.Name||ch.name||"").replace(/&/g,"&amp;");
      self.innerHTML = '<div class="themeix-section-h"><h3>Edit channel</h3><hr /></div><form><div class="mb-3"><label class="form-label">Name</label><input name="name" class="form-control" value="'+name+'" required></div><button type="submit" class="btn btn-primary">Save</button></form>';
      self.querySelector("form").onsubmit = function(e) {
        e.preventDefault();
        fetch("/api/channel/"+id, { method: "PUT", body: JSON.stringify({ name: e.target.name.value }), headers: { "Content-Type": "application/json" }, credentials: "same-origin" })
          .then(function() { window.navigate("/channel/"+id); });
      };
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() { self.innerHTML = '<div class="alert alert-danger">Not found.</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); });
};
customElements.define("zt-channel-edit", ZtChannelEdit);
})();
