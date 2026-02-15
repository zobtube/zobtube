(function() {
"use strict";
function ZtCategorySelection() {
  var el = Reflect.construct(HTMLElement, [], ZtCategorySelection);
  return el;
}
ZtCategorySelection.prototype = Object.create(HTMLElement.prototype);
ZtCategorySelection.prototype.connectedCallback = function() {
  var cats = JSON.parse(this.getAttribute("data-categories") || "[]");
  var selected = JSON.parse(this.getAttribute("data-selected") || "[]");
  var id = this.getAttribute("data-id") || "category_chip_selector";
  var html = '<div class="form-floating"><div class="form-control chip-selector" id="' + id + '" style="height:unset;display:flex;"><div class="chips">';
  cats.forEach(function(c) {
    (c.Sub || c.sub || []).forEach(function(s) {
      var sid = s.ID || s.id;
      var show = selected.indexOf(sid) >= 0 ? "" : "display:none";
      var thumb = (s.Thumbnail || s.thumbnail) ? '<img src="/api/category-sub/' + encodeURIComponent(sid) + '/thumb" width="100" height="50">' : "";
      var name = (s.Name || s.name || "").replace(/&/g,"&amp;").replace(/</g,"&lt;");
      html += '<div class="chip video-category-list" category-id="' + sid + '" style="' + show + '">' + thumb + name +
        '<button class="btn btn-danger" onclick="window.zt.categorySelection&&window.zt.categorySelection.categoryDeselect(\'' + sid + '\');"><i class="fa fa-trash-alt"></i></button></div>';
    });
  });
  html += '<div class="chip">Add a category<button class="btn btn-success" data-bs-toggle="modal" data-bs-target="#categorySelectionModal"><i class="fa fa-plus-circle"></i></button></div></div></div><label for="categories">Categories</label></div>';
  this.innerHTML = html;
};
customElements.define("zt-category-selection", ZtCategorySelection);
})();
