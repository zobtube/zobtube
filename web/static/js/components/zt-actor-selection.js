(function() {
"use strict";
function ZtActorSelection() {
  var el = Reflect.construct(HTMLElement, [], ZtActorSelection);
  return el;
}
ZtActorSelection.prototype = Object.create(HTMLElement.prototype);
ZtActorSelection.prototype.connectedCallback = function() {
  var actors = JSON.parse(this.getAttribute("data-actors") || "[]");
  var selected = JSON.parse(this.getAttribute("data-selected") || "[]");
  var id = this.getAttribute("data-id") || "actor_chip_selector";
  var html = '<div class="form-floating"><div class="form-control chip-selector" id="' + id + '" style="height:unset;display:flex;"><div class="chips">';
  actors.forEach(function(a) {
    var aid = a.ID || a.id;
    var show = selected.indexOf(aid) >= 0 ? "" : "display:none";
    var name = (a.Name || a.name || "").replace(/&/g,"&amp;").replace(/</g,"&lt;");
    html += '<div class="chip video-actor-list" actor-id="' + aid + '" style="' + show + '"><img src="/api/actor/' + aid + '/thumb" width="50" height="50">' + name +
      '<button class="btn btn-danger" onclick="window.zt.actorSelection&&window.zt.actorSelection.actorDeselect(\'' + aid + '\');"><i class="fa fa-trash-alt"></i></button></div>';
  });
  html += '<div class="chip">Add an actor<button class="btn btn-success" data-bs-toggle="modal" data-bs-target="#actorSelectionModal"><i class="fa fa-plus-circle"></i></button></div></div></div><label for="actors">Actors</label></div>';
  this.innerHTML = html;
};
customElements.define("zt-actor-selection", ZtActorSelection);
})();
