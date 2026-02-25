(function() {
"use strict";
function escapeHtml(s) {
  return String(s).replace(/&/g, "&amp;").replace(/"/g, "&quot;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}
function ZtProfileTokens() {
  var el = Reflect.construct(HTMLElement, [], ZtProfileTokens);
  return el;
}
ZtProfileTokens.prototype = Object.create(HTMLElement.prototype);
ZtProfileTokens.prototype.connectedCallback = function() {
  var self = this;
  var html = '<div class="row"><div class="col-12"><zt-profile-tabs data-active="tokens"></zt-profile-tabs></div></div>';
  html += '<div class="row"><div class="col-md-12"><div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-key"></i></span><h3>API tokens</h3><hr /></div>';
  html += '<p class="text-muted">Use API tokens to authenticate script or client requests with <code>Authorization: Bearer &lt;token&gt;</code>. Create a token and copy it now; it will not be shown again.</p>';
  html += '<form id="zt-profile-token-create-form" class="mb-4" style="max-width:400px"><div class="input-group"><input type="text" class="form-control" id="zt-token-name" placeholder="Token name (e.g. My script)" required><button type="submit" class="btn btn-primary">Create token</button></div><div id="zt-token-form-error" class="alert alert-danger mt-2" style="display:none" role="alert"></div></form>';
  html += '<div id="zt-tokens-list"></div>';
  html += '<div class="modal fade" id="zt-token-show-modal" tabindex="-1"><div class="modal-dialog"><div class="modal-content"><div class="modal-header"><h5 class="modal-title">Token created</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div><div class="modal-body"><p class="text-warning">Copy this token now. It will not be shown again.</p><div class="input-group"><input type="text" class="form-control font-monospace" id="zt-token-show-value" readonly><button type="button" class="btn btn-outline-secondary" id="zt-token-copy-btn">Copy</button></div></div><div class="modal-footer"><button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button></div></div></div></div>';
  self.innerHTML = html;

  var listEl = self.querySelector("#zt-tokens-list");
  var form = self.querySelector("#zt-profile-token-create-form");
  var errEl = self.querySelector("#zt-token-form-error");
  var modalEl = self.querySelector("#zt-token-show-modal");
  var tokenValueEl = self.querySelector("#zt-token-show-value");
  var copyBtn = self.querySelector("#zt-token-copy-btn");

  function loadList() {
    fetch("/api/profile/tokens", { credentials: "same-origin" })
      .then(function(r) {
        if (r.status === 401) {
          window.location.href = "/auth/login?next=" + encodeURIComponent(window.location.pathname);
          return;
        }
        return r.json();
      })
      .then(function(data) {
        if (!data) return;
        var tokens = data.tokens || [];
        if (tokens.length === 0) {
          listEl.innerHTML = '<p class="text-muted">No tokens yet. Create one above.</p>';
          return;
        }
        var table = '<table class="table table-striped"><thead><tr><th>Name</th><th>Created</th><th></th></tr></thead><tbody>';
        tokens.forEach(function(t) {
          var name = escapeHtml(t.name || "");
          var created = t.created_at ? new Date(t.created_at).toLocaleString() : "";
          var id = t.id || "";
          table += '<tr><td>' + name + '</td><td>' + escapeHtml(created) + '</td><td><button type="button" class="btn btn-sm btn-outline-danger zt-token-delete" data-id="' + escapeHtml(id) + '">Delete</button></td></tr>';
        });
        table += "</tbody></table>";
        listEl.innerHTML = table;
        listEl.querySelectorAll(".zt-token-delete").forEach(function(btn) {
          btn.addEventListener("click", function() {
            var id = btn.getAttribute("data-id");
            if (!id) return;
            if (!confirm("Delete this token? It will stop working immediately.")) return;
            fetch("/api/profile/tokens/" + encodeURIComponent(id), { method: "DELETE", credentials: "same-origin" })
              .then(function(r) {
                if (r.status === 401) {
                  window.location.href = "/auth/login?next=" + encodeURIComponent(window.location.pathname);
                  return;
                }
                if (r.status === 204 || r.ok) loadList();
              });
          });
        });
      })
      .catch(function() {
        listEl.innerHTML = '<div class="alert alert-danger">Failed to load tokens.</div>';
      });
  }

  form.addEventListener("submit", function(e) {
    e.preventDefault();
    var name = self.querySelector("#zt-token-name").value.trim();
    errEl.style.display = "none";
    errEl.textContent = "";
    if (!name) {
      errEl.textContent = "Name is required.";
      errEl.style.display = "block";
      return;
    }
    fetch("/api/profile/tokens", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "same-origin",
      body: JSON.stringify({ name: name })
    }).then(function(r) {
      if (r.status === 401) {
        window.location.href = "/auth/login?next=" + encodeURIComponent(window.location.pathname);
        return Promise.resolve();
      }
      return r.json().then(function(data) {
        if (r.ok && data && data.token) {
          form.reset();
          tokenValueEl.value = data.token;
          if (window.bootstrap && window.bootstrap.Modal) {
            var modal = new window.bootstrap.Modal(modalEl);
            modal.show();
          } else {
            modalEl.classList.add("show");
            modalEl.style.display = "block";
          }
          loadList();
        } else {
          errEl.textContent = (data && data.error) || "Failed to create token.";
          errEl.style.display = "block";
        }
      });
    }).catch(function() {
      errEl.textContent = "Request failed.";
      errEl.style.display = "block";
    });
  });

  copyBtn.addEventListener("click", function() {
    tokenValueEl.select();
    document.execCommand("copy");
    if (typeof sendToast === "function") sendToast("Copied", "", "bg-success", "Token copied to clipboard.");
  });

  loadList();
  if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
};
customElements.define("zt-profile-tokens", ZtProfileTokens);
})();
