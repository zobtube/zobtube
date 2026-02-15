(function() {
"use strict";
function ZtActorCreate() {
  var el = Reflect.construct(HTMLElement, [], ZtActorCreate);
  return el;
}
ZtActorCreate.prototype = Object.create(HTMLElement.prototype);
ZtActorCreate.prototype.connectedCallback = function() {
  var admin = (window.__USER__ && window.__USER__.admin);
  if (!admin) { this.innerHTML = '<div class="alert alert-danger">Forbidden</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(this); return; }
  this.innerHTML = '<div class="themeix-section-h"><h3>New actor</h3><hr /></div><form id="actor-create-form"><div class="mb-3"><label class="form-label">Name</label><input name="name" class="form-control" required></div><button type="submit" class="btn btn-primary">Create</button></form>';
  if (window.zt && window.zt.pageReady) window.zt.pageReady(this);
  this.querySelector("form").onsubmit = function(e) {
    e.preventDefault();
    var fd = new FormData(e.target);
    fd.set("name", fd.get("name"));
    fetch("/api/actor/", { method: "POST", body: fd, credentials: "same-origin" })
      .then(function(r) { return r.json(); })
      .then(function(d) { window.navigate("/actor/" + (d.result || d.id || d.ID)); })
      .catch(function() { alert("Failed"); });
  };
};
customElements.define("zt-actor-create", ZtActorCreate);
})();
