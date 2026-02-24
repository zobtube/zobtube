(function() {
"use strict";
function ZtProfileTabs() {
  var el = Reflect.construct(HTMLElement, [], ZtProfileTabs);
  return el;
}
ZtProfileTabs.prototype = Object.create(HTMLElement.prototype);
ZtProfileTabs.prototype.connectedCallback = function() {
  var active = this.getAttribute("data-active") || "";
  function link(label, href, tab) {
    var c = tab === active ? "nav-link active" : "nav-link";
    return '<li class="nav-item"><a class="' + c + '" href="' + (tab === active ? "#" : href) + '">' + label + '</a></li>';
  }
  this.innerHTML = '<div class="col-md-12"><div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-user"></i></span><h3>Your account</h3><hr /></div></div>' +
    '<div class="col-md-12 mb-4"><ul class="nav nav-tabs">' +
    link("Most viewed", "/profile", "most-viewed") +
    link("Settings", "/profile/settings", "settings") +
    '</ul></div>';
};
customElements.define("zt-profile-tabs", ZtProfileTabs);
})();
