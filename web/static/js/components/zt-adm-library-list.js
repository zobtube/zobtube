(function() {
"use strict";
function esc(s) { return String(s == null ? "" : s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/"/g, "&quot;"); }
function ZtAdmLibraryList() {
  var el = Reflect.construct(HTMLElement, [], ZtAdmLibraryList);
  return el;
}
ZtAdmLibraryList.prototype = Object.create(HTMLElement.prototype);
ZtAdmLibraryList.prototype.connectedCallback = function() {
  var self = this;
  fetch("/api/adm/libraries", { credentials: "same-origin" })
    .then(function(r) { return r.json(); })
    .then(function(d) {
      var items = d.items || [];
      var html = '<div class="row"><div class="col-md-3 col-lg-3"><zt-adm-tabs data-active="libraries"></zt-adm-tabs></div><div class="col-md-9 col-lg-9"><div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-folder-open"></i></span><h3>Libraries</h3><hr /></div>';
      html += '<p class="text-muted">Libraries define where media is stored (filesystem or S3). The default library is used for actor/channel/category thumbnails and as the default upload target.</p>';
      if (items.length === 0) {
        html += '<p class="text-muted">No libraries. Use the form below to add one.</p>';
      } else {
        html += '<table class="table table-striped table-hover"><thead><tr><th>Name</th><th>Type</th><th>Config</th><th>Default</th><th></th></tr></thead><tbody>';
        items.forEach(function(lib) {
          var id = lib.id || lib.ID;
          var name = esc(lib.name || lib.Name);
          var typeStr = esc(lib.type || lib.Type || "—");
          var isDefault = !!(lib.is_default || lib.IsDefault);
          var defaultBadge = isDefault ? '<span class="badge bg-primary">default</span>' : '';
          var configStr = "—";
          if (lib.config && lib.config.filesystem && lib.config.filesystem.path) {
            configStr = esc(lib.config.filesystem.path);
          } else if (lib.config && lib.config.s3) {
            var s3 = lib.config.s3;
            configStr = esc((s3.bucket || s3.Bucket) + " / " + (s3.region || s3.Region || ""));
          }
          var deleteBtn = isDefault
            ? '<button type="button" class="btn btn-sm btn-danger" disabled title="The default library cannot be deleted">Delete</button>'
            : '<button type="button" class="btn btn-sm btn-danger" data-library-id="' + esc(id) + '" data-library-name="' + name + '">Delete</button>';
          html += '<tr><td>' + name + '</td><td>' + typeStr + '</td><td><code>' + configStr + '</code></td><td>' + defaultBadge + '</td><td style="text-align:end">' + deleteBtn + '</td></tr>';
        });
        html += '</tbody></table>';
      }
      html += '<hr /><h5>Add library</h5><form id="zt-adm-library-form" class="mb-4"><div class="mb-2"><label class="form-label">Name</label><input type="text" class="form-control" name="name" required placeholder="e.g. Local SSD"></div><div class="mb-2"><label class="form-label">Type</label><select class="form-select" name="type"><option value="filesystem">Filesystem</option><option value="s3">S3</option></select></div><div id="zt-adm-library-fs" class="mb-2"><label class="form-label">Path</label><input type="text" class="form-control" name="path" placeholder="/path/to/media"></div><div id="zt-adm-library-s3" class="mb-2" style="display:none"><label class="form-label">Bucket</label><input type="text" class="form-control" name="bucket" placeholder="my-bucket"><label class="form-label mt-1">Region</label><input type="text" class="form-control" name="region" placeholder="us-east-1"><label class="form-label mt-1">Prefix (optional)</label><input type="text" class="form-control" name="prefix" placeholder="media/"><label class="form-label mt-1">Endpoint (optional, for Minio)</label><input type="text" class="form-control" name="endpoint" placeholder="http://localhost:9000"></div><div class="mb-2"><div class="form-check"><input type="checkbox" class="form-check-input" name="default" id="zt-lib-default"><label class="form-check-label" for="zt-lib-default">Set as default library</label></div></div><button type="submit" class="btn btn-primary">Add library</button></form>';
      html += '</div></div>';
      self.innerHTML = html;
      var typeSel = self.querySelector("select[name=type]");
      var fsDiv = self.querySelector("#zt-adm-library-fs");
      var s3Div = self.querySelector("#zt-adm-library-s3");
      function toggleConfig() {
        var t = typeSel.value;
        fsDiv.style.display = t === "filesystem" ? "block" : "none";
        s3Div.style.display = t === "s3" ? "block" : "none";
      }
      typeSel.addEventListener("change", toggleConfig);
      toggleConfig();
      self.querySelector("#zt-adm-library-form").addEventListener("submit", function(e) {
        e.preventDefault();
        var form = e.target;
        var name = form.name.value.trim();
        var typeVal = form.type.value;
        var payload = { name: name, type: typeVal, config: {} };
        if (typeVal === "filesystem") {
          payload.config = { filesystem: { path: form.path.value.trim() } };
        } else {
          payload.config = { s3: { bucket: form.bucket.value.trim(), region: form.region.value.trim() || "us-east-1", prefix: form.prefix.value.trim() || "", endpoint: form.endpoint.value.trim() || "" } };
        }
        payload.default = form.querySelector("#zt-lib-default").checked;
        fetch("/api/adm/libraries", { method: "POST", credentials: "same-origin", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload) })
          .then(function(r) { return r.json().then(function(data) { return { status: r.status, data: data }; }); })
          .then(function(res) {
            if (res.status === 201 && typeof loadPage === "function") loadPage("/adm/libraries");
            else if (res.status === 201) window.location.reload();
            else if (typeof sendToast === "function") sendToast("Error", "", "bg-danger", (res.data && res.data.error) || "Failed to add library");
          });
      });
      self.querySelectorAll("button[data-library-id]").forEach(function(btn) {
        btn.onclick = function() {
          if (!confirm("Delete library \"" + btn.getAttribute("data-library-name") + "\"? This is only allowed if it has no videos and is not the default.")) return;
          var id = btn.getAttribute("data-library-id");
          fetch("/api/adm/libraries/" + encodeURIComponent(id), { method: "DELETE", credentials: "same-origin" })
            .then(function(r) {
              if (r.status === 204 && typeof loadPage === "function") loadPage("/adm/libraries");
              else if (r.status === 204) window.location.reload();
              else return r.json().then(function(data) { if (typeof sendToast === "function") sendToast("Error", "", "bg-danger", (data && data.error) || "Delete failed"); });
            });
        };
      });
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() { self.innerHTML = '<div class="alert alert-danger">Failed to load.</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); });
};
customElements.define("zt-adm-library-list", ZtAdmLibraryList);
})();
