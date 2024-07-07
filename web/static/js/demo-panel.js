(function($){

 "use strict";

// Demo Panel 

 $(".setting-icons span").on("click", function(){

        $(".demo-wrapper").toggleClass("collapse-left");
 });




 $(".styleswitcher-list li a").on("click", function(e){

     var colorCode =   $(this).data("color");

     e.preventDefault();

       $(".header-top,.hover-bg,.heading-icon,.watch-btn,.subscribe-form button").css({"background-color":colorCode});


 });



 

}(jQuery));