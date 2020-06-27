// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

namespace Visit {
	export async function list() {
		let tbody = Deskop.table(['ID', 'Pseudo', 'Email', 'Admin', 'Application']);
		(await (await fetch('/visit/list')).json())
			.forEach(v => tbody.insertAdjacentHTML('beforeend', `<tr>
				<td>${v.id}</td>
				<td>${v.pseudo}</td>
				<td>${v.email}</td>
				<td>${v.author}</td>
				<td>${v.app}</td>
			</tr>`));
	}
	// Create a new visit ticket.
	export async function create() {
		Deskop.edit(`Création d'un ticket de visite`, list);
		let app = await Edit.createP('Application ID');
		let pseudo = await Edit.createP('Pseudo');
		let email = await Edit.createP('Email');

		Deskop.errorRep(await fetch('/visit/add', {
			method: 'POST',
			headers: new Headers({
				'Content-Type': 'application/x-www-form-urlencoded',
			}),
			body: encodeForm({
				app: app,
				pseudo: pseudo,
				email: email,
			}),
		}));
	}
}

function encodeForm(o): string {
	let body: string[] = [];
	for (let key in o) {
		body.push(`${key}=${encodeURI(o[key])}`)
	}
	return body.join("&");
}
