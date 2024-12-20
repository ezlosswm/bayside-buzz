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
