{{ define "actor/create.js" }}

function validateMinimalDetails() {
  ActorNameInput = document.getElementById('actor-name');
  ActorSexSelect = document.getElementById('actor-sex');
  PreFormNextButton = document.getElementById('preform-next');

  if ((ActorNameInput.value.length > 0) && (ActorSexSelect.value != 'xx')) {
    PreFormNextButton.classList.remove('disabled');
  } else {
    PreFormNextButton.classList.add('disabled');
  }
}
validateMinimalDetails();

{{ end }}
