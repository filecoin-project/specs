// iife to avoid polluting the global.
(function () {
  // Run me as soon as possible after the the css links are in the dom.
  // This assumes this js file is added to the page after the css links.
  const lightMode = document.getElementById('light-mode-link')
  const darkMode = document.getElementById('dark-mode-link')
  const btn = document.querySelector('.dark-mode-toggle')
  const prefersDarkScheme = window.matchMedia('(prefers-color-scheme: dark)').matches 
  let theme = prefersDarkScheme ? 'dark' : 'light'

  function enableLightMode () {
    lightMode.disabled = false
    darkMode.disabled = true
    btn.setAttribute('aria-pressed', theme === 'dark')
  }

  function enableDarkMode () {
    darkMode.disabled = false
    lightMode.disabled = true
    btn.setAttribute('aria-pressed', theme === 'dark')
  }
  
  // enable dark theme optimistically on OS with dark theme enabled to reduce flashing of white theme.
  if (prefersDarkScheme) {
    enableDarkMode()
  }

  // wait for localstorage...
  const previousChoice = localStorage.getItem('theme')
  theme = previousChoice || theme
  
  // Light is default, so enable dark if user previously chose it but their OS pref is light.
  if (theme === 'dark') {
    enableDarkMode()
  } else {
    // needed to catch those OS darkmoders who want their specs to be light mode.
    enableLightMode()
  }

  // set up the toggle once the DOM is ready.
  document.addEventListener("DOMContentLoaded", function(event) {
    
    btn.addEventListener('click', function () {
      theme = (theme === 'light' ? 'dark' : 'light')
      if (theme === 'dark') {
        enableDarkMode()
      } else {
        enableLightMode()
      }
      localStorage.setItem('theme', theme)
    })
    // init the button state to match the currently selected theme.
    btn.setAttribute('aria-pressed', theme === 'dark')
  });  
})()
