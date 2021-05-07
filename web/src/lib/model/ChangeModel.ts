import {action, makeObservable, observable} from 'mobx';

type Entries<T> = {
    [K in keyof T]: [K, T[K]]
}[keyof T][];

function filterNulls<T>(value: T): T {
    const entries = Object.entries(value) as Entries<T>;
    return Object.fromEntries(entries.filter(([, value]) => value !== null)) as unknown as T;
}

export default class ChangeModel<T> {
    @observable protected readonly _originalValues: T;
    @observable protected readonly _values: Partial<T> = {};

    constructor(data: T) {
        makeObservable(this);
        this._originalValues = filterNulls(data);
    }

    get changes() {
        return this._values;
    }

    get<K extends keyof T>(key: K): T[K] {
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        return key in this._values ? this._values[key]! : this._originalValues[key];
    }

    @action
    set<K extends keyof T>(key: K, value: Partial<T>[K]) {
        if (this._originalValues[key] === value) delete this._values[key];
        else this._values[key] = value;
    }

    @action
    undo(key: keyof T) {
        delete this._values[key];
    }

    get isChanged() {
        return Object.keys(this._values).length > 0;
    }

    @action
    commit() {
        const entries = Object.entries(this._values) as Entries<T>;
        entries.forEach(([key, value]) => {
            this._originalValues[key] = value;
            delete this._values[key];
        });
    }

    @action
    revert() {
        const keys = Object.keys(this._values) as (keyof T)[];
        keys.forEach((key) => delete this._values[key]);
    }
}
