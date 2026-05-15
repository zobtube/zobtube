(function() {
"use strict";
function esc(s) { return String(s == null ? "" : s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/"/g, "&quot;"); }
function ZtAdmMetadataStorage() {
  var el = Reflect.construct(HTMLElement, [], ZtAdmMetadataStorage);
  return el;
}
ZtAdmMetadataStorage.prototype = Object.create(HTMLElement.prototype);
ZtAdmMetadataStorage.prototype.connectedCallback = function() {
  var self = this;
  fetch("/api/adm/metadata-storage", { credentials: "same-origin" })
    .then(function(r) { return r.json(); })
    .then(function(d) {
      var typeRaw = String(d.type || "").toLowerCase();
      var typeBadge;
      if (typeRaw === "s3") {
        typeBadge = '<span class="badge text-bg-info">S3</span>';
      } else if (typeRaw === "filesystem") {
        typeBadge = '<span class="badge text-bg-secondary">Filesystem</span>';
      } else {
        typeBadge = '<span class="badge text-bg-light text-dark border">' + esc(typeRaw || "—") + '</span>';
      }
      var html = '<div class="row"><div class="col-md-3 col-lg-3"><zt-adm-tabs data-active="metadata-storage"></zt-adm-tabs></div><div class="col-md-9 col-lg-9">';
      html += '<div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-database"></i></span><h3>Metadata storage</h3><hr /></div>';
      html += '<p class="text-muted">Dedicated storage for actor, channel, and category thumbnails and related metadata assets. Configuration is read-only here; set values via CLI flags, environment variables, or <code>config.yml</code>, then restart the server.</p>';
      html += '<div class="card shadow-sm mb-3"><div class="card-body">';
      html += '<div class="d-flex flex-wrap align-items-center gap-2 mb-3">' + typeBadge + '</div>';
      if (typeRaw === "filesystem") {
        html += '<dl class="row mb-0"><dt class="col-sm-3">Path</dt><dd class="col-sm-9"><code class="text-break">' + esc(d.path || "—") + '</code></dd></dl>';
      } else if (typeRaw === "s3") {
        html += '<dl class="row mb-0">';
        html += '<dt class="col-sm-3">Bucket</dt><dd class="col-sm-9"><code>' + esc(d.bucket || "—") + '</code></dd>';
        html += '<dt class="col-sm-3">Region</dt><dd class="col-sm-9"><code>' + esc(d.region || "—") + '</code></dd>';
        if (d.prefix) {
          html += '<dt class="col-sm-3">Prefix</dt><dd class="col-sm-9"><code class="text-break">' + esc(d.prefix) + '</code></dd>';
        }
        if (d.endpoint) {
          html += '<dt class="col-sm-3">Endpoint</dt><dd class="col-sm-9"><code class="text-break">' + esc(d.endpoint) + '</code></dd>';
        }
        if (d.access_key_configured) {
          html += '<dt class="col-sm-3">Access key</dt><dd class="col-sm-9"><span class="text-muted">Configured</span></dd>';
        }
        if (d.secret_access_key_configured) {
          html += '<dt class="col-sm-3">Secret key</dt><dd class="col-sm-9"><span class="text-muted">Configured (hidden)</span></dd>';
        }
        html += '</dl>';
      }
      html += '</div></div>';
      html += '<div class="d-flex flex-wrap gap-2 mb-3">';
      html += '<button type="button" class="btn btn-primary" id="zt-metadata-migrate-btn"><i class="fa fa-exchange me-1"></i> Migrate thumbnails to metadata storage</button>';
      html += '<a class="btn btn-outline-secondary" href="/adm/tasks">View tasks</a>';
      html += '</div>';
      if (d.message) {
        html += '<div class="alert alert-info mb-0"><i class="fa fa-info-circle me-1"></i> ' + esc(d.message) + '</div>';
      }
      html += '</div></div>';
      self.innerHTML = html;
      var migrateBtn = self.querySelector("#zt-metadata-migrate-btn");
      if (migrateBtn) {
        migrateBtn.onclick = function() {
          if (!confirm("Queue migration of all legacy thumbnails to metadata storage? This may take a while.")) return;
          migrateBtn.disabled = true;
          fetch("/api/adm/metadata-storage/migrate", { method: "POST", credentials: "same-origin" })
            .then(function(r) { return r.json().then(function(data) { return { status: r.status, data: data }; }); })
            .then(function(res) {
              if (res.status === 202) {
                if (typeof sendToast === "function") sendToast("Queued", "", "bg-success", res.data.message || "Migration task queued");
                if (res.data.redirect && typeof loadPage === "function") loadPage(res.data.redirect);
              } else if (typeof sendToast === "function") {
                sendToast("Error", "", "bg-danger", (res.data && res.data.error) || "Migration failed");
                migrateBtn.disabled = false;
              }
            })
            .catch(function() {
              if (typeof sendToast === "function") sendToast("Error", "", "bg-danger", "Migration request failed");
              migrateBtn.disabled = false;
            });
        };
      }
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() {
      self.innerHTML = '<div class="alert alert-danger">Failed to load metadata storage configuration.</div>';
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    });
};
customElements.define("zt-adm-metadata-storage", ZtAdmMetadataStorage);
})();
