(function() {
"use strict";
function esc(s) {
  return String(s).replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/"/g,"&quot;");
}
function toast(title, body, cls) {
  if (typeof sendToast === "function") sendToast(title, "", cls || "bg-success", body);
}
function actorProfileLink(actorId, name) {
  var label = esc(name || actorId || "");
  if (!actorId) return label;
  return '<a href="/actor/' + encodeURIComponent(actorId) + '" target="_blank" rel="noopener noreferrer">' + label + '</a>';
}
function ZtAdmActorDuplicates() {
  var el = Reflect.construct(HTMLElement, [], ZtAdmActorDuplicates);
  return el;
}
ZtAdmActorDuplicates.prototype = Object.create(HTMLElement.prototype);

ZtAdmActorDuplicates.prototype.connectedCallback = function() {
  var self = this;
  if (!(window.__USER__ && window.__USER__.admin)) {
    self.innerHTML = '<div class="alert alert-danger">Forbidden</div>';
    if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    return;
  }
  self.innerHTML = '<p class="text-muted">Loading duplicate actors...</p>';
  self._loadAll();
};

ZtAdmActorDuplicates.prototype._loadAll = function() {
  var self = this;
  Promise.all([
    fetch("/api/adm/actor/duplicates", { credentials: "same-origin" }).then(function(r) {
      if (!r.ok) throw new Error("duplicates");
      return r.json();
    }),
    fetch("/api/adm/actor/duplicates/dismissed", { credentials: "same-origin" }).then(function(r) {
      if (!r.ok) throw new Error("dismissed");
      return r.json();
    }),
    fetch("/api/adm/actor", { credentials: "same-origin" }).then(function(r) {
      if (!r.ok) throw new Error("actors");
      return r.json();
    })
  ]).then(function(results) {
    self._groups = results[0].groups || [];
    self._dismissed = results[1].items || [];
    self._allActors = results[2].items || [];
    self._render();
    if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
  }).catch(function() {
    self.innerHTML = '<div class="alert alert-danger">Failed to load duplicate actors.</div>';
    if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
  });
};

ZtAdmActorDuplicates.prototype._mergeActor = function(sourceId, targetId) {
  return fetch("/api/actor/" + encodeURIComponent(sourceId) + "/merge", {
    method: "POST",
    credentials: "same-origin",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ target_id: targetId })
  }).then(function(r) {
    if (!r.ok) {
      return r.json().catch(function() { return {}; }).then(function(d) {
        throw new Error((d && d.error) || "merge failed");
      });
    }
    return r.json();
  });
};

ZtAdmActorDuplicates.prototype._dismissPair = function(id1, id2) {
  return fetch("/api/adm/actor/duplicates/dismiss", {
    method: "POST",
    credentials: "same-origin",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ actor_id_1: id1, actor_id_2: id2 })
  }).then(function(r) {
    if (!r.ok && r.status !== 409) {
      return r.json().catch(function() { return {}; }).then(function(d) {
        throw new Error((d && d.error) || "dismiss failed");
      });
    }
    return r.json().catch(function() { return {}; });
  });
};

ZtAdmActorDuplicates.prototype._showConfirmModal = function(title, message, onConfirm) {
  var modalEl = this.querySelector("#zt-merge-confirm-modal");
  if (!modalEl || typeof bootstrap === "undefined") {
    if (typeof onConfirm === "function") onConfirm();
    return;
  }
  this._confirmCallback = onConfirm;
  var titleEl = this.querySelector("#zt-merge-confirm-title");
  var bodyEl = this.querySelector("#zt-merge-confirm-body");
  if (titleEl) titleEl.textContent = title || "Confirm";
  if (bodyEl) bodyEl.textContent = message || "";
  bootstrap.Modal.getOrCreateInstance(modalEl).show();
};

ZtAdmActorDuplicates.prototype._bindConfirmModal = function() {
  var self = this;
  var modalEl = self.querySelector("#zt-merge-confirm-modal");
  var confirmBtn = self.querySelector("#zt-merge-confirm-btn");
  if (!modalEl || !confirmBtn) return;
  confirmBtn.addEventListener("click", function() {
    var cb = self._confirmCallback;
    self._confirmCallback = null;
    var inst = bootstrap.Modal.getInstance(modalEl);
    if (inst) inst.hide();
    if (typeof cb === "function") cb();
  });
};

