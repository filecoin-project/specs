// iife to avoid polluting the global.
(function () {
  // Run me as soon as possible after the the css links are in the dom.
  // This assumes this js file is added to the page after the css links.
  const lightMode = document.getElementById('light-mode-link')
  const darkMode = document.getElementById('dark-mode-link')
  const btn = document.querySelector('.dark-mode-toggle')
  let theme = 'light'

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
  
  // wait for localstorage...
  const previousChoice = localStorage.getItem('theme')
  theme = previousChoice || theme
  
  // Light is default, so enable dark if user previously chose it.
  if (theme === 'dark') {
    enableDarkMode()
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
