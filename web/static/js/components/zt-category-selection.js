(function() {
"use strict";
function ZtCategorySelection() {
  var el = Reflect.construct(HTMLElement, [], ZtCategorySelection);
  return el;
}
ZtCategorySelection.prototype = Object.create(HTMLElement.prototype);
ZtCategorySelection.prototype.connectedCallback = function() {
  var esc = window.ztEscHtml;
  var cats = JSON.parse(this.getAttribute("data-categories") || "[]");
  var selected = JSON.parse(this.getAttribute("data-selected") || "[]");
  var id = this.getAttribute("data-id") || "category_chip_selector";
  var html = '<div class="form-floating"><div class="form-control chip-selector" id="' + esc(id) + '" style="height:unset;display:flex;"><div class="chips">';
  cats.forEach(function(c) {
    (c.Sub || c.sub || []).forEach(function(s) {
      var sid = s.ID || s.id;
      var show = selected.indexOf(sid) >= 0 ? "" : "display:none";
      var thumb = (s.Thumbnail || s.thumbnail) ? '<img src="/api/category-sub/' + encodeURIComponent(sid) + '/thumb" width="100" height="50">' : "";
      var name = esc(s.Name || s.name || "");
      // sid goes through the deselect button's data-category-id below rather
      // than an inline handler, so it never has to survive an HTML-then-JS parse.
      html += '<div class="chip video-category-list" category-id="' + esc(sid) + '" style="' + show + '">' + thumb + name +
        '<button type="button" class="btn btn-danger zt-category-deselect" data-category-id="' + esc(sid) + '"><i class="fa fa-trash-alt"></i></button></div>';
    });
  });
  html += '<div class="chip">Add a category<button class="btn btn-success" data-bs-toggle="modal" data-bs-target="#categorySelectionModal"><i class="fa fa-plus-circle"></i></button></div></div></div><label for="categories">Categories</label></div>';
  this.innerHTML = html;
  this.querySelectorAll(".zt-category-deselect").forEach(function(btn) {
    btn.addEventListener("click", function() {
      if (window.zt && window.zt.categorySelection) {
        window.zt.categorySelection.categoryDeselect(btn.getAttribute("data-category-id"));
      }
    });
  });
};
customElements.define("zt-category-selection", ZtCategorySelection);
})();
