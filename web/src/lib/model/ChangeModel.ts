import {action, makeObservable, observable, ObservableMap} from 'mobx';
import {Nullable} from './entityUtils';

type Entry<T> = {
    [K in keyof T]: [K, T[K]]
}[keyof T];

interface MapOf<T> extends Omit<ObservableMap<keyof T, unknown>, 'get' | 'set' | 'entries'> {
    get<K extends keyof T>(key: K): T[K];
    set<K extends keyof T>(key: K, value: T[K]): void;
    entries(): IterableIterator<Entry<T>>;
}

function filterNulls<T>(value: T): T {
    const entries = Object.entries(value) as Entry<T>[];
    return Object.fromEntries(entries.filter(([, value]) => value !== null)) as unknown as T;
}

export default class ChangeModel<T> {
    protected readonly _originalValues: T;
    protected readonly _values = observable.map() as MapOf<Partial<T>>;

    constructor(data: T) {
        makeObservable(this);
        this._originalValues = filterNulls(data);
    }

    get changes() {
        return Array.from(this._values.entries()).reduce<Nullable<Partial<T>>>((changes, [key, value]) => ({
            ...changes,
            [key]: value ?? null,
        }), {});
    }

    get<K extends keyof T>(key: K): T[K] { // TODO should return | undefined
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        return this._values.has(key) ? this._values.get(key)! : this._originalValues[key];
    }

    @action
    set<K extends keyof T>(key: K, value: Partial<T>[K]) {
        if (this._originalValues[key] === value) this._values.delete(key);
        else this._values.set(key, value);
    }

    @action
    undo(key: keyof T) {
        this._values.delete(key);
    }

    get isChanged() {
        return this._values.size > 0;
    }

    @action
    commit() {
        const entries = Array.from(this._values.entries());
        entries.forEach(([key, value]) => {
            // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
            this._originalValues[key] = value!;
            this._values.delete(key);
        });
    }

    @action
    revert() {
        const keys = Array.from(this._values.keys());
        keys.forEach((key) => this._values.delete(key));
    }
}
