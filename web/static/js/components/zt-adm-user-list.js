(function() {
"use strict";
function ZtAdmUserList() {
  var el = Reflect.construct(HTMLElement, [], ZtAdmUserList);
  return el;
}
ZtAdmUserList.prototype = Object.create(HTMLElement.prototype);
ZtAdmUserList.prototype.connectedCallback = function() {
  var self = this;
  var admin = (window.__USER__ && window.__USER__.admin);
  fetch("/api/adm/user", { credentials: "same-origin" })
    .then(function(r) { return r.json(); })
    .then(function(d) {
      var items = d.items || [];
      var addLink = admin ? ' <a href="/adm/user"><i class="fas fa-plus-circle"></i></a>' : '';
      var html = '<div class="row"><div class="col-md-3 col-lg-3"><zt-adm-tabs data-active="users"></zt-adm-tabs></div><div class="col-md-9 col-lg-9"><div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-user-circle"></i></span><h3>User list'+addLink+'</h3><hr /></div><table class="table table-striped table-hover"><thead><tr><th>Username</th><th>Has admin rights</th><th></th></tr></thead><tbody>';
      items.forEach(function(u) {
        var id = u.ID || u.id;
        var un = (u.Username||u.username||"").replace(/&/g,"&amp;").replace(/</g,"&lt;");
        var adminIcon = (u.Admin || u.admin) ? 'far fa-check-circle' : 'fas fa-ban';
        var deleteUrl = '/api/adm/user/' + id;
        html += '<tr><td>'+un+'</td><td><i class="'+adminIcon+'"></i></td><td style="text-align:end"><button type="button" class="btn btn-sm btn-danger" data-user-id="'+id+'">Delete</button></td></tr>';
      });
      html += '</tbody></table></div></div>';
      self.innerHTML = html;
      self.querySelectorAll("button[data-user-id]").forEach(function(btn) {
        btn.onclick = function() {
          if (!confirm("Delete user?")) return;
          fetch("/api/adm/user/" + btn.dataset.userId, { method: "DELETE", credentials: "same-origin" })
            .then(function(r) { if (r.ok) window.navigate("/adm/users"); });
        };
      });
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() { self.innerHTML = '<div class="alert alert-danger">Failed.</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); });
};
customElements.define("zt-adm-user-list", ZtAdmUserList);
})();
