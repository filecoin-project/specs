
document.addEventListener("DOMContentLoaded", function(event) {

  // Swap the active css based on the users color-scheme preference.
  const btn = document.querySelector('.dark-mode-toggle')
  const lightMode = document.getElementById('light-mode-link')
  const darkMode = document.getElementById('dark-mode-link')

  const prefersDarkScheme = window.matchMedia('(prefers-color-scheme: dark)') 
  const previousChoice = localStorage.getItem('theme')
  let theme = previousChoice || (prefersDarkScheme.matches ? 'dark' : 'light')

  // Listen for a click on the button 
  btn.addEventListener('click', function () {
    // toggle all the things
    theme = (theme === 'light' ? 'dark' : 'light')
    console.log(theme)
    lightMode.disabled = theme === 'dark'
    darkMode.disabled = theme !== 'dark'
    btn.setAttribute('aria-pressed', theme === 'dark')
    localStorage.setItem('theme', theme)
  })

  function enableDarkMode () {
    darkMode.disabled = false
    lightMode.disabled = true
    btn.setAttribute('aria-pressed', true)
    localStorage.setItem('theme', 'dark')
  }

  // Light is default, so enable dark if user previously chose it.
  if (theme === 'dark') {
    console.log('snap')
    enableDarkMode()
  }

});
