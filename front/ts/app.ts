// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

namespace App {
	interface App {
		ID: string;
		Name: string;
		URL: string;
	}

	// The list of the apps.
	export var l: Array<App> = [];

	export async function download() {
		l = await (await fetch('/app/list')).json();
	}

	// List all the application
	export async function list() {
		(User.admin ? listAdmin : listSimple)()
	}
	async function listSimple() {
		let table = Deskop.table(['ID', 'Name', 'URL']);
		await download();
		l.forEach(app => {
			let e = $(`<tr>
					<td>${app.ID}</td>
					<td>${app.Name}</td>
					<td>${app.URL}</td>
				</tr>`);
			table.append(e);
		})
	}
	async function listAdmin() {
		let table = Deskop.table(['ID', 'Name', 'URL', '']);
		let l: Array<App> = await (await fetch('/app/list')).json();
		l.forEach(app => {
			let e = $(`<tr>
					<td>${app.ID}</td>
					<td>${app.Name}</td>
					<td>${app.URL}</td>
					<td>
						<button type=button id=edit
							class="btn btn-sm btn-warning">Modifier</button>
						<button type=button id=rm
							class="btn btn-sm btn-danger ml-1">Supprimer</button>
					</td>
				</tr>`);
			e.querySelector('#edit').addEventListener('click', () => edit(app));
			e.querySelector('#rm').addEventListener('click', () => rm(app));
			table.append(e);
		})
	}
	function edit(app: App) {
		Deskop.edit(`ID&nbsp;: ${app.ID}`, list);
		Edit.text(app.Name, 'Nom', `/app/edit/name?id=${app.ID}`);
		Edit.text(app.URL, 'URL', `/app/edit/url?id=${app.ID}`);
	}
	// Supprimer l'application
	async function rm(app: App) {
		Deskop.edit(`Supprétion de l'application&nbsp;: «&#8239;${app.Name}&#8239;»`, list);
		Edit.confirm(app.ID, async () => {
			await fetch(`/app/rm?id=${app.ID}`, { method: 'POST' });
			list();
		});
	}
	// Affiche un écran pour éditer l'application
	export function create() {
		Deskop.edit(`Création d'une application`, list);
		Edit.create('ID', async id => {
			await fetch(`/app/add?id=${id}`, { method: 'POST' });
			edit({
				ID: id,
				Name: '',
				URL: '',
			});
		});
	}
}
