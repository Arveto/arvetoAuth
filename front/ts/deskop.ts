// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

namespace Deskop {
	// Clear all grousp and display #list.
	export function list() {
		reset();
		document.getElementById('list').hidden = false;
	}
	// Clear all groups and display #edit, append an title and a button
	// to return to a before state.
	export function edit(title: string, endCB: () => void) {
		reset();
		let editGroup = document.getElementById('edit');
		editGroup.hidden = false;
		editGroup.insertAdjacentHTML('beforeend', `<h1>${title}</h1>`);

		let end = $(`<button type=button class="btn btn-danger mb-3">Fin</button>`);
		end.addEventListener('click', endCB);
		editGroup.append(end);
	}
	function reset() {
		const groups = ['list', 'edit'];

		groups
			.map(g => `#${g}>*`)
			.forEach(g => document.querySelectorAll(g).forEach(e => e.remove()))

		groups
			.map(g => document.getElementById(g))
			.forEach(e => e.hidden = true);
	}
	// Display an error
	export function error(m: string) { }
}

document.addEventListener("DOMContentLoaded", () => {
	document.getElementById('appGo').addEventListener('click', App.list, false);
}, { once: true, });

// Create an element and return it.
function $(html: string): Element {
	let div = document.createElement('div');
	div.innerHTML = html;
	let e = div.firstElementChild;
	div.remove();
	return e;
}
