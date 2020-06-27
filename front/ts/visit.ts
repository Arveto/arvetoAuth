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
}
