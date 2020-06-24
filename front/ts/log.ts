// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

namespace Log {
	// Display log in #table.
	export async function list() {
		let tbody = Deskop.table(['Opération', 'Acteur', 'Date', 'Valeurs']);
		let l = await (await fetch('/log/list')).json();
		l.map(i => `<tr>
				<td>${i.operation}</td>
				<td>${i.actor}</td>
				<td>${new Date(i.date).toLocaleString()}</td>
				<td>${(i.value || []).map(v => `'${v}'`).join('&#8239;; ')}</td>
			</tr>`).forEach(i => {
				tbody.insertAdjacentHTML('beforeend', i)
			});
	}
}
