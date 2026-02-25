(function() {
"use strict";
function esc(s) { return String(s == null ? "" : s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/"/g, "&quot;"); }
function formatDate(createdAt) {
  if (!createdAt) return "—";
  var d = new Date(createdAt);
  return isNaN(d.getTime()) ? esc(createdAt) : d.toLocaleString();
}
function ZtAdmTokenList() {
  var el = Reflect.construct(HTMLElement, [], ZtAdmTokenList);
  return el;
}
ZtAdmTokenList.prototype = Object.create(HTMLElement.prototype);
ZtAdmTokenList.prototype.connectedCallback = function() {
  var self = this;
  fetch("/api/adm/tokens", { credentials: "same-origin" })
    .then(function(r) { return r.json(); })
    .then(function(d) {
      var tokens = d.tokens || [];
      var html = '<div class="row"><div class="col-md-3 col-lg-3"><zt-adm-tabs data-active="tokens"></zt-adm-tabs></div><div class="col-md-9 col-lg-9"><div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-key"></i></span><h3>API tokens</h3><hr /></div>';
      if (tokens.length === 0) {
        html += '<p class="text-muted">No API tokens.</p>';
      } else {
        html += '<table class="table table-striped table-hover"><thead><tr><th>Token name</th><th>User</th><th>Created at</th><th></th></tr></thead><tbody>';
        tokens.forEach(function(t) {
          var id = t.id || t.ID;
          var name = esc(t.name || t.Name);
          var username = esc(t.username || t.Username || "—");
          var createdAt = formatDate(t.created_at || t.CreatedAt);
          html += '<tr><td>' + name + '</td><td>' + username + '</td><td>' + createdAt + '</td><td style="text-align:end"><button type="button" class="btn btn-sm btn-danger" data-token-id="' + esc(id) + '">Delete</button></td></tr>';
        });
        html += '</tbody></table>';
      }
      html += '</div></div>';
      self.innerHTML = html;
      self.querySelectorAll("button[data-token-id]").forEach(function(btn) {
        btn.onclick = function() {
          if (!confirm("Delete this API token? The token will stop working immediately.")) return;
          var tokenId = btn.getAttribute("data-token-id");
          fetch("/api/adm/tokens/" + encodeURIComponent(tokenId), { method: "DELETE", credentials: "same-origin" })
            .then(function(r) {
              if (r.status === 204 && typeof loadPage === "function") loadPage("/adm/tokens");
              else if (r.status === 204) window.location.reload();
            });
        };
      });
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() { self.innerHTML = '<div class="alert alert-danger">Failed to load.</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); });
};
customElements.define("zt-adm-token-list", ZtAdmTokenList);
})();
