(function() {
"use strict";
var categoryStyles = ":root{--separator-width:4px}#zt-adm-categories-container section{padding:0 0 20px}#zt-adm-categories-container .category-values{overflow-y:hidden;grid-auto-flow:row;grid-template-columns:repeat(6,minmax(0,1fr));gap:var(--separator-width);display:grid;grid-auto-rows:100px;overflow-x:auto;min-width:0}#zt-adm-categories-container a{border-radius:8px;display:flex;overflow:hidden;position:relative;min-width:80px}#zt-adm-categories-container img{height:100%;left:0;object-fit:cover;position:absolute;top:0;width:100%}#zt-adm-categories-container a:not(.category-new):after{content:'';position:absolute;inset:0;background:linear-gradient(to bottom,transparent 40%,black 100%)}#zt-adm-categories-container .category-value-header{background:none;bottom:var(--separator-width);display:flex;flex-direction:column;height:auto;left:var(--separator-width);max-width:90%;position:absolute;width:auto;z-index:2}#zt-adm-categories-container h5{color:white;margin:unset}";

function subName(s) {
  var n = s.Name || s.name || s.Title || s.title || "";
  if (typeof n !== "string" || !n.trim()) return "Uncategorized";
  return n.replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/"/g,"&quot;");
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
      var html = '<style>' + categoryStyles + '</style>';
      html += '<h1>Categories</h1><hr /><div style="display:flex;justify-content:space-between;"><h5>Action</h5><div>';
      html += '<button class="btn btn-primary" id="zt-add-category-btn">Add category</button>';
      html += '</div></div><hr /><div class="row"><div class="col-12" id="zt-adm-categories-container">';
      items.forEach(function(c) {
        var subs = c.Sub || c.sub || [];
        var catId = c.ID || c.id;
        var catName = (c.Name || c.name || "Other").replace(/&/g,"&amp;").replace(/</g,"&lt;");
        html += '<section><h4>' + catName + '</h4><div class="category-values mb-3">';
        subs.forEach(function(s) {
          var sid = s.ID || s.id;
          var name = subName(s);
          var thumbUrl = "/api/category-sub/" + encodeURIComponent(sid) + "/thumb";
          html += '<a data-sub-id="' + sid + '" data-sub-name="' + name.replace(/"/g,"&quot;") + '"><img src="' + thumbUrl + '" alt=""><div class="category-value-header"><h5>' + name + '</h5></div></a>';
        });
        html += '<a class="category-new" data-parent-id="' + catId + '" style="cursor:pointer"><img src="/static/images/category-add.svg" alt=""><div class="category-value-header"><h5>New</h5></div></a>';
        html += '</div></section>';
      });
      html += "</div></div>";
      self.innerHTML = html;

      self.querySelectorAll("a.category-new").forEach(function(link) {
        link.onclick = function(e) {
          e.preventDefault();
          var parentId = link.getAttribute("data-parent-id");
          var name = prompt("Name of the new category item");
          if (!name) return;
          var fd = new FormData();
          fd.set("Name", name);
          fd.set("Parent", parentId);
          fetch("/api/category-sub", { method: "POST", credentials: "same-origin", body: fd })
            .then(function(r) { if (r.ok) self.connectedCallback(); });
        };
      });

      var addBtn = self.querySelector("#zt-add-category-btn");
      addBtn.onclick = function() {
        var name = prompt("Name of the new category");
        if (!name) return;
        var fd = new FormData();
        fd.set("Name", name);
        fetch("/api/category", { method: "POST", credentials: "same-origin", body: fd })
          .then(function(r) { if (r.ok) { self.connectedCallback(); } });
      };
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() { self.innerHTML = '<div class="alert alert-danger">Failed.</div>'; if (window.zt && window.zt.pageReady) window.zt.pageReady(self); });
};
customElements.define("zt-adm-category", ZtAdmCategory);
})();
