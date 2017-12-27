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
                            <th>ID</th><th>Created</th><th>Message</th><th>Action</th>
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
