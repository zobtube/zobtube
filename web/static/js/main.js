(function ($) {

   "use strict";


   $(window).on("load", function () {
      $('.preloader').fadeOut(1000);
   });

   $(".slider").pgwSlider({

      autoSlide: false
   });

   $(".slider-2").pgwSlideshow();


   $(".video-carousel").owlCarousel({
      loop: true,
      margin: 30,
      nav: true,
      navText: ['<i class="fas fa-angle-left"></i>', '<i class="fas fa-angle-right"></i>'],
      responsive: {

         0: {

            items: 1
         },

         480: {

            items: 2
         },

         768: {

            items: 3
         },
         992: {
            items: 4
         }
      }
   });

   $(".feature-carousel").owlCarousel({
      loop: true,
      margin: 15,
      nav: false,
      items: 1,
      autoplay: true,
      autoplayTimeout: 1000,
      autoplayHoverPause: true

   });


   // Youtube Video Bg

   var bgvideo = $(".bgvideo");

   if (bgvideo.length > 0) {

      bgvideo.YTPlayer({

         videoURL: 'https://www.youtube.com/watch?v=RoKeSWzZAwA',
         containment: '.video-area',
         quality: 'large',
         autoPlay: true,
         mute: true,
         opacity: 1

      });
   }




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




   $('.lazy').Lazy({
      placeholder: "data:image/gif;base64,R0lGODlhEALAPQAPzl5uLr9Nrl8e7..."
   });

  /*  Footer full year */
  $('#spanYear').html(new Date().getFullYear());



})(jQuery);
