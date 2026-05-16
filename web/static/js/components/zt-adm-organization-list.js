(function() {
"use strict";
function esc(s) { return String(s == null ? "" : s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/"/g, "&quot;"); }
function ZtAdmOrganizationList() {
  var el = Reflect.construct(HTMLElement, [], ZtAdmOrganizationList);
  return el;
}
ZtAdmOrganizationList.prototype = Object.create(HTMLElement.prototype);
ZtAdmOrganizationList.prototype.connectedCallback = function() {
  var self = this;
  fetch("/api/adm/organizations", { credentials: "same-origin" })
    .then(function(r) { return r.json(); })
    .then(function(d) {
      var items = d.items || [];
      var reorganizeOnImport = !!d.reorganize_on_import;
      var html = '<div class="row"><div class="col-md-3 col-lg-3"><zt-adm-tabs data-active="organizations"></zt-adm-tabs></div><div class="col-md-9 col-lg-9"><div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-sitemap"></i></span><h3>Organization</h3><hr /></div>';
      html += '<p class="text-muted">An Organization defines the on-disk layout (path template) used when a new video is imported. Only one version is <em>Active</em> at a time. Existing videos keep their stored path until you trigger a reorganize.</p>';

      html += '<div class="alert alert-secondary d-flex flex-wrap align-items-center justify-content-between mb-3">'
        + '<div class="me-3"><strong>Reorganize on import</strong><br /><small class="text-muted">Global default: when off, freshly imported videos keep their original triage path. The upload UI can override this per import.</small></div>'
        + '<div class="form-check form-switch m-0"><input class="form-check-input" type="checkbox" role="switch" id="zt-reorg-on-import-toggle"' + (reorganizeOnImport ? " checked" : "") + '><label class="form-check-label" for="zt-reorg-on-import-toggle">' + (reorganizeOnImport ? "Enabled" : "Disabled") + '</label></div>'
        + '</div>';

      if (items.length === 0) {
        html += '<p class="text-muted mb-3">No organization yet. Use <strong>Add organization</strong> to create one.</p>';
      }

      html += '<div class="row row-cols-1 row-cols-sm-2 row-cols-lg-3 row-cols-xl-4 g-3 mb-3">';
      items.forEach(function(o) {
        var id = o.id || o.ID;
        var name = esc(o.name || o.Name);
        var tmpl = esc(o.template || o.Template);
        var isActive = !!(o.active || o.Active);
        var count = (o.video_count != null) ? o.video_count : 0;
        var activeBadge = isActive ? ' <span class="badge text-bg-primary">Active</span>' : "";
        var deleteBtn = isActive
          ? '<button type="button" class="btn btn-sm btn-danger" disabled title="Activate another organization first">Delete</button>'
          : '<button type="button" class="btn btn-sm btn-danger" data-org-id="' + esc(id) + '" data-org-name="' + name + '">Delete</button>';
        var activateBtn = isActive ? "" : '<button type="button" class="btn btn-sm btn-outline-primary" data-activate-id="' + esc(id) + '">Activate</button>';
        var reorganizeBtn = '<button type="button" class="btn btn-sm btn-outline-warning" data-reorganize-id="' + esc(id) + '" data-org-name="' + name + '" title="Move videos that follow another organization to this layout">Reorganize</button>';
        html += '<div class="col">' +
          '<div class="card h-100 shadow-sm">' +
          '<div class="card-body py-3 px-3 d-flex flex-column">' +
          '<div class="fw-semibold text-break mb-2">' + name + activeBadge + '</div>' +
          '<div class="small text-muted mb-2"><code class="user-select-all">' + tmpl + '</code></div>' +
          '<div class="small text-muted mb-3">' + count + ' video' + (count === 1 ? '' : 's') + '</div>' +
          '<div class="mt-auto d-flex flex-wrap gap-1 justify-content-end">' +
          '<button type="button" class="btn btn-sm btn-outline-primary" data-edit-id="' + esc(id) + '">Edit</button>' +
          activateBtn +
          reorganizeBtn +
          deleteBtn +
          "</div></div></div></div>";
      });
      html += '<div class="col">' +
        '<div class="card h-100 shadow-sm border border-2 text-secondary" id="zt-adm-org-add-card" role="button" tabindex="0" style="cursor:pointer;border-style:dashed;min-height:8.5rem">' +
        '<div class="card-body py-3 px-3 d-flex flex-column align-items-center justify-content-center text-center flex-grow-1">' +
        '<i class="fa fa-plus-circle fa-2x mb-2"></i>' +
        '<div class="fw-semibold">Add organization</div>' +
        '<div class="small mt-1">Define a new layout</div>' +
        "</div></div></div>";
      html += "</div></div></div>";

      var helpVars = ''
        + '<ul class="mb-0 small text-muted">'
        + '<li><code>$ID</code> &mdash; video UUID (required for unique paths)</li>'
        + '<li><code>$TYPE</code> &mdash; plural folder: clips, videos, movies</li>'
        + '<li><code>$TYPE_NAME</code> &mdash; clip, video, movie</li>'
        + '<li><code>$TYPE_LETTER</code> &mdash; c, v, m</li>'
        + '<li><code>$FILENAME</code> &mdash; original filename (with extension)</li>'
        + '<li><code>$BASENAME</code> &mdash; original filename without extension</li>'
        + '<li><code>$EXT</code> &mdash; file extension (e.g. <code>.mp4</code>)</li>'
        + '</ul>';

      html += '<div class="modal fade" id="zt-add-org-modal" tabindex="-1"><div class="modal-dialog modal-lg"><div class="modal-content"><div class="modal-header"><h5 class="modal-title">Add organization</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div><div class="modal-body"><form id="zt-adm-org-form">'
        + '<div class="mb-2"><label class="form-label">Name</label><input type="text" class="form-control" name="name" required placeholder="e.g. v2 - flat layout"></div>'
        + '<div class="mb-2"><label class="form-label">Path template</label><input type="text" class="form-control" name="template" required placeholder="$TYPE/$ID/video.mp4"><div class="form-text">' + helpVars + '</div></div>'
        + '<div class="mb-2"><div class="form-check"><input type="checkbox" class="form-check-input" name="active" id="zt-org-active"><label class="form-check-label" for="zt-org-active">Set as active (use for new imports)</label></div></div>'
        + '</form></div><div class="modal-footer"><button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button><button type="submit" class="btn btn-primary" form="zt-adm-org-form">Add organization</button></div></div></div></div>';

      html += '<div class="modal fade" id="zt-edit-org-modal" tabindex="-1"><div class="modal-dialog modal-lg"><div class="modal-content"><div class="modal-header"><h5 class="modal-title">Edit organization</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div><div class="modal-body"><form id="zt-adm-org-edit-form">'
        + '<input type="hidden" id="zt-edit-org-id">'
        + '<div class="mb-2"><label class="form-label">Name</label><input type="text" class="form-control" id="zt-edit-org-name" required></div>'
        + '<div class="mb-2"><label class="form-label">Path template</label><input type="text" class="form-control" id="zt-edit-org-template" required><div class="form-text" id="zt-edit-org-template-help"><span class="text-warning small d-none" id="zt-edit-org-template-locked">Locked: organization is already used by existing videos. Create a new organization instead.</span>' + helpVars + '</div></div>'
        + '<div class="mb-2"><div class="form-check"><input type="checkbox" class="form-check-input" id="zt-edit-org-active"><label class="form-check-label" for="zt-edit-org-active">Set as active</label></div></div>'
        + '</form></div><div class="modal-footer"><button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button><button type="button" class="btn btn-primary" id="zt-edit-org-submit">Save</button></div></div></div></div>';

      self.innerHTML = html;

      var addCard = self.querySelector("#zt-adm-org-add-card");
      var addFormEl = self.querySelector("#zt-adm-org-form");
      function openAddModal() {
        if (addFormEl) addFormEl.reset();
        var addModalEl = self.querySelector("#zt-add-org-modal");
        if (addModalEl) {
          bootstrap.Modal.getOrCreateInstance(addModalEl).show();
          setTimeout(function() {
            var nameInput = addFormEl && addFormEl.querySelector('input[name="name"]');
            if (nameInput) nameInput.focus();
          }, 400);
        }
      }
      if (addCard) {
        addCard.addEventListener("click", openAddModal);
        addCard.addEventListener("keydown", function(e) {
          if (e.key === "Enter" || e.key === " ") { e.preventDefault(); openAddModal(); }
        });
      }

      if (addFormEl) {
        addFormEl.addEventListener("submit", function(e) {
          e.preventDefault();
          var name = addFormEl.querySelector('input[name="name"]').value.trim();
          var template = addFormEl.querySelector('input[name="template"]').value.trim();
          var active = addFormEl.querySelector('input[name="active"]').checked;
          var payload = { name: name, template: template, active: active };
          fetch("/api/adm/organizations", {
            method: "POST",
            credentials: "same-origin",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload)
          }).then(function(r) { return r.json().then(function(data) { return { status: r.status, data: data }; }); })
            .then(function(res) {
              if (res.status === 201) {
                var addModalEl = self.querySelector("#zt-add-org-modal");
                var addInst = addModalEl && bootstrap.Modal.getInstance(addModalEl);
                if (addInst) addInst.hide();
                if (typeof loadPage === "function") loadPage("/adm/organizations");
                else window.location.reload();
              } else if (typeof sendToast === "function") sendToast("Error", "", "bg-danger", (res.data && res.data.error) || "Failed to add organization");
            });
        });
      }

      self.querySelectorAll("button[data-org-id]").forEach(function(btn) {
        btn.onclick = function() {
          var name = btn.getAttribute("data-org-name") || "this organization";
          if (!confirm("Delete \"" + name + "\"? This is only allowed if it is not active and has no videos.")) return;
          var id = btn.getAttribute("data-org-id");
          fetch("/api/adm/organizations/" + encodeURIComponent(id), { method: "DELETE", credentials: "same-origin" })
            .then(function(r) {
              if (r.status === 204) {
                if (typeof loadPage === "function") loadPage("/adm/organizations");
                else window.location.reload();
              } else r.json().then(function(data) { if (typeof sendToast === "function") sendToast("Error", "", "bg-danger", (data && data.error) || "Delete failed"); });
            });
        };
      });

      self.querySelectorAll("button[data-activate-id]").forEach(function(btn) {
        btn.onclick = function() {
          var id = btn.getAttribute("data-activate-id");
          fetch("/api/adm/organizations/" + encodeURIComponent(id) + "/activate", { method: "POST", credentials: "same-origin" })
            .then(function(r) { return r.json().then(function(d) { return { status: r.status, data: d }; }); })
            .then(function(res) {
              if (res.status === 200) {
                if (typeof loadPage === "function") loadPage("/adm/organizations");
                else window.location.reload();
              } else if (typeof sendToast === "function") sendToast("Error", "", "bg-danger", (res.data && res.data.error) || "Activation failed");
            });
        };
      });

      self.querySelectorAll("button[data-reorganize-id]").forEach(function(btn) {
        btn.onclick = function() {
          var id = btn.getAttribute("data-reorganize-id");
          var name = btn.getAttribute("data-org-name") || "this organization";
          if (!confirm("Queue reorganize tasks moving all videos onto \"" + name + "\"? Files will be copied and the old ones removed in the background.")) return;
          fetch("/api/adm/organizations/" + encodeURIComponent(id) + "/reorganize", { method: "POST", credentials: "same-origin" })
            .then(function(r) { return r.json().then(function(d) { return { status: r.status, data: d }; }); })
            .then(function(res) {
              if (res.status === 202 || res.status === 200) {
                if (typeof sendToast === "function") sendToast("Reorganize", "", "bg-success", (res.data && res.data.message) || "Reorganize queued");
              } else if (typeof sendToast === "function") sendToast("Error", "", "bg-danger", (res.data && res.data.error) || "Reorganize failed");
            });
        };
      });

      self.querySelectorAll("button[data-edit-id]").forEach(function(btn) {
        btn.onclick = function() {
          var id = btn.getAttribute("data-edit-id");
          var org = items.find(function(o) { return (o.id || o.ID) === id; });
          if (!org) return;
          self.querySelector("#zt-edit-org-id").value = id;
          self.querySelector("#zt-edit-org-name").value = org.name || org.Name || "";
          var tmplInput = self.querySelector("#zt-edit-org-template");
          tmplInput.value = org.template || org.Template || "";
          var locked = (org.video_count || 0) > 0;
          tmplInput.disabled = locked;
          self.querySelector("#zt-edit-org-template-locked").classList.toggle("d-none", !locked);
          self.querySelector("#zt-edit-org-active").checked = !!(org.active || org.Active);
          var modalEl = self.querySelector("#zt-edit-org-modal");
          new bootstrap.Modal(modalEl).show();
        };
      });

      var editSubmit = self.querySelector("#zt-edit-org-submit");
      if (editSubmit) {
        editSubmit.onclick = function() {
          var id = self.querySelector("#zt-edit-org-id").value;
          if (!id) return;
          var payload = {
            name: self.querySelector("#zt-edit-org-name").value.trim(),
            active: self.querySelector("#zt-edit-org-active").checked
          };
          var tmplInput = self.querySelector("#zt-edit-org-template");
          if (!tmplInput.disabled) payload.template = tmplInput.value.trim();
          fetch("/api/adm/organizations/" + encodeURIComponent(id), {
            method: "PUT",
            credentials: "same-origin",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload)
          }).then(function(r) { return r.json().then(function(d) { return { status: r.status, data: d }; }); })
            .then(function(res) {
              if (res.status === 200) {
                bootstrap.Modal.getInstance(self.querySelector("#zt-edit-org-modal")).hide();
                if (typeof loadPage === "function") loadPage("/adm/organizations");
                else window.location.reload();
              } else if (typeof sendToast === "function") sendToast("Error", "", "bg-danger", (res.data && res.data.error) || "Update failed");
            });
        };
      }

      var toggle = self.querySelector("#zt-reorg-on-import-toggle");
      if (toggle) {
        toggle.addEventListener("change", function() {
          var action = toggle.checked ? "enable" : "disable";
          fetch("/api/adm/config/reorganize-on-import/" + action, { credentials: "same-origin" })
            .then(function(r) { return r.json().then(function(d) { return { status: r.status, data: d }; }); })
            .then(function(res) {
              if (res.status === 200) {
                var label = toggle.parentElement.querySelector("label");
                if (label) label.textContent = toggle.checked ? "Enabled" : "Disabled";
              } else if (typeof sendToast === "function") sendToast("Error", "", "bg-danger", (res.data && res.data.error) || "Update failed");
            });
        });
      }

      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() {
      self.innerHTML = '<div class="alert alert-danger">Failed to load.</div>';
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    });
};
customElements.define("zt-adm-organization-list", ZtAdmOrganizationList);
})();
