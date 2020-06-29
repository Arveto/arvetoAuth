// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

namespace Log {
	// One event.
	interface Event {
		id: string;
		actor: string;
		operation: string;
		value: string[];
		date: Date;
	}

	// Display log in #table.
	export async function list() {
		let tbody = Deskop.table(['Opération', 'Acteur', 'Date', 'Valeurs']);
		Deskop.activeDate(list);
		let l: Event[] = (await (await fetch('/log/list' + selectDate())).json())
			.map(e => {
				e.date = new Date(e.date);
				return e;
			})
			.sort((e1, e2) => e1.date < e2.date);

		l.forEach(i => tbody.insertAdjacentHTML('beforeend', `<tr>
			<td>${i.operation}</td>
			<td>${i.actor}</td>
			<td>${i.date.toLocaleString()}</td>
			<td>${(i.value || []).map(v => `'${v}'`).join('&#8239;; ')}</td>
		</tr>`));
	}
}
