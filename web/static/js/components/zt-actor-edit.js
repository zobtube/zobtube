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
    fetch("/api/category", { credentials: "same-origin" }).then(function(r) { if (!r.ok) throw new Error(r.status); return r.json(); }),
    fetch("/api/adm/config/provider", { credentials: "same-origin" }).then(function(r) { if (!r.ok) throw new Error(r.status); return r.json(); })
  ]).then(function(results) {
    var a = results[0];
    var catsData = results[1];
    var providerData = results[2];
    var categories = catsData.items || catsData.categories || [];
    var urlView = "/actor/" + id;
    var urlThumb = "/api/actor/" + id + "/thumb";
    var name = esc(a.Name || a.name || "");
    var sex = a.Sex || a.sex || "m";
    var desc = esc(a.Description || a.description || "");
    var aliases = a.Aliases || a.aliases || [];
    var actorCats = a.Categories || a.categories || [];
    var catSelectedIds = actorCats.map(function(c) { return c.ID || c.id; });
    var linksMap = {};
    (a.Links || a.links || []).forEach(function(l) {
      var pid = l.Provider || l.provider;
      if (pid) linksMap[pid] = { link_url: l.URL || l.url || "", link_id: l.ID || l.id || "" };
    });
    var providers = providerData.providers || [];
    var offlineMode = !!providerData.offline_mode;
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
    html += '<a class="btn btn-danger" style="margin-top:15px;width:100%" href="/actor/' + id + '/delete" data-full-reload>Delete actor profile</a>';
    html += '<button type="button" class="btn btn-info" style="margin-top:15px;width:100%" id="zt-view-link-pictures">View link pictures</button></div></div>';
    html += '<div class="col-9"><h3>Profile details</h3><br/>';
    html += '<div class="mb-3"><div class="form-floating input-group"><input type="text" class="form-control" id="actor-name" disabled value="' + name + '"><label for="actor-name">Name</label><button class="btn btn-outline-warning" type="button" id="actor-name-edit">Edit</button></div></div>';
    html += '<div class="mb-3"><label class="form-label">Sex:</label><select class="form-select" id="actor-sex" disabled><option value="m"' + (sex==="m"?" selected":"") + '>Male</option><option value="f"' + (sex==="f"?" selected":"") + '>Female</option><option value="tw"' + (sex==="tw"?" selected":"") + '>Trans women</option></select></div>';
    html += '<div class="mb-3"><label class="form-label" for="actor-description">Description&nbsp;<i id="actor-description-btn-edit" class="far fa-edit text-warning clickable" title="Edit"></i><i id="actor-description-btn-discard" class="fas fa-history text-danger clickable" title="Discard" style="display:none;"></i><i id="actor-description-btn-save" class="far fa-save text-success clickable" title="Save" style="display:none;"></i></label><textarea class="form-control" id="actor-description" readonly disabled>' + desc + '</textarea></div>';
    html += '<h4 class="mt-3 mb-3">Aliases</h4><div class="form-control chip-selector" style="height:unset;display:flex"><div class="chips">' + aliasChips;
    html += '<div class="chip">Add an alias<button class="btn btn-success" data-bs-toggle="modal" data-bs-target="#addActorAliasModal"><i class="fa fa-plus-circle"></i></button></div></div></div>';
    html += '<h4 class="mt-3">Categories</h4><div class="form-control chip-selector" style="height:unset;display:flex"><div class="chips">' + categoryChips.join("") + '<div class="chip">Add a category<button class="btn btn-success" data-bs-toggle="modal" data-bs-target="#addCategoryModal"><i class="fa fa-plus-circle"></i></button></div></div></div>';
    html += '<h4 class="mt-4">Actor links</h4><table class="table"><thead><tr><th>Provider</th><th style="text-align:right"></th></tr></thead><tbody>';
    providers.forEach(function(p) {
      var pid = p.ID || p.id;
      var niceName = esc(p.NiceName || p.nice_name || pid);
      var enabled = p.Enabled !== false && p.enabled !== false;
      var disabledCls = enabled ? "" : " disabled";
      html += '<tr id="provider-' + esc(pid) + '"><td style="vertical-align:middle">' + niceName + (enabled ? "" : ' <span class="badge text-bg-warning">Disabled</span>') + '</td><td style="text-align:right;vertical-align:middle">';
      html += '<button id="provider-action-' + esc(pid) + '-search" class="btn btn-primary' + disabledCls + '"><i class="fas fa-search"></i> Automatic search</button> ';
      html += '<button id="provider-action-' + esc(pid) + '-add" class="btn btn-primary"><i class="fas fa-plus"></i> Manual add</button> ';
      html += '<a class="btn btn-success" href="" id="provider-action-' + esc(pid) + '-view" target="_blank" rel="noopener noreferrer"><i class="fas fa-globe"></i> View page</a> ';
      html += '<button class="btn btn-danger" id="provider-action-' + esc(pid) + '-delete"><i class="fas fa-trash"></i> Remove link</button> ';
      html += '<span id="provider-text-' + esc(pid) + '-first-time" style="display:none"><i>Automatic research will begin shortly...</i></span>';
      html += '</td></tr>';
    });
    html += '</tbody></table></div></div></div>';

    html += '<div class="row" id="profile-picture-propositions-row" style="display:none"><div class="col-12"><h4 class="mt-4">Actor pictures from external links</h4><hr/><div class="row" id="profile-picture-propositions"><div class="col-12"><div class="alert alert-warning" id="profile-picture-suggestion-welcome">Picture suggestion loading...</div><div class="alert alert-info" id="profile-picture-suggestion-click-info" style="display:none"><i>Click on one picture to set it as profile picture</i></div></div><div class="col-3" id="profile-picture-template" style="display:none"><div class="card"><div class="card-body"><img class="rounded-start" style="width:100%;cursor:pointer" data-bs-toggle="modal" data-bs-target="#profilePictureModal"></div><div class="card-footer" style="text-align:center">Provider</div></div></div></div></div></div>';

    html += '<div class="modal fade" id="profilePictureModal" tabindex="-1"><div class="modal-dialog modal-xl"><div class="modal-content"><div class="modal-header"><h5 class="modal-title">Profile picture preview</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div><div class="modal-body" id="profilePictureModalBody"><h5 id="cropper-status"></h5><div id="profile-picture-drop-area" style="width:100%;height:400px;border:1px solid #000;background-position:center;background-size:cover;box-sizing:border-box"><img id="profile-picture-crop-img" src="" style="height:100%;width:auto;display:none"></div></div><div class="modal-footer"><button type="button" class="btn btn-primary" id="profile-picture-set-btn">Set as new profile picture</button></div></div></div></div>';

    html += '<div class="modal fade" id="manualLinkModal" tabindex="-1"><div class="modal-dialog"><div class="modal-content"><div class="modal-header"><h5 class="modal-title">Add a new link for ' + name + '</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div><div class="modal-body"><div class="form-floating mb-3"><input type="url" class="form-control" id="newLinkInput" placeholder=" "><label for="newLinkInput">URL to the profile</label></div></div><div class="modal-footer"><button type="button" class="btn btn-primary" id="manualLinkSubmit">Add link</button><button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button></div></div></div></div>';
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
    var descBtnEdit = self.querySelector("#actor-description-btn-edit");
    var descBtnSave = self.querySelector("#actor-description-btn-save");
    var descBtnDiscard = self.querySelector("#actor-description-btn-discard");
    descBtnEdit.addEventListener("click", function() {
      descEl.old_value = descEl.value;
      descEl.disabled = false; descEl.removeAttribute("readonly");
      descBtnEdit.style.display = "none"; descBtnSave.style.display = ""; descBtnDiscard.style.display = "";
    });
    descBtnDiscard.addEventListener("click", function() {
      descEl.value = descEl.old_value || "";
      descEl.disabled = true; descEl.setAttribute("readonly", "");
      descBtnEdit.style.display = ""; descBtnSave.style.display = "none"; descBtnDiscard.style.display = "none";
    });
    descBtnSave.addEventListener("click", function() {
      var fd = new FormData(); fd.set("description", descEl.value);
      fetch("/api/actor/" + id + "/description", { method: "POST", credentials: "same-origin", body: fd })
        .then(function(r) {
          if (r.ok) {
            descEl.disabled = true; descEl.setAttribute("readonly", "");
            descBtnEdit.style.display = ""; descBtnSave.style.display = "none"; descBtnDiscard.style.display = "none";
            if (typeof sendToast === "function") sendToast("Actor description updated", "", "bg-success", "Description changed for " + (nameEl.value || ""));
          } else {
            r.json().catch(function(){return{};}).then(function(d){ if (typeof sendToast === "function") sendToast("Actor description update failed", "", "bg-danger", (d && d.error) || "Unexpected error during actor description update"); });
          }
        });
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

    function providerLinkPresent(providerSlug, linkUrl, linkId) {
      var row = self.querySelector("#provider-" + providerSlug);
      if (!row) return;
      row.querySelector("#provider-action-" + providerSlug + "-search").style.display = "none";
      row.querySelector("#provider-action-" + providerSlug + "-add").style.display = "none";
      var firstTime = row.querySelector("#provider-text-" + providerSlug + "-first-time");
      if (firstTime) firstTime.style.display = "none";
      var btnView = row.querySelector("#provider-action-" + providerSlug + "-view");
      var btnDelete = row.querySelector("#provider-action-" + providerSlug + "-delete");
      if (btnView) { btnView.style.display = ""; btnView.href = linkUrl; }
      if (btnDelete) btnDelete.style.display = "";
    }
    function providerLinkAbsent(providerSlug) {
      var row = self.querySelector("#provider-" + providerSlug);
      if (!row) return;
      row.querySelector("#provider-action-" + providerSlug + "-view").style.display = "none";
      row.querySelector("#provider-action-" + providerSlug + "-delete").style.display = "none";
      var firstTime = row.querySelector("#provider-text-" + providerSlug + "-first-time");
      if (firstTime) firstTime.style.display = "none";
      row.querySelector("#provider-action-" + providerSlug + "-search").style.display = "";
      row.querySelector("#provider-action-" + providerSlug + "-add").style.display = "";
    }
    function providerFirstTime(providerSlug) {
      var row = self.querySelector("#provider-" + providerSlug);
      if (!row) return;
      row.querySelector("#provider-action-" + providerSlug + "-view").style.display = "none";
      row.querySelector("#provider-action-" + providerSlug + "-delete").style.display = "none";
      row.querySelector("#provider-action-" + providerSlug + "-search").style.display = "none";
      row.querySelector("#provider-action-" + providerSlug + "-add").style.display = "none";
      var firstTime = row.querySelector("#provider-text-" + providerSlug + "-first-time");
      if (firstTime) firstTime.style.display = "";
    }
    function updateProviders() {
      providers.forEach(function(p) {
        var pid = p.ID || p.id;
        if (linksMap[pid]) {
          providerLinkPresent(pid, linksMap[pid].link_url, linksMap[pid].link_id);
        } else {
          providerLinkAbsent(pid);
        }
      });
    }
    updateProviders();

    providers.forEach(function(p) {
      var pid = p.ID || p.id;
      var niceName = p.NiceName || p.nice_name || pid;
      var enabled = p.Enabled !== false && p.enabled !== false;
      var searchBtn = self.querySelector("#provider-action-" + pid + "-search");
      var addBtn = self.querySelector("#provider-action-" + pid + "-add");
      var deleteBtn = self.querySelector("#provider-action-" + pid + "-delete");
      if (searchBtn) {
        searchBtn.addEventListener("click", function() {
          if (!enabled) return;
          if (typeof sendToast === "function") sendToast("Automatic actor search", "", "bg-info", "Searching for " + name + " on " + niceName);
          providerFirstTime(pid);
          fetch("/api/actor/" + id + "/provider/" + encodeURIComponent(pid), { credentials: "same-origin" })
            .then(function(r) { return r.json(); })
            .then(function(data) {
              if (data.link_url) {
                linksMap[pid] = { link_url: data.link_url, link_id: data.link_id || "" };
                updateProviders();
                if (typeof sendToast === "function") sendToast("Automatic actor search", "", "bg-success", name + " found on " + niceName + "!");
              } else {
                updateProviders();
                if (typeof sendToast === "function") sendToast("Automatic actor search", "", "bg-warning", name + " not found on " + niceName);
              }
            })
            .catch(function() {
              updateProviders();
              if (typeof sendToast === "function") sendToast("Automatic actor search", "", "bg-warning", name + " not found on " + niceName);
            });
        });
      }
      if (addBtn) {
        addBtn.addEventListener("click", function() {
          self._manualLinkActorId = id;
          self._manualLinkProviderSlug = pid;
          self.querySelector("#newLinkInput").value = "";
          bootstrap.Modal.getOrCreateInstance(self.querySelector("#manualLinkModal")).show();
        });
      }
      if (deleteBtn) {
        deleteBtn.addEventListener("click", function() {
          var linkId = linksMap[pid] && linksMap[pid].link_id;
          if (!linkId) return;
          fetch("/api/actor/link/" + linkId, { method: "DELETE", credentials: "same-origin" })
            .then(function(r) {
              if (r.ok) {
                delete linksMap[pid];
                updateProviders();
                if (typeof sendToast === "function") sendToast("Actor's link deleted", "", "bg-success", "Success");
              } else {
                if (typeof sendToast === "function") sendToast("Actor's link deletion", "", "bg-warning", "Unable to delete");
              }
            })
            .catch(function() {
              if (typeof sendToast === "function") sendToast("Actor's link deletion", "", "bg-warning", "Unable to delete");
            });
        });
      }
    });

    function handlePictureError(ev) {
      ev.target.src = offlineMode ? "/static/images/actor-picture-error-offline-mode.svg" : "/static/images/actor-picture-error-not-found.svg";
    }
    function showActorPictures() {
      var propRow = self.querySelector("#profile-picture-propositions-row");
      var propRoot = self.querySelector("#profile-picture-propositions");
      var template = self.querySelector("#profile-picture-template");
      var welcome = self.querySelector("#profile-picture-suggestion-welcome");
      var clickInfo = self.querySelector("#profile-picture-suggestion-click-info");
      if (!propRow || !propRoot || !template) return;
      propRow.style.display = "";
      propRoot.querySelectorAll("[id^=profile-picture-provider-]").forEach(function(el) { el.remove(); });
      var hasAny = false;
      for (var pid in linksMap) {
        var prov = providers.find(function(p) { return (p.ID || p.id) === pid; });
        if (prov && (prov.Enabled === false || prov.enabled === false)) continue;
        var linkId = linksMap[pid].link_id;
        if (!linkId) continue;
        var niceName = prov ? (prov.NiceName || prov.nice_name || pid) : pid;
        var thumbUrl = "/api/actor/link/" + linkId + "/thumb";
        if (welcome) welcome.style.display = "none";
        if (clickInfo) clickInfo.style.display = "";
        hasAny = true;
        var clone = template.cloneNode(true);
        clone.style.display = "";
        clone.id = "profile-picture-provider-" + pid;
        var img = clone.querySelector("img");
        img.addEventListener("error", handlePictureError);
        img.src = thumbUrl;
        clone.querySelector(".card-footer").textContent = niceName;
        propRoot.appendChild(clone);
      }
      if (!hasAny && welcome) welcome.style.display = "";
    }
    self.querySelector("#zt-view-link-pictures").addEventListener("click", showActorPictures);

    var profilePictureModal = self.querySelector("#profilePictureModal");
    var cropImg = self.querySelector("#profile-picture-crop-img");
    var profilePictureCropper = null;
    if (profilePictureModal) {
      profilePictureModal.addEventListener("shown.bs.modal", function(ev) {
        var trigger = ev.relatedTarget;
        if (trigger && trigger.src) {
          self.querySelector("#cropper-status").textContent = "Cropper library loading...";
          cropImg.style.display = "none";
          cropImg.src = trigger.src;
        }
      });
      cropImg.onload = function() {
        if (profilePictureCropper) { profilePictureCropper.destroy(); profilePictureCropper = null; }
        cropImg.style.display = "";
        var minAr = 1, maxAr = 1;
        profilePictureCropper = new Cropper(cropImg, {
          viewMode: 2,
          aspectRatio: 1,
          ready: function() {
            self.querySelector("#cropper-status").textContent = "Cropper library ready";
          },
          cropmove: function() {
            var cb = this.cropper.getCropBoxData();
            var ar = cb.width / cb.height;
            if (ar < minAr) this.cropper.setCropBoxData({ width: cb.height * minAr });
            else if (ar > maxAr) this.cropper.setCropBoxData({ width: cb.height * maxAr });
          }
        });
      };
      profilePictureModal.addEventListener("hide.bs.modal", function() {
        if (profilePictureCropper) { profilePictureCropper.destroy(); profilePictureCropper = null; }
      });
    }
    self.querySelector("#profile-picture-set-btn").addEventListener("click", function() {
      if (!profilePictureCropper) return;
      profilePictureCropper.getCroppedCanvas().toBlob(function(blob) {
        var fd = new FormData();
        fd.set("pp", blob);
        fetch("/api/actor/" + id + "/thumb", { method: "POST", credentials: "same-origin", body: fd })
          .then(function(r) {
            if (r.ok) { bootstrap.Modal.getInstance(profilePictureModal).hide(); refresh(); if (typeof sendToast === "function") sendToast("Profile picture updated", "", "bg-success", "Success"); }
          });
      });
    });

    self.querySelector("#manualLinkSubmit").addEventListener("click", function() {
      var url = self.querySelector("#newLinkInput").value.trim();
      var actorId = self._manualLinkActorId || id;
      var providerSlug = self._manualLinkProviderSlug || "";
      if (!url || !providerSlug) return;
      var fd = new FormData();
      fd.set("url", url);
      fd.set("provider", providerSlug);
      fetch("/api/actor/" + actorId + "/link", { method: "POST", credentials: "same-origin", body: fd })
        .then(function(r) { return r.json(); })
        .then(function(data) {
          if (data.link_url) {
            linksMap[providerSlug] = { link_url: data.link_url, link_id: data.link_id || "" };
            updateProviders();
            bootstrap.Modal.getInstance(self.querySelector("#manualLinkModal")).hide();
            if (typeof sendToast === "function") sendToast("New link added for " + name, "", "bg-success", "Success");
          } else {
            if (typeof sendToast === "function") sendToast("Automatic actor search", "", "bg-warning", "Adding failed");
          }
        })
        .catch(function() {
          if (typeof sendToast === "function") sendToast("Automatic actor search", "", "bg-warning", "Adding failed");
        });
    });

    if (!offlineMode && (a.Links || a.links || []).length === 0) {
      var firstEnabled = providers.find(function(p) { return (p.Enabled !== false && p.enabled !== false) && !linksMap[p.ID || p.id]; });
      if (firstEnabled) {
        var slug = firstEnabled.ID || firstEnabled.id;
        providerFirstTime(slug);
        fetch("/api/actor/" + id + "/provider/" + encodeURIComponent(slug), { credentials: "same-origin" })
          .then(function(r) { return r.json(); })
          .then(function(data) {
            if (data.link_url) {
              linksMap[slug] = { link_url: data.link_url, link_id: data.link_id || "" };
              updateProviders();
            } else { updateProviders(); }
          })
          .catch(function() { updateProviders(); });
      }
    }

    if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
  }).catch(function(e) {
    console.error('error loading actor edit', e);
    if (e && e.message === "403") self.innerHTML = '<div class="alert alert-danger">Forbidden</div>';
    else if (e && e.message === "404") self.innerHTML = '<div class="alert alert-danger">Not found</div>';
    else self.innerHTML = '<div class="alert alert-danger">Failed to load.</div>';
    if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
  });
};
customElements.define("zt-actor-edit", ZtActorEdit);
})();
