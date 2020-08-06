/**
 * Tocbot
 * Tocbot creates a toble of contents based on HTML headings on a page,
 * this allows users to easily jump to different sections of the document.
 * Tocbot was inspired by tocify (http://gregfranko.com/jquery.tocify.js/).
 * The main differences are that it works natively without any need for jquery or jquery UI).
 *
 * @author Tim Scanlin
 */

/* globals define */

// Default options.
var defaultOptions = require('./default-options.js')
var debounce = require('debounce-fn')
// Object to store current options.
var options = {}
// Object for public APIs.
var tocbot = {}

var BuildHtml = require('./build-html.js')
var ParseContent = require('./parse-content.js')
var updateTocScroll = require('./update-toc-scroll.js')
// Keep these variables at top scope once options are passed in.
var buildHtml
var parseContent

// Just return if its not a browser.
var supports = true
var headingsArray

// From: https://github.com/Raynos/xtend
var hasOwnProperty = Object.prototype.hasOwnProperty
function extend () {
    var target = {}
    for (var i = 0; i < arguments.length; i++) {
    var source = arguments[i]
    for (var key in source) {
        if (hasOwnProperty.call(source, key)) {
        target[key] = source[key]
        }
    }
    }
    return target
}

/**
 * Destroy tocbot.
 */
tocbot.destroy = function () {
    if (!options.skipRendering) {
    // Clear HTML.
    try {
        document.querySelector(options.tocSelector).innerHTML = ''
    } catch (e) {
        console.warn('Element not found: ' + options.tocSelector); // eslint-disable-line
    }
    }

    // Remove event listeners.
    if (options.scrollContainer && document.querySelector(options.scrollContainer)) {
    document.querySelector(options.scrollContainer).removeEventListener('scroll', this._scrollListener, false)
    document.querySelector(options.scrollContainer).removeEventListener('resize', this._scrollListener, false)
    if (buildHtml) {
        document.querySelector(options.scrollContainer).removeEventListener('click', this._clickListener, false)
    }
    } else {
    document.removeEventListener('scroll', this._scrollListener, false)
    document.removeEventListener('resize', this._scrollListener, false)
    if (buildHtml) {
        document.removeEventListener('click', this._clickListener, false)
    }
    }
}

/**
 * Initialize tocbot.
 * @param {object} customOptions
 */
tocbot.init = function (customOptions) {
    // feature test
    if (!supports) {
    return
    }

    // Merge defaults with user options.
    // Set to options variable at the top.
    options = extend(defaultOptions, customOptions || {})
    this.options = options
    this.state = {}
    this.tocEl = document.querySelector(options.tocSelector)

    // Pass options to these modules.
    buildHtml = BuildHtml(options)
    parseContent = ParseContent(options)

    // For testing purposes.
    this._buildHtml = buildHtml
    this._parseContent = parseContent

    // Destroy it if it exists first.
    tocbot.destroy()

    // Get headings array.
    headingsArray = parseContent.selectHeadings(options.contentSelector, options.headingSelector)
    // Return if no headings are found.
    if (headingsArray === null) {
        return
    }

    // Build nested headings array.
    var nestedHeadingsObj = parseContent.nestHeadingsArray(headingsArray)
    var nestedHeadings = nestedHeadingsObj.nest

    // Render.
    if (!options.skipRendering) {
        buildHtml.render(options.tocSelector, nestedHeadings)
    }

    var timeout = false;

    // Update Sidebar and bind listeners.
    this._scrollListener = debounce(() => {
        buildHtml.updateToc(headingsArray, this.tocEl)
        updateTocScroll(options, this.tocEl)
    }, {wait: 200, before: true})

    this._scrollListener()

    if (options.scrollContainer && document.querySelector(options.scrollContainer)) {
        document.querySelector(options.scrollContainer).addEventListener('scroll', this._scrollListener, {passive: true})
        document.querySelector(options.scrollContainer).addEventListener('resize', this._scrollListener, {passive: true})
    } else {
        document.addEventListener('scroll', this._scrollListener, {passive: true})
        document.addEventListener('resize', this._scrollListener, {passive: true})
    }

    return this
}

/**
 * Refresh tocbot.
 */
tocbot.refresh = function (customOptions) {
    tocbot.destroy()
    tocbot.init(customOptions || this.options)
}

module.exports = tocbot