/**
 * @global m object
 */

const store = { items: [] };

m.request({ url: '/api/messages' }).then((r) => store.items = r);


const MessageCardComp = {
    view(vnode) {
        let message = vnode.attrs.message;
        let children = [
            m('div', { class: 'card-body' }, m('p', { class: 'card-text' }, message.Original)),
            m('div', { class: 'card-footer text-muted' }, message.CreatedAt)
        ]
        return m('div', { class: 'card' }, children);
    }
}


const RootListComp = {
    view() {
        let children = store.items.map(message => m(MessageCardComp, { message: message }))
        return m('div', { class: 'row' }, m('div', { class: 'col-xs-12' }, children));
    }
}




const root = document.getElementById('items-section-container');
m.mount(root, RootListComp);