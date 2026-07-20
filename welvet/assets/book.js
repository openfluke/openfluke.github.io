(function () {
  var btn = document.getElementById("menu-btn");
  if (btn) {
    btn.addEventListener("click", function () {
      document.body.classList.toggle("nav-open");
    });
  }
  document.addEventListener("click", function (e) {
    if (!document.body.classList.contains("nav-open")) return;
    var side = document.querySelector(".sidebar");
    if (side && !side.contains(e.target) && e.target !== btn) {
      document.body.classList.remove("nav-open");
    }
  });
  var path = location.pathname.split("/").pop() || "index.html";
  document.querySelectorAll(".nav-group a").forEach(function (a) {
    var href = a.getAttribute("href") || "";
    if (href.endsWith(path)) a.classList.add("active");
  });
})();
