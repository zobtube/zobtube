{{ define "adm/category.js" }}

function modal_category_delete_show(id, value) {
  document.getElementById('modal-category-delete-id').value = id;
  document.getElementById('modal-category-delete-value').innerText = value;
  window.ModalCategoryDelete.show();
}

function modal_category_delete_send(category_id) {
  category_id = document.getElementById('modal-category-delete-id').value;

  $.ajax("/api/adm/category/"+category_id, {
    method: 'DELETE',
    processData: false,
    contentType: false,

    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },

    success: function () {
      console.debug('success, redirecting');
      window.location.reload();
    },

    error: function () {
      console.debug('failed');
    },
  });
}

function modal_category_value_show(parent) {
  document.getElementById('modal-category-value-parent').value = parent;
  window.ModalCategoryValue.show();
}

function modal_category_send() {
  var formData = new FormData(document.getElementById("modal-category-form"));

  $.ajax("/api/adm/category", {
    method: 'POST',
    data: formData,
    processData: false,
    contentType: false,

    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },

    success: function () {
      console.debug('success, redirecting');
      window.location.reload();
    },

    error: function () {
      console.debug('failed');
    },
  });
}

function modal_category_value_send() {
  var formData = new FormData(document.getElementById("modal-category-value-form"));

  $.ajax("/api/adm/category-sub", {
    method: 'POST',
    data: formData,
    processData: false,
    contentType: false,

    xhr: function () {
      var xhr = new XMLHttpRequest();
      return xhr;
    },

    success: function () {
      console.debug('success, redirecting');
      // erase form content
      document.getElementById('modal-category-value-name').value = '';
      window.location.reload();
    },

    error: function () {
      console.debug('failed');
    },
  });
}

/* -- start subcategory edit -- */
function offcanvas_subcategory_edit_show(id, name, thumbnail) {
  window.globals['sub_id'] = id;

  // set subcategory name
  document.getElementById('offcanvas-subcategory-edit-title').value = name;

  // display thumbnail if set
  document.getElementById('offcanvas-subcategory-edit-thumbnail-div').style.display = thumbnail ? '' : 'none';
  if (thumbnail) {
    document.getElementById('offcanvas-subcategory-edit-thumbnail-img').src = '/category/'+id+'/thumb';
    document.getElementById('offcanvas-subcategory-edit-thumbnail-delete').onclick = function() {
      offcanvas_subcategory_edit_delete_thumbnail(id);
    };
  }

  // display new thumbnail buttons
  document.getElementById('offcanvas-subcategory-edit-thumbnail-new').style.display = thumbnail ? 'none' : '';

  // hide new thumbnail content
  document.getElementById('offcanvas-subcategory-edit-thumbnail-new-from-url').style.display = 'none';
  document.getElementById('offcanvas-subcategory-edit-thumbnail-new-from-upload').style.display = 'none';

  // show
  window.OffCanvasSubCategoryEdit.show();
}

function offcanvas_subcategory_edit_delete_thumbnail(id) {
  console.debug('delete subcategory id '+id+' thumbnail');
  $.ajax('/api/category-sub/'+id+'/thumb', {
    method: 'DELETE',
    success: function () {
      console.debug('success, redirecting');
      window.location.reload();
    },
    error: function (data) {
      sendToast('Deletion of sub-category thumbnail', 'failed', 'bg-danger', data.responseJSON.error);
    },
  });
}

function offcanvas_subcategory_edit_new_from_url_show() {
  document.getElementById('offcanvas-subcategory-edit-thumbnail-new').style.display = 'none';
  document.getElementById('offcanvas-subcategory-edit-thumbnail-new-from-url').style.display = '';
}

function offcanvas_subcategory_edit_new_from_url_send() {
  url = document.getElementById('offcanvas-subcategory-edit-thumbnail-new-from-url-input').value;
  $.ajax(url, {
    xhr: function() {
      var xhr = new XMLHttpRequest();
      xhr.onreadystatechange = function() {
        if (xhr.readyState == 2) {
          if (xhr.status == 200) {
            xhr.responseType = "blob";
          } else {
            xhr.responseType = "text";
          }
        }
      };
      return xhr;
    },

    success: function (data) {
      var formData = new FormData();
      formData.append('pp', data);
      $.ajax("/api/category-sub/"+window.globals['sub_id']+"/thumb", {
        method: 'POST',
        data: formData,
        processData: false,
        contentType: false,

        success: function () {
          console.debug('success, redirecting');
          window.location.reload();
        },

        error: function (data) {
          console.debug('failed');
          console.error(data);
        },
      });
    },
    error: function (data) {
      sendToast('Use an URL as new sub-category thumbnail', 'failed', 'bg-danger', data.responseJSON.error);
    },
  });

}

function offcanvas_subcategory_edit_new_from_upload_show() {
  document.getElementById('offcanvas-subcategory-edit-thumbnail-new').style.display = 'none';
  document.getElementById('offcanvas-subcategory-edit-thumbnail-new-from-upload').style.display = '';
}

function offcanvas_subcategory_edit_new_from_upload_send() {
  input = document.getElementById('offcanvas-subcategory-edit-thumbnail-new-from-upload-input');
  var formData = new FormData();
  formData.append('pp', input.files[0]);
  $.ajax("/api/category-sub/"+window.globals['sub_id']+"/thumb", {
    method: 'POST',
    data: formData,
    processData: false,
    contentType: false,

    success: function () {
      console.debug('success, redirecting');
      window.location.reload();
    },

    error: function (data) {
      console.debug('failed');
      console.error(data);
    },
  });
}

function offcanvas_subcategory_title_edit() {
  title = document.getElementById('offcanvas-subcategory-edit-title-input');
  title.disabled = false;

  btn = document.getElementById('offcanvas-subcategory-edit-title-button');
  btn.classList.remove('btn-outline-warning');
  btn.classList.add('btn-outline-success');
  btn.innerText = 'Send';
  btn.onclick = offcanvas_subcategory_title_send;
}

function offcanvas_subcategory_title_send() {
  $.ajax("/api/category-sub/"+window.globals['sub_id']+"/rename", {
    method: 'POST',
    data: {
      'title': document.getElementById('offcanvas-subcategory-edit-title-input').value,
    },
    success: function () {
      console.debug('success, redirecting');
      window.location.reload();
    },
    error: function (data) {
      console.debug('failed');
      console.error(data);
    },
  });
}
/* -- end subcategory edit -- */

function modal_register() {
  // register modal category
  window.ModalCategory = new bootstrap.Modal(document.getElementById('modal-category'));

  // register enter cmd
  document.getElementById('modal-category-name').addEventListener("keydown", function (e) {
    if (e.code === "Enter") {
      modal_category_send();
    }
  });

  // register modal category value
  window.ModalCategoryValue = new bootstrap.Modal(document.getElementById('modal-category-value'));
  // register enter cmd
  document.getElementById('modal-category-value-name').addEventListener("keydown", function (e) {
    if (e.code === "Enter") {
      modal_category_value_send();
    }
  });

  // register modal category delete
  window.ModalCategoryDelete = new bootstrap.Modal(document.getElementById('modal-category-delete'));
}

function offcanvas_register() {
  // register offcanvas for subcategory edition
  window.OffCanvasSubCategoryEdit = new bootstrap.Offcanvas(document.getElementById('offcanvas-subcategory-edit'));
}

window.onload = function() {
  window.globals = {};
  modal_register();
  offcanvas_register();
}

{{ end }}
