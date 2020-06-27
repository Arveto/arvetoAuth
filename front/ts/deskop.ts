// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

namespace Deskop {
	// Clear all grousp and display #list.
	export function create() {
		reset();
		document.getElementById('create').hidden = false;
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
		t.insertAdjacentHTML('beforeend', '<thead class="thead-light"><tr>' +
			head.map(h => `<th scope="col">${h}</th>`).join('') +
			'</tr></thead><tbody></tbody>');
		return t.querySelector('tbody');
	}
	function reset() {
		const groups = ['edit', 'table'];

		groups
			.map(g => `#${g}>*`)
			.forEach(g => document.querySelectorAll(g).forEach(e => e.remove()))

		groups
			.map(g => document.getElementById(g))
			.forEach(e => e.hidden = true);

		document.getElementById('create').hidden = true;
	}
	// Display an error
	export function error(m: string) {
		if (!m) { return; }
		let e = $(`<div class="container-sm alert alert-dismissible alert-danger"><button type=button class="close" data-dismiss="alert">&times;</button>${m}</div>`);
		e.querySelector('button.close').addEventListener('click', () => e.remove());
		document.querySelector('nav').insertAdjacentElement('afterend', e);
	}
	// Display the error if the request fail.
	export async function errorRep(rep: Response) {
		if (rep.status !== 200) {
			error(await rep.text());
		}
	}
}

document.addEventListener("DOMContentLoaded", () => {
	User.list();
	document.getElementById('myconfig').addEventListener('click', User.editMe);
	document.getElementById('userGo').addEventListener('click', User.list);
	document.getElementById('logGo').addEventListener('click', Log.list);
	document.getElementById('visitGo').addEventListener('click', Visit.list);
	document.getElementById('appGo').addEventListener('click', App.list);
	document.getElementById('createGo').addEventListener('click', Deskop.create);
	document.getElementById('createApplication').addEventListener('click', App.create);
	let s: HTMLInputElement = document.querySelector('input[type=search]');
	s.addEventListener('input', () => search(s.value));
}, { once: true, });

// Make a search into #table.
function search(pattern: string) {
	pattern = pattern.toLowerCase();
	Array.from(document.querySelectorAll('tr'))
		.filter(tr => tr.parentElement.tagName !== 'THEAD')
		.forEach(tr => {
			tr.hidden = !Array.from(tr.querySelectorAll('td'))
				.some(td => td.innerText.toLowerCase().includes(pattern));
		});
}

// Create an element and return it.
function $(html: string): Element {
	let table = /^\s*<tr/.test(html);
	let div = document.createElement(table ? 'table' : 'div');
	div.innerHTML = html;
	let e = table ? div.querySelector('tr') : div.firstElementChild;
	div.remove();
	return e;
}
