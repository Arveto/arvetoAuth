// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

namespace Deskop {
	export var lister: () => void;
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
	// Enable selecting date
	export function activeDate(f: () => void) {
		document.getElementById('dateSelect').hidden = false;
		if (lister !== f) {
			(<HTMLInputElement>document.getElementById('dateYear')).checked = true;
			(<HTMLInputElement>document.getElementById('dateMount')).checked = true;
			(<HTMLInputElement>document.getElementById('dateDay')).checked = false;
			lister = f;
		}
	}
	// Clean the deskop
	function reset() {
		const groups = ['edit', 'table'];

		groups
			.map(g => `#${g}>*`)
			.forEach(g => document.querySelectorAll(g).forEach(e => e.remove()))

		groups
			.map(g => document.getElementById(g))
			.forEach(e => e.hidden = true);

		document.getElementById('create').hidden = true;
		document.getElementById('dateSelect').hidden = true;
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
	document.getElementById('createVisit').addEventListener('click', Visit.create);
	let s: HTMLInputElement = document.querySelector('input[type=search]');
	s.addEventListener('input', () => search(s.value));
	document.getElementById('dateGo').addEventListener('click', () => {
		if (Deskop.lister) {
			Deskop.lister();
		}
	});
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

// Select a date to make a request params
function selectDate(): string {
	let r = '?';

	const d = new Date(
		(<HTMLInputElement>document.getElementById('dateDate')).value
		|| new Date()
	);

	function getUncheck(id: string): boolean {
		return !(<HTMLInputElement>document.getElementById(id))
			.checked;
	}

	if (getUncheck('dateYear')) {
		return r;
	}
	r += 'y=' + d.getFullYear();

	if (getUncheck('dateMount')) {
		return r;
	}
	r += '&m=' + (d.getMonth() + 1);

	if (getUncheck('dateDay')) {
		return r;
	}
	r += '&d=' + d.getDate();

	return r;
}
