(function ($) {

   "use strict";


   $(window).on("load", function () {
      $('.preloader').fadeOut(1000);
   });

   // Scroll Top
   function scrolltop() {
      var wind = $(window);
      wind.on("scroll", function () {
         var scrollTop = $(window).scrollTop();
         if (scrollTop >= 500) {
            $(".scroll-top").fadeIn("slow");
         } else {
            $(".scroll-top").fadeOut("slow");
         }

      });

      $(".scroll-top").on("click", function () {
         var bodyTop = $("html, body");
         bodyTop.animate({
            scrollTop: 0
         }, 800, "swing");
      });

   }
   scrolltop();

   window.lazyLoadInstance = new LazyLoad({});

})(jQuery);

// common aynsc ajax request
function ajax(url, method, data) {
  return new Promise(function(resolve, reject) {
    $.ajax(url, {
      method: method,
      xhr: function () {
        var xhr = new XMLHttpRequest();
        return xhr;
      },
      success: function (res) {
        resolve(res);
      },
      error: function (res) {
        reject(res);
      },
    });
  });
}

// common toast function
function sendToast(title, subtitle, title_color, message) {
   // get const
   toastContainer = document.getElementById('zt-toast-container');
   toastTemplate = document.getElementById('toastTemplate');

   // create new toast
   newtoast = toastTemplate.cloneNode(true);

   // set title
   n_title = newtoast.getElementsByClassName('zt-toast-title');
   n_title[0].innerText = title;

   // set subtitle
   n_subtitle = newtoast.getElementsByClassName('zt-toast-subtitle');
   n_subtitle[0].innerText = subtitle;

   // set color
   n_header = newtoast.getElementsByClassName('toast-header');
   n_header[0].className += " "+title_color

   // set content
   n_content = newtoast.getElementsByClassName('zt-toast-body');
   n_content[0].innerHTML = message;

   toastContainer.appendChild(newtoast);
   toast = new bootstrap.Toast(newtoast);
   toast.show();
}

// at last, calling all boot functions
window.zt.onload.forEach(function(item){
  // call each onload function
  item();
});
