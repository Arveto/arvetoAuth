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
	// Clear all groups, display table with the headers and return <tbody>.
	export function table(head: string[]): HTMLTableSectionElement {
		reset();
		let t = document.getElementById('table');
		t.hidden = false;
		t.insertAdjacentHTML('beforeend', '<thead><tr>' +
			head.map(h => `<th scope="col">${h}</th>`).join('') +
			'</tr></thead><tbody></tbody>');
		return t.querySelector('tbody');
	}
	function reset() {
		const groups = ['list', 'edit', 'table'];

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
	document.getElementById('appGo').addEventListener('click', App.list);
	document.getElementById('logGo').addEventListener('click', Log.list);
}, { once: true, });

// Create an element and return it.
function $(html: string): Element {
	let div = document.createElement('div');
	div.innerHTML = html;
	let e = div.firstElementChild;
	div.remove();
	return e;
}
