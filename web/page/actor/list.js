{{ define "actor/list.js" }}
function filterActorsByName(e) {
  console.log(e.target.value);

  // create regex
  var re = new RegExp(e.target.value, 'i');

  // iterate on all cols
  cols = document.getElementsByClassName('col');
  for (const col of cols) {
    a = col.querySelector('.stretched-link');
    console.debug(a);
    if (re.test(a.innerText)) {
      col.style.display = '';
    } else {
      col.style.display = 'none';
    }
  }
}

window.zt.onload.push(function actorListSetupNameFilter() {
  // reset
  document.getElementById('actor-filter').value = '';

  // bind input
  document.getElementById('actor-filter').addEventListener('input', filterActorsByName);
});

{{ end }}
