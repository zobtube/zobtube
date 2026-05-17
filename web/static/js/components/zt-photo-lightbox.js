(function() {
"use strict";
function esc(s) { return String(s).replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/"/g,"&quot;").replace(/>/g,"&gt;"); }
function escAttr(s) { return String(s).replace(/&/g,"&amp;").replace(/"/g,"&quot;").replace(/</g,"&lt;").replace(/>/g,"&gt;"); }

function photoId(p) { return p.ID || p.id; }

function ZtPhotoLightbox() {
  return Reflect.construct(HTMLElement, [], ZtPhotoLightbox);
}
ZtPhotoLightbox.prototype = Object.create(HTMLElement.prototype);

ZtPhotoLightbox.prototype.connectedCallback = function() {
  var self = this;
  this._photos = [];
  try { this._photos = JSON.parse(this.getAttribute("data-photos") || "[]"); } catch (e) {}
  this._index = parseInt(this.getAttribute("data-index") || "0", 10) || 0;
  this._photosetId = this.getAttribute("data-photoset-id") || "";
  this._admin = !!(window.__USER__ && window.__USER__.admin);
  this._catalog = null;
  this._editorOpen = false;
  this._actorSel = null;
  this._catSel = null;
  this.render();
  document.addEventListener("keydown", this._onKey = function(e) {
    if (!self.parentElement) return;
    if (e.key === "Escape") {
      var openModal = self.querySelector(".modal.show");
      if (openModal) return;
      self.close();
    }
    if (e.key === "ArrowLeft") self.prev();
    if (e.key === "ArrowRight") self.next();
  });
};

ZtPhotoLightbox.prototype.disconnectedCallback = function() {
  if (this._onKey) document.removeEventListener("keydown", this._onKey);
};

ZtPhotoLightbox.prototype.close = function() {
  this.remove();
};

ZtPhotoLightbox.prototype.prev = function() {
  if (this._photos.length < 2) return;
  this._index = (this._index - 1 + this._photos.length) % this._photos.length;
  this.render();
};

ZtPhotoLightbox.prototype.next = function() {
  if (this._photos.length < 2) return;
  this._index = (this._index + 1) % this._photos.length;
  this.render();
};

ZtPhotoLightbox.prototype.current = function() {
  return this._photos[this._index] || {};
};

ZtPhotoLightbox.prototype.refreshPhotos = function() {
  var self = this;
  var keepEditor = self._editorOpen;
  if (!self._photosetId) return Promise.resolve();
  return fetch("/api/photoset/" + encodeURIComponent(self._photosetId), { credentials: "same-origin" })
    .then(function(r) { if (!r.ok) throw new Error(r.status); return r.json(); })
    .then(function(data) {
      self._photos = data.photos || [];
      if (self._index >= self._photos.length) self._index = 0;
      self._editorOpen = keepEditor;
      self.render();
    });
};

ZtPhotoLightbox.prototype.loadCatalog = function() {
  var self = this;
  if (self._catalog) return Promise.resolve(self._catalog);
  if (!self._photosetId) return Promise.reject(new Error("no photoset"));
  return fetch("/api/photoset/" + encodeURIComponent(self._photosetId) + "/edit", { credentials: "same-origin" })
    .then(function(r) { if (!r.ok) throw new Error(r.status); return r.json(); })
    .then(function(data) {
      var ps = data.photoset || data;
      self._catalog = {
        actors: data.actors || data.Actors || [],
        categories: data.categories || data.Categories || [],
        photosetActors: ps.Actors || ps.actors || [],
        photosetCategories: ps.Categories || ps.categories || [],
        photosetChannel: ps.Channel || ps.channel || null
      };
      return self._catalog;
    });
};

ZtPhotoLightbox.prototype.photoActors = function(p) {
  return p.Actors || p.actors || [];
};

ZtPhotoLightbox.prototype.photoCategories = function(p) {
  return p.Categories || p.categories || [];
};

ZtPhotoLightbox.prototype.photoChannel = function(p) {
  return p.Channel || p.channel || null;
};

ZtPhotoLightbox.prototype.buildEditorHtml = function(p, catalog) {
  var self = this;
  var pid = photoId(p);
  var actors = catalog.actors;
  var categories = catalog.categories;
  var photoActors = self.photoActors(p);
  var photoCats = self.photoCategories(p);
  var photoCh = self.photoChannel(p);
  var actorSelectedIds = photoActors.map(photoId);
  var catSelectedIds = photoCats.map(photoId);
  var channelName = photoCh ? (photoCh.Name || photoCh.name || "") : "None (inherit photoset)";

  var actorSelectable = {};
  actors.forEach(function(a) { actorSelectable[photoId(a)] = { name: a.Name || a.name || "" }; });

  var categorySelectable = {};
  categories.forEach(function(c) {
    (c.Sub || c.sub || []).forEach(function(s) {
      categorySelectable[photoId(s)] = { name: s.Name || s.name || "" };
    });
  });

  var actorSel = self._actorSel = self._actorSel || {};
  actorSel.actorSelectable = actorSelectable;
  actorSel.actorSelected = {};
  actorSelectedIds.forEach(function(aid) { actorSel.actorSelected[aid] = undefined; });
  actorSel._root = self;
  actorSel._photoId = pid;
  actorSel._updateSelectedActors = function() {
    self.querySelectorAll(".photo-actor-list").forEach(function(chip) {
      var aid = chip.getAttribute("actor-id");
      if (!(aid in actorSel.actorSelected)) chip.remove();
      else chip.style.display = "";
    });
    self.querySelectorAll("#zt-lb-actorSelectionModal .add-actor-list").forEach(function(chip) {
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
    fetch("/api/photoset/photo/" + pid + "/actor/" + aid, { method: "DELETE", credentials: "same-origin" })
      .then(function(r) {
        if (r.ok) {
          return self.refreshPhotos();
        }
        if (typeof sendToast === "function") sendToast("Actor not removed", "", "bg-danger", "Request failed.");
      });
  };
  window.zt = window.zt || {};
  window.zt.lbActorSelection = actorSel;

  var catSel = self._catSel = self._catSel || {};
  catSel.categorySelectable = categorySelectable;
  catSel.categorySelected = {};
  catSelectedIds.forEach(function(cid) { catSel.categorySelected[cid] = undefined; });
  catSel._root = self;
  catSel._photoId = pid;
  catSel._updateSelectedCategories = function() {
    self.querySelectorAll(".photo-category-list").forEach(function(chip) {
      var cid = chip.getAttribute("category-id");
      if (!(cid in catSel.categorySelected)) chip.remove();
      else chip.style.display = "";
    });
    self.querySelectorAll("#zt-lb-categorySelectionModal .add-category-list").forEach(function(chip) {
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
    fetch("/api/photoset/photo/" + pid + "/category/" + cid, { method: "DELETE", credentials: "same-origin" })
      .then(function(r) {
        if (r.ok) return self.refreshPhotos();
        if (typeof sendToast === "function") sendToast("Category not removed", "", "bg-danger", "Request failed.");
      });
  };
  window.zt.lbCategorySelection = catSel;

  var actorChips = actors.map(function(a) {
    var aid = photoId(a);
    var show = actorSelectedIds.indexOf(aid) >= 0 ? "" : "display:none";
    return '<div class="chip photo-actor-list" actor-id="' + aid + '" style="' + show + '"><img src="/api/actor/' + encodeURIComponent(aid) + '/thumb" width="50" height="50">' + esc(a.Name || a.name) +
      '<button type="button" class="btn btn-danger" onclick="window.zt.lbActorSelection.actorDeselect(\'' + aid + '\');"><i class="fa fa-trash-alt"></i></button></div>';
  }).join("");

  var categoryChips = [];
  categories.forEach(function(c) {
    (c.Sub || c.sub || []).forEach(function(s) {
      var sid = photoId(s);
      var show = catSelectedIds.indexOf(sid) >= 0 ? "" : "display:none";
      var thumb = '<img src="/api/category-sub/' + encodeURIComponent(sid) + '/thumb" width="50" height="50">';
      categoryChips.push('<div class="chip photo-category-list" category-id="' + sid + '" style="' + show + '">' + thumb + esc(s.Name || s.name) +
        '<button type="button" class="btn btn-danger" onclick="window.zt.lbCategorySelection.categoryDeselect(\'' + sid + '\');"><i class="fa fa-trash-alt"></i></button></div>');
    });
  });

  var actorModalChips = actors.map(function(a) {
    var aid = photoId(a);
    var sel = actorSelectedIds.indexOf(aid) >= 0;
    return '<div class="chip add-actor-list" actor-id="' + aid + '" style="' + (sel ? "display:none;" : "") + '"><img src="/api/actor/' + encodeURIComponent(aid) + '/thumb" width="50" height="50">' + esc(a.Name || a.name) +
      '<button type="button" class="btn btn-success zt-lb-add-actor"><i class="fa fa-plus-circle"></i></button><button type="button" class="btn btn-danger zt-lb-remove-actor" style="' + (sel ? "" : "display:none") + '"><i class="fa fa-trash-alt"></i></button></div>';
  }).join("");

  var categoryModalHtml = "";
  categories.forEach(function(c) {
    var subs = c.Sub || c.sub || [];
    if (!subs.length) return;
    categoryModalHtml += '<h4 class="mt-3">' + esc(c.Name || c.name) + '</h4><div class="chips zt-lb-chips">';
    subs.forEach(function(s) {
      var sid = photoId(s);
      var sel = catSelectedIds.indexOf(sid) >= 0;
      var thumb = '<img src="/api/category-sub/' + encodeURIComponent(sid) + '/thumb" width="50" height="50">';
      categoryModalHtml += '<div class="chip add-category-list" category-id="' + sid + '" style="' + (sel ? "display:none;" : "") + '">' + thumb + esc(s.Name || s.name) +
        '<button type="button" class="btn btn-success zt-lb-add-category"><i class="fa fa-plus-circle"></i></button><button type="button" class="btn btn-danger zt-lb-remove-category" style="' + (sel ? "" : "display:none") + '"><i class="fa fa-trash-alt"></i></button></div>';
    });
    categoryModalHtml += "</div>";
  });

  var inherited = "";
  var psActors = catalog.photosetActors || [];
  var psCats = catalog.photosetCategories || [];
  if (psActors.length || psCats.length || catalog.photosetChannel) {
    inherited = '<p class="small text-muted mb-2">From photoset: ';
    var parts = [];
    if (catalog.photosetChannel) {
      parts.push("channel " + esc(catalog.photosetChannel.Name || catalog.photosetChannel.name || ""));
    }
    psActors.forEach(function(a) { parts.push("@" + esc(a.Name || a.name || "")); });
    psCats.forEach(function(c) { parts.push("#" + esc(c.Name || c.name || "")); });
    inherited += parts.join(", ") + "</p>";
  }

  var html = '<p class="small text-muted mb-2">Overrides on this photo only (layered on photoset defaults).</p>';
  html += inherited;
  html += '<div class="mb-2"><div class="input-group input-group-sm"><span class="input-group-text">Channel</span>';
  html += '<input type="text" class="form-control" id="zt-lb-channel-display" readonly value="' + esc(channelName) + '">';
  html += '<button type="button" class="btn btn-outline-warning" id="zt-lb-channel-edit">Change</button></div></div>';
  html += '<div class="mb-2"><label class="form-label small text-light">Actors on this photo</label>';
  html += '<div class="chip-selector zt-lb-chips" style="display:flex;flex-wrap:wrap;gap:0.25rem"><div class="chips zt-lb-chips">' + actorChips;
  html += '<div class="chip">Add<button type="button" class="btn btn-success btn-sm ms-1" id="zt-lb-open-actor-modal"><i class="fa fa-plus-circle"></i></button></div></div></div></div>';
  html += '<div class="mb-2"><label class="form-label small text-light">Categories on this photo</label>';
  html += '<div class="chip-selector zt-lb-chips" style="display:flex;flex-wrap:wrap;gap:0.25rem"><div class="chips zt-lb-chips">' + categoryChips.join("");
  html += '<div class="chip">Add<button type="button" class="btn btn-success btn-sm ms-1" id="zt-lb-open-category-modal"><i class="fa fa-plus-circle"></i></button></div></div></div>';

  html += '<div class="modal fade zt-lb-modal" id="zt-lb-actorSelectionModal" tabindex="-1"><div class="modal-dialog modal-xl"><div class="modal-content">';
  html += '<div class="modal-header"><h5 class="modal-title">Add actor to photo</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div>';
  html += '<div class="modal-body"><div class="mb-3"><input type="text" class="form-control form-control-sm" id="zt-lb-actor-filter" placeholder="Filter by name" autocomplete="off"></div>';
  html += '<div class="chips zt-lb-chips">' + actorModalChips + '</div></div><div class="modal-footer"><button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button></div></div></div></div>';

  html += '<div class="modal fade zt-lb-modal" id="zt-lb-categorySelectionModal" tabindex="-1"><div class="modal-dialog modal-xl"><div class="modal-content">';
  html += '<div class="modal-header"><h5 class="modal-title">Add category to photo</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div>';
  html += '<div class="modal-body">' + categoryModalHtml + '</div><div class="modal-footer"><button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button></div></div></div></div>';

  html += '<div class="modal fade zt-lb-modal" id="zt-lb-editChannelModal" tabindex="-1"><div class="modal-dialog"><div class="modal-content">';
  html += '<div class="modal-header"><h5 class="modal-title">Photo channel override</h5><button type="button" class="btn-close" data-bs-dismiss="modal"></button></div>';
  html += '<div class="modal-body"><select class="form-select form-select-sm" id="zt-lb-channel-list"></select>';
  html += '<p class="small text-muted mt-2 mb-0">Choose None to inherit the photoset channel.</p></div>';
  html += '<div class="modal-footer"><button type="button" class="btn btn-success" id="zt-lb-channel-send">Save</button>';
  html += '<button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button></div></div></div></div>';

  return html;
};

ZtPhotoLightbox.prototype.wireEditor = function() {
  var self = this;
  var p = self.current();
  var pid = photoId(p);
  var actorSel = self._actorSel;
  var catSel = self._catSel;
  var backdrop = self.querySelector(".zt-lightbox-backdrop");
  if (!backdrop) return;

  function showModal(el) {
    if (!el || typeof bootstrap === "undefined") return;
    bootstrap.Modal.getOrCreateInstance(el, { container: backdrop, focus: true }).show();
  }

  self.querySelectorAll(".zt-lb-modal").forEach(function(el) {
    if (typeof bootstrap !== "undefined") bootstrap.Modal.getOrCreateInstance(el, { container: backdrop });
  });

  var openActor = self.querySelector("#zt-lb-open-actor-modal");
  if (openActor) openActor.onclick = function() { showModal(self.querySelector("#zt-lb-actorSelectionModal")); };
  var openCat = self.querySelector("#zt-lb-open-category-modal");
  if (openCat) openCat.onclick = function() { showModal(self.querySelector("#zt-lb-categorySelectionModal")); };

  var chEdit = self.querySelector("#zt-lb-channel-edit");
  if (chEdit) {
    chEdit.onclick = function() {
      fetch("/api/channel/map", { credentials: "same-origin" }).then(function(r) { return r.json(); })
        .then(function(mapData) {
          var ch = mapData.channels || {};
          var sel = self.querySelector("#zt-lb-channel-list");
          var cur = self.photoChannel(p);
          var curId = cur ? photoId(cur) : "";
          sel.innerHTML = '<option value="">None (inherit photoset)</option>' + Object.keys(ch).map(function(k) {
            return '<option value="' + escAttr(k) + '"' + (k === curId ? " selected" : "") + ">" + esc(ch[k] || "") + "</option>";
          }).join("");
          showModal(self.querySelector("#zt-lb-editChannelModal"));
        });
    };
  }

  var chSend = self.querySelector("#zt-lb-channel-send");
  if (chSend) {
    chSend.onclick = function() {
      var cid = self.querySelector("#zt-lb-channel-list").value;
      var fd = new FormData();
      fd.set("channel_id", cid);
      fetch("/api/photoset/photo/" + pid + "/channel", { method: "POST", credentials: "same-origin", body: fd })
        .then(function(r) {
          if (!r.ok) {
            if (typeof sendToast === "function") sendToast("Channel not updated", "", "bg-danger", "Request failed.");
            return;
          }
          var modalEl = self.querySelector("#zt-lb-editChannelModal");
          var inst = bootstrap.Modal.getInstance(modalEl);
          if (inst) inst.hide();
          return self.refreshPhotos();
        });
    };
  }

  self.querySelectorAll(".zt-lb-add-actor").forEach(function(btn) {
    btn.onclick = function() {
      var aid = this.closest(".add-actor-list").getAttribute("actor-id");
      var aname = (actorSel.actorSelectable[aid] && actorSel.actorSelectable[aid].name) || "Actor";
      fetch("/api/photoset/photo/" + pid + "/actor/" + aid, { method: "PUT", credentials: "same-origin" })
        .then(function(r) {
          if (r.ok) {
            if (typeof sendToast === "function") sendToast("Actor added", "", "bg-success", aname + " added.");
            return self.refreshPhotos();
          }
          if (typeof sendToast === "function") sendToast("Actor not added", "", "bg-danger", "Request failed.");
        });
    };
  });
  self.querySelectorAll(".zt-lb-remove-actor").forEach(function(btn) {
    btn.onclick = function() {
      actorSel.actorDeselect(this.closest(".add-actor-list").getAttribute("actor-id"));
    };
  });

  self.querySelectorAll(".zt-lb-add-category").forEach(function(btn) {
    btn.onclick = function() {
      var cid = this.closest(".add-category-list").getAttribute("category-id");
      var cname = (catSel.categorySelectable[cid] && catSel.categorySelectable[cid].name) || "Category";
      fetch("/api/photoset/photo/" + pid + "/category/" + cid, { method: "PUT", credentials: "same-origin" })
        .then(function(r) {
          if (r.ok) {
            if (typeof sendToast === "function") sendToast("Category added", "", "bg-success", cname + " added.");
            return self.refreshPhotos();
          }
          if (typeof sendToast === "function") sendToast("Category not added", "", "bg-danger", "Request failed.");
        });
    };
  });
  self.querySelectorAll(".zt-lb-remove-category").forEach(function(btn) {
    btn.onclick = function() {
      catSel.categoryDeselect(this.closest(".add-category-list").getAttribute("category-id"));
    };
  });

  var filterInput = self.querySelector("#zt-lb-actor-filter");
  if (filterInput) {
    filterInput.oninput = function() {
      var q = filterInput.value.trim().toLowerCase();
      self.querySelectorAll("#zt-lb-actorSelectionModal .add-actor-list").forEach(function(chip) {
        var aid = chip.getAttribute("actor-id");
        var aname = ((actorSel.actorSelectable[aid] && actorSel.actorSelectable[aid].name) || "").toLowerCase();
        chip.style.display = !q || aname.indexOf(q) >= 0 ? "" : "none";
      });
    };
  }
};

ZtPhotoLightbox.prototype.render = function() {
  var self = this;
  var p = self.current();
  var id = photoId(p);
  var stream = "/api/photo/" + id + "/stream";
  var name = esc(p.Filename || p.filename || "");
  var actors = p.effective_actors || p.EffectiveActors || [];
  var cats = p.effective_categories || p.EffectiveCategories || [];
  var channel = p.effective_channel || p.EffectiveChannel;
  var actorsHtml = actors.map(function(a) {
    var aid = photoId(a);
    var aname = esc(a.Name || a.name || "");
    var tags = (a.Categories || a.categories || []).map(function(c) {
      var cid = photoId(c);
      return '<a class="btn btn-sm btn-secondary me-1" href="/category/' + cid + '">#' + esc(c.Name || c.name || "") + '</a>';
    }).join("");
    return '<span class="me-2 d-inline-flex flex-wrap align-items-center gap-1"><a class="btn btn-sm btn-danger" href="/actor/' + aid + '">@' + aname + '</a>' + tags + '</span>';
  }).join("");
  var catsHtml = cats.map(function(c) {
    return '<a class="btn btn-sm btn-secondary me-1" href="/category/' + photoId(c) + '">#' + esc(c.Name || c.name || "") + '</a>';
  }).join("");
  var chHtml = channel ? '<a class="btn btn-sm btn-dark me-1" href="/channel/' + photoId(channel) + '">' + esc(channel.Name || channel.name || "") + '</a>' : "";

  var editorDisplay = self._editorOpen ? "block" : "none";
  var editorBody = "";
  if (self._admin && self._editorOpen) {
    if (self._catalog) {
      editorBody = self.buildEditorHtml(p, self._catalog);
    } else {
      editorBody = '<p class="small text-muted">Loading editor…</p>';
    }
  }

  var html = '<style>'
    + '.zt-lightbox-backdrop .modal{z-index:2100}'
    + '.zt-lightbox-backdrop .modal-backdrop{z-index:2090}'
    + '#zt-lb-editor .zt-lb-chips,#zt-lb-editor .zt-lb-chips .chip{color:#212529}'
    + '.zt-lb-modal .modal-content,.zt-lb-modal .modal-content .chip,.zt-lb-modal .modal-content h4{color:#212529}'
    + '#zt-lb-editor .form-label.text-light{color:#f8f9fa}'
    + '</style>';
  html += '<div class="zt-lightbox-backdrop" style="position:fixed;inset:0;background:rgba(0,0,0,0.92);z-index:2000;display:flex;flex-direction:column">';
  html += '<div style="display:flex;justify-content:space-between;align-items:center;padding:0.75rem 1rem;color:#fff">';
  html += '<span>' + name + ' (' + (this._index + 1) + '/' + this._photos.length + ')</span>';
  html += '<div><button type="button" class="btn btn-sm btn-outline-light me-2" id="zt-lb-prev">&larr;</button>';
  html += '<button type="button" class="btn btn-sm btn-outline-light me-2" id="zt-lb-next">&rarr;</button>';
  if (this._admin) html += '<button type="button" class="btn btn-sm btn-warning me-2" id="zt-lb-edit">' + (this._editorOpen ? "Hide editor" : "Edit photo") + '</button>';
  html += '<button type="button" class="btn btn-sm btn-light" id="zt-lb-close">&times;</button></div></div>';
  html += '<div style="flex:1;display:flex;align-items:center;justify-content:center;overflow:hidden;padding:0 3rem">';
  html += '<img src="' + stream + '" alt="' + name + '" style="max-width:100%;max-height:calc(100vh - 140px);object-fit:contain">';
  html += '</div>';
  html += '<div style="padding:0.5rem 1rem;color:#fff">' + chHtml + actorsHtml + catsHtml + '</div>';
  if (this._admin) {
    html += '<div id="zt-lb-editor" style="display:' + editorDisplay + ';padding:1rem;background:#222;color:#fff;max-height:40vh;overflow:auto">' + editorBody + '</div>';
  }
  html += '</div>';

  this.innerHTML = html;

  self.addEventListener("click", function(e) {
    var a = e.target.closest("a[href]");
    if (!a || !self.contains(a)) return;
    var href = a.getAttribute("href");
    if (!href || href === "#" || a.target === "_blank" || a.hasAttribute("download")) return;
    if (href.indexOf("/") === 0) self.close();
  }, true);

  this.querySelector("#zt-lb-close").onclick = function() { self.close(); };
  this.querySelector("#zt-lb-prev").onclick = function() { self.prev(); };
  this.querySelector("#zt-lb-next").onclick = function() { self.next(); };

  var editBtn = this.querySelector("#zt-lb-edit");
  if (editBtn) {
    editBtn.onclick = function() {
      if (self._editorOpen) {
        self._editorOpen = false;
        self.render();
        return;
      }
      self._editorOpen = true;
      self.loadCatalog()
        .then(function() { self.render(); })
        .catch(function() {
          if (typeof sendToast === "function") sendToast("Editor", "", "bg-danger", "Failed to load editor data.");
          self._editorOpen = false;
        });
    };
  }

  if (self._admin && self._editorOpen && self._catalog) {
    self.wireEditor();
  }
};

customElements.define("zt-photo-lightbox", ZtPhotoLightbox);
})();
