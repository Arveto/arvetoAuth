var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
var App;
(function (App) {
    function list() {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                (User.admin ? listAdmin : listSimple)();
                return [2];
            });
        });
    }
    App.list = list;
    function listSimple() {
        return __awaiter(this, void 0, void 0, function () {
            var table, l;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        table = Deskop.table(['ID', 'Name', 'URL']);
                        return [4, fetch('/app/list')];
                    case 1: return [4, (_a.sent()).json()];
                    case 2:
                        l = _a.sent();
                        l.forEach(function (app) {
                            var e = $("<tr>\n\t\t\t\t\t<td>" + app.ID + "</td>\n\t\t\t\t\t<td>" + app.Name + "</td>\n\t\t\t\t\t<td>" + app.URL + "</td>\n\t\t\t\t</tr>");
                            table.append(e);
                        });
                        return [2];
                }
            });
        });
    }
    function listAdmin() {
        return __awaiter(this, void 0, void 0, function () {
            var table, l;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        table = Deskop.table(['ID', 'Name', 'URL', '']);
                        return [4, fetch('/app/list')];
                    case 1: return [4, (_a.sent()).json()];
                    case 2:
                        l = _a.sent();
                        l.forEach(function (app) {
                            var e = $("<tr>\n\t\t\t\t\t<td>" + app.ID + "</td>\n\t\t\t\t\t<td>" + app.Name + "</td>\n\t\t\t\t\t<td>" + app.URL + "</td>\n\t\t\t\t\t<td>\n\t\t\t\t\t\t<button type=button id=edit\n\t\t\t\t\t\t\tclass=\"btn btn-sm btn-warning\">Modifier</button>\n\t\t\t\t\t\t<button type=button id=rm\n\t\t\t\t\t\t\tclass=\"btn btn-sm btn-danger ml-1\">Supprimer</button>\n\t\t\t\t\t</td>\n\t\t\t\t</tr>");
                            e.querySelector('#edit').addEventListener('click', function () { return edit(app); });
                            e.querySelector('#rm').addEventListener('click', function () { return rm(app); });
                            table.append(e);
                        });
                        return [2];
                }
            });
        });
    }
    function edit(app) {
        Deskop.edit("ID&nbsp;: " + app.ID, list);
        Edit.text(app.Name, 'Nom', "/app/edit/name?id=" + app.ID);
        Edit.text(app.URL, 'URL', "/app/edit/url?id=" + app.ID);
    }
    function rm(app) {
        return __awaiter(this, void 0, void 0, function () {
            var _this = this;
            return __generator(this, function (_a) {
                Deskop.edit("Suppr\u00E9tion de l'application&nbsp;: \u00AB&#8239;" + app.Name + "&#8239;\u00BB", list);
                Edit.confirm(app.ID, function () { return __awaiter(_this, void 0, void 0, function () {
                    return __generator(this, function (_a) {
                        switch (_a.label) {
                            case 0: return [4, fetch("/app/rm?id=" + app.ID, { method: 'POST' })];
                            case 1:
                                _a.sent();
                                list();
                                return [2];
                        }
                    });
                }); });
                return [2];
            });
        });
    }
    function create() {
        var _this = this;
        Deskop.edit("Cr\u00E9ation d'une application", list);
        Edit.create('ID', function (id) { return __awaiter(_this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4, fetch("/app/add?id=" + id, { method: 'POST' })];
                    case 1:
                        _a.sent();
                        edit({
                            ID: id,
                            Name: '',
                            URL: ''
                        });
                        return [2];
                }
            });
        }); });
    }
    App.create = create;
})(App || (App = {}));
var Deskop;
(function (Deskop) {
    function create() {
        reset();
        document.getElementById('create').hidden = false;
    }
    Deskop.create = create;
    function edit(title, endCB) {
        reset();
        var editGroup = document.getElementById('edit');
        editGroup.hidden = false;
        editGroup.insertAdjacentHTML('beforeend', "<h1>" + title + "</h1>");
        var end = $("<button type=button class=\"btn btn-danger mb-3\">Fin</button>");
        end.addEventListener('click', endCB);
        editGroup.append(end);
    }
    Deskop.edit = edit;
    function table(head) {
        reset();
        var t = document.getElementById('table');
        t.hidden = false;
        t.insertAdjacentHTML('beforeend', '<thead class="thead-light"><tr>' +
            head.map(function (h) { return "<th scope=\"col\">" + h + "</th>"; }).join('') +
            '</tr></thead><tbody></tbody>');
        return t.querySelector('tbody');
    }
    Deskop.table = table;
    function reset() {
        var groups = ['edit', 'table'];
        groups
            .map(function (g) { return "#" + g + ">*"; })
            .forEach(function (g) { return document.querySelectorAll(g).forEach(function (e) { return e.remove(); }); });
        groups
            .map(function (g) { return document.getElementById(g); })
            .forEach(function (e) { return e.hidden = true; });
        document.getElementById('create').hidden = true;
    }
    function error(m) {
        if (!m) {
            return;
        }
        var e = $("<div class=\"container-sm alert alert-dismissible alert-danger\"><button type=button class=\"close\" data-dismiss=\"alert\">&times;</button>" + m + "</div>");
        e.querySelector('button.close').addEventListener('click', function () { return e.remove(); });
        document.querySelector('nav').insertAdjacentElement('afterend', e);
    }
    Deskop.error = error;
    function errorRep(rep) {
        return __awaiter(this, void 0, void 0, function () {
            var _a;
            return __generator(this, function (_b) {
                switch (_b.label) {
                    case 0:
                        if (!(rep.status !== 200)) return [3, 2];
                        _a = error;
                        return [4, rep.text()];
                    case 1:
                        _a.apply(void 0, [_b.sent()]);
                        _b.label = 2;
                    case 2: return [2];
                }
            });
        });
    }
    Deskop.errorRep = errorRep;
})(Deskop || (Deskop = {}));
document.addEventListener("DOMContentLoaded", function () {
    User.list();
    document.getElementById('myconfig').addEventListener('click', User.editMe);
    document.getElementById('userGo').addEventListener('click', User.list);
    document.getElementById('logGo').addEventListener('click', Log.list);
    document.getElementById('visitGo').addEventListener('click', Visit.list);
    document.getElementById('appGo').addEventListener('click', App.list);
    document.getElementById('createGo').addEventListener('click', Deskop.create);
    document.getElementById('createApplication').addEventListener('click', App.create);
    document.getElementById('createVisit').addEventListener('click', Visit.create);
    var s = document.querySelector('input[type=search]');
    s.addEventListener('input', function () { return search(s.value); });
}, { once: true });
function search(pattern) {
    pattern = pattern.toLowerCase();
    Array.from(document.querySelectorAll('tr'))
        .filter(function (tr) { return tr.parentElement.tagName !== 'THEAD'; })
        .forEach(function (tr) {
        tr.hidden = !Array.from(tr.querySelectorAll('td'))
            .some(function (td) { return td.innerText.toLowerCase().includes(pattern); });
    });
}
function $(html) {
    var table = /^\s*<tr/.test(html);
    var div = document.createElement(table ? 'table' : 'div');
    div.innerHTML = html;
    var e = table ? div.querySelector('tr') : div.firstElementChild;
    div.remove();
    return e;
}
var Edit;
(function (Edit) {
    function create(name, end) {
        var g = $("<div class=\"input-group input-group-lg mb-3\">\n\t\t\t<div class=\"input-group-prepend\">\n\t\t\t\t<span class=\"input-group-text\">" + name + "&nbsp;: </span>\n\t\t\t</div>\n\t\t\t<input type=text required class=\"form-control\">\n\t\t\t<div class=\"input-group-append\">\n\t\t\t\t<button type=submit class=\"btn btn-success\">Creer</button>\n\t\t\t</div>\n\t\t</div>");
        document.getElementById('edit').appendChild(g);
        var input = g.querySelector('input');
        function call() {
            if (!input.value) {
                return;
            }
            end(input.value);
        }
        g.querySelector('button[type=submit]').addEventListener('click', call);
        input.addEventListener('keydown', function (event) {
            if (event.key !== 'Enter') {
                return;
            }
            call();
        });
    }
    Edit.create = create;
    function createP(name) {
        var p = new Promise(function (resolve) { return create(name, function (s) {
            document.getElementById('edit')
                .querySelector('div.input-group.input-group-lg.mb-3')
                .remove();
            resolve(s);
        }); });
        return p;
    }
    Edit.createP = createP;
    function confirm(v, end) {
        var g = $("<div class=\"input-group input-group-lg mb-3\">\n\t\t\t<div class=\"input-group-prepend\">\n\t\t\t\t<span class=\"input-group-text\">Recopier&nbsp;: </span>\n\t\t\t</div>\n\t\t\t<input type=text required class=\"form-control\">\n\t\t\t<div class=\"input-group-append\">\n\t\t\t\t<button type=submit class=\"btn btn-success\" disabled>Confirmer</button>\n\t\t\t</div>\n\t\t</div>");
        document.getElementById('edit').append(g);
        var input = g.querySelector('input');
        var go = g.querySelector('button');
        input.placeholder = v;
        input.addEventListener('input', function () { return go.disabled = input.value !== v; });
        go.addEventListener('click', end);
    }
    Edit.confirm = confirm;
    function text(value, name, to) {
        var g = $("<div class=\"input-group input-group-lg mb-3\">\n\t\t\t<div class=\"input-group-prepend\">\n\t\t\t\t<span class=\"input-group-text\">" + name + "&nbsp;: </span>\n\t\t\t</div>\n\t\t\t<input type=text required class=\"form-control\">\n\t\t\t<div class=\"input-group-append\">\n\t\t\t\t<button type=submit class=\"btn btn-success\">\n\t\t\t\t\t<span hidden class=\"spinner-border spinner-border-sm mr-2\"></span>\n\t\t\t\t\tEnvoyer\n\t\t\t\t</button>\n\t\t\t</div>\n\t\t</div>");
        document.getElementById('edit').appendChild(g);
        var input = g.querySelector('input');
        input.value = value;
        var spinner = g.querySelector('.spinner-border');
        function send() {
            return __awaiter(this, void 0, void 0, function () {
                var rep;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0:
                            if (!input.value || input.value === value) {
                                return [2];
                            }
                            spinner.hidden = false;
                            return [4, fetch(to, {
                                    method: 'PATCH',
                                    headers: new Headers({
                                        'Content-Type': 'text/plain; charset=utf-8'
                                    }),
                                    body: input.value
                                })];
                        case 1:
                            rep = _a.sent();
                            spinner.hidden = true;
                            Deskop.errorRep(rep);
                            return [2];
                    }
                });
            });
        }
        input.addEventListener('keydown', function (event) {
            if (event.key !== 'Enter') {
                return;
            }
            send();
            nextFocus(input);
        });
        g.querySelector('button[type=submit]').addEventListener('click', send);
    }
    Edit.text = text;
    function nextFocus(input) {
        var list = document.querySelectorAll('input[type=text]');
        for (var i = 0; i < list.length; i++) {
            if (list[i] === input && i + 1 < list.length) {
                list[i + 1].focus();
            }
        }
    }
    function options(values, current, name, end) {
        var opt = $("<div class=\"form-group\"></div>");
        document.getElementById('edit').append(opt);
        values.forEach(function (v) {
            var g = $("<div class=\"custom-control custom-radio\">\n\t\t\t\t<input type=radio class=\"custom-control-input\" id=\"" + btoa(v) + "\" name=\"" + btoa(name) + "\">\n\t\t\t\t<label class=\"custom-control-label\" for=\"" + btoa(v) + "\">" + v + "</label>\n\t\t\t</div>");
            opt.append(g);
            var input = g.querySelector('input');
            input.checked = current === v;
            input.addEventListener('change', function () {
                if (input.checked) {
                    end(v);
                }
            });
        });
    }
    Edit.options = options;
    function avatar() {
        var g = $("<div class=\"input-group mb-3\">\n\t\t\t<input type=file id=avatarInput hidden accept=\"image/png,image/jpeg,image/gif,image/webp\">\n\t\t\t<label for=avatarInput>\n\t\t\t\tEnvoyer une image (512\u00D7512)&nbsp;:<br>\n\t\t\t\t<img class=\"rounded\" src=\"" + User.me.avatar + "\">\n\t\t\t</label>\n\t\t</div>");
        document.getElementById('edit').append(g);
        var label = g.querySelector('label');
        var input = g.querySelector('input');
        function send(f) {
            return __awaiter(this, void 0, void 0, function () {
                var rep, _a, _b, _c;
                return __generator(this, function (_d) {
                    switch (_d.label) {
                        case 0:
                            _a = fetch;
                            _b = ['/avatar/edit'];
                            _c = {
                                method: 'PATCH',
                                headers: new Headers({
                                    'Content-Type': f.type
                                })
                            };
                            return [4, f.arrayBuffer()];
                        case 1: return [4, _a.apply(void 0, _b.concat([(_c.body = _d.sent(),
                                    _c)]))];
                        case 2:
                            rep = _d.sent();
                            if (rep.status === 200) {
                                document.location.reload();
                            }
                            Deskop.errorRep(rep);
                            return [2];
                    }
                });
            });
        }
        input.addEventListener('input', function () { return send(input.files[0]); });
        label.addEventListener('drop', function (event) {
            send(event.dataTransfer.files[0]);
        }, false);
        label.addEventListener('dragenter', function (event) {
            label.classList.add('border', 'border-primary');
        }, false);
        label.addEventListener('dragleave', function (event) {
            label.classList.remove('border', 'border-primary');
        }, false);
        ['dragover', 'drop'].forEach(function (name) {
            label.addEventListener(name, function (event) {
                event.stopPropagation();
                event.preventDefault();
            }, false);
        });
    }
    Edit.avatar = avatar;
})(Edit || (Edit = {}));
var Log;
(function (Log) {
    function list() {
        return __awaiter(this, void 0, void 0, function () {
            var tbody, l;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        tbody = Deskop.table(['Opération', 'Acteur', 'Date', 'Valeurs']);
                        return [4, fetch('/log/list')];
                    case 1: return [4, (_a.sent()).json()];
                    case 2:
                        l = (_a.sent())
                            .map(function (e) {
                            e.date = new Date(e.date);
                            return e;
                        })
                            .sort(function (e1, e2) { return e1.date < e2.date; });
                        l.forEach(function (i) { return tbody.insertAdjacentHTML('beforeend', "<tr>\n\t\t\t<td>" + i.operation + "</td>\n\t\t\t<td>" + i.actor + "</td>\n\t\t\t<td>" + i.date.toLocaleString() + "</td>\n\t\t\t<td>" + (i.value || []).map(function (v) { return "'" + v + "'"; }).join('&#8239;; ') + "</td>\n\t\t</tr>"); });
                        return [2];
                }
            });
        });
    }
    Log.list = list;
})(Log || (Log = {}));
var User;
(function (User) {
    (function () {
        return __awaiter(this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4, fetch('/me')];
                    case 1: return [4, (_a.sent()).json()];
                    case 2:
                        User.me = _a.sent();
                        User.admin = User.me.level === 'Admin';
                        document.getElementById('createGo').hidden = User.me.level !== 'Admin';
                        loadMyAvatar();
                        return [2];
                }
            });
        });
    })();
    function loadMyAvatar() {
        document.querySelector('#myconfig>img')
            .src = "/avatar/get?u=" + User.me.id;
    }
    function list() {
        return __awaiter(this, void 0, void 0, function () {
            var l;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4, fetch('/user/list')];
                    case 1: return [4, (_a.sent()).json()];
                    case 2:
                        l = (_a.sent())
                            .filter(function (u) { return u.id != User.me.id; });
                        (User.admin ? listAdmin : listSimple)(l);
                        return [2];
                }
            });
        });
    }
    User.list = list;
    function listSimple(l) {
        return __awaiter(this, void 0, void 0, function () {
            var t;
            return __generator(this, function (_a) {
                t = Deskop.table(['', 'ID', 'Pseudo', 'Email', 'Level']);
                l.forEach(function (u) {
                    t.insertAdjacentHTML('beforeend', "<tr>\n\t\t\t\t<td class=\"align-middle\"><img src=\"" + u.avatar + "\"\n\t\t\t\t\tclass=\"rounded\" style=\"width:3em;\"></td>\n\t\t\t\t<td class=\"align-middle\">" + u.id + "</td>\n\t\t\t\t<td class=\"align-middle\">" + u.pseudo + "</td>\n\t\t\t\t<td class=\"align-middle\">" + u.email + "</td>\n\t\t\t\t<td class=\"align-middle\">" + u.level + "</td>\n\t\t\t</tr>");
                });
                return [2];
            });
        });
    }
    function listAdmin(l) {
        return __awaiter(this, void 0, void 0, function () {
            var t;
            return __generator(this, function (_a) {
                t = Deskop.table(['', 'ID', 'Pseudo', 'Email', 'Level', '']);
                l.forEach(function (u) {
                    var tr = $("<tr>\n\t\t\t\t<td class=\"align-middle\"><img src=\"" + u.avatar + "\"\n\t\t\t\t\tclass=\"rounded\" style=\"width:3em;\"></td>\n\t\t\t\t<td class=\"align-middle\">" + u.id + "</td>\n\t\t\t\t<td class=\"align-middle\">" + u.pseudo + "</td>\n\t\t\t\t<td class=\"align-middle\">" + u.email + "</td>\n\t\t\t\t<td class=\"align-middle\">" + u.level + "</td>\n\t\t\t\t<td class=\"align-middle\">\n\t\t\t\t\t<button type=button id=level\n\t\t\t\t\t\tclass=\"btn btn-sm btn-warning\">Modifier</button>\n\t\t\t\t\t<button type=button id=rm\n\t\t\t\t\t\tclass=\"btn btn-sm btn-danger ml-1\">Supprimer</button>\n\t\t\t\t</td>\n\t\t\t</tr>");
                    tr.querySelector('#level').addEventListener('click', function () { return editLevel(u); });
                    tr.querySelector('#rm').addEventListener('click', function () { return rm(u); });
                    t.append(tr);
                });
                return [2];
            });
        });
    }
    function editLevel(u) {
        var _this = this;
        Deskop.edit("Modification de l'accr\u00E9ditation de " + u.pseudo, list);
        Edit.options(['Ban', 'Candidate', 'Visitor', 'Std', 'Admin'], u.level, 'name', function (l) { return __awaiter(_this, void 0, void 0, function () {
            var _a, _b;
            return __generator(this, function (_c) {
                switch (_c.label) {
                    case 0:
                        _b = (_a = Deskop).errorRep;
                        return [4, fetch("/user/edit/level?u=" + u.id, {
                                method: 'PATCH',
                                headers: new Headers({ 'Content-Type': 'application/json' }),
                                body: JSON.stringify(l)
                            })];
                    case 1: return [4, _b.apply(_a, [_c.sent()])];
                    case 2:
                        _c.sent();
                        list();
                        return [2];
                }
            });
        }); });
    }
    function rm(u) {
        var _this = this;
        Deskop.edit("Modification de l'accr\u00E9ditation de " + u.pseudo, list);
        Edit.confirm(u.id, function () { return __awaiter(_this, void 0, void 0, function () {
            var _a, _b;
            return __generator(this, function (_c) {
                switch (_c.label) {
                    case 0:
                        _b = (_a = Deskop).errorRep;
                        return [4, fetch("/user/rm/other?u=" + u.id, {
                                method: 'PATCH'
                            })];
                    case 1: return [4, _b.apply(_a, [_c.sent()])];
                    case 2:
                        _c.sent();
                        list();
                        return [2];
                }
            });
        }); });
    }
    function editMe() {
        var _this = this;
        Deskop.edit("Modification de son compte", list);
        Edit.text(User.me.pseudo, 'Pseudo', '/user/edit/pseudo');
        Edit.text(User.me.email, 'Email', '/user/edit/email');
        Edit.avatar();
        var rmButton = $('<button type=button class="btn btn-danger btn-lg">Supprime mon compte</button>');
        rmButton.addEventListener('click', function () {
            Deskop.edit('Supprime mon compte', editMe);
            Edit.confirm(User.me.id, function () { return __awaiter(_this, void 0, void 0, function () {
                var rep;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [4, fetch('/user/rm/me')];
                        case 1:
                            rep = _a.sent();
                            if (rep.status !== 200) {
                                Deskop.errorRep(rep);
                                return [2];
                            }
                            document.location.replace('/');
                            return [2];
                    }
                });
            }); });
        });
        document.getElementById('edit').append(rmButton);
    }
    User.editMe = editMe;
})(User || (User = {}));
var Visit;
(function (Visit) {
    function list() {
        return __awaiter(this, void 0, void 0, function () {
            var tbody;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        tbody = Deskop.table(['ID', 'Pseudo', 'Email', 'Admin', 'Application']);
                        return [4, fetch('/visit/list')];
                    case 1: return [4, (_a.sent()).json()];
                    case 2:
                        (_a.sent())
                            .forEach(function (v) { return tbody.insertAdjacentHTML('beforeend', "<tr>\n\t\t\t\t<td>" + v.id + "</td>\n\t\t\t\t<td>" + v.pseudo + "</td>\n\t\t\t\t<td>" + v.email + "</td>\n\t\t\t\t<td>" + v.author + "</td>\n\t\t\t\t<td>" + v.app + "</td>\n\t\t\t</tr>"); });
                        return [2];
                }
            });
        });
    }
    Visit.list = list;
    function create() {
        return __awaiter(this, void 0, void 0, function () {
            var app, pseudo, email, _a, _b;
            return __generator(this, function (_c) {
                switch (_c.label) {
                    case 0:
                        Deskop.edit("Cr\u00E9ation d'un ticket de visite", list);
                        return [4, Edit.createP('Application ID')];
                    case 1:
                        app = _c.sent();
                        return [4, Edit.createP('Pseudo')];
                    case 2:
                        pseudo = _c.sent();
                        return [4, Edit.createP('Email')];
                    case 3:
                        email = _c.sent();
                        _b = (_a = Deskop).errorRep;
                        return [4, fetch('/visit/add', {
                                method: 'POST',
                                headers: new Headers({
                                    'Content-Type': 'application/x-www-form-urlencoded'
                                }),
                                body: encodeForm({
                                    app: app,
                                    pseudo: pseudo,
                                    email: email
                                })
                            })];
                    case 4:
                        _b.apply(_a, [_c.sent()]);
                        return [2];
                }
            });
        });
    }
    Visit.create = create;
})(Visit || (Visit = {}));
function encodeForm(o) {
    var body = [];
    for (var key in o) {
        body.push(key + "=" + encodeURI(o[key]));
    }
    return body.join("&");
}