ZtAdmActorDuplicates.prototype._bindTabs = function() {
  var self = this;
  var tabDup = self.querySelector("#zt-tab-duplicates");
  var tabDis = self.querySelector("#zt-tab-dismissed");
  var tabMan = self.querySelector("#zt-tab-manual");
  self.querySelectorAll(".nav-link[data-tab]").forEach(function(btn) {
    btn.addEventListener("click", function() {
      self.querySelectorAll(".nav-link[data-tab]").forEach(function(b) { b.classList.remove("active"); });
      btn.classList.add("active");
      var tab = btn.getAttribute("data-tab");
      if (tabDup) tabDup.style.display = tab === "duplicates" ? "" : "none";
      if (tabDis) tabDis.style.display = tab === "dismissed" ? "" : "none";
      if (tabMan) tabMan.style.display = tab === "manual" ? "" : "none";
    });
  });
};

ZtAdmActorDuplicates.prototype._selectDefaultMergeTargets = function() {
  this.querySelectorAll(".zt-dup-group").forEach(function(card) {
    var radios = card.querySelectorAll(".zt-merge-target");
    if (radios.length === 0) return;
    radios.forEach(function(r) { r.checked = false; });
    radios[0].checked = true;
  });
};

ZtAdmActorDuplicates.prototype._bindGroups = function() {
  var self = this;
  self.querySelectorAll(".zt-dup-group").forEach(function(card) {
    var gi = card.getAttribute("data-group-index");
    var group = (self._groups || [])[parseInt(gi, 10)];
    if (!group) return;

    card.querySelector(".zt-merge-group").addEventListener("click", function() {
      var targetInput = card.querySelector(".zt-merge-target:checked");
      if (!targetInput) {
        toast("Merge actors", "Select which actor to keep.", "bg-warning");
        return;
      }
      var targetId = targetInput.value;
      var sources = (group.actors || []).map(function(a) { return a.id || a.ID; }).filter(function(id) { return id !== targetId; });
      if (sources.length === 0) return;
      self._showConfirmModal(
        "Confirm merge",
        "Merge " + sources.length + " actor(s) into the selected profile? This cannot be undone.",
        function() {
          var chain = Promise.resolve();
          sources.forEach(function(sourceId) {
            chain = chain.then(function() { return self._mergeActor(sourceId, targetId); });
          });
          chain.then(function() {
            toast("Actors merged", "Profiles combined successfully.", "bg-success");
            self._loadAll();
          }).catch(function(e) {
            toast("Merge failed", e.message || "Unexpected error", "bg-danger");
            self._loadAll();
          });
        }
      );
    });

    card.querySelectorAll(".zt-dismiss-pair").forEach(function(btn) {
      btn.addEventListener("click", function() {
        self._dismissPair(btn.getAttribute("data-id1"), btn.getAttribute("data-id2")).then(function() {
          toast("Marked as different", "This pair will not be suggested again.", "bg-success");
          self._loadAll();
        }).catch(function(e) {
          toast("Failed", e.message || "Could not dismiss pair", "bg-danger");
        });
      });
    });

    var dismissSelectBtn = card.querySelector(".zt-dismiss-pair-select");
    if (dismissSelectBtn) {
      dismissSelectBtn.addEventListener("click", function() {
        var id1 = card.querySelector(".zt-dismiss-a").value;
        var id2 = card.querySelector(".zt-dismiss-b").value;
        if (id1 === id2) {
          toast("Mark as different", "Choose two different actors.", "bg-warning");
          return;
        }
        self._dismissPair(id1, id2).then(function() {
          toast("Marked as different", "This pair will not be suggested again.", "bg-success");
          self._loadAll();
        }).catch(function(e) {
          toast("Failed", e.message || "Could not dismiss pair", "bg-danger");
        });
      });
    }
  });
};

ZtAdmActorDuplicates.prototype._bindDismissed = function() {
  var self = this;
  self.querySelectorAll(".zt-reenable-dismiss").forEach(function(btn) {
    btn.addEventListener("click", function() {
      fetch("/api/adm/actor/duplicates/dismiss/" + encodeURIComponent(btn.getAttribute("data-dismiss-id")), {
        method: "DELETE",
        credentials: "same-origin"
      }).then(function(r) {
        if (r.ok) {
          toast("Pair re-enabled", "These actors may appear as duplicates again.", "bg-success");
          self._loadAll();
        } else {
          toast("Failed", "Could not remove dismissal", "bg-danger");
        }
      });
    });
  });
};

