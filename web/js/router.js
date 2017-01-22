module.exports = {
	firstA,
	getRoute,
	init,
	next,
	meetRequirements
}

function firstA(element) {
	if(element.tagName.toUpperCase() == "A") {
		return element
	}
	if(element.parentElement == null) {
		return null
	}
	return firstA(element.parentElement)
}

function next(a) {
	return getRoute(a).then(text => {
		window.history.pushState(text, window.document.title, a.href)
		return Promise.resolve(text)
	})
}

function getRoute(a) {
	if(a && a.tagName.toUpperCase() != "A") {
		return Promise.reject("Provided value is not an A element")
	}

	return Bliss.fetch(a.href + "?theme=none").then(response => {
		return Promise.resolve(response.responseText)
	})
}

function init() {
	const page = window.location.origin + window.location.pathname
	window.history.replaceState(Bliss("main").innerHTML, window.document.title, page)
}

function meetRequirements() {
	return !!(DOMParser && window.history && window.history.pushState)
}