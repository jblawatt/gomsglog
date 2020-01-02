
import * as m from "mithril";

// console.log(Vue)

interface IMessage {
    ID: string;
    Message: string;
    CreatedAt: string;
    HTML: string;
}



let AllStorage = {
    data: [],
    load: () => {
        m.request<IMessage[]>('/api/messages').then(data => {
            data.forEach(message => {
                AllStorage.data.push(message);
            })
        });
    }
}

window["AllStorage"] = AllStorage;

interface IStorage {
    data: IMessage[]
}

interface ICardComponentAttrs {
    message: string;
    storage: IStorage;
    isPrimary: boolean;
}

interface ICardItemComponentAttrs {
    message: IMessage;
}

class CardItemComponent implements m.ClassComponent<ICardItemComponentAttrs>{
    view(vnode: m.Vnode<ICardItemComponentAttrs>): m.Vnode {
        return <li class="list-group-item" data-item-id={vnode.attrs.message.ID}>
            <div class="card-text">
                <small class="text-muted">{vnode.attrs.message.CreatedAt}</small>
                <p>
                    {m.trust(vnode.attrs.message.HTML)}
                </p>
                <div class="text-center">
                    <div class="btn-group btn-group-sm">
                        <button class="btn btn-light btn-sm">
                            <i class="fa fa-trash"></i> delete
                                                </button>
                        <button class="btn btn-light">
                            <i class="fa fa-retweet"></i> repost
                                                </button>
                    </div>
                </div>
            </div>
        </li>;
    }
}

class CardComponent implements m.ClassComponent<ICardComponentAttrs> {
    view(vnode: m.Vnode<ICardComponentAttrs>): m.Vnode {
        let classNames = ['max-w-sm', 'rounded', 'overflow-hidden', 'shadow-lg'];
        
        return <div class={classNames.join(" ")}>
            <p className="text-gray-700 text-base">
                {vnode.attrs.message}
            </p>
            <ul className="list-group list-group-flush">
                {vnode.attrs.storage.data.map(item => {
                    return <CardItemComponent message={item} />
                })}
            </ul>
        </div>;
    }
}



m.mount(document.getElementById('latest-posts'), {
    view() {
        return <CardComponent message="Latest Posts" storage={AllStorage} isPrimary={true} />
    }
});


// m.mount(document.getElementById('today-posts'), {
//     view() {
//         return <CardComponent message="Today Posts" storage={AllStorage} />
//     }
// });


// m.mount(document.getElementById('project-posts'), {
//     view() {
//         return <CardComponent message="Project Posts" storage={AllStorage} />
//     }
// });

// m.mount(document.getElementById('todo-posts'), {
//     view() {
//         return <CardComponent message="Todo Posts" storage={AllStorage} />
//     }
// });

// m.mount(document.getElementById('another-posts'), {
//     view() {
//         return <CardComponent message="Another Posts" storage={AllStorage} />
//     }
// });



document.addEventListener("DOMContentLoaded", function () {
    AllStorage.load();
})