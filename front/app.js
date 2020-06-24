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
                            var e = $("<table><tr>\n\t\t\t\t\t<td>" + app.ID + "</td>\n\t\t\t\t\t<td>" + app.Name + "</td>\n\t\t\t\t\t<td>" + app.URL + "</td>\n\t\t\t\t\t<td>\n\t\t\t\t\t\t<button type=button id=edit\n\t\t\t\t\t\t\tclass=\"btn btn-sm btn-warning\">Modifier</button>\n\t\t\t\t\t\t<button type=button id=rm\n\t\t\t\t\t\t\tclass=\"btn btn-sm btn-danger ml-1\">Supprimer</button>\n\t\t\t\t\t</td>\n\t\t\t\t</tr></table>");
                            e.querySelector('#edit').addEventListener('click', function () { return edit(app); });
                            e.querySelector('#rm').addEventListener('click', function () { return rm(app); });
                            table.append(e.querySelector('tr'));
                        });
                        return [2];
                }
            });
        });
    }
    App.list = list;
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
    function list() {
        reset();
        document.getElementById('list').hidden = false;
    }
    Deskop.list = list;
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
        var groups = ['list', 'edit', 'table'];
        groups
            .map(function (g) { return "#" + g + ">*"; })
            .forEach(function (g) { return document.querySelectorAll(g).forEach(function (e) { return e.remove(); }); });
        groups
            .map(function (g) { return document.getElementById(g); })
            .forEach(function (e) { return e.hidden = true; });
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
})(Deskop || (Deskop = {}));
document.addEventListener("DOMContentLoaded", function () {
    document.getElementById('appGo').addEventListener('click', App.list);
    document.getElementById('logGo').addEventListener('click', Log.list);
}, { once: true });
function $(html) {
    var div = document.createElement('div');
    div.innerHTML = html;
    var e = div.firstElementChild;
    div.remove();
    return e;
}
var Edit;
(function (Edit) {
    function create(name, end) {
        var g = $("<div class=\"input-group input-group-lg mb-3\">\n\t\t\t<div class=\"input-group-prepend\">\n\t\t\t\t<span class=\"input-group-text\">" + name + "&nbsp;: </span>\n\t\t\t</div>\n\t\t\t<input type=text required class=\"form-control\">\n\t\t\t<div class=\"input-group-append\">\n\t\t\t\t<button type=submit class=\"btn btn-success\">Creer</button>\n\t\t\t</div>\n\t\t</div>");
        document.getElementById('edit').appendChild(g);
        g.querySelector('button[type=submit]').addEventListener('click', function () {
            end(g.querySelector('input').value);
        });
    }
    Edit.create = create;
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
        var _this = this;
        var g = $("<div class=\"input-group input-group-lg mb-3\">\n\t\t\t<div class=\"input-group-prepend\">\n\t\t\t\t<span class=\"input-group-text\">" + name + "&nbsp;: </span>\n\t\t\t</div>\n\t\t\t<input type=text required class=\"form-control\">\n\t\t\t<div class=\"input-group-append\">\n\t\t\t\t<button type=submit class=\"btn btn-success\">\n\t\t\t\t\t<span hidden class=\"spinner-border spinner-border-sm mr-2\"></span>\n\t\t\t\t\tEnvoyer\n\t\t\t\t</button>\n\t\t\t</div>\n\t\t</div>");
        document.getElementById('edit').appendChild(g);
        var input = g.querySelector('input');
        input.value = value;
        var spinner = g.querySelector('.spinner-border');
        g.querySelector('button[type=submit]').addEventListener('click', function () { return __awaiter(_this, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        spinner.hidden = false;
                        return [4, fetch(to, {
                                method: 'PATCH',
                                headers: new Headers({
                                    'Content-Type': 'text/plain; charset=utf-8'
                                }),
                                body: input.value
                            })];
                    case 1:
                        _a.sent();
                        spinner.hidden = true;
                        return [2];
                }
            });
        }); });
    }
    Edit.text = text;
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
                        l = _a.sent();
                        l.map(function (i) { return "<tr>\n\t\t\t\t<td>" + i.operation + "</td>\n\t\t\t\t<td>" + i.actor + "</td>\n\t\t\t\t<td>" + new Date(i.date).toLocaleString() + "</td>\n\t\t\t\t<td>" + (i.value || []).map(function (v) { return "'" + v + "'"; }).join('&#8239;; ') + "</td>\n\t\t\t</tr>"; }).forEach(function (i) {
                            tbody.insertAdjacentHTML('beforeend', i);
                        });
                        return [2];
                }
            });
        });
    }
    Log.list = list;
})(Log || (Log = {}));
var User;
(function (User) {
    function list() {
        return __awaiter(this, void 0, void 0, function () {
            var t, l;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        t = Deskop.table(['', 'ID', 'Pseudo', 'Email', 'Level']);
                        return [4, fetch('/user/list')];
                    case 1: return [4, (_a.sent()).json()];
                    case 2:
                        l = _a.sent();
                        console.log("l:", l);
                        l.forEach(function (u) { return t.insertAdjacentHTML('beforeend', "<tr>\n\t\t\t<td></td>\n\t\t\t<td>" + u.login + "</td>\n\t\t\t<td>" + u.pseudo + "</td>\n\t\t\t<td>" + u.email + "</td>\n\t\t\t<td>" + u.level + "</td>\n\t\t</tr>"); });
                        return [2];
                }
            });
        });
    }
    User.list = list;
})(User || (User = {}));
