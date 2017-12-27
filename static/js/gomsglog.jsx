
let _itemStorage = [];

let FormComp = {
    onSubmit(evt) {
        evt.preventDefault();
        let input = evt.currentTarget.querySelector('.message');
        input.setAttribute("disabled", true);
        let message = input.value;
        let method = evt.currentTarget.getAttribute("method");
        let action = evt.currentTarget.getAttribute("action");
        m.request({
            url: action,
            method,
            headers: {
                "Content-Type": "application/json"
            },
            data: {
                message
            }
        }).then(r => {
            input.value = "";
            input.removeAttribute("disabled");
            _itemStorage.unshift(r);
        });
        return false;
    },
    view() {
        return (
            <div >
                <form method="POST" action="/api/messages" onsubmit={this.onSubmit}>
                    <input type="text" class="message u-full-width" placeholder="your message"></input>
                    <br />
                    <input type="submit" value="weg schicken" class="button-primary"></input>
                </form>
            </div>
        );
    }
};


let RemoveLinkComponent = {
    view(vnode) {
        let item = vnode.attrs.item;
        return <a href={vnode.attrs.href} onclick={(evt) => {
            evt.preventDefault();
            if (!confirm("What, really delete??")) return false;
            _itemStorage = _itemStorage.filter(i => i != item);
            m.request({
                url: vnode.attrs.href,
                method: 'DELETE',
                data: { ID: item.ID },
            }).catch(e => {
                alert(e);
                console.error(e);
            });
            return false;
        }}>x</a>
    }
}


let ListItemComponent = {
    view(vnode) {
        let item = vnode.attrs.item;
        return (
            <tr key={item.ID}>
                <td>#{item.ID}</td>
                <td>{item.CreatedAt}</td>
                <td>{m.trust(item.HTML)}</td>
                <td>
                    <RemoveLinkComponent item={item} href="/api/messages/:ID" />
                </td>
            </tr>
        );
    }
}

class ListComponent {
    loadData(url) {
        let rargs = {
            url: url,
            method: 'GET',
        };
        m.request(rargs).then(r => {
            _itemStorage = r;
        });
    }
    oncreate(vnode) {
        this.loadData(vnode.attrs.url);
    }
    view(vnode) {
        return (
            <div className="container">
                <table>
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>Created</th>
                            <th>Message</th>
                            <th>Action</th>
                        </tr>
                    </thead>
                    <tbody>
                        {
                            _itemStorage.map(i => {
                                return <ListItemComponent item={i} />
                            })
                        }
                        <tr>
                            <td colspan="4">
                                <FormComp />
                            </td>

                        </tr>
                    </tbody>
                </table>
            </div >
        );
    }
};

class RootComp {
    view(vnode) {
        return <ListComponent url="/api/messages" />;
    }
}

m.mount(document.getElementById("app"), RootComp);
