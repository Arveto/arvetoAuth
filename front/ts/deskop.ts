// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

namespace Deskop {
	// Clear #list to display other elements.
	export function list() {
		reset();
	}
	function reset() {
		document.querySelectorAll("#list>*").forEach(e => e.remove());
	}
	// Display an error
	export function error(m: string) { }
}

document.addEventListener("DOMContentLoaded", () => {
	document.getElementById('appGo').addEventListener('click', App.list, false);
}, { once: true, });
