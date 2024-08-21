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
         }, 800, "easeOutCubic");
      });

   }
   scrolltop();

   window.lazyLoadInstance = new LazyLoad({});

})(jQuery);
