(function() {
"use strict";
function ZtProfileSettings() {
  var el = Reflect.construct(HTMLElement, [], ZtProfileSettings);
  return el;
}
ZtProfileSettings.prototype = Object.create(HTMLElement.prototype);
ZtProfileSettings.prototype.connectedCallback = function() {
  var self = this;
  var html = '<div class="row"><div class="col-12"><zt-profile-tabs data-active="settings"></zt-profile-tabs></div></div>';
  html += '<div class="row"><div class="col-md-12"><div class="themeix-section-h"><span class="heading-icon"><i class="fa fa-cog"></i></span><h3>Settings</h3><hr /></div>';
  html += '<form id="zt-profile-change-password-form" style="margin-top:1rem;max-width:400px">';
  html += '<div class="mb-3"><label for="zt-current-password" class="form-label">Current password</label><input type="password" class="form-control" id="zt-current-password" name="current_password" required autocomplete="current-password"></div>';
  html += '<div class="mb-3"><label for="zt-new-password" class="form-label">New password</label><input type="password" class="form-control" id="zt-new-password" name="new_password" required autocomplete="new-password"></div>';
  html += '<div class="mb-3"><label for="zt-confirm-password" class="form-label">Confirm new password</label><input type="password" class="form-control" id="zt-confirm-password" name="confirm_password" required autocomplete="new-password"></div>';
  html += '<div id="zt-password-form-error" class="alert alert-danger mb-3" style="display:none" role="alert"></div>';
  html += '<button type="submit" class="btn btn-primary">Change password</button></form></div></div>';
  self.innerHTML = html;

  var form = self.querySelector("#zt-profile-change-password-form");
  var errEl = self.querySelector("#zt-password-form-error");
  form.addEventListener("submit", function(e) {
    e.preventDefault();
    var current = self.querySelector("#zt-current-password").value;
    var newPw = self.querySelector("#zt-new-password").value;
    var confirmPw = self.querySelector("#zt-confirm-password").value;
    errEl.style.display = "none";
    errEl.textContent = "";
    if (newPw !== confirmPw) {
      errEl.textContent = "New password and confirmation do not match.";
      errEl.style.display = "block";
      return;
    }
    if (!newPw || !current) {
      errEl.textContent = "Current and new password are required.";
      errEl.style.display = "block";
      return;
    }
    fetch("/api/profile/password", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "same-origin",
      body: JSON.stringify({ current_password: current, new_password: newPw })
    }).then(function(r) {
      if (r.status === 401) {
        window.location.href = "/auth/login?next=" + encodeURIComponent(window.location.pathname);
        return Promise.resolve();
      }
      return r.json().then(function(data) {
        if (r.ok) {
          if (typeof sendToast === "function") sendToast("Password changed", "", "bg-success", "Your password has been updated.");
          form.reset();
          errEl.style.display = "none";
        } else {
          errEl.textContent = (data && data.error) || "Failed to change password.";
          errEl.style.display = "block";
        }
      });
    }).catch(function() {
      errEl.textContent = "Request failed.";
      errEl.style.display = "block";
    });
  });

  if (window.zt && window.zt.pageReady) window.zt.pageReady(self);
};
customElements.define("zt-profile-settings", ZtProfileSettings);
})();
