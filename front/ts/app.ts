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
	export async function edit(app: App) { }
	async function rm(app: App) { }
}

function $(html: string): Element {
	let div = document.createElement('div');
	div.innerHTML = html;
	let e = div.firstElementChild;
	div.remove();
	return e;
}
