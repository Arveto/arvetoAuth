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

		let input = g.querySelector('input');
		function call() {
			if (!input.value) {
				return;
			}
			end(input.value);
		}
		g.querySelector('button[type=submit]').addEventListener('click', call);
		input.addEventListener('keydown', event => {
			if (event.key !== 'Enter') {
				return;
			}
			call();
		});
	}
	// Return a promise resolved with the string.
	export function createP(name: string): Promise<string> {
		let p = new Promise<string>(resolve => create(name, s => {
			document.getElementById('edit')
				.querySelector('div.input-group.input-group-lg.mb-3')
				.remove();
			resolve(s);
		}));
		return p;
	}

	// Create a text confirm. When done, execute the callBack.
	export function confirm(v: string, end: () => void) {
		let g = $(`<div class="input-group input-group-lg mb-3">
			<div class="input-group-prepend">
				<span class="input-group-text">Recopier&nbsp;: </span>
			</div>
			<input type=text required class="form-control">
			<div class="input-group-append">
				<button type=submit class="btn btn-success" disabled>Confirmer</button>
			</div>
		</div>`);
		document.getElementById('edit').append(g);
		let input = g.querySelector('input');
		let go = g.querySelector('button');
		input.placeholder = v;
		input.addEventListener('input', () => go.disabled = input.value !== v);
		go.addEventListener('click', end);
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

		async function send() {
			if (!input.value || input.value === value) {
				return;
			}
			spinner.hidden = false;
			let rep = await fetch(to, {
				method: 'PATCH',
				headers: new Headers({
					'Content-Type': 'text/plain; charset=utf-8'
				}),
				body: input.value,
			});
			spinner.hidden = true;
			Deskop.errorRep(rep);
		}

		input.addEventListener('keydown', event => {
			if (event.key !== 'Enter') {
				return;
			}
			send();
			nextFocus(input);
		});
		g.querySelector('button[type=submit]').addEventListener('click', send);
	}
	// Focus on the next input
	function nextFocus(input: HTMLInputElement) {
		let list = <NodeListOf<HTMLInputElement>>document.querySelectorAll('input[type=text]');
		for (let i = 0; i < list.length; i++) {
			if (list[i] === input && i + 1 < list.length) {
				list[i + 1].focus();
			}
		}
	}
	// Create an option list.
	export function options(values: string[], current: string, name: string, end: (string) => void) {
		let opt = $(`<div class="form-group"></div>`);
		document.getElementById('edit').append(opt);
		values.forEach(v => {
			let g = $(`<div class="custom-control custom-radio">
				<input type=radio class="custom-control-input" id="${btoa(v)}" name="${btoa(name)}">
				<label class="custom-control-label" for="${btoa(v)}">${v}</label>
			</div>`);
			opt.append(g);
			let input = g.querySelector('input');
			input.checked = current === v;
			input.addEventListener('change', () => {
				if (input.checked) {
					end(v);
				}
			});
		});
	}
	export function avatar() {
		let g = $(`<div class="input-group mb-3">
			<input type=file id=avatarInput hidden accept="image/png,image/jpeg,image/gif,image/webp">
			<label for=avatarInput>
				Envoyer une image (512×512)&nbsp;:<br>
				<img class="rounded" src="${User.me.avatar}">
			</label>
		</div>`);
		document.getElementById('edit').append(g);
		let label = g.querySelector('label');
		let input = g.querySelector('input');
		async function send(f) {
			let rep = await fetch('/avatar/edit', {
				method: 'PATCH',
				headers: new Headers({
					'Content-Type': f.type,
				}),
				body: await f.arrayBuffer(),
			});
			if (rep.status === 200) {
				document.location.reload();
			}
			Deskop.errorRep(rep);
		}
		input.addEventListener('input', () => send(input.files[0]));
		label.addEventListener('drop', event => {
			send(event.dataTransfer.files[0]);
		}, false);

		label.addEventListener('dragenter', event => {
			label.classList.add('border', 'border-primary')
		}, false);
		label.addEventListener('dragleave', event => {
			label.classList.remove('border', 'border-primary')
		}, false);

		['dragover', 'drop'].forEach(name => {
			label.addEventListener(name, event => {
				event.stopPropagation();
				event.preventDefault();
			}, false);
		});
	}
}
