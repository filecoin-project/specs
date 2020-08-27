import zoomable from 'd3-zoomable';

function lightbox () {
    const transitionSpeedInMilliseconds = 250;

    // template
    const fragment = new DocumentFragment()
    const container = document.createElement('div')
    container.classList.add('lightbox-container')
  
    const zoom = document.createElement('div')
    zoom.classList.add('lightbox-zoom')
    
    const img = document.createElement('img')
    img.src = ''
    img.classList.add('lightbox-image')
    
    container.appendChild(zoom)
    zoom.appendChild(img)
    fragment.appendChild(container)
    document.body.appendChild(fragment)
  
    // init zoomable in the template
    const myZoom = zoomable()
    myZoom(container).htmlEl(zoom)
    
    // hook events
    const elements = document.querySelectorAll('.zoomable img')
    elements.forEach((element) => {
      element.addEventListener('click', () => {
          handleElementClick(element);
      });
    });
    container.addEventListener('click', hideLightbox);
    window.addEventListener('keyup', (event) => {
        if (event.key === 'Escape') {
            hideLightbox();
        }
    });
    
    function handleElementClick(htmlElement) {
      img.attributes['src'].value = htmlElement.attributes['src'].value;
      if (!isVisible()) {
        container.classList.remove('hidden');
        container.classList.add('visible');
        document.body.classList.add('lightbox-body-scroll-stop')
      }
    }
  
    function hideLightbox() {
      if (isVisible()) {
          container.classList.add('hidden');
  
          setTimeout(() => {
              container.classList.remove('visible');
              container.classList.remove('hidden');
              document.body.classList.remove('lightbox-body-scroll-stop')
              img.attributes['src'].value = '';
              myZoom.zoomReset()
          }, transitionSpeedInMilliseconds);
      }
    }
  
    function isVisible() {
      return container.classList.contains('visible');
    }
  }

  export {
      lightbox
  }