ZtAdmActorDuplicates.prototype._bindManual = function() {
  var self = this;
  var mergeBtn = self.querySelector(".zt-manual-merge");
  if (!mergeBtn) return;
  mergeBtn.addEventListener("click", function() {
    var sourceId = self.querySelector(".zt-manual-source").value;
    var targetId = self.querySelector(".zt-manual-target").value;
    if (sourceId === targetId) {
      toast("Manual merge", "Source and target must be different.", "bg-warning");
      return;
    }
    self._showConfirmModal(
      "Confirm merge",
      "Merge source actor into target? The source profile will be deleted.",
      function() {
        self._mergeActor(sourceId, targetId).then(function() {
          toast("Actors merged", "Profiles combined successfully.", "bg-success");
          self._loadAll();
        }).catch(function(e) {
          toast("Merge failed", e.message || "Unexpected error", "bg-danger");
        });
      }
    );
  });
};

ZtAdmActorDuplicates.prototype._renderGroupActors = function(group, gi) {
  var parts = ['<div class="row">'];
  (group.actors || []).forEach(function(a, ai) {
    var id = a.id || a.ID;
    var checked = ai === 0 ? ' checked="checked"' : '';
    var thumb = "/api/actor/" + encodeURIComponent(id) + "/thumb";
    parts.push('<div class="col-md-4 mb-3">');
    parts.push('<img class="rounded mb-2" src="' + thumb + '" style="width:80px;height:80px;object-fit:cover" onerror="this.style.display=\'none\'">');
    parts.push('<div class="form-check"><input class="form-check-input zt-merge-target" type="radio" name="merge-target-' + gi + '" value="' + esc(id) + '" id="target-' + gi + '-' + esc(id) + '"' + checked + '">');
    parts.push('<label class="form-check-label" for="target-' + gi + '-' + esc(id) + '">Keep: ' + esc(a.name || a.Name || "") + '</label></div>');
    parts.push('<div class="small text-muted">' + (a.video_count || 0) + ' video(s)</div>');
    parts.push('<br><a href="/actor/' + encodeURIComponent(id) + '/edit" class="small">Edit profile</a>');
    parts.push('</div>');
  });
  parts.push('</div>');
  return parts.join('');
};

ZtAdmActorDuplicates.prototype._renderGroupActions = function(group) {
  var parts = ['<div class="mt-2 d-flex flex-wrap gap-2 align-items-end">'];
  parts.push('<button type="button" class="btn btn-primary zt-merge-group">Merge into selected</button>');
  if ((group.actors || []).length === 2) {
    var a0 = group.actors[0], a1 = group.actors[1];
    var id0 = a0.id || a0.ID, id1 = a1.id || a1.ID;
    parts.push('<button type="button" class="btn btn-outline-secondary zt-dismiss-pair" data-id1="' + esc(id0) + '" data-id2="' + esc(id1) + '">Mark as different actors</button>');
  } else {
    parts.push('<div class="d-flex gap-2 align-items-center flex-wrap"><select class="form-select form-select-sm zt-dismiss-a" style="max-width:200px">');
    (group.actors || []).forEach(function(a) {
      var id = a.id || a.ID;
      parts.push('<option value="' + esc(id) + '">' + esc(a.name || a.Name || id) + '</option>');
    });
    parts.push('</select><span>and</span><select class="form-select form-select-sm zt-dismiss-b" style="max-width:200px">');
    (group.actors || []).forEach(function(a) {
      var id = a.id || a.ID;
      parts.push('<option value="' + esc(id) + '">' + esc(a.name || a.Name || id) + '</option>');
    });
    parts.push('</select><button type="button" class="btn btn-outline-secondary btn-sm zt-dismiss-pair-select">Mark pair as different</button></div>');
  }
  parts.push('</div>');
  return parts.join('');
};

