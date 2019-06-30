export default class Meta {
    constructor (data = {}) {
        Object.keys(data).forEach(key => {
            this[key] = data[key]
        })

        this.menu = !data.menu ? false : Object.assign({
            id: null,
            i18n: null,
            path: null,
            order: 1,
            parent: null,
            adminView: true,
            businessView: true
        }, data.menu)

        this.auth = Object.assign({
            view: null,
            operation: []
        }, data.auth)

        this.authScope = 'global'
    }
}
