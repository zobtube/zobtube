(function() {
"use strict";
var categoryListStyles = ":root{--separator-width:4px}#zt-categories-container section{padding:0 0 20px}#zt-categories-container .category-values{overflow-y:hidden;grid-auto-columns:calc(16.66667% - var(--separator-width)*5/6);grid-auto-flow:row;grid-template-columns:repeat(6,calc(16.66667% - var(--separator-width)*5/6));gap:var(--separator-width);display:grid;grid-auto-rows:100px;overflow-x:scroll}#zt-categories-container a{border-radius:8px;display:flex;overflow:hidden;position:relative}#zt-categories-container img{height:100%;left:0;object-fit:cover;position:absolute;top:0;width:100%}#zt-categories-container a:not(.category-new):after{content:'';position:absolute;inset:0;background:linear-gradient(to bottom,transparent 40%,black 100%)}#zt-categories-container .category-value-header{background:none;bottom:var(--separator-width);display:flex;flex-direction:column;height:auto;left:var(--separator-width);max-width:90%;position:absolute;width:auto;z-index:2}#zt-categories-container h5{color:white;margin:unset}";
function subName(s) {
  var n = s.Name || s.name || s.Title || s.title || "";
  if (typeof n !== "string" || !n.trim()) return "Uncategorized";
  return n.replace(/&/g,"&amp;").replace(/</g,"&lt;").replace(/"/g,"&quot;");
}
function ZtCategoryList() {
  var el = Reflect.construct(HTMLElement, [], ZtCategoryList);
  return el;
}
ZtCategoryList.prototype = Object.create(HTMLElement.prototype);
ZtCategoryList.prototype.connectedCallback = function() {
  var self = this;
  var admin = (window.__USER__ && window.__USER__.admin);
  fetch("/api/category", { credentials: "same-origin" })
    .then(function(r) { return r.json(); })
    .then(function(data) {
      var items = data.items || [];
      var html = '<style>' + categoryListStyles + '</style>';
      html += '<div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-th-large"></i></span><h3>Categories' +
        (admin ? ' <a href="/adm/categories"><i class="fas fa-cog" title="Manage categories"></i></a>' : '') + '</h3><hr /></div>';
      html += '<div class="row"><div id="zt-categories-container">';
      items.forEach(function(c) {
        var subs = c.Sub || c.sub || [];
        if (subs.length === 0) return;
        var catName = (c.Name || c.name || "Other").replace(/&/g,"&amp;").replace(/</g,"&lt;");
        html += '<section><h4>' + catName + '</h4><div class="category-values mb-3">';
        subs.forEach(function(s) {
          var sid = s.ID || s.id;
          var name = subName(s);
          var thumbUrl = "/api/category-sub/" + encodeURIComponent(sid) + "/thumb";
          html += '<a href="/category/' + encodeURIComponent(sid) + '"><img class="lazy" data-src="' + thumbUrl + '" alt=""><div class="category-value-header"><h5>' + name + '</h5></div></a>';
        });
        html += '</div></section>';
      });
      html += "</div></div>";
      self.innerHTML = html;
      if (window.zt && window.zt.loadLazyIn) window.zt.loadLazyIn(self);
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    })
    .catch(function() {
      self.innerHTML = '<div class="alert alert-danger">Failed to load.</div>';
      if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
    });
};
customElements.define("zt-category-list", ZtCategoryList);
})();
