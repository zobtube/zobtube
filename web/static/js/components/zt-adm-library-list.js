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
          html += '<tr><td>' + name + '</td><td>' + typeStr + '</td><td><code>' + configStr + '</code></td><td>' + defaultBadge + '</td><td style="text-align:end"><button type="button" class="btn btn-sm btn-outline-primary me-1" data-edit-library-id="' + esc(id) + '">Edit</button>' + deleteBtn + '</td></tr>';
        });
        html += '</tbody></table>';
      }
      html += '<hr /><h5>Add library</h5><form id="zt-adm-library-form" class="mb-4"><div class="mb-2"><label class="form-label">Name</label><input type="text" class="form-control" name="name" required placeholder="e.g. Local SSD"></div><div class="mb-2"><label class="form-label">Type</label><select class="form-select" name="type"><option value="filesystem">Filesystem</option><option value="s3">S3</option></select></div><div id="zt-adm-library-fs" class="mb-2"><label class="form-label">Path</label><input type="text" class="form-control" name="path" placeholder="/path/to/media"></div><div id="zt-adm-library-s3" class="mb-2" style="display:none"><label class="form-label">Bucket</label><input type="text" class="form-control" name="bucket" placeholder="my-bucket"><label class="form-label mt-1">Region</label><input type="text" class="form-control" name="region" placeholder="us-east-1"><label class="form-label mt-1">Prefix (optional)</label><input type="text" class="form-control" name="prefix" placeholder="media/"><label class="form-label mt-1">Endpoint (optional, for Minio)</label><input type="text" class="form-control" name="endpoint" placeholder="http://localhost:9000"><label class="form-label mt-1">Access Key ID (optional)</label><input type="text" class="form-control" name="access_key_id" placeholder="Leave empty for env/IAM"><label class="form-label mt-1">Secret Access Key (optional)</label><input type="password" class="form-control" name="secret_access_key" placeholder="Leave empty for env/IAM" autocomplete="new-password"></div><div class="mb-2"><div class="form-check"><input type="checkbox" class="form-check-input" name="default" id="zt-lib-default"><label class="form-check-label" for="zt-lib-default">Set as default library</label></div></div><button type="submit" class="btn btn-primary">Add library</button></form>';
      html += '</div></div>';
      html += '<div class="modal fade" id="zt-edit-library-modal" tabindex="-1"><div class="modal-dialog"><div class="modal-content"><div class="modal-header"><h5 class="modal-title">Edit library</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div><div class="modal-body"><form id="zt-adm-library-edit-form"><input type="hidden" id="zt-edit-library-id"><div class="mb-2"><label class="form-label">Name</label><input type="text" class="form-control" id="zt-edit-name" required></div><div class="mb-2"><label class="form-label">Type</label><select class="form-select" id="zt-edit-type"><option value="filesystem">Filesystem</option><option value="s3">S3</option></select></div><div id="zt-edit-fs" class="mb-2"><label class="form-label">Path</label><input type="text" class="form-control" id="zt-edit-path" placeholder="/path/to/media"></div><div id="zt-edit-s3" class="mb-2" style="display:none"><label class="form-label">Bucket</label><input type="text" class="form-control" id="zt-edit-bucket" placeholder="my-bucket"><label class="form-label mt-1">Region</label><input type="text" class="form-control" id="zt-edit-region" placeholder="us-east-1"><label class="form-label mt-1">Prefix (optional)</label><input type="text" class="form-control" id="zt-edit-prefix" placeholder="media/"><label class="form-label mt-1">Endpoint (optional, for Minio)</label><input type="text" class="form-control" id="zt-edit-endpoint" placeholder="http://localhost:9000"><label class="form-label mt-1">Access Key ID (optional)</label><input type="text" class="form-control" id="zt-edit-access-key-id" placeholder="Leave empty for env/IAM"><label class="form-label mt-1">Secret Access Key (optional)</label><input type="password" class="form-control" id="zt-edit-secret-access-key" placeholder="Leave blank to keep current" autocomplete="new-password"></div><div class="mb-2"><div class="form-check"><input type="checkbox" class="form-check-input" id="zt-edit-default"><label class="form-check-label" for="zt-edit-default">Set as default library</label></div></div></form></div><div class="modal-footer"><button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button><button type="button" class="btn btn-primary" id="zt-edit-library-submit">Save</button></div></div></div></div>';
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
          payload.config = { s3: { bucket: form.bucket.value.trim(), region: form.region.value.trim() || "us-east-1", prefix: form.prefix.value.trim() || "", endpoint: form.endpoint.value.trim() || "", access_key_id: (form.querySelector('input[name="access_key_id"]') || {}).value.trim() || undefined, secret_access_key: (form.querySelector('input[name="secret_access_key"]') || {}).value.trim() || undefined } };
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

      self.querySelectorAll("button[data-edit-library-id]").forEach(function(btn) {
        btn.onclick = function() {
          var id = btn.getAttribute("data-edit-library-id");
          var lib = items.find(function(l) { return (l.id || l.ID) === id; });
          if (!lib) return;
          var config = lib.config || lib.Config || {};
          var fsConfig = config.filesystem || config.Filesystem;
          var s3Config = config.s3 || config.S3;
          self.querySelector("#zt-edit-library-id").value = id;
          self.querySelector("#zt-edit-name").value = lib.name || lib.Name || "";
          var typeVal = (lib.type || lib.Type || "filesystem").toLowerCase();
          if (typeVal !== "filesystem" && typeVal !== "s3") typeVal = "filesystem";
          var editTypeSel = self.querySelector("#zt-edit-type");
          editTypeSel.value = typeVal;
          var editFs = self.querySelector("#zt-edit-fs");
          var editS3 = self.querySelector("#zt-edit-s3");
          editFs.style.display = typeVal === "filesystem" ? "block" : "none";
          editS3.style.display = typeVal === "s3" ? "block" : "none";
          self.querySelector("#zt-edit-path").value = (fsConfig && (fsConfig.path || fsConfig.Path)) ? (fsConfig.path || fsConfig.Path || "") : "";
          if (s3Config) {
            self.querySelector("#zt-edit-bucket").value = s3Config.bucket || s3Config.Bucket || "";
            self.querySelector("#zt-edit-region").value = s3Config.region || s3Config.Region || "us-east-1";
            self.querySelector("#zt-edit-prefix").value = s3Config.prefix || s3Config.Prefix || "";
            self.querySelector("#zt-edit-endpoint").value = s3Config.endpoint || s3Config.Endpoint || "";
            self.querySelector("#zt-edit-access-key-id").value = s3Config.access_key_id || s3Config.AccessKeyID || "";
            self.querySelector("#zt-edit-secret-access-key").value = "";
          } else {
            self.querySelector("#zt-edit-bucket").value = "";
            self.querySelector("#zt-edit-region").value = "us-east-1";
            self.querySelector("#zt-edit-prefix").value = "";
            self.querySelector("#zt-edit-endpoint").value = "";
            self.querySelector("#zt-edit-access-key-id").value = "";
            self.querySelector("#zt-edit-secret-access-key").value = "";
          }
          self.querySelector("#zt-edit-default").checked = !!(lib.is_default || lib.IsDefault);
          var modalEl = self.querySelector("#zt-edit-library-modal");
          var modal = new bootstrap.Modal(modalEl);
          modal.show();
        };
      });

      var editTypeSel = self.querySelector("#zt-edit-type");
      if (editTypeSel) {
        editTypeSel.addEventListener("change", function() {
          var t = editTypeSel.value;
          self.querySelector("#zt-edit-fs").style.display = t === "filesystem" ? "block" : "none";
          self.querySelector("#zt-edit-s3").style.display = t === "s3" ? "block" : "none";
        });
      }

      var editSubmitBtn = self.querySelector("#zt-edit-library-submit");
      if (editSubmitBtn) {
        editSubmitBtn.onclick = function() {
          var id = self.querySelector("#zt-edit-library-id").value;
          if (!id) return;
          var name = self.querySelector("#zt-edit-name").value.trim();
          if (!name) return;
          var typeVal = self.querySelector("#zt-edit-type").value;
          var payload = { name: name, type: typeVal };
          if (typeVal === "filesystem") {
            payload.config = { filesystem: { path: self.querySelector("#zt-edit-path").value.trim() } };
          } else if (typeVal === "s3") {
            var s3Payload = { bucket: self.querySelector("#zt-edit-bucket").value.trim(), region: self.querySelector("#zt-edit-region").value.trim() || "us-east-1", prefix: self.querySelector("#zt-edit-prefix").value.trim() || "", endpoint: self.querySelector("#zt-edit-endpoint").value.trim() || "", access_key_id: self.querySelector("#zt-edit-access-key-id").value.trim() || undefined };
            var secretVal = self.querySelector("#zt-edit-secret-access-key").value;
            if (secretVal) s3Payload.secret_access_key = secretVal;
            payload.config = { s3: s3Payload };
          }
          payload.default = self.querySelector("#zt-edit-default").checked;
          fetch("/api/adm/libraries/" + encodeURIComponent(id), { method: "PUT", credentials: "same-origin", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload) })
            .then(function(r) { return r.json().then(function(data) { return { status: r.status, data: data }; }); })
            .then(function(res) {
              if (res.status === 200) {
                bootstrap.Modal.getInstance(self.querySelector("#zt-edit-library-modal")).hide();
                if (typeof loadPage === "function") loadPage("/adm/libraries");
                else window.location.reload();
              } else if (typeof sendToast === "function") sendToast("Error", "", "bg-danger", (res.data && res.data.error) || "Update failed");
            });
        };
      }

      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() { self.innerHTML = '<div class="alert alert-danger">Failed to load.</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); });
};
customElements.define("zt-adm-library-list", ZtAdmLibraryList);
})();
