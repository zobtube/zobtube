(function() {
"use strict";
function esc(s) { return String(s).replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/"/g,"&quot;").replace(/>/g,"&gt;"); }
function escAttr(s) { return String(s).replace(/&/g,"&amp;").replace(/"/g,"&quot;").replace(/</g,"&lt;").replace(/>/g,"&gt;"); }

function ZtPhotosetEdit() {
  return Reflect.construct(HTMLElement, [], ZtPhotosetEdit);
}
ZtPhotosetEdit.prototype = Object.create(HTMLElement.prototype);
ZtPhotosetEdit.prototype.connectedCallback = function() {
  var self = this;
  var id = this.getAttribute("data-id");
  if (!id) {
    self.innerHTML = '<div class="alert alert-danger">Missing ID</div>';
    if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    return;
  }
  fetch("/api/photoset/" + encodeURIComponent(id) + "/edit", { credentials: "same-origin" })
    .then(function(r) { if (!r.ok) throw new Error(r.status); return r.json(); })
    .then(function(data) {
      var ps = data.photoset || data;
      var photos = ps.Photos || ps.photos || [];
      var actors = data.actors || data.Actors || [];
      var categories = data.categories || data.Categories || [];
      var name = esc(ps.Name || ps.name || "");
      var channel = ps.Channel || ps.channel;
      var channelName = channel ? (channel.Name || channel.name || "") : "None";
      var psActors = ps.Actors || ps.actors || [];
      var psCats = ps.Categories || ps.categories || [];

      var actorSelectedIds = psActors.map(function(a) { return a.ID || a.id; });
      var actorSelectable = {};
      actors.forEach(function(a) { actorSelectable[a.ID || a.id] = { name: a.Name || a.name || "" }; });
      var categorySelectedIds = psCats.map(function(c) { return c.ID || c.id; });
      var categorySelectable = {};
      categories.forEach(function(c) {
        (c.Sub || c.sub || []).forEach(function(s) {
          categorySelectable[s.ID || s.id] = { name: s.Name || s.name || "" };
        });
      });

      window.zt.actorSelection = window.zt.actorSelection || {};
      var actorSel = window.zt.actorSelection;
      actorSel.actorSelectable = actorSelectable;
      actorSel.actorSelected = {};
      actorSelectedIds.forEach(function(aid) { actorSel.actorSelected[aid] = undefined; });
      actorSel._updateSelectedActors = function() {
        self.querySelectorAll(".photoset-actor-list").forEach(function(chip) {
          var aid = chip.getAttribute("actor-id");
          if (!(aid in actorSel.actorSelected)) chip.remove();
          else chip.style.display = "";
        });
        self.querySelectorAll(".add-actor-list").forEach(function(chip) {
          var aid = chip.getAttribute("actor-id");
          var addBtn = chip.querySelector(".btn-success");
          var delBtn = chip.querySelector(".btn-danger");
          if (aid in actorSel.actorSelected) {
            if (addBtn) addBtn.style.display = "none";
            if (delBtn) delBtn.style.display = "";
            chip.style.display = "none";
          } else {
            if (addBtn) addBtn.style.display = "";
            if (delBtn) delBtn.style.display = "none";
            chip.style.display = "";
          }
        });
      };
      actorSel.actorSelect = function(aid) {
        actorSel.actorSelected[aid] = undefined;
        actorSel._updateSelectedActors();
      };
      actorSel.actorDeselect = function(aid) {
        var aname = (actorSelectable[aid] && actorSelectable[aid].name) || "Actor";
        fetch("/api/photoset/" + id + "/actor/" + aid, { method: "DELETE", credentials: "same-origin" })
          .then(function(r) {
            if (r.ok) {
              delete actorSel.actorSelected[aid];
              actorSel._updateSelectedActors();
              if (typeof sendToast === "function") sendToast("Actor removed", "", "bg-success", aname + " removed.");
            } else if (typeof sendToast === "function") sendToast("Actor not removed", "", "bg-danger", "Request failed.");
          });
      };

      window.zt.categorySelection = window.zt.categorySelection || {};
      var catSel = window.zt.categorySelection;
      catSel.categorySelectable = categorySelectable;
      catSel.categorySelected = {};
      categorySelectedIds.forEach(function(cid) { catSel.categorySelected[cid] = undefined; });
      catSel._updateSelectedCategories = function() {
        self.querySelectorAll(".photoset-category-list").forEach(function(chip) {
          var cid = chip.getAttribute("category-id");
          if (!(cid in catSel.categorySelected)) chip.remove();
          else chip.style.display = "";
        });
        self.querySelectorAll(".add-category-list").forEach(function(chip) {
          var cid = chip.getAttribute("category-id");
          var addBtn = chip.querySelector(".btn-success");
          var delBtn = chip.querySelector(".btn-danger");
          if (cid in catSel.categorySelected) {
            if (addBtn) addBtn.style.display = "none";
            if (delBtn) delBtn.style.display = "";
            chip.style.display = "none";
          } else {
            if (addBtn) addBtn.style.display = "";
            if (delBtn) delBtn.style.display = "none";
            chip.style.display = "";
          }
        });
      };
      catSel.categorySelect = function(cid) {
        catSel.categorySelected[cid] = undefined;
        catSel._updateSelectedCategories();
      };
      catSel.categoryDeselect = function(cid) {
        var cname = (categorySelectable[cid] && categorySelectable[cid].name) || "Category";
        fetch("/api/photoset/" + id + "/category/" + cid, { method: "DELETE", credentials: "same-origin" })
          .then(function(r) {
            if (r.ok) {
              delete catSel.categorySelected[cid];
              catSel._updateSelectedCategories();
              if (typeof sendToast === "function") sendToast("Category removed", "", "bg-success", cname + " removed.");
            } else if (typeof sendToast === "function") sendToast("Category not removed", "", "bg-danger", "Request failed.");
          });
      };

      var actorChips = actors.map(function(a) {
        var aid = a.ID || a.id;
        var show = actorSelectedIds.indexOf(aid) >= 0 ? "" : "display:none";
        return '<div class="chip photoset-actor-list" actor-id="' + aid + '" style="' + show + '"><img src="/api/actor/' + encodeURIComponent(aid) + '/thumb" width="50" height="50">' + esc(a.Name || a.name) +
          '<button type="button" class="btn btn-danger" onclick="window.zt.actorSelection.actorDeselect(\'' + aid + '\');"><i class="fa fa-trash-alt"></i></button></div>';
      }).join("");

      var categoryChips = [];
      categories.forEach(function(c) {
        (c.Sub || c.sub || []).forEach(function(s) {
          var sid = s.ID || s.id;
          var show = categorySelectedIds.indexOf(sid) >= 0 ? "" : "display:none";
          var thumb = '<img src="/api/category-sub/' + encodeURIComponent(sid) + '/thumb" width="50" height="50">';
          categoryChips.push('<div class="chip photoset-category-list" category-id="' + sid + '" style="' + show + '">' + thumb + esc(s.Name || s.name) +
            '<button type="button" class="btn btn-danger" onclick="window.zt.categorySelection.categoryDeselect(\'' + sid + '\');"><i class="fa fa-trash-alt"></i></button></div>');
        });
      });

      var actorModalChips = actors.map(function(a) {
        var aid = a.ID || a.id;
        var sel = actorSelectedIds.indexOf(aid) >= 0;
        return '<div class="chip add-actor-list" actor-id="' + aid + '" style="' + (sel ? "display:none;" : "") + '"><img src="/api/actor/' + encodeURIComponent(aid) + '/thumb" width="50" height="50">' + esc(a.Name || a.name) +
          '<button type="button" class="btn btn-success add-actor-add"><i class="fa fa-plus-circle"></i></button><button type="button" class="btn btn-danger add-actor-remove" style="' + (sel ? "" : "display:none") + '"><i class="fa fa-trash-alt"></i></button></div>';
      }).join("");

      var categoryModalHtml = "";
      categories.forEach(function(c) {
        var subs = c.Sub || c.sub || [];
        if (!subs.length) return;
        categoryModalHtml += '<h4 class="mt-3">' + esc(c.Name || c.name) + '</h4><div class="chips">';
        subs.forEach(function(s) {
          var sid = s.ID || s.id;
          var sel = categorySelectedIds.indexOf(sid) >= 0;
          var thumb = '<img src="/api/category-sub/' + encodeURIComponent(sid) + '/thumb" width="50" height="50">';
          categoryModalHtml += '<div class="chip add-category-list" category-id="' + sid + '" style="' + (sel ? "display:none;" : "") + '">' + thumb + esc(s.Name || s.name) +
            '<button type="button" class="btn btn-success add-category-add"><i class="fa fa-plus-circle"></i></button><button type="button" class="btn btn-danger add-category-remove" style="' + (sel ? "" : "display:none") + '"><i class="fa fa-trash-alt"></i></button></div>';
        });
        categoryModalHtml += "</div>";
      });

      var html = '<div class="themeix-section-h"><h3>Edit photoset</h3><hr /></div>';
      html += '<form id="zt-ps-rename" class="mb-3"><div class="input-group"><input class="form-control" name="name" value="' + name + '"><button class="btn btn-primary" type="submit">Rename</button></div></form>';
      html += '<p><a href="/photoset/' + id + '">Back to photoset</a></p>';

      html += '<div class="row mb-4"><div class="col-12">';
      html += '<div class="mb-3"><div class="form-floating input-group"><input type="text" disabled class="form-control" id="ps-channel" value="' + esc(channelName) + '"><label for="ps-channel">Channel</label><button class="btn btn-outline-warning" type="button" id="ps-channel-edit">Change</button></div></div>';
      html += '<div class="mb-3"><div class="form-floating"><div class="form-control chip-selector" style="height:unset;display:flex;"><div class="chips">' + actorChips +
        '<div class="chip">Add an actor<button type="button" class="btn btn-success" data-bs-toggle="modal" data-bs-target="#actorSelectionModal"><i class="fa fa-plus-circle"></i></button></div></div></div><label>Actors</label></div></div>';
      html += '<div class="mb-3"><div class="form-floating"><div class="form-control chip-selector" style="height:unset;display:flex;"><div class="chips">' + categoryChips.join("") +
        '<div class="chip">Add a category<button type="button" class="btn btn-success" data-bs-toggle="modal" data-bs-target="#categorySelectionModal"><i class="fa fa-plus-circle"></i></button></div></div></div><label>Categories</label></div></div>';
      html += "</div></div>";

      html += '<h5>Photos (' + photos.length + ')</h5><div class="row g-2">';
      photos.forEach(function(p) {
        var pid = p.ID || p.id;
        var thumb = "/api/photo/" + pid + "/thumb_mini";
        html += '<div class="col-4 col-md-2 text-center"><img class="img-fluid rounded lazy" data-src="' + thumb + '" style="aspect-ratio:1;object-fit:cover" alt="">';
        html += '<div class="mt-1"><button type="button" class="btn btn-sm btn-outline-primary zt-set-cover" data-id="' + pid + '">Set cover</button></div></div>';
      });
      html += "</div>";
      html += '<div class="mt-3"><button type="button" class="btn btn-danger" id="zt-ps-delete">Delete photoset</button></div>';

      html += '<div class="modal fade" id="actorSelectionModal" tabindex="-1"><div class="modal-dialog modal-xl"><div class="modal-content"><div class="modal-header"><h5 class="modal-title">Add actor to photoset</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div>';
      html += '<div class="modal-body"><div class="form-floating mb-3"><input type="text" class="form-control" id="actorSelectionModalInput" autocomplete="off"><label for="actorSelectionModalInput">Filter by name</label></div><div class="chips">' + actorModalChips +
        '<div class="chip">Add a new actor<button type="button" class="btn btn-success" onclick="window.open(\'/actor/new\',\'_blank\');"><i class="fa fa-plus-circle"></i></button></div></div></div>';
      html += '<div class="modal-footer"><button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button></div></div></div></div>';

      html += '<div class="modal fade" id="categorySelectionModal" tabindex="-1"><div class="modal-dialog modal-xl"><div class="modal-content"><div class="modal-header"><h5 class="modal-title">Add category to photoset</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div>';
      html += '<div class="modal-body">' + categoryModalHtml + '</div><div class="modal-footer"><button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button></div></div></div></div>';

      html += '<div class="modal fade" id="editChannelModal" tabindex="-1"><div class="modal-dialog modal-xl"><div class="modal-content"><div class="modal-header"><h5 class="modal-title">Change photoset channel</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div>';
      html += '<div class="modal-body"><div class="form-floating mb-3"><select class="form-select" id="channel-list"></select><label for="channel-list">Channel</label></div></div>';
      html += '<div class="modal-footer"><button type="button" class="btn btn-success" id="channel-send">Change</button><button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button></div></div></div></div>';

      self.innerHTML = html;

      self.querySelector("#zt-ps-rename").onsubmit = function(e) {
        e.preventDefault();
        var fd = new FormData(e.target);
        fetch("/api/photoset/" + id + "/rename", { method: "POST", credentials: "same-origin", body: fd })
          .then(function(r) { if (r.ok && typeof navigate === "function") navigate("/photoset/" + id + "/edit"); });
      };

      self.querySelector("#ps-channel-edit").addEventListener("click", function() {
        fetch("/api/channel/map", { credentials: "same-origin" }).then(function(r) { return r.json(); })
          .then(function(mapData) {
            var ch = mapData.channels || {};
            var sel = self.querySelector("#channel-list");
            var cur = channel ? (channel.ID || channel.id) : "";
            sel.innerHTML = '<option value="">None</option>' + Object.keys(ch).map(function(k) {
              return '<option value="' + escAttr(k) + '"' + (k === cur ? " selected" : "") + ">" + esc(ch[k] || "") + "</option>";
            }).join("");
            new bootstrap.Modal(self.querySelector("#editChannelModal")).show();
          });
      });

      self.querySelector("#channel-send").addEventListener("click", function() {
        var cid = self.querySelector("#channel-list").value;
        var fd = new FormData();
        fd.set("channel_id", cid);
        fetch("/api/photoset/" + id + "/channel", { method: "POST", credentials: "same-origin", body: fd })
          .then(function(r) {
            if (r.ok) {
              bootstrap.Modal.getInstance(self.querySelector("#editChannelModal")).hide();
              self.querySelector("#ps-channel").value = cid ? (self.querySelector("#channel-list option:checked").textContent || cid) : "None";
              if (typeof sendToast === "function") sendToast("Channel updated", "", "bg-success", "Channel saved.");
            } else if (typeof sendToast === "function") sendToast("Channel not updated", "", "bg-danger", "Request failed.");
          });
      });

      self.querySelectorAll(".add-actor-list .add-actor-add").forEach(function(btn) {
        btn.addEventListener("click", function() {
          var aid = this.closest(".add-actor-list").getAttribute("actor-id");
          var aname = (actorSelectable[aid] && actorSelectable[aid].name) || "Actor";
          fetch("/api/photoset/" + id + "/actor/" + aid, { method: "PUT", credentials: "same-origin" })
            .then(function(r) {
              if (r.ok) {
                window.zt.actorSelection.actorSelect(aid);
                if (typeof sendToast === "function") sendToast("Actor added", "", "bg-success", aname + " added.");
              } else if (typeof sendToast === "function") sendToast("Actor not added", "", "bg-danger", "Request failed.");
            });
        });
      });
      self.querySelectorAll(".add-actor-list .add-actor-remove").forEach(function(btn) {
        btn.addEventListener("click", function() {
          window.zt.actorSelection.actorDeselect(this.closest(".add-actor-list").getAttribute("actor-id"));
        });
      });

      self.querySelectorAll(".add-category-list .add-category-add").forEach(function(btn) {
        btn.addEventListener("click", function() {
          var cid = this.closest(".add-category-list").getAttribute("category-id");
          var cname = (categorySelectable[cid] && categorySelectable[cid].name) || "Category";
          fetch("/api/photoset/" + id + "/category/" + cid, { method: "PUT", credentials: "same-origin" })
            .then(function(r) {
              if (r.ok) {
                window.zt.categorySelection.categorySelect(cid);
                if (typeof sendToast === "function") sendToast("Category added", "", "bg-success", cname + " added.");
              } else if (typeof sendToast === "function") sendToast("Category not added", "", "bg-danger", "Request failed.");
            });
        });
      });
      self.querySelectorAll(".add-category-list .add-category-remove").forEach(function(btn) {
        btn.addEventListener("click", function() {
          window.zt.categorySelection.categoryDeselect(this.closest(".add-category-list").getAttribute("category-id"));
        });
      });

      var filterInput = self.querySelector("#actorSelectionModalInput");
      if (filterInput) {
        filterInput.addEventListener("input", function() {
          var q = filterInput.value.trim().toLowerCase();
          self.querySelectorAll("#actorSelectionModal .add-actor-list").forEach(function(chip) {
            var aid = chip.getAttribute("actor-id");
            var aname = ((actorSelectable[aid] && actorSelectable[aid].name) || "").toLowerCase();
            chip.style.display = !q || aname.indexOf(q) >= 0 ? "" : "none";
          });
        });
      }

      self.querySelectorAll(".zt-set-cover").forEach(function(btn) {
        btn.onclick = function() {
          var pid = btn.getAttribute("data-id");
          fetch("/api/photoset/" + id + "/cover/" + pid, { method: "POST", credentials: "same-origin" })
            .then(function(r) { if (r.ok && typeof sendToast === "function") sendToast("Cover set", "", "bg-success", "Cover updated."); });
        };
      });

      var del = self.querySelector("#zt-ps-delete");
      if (del) del.onclick = function() {
        if (!confirm("Delete this photoset?")) return;
        fetch("/api/photoset/" + id, { method: "DELETE", credentials: "same-origin" })
          .then(function(r) { if (r.ok && typeof navigate === "function") navigate("/adm/tasks"); });
      };

      if (window.zt && window.zt.loadLazyIn) window.zt.loadLazyIn(self);
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() {
      self.innerHTML = '<div class="alert alert-danger">Failed to load.</div>';
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    });
};
customElements.define("zt-photoset-edit", ZtPhotosetEdit);
})();
