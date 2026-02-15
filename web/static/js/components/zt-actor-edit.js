(function() {
"use strict";
function esc(s) { return String(s).replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/"/g,"&quot;"); }
function ZtActorEdit() {
  var el = Reflect.construct(HTMLElement, [], ZtActorEdit);
  return el;
}
ZtActorEdit.prototype = Object.create(HTMLElement.prototype);
ZtActorEdit.prototype.connectedCallback = function() {
  var self = this;
  var id = this.getAttribute("data-id");
  if (!id) { self.innerHTML = '<div class="alert alert-danger">Missing id</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); return; }
  if (!(window.__USER__ && window.__USER__.admin)) {
    self.innerHTML = '<div class="alert alert-danger">Forbidden</div>';
    if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    return;
  }
  Promise.all([
    fetch("/api/actor/" + encodeURIComponent(id), { credentials: "same-origin" }).then(function(r) { if (!r.ok) throw new Error(r.status); return r.json(); }),
    fetch("/api/category", { credentials: "same-origin" }).then(function(r) { if (!r.ok) throw new Error(r.status); return r.json(); })
  ]).then(function(results) {
    var a = results[0];
    var catsData = results[1];
    var categories = catsData.items || catsData.categories || [];
    var urlView = "/actor/" + id;
    var urlThumb = "/api/actor/" + id + "/thumb";
    var name = esc(a.Name || a.name || "");
    var sex = a.Sex || a.sex || "m";
    var desc = esc(a.Description || a.description || "");
    var aliases = a.Aliases || a.aliases || [];
    var actorCats = a.Categories || a.categories || [];
    var catSelectedIds = actorCats.map(function(c) { return c.ID || c.id; });
    var categorySelectable = {};
    categories.forEach(function(c) {
      (c.Sub || c.sub || []).forEach(function(s) { categorySelectable[s.ID || s.id] = s.Name || s.name || ""; });
    });

    function refresh() {
      if (typeof loadPage === "function") loadPage("/actor/" + id + "/edit");
      else window.location.reload();
    }

    var categoryChips = [];
    categories.forEach(function(c) {
      (c.Sub || c.sub || []).forEach(function(s) {
        var sid = s.ID || s.id;
        var show = catSelectedIds.indexOf(sid) >= 0 ? "" : "display:none";
        var thumb = (s.Thumbnail || s.thumbnail) ? '<img src="/api/category-sub/' + encodeURIComponent(sid) + '/thumb" width="100" height="50">' : "";
        categoryChips.push('<div class="chip actor-category-list" category-id="' + sid + '" style="' + show + '">' + thumb + esc(s.Name || s.name || "") + '<button class="btn btn-danger"><i class="fa fa-trash-alt"></i></button></div>');
      });
    });
    var categoryModalHtml = "";
    categories.forEach(function(c) {
      var subs = c.Sub || c.sub || [];
      if (subs.length === 0) return;
      categoryModalHtml += '<h4 class="mt-3">' + esc(c.Name || c.name) + '</h4><div class="chips">';
      subs.forEach(function(s) {
        var sid = s.ID || s.id;
        var sel = catSelectedIds.indexOf(sid) >= 0;
        var addStyle = sel ? "display:none" : "";
        var remStyle = sel ? "" : "display:none";
        var thumb = (s.Thumbnail || s.thumbnail) ? '<img class="lazy" data-src="/api/category-sub/' + encodeURIComponent(sid) + '/thumb" width="100" height="50">' : "";
        categoryModalHtml += '<div class="chip add-category-list" category-id="' + sid + '">' + thumb + esc(s.Name || s.name) + '<button class="btn btn-success add-category-add" style="' + addStyle + '"><i class="fa fa-plus-circle"></i></button><button class="btn btn-danger add-category-remove" style="' + remStyle + '"><i class="fa fa-trash-alt"></i></button></div>';
      });
      categoryModalHtml += "</div>";
    });

    var aliasChips = aliases.map(function(al) {
      var aid = al.ID || al.id;
      return '<div class="chip" alias-id="' + aid + '">' + esc(al.Name || al.name || "") + '<button class="btn btn-danger zt-alias-remove" data-id="' + aid + '"><i class="fa fa-trash-alt"></i></button></div>';
    }).join("");

    var html = '<h2>Edit actor information</h2><hr/><br/><div class="row">';
    html += '<div class="col-3"><img class="rounded" src="' + urlThumb + '" style="width:100%">';
    html += '<div style="margin-top:15px"><a class="btn btn-success" style="width:100%" href="' + urlView + '">View profile</a>';
    html += '<a class="btn btn-danger" style="margin-top:15px;width:100%" href="/actor/' + id + '/delete" data-full-reload>Delete actor profile</a></div></div>';
    html += '<div class="col-9"><h3>Profile details</h3><br/>';
    html += '<div class="mb-3"><div class="form-floating input-group"><input type="text" class="form-control" id="actor-name" disabled value="' + name + '"><label for="actor-name">Name</label><button class="btn btn-outline-warning" type="button" id="actor-name-edit">Edit</button></div></div>';
    html += '<div class="mb-3"><label class="form-label">Sex:</label><select class="form-select" id="actor-sex" disabled><option value="m"' + (sex==="m"?" selected":"") + '>Male</option><option value="f"' + (sex==="f"?" selected":"") + '>Female</option><option value="tw"' + (sex==="tw"?" selected":"") + '>Trans women</option></select></div>';
    html += '<div class="mb-3"><label class="form-label">Description</label><textarea class="form-control" id="actor-description" readonly disabled>' + desc + '</textarea><button class="btn btn-outline-warning btn-sm mt-1" id="actor-desc-edit">Edit</button></div>';
    html += '<h4 class="mt-3 mb-3">Aliases</h4><div class="form-control chip-selector" style="height:unset;display:flex"><div class="chips">' + aliasChips;
    html += '<div class="chip">Add an alias<button class="btn btn-success" data-bs-toggle="modal" data-bs-target="#addActorAliasModal"><i class="fa fa-plus-circle"></i></button></div></div></div>';
    html += '<h4 class="mt-3">Categories</h4><div class="form-control chip-selector" style="height:unset;display:flex"><div class="chips">' + categoryChips.join("") + '<div class="chip">Add a category<button class="btn btn-success" data-bs-toggle="modal" data-bs-target="#addCategoryModal"><i class="fa fa-plus-circle"></i></button></div></div></div></div></div>';

    html += '<div class="modal fade" id="addActorAliasModal" tabindex="-1"><div class="modal-dialog"><div class="modal-content"><div class="modal-header"><h5 class="modal-title">Add alias</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div><div class="modal-body"><div class="form-floating"><input type="text" class="form-control" id="new-alias-input" placeholder="Alias"><label for="new-alias-input">Alias</label></div></div><div class="modal-footer"><button type="button" class="btn btn-primary" id="zt-alias-add">Add</button><button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button></div></div></div></div>';
    html += '<div class="modal fade" id="addCategoryModal" tabindex="-1"><div class="modal-dialog modal-xl"><div class="modal-content"><div class="modal-header"><h5 class="modal-title">Add category</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div><div class="modal-body">' + categoryModalHtml + '</div><div class="modal-footer"><button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button></div></div></div></div>';

    self.innerHTML = html;

    var nameEl = self.querySelector("#actor-name");
    var nameBtn = self.querySelector("#actor-name-edit");
    var nameEdit = false;
    nameBtn.addEventListener("click", function() {
      if (!nameEdit) { nameEl.disabled = false; nameBtn.textContent = "Send"; nameEdit = true; }
      else {
        var fd = new FormData(); fd.set("name", nameEl.value);
        fetch("/api/actor/" + id + "/rename", { method: "POST", credentials: "same-origin", body: fd })
          .then(function(r) { if (r.ok) { nameEl.disabled = true; nameBtn.textContent = "Edit"; nameEdit = false; } });
      }
    });

    var descEl = self.querySelector("#actor-description");
    var descBtn = self.querySelector("#actor-desc-edit");
    var descEdit = false;
    descBtn.addEventListener("click", function() {
      if (!descEdit) { descEl.disabled = false; descEl.readOnly = false; descBtn.textContent = "Save"; descEdit = true; }
      else {
        var fd = new FormData(); fd.set("description", descEl.value);
        fetch("/api/actor/" + id + "/description", { method: "POST", credentials: "same-origin", body: fd })
          .then(function(r) {
            if (r.ok) { descEl.disabled = true; descEl.readOnly = true; descBtn.textContent = "Edit"; descEdit = false; if (typeof sendToast === "function") sendToast("Actor description updated", "", "bg-success", "Description changed for " + (nameEl.value || "")); }
            else r.json().catch(function(){return{};}).then(function(d){ if (typeof sendToast === "function") sendToast("Actor description update failed", "", "bg-danger", (d && d.error) || "Unexpected error during actor description update"); });
          });
      }
    });

    self.querySelector("#zt-alias-add").addEventListener("click", function() {
      var val = self.querySelector("#new-alias-input").value;
      if (!val.trim()) return;
      var fd = new FormData(); fd.set("alias", val.trim());
      fetch("/api/actor/" + id + "/alias", { method: "POST", credentials: "same-origin", body: fd })
        .then(function(r) {
          if (r.ok) { bootstrap.Modal.getInstance(self.querySelector("#addActorAliasModal")).hide(); refresh(); if (typeof sendToast === "function") sendToast("Actor's new alias", "", "bg-success", "New alias recorded!"); }
          else r.json().catch(function(){return{};}).then(function(d){ if (typeof sendToast === "function") sendToast("Actor's new alias", "", "bg-warning", (d && d.error) || "Unable to create this new alias"); });
        });
    });

    self.querySelectorAll(".zt-alias-remove").forEach(function(btn) {
      btn.addEventListener("click", function() {
        var aid = btn.getAttribute("data-id");
        fetch("/api/actor/alias/" + aid, { method: "DELETE", credentials: "same-origin" })
          .then(function(r) {
            if (r.ok) { refresh(); if (typeof sendToast === "function") sendToast("Actor's alias removal", "", "bg-success", "Successfully removed!"); }
            else r.json().catch(function(){return{};}).then(function(d){ if (typeof sendToast === "function") sendToast("Actor's alias removal", "", "bg-warning", (d && d.error) || "Unable to remove alias"); });
          });
      });
    });

    self.querySelectorAll(".actor-category-list .btn-danger").forEach(function(btn) {
      btn.addEventListener("click", function() {
        var cid = btn.closest(".chip").getAttribute("category-id");
        var cname = categorySelectable[cid] || "Category";
        fetch("/api/actor/" + id + "/category/" + cid, { method: "DELETE", credentials: "same-origin" })
          .then(function(r) {
            if (r.ok) { refresh(); if (typeof sendToast === "function") sendToast("Category removed", "", "bg-success", cname + " removed."); }
            else r.json().catch(function(){return{};}).then(function(d){ if (typeof sendToast === "function") sendToast("Category not removed", "", "bg-danger", (d && d.error) || cname + " not removed, call failed."); });
          });
      });
    });
    self.querySelectorAll(".add-category-list .add-category-add").forEach(function(btn) {
      btn.addEventListener("click", function() {
        var cid = btn.closest(".add-category-list").getAttribute("category-id");
        var cname = categorySelectable[cid] || "Category";
        fetch("/api/actor/" + id + "/category/" + cid, { method: "PUT", credentials: "same-origin" })
          .then(function(r) {
            if (r.ok) { refresh(); if (typeof sendToast === "function") sendToast("Category added", "", "bg-success", cname + " added."); }
            else r.json().catch(function(){return{};}).then(function(d){ if (typeof sendToast === "function") sendToast("Category not added", "", "bg-danger", (d && d.error) || cname + " not added, call failed."); });
          });
      });
    });
    self.querySelectorAll(".add-category-list .add-category-remove").forEach(function(btn) {
      btn.addEventListener("click", function() {
        var cid = btn.closest(".add-category-list").getAttribute("category-id");
        var cname = categorySelectable[cid] || "Category";
        fetch("/api/actor/" + id + "/category/" + cid, { method: "DELETE", credentials: "same-origin" })
          .then(function(r) {
            if (r.ok) { refresh(); if (typeof sendToast === "function") sendToast("Category removed", "", "bg-success", cname + " removed."); }
            else r.json().catch(function(){return{};}).then(function(d){ if (typeof sendToast === "function") sendToast("Category not removed", "", "bg-danger", (d && d.error) || cname + " not removed, call failed."); });
          });
      });
    });

    if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
  }).catch(function(e) {
    if (e && e.message === "403") self.innerHTML = '<div class="alert alert-danger">Forbidden</div>';
    else if (e && e.message === "404") self.innerHTML = '<div class="alert alert-danger">Not found</div>';
    else self.innerHTML = '<div class="alert alert-danger">Failed to load.</div>';
    if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
  });
};
customElements.define("zt-actor-edit", ZtActorEdit);
})();
