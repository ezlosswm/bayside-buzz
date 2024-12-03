document.addEventListener("DOMContentLoaded", () => {
  const mobileMenu = document.getElementById("mobile__menu");
  const menuButton = document.getElementById("menu__btn");
  const closeButton = document.getElementById("close__btn");

  function toggleMenu() {
    const isOpen = mobileMenu.classList.contains("translate-x-0");
    mobileMenu.classList.toggle("translate-x-0", !isOpen);
    mobileMenu.classList.toggle("-translate-x-full", isOpen);
    mobileMenu.setAttribute("aria-hidden", isOpen);
    menuButton.setAttribute("aria-expanded", !isOpen);
  }

  menuButton.addEventListener("click", toggleMenu);
  closeButton.addEventListener("click", toggleMenu);

  // footer
  const year = document.getElementById("year");

  year.innerText = new Date().getFullYear().toString();
});
