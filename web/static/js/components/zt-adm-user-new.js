(function() {
"use strict";
function ZtAdmUserNew() {
  var el = Reflect.construct(HTMLElement, [], ZtAdmUserNew);
  return el;
}
ZtAdmUserNew.prototype = Object.create(HTMLElement.prototype);
ZtAdmUserNew.prototype.connectedCallback = function() {
  var self = this;
  self.innerHTML = '<div class="row"><div class="col-md-3 col-lg-3"><zt-adm-tabs data-active="users"></zt-adm-tabs></div><div class="col-md-9 col-lg-9"><div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-user-plus"></i></span><h3>Add a new user</h3><hr /></div><p class="small text-muted mb-2"><a href="/adm/users">← User list</a></p><form id="adm-user-new-form"><div class="mb-3"><label for="adm-username" class="form-label">Username</label><input type="text" class="form-control" id="adm-username" name="username" placeholder="my-new-user" required></div><div class="mb-3"><label for="adm-password" class="form-label">Password</label><input type="password" class="form-control" id="adm-password" name="password" required></div><div class="form-check form-switch mb-3"><label class="form-check-label" for="adm-admin">Has admin rights</label><input class="form-check-input" type="checkbox" id="adm-admin" name="admin" value="x"></div><div><input class="btn btn-primary" type="submit" value="Create new user"/></div></form></div></div>';
  self.querySelector("#adm-user-new-form").onsubmit = function(e) {
    e.preventDefault();
    var form = e.target;
    fetch("/api/adm/user", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "same-origin",
      body: JSON.stringify({
        username: form.username.value,
        password: form.password.value,
        admin: form.admin.checked
      })
    }).then(function(r) {
      if (r.ok) return r.json();
      throw new Error(r.status);
    }).then(function(d) {
      if (d.redirect) window.navigate(d.redirect);
      else window.navigate("/adm/users");
    }).catch(function(err) {
      var alert = document.createElement("div");
      alert.className = "alert alert-danger mt-3";
      alert.textContent = "Failed to create user.";
      self.querySelector("form").appendChild(alert);
    });
  };
  if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
};
customElements.define("zt-adm-user-new", ZtAdmUserNew);
})();
