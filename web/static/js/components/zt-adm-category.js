(function() {
"use strict";
var categoryStyles = ":root{--separator-width:4px}#zt-adm-categories-container section{padding:0 0 20px}#zt-adm-categories-container .category-values{overflow-y:hidden;grid-auto-flow:row;grid-template-columns:repeat(6,minmax(0,1fr));gap:var(--separator-width);display:grid;grid-auto-rows:100px;overflow-x:auto;min-width:0}#zt-adm-categories-container a.category-item{cursor:pointer}#zt-adm-categories-container a{border-radius:8px;display:flex;overflow:hidden;position:relative;min-width:80px}#zt-adm-categories-container img{height:100%;left:0;object-fit:cover;position:absolute;top:0;width:100%}#zt-adm-categories-container a:not(.category-new):after{content:'';position:absolute;inset:0;background:linear-gradient(to bottom,transparent 40%,black 100%)}#zt-adm-categories-container .category-value-header{background:none;bottom:var(--separator-width);display:flex;flex-direction:column;height:auto;left:var(--separator-width);max-width:90%;position:absolute;width:auto;z-index:2}#zt-adm-categories-container h5{color:white;margin:unset}";

function subName(s) {
  var n = s.Name || s.name || s.Title || s.title || "";
  if (typeof n !== "string" || !n.trim()) return "Uncategorized";
  return n.replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/"/g,"&quot;");
}

function toastError(msg) {
  if (typeof sendToast === "function") sendToast("Error", "", "bg-danger", msg);
}

function ZtAdmCategory() {
  var el = Reflect.construct(HTMLElement, [], ZtAdmCategory);
  return el;
}
ZtAdmCategory.prototype = Object.create(HTMLElement.prototype);
ZtAdmCategory.prototype.connectedCallback = function() {
  var self = this;
  fetch("/api/adm/category", { credentials: "same-origin" })
    .then(function(r) { return r.json(); })
    .then(function(data) {
      var items = data.items || [];
      var html = '<style>' + categoryStyles + '</style><div class="row"><div class="col-md-3 col-lg-3"><zt-adm-tabs data-active="categories"></zt-adm-tabs></div><div class="col-md-9 col-lg-9">';
      html += '<div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-certificate"></i></span><h3>Categories</h3><hr /></div>';
      html += '<div style="display:flex;justify-content:space-between;"><h5>Action</h5><div>';
      html += '<button class="btn btn-primary" id="zt-add-category-btn">Add category</button>';
      html += '</div></div><hr /><div class="row"><div class="col-12" id="zt-adm-categories-container">';
      items.forEach(function(c) {
        var subs = c.Sub || c.sub || [];
        var catId = c.ID || c.id;
        var catName = (c.Name || c.name || "Other").replace(/&/g,"&amp;").replace(/</g,"&lt;");
        html += '<section data-cat-id="' + catId + '"><div class="d-flex align-items-center gap-2 mb-2"><h4 class="mb-0">' + catName + '</h4><button type="button" class="btn btn-link p-0 border-0 zt-edit-category-btn text-secondary" data-cat-id="' + catId + '" data-cat-name="' + catName.replace(/"/g,"&quot;") + '" title="Edit category" aria-label="Edit category"><i class="fas fa-pen"></i></button></div><div class="category-values mb-3">';
        subs.forEach(function(s) {
          var sid = s.ID || s.id;
          var name = subName(s);
          var thumbUrl = "/api/category-sub/" + encodeURIComponent(sid) + "/thumb";
          html += '<a href="#" class="category-item" data-sub-id="' + sid + '" data-sub-name="' + name.replace(/"/g,"&quot;") + '"><img src="' + thumbUrl + '" alt=""><div class="category-value-header"><h5>' + name + '</h5></div></a>';
        });
        html += '<a class="category-new" data-parent-id="' + catId + '" style="cursor:pointer"><img src="/static/images/category-add.svg" alt=""><div class="category-value-header"><h5>New</h5></div></a>';
        html += '</div></section>';
      });
      html += "</div></div></div></div>";
      html += '<div class="modal fade" id="zt-add-category-modal" tabindex="-1" aria-hidden="true">'
        + '<div class="modal-dialog"><div class="modal-content">'
        + '<div class="modal-header"><h5 class="modal-title">Add category</h5>'
        + '<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button></div>'
        + '<div class="modal-body"><form id="zt-add-category-form">'
        + '<div class="mb-2"><label class="form-label" for="zt-add-category-name">Name</label>'
        + '<input type="text" class="form-control" id="zt-add-category-name" name="name" required placeholder="Category name"></div>'
        + '</form></div>'
        + '<div class="modal-footer">'
        + '<button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>'
        + '<button type="submit" class="btn btn-primary" form="zt-add-category-form" id="zt-add-category-submit">Create</button>'
        + '</div></div></div></div>';
      html += '<div class="modal fade" id="zt-add-sub-modal" tabindex="-1" aria-hidden="true">'
        + '<div class="modal-dialog"><div class="modal-content">'
        + '<div class="modal-header"><h5 class="modal-title">Add category item</h5>'
        + '<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button></div>'
        + '<div class="modal-body"><form id="zt-add-sub-form">'
        + '<input type="hidden" id="zt-add-sub-parent-id" name="parent">'
        + '<div class="mb-2"><label class="form-label" for="zt-add-sub-name">Name</label>'
        + '<input type="text" class="form-control" id="zt-add-sub-name" name="name" required placeholder="Category item name"></div>'
        + '<div class="mb-2"><label class="form-label" for="zt-add-sub-thumb">Thumbnail</label>'
        + '<input type="file" class="form-control" id="zt-add-sub-thumb" name="thumb" accept="image/*">'
        + '<div class="form-text">Optional. Shown on category tiles.</div></div>'
        + '<div class="mb-2" id="zt-add-sub-thumb-preview-wrap" style="display:none">'
        + '<img id="zt-add-sub-thumb-preview" class="img-thumbnail" alt="" style="max-height:120px;max-width:100%"></div>'
        + '</form></div>'
        + '<div class="modal-footer">'
        + '<button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>'
        + '<button type="submit" class="btn btn-primary" form="zt-add-sub-form" id="zt-add-sub-submit">Create</button>'
        + '</div></div></div></div>';
      html += '<div class="modal fade" id="zt-edit-category-modal" tabindex="-1" aria-hidden="true">'
        + '<div class="modal-dialog"><div class="modal-content">'
        + '<div class="modal-header"><h5 class="modal-title">Edit category</h5>'
        + '<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button></div>'
        + '<div class="modal-body"><form id="zt-edit-category-form">'
        + '<input type="hidden" id="zt-edit-category-id">'
        + '<div class="mb-2"><label class="form-label" for="zt-edit-category-name">Name</label>'
        + '<input type="text" class="form-control" id="zt-edit-category-name" name="name" required></div>'
        + '</form></div>'
        + '<div class="modal-footer">'
        + '<button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>'
        + '<button type="submit" class="btn btn-primary" form="zt-edit-category-form" id="zt-edit-category-submit">Save</button>'
        + '</div></div></div></div>';
      html += '<div class="modal fade" id="zt-edit-sub-modal" tabindex="-1" aria-hidden="true">'
        + '<div class="modal-dialog"><div class="modal-content">'
        + '<div class="modal-header"><h5 class="modal-title">Edit category item</h5>'
        + '<button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button></div>'
        + '<div class="modal-body"><form id="zt-edit-sub-form">'
        + '<input type="hidden" id="zt-edit-sub-id">'
        + '<div class="mb-2"><label class="form-label" for="zt-edit-sub-name">Name</label>'
        + '<input type="text" class="form-control" id="zt-edit-sub-name" name="name" required></div>'
        + '<div class="mb-2"><label class="form-label">Current thumbnail</label>'
        + '<div id="zt-edit-sub-thumb-current-wrap"><img id="zt-edit-sub-thumb-current" class="img-thumbnail" alt="" style="max-height:120px;max-width:100%"></div></div>'
        + '<div class="mb-2"><label class="form-label" for="zt-edit-sub-thumb">Replace thumbnail</label>'
        + '<input type="file" class="form-control" id="zt-edit-sub-thumb" name="thumb" accept="image/*">'
        + '<div class="form-text">Optional. Upload a new image to replace the current one.</div></div>'
        + '<div class="mb-2" id="zt-edit-sub-thumb-preview-wrap" style="display:none">'
        + '<label class="form-label">New thumbnail preview</label>'
        + '<img id="zt-edit-sub-thumb-preview" class="img-thumbnail" alt="" style="max-height:120px;max-width:100%"></div>'
        + '<button type="button" class="btn btn-sm btn-outline-danger" id="zt-edit-sub-thumb-remove">Remove thumbnail</button>'
        + '</form></div>'
        + '<div class="modal-footer">'
        + '<button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancel</button>'
        + '<button type="submit" class="btn btn-primary" form="zt-edit-sub-form" id="zt-edit-sub-submit">Save</button>'
        + '</div></div></div></div>';
      self.innerHTML = html;

      var addCategoryModal = self.querySelector("#zt-add-category-modal");
      var addCategoryForm = self.querySelector("#zt-add-category-form");
      var addSubModal = self.querySelector("#zt-add-sub-modal");
      var addSubForm = self.querySelector("#zt-add-sub-form");
      var subThumbInput = self.querySelector("#zt-add-sub-thumb");
      var subThumbPreview = self.querySelector("#zt-add-sub-thumb-preview");
      var subThumbPreviewWrap = self.querySelector("#zt-add-sub-thumb-preview-wrap");

      function openAddCategoryModal() {
        if (addCategoryForm) addCategoryForm.reset();
        if (addCategoryModal && typeof bootstrap !== "undefined") {
          bootstrap.Modal.getOrCreateInstance(addCategoryModal).show();
          setTimeout(function() {
            var nameInput = self.querySelector("#zt-add-category-name");
            if (nameInput) nameInput.focus();
          }, 400);
        }
      }

      function openAddSubModal(parentId) {
        if (addSubForm) addSubForm.reset();
        var parentInput = self.querySelector("#zt-add-sub-parent-id");
        if (parentInput) parentInput.value = parentId || "";
        if (subThumbPreviewWrap) subThumbPreviewWrap.style.display = "none";
        if (subThumbPreview) subThumbPreview.removeAttribute("src");
        if (addSubModal && typeof bootstrap !== "undefined") {
          bootstrap.Modal.getOrCreateInstance(addSubModal).show();
          setTimeout(function() {
            var nameInput = self.querySelector("#zt-add-sub-name");
            if (nameInput) nameInput.focus();
          }, 400);
        }
      }

      if (subThumbInput) {
        subThumbInput.addEventListener("change", function() {
          var file = subThumbInput.files && subThumbInput.files[0];
          if (!file || !subThumbPreview || !subThumbPreviewWrap) {
            if (subThumbPreviewWrap) subThumbPreviewWrap.style.display = "none";
            return;
          }
          var reader = new FileReader();
          reader.onload = function(ev) {
            subThumbPreview.src = ev.target.result;
            subThumbPreviewWrap.style.display = "block";
          };
          reader.readAsDataURL(file);
        });
      }

      self.querySelectorAll("a.category-new").forEach(function(link) {
        link.addEventListener("click", function(e) {
          e.preventDefault();
          openAddSubModal(link.getAttribute("data-parent-id"));
        });
      });

      var addBtn = self.querySelector("#zt-add-category-btn");
      if (addBtn) addBtn.addEventListener("click", openAddCategoryModal);

      if (addCategoryForm) {
        addCategoryForm.addEventListener("submit", function(e) {
          e.preventDefault();
          var nameInput = self.querySelector("#zt-add-category-name");
          var submitBtn = self.querySelector("#zt-add-category-submit");
          var name = nameInput && nameInput.value.trim();
          if (!name) return;
          if (submitBtn) submitBtn.disabled = true;
          var fd = new FormData();
          fd.set("Name", name);
          fetch("/api/category", { method: "POST", credentials: "same-origin", body: fd })
            .then(function(r) {
              if (r.ok) {
                var inst = addCategoryModal && bootstrap.Modal.getInstance(addCategoryModal);
                if (inst) inst.hide();
                self.connectedCallback();
              } else {
                return r.json().then(function(data) {
                  toastError((data && data.error) || "Failed to add category");
                }).catch(function() { toastError("Failed to add category"); });
              }
            })
            .catch(function() { toastError("Failed to add category"); })
            .finally(function() { if (submitBtn) submitBtn.disabled = false; });
        });
      }

      if (addSubForm) {
        addSubForm.addEventListener("submit", function(e) {
          e.preventDefault();
          var nameInput = self.querySelector("#zt-add-sub-name");
          var parentInput = self.querySelector("#zt-add-sub-parent-id");
          var submitBtn = self.querySelector("#zt-add-sub-submit");
          var name = nameInput && nameInput.value.trim();
          var parentId = parentInput && parentInput.value;
          if (!name || !parentId) return;
          if (submitBtn) submitBtn.disabled = true;
          var thumbFile = subThumbInput && subThumbInput.files && subThumbInput.files[0];
          var fd = new FormData();
          fd.set("Name", name);
          fd.set("Parent", parentId);
          fetch("/api/category-sub", { method: "POST", credentials: "same-origin", body: fd })
            .then(function(r) {
              return r.json().then(function(data) {
                return { ok: r.ok, data: data };
              });
            })
            .then(function(res) {
              if (!res.ok) {
                toastError((res.data && res.data.error) || "Failed to add category item");
                return;
              }
              var subId = res.data && (res.data.id || res.data.ID);
              if (!thumbFile || !subId) {
                var inst = addSubModal && bootstrap.Modal.getInstance(addSubModal);
                if (inst) inst.hide();
                self.connectedCallback();
                return;
              }
              var thumbFd = new FormData();
              thumbFd.set("pp", thumbFile);
              return fetch("/api/category-sub/" + encodeURIComponent(subId) + "/thumb", {
                method: "POST",
                credentials: "same-origin",
                body: thumbFd
              }).then(function(r) {
                if (r.ok) {
                  var inst = addSubModal && bootstrap.Modal.getInstance(addSubModal);
                  if (inst) inst.hide();
                  self.connectedCallback();
                } else {
                  return r.json().then(function(data) {
                    toastError((data && data.human_error) || (data && data.error) || "Category item created but thumbnail upload failed");
                    self.connectedCallback();
                  });
                }
              });
            })
            .catch(function() { toastError("Failed to add category item"); })
            .finally(function() { if (submitBtn) submitBtn.disabled = false; });
        });
      }

      var editCategoryModal = self.querySelector("#zt-edit-category-modal");
      var editCategoryForm = self.querySelector("#zt-edit-category-form");
      var editSubModal = self.querySelector("#zt-edit-sub-modal");
      var editSubForm = self.querySelector("#zt-edit-sub-form");
      var editSubThumbInput = self.querySelector("#zt-edit-sub-thumb");
      var editSubThumbPreview = self.querySelector("#zt-edit-sub-thumb-preview");
      var editSubThumbPreviewWrap = self.querySelector("#zt-edit-sub-thumb-preview-wrap");
      var editSubThumbCurrent = self.querySelector("#zt-edit-sub-thumb-current");
      var editSubThumbRemoveBtn = self.querySelector("#zt-edit-sub-thumb-remove");
      var editSubRemoveThumb = false;

      function openEditCategoryModal(catId, name) {
        var idInput = self.querySelector("#zt-edit-category-id");
        var nameInput = self.querySelector("#zt-edit-category-name");
        if (idInput) idInput.value = catId || "";
        if (nameInput) nameInput.value = name || "";
        if (editCategoryModal && typeof bootstrap !== "undefined") {
          bootstrap.Modal.getOrCreateInstance(editCategoryModal).show();
          setTimeout(function() { if (nameInput) nameInput.focus(); }, 400);
        }
      }

      function openEditSubModal(subId, name) {
        editSubRemoveThumb = false;
        if (editSubForm) editSubForm.reset();
        var idInput = self.querySelector("#zt-edit-sub-id");
        var nameInput = self.querySelector("#zt-edit-sub-name");
        if (idInput) idInput.value = subId || "";
        if (nameInput) nameInput.value = name || "";
        if (editSubThumbPreviewWrap) editSubThumbPreviewWrap.style.display = "none";
        if (editSubThumbPreview) editSubThumbPreview.removeAttribute("src");
        if (editSubThumbCurrent && subId) {
          editSubThumbCurrent.style.opacity = "1";
          editSubThumbCurrent.src = "/api/category-sub/" + encodeURIComponent(subId) + "/thumb?t=" + Date.now();
        }
        if (editSubModal && typeof bootstrap !== "undefined") {
          bootstrap.Modal.getOrCreateInstance(editSubModal).show();
          setTimeout(function() { if (nameInput) nameInput.focus(); }, 400);
        }
      }

      self.querySelectorAll(".zt-edit-category-btn").forEach(function(btn) {
        btn.addEventListener("click", function(e) {
          e.preventDefault();
          e.stopPropagation();
          openEditCategoryModal(btn.getAttribute("data-cat-id"), btn.getAttribute("data-cat-name"));
        });
      });

      self.querySelectorAll("a.category-item").forEach(function(link) {
        link.addEventListener("click", function(e) {
          e.preventDefault();
          openEditSubModal(link.getAttribute("data-sub-id"), link.getAttribute("data-sub-name"));
        });
      });

      if (editSubThumbInput) {
        editSubThumbInput.addEventListener("change", function() {
          editSubRemoveThumb = false;
          var file = editSubThumbInput.files && editSubThumbInput.files[0];
          if (!file || !editSubThumbPreview || !editSubThumbPreviewWrap) {
            if (editSubThumbPreviewWrap) editSubThumbPreviewWrap.style.display = "none";
            return;
          }
          var reader = new FileReader();
          reader.onload = function(ev) {
            editSubThumbPreview.src = ev.target.result;
            editSubThumbPreviewWrap.style.display = "block";
          };
          reader.readAsDataURL(file);
        });
      }

      if (editSubThumbRemoveBtn) {
        editSubThumbRemoveBtn.addEventListener("click", function() {
          editSubRemoveThumb = true;
          if (editSubThumbInput) editSubThumbInput.value = "";
          if (editSubThumbPreviewWrap) editSubThumbPreviewWrap.style.display = "none";
          if (editSubThumbPreview) editSubThumbPreview.removeAttribute("src");
          if (editSubThumbCurrent) editSubThumbCurrent.style.opacity = "0.35";
        });
      }

      if (editCategoryForm) {
        editCategoryForm.addEventListener("submit", function(e) {
          e.preventDefault();
          var catId = self.querySelector("#zt-edit-category-id") && self.querySelector("#zt-edit-category-id").value;
          var name = self.querySelector("#zt-edit-category-name") && self.querySelector("#zt-edit-category-name").value.trim();
          var submitBtn = self.querySelector("#zt-edit-category-submit");
          if (!catId || !name) return;
          if (submitBtn) submitBtn.disabled = true;
          var fd = new FormData();
          fd.set("title", name);
          fetch("/api/category/" + encodeURIComponent(catId) + "/rename", { method: "POST", credentials: "same-origin", body: fd })
            .then(function(r) {
              if (r.ok) {
                var inst = editCategoryModal && bootstrap.Modal.getInstance(editCategoryModal);
                if (inst) inst.hide();
                self.connectedCallback();
              } else {
                return r.json().then(function(data) {
                  toastError((data && data.error) || "Failed to update category");
                }).catch(function() { toastError("Failed to update category"); });
              }
            })
            .catch(function() { toastError("Failed to update category"); })
            .finally(function() { if (submitBtn) submitBtn.disabled = false; });
        });
      }

      if (editSubForm) {
        editSubForm.addEventListener("submit", function(e) {
          e.preventDefault();
          var subId = self.querySelector("#zt-edit-sub-id") && self.querySelector("#zt-edit-sub-id").value;
          var name = self.querySelector("#zt-edit-sub-name") && self.querySelector("#zt-edit-sub-name").value.trim();
          var submitBtn = self.querySelector("#zt-edit-sub-submit");
          if (!subId || !name) return;
          if (submitBtn) submitBtn.disabled = true;
          var thumbFile = editSubThumbInput && editSubThumbInput.files && editSubThumbInput.files[0];
          var removeThumb = editSubRemoveThumb;
          var fd = new FormData();
          fd.set("title", name);
          fetch("/api/category-sub/" + encodeURIComponent(subId) + "/rename", { method: "POST", credentials: "same-origin", body: fd })
            .then(function(r) {
              if (!r.ok) {
                return r.json().then(function(data) {
                  toastError((data && data.error) || "Failed to update category item");
                  return Promise.reject();
                });
              }
              if (removeThumb) {
                return fetch("/api/category-sub/" + encodeURIComponent(subId) + "/thumb", {
                  method: "DELETE",
                  credentials: "same-origin"
                }).then(function(delR) {
                  if (!delR.ok) {
                    return delR.json().then(function(data) {
                      toastError((data && data.human_error) || (data && data.error) || "Name saved but thumbnail removal failed");
                      return Promise.reject();
                    });
                  }
                });
              }
              if (thumbFile) {
                var thumbFd = new FormData();
                thumbFd.set("pp", thumbFile);
                return fetch("/api/category-sub/" + encodeURIComponent(subId) + "/thumb", {
                  method: "POST",
                  credentials: "same-origin",
                  body: thumbFd
                }).then(function(upR) {
                  if (!upR.ok) {
                    return upR.json().then(function(data) {
                      toastError((data && data.human_error) || (data && data.error) || "Name saved but thumbnail upload failed");
                      return Promise.reject();
                    });
                  }
                });
              }
            })
            .then(function() {
              var inst = editSubModal && bootstrap.Modal.getInstance(editSubModal);
              if (inst) inst.hide();
              self.connectedCallback();
            })
            .catch(function(err) {
              if (err) toastError("Failed to update category item");
            })
            .finally(function() { if (submitBtn) submitBtn.disabled = false; });
        });
      }

      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() { self.innerHTML = '<div class="alert alert-danger">Failed.</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); });
};
customElements.define("zt-adm-category", ZtAdmCategory);
})();
