// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

namespace App {
	interface App {
		ID: string;
		Name: string;
		URL: string;
	}

	// List all the application
	export async function list() {
		Deskop.list();
		let listGroup = document.getElementById('list');
		let l: Array<App> = await (await fetch('/app/list')).json();
		for (let app of l) {
			let e = $(`<div class="card card-body">
				<h3 class="card-title">${app.Name}</h3>
				<div class="my-1">
					ID&nbsp;: <mark class="text-monospace">${app.ID}</mark>
				</div>
				<div class="my-1">
					URL&nbsp;: <mark class="text-monospace">${app.URL}</mark>
				</div>
				<div class="mt-1">
					<button type=button id=edit
							class="btn btn-sm btn-warning">Modifier</button>
					<button type=button id=rm
							class="btn btn-sm btn-danger ml-1">Supprimer</button>
				</div>
			</div>`);
			e.querySelector('#edit').addEventListener('click', () => edit(app));
			e.querySelector('#rm').addEventListener('click', () => rm(app));
			listGroup.append(e);
		}
	}
	function edit(app: App) {
		Deskop.edit(`ID&nbsp;: ${app.ID}`, list);
		Edit.text(app.Name, 'Nom', `/app/edit/name?id=${app.ID}`);
		Edit.text(app.URL, 'URL', `/app/edit/url?id=${app.ID}`);
	}
	async function rm(app: App) { }

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
