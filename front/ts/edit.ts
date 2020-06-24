// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

namespace Edit {
	// Ajoute un champ pour créer un élément, une fois validé, le call-back
	// end est appelé.
	export function create(name: string, end: (string) => void) {
		let g = $(`<div class="input-group input-group-lg mb-3">
			<div class="input-group-prepend">
				<span class="input-group-text">${name}&nbsp;: </span>
			</div>
			<input type=text required class="form-control">
			<div class="input-group-append">
				<button type=submit class="btn btn-success">Creer</button>
			</div>
		</div>`);
		document.getElementById('edit').appendChild(g);
		g.querySelector('button[type=submit]').addEventListener('click', () => {
			end(g.querySelector('input').value);
		});
	}
	// Add and text input.
	export function text(value: string, name: string, to: string) {
		let g = $(`<div class="input-group input-group-lg mb-3">
			<div class="input-group-prepend">
				<span class="input-group-text">${name}&nbsp;: </span>
			</div>
			<input type=text required class="form-control">
			<div class="input-group-append">
				<button type=submit class="btn btn-success">
					<span hidden class="spinner-border spinner-border-sm mr-2"></span>
					Envoyer
				</button>
			</div>
		</div>`);
		document.getElementById('edit').appendChild(g);

		let input = g.querySelector('input');
		input.value = value;
		let spinner: HTMLElement = g.querySelector('.spinner-border');

		g.querySelector('button[type=submit]').addEventListener('click', async () => {
			spinner.hidden = false;
			await fetch(to, {
				method: 'PATCH',
				headers: new Headers({
					'Content-Type': 'text/plain; charset=utf-8'
				}),
				body: input.value,
			});
			spinner.hidden = true;
		});
	}
}
