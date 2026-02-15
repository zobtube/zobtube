(function() {
"use strict";
function ZtChannelCreate() {
  var el = Reflect.construct(HTMLElement, [], ZtChannelCreate);
  return el;
}
ZtChannelCreate.prototype = Object.create(HTMLElement.prototype);
ZtChannelCreate.prototype.connectedCallback = function() {
  var admin = (window.__USER__ && window.__USER__.admin);
  if (!admin) { this.innerHTML = '<div class="alert alert-danger">Forbidden</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(this); return; }
  this.innerHTML = '<div class="themeix-section-h"><h3>New channel</h3><hr /></div><form id="channel-create-form"><div class="mb-3"><label class="form-label">Name</label><input name="name" class="form-control" required></div><button type="submit" class="btn btn-primary">Create</button></form>';
  if (window.zt && window.zt.pageReady) window.zt.pageReady(this);
  this.querySelector("form").onsubmit = function(e) {
    e.preventDefault();
    fetch("/api/channel", { method: "POST", body: JSON.stringify({ name: e.target.name.value }), headers: { "Content-Type": "application/json" }, credentials: "same-origin" })
      .then(function(r) { return r.json(); })
      .then(function(d) { if (d.redirect) window.navigate(d.redirect); else window.navigate("/channel/" + d.id); });
  };
};
customElements.define("zt-channel-create", ZtChannelCreate);
})();