ZtAdmActorDuplicates.prototype._render = function() {
  var self = this;
  var groups = self._groups || [];
  var dismissed = self._dismissed || [];
  var allActors = self._allActors || [];
  var parts = [];

  parts.push('<div class="row"><div class="col-md-3 col-lg-3"><zt-adm-tabs data-active="actor-duplicates"></zt-adm-tabs></div><div class="col-md-9 col-lg-9">');
  parts.push('<div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-user"></i></span><h3>Actor duplicates</h3><hr /></div>');
  parts.push('<p class="text-muted">Actors with identical names (case-insensitive). Merge to keep one profile and transfer videos, or mark pairs as different actors.</p>');
  parts.push('<ul class="nav nav-tabs mb-3"><li class="nav-item"><button class="nav-link active" type="button" data-tab="duplicates">Duplicates (' + groups.length + ')</button></li>');
  parts.push('<li class="nav-item"><button class="nav-link" type="button" data-tab="dismissed">Dismissed (' + dismissed.length + ')</button></li>');
  parts.push('<li class="nav-item"><button class="nav-link" type="button" data-tab="manual">Manual merge</button></li></ul>');

  parts.push('<div id="zt-tab-duplicates">');
  if (groups.length === 0) {
    parts.push('<div class="alert alert-info">No duplicate actor names pending review.</div>');
  }
  groups.forEach(function(group, gi) {
    parts.push('<div class="card mb-3 zt-dup-group" data-group-index="' + gi + '"><div class="card-header"><strong>' + esc(group.name) + '</strong> <span class="badge text-bg-secondary">' + (group.actors || []).length + ' actors</span></div><div class="card-body">');
    parts.push(self._renderGroupActors(group, gi));
    parts.push(self._renderGroupActions(group));
    parts.push('</div></div>');
  });
  parts.push('</div>');

  parts.push('<div id="zt-tab-dismissed" style="display:none">');
  if (dismissed.length === 0) {
    parts.push('<div class="alert alert-info">No dismissed pairs.</div>');
  } else {
    parts.push('<table class="table table-striped"><thead><tr><th>Actor 1</th><th>Actor 2</th><th>Dismissed</th><th></th></tr></thead><tbody>');
    dismissed.forEach(function(item) {
      var id = item.id || item.ID;
      var actorId1 = item.actor_id_1 || item.ActorID1 || "";
      var actorId2 = item.actor_id_2 || item.ActorID2 || "";
      parts.push('<tr><td>' + actorProfileLink(actorId1, item.actor_name_1 || item.ActorName1) + '</td><td>' + actorProfileLink(actorId2, item.actor_name_2 || item.ActorName2) + '</td><td>' + esc(item.created_at || '') + '</td><td style="text-align:end"><button type="button" class="btn btn-sm btn-warning zt-reenable-dismiss" data-dismiss-id="' + esc(id) + '">Re-enable</button></td></tr>');
    });
    parts.push('</tbody></table>');
  }
  parts.push('</div>');

  parts.push('<div id="zt-tab-manual" style="display:none">');
  parts.push('<p class="text-muted">Merge any two actors manually (not limited to duplicate names).</p>');
  parts.push('<div class="row g-2 align-items-end"><div class="col-md-4"><label class="form-label">Source (removed)</label><select class="form-select zt-manual-source">');
  allActors.forEach(function(a) {
    var id = a.ID || a.id;
    parts.push('<option value="' + esc(id) + '">' + esc(a.Name || a.name || id) + '</option>');
  });
  parts.push('</select></div><div class="col-md-4"><label class="form-label">Target (kept)</label><select class="form-select zt-manual-target">');
  allActors.forEach(function(a) {
    var id = a.ID || a.id;
    parts.push('<option value="' + esc(id) + '">' + esc(a.Name || a.name || id) + '</option>');
  });
  parts.push('</select></div><div class="col-md-4"><button type="button" class="btn btn-primary zt-manual-merge">Merge</button></div></div>');
  parts.push('</div>');

  parts.push('</div></div>');

  parts.push('<div class="modal fade" id="zt-merge-confirm-modal" tabindex="-1" aria-hidden="true">');
  parts.push('<div class="modal-dialog"><div class="modal-content">');
  parts.push('<div class="modal-header"><h5 class="modal-title" id="zt-merge-confirm-title">Confirm merge</h5>');
  parts.push('<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button></div>');
  parts.push('<div class="modal-body" id="zt-merge-confirm-body"></div>');
  parts.push('<div class="modal-footer">');
  parts.push('<button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>');
  parts.push('<button type="button" class="btn btn-danger" id="zt-merge-confirm-btn">Merge</button>');
  parts.push('</div></div></div></div>');

  self.innerHTML = parts.join('');


  self._selectDefaultMergeTargets();
  self._bindConfirmModal();
  self._bindTabs();
  self._bindGroups();
  self._bindDismissed();
  self._bindManual();
};

customElements.define("zt-adm-actor-duplicates", ZtAdmActorDuplicates);
})();
