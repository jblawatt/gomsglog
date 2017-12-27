"use strict";

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

var _itemStorage = [];

var FormComp = {
    onSubmit: function onSubmit(evt) {
        evt.preventDefault();
        var input = evt.currentTarget.querySelector('.message');
        input.setAttribute("disabled", true);
        var message = input.value;
        var method = evt.currentTarget.getAttribute("method");
        var action = evt.currentTarget.getAttribute("action");
        m.request({
            url: action,
            method: method,
            headers: {
                "Content-Type": "application/json"
            },
            data: {
                message: message
            }
        }).then(function (r) {
            input.value = "";
            input.removeAttribute("disabled");
            _itemStorage.unshift(r);
        });
        return false;
    },
    view: function view() {
        return m(
            "div",
            null,
            m(
                "form",
                { method: "POST", action: "/api/messages", onsubmit: this.onSubmit },
                m("input", { type: "text", "class": "message u-full-width", placeholder: "your message" }),
                m("br", null),
                m("input", { type: "submit", value: "weg schicken", "class": "button-primary" })
            )
        );
    }
};

var RemoveLinkComponent = {
    view: function view(vnode) {
        var item = vnode.attrs.item;
        return m(
            "a",
            { href: vnode.attrs.href, onclick: function onclick(evt) {
                    evt.preventDefault();
                    if (!confirm("What, really delete??")) return false;
                    _itemStorage = _itemStorage.filter(function (i) {
                        return i != item;
                    });
                    m.request({
                        url: vnode.attrs.href,
                        method: 'DELETE',
                        data: { ID: item.ID }
                    }).catch(function (e) {
                        alert(e);
                        console.error(e);
                    });
                    return false;
                } },
            "x"
        );
    }
};

var ListItemComponent = {
    view: function view(vnode) {
        var item = vnode.attrs.item;
        return m(
            "tr",
            { key: item.ID },
            m(
                "td",
                null,
                "#",
                item.ID
            ),
            m(
                "td",
                null,
                item.CreatedAt
            ),
            m(
                "td",
                null,
                m.trust(item.HTML)
            ),
            m(
                "td",
                null,
                m(RemoveLinkComponent, { item: item, href: "/api/messages/:ID" })
            )
        );
    }
};

var ListComponent = function () {
    function ListComponent() {
        _classCallCheck(this, ListComponent);
    }

    _createClass(ListComponent, [{
        key: "loadData",
        value: function loadData(url) {
            var rargs = {
                url: url,
                method: 'GET'
            };
            m.request(rargs).then(function (r) {
                _itemStorage = r;
            });
        }
    }, {
        key: "oncreate",
        value: function oncreate(vnode) {
            this.loadData(vnode.attrs.url);
        }
    }, {
        key: "view",
        value: function view(vnode) {
            return m(
                "div",
                { className: "container" },
                m(
                    "table",
                    null,
                    m(
                        "thead",
                        null,
                        m(
                            "tr",
                            null,
                            m(
                                "th",
                                null,
                                "ID"
                            ),
                            m(
                                "th",
                                null,
                                "Created"
                            ),
                            m(
                                "th",
                                null,
                                "Message"
                            ),
                            m(
                                "th",
                                null,
                                "Action"
                            )
                        )
                    ),
                    m(
                        "tbody",
                        null,
                        _itemStorage.map(function (i) {
                            return m(ListItemComponent, { item: i });
                        }),
                        m(
                            "tr",
                            null,
                            m(
                                "td",
                                { colspan: "4" },
                                m(FormComp, null)
                            )
                        )
                    )
                )
            );
        }
    }]);

    return ListComponent;
}();

;

var RootComp = function () {
    function RootComp() {
        _classCallCheck(this, RootComp);
    }

    _createClass(RootComp, [{
        key: "view",
        value: function view(vnode) {
            return m(ListComponent, { url: "/api/messages" });
        }
    }]);

    return RootComp;
}();

m.mount(document.getElementById("app"), RootComp);
//# sourceMappingURL=gomsglog.js.map