(function() {
"use strict";
function ZtActorSelection() {
  var el = Reflect.construct(HTMLElement, [], ZtActorSelection);
  return el;
}
ZtActorSelection.prototype = Object.create(HTMLElement.prototype);
ZtActorSelection.prototype.connectedCallback = function() {
  var esc = window.ztEscHtml;
  var actors = JSON.parse(this.getAttribute("data-actors") || "[]");
  var selected = JSON.parse(this.getAttribute("data-selected") || "[]");
  var id = this.getAttribute("data-id") || "actor_chip_selector";
  var html = '<div class="form-floating"><div class="form-control chip-selector" id="' + esc(id) + '" style="height:unset;display:flex;"><div class="chips">';
  actors.forEach(function(a) {
    var aid = a.ID || a.id;
    var show = selected.indexOf(aid) >= 0 ? "" : "display:none";
    var name = esc(a.Name || a.name || "");
    // aid goes through the deselect button's data-actor-id below rather than an
    // inline handler, so it never has to survive an HTML-then-JS parse.
    html += '<div class="chip video-actor-list" actor-id="' + esc(aid) + '" style="' + show + '"><img src="/api/actor/' + encodeURIComponent(aid) + '/thumb" width="50" height="50">' + name +
      '<button type="button" class="btn btn-danger zt-actor-deselect" data-actor-id="' + esc(aid) + '"><i class="fa fa-trash-alt"></i></button></div>';
  });
  html += '<div class="chip">Add an actor<button class="btn btn-success" data-bs-toggle="modal" data-bs-target="#actorSelectionModal"><i class="fa fa-plus-circle"></i></button></div></div></div><label for="actors">Actors</label></div>';
  this.innerHTML = html;
  this.querySelectorAll(".zt-actor-deselect").forEach(function(btn) {
    btn.addEventListener("click", function() {
      if (window.zt && window.zt.actorSelection) {
        window.zt.actorSelection.actorDeselect(btn.getAttribute("data-actor-id"));
      }
    });
  });
};
customElements.define("zt-actor-selection", ZtActorSelection);
})();
