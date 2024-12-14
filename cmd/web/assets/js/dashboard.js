document.addEventListener("DOMContentLoaded", () => {
  // Get DOM elements
  const mobileMenu = document.getElementById("dashboard__mobile__menu");
  const menuButton = document.getElementById("dashboard__menu__btn");
  const closeButton = document.getElementById("dashboard__close__btn");

  // Toggle menu function
  function toggleMenu() {
    const isOpen = !mobileMenu.classList.contains("-translate-x-full");

    // Toggle translation class
    mobileMenu.classList.toggle("-translate-x-full", isOpen);

    // Update ARIA attributes
    mobileMenu.setAttribute("aria-hidden", isOpen);
    menuButton.setAttribute("aria-expanded", !isOpen);

    // Prevent body scroll when menu is open
    document.body.style.overflow = isOpen ? "auto" : "hidden";
  }

  // Add event listeners if elements exist
  if (menuButton && closeButton && mobileMenu) {
    menuButton.addEventListener("click", toggleMenu);
    closeButton.addEventListener("click", toggleMenu);

    // Close menu on escape key
    document.addEventListener("keydown", (e) => {
      if (
        e.key === "Escape" &&
        !mobileMenu.classList.contains("-translate-x-full")
      ) {
        toggleMenu();
      }
    });

    // Close menu when clicking outside
    mobileMenu.addEventListener("click", (e) => {
      if (e.target === mobileMenu) {
        toggleMenu();
      }
    });
  }
});
