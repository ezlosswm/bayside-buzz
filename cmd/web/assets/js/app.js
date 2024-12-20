// Mobile menu
document.addEventListener("DOMContentLoaded", () => {
  // Get DOM elements
  const mobileMenu = document.getElementById("mobile__menu");
  const menuButton = document.getElementById("menu__btn");
  const closeButton = document.getElementById("close__btn");

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

    // footer
    const year = document.getElementById("year");
    if (year) {
        year.innerText = new Date().getFullYear().toString();
    }

    // shareable link
    const shareBtn = document.getElementById("share-button");

    async function copyText(e) {
        e.preventDefault();

        try {
            await navigator.clipboard.writeText(window.location.href);
            alert("Link copied to clipboard!");
        } catch (err) {
            console.error(err);
        }
    }

    shareBtn?.addEventListener("click", copyText);
});
