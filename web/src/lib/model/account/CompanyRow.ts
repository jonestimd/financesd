import {action, makeObservable, observable} from 'mobx';
import {IRow} from '../../component/table/Table';
import {CompanyModel} from './CompanyModel';

export default class CompanyRow {
    @observable private _name?: string;

    constructor(private data: CompanyModel | IRow & Partial<Pick<CompanyModel, 'name' | 'accounts' | 'version'>>) {
        makeObservable(this);
    }

    get id() {
        return this.data.id;
    }

    get version() {
        return this.data.version ?? 0;
    }

    get name() {
        return this._name ?? this.data.name ?? '';
    }

    @action
    setName(name: string) {
        name = name.trim();
        this._name = name === this.data.name ? undefined : name;
    }

    get accounts() {
        return this.data.accounts?.length ?? 0;
    }

    get isChanged() {
        return !('version' in this.data) || !!this._name;
    }

    get isValid() {
        return !!this.name;
    }
}
