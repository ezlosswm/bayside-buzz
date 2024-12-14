document.addEventListener("DOMContentLoaded", () => {
  // Mobile menu
  const menuButton = document.querySelector("[data-collapse-toggle]");
  const closeButton = document.querySelector("[data-drawer-close]");
  const menu = document.getElementById("navbar-sticky");

  function toggleMenu() {
    menu?.classList.toggle("hidden");
  }

  menuButton?.addEventListener("click", toggleMenu);
  closeButton?.addEventListener("click", toggleMenu);

  // footer
  const year = document.getElementById("year");
  if (year) {
    year.innerText = new Date().getFullYear().toString();
  }
});
