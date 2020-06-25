// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

namespace User {
	interface User {
		id: string;
		pseudo: string;
		email: string;
		avatar: string;
		level: string;
	}

	export var me: User;
	export var admin: boolean;
	(async function() {
		me = await (await fetch('/me')).json();
		admin = me.level === 'Admin';
		loadMyAvatar();
	})();

	function loadMyAvatar() {
		(<HTMLImageElement>document.querySelector('#myconfig>img'))
			.src = `/avatar/get?u=${me.id}`;
	}

	// List the user in #table
	export async function list() {
		let l: User[] = (await (await fetch('/user/list')).json())
			.filter(u => u.id != me.id);
		(admin ? listAdmin : listSimple)(l)
	}
	async function listSimple(l: User[]) {
		let t = Deskop.table(['', 'ID', 'Pseudo', 'Email', 'Level']);
		l.forEach(u => {
			t.insertAdjacentHTML('beforeend', `<tr>
				<td class="align-middle"><img src="${u.avatar}"
					class="rounded" style="width:3em;"></td>
				<td class="align-middle">${u.id}</td>
				<td class="align-middle">${u.pseudo}</td>
				<td class="align-middle">${u.email}</td>
				<td class="align-middle">${u.level}</td>
			</tr>`);
		});
	}
	async function listAdmin(l: User[]) {
		let t = Deskop.table(['', 'ID', 'Pseudo', 'Email', 'Level', '']);
		l.forEach(u => {
			let tr = $(`<tr>
				<td class="align-middle"><img src="${u.avatar}"
					class="rounded" style="width:3em;"></td>
				<td class="align-middle">${u.id}</td>
				<td class="align-middle">${u.pseudo}</td>
				<td class="align-middle">${u.email}</td>
				<td class="align-middle">${u.level}</td>
				<td class="align-middle">
					<button type=button id=level
						class="btn btn-sm btn-warning">Modifier</button>
					<button type=button id=rm
						class="btn btn-sm btn-danger ml-1">Supprimer</button>
				</td>
			</tr>`);
			tr.querySelector('#level').addEventListener('click', () => editLevel(u));
			tr.querySelector('#rm').addEventListener('click', () => rm(u));
			t.append(tr);
		});
	}
	function editLevel(u: User) {
		Deskop.edit(`Modification de l'accréditation de ${u.pseudo}`, list);
		Edit.options(['Ban', 'Candidate', 'Visitor', 'Std', 'Admin'],
			u.level, 'name', async l => {
				await Deskop.errorRep(await fetch(`/user/edit/level?u=${u.id}`, {
					method: 'PATCH',
					headers: new Headers({ 'Content-Type': 'application/json' }),
					body: JSON.stringify(l),
				}));
				list();
			});
	}
	function rm(u: User) {
		Deskop.edit(`Modification de l'accréditation de ${u.pseudo}`, list);
		Edit.confirm(u.id, async () => {
			await Deskop.errorRep(await fetch(`/user/rm/other?u=${u.id}`, {
				method: 'PATCH',
			}));
			list();
		});
	}
	// Edit her own configuration
	export function editMe() {
		Deskop.edit(`Modification de son compte`, list);
		Edit.text(me.pseudo, 'Pseudo', '/user/edit/pseudo');
		Edit.text(me.email, 'Email', '/user/edit/email');
		let rmButton = $('<button type=button class="btn btn-danger btn-lg">Supprime mon compte</button>');
		rmButton.addEventListener('click', () => {
			Deskop.edit('Supprime mon compte', editMe);
			Edit.confirm(me.id, async () => {
				let rep = await fetch('/user/rm/me');
				if (rep.status !== 200) {
					Deskop.errorRep(rep);
					return;
				}
				document.location.replace('/');
			});
		});
		document.getElementById('edit').append(rmButton)
	}
}